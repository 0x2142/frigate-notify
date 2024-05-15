package notifier

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

// SendNftyPush forwards alert messages to Nfty server
func SendNftyPush(event models.Event, snapshot io.Reader, eventid string) {
	// Build notification
	message := renderMessage("plaintext", event)

	NftyURL := fmt.Sprintf("%s/%s", config.ConfigData.Alerts.Nfty.Server, config.ConfigData.Alerts.Nfty.Topic)

	// Set headers
	var headers []map[string]string
	headers = append(headers, map[string]string{"Content-Type": "text/markdown"})
	headers = append(headers, map[string]string{"X-Title": config.ConfigData.Alerts.General.Title})

	// Set action link to the recorded clip
	clip := fmt.Sprintf("%s/api/events/%s/clip.mp4", config.ConfigData.Frigate.Server, eventid)
	headers = append(headers, map[string]string{"X-Actions": "view, View Clip, " + clip + ", clear=true"})

	var attachment []byte
	if snapshot != nil {
		headers = append(headers, map[string]string{"X-Filename": "snapshot.jpg"})
		attachment, _ = io.ReadAll(snapshot)
	} else {
		message += "\n\nNo snapshot saved."
	}

	// Escape newlines in message
	message = strings.ReplaceAll(message, "\n", "\\n")
	headers = append(headers, map[string]string{"X-Message": message})

	_, err := util.HTTPPost(NftyURL, config.ConfigData.Alerts.Nfty.Insecure, attachment, headers...)
	if err != nil {
		log.Print("Failed to send Nfty notification: ", err)
		return
	}

	log.Printf("Event ID %v - Nfty alert sent", eventid)
}
