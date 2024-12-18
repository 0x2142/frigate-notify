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
func SendNtfyPush(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Ntfy[provider.index]

	// Build notification
	var message string
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Ntfy")
	} else {
		message = renderMessage("plaintext", event, "message", "Ntfy")
	}

	NtfyURL := fmt.Sprintf("%s/%s", profile.Server, profile.Topic)

	// Set headers
	title := renderMessage(config.ConfigData.Alerts.General.Title, event, "title", "Ntfy")
	var headers []map[string]string
	headers = append(headers, map[string]string{"Content-Type": "text/markdown"})
	headers = append(headers, map[string]string{"X-Title": title})
	headers = append(headers, profile.Headers...)

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
		if event.Extra.ReviewLink != "" {
			headers = append(headers, map[string]string{"X-Actions": "view, Review Event, " + event.Extra.ReviewLink + ", clear=true"})
		} else {
			headers = append(headers, map[string]string{"X-Actions": "view, View Clip, " + event.Extra.EventLink + ", clear=true"})
		}
	}

	headers = renderHTTPKV(headers, event, "headers", "Ntfy")

	resp, err := util.HTTPPost(NtfyURL, profile.Insecure, attachment, "", headers...)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Ntfy").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Ntfy[0].NotifFailure(err.Error())
		return
	}

	// Ntfy returns HTTP 200 even if there is an error, so we need to inspect returned body
	if strings.Contains(string(resp), "error") {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Ntfy").
			Int("provider_id", provider.index).
			Str("error", string(resp)).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Ntfy[0].NotifFailure(string(resp))

	}

	log.Info().
		Str("event_id", event.ID).
		Int("provider_id", provider.index).
		Str("provider", "Ntfy").
		Msg("Alert sent")
	config.Internal.Status.Notifications.Ntfy[0].NotifSuccess()
}
