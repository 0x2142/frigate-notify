package notifier

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

type SignalPayload struct {
	Number      string   `json:"number"`
	Recipients  []string `json:"recipients,omitempty"`
	Message     string   `json:"message"`
	Attachments []string `json:"base64_attachments"`
}

// SendSignalMessage pushes alert message to Signal via webhook
func SendSignalMessage(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Signal[provider.index]
	status := &config.Internal.Status.Notifications.Signal[provider.index]

	var err error
	var message string
	// Build notification
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Signal")
	} else {
		message = renderMessage("plaintext", event, "message", "Signal")
	}

	// Process recipients list & prefix + to numbers that don't include it
	var recipients []string
	r, _ := regexp.Compile(`^\d+$`)
	for _, recipient := range profile.Recipients {
		if r.Match([]byte(recipient)) {
			recipients = append(recipients, "+"+recipient)
		} else {
			recipients = append(recipients, recipient)
		}
	}

	// Check if sending account has + prefix if needed
	if r.Match([]byte(profile.Account)) {
		profile.Account = "+" + profile.Account
	}

	// Build payload
	payload := SignalPayload{Message: message, Number: profile.Account, Recipients: recipients}
	if event.HasSnapshot {
		img, _ := io.ReadAll(snapshot)
		attach := base64.StdEncoding.EncodeToString(img)
		payload.Attachments = append(payload.Attachments, attach)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Signal").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	url := profile.Server + "/v2/send"
	_, err = util.HTTPPost(url, profile.Insecure, []byte(data), "")

	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Signal").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Signal").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
