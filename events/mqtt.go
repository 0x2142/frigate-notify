package frigate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/0x2142/frigate-notify/notifier"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MQTTServer string
var MQTTPort int
var MQTTUser string
var MQTTPass string

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
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", MQTTServer, MQTTPort))
	opts.SetClientID("frigate-notify")
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetOnConnectHandler(connectHandler)
	if MQTTUser != "" && MQTTPass != "" {
		opts.SetUsername(MQTTUser)
		opts.SetPassword(MQTTPass)
	}

	var subscribed = false
	var retry = 0
	for !subscribed {
		if retry >= 3 {
			log.Fatalf("ERROR: Max retries exceeded. Failed to establish MQTT session to %s", MQTTServer)
		}
		// Connect to MQTT broker
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			retry += 1
			log.Printf("Could not connect to MQTT at %v: %v", MQTTServer, token.Error())
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

	if event.Type == "new" {
		// Convert to human-readable timestamp
		eventTime := time.Unix(int64(event.After.StartTime), 0)

		log.Printf("Event ID %v detected %v in zone(s): %v", event.After.ID, event.After.Label, event.After.Zones)
		log.Println("Event Start time: ", eventTime)

		// If snapshot was collected, pull down image to send with alert
		var snapshot io.Reader
		var snapshotURL string
		if event.After.HasSnapshot {
			snapshotURL = FrigateServerURL + eventsURI + "/" + event.After.ID + snapshotURI
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
