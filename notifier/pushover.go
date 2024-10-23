package notifier

import (
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/gregdel/pushover"
)

// SendPushoverMessage sends alert message through Pushover service
func SendPushoverMessage(event models.Event, snapshot io.Reader) {
	// Build notification
	var message string
	if config.ConfigData.Alerts.Pushover.Template != "" {
		message = renderMessage(config.ConfigData.Alerts.Pushover.Template, event, "message", "Pushover")
	} else {
		message = renderMessage("html", event, "message", "Pushover")
		message = strings.ReplaceAll(message, "<br />", "")
	}

	push := pushover.New(config.ConfigData.Alerts.Pushover.Token)
	recipient := pushover.NewRecipient(config.ConfigData.Alerts.Pushover.Userkey)

	// Create new message
	title := renderMessage(config.ConfigData.Alerts.General.Title, event, "title", "Pushover")
	notif := &pushover.Message{
		Message:  message,
		Title:    title,
		Priority: config.ConfigData.Alerts.Pushover.Priority,
		Sound:    config.ConfigData.Alerts.Pushover.Sound,
		HTML:     true,
		TTL:      time.Duration(config.ConfigData.Alerts.Pushover.TTL) * time.Second,
	}

	// If emergency priority, set retry / expiration
	if notif.Priority == 2 {
		notif.Retry = time.Duration(config.ConfigData.Alerts.Pushover.Retry) * time.Second
		notif.Expire = time.Duration(config.ConfigData.Alerts.Pushover.Expire) * time.Second
	}

	// Add target devices if specified
	if config.ConfigData.Alerts.Pushover.Devices != "" {
		devices := strings.ReplaceAll(config.ConfigData.Alerts.Pushover.Devices, " ", "")
		notif.DeviceName = devices
	}

	log.Trace().
		Interface("payload", notif).
		Interface("recipient", "--secret removed--").
		Msg("Send Pushover alert")

	// Send notification
	if event.HasSnapshot {
		notif.AddAttachment(snapshot)
		response, err := push.SendMessage(notif, recipient)
		log.Trace().
			Interface("payload", response).
			Msg("Pushover response")
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Pushover").
				Msgf("Unable to send alert: %v", err)
			return
		}
	} else {
		response, err := push.SendMessage(notif, recipient)
		log.Trace().
			Interface("payload", response).
			Msg("Pushover response")
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Pushover").
				Msgf("Unable to send alert: %v", err)
			return
		}
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Pushover").
		Msgf("Alert sent")
}
