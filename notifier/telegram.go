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
func SendTelegramMessage(event models.Event, snapshot io.Reader) {
	// Build notification
	var message string
	if config.ConfigData.Alerts.Telegram.Template != "" {
		message = renderMessage(config.ConfigData.Alerts.Telegram.Template, event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "Telegram").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("html", event)
		message = strings.ReplaceAll(message, "<br />", "")
	}

	bot, err := tgbotapi.NewBotAPI(config.ConfigData.Alerts.Telegram.Token)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Telegram").
			Err(err).
			Msg("Unable to send alert")
		return
	}

	if event.HasSnapshot {
		// Attach & send snapshot
		photo := tgbotapi.NewPhoto(config.ConfigData.Alerts.Telegram.ChatID, tgbotapi.FileReader{Name: "Snapshot", Reader: snapshot})
		photo.Caption = message
		photo.ParseMode = "HTML"
		if _, err := bot.Send(photo); err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Telegram").
				Err(err).
				Msg("Unable to send alert")
			return
		}
	} else {
		// Send plain text message if no snapshot available
		message += "No snapshot saved."
		msg := tgbotapi.NewMessage(config.ConfigData.Alerts.Telegram.ChatID, message)
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Telegram").
				Err(err).
				Msg("Unable to send alert")
			return
		}
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Telegram").
		Msg("Alert sent")
}
