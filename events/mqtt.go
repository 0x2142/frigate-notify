package frigate

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/exp/slices"
)

// SubscribeMQTT establishes subscription to MQTT server & listens for messages
func SubscribeMQTT() {
	// MQTT client configuration
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.ConfigData.Frigate.MQTT.Server, config.ConfigData.Frigate.MQTT.Port))
	opts.SetClientID(config.ConfigData.Frigate.MQTT.ClientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(connectHandler)
	if config.ConfigData.Frigate.MQTT.Username != "" && config.ConfigData.Frigate.MQTT.Password != "" {
		opts.SetUsername(config.ConfigData.Frigate.MQTT.Username)
		opts.SetPassword(config.ConfigData.Frigate.MQTT.Password)
	}

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

// processEvent handles incoming MQTT messages & pulls out relevant info for alerting
func processEvent(client mqtt.Client, msg mqtt.Message) {
	// Parse incoming MQTT message
	var event models.MQTTEvent
	json.Unmarshal(msg.Payload(), &event)

	if event.Type == "new" || event.Type == "update" {
		if event.Type == "new" {
			log.Info().
				Str("event_id", event.After.ID).
				Msg("New event received")
		} else if event.Type == "update" {
			log.Info().
				Str("event_id", event.After.ID).
				Msg("Event update received")
		}
		// Skip excluded cameras
		if slices.Contains(config.ConfigData.Frigate.Cameras.Exclude, event.After.Camera) {
			log.Info().
				Str("event_id", event.After.ID).
				Str("camera", event.After.Camera).
				Msg("Event dropped - Camera Excluded")
			return
		}

		// Convert to human-readable timestamp
		eventTime := time.Unix(int64(event.After.StartTime), 0)
		log.Info().
			Str("event_id", event.After.ID).
			Str("camera", event.After.Camera).
			Str("label", event.After.Label).
			Str("zones", strings.Join(event.After.CurrentZones, ",")).
			Msg("Event Detected")
		log.Debug().
			Str("event_id", event.After.ID).
			Msgf("Event start time: %s", eventTime)

		// Check that event passes configured filters
		if !checkEventFilters(event.After.Event) {
			return
		}

		// Check if already notified on zones
		if zoneAlreadyAlerted(event.After.Event) {
			log.Info().
				Str("event_id", event.After.ID).
				Str("camera", event.After.Camera).
				Str("label", event.After.Label).
				Str("zones", strings.Join(event.After.CurrentZones, ",")).
				Msg("Event dropped - Already notified on this zone")
			return
		} else {
			log.Debug().
				Str("event_id", event.After.ID).
				Str("camera", event.After.Camera).
				Str("label", event.After.Label).
				Str("zones", strings.Join(event.After.CurrentZones, ",")).
				Msg("Object entered new zone")
		}

		// If snapshot was collected, pull down image to send with alert
		var snapshot io.Reader
		var snapshotURL string
		if event.After.HasSnapshot {
			snapshotURL = config.ConfigData.Frigate.Server + eventsURI + "/" + event.After.ID + snapshotURI
			snapshot = GetSnapshot(snapshotURL, event.After.ID)
		}

		// Send alert with snapshot
		notifier.SendAlert(event.After.Event, snapshotURL, snapshot, event.After.ID)
	}

	// Clear event cache entry when event ends
	if event.Type == "end" {
		log.Debug().
			Str("event_id", event.After.ID).
			Msg("Event ended")
		delZoneAlerted(event.After.Event)
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
	topic := fmt.Sprintf(config.ConfigData.Frigate.MQTT.TopicPrefix + "/events")
	if subscription := client.Subscribe(topic, 0, processEvent); subscription.Wait() && subscription.Error() != nil {
		log.Error().Msgf("Failed to subscribe to topic: %s", topic)
		time.Sleep(10 * time.Second)
	}
	log.Info().Msgf("Subscribed to MQTT topic: %s", topic)
}
