package notifier

import (
	"io"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendTelegramMessage sends alert through Telegram to individual users
func SendTelegramMessage(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Telegram[provider.index]

	// Build notification
	var message string
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Telegram")
	} else {
		message = renderMessage("html", event, "message", "Telegram")
		message = strings.ReplaceAll(message, "<br />", "")
	}

	bot, err := tgbotapi.NewBotAPI(profile.Token)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Telegram").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Telegram[0].NotifFailure(err.Error())

		return
	}

	if event.HasSnapshot {
		// Attach & send snapshot
		photo := tgbotapi.NewPhoto(profile.ChatID, tgbotapi.FileReader{Name: "Snapshot", Reader: snapshot})
		photo.Caption = message
		photo.ParseMode = "HTML"
		response, err := bot.Send(photo)
		log.Trace().
			Interface("content", response).
			Int("provider_id", provider.index).
			Msg("Send Telegram Alert")
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Telegram").
				Int("provider_id", provider.index).
				Err(err).
				Msg("Unable to send alert")
			config.Internal.Status.Notifications.Telegram[0].NotifFailure(err.Error())
			return
		}
	} else {
		// Send plain text message if no snapshot available
		msg := tgbotapi.NewMessage(profile.ChatID, message)
		msg.ParseMode = "HTML"
		response, err := bot.Send(msg)
		log.Trace().
			Interface("content", response).
			Int("provider_id", provider.index).
			Msg("Send Telegram Alert")
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Telegram").
				Int("provider_id", provider.index).
				Err(err).
				Msg("Unable to send alert")
			config.Internal.Status.Notifications.Telegram[0].NotifFailure(err.Error())
			return
		}
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Telegram").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	config.Internal.Status.Notifications.Telegram[0].NotifSuccess()
}
