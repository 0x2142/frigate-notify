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
func SendPushoverMessage(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Pushover[provider.index]
	status := &config.Internal.Status.Notifications.Pushover[provider.index]

	// Build notification
	var message string
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Pushover")
	} else {
		message = renderMessage("html", event, "message", "Pushover")
		message = strings.ReplaceAll(message, "<br />", "")
	}

	push := pushover.New(profile.Token)
	recipient := pushover.NewRecipient(profile.Userkey)

	// Create new message
	var title string
	if profile.Title != "" {
		title = renderMessage(profile.Title, event, "title", "pushover")
	} else {
		title = renderMessage(config.ConfigData.Alerts.General.Title, event, "title", "pushover")
	}
	notif := &pushover.Message{
		Message:  message,
		Title:    title,
		Priority: profile.Priority,
		Sound:    profile.Sound,
		HTML:     true,
		TTL:      time.Duration(profile.TTL) * time.Second,
	}

	// Add links
	if event.Extra.ReviewLink != "" {
		notif.URL = event.Extra.ReviewLink
		notif.URLTitle = "Review Event"
	} else {
		notif.URL = event.Extra.EventLink
		notif.URLTitle = "View Clip"
	}

	// If emergency priority, set retry / expiration
	if notif.Priority == 2 {
		notif.Retry = time.Duration(profile.Retry) * time.Second
		notif.Expire = time.Duration(profile.Expire) * time.Second
	}

	// Add target devices if specified
	if profile.Devices != "" {
		devices := strings.ReplaceAll(profile.Devices, " ", "")
		notif.DeviceName = devices
	}

	log.Trace().
		Interface("payload", notif).
		Interface("recipient", "--secret removed--").
		Int("provider_id", provider.index).
		Msg("Send Pushover alert")

	// Send notification
	if event.HasSnapshot {
		notif.AddAttachment(snapshot)
		response, err := push.SendMessage(notif, recipient)
		log.Trace().
			Interface("payload", response).
			Int("provider_id", provider.index).
			Msg("Pushover response")
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Pushover").
				Int("provider_id", provider.index).
				Msgf("Unable to send alert: %v", err)
			status.NotifFailure(err.Error())
			return
		}
	} else {
		response, err := push.SendMessage(notif, recipient)
		log.Trace().
			Interface("payload", response).
			Int("provider_id", provider.index).
			Msg("Pushover response")
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Pushover").
				Int("provider_id", provider.index).
				Msgf("Unable to send alert: %v", err)
			status.NotifFailure(err.Error())
			return
		}
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Pushover").
		Int("provider_id", provider.index).
		Msgf("Alert sent")
	status.NotifSuccess()

}
