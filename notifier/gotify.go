package notifier

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
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
func SendGotifyPush(event models.Event, provider notifMeta) {
	profile := config.ConfigData.Alerts.Gotify[provider.index]
	status := &config.Internal.Status.Notifications.Gotify[provider.index]

	var snapshotURL string
	if config.ConfigData.Frigate.PublicURL != "" {
		snapshotURL = config.ConfigData.Frigate.PublicURL + "/api/events/" + event.ID + "/snapshot.jpg"
	} else {
		snapshotURL = config.ConfigData.Frigate.Server + "/api/events/" + event.ID + "/snapshot.jpg"
	}
	// Build notification
	var message string
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Gotify")
	} else {
		message = renderMessage("markdown", event, "message", "Gotify")
	}

	if event.HasSnapshot {
		message += fmt.Sprintf("\n\n![](%s)", snapshotURL)
	}
	var title string
	if profile.Title != "" {
		title = renderMessage(profile.Title, event, "title", "gotify")
	} else {
		title = renderMessage(config.ConfigData.Alerts.General.Title, event, "title", "gotify")
	}
	payload := gotifyPayload{
		Message:  message,
		Title:    title,
		Priority: config.ConfigData.Alerts.Gotify[provider.index].Priority,
	}
	payload.Extras.ClientDisplay.ContentType = "text/markdown"
	payload.Extras.ClientNotification.BigImageURL = snapshotURL

	data, err := json.Marshal(payload)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Gotify").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	gotifyURL := fmt.Sprintf("%s/message?token=%s&", profile.Server, profile.Token)

	header := map[string]string{"Content-Type": "application/json"}
	response, err := util.HTTPPost(gotifyURL, profile.Insecure, data, "", header)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Gotify").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	// Check for errors:
	if strings.Contains(string(response), "error") {
		var errorMessage gotifyError
		json.Unmarshal(response, &errorMessage)
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Gotify").
			Int("provider_id", provider.index).
			Msgf("Unable to send alert: %v - %v", errorMessage.Error, errorMessage.ErrorDescription)
		status.NotifFailure(errorMessage.ErrorDescription)
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Gotify").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
