package notifier

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
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
func SendGotifyPush(message, snapshotURL string) {
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
		log.Println("Unable to build Gotify payload: ", err)
		return
	}

	response, err := HTTPPost(data)
	if err != nil {
		log.Print("Failed to send Gotify notification: ", err)
		return
	}
	// Check for errors:
	if strings.Contains(string(response), "error") {
		var errorMessage gotifyError
		json.Unmarshal(response, &errorMessage)
		log.Printf("Failed to send Gotify notification: %s - %s", errorMessage.Error, errorMessage.ErrorDescription)
		return
	}
	log.Print("Gotify message sent")
}

// HTTPPost performs an HTTP POST to the target URL
// and includes auth parameters, ignoring certificates, etc
func HTTPPost(payload []byte) ([]byte, error) {
	gotifyURL := fmt.Sprintf("%s/message?token=%s&", config.ConfigData.Alerts.Gotify.Server, config.ConfigData.Alerts.Gotify.Token)

	// New HTTP Client
	client := http.Client{Timeout: 10 * time.Second}

	// Ignore SSL verification if set
	if config.ConfigData.Alerts.Gotify.Insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	// Setup new HTTP Request
	req, err := http.NewRequest(http.MethodPost, gotifyURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP POST
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
