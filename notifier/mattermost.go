package notifier

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

type MattermostPayload struct {
	Text        string                 `json:"text"`
	Channel     string                 `json:"channel,omitempty"`
	Attachments []MattermostAttachment `json:"attachments"`
	Username    string                 `json:"username"`
	Priority    struct {
		Priority string `json:"priority,omitempty"`
	} `json:"priority,omitempty"`
}

type MattermostAttachment struct {
	ImageURL string `json:"image_url,omitempty"`
}

// SendMatterMost pushes alert message to Mattermost via webhook
func SendMattermost(event models.Event, provider notifMeta) {
	profile := config.ConfigData.Alerts.Mattermost[provider.index]
	status := &config.Internal.Status.Notifications.Mattermost[provider.index]

	var snapshotURL string
	if config.ConfigData.Frigate.PublicURL != "" {
		snapshotURL = config.ConfigData.Frigate.PublicURL + "/api/events/" + event.ID + "/snapshot.jpg"
	} else {
		snapshotURL = config.ConfigData.Frigate.Server + "/api/events/" + event.ID + "/snapshot.jpg"
	}

	var err error
	var message string
	// Build notification
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Mattermost")
	} else {
		message = renderMessage("markdown", event, "message", "Mattermost")
	}

	headers := renderHTTPKV(profile.Headers, event, "headers", "Mattermost")

	payload := MattermostPayload{Text: message, Channel: profile.Channel, Username: profile.Username}
	payload.Priority.Priority = profile.Priority
	attach := MattermostAttachment{ImageURL: snapshotURL}
	payload.Attachments = append(payload.Attachments, attach)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Mattermost").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	_, err = util.HTTPPost(profile.Webhook, profile.Insecure, []byte(data), "", headers...)

	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Mattermost").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Mattermost").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
