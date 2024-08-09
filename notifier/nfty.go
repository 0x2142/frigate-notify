package notifier

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

// SendNtfyPush forwards alert messages to Ntfy server
func SendNtfyPush(event models.Event, snapshot io.Reader) {
	// Build notification
	var message string
	if config.ConfigData.Alerts.Ntfy.Template != "" {
		message = renderMessage(config.ConfigData.Alerts.Ntfy.Template, event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "Ntfy").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("plaintext", event)
	}

	NtfyURL := fmt.Sprintf("%s/%s", config.ConfigData.Alerts.Ntfy.Server, config.ConfigData.Alerts.Ntfy.Topic)

	// Set headers
	var headers []map[string]string
	headers = append(headers, map[string]string{"Content-Type": "text/markdown"})
	headers = append(headers, map[string]string{"X-Title": config.ConfigData.Alerts.General.Title})
	headers = append(headers, config.ConfigData.Alerts.Ntfy.Headers...)

	// Set action link to the recorded clip
	var clip string
	if config.ConfigData.Frigate.PublicURL != "" {
		clip = fmt.Sprintf("%s/api/events/%s/clip.mp4", config.ConfigData.Frigate.PublicURL, event.ID)
	} else {
		clip = fmt.Sprintf("%s/api/events/%s/clip.mp4", config.ConfigData.Frigate.Server, event.ID)
	}

	var attachment []byte
	if event.HasSnapshot {
		headers = append(headers, map[string]string{"X-Filename": "snapshot.jpg"})
		attachment, _ = io.ReadAll(snapshot)
	}

	// Escape newlines in message
	message = strings.ReplaceAll(message, "\n", "\\n")
	headers = append(headers, map[string]string{"X-Message": message})

	// Check if custom action has been added, otherwise include default
	var hasAction bool
	for _, header := range headers {
		if _, ok := header["X-Actions"]; ok {
			hasAction = true
			break
		}
	}
	if !hasAction {
		headers = append(headers, map[string]string{"X-Actions": "view, View Clip, " + clip + ", clear=true"})
	}

	headers = renderHeaders(headers, event)

	resp, err := util.HTTPPost(NtfyURL, config.ConfigData.Alerts.Ntfy.Insecure, attachment, headers...)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Ntfy").
			Err(err).
			Msg("Unable to send alert")
		return
	}

	// Ntfy returns HTTP 200 even if there is an error, so we need to inspect returned body
	if strings.Contains(string(resp), "error") {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Ntfy").
			Str("error", string(resp)).
			Msg("Unable to send alert")
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Ntfy").
		Msg("Alert sent")
}
