package notifier

import (
	"io"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

// SendTelegramMessage sends alert through Telegram to individual users
func SendTelegramMessage(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Telegram[provider.index]
	status := &config.Internal.Status.Notifications.Telegram[provider.index]

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
		status.NotifFailure(err.Error())

		return
	}

	var response tgbotapi.Message
	if event.HasClip && profile.SendClip {
		clip := GetClip(event)
		msg := tgbotapi.NewVideo(profile.ChatID, tgbotapi.FileReader{Name: "Clip", Reader: clip})
		if profile.MessageThreadID != 0 {
			msg.MessageThreadID = profile.MessageThreadID
		}
		msg.Caption = message
		msg.ParseMode = "HTML"
		response, err = bot.Send(msg)
	} else if event.HasSnapshot {
		// Attach & send snapshot
		msg := tgbotapi.NewPhoto(profile.ChatID, tgbotapi.FileReader{Name: "Snapshot", Reader: snapshot})
		if profile.MessageThreadID != 0 {
			msg.MessageThreadID = profile.MessageThreadID
		}
		msg.Caption = message
		msg.ParseMode = "HTML"
		response, err = bot.Send(msg)
	} else {
		// Send plain text message if no snapshot available
		msg := tgbotapi.NewMessage(profile.ChatID, message)
		if profile.MessageThreadID != 0 {
			msg.MessageThreadID = profile.MessageThreadID
		}
		msg.ParseMode = "HTML"
		response, err = bot.Send(msg)
	}
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
		status.NotifFailure(err.Error())
		return
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Telegram").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
