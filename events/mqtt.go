package frigate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/notifier"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/exp/slices"
)

// MQTTEvent stores incoming MQTT payloads from Frigate
type MQTTEvent struct {
	Before struct {
		Event
	} `json:"before,omitempty"`
	After struct {
		Event
	} `json:"after,omitempty"`
	Type string `json:"type"`
}

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
			log.Fatalf("ERROR: Max retries exceeded. Failed to establish MQTT session to %s", config.ConfigData.Frigate.MQTT.Server)
		}
		// Connect to MQTT broker
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			retry += 1
			log.Printf("Could not connect to MQTT at %v: %v", config.ConfigData.Frigate.MQTT.Server, token.Error())
			log.Printf("Retrying in 10 seconds. Attempt %v of 3.", retry)
			time.Sleep(10 * time.Second)
			continue
		}
		return
	}
}

// processEvent handles incoming MQTT messages & pulls out relevant info for alerting
func processEvent(client mqtt.Client, msg mqtt.Message) {
	// Parse incoming MQTT message
	var event MQTTEvent
	json.Unmarshal(msg.Payload(), &event)

	if event.Type == "new" || event.Type == "update" {
		// Skip excluded cameras
		if slices.Contains(config.ConfigData.Frigate.Cameras.Exclude, event.After.Camera) {
			log.Printf("Skipping event from excluded camera: %v", event.After.Camera)
			return
		}

		// Convert to human-readable timestamp
		eventTime := time.Unix(int64(event.After.StartTime), 0)

		log.Printf("Event ID %v - Camera %v detected %v in zone(s): %v", event.After.ID, event.After.Camera, event.After.Label, event.After.CurrentZones)
		log.Printf("Event ID %v - Start time: %s", event.After.ID, eventTime)

		// Check that event passes the zone filter
		if !isAllowedZone(event.After.ID, event.After.CurrentZones) {
			return
		}

		// Skip update events where zone didn't change
		// Compares current detected zone to previous list of zones entered
		zoneChanged := false
		for _, zone := range event.After.CurrentZones {
			if !slices.Contains(event.Before.EnteredZones, zone) {
				zoneChanged = true
				log.Printf("Event ID %v - Entered new zone: %s", event.After.ID, zone)
			}
		}
		if event.Type == "update" && !zoneChanged {
			log.Printf("Event ID %v - Zone already alerted, skipping...", event.After.ID)
			return
		}

		// If snapshot was collected, pull down image to send with alert
		var snapshot io.Reader
		var snapshotURL string
		if event.After.HasSnapshot {
			snapshotURL = config.ConfigData.Frigate.Server + eventsURI + "/" + event.After.ID + snapshotURI
			snapshot = GetSnapshot(snapshotURL, event.After.ID)
		}

		message := buildMessage(eventTime, event.After.Event)

		// Send alert with snapshot
		notifier.SendAlert(message, snapshotURL, snapshot)
	}
}

// connectionLostHandler logs error message on MQTT connection loss
func connectionLostHandler(c mqtt.Client, err error) {
	log.Println("Lost connection to MQTT broker. Error: ", err)
}

// connectHandler logs message on MQTT connection
func connectHandler(client mqtt.Client) {
	log.Println("Connected to MQTT.")
	if subscription := client.Subscribe("frigate/events", 0, processEvent); subscription.Wait() && subscription.Error() != nil {
		log.Printf("Failed to subscribe to topic frigate/events")
		time.Sleep(10 * time.Second)
	}
	log.Printf("Subscribed to MQTT topic frigate/events")
}
