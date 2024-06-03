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
		message = renderMessage(config.ConfigData.Alerts.Pushover.Template, event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "Pushover").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("html", event)
		message = strings.ReplaceAll(message, "<br />", "")
	}

	push := pushover.New(config.ConfigData.Alerts.Pushover.Token)
	recipient := pushover.NewRecipient(config.ConfigData.Alerts.Pushover.Userkey)

	// Create new message
	notif := &pushover.Message{
		Message:  message,
		Title:    config.ConfigData.Alerts.General.Title,
		Priority: config.ConfigData.Alerts.Pushover.Priority,
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

	// Send notification
	if event.HasSnapshot {
		notif.AddAttachment(snapshot)
		if _, err := push.SendMessage(notif, recipient); err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Pushover").
				Msgf("Unable to send alert: %v", err)
			return
		}
	} else {
		if _, err := push.SendMessage(notif, recipient); err != nil {
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
