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

// SendNftyPush forwards alert messages to Nfty server
func SendNftyPush(event models.Event, snapshot io.Reader) {
	// Build notification
	var message string
	if config.ConfigData.Alerts.Nfty.Template != "" {
		message = renderMessage(config.ConfigData.Alerts.Nfty.Template, event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "Nfty").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("plaintext", event)
	}

	NftyURL := fmt.Sprintf("%s/%s", config.ConfigData.Alerts.Nfty.Server, config.ConfigData.Alerts.Nfty.Topic)

	// Set headers
	var headers []map[string]string
	headers = append(headers, map[string]string{"Content-Type": "text/markdown"})
	headers = append(headers, map[string]string{"X-Title": config.ConfigData.Alerts.General.Title})
	headers = append(headers, config.ConfigData.Alerts.Nfty.Headers...)

	// Set action link to the recorded clip
	var clip string
	if config.ConfigData.Frigate.PublicURL != "" {
		clip = fmt.Sprintf("%s/api/events/%s/clip.mp4", config.ConfigData.Frigate.PublicURL, event.ID)
	} else {
		clip = fmt.Sprintf("%s/api/events/%s/clip.mp4", config.ConfigData.Frigate.Server, event.ID)
	}

	headers = append(headers, map[string]string{"X-Actions": "view, View Clip, " + clip + ", clear=true"})

	var attachment []byte
	if event.HasSnapshot {
		headers = append(headers, map[string]string{"X-Filename": "snapshot.jpg"})
		attachment, _ = io.ReadAll(snapshot)
	} else {
		message += "\n\nNo snapshot available."
	}

	// Escape newlines in message
	message = strings.ReplaceAll(message, "\n", "\\n")
	headers = append(headers, map[string]string{"X-Message": message})

	resp, err := util.HTTPPost(NftyURL, config.ConfigData.Alerts.Nfty.Insecure, attachment, headers...)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Nfty").
			Err(err).
			Msg("Unable to send alert")
		return
	}

	// Nfty returns HTTP 200 even if there is an error, so we need to inspect returned body
	if strings.Contains(string(resp), "error") {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Nfty").
			Str("error", string(resp)).
			Msg("Unable to send alert")
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Nfty").
		Msg("Alert sent")
}
