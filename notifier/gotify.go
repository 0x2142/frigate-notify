package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/util"
)

// gotifyError defines structure of Gotify error messages
type gotifyError struct {
	Error            string `json:"error"`
	ErrorCode        int    `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
}

// gotifyPayload defines structure of Gotify push messages
type gotifyPayload struct {
	Message  string `json:"message"`
	Title    string `json:"title,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Extras   struct {
		ClientDisplay struct {
			ContentType string `json:"contentType,omitempty"`
		} `json:"client::display"`
		ClientNotification struct {
			BigImageURL string `json:"bigImageUrl,omitempty"`
		} `json:"client::notification"`
	} `json:"extras,omitempty"`
}

// SendGotifyPush forwards alert messages to Gotify push notification server
func SendGotifyPush(message, snapshotURL string, eventid string) {
	if snapshotURL != "" {
		message += fmt.Sprintf("\n\n![](%s)", snapshotURL)
	} else {
		message += "\n\nNo snapshot saved."
	}
	payload := gotifyPayload{
		Message:  message,
		Title:    config.ConfigData.Alerts.General.Title,
		Priority: 5,
	}
	payload.Extras.ClientDisplay.ContentType = "text/markdown"
	payload.Extras.ClientNotification.BigImageURL = snapshotURL

	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Event ID %v - Unable to build Gotify payload: ", eventid, err)
		return
	}

	gotifyURL := fmt.Sprintf("%s/message?token=%s&", config.ConfigData.Alerts.Gotify.Server, config.ConfigData.Alerts.Gotify.Token)

	response, err := util.HTTPPost(gotifyURL, config.ConfigData.Alerts.Gotify.Insecure, data)
	if err != nil {
		log.Print("Failed to send Gotify notification: ", err)
		return
	}
	// Check for errors:
	if strings.Contains(string(response), "error") {
		var errorMessage gotifyError
		json.Unmarshal(response, &errorMessage)
		log.Printf("Event ID %v - Failed to send Gotify notification: %s - %s", eventid, errorMessage.Error, errorMessage.ErrorDescription)
		return
	}
	log.Printf("Event ID %v - Gotify alert sent", eventid)
}
