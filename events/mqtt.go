package frigate

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqtt_topic string

// SubscribeMQTT establishes subscription to MQTT server & listens for messages
func SubscribeMQTT() {
	mqtt_topic = fmt.Sprintf("%s/%s", config.ConfigData.Frigate.MQTT.TopicPrefix, strings.ToLower(config.ConfigData.Frigate.MQTT.Mode))
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
		if retry >= 3 {
			log.Fatal().Msgf("Max retries exceeded. Failed to establish MQTT session to %s", config.ConfigData.Frigate.MQTT.Server)
		}
		// Connect to MQTT broker
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			retry += 1
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
	log.Error().
		Err(err).
		Msg("Lost connection to MQTT broker")
}

// connectHandler logs message on MQTT connection
func connectHandler(client mqtt.Client) {
	log.Info().Msg("Connected to MQTT.")
	if strings.ToLower(config.ConfigData.Frigate.MQTT.Mode) == "events" {
		if subscription := client.Subscribe(mqtt_topic, 0, processEvent); subscription.Wait() && subscription.Error() != nil {
			log.Error().Msgf("Failed to subscribe to topic: %s", mqtt_topic)
			time.Sleep(10 * time.Second)
		}
	}
	if strings.ToLower(config.ConfigData.Frigate.MQTT.Mode) == "reviews" {
		if subscription := client.Subscribe(mqtt_topic, 0, processReview); subscription.Wait() && subscription.Error() != nil {
			log.Error().Msgf("Failed to subscribe to topic: %s", mqtt_topic)
			time.Sleep(10 * time.Second)
		}
	}
	log.Info().Msgf("Subscribed to MQTT topic: %s", mqtt_topic)
}
