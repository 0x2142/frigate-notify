package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqtt_topic string

// SubscribeMQTT establishes subscription to MQTT server & listens for messages
func SubscribeMQTT() {
	config.Internal.Status.Frigate.MQTT = "connecting"
	mqtt_topic = fmt.Sprintf("%s/%s", config.ConfigData.Frigate.MQTT.TopicPrefix, strings.ToLower(config.ConfigData.App.Mode))
	// MQTT client configuration
	mqttServer := fmt.Sprintf("tcp://%s:%d", config.ConfigData.Frigate.MQTT.Server, config.ConfigData.Frigate.MQTT.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttServer)
	opts.SetClientID(config.ConfigData.Frigate.MQTT.ClientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(connectHandler)
	if config.ConfigData.Frigate.MQTT.Username != "" && config.ConfigData.Frigate.MQTT.Password != "" {
		opts.SetUsername(config.ConfigData.Frigate.MQTT.Username)
		opts.SetPassword(config.ConfigData.Frigate.MQTT.Password)
	}

	log.Trace().
		Str("server", mqttServer).
		Str("client_id", config.ConfigData.Frigate.MQTT.ClientID).
		Str("username", config.ConfigData.Frigate.MQTT.Username).
		Str("password", "--secret removed--").
		Str("topic", mqtt_topic).
		Bool("auto_reconnect", true).
		Msg("Init MQTT connection")

	var subscribed = false
	var retry = 0
	for !subscribed {
		config.Internal.Status.Frigate.MQTT = "unreachable"
		if retry >= 3 {
			log.Fatal().Msgf("Max retries exceeded. Failed to establish MQTT session to %s", config.ConfigData.Frigate.MQTT.Server)
		}
		// Connect to MQTT broker
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			retry += 1
			config.Internal.Status.Frigate.MQTT = "unreachable"

			log.Warn().Msgf("Could not connect to MQTT at %v: %v", config.ConfigData.Frigate.MQTT.Server, token.Error())
			log.Warn().Msgf("Retrying in 10 seconds. Attempt %v of 3.", retry)
			time.Sleep(10 * time.Second)
			continue
		}
		return
	}
}

// connectionLostHandler logs error message on MQTT connection loss
func connectionLostHandler(c mqtt.Client, err error) {
	config.Internal.Status.Frigate.MQTT = "unreachable"
	log.Error().
		Err(err).
		Msg("Lost connection to MQTT broker")
}

// connectHandler logs message on MQTT connection
func connectHandler(client mqtt.Client) {
	log.Info().Msg("Connected to MQTT.")
	config.Internal.Status.Frigate.MQTT = "ok"
	if subscription := client.Subscribe(mqtt_topic, 0, handleMQTTMsg); subscription.Wait() && subscription.Error() != nil {
		config.Internal.Status.Frigate.MQTT = "unreachable"
		log.Error().Msgf("Failed to subscribe to topic: %s", mqtt_topic)
		time.Sleep(10 * time.Second)
	}

	log.Info().Msgf("Subscribed to MQTT topic: %s", mqtt_topic)
}

// handleMQTTMsg processes incoming MQTT messages depending on topic
func handleMQTTMsg(client mqtt.Client, msg mqtt.Message) {
	topic := strings.Split(msg.Topic(), "/")[1]

	log.Trace().
		RawJSON("payload", msg.Payload()).
		Msg("New MQTT message received")

	switch topic {
	case "reviews":
		var review models.MQTTReview
		json.Unmarshal(msg.Payload(), &review)

		switch review.Type {
		case "new":
			log.Debug().
				Str("review_id", review.After.ID).
				Msg("New review received")
			processReview(review.After.Review)
		case "update":
			log.Debug().
				Str("review_id", review.After.ID).
				Msg("Review update received")
			processReview(review.After.Review)
		case "end":
			log.Debug().
				Str("review_id", review.After.ID).
				Msg("Review ended")
			for _, detection := range review.After.Data.Detections {
				delZoneAlerted(models.Event{
					ID:           detection,
					Camera:       review.After.Camera,
					CurrentZones: review.After.Data.Zones,
				})
			}
		}
	case "events":
		var event models.MQTTEvent
		json.Unmarshal(msg.Payload(), &event)

		switch event.Type {
		case "new":
			log.Info().
				Str("event_id", event.After.ID).
				Msg("New event received")
			processEvent(event.After.Event)
		case "update":
			log.Info().
				Str("event_id", event.After.ID).
				Msg("Event update received")
			processEvent(event.After.Event)
		case "end":
			log.Debug().
				Str("event_id", event.After.ID).
				Msg("Event ended")
			delZoneAlerted(event.After.Event)
		}

	}
}
