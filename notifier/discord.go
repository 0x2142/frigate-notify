package notifier

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

// SendDiscordMessage pushes alert message to Discord via webhook
func SendDiscordMessage(event models.Event, snapshot io.Reader) {
	var err error
	var message string
	// Build notification
	if config.ConfigData.Alerts.Discord.Template != "" {
		message = renderMessage(config.ConfigData.Alerts.Discord.Template, event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "Discord").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("markdown", event)
	}

	// Connect to Discord
	client, err := webhook.NewWithURL(config.ConfigData.Alerts.Discord.Webhook)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Discord").
			Err(err).
			Msg("Unable to send alert")
	}
	defer client.Close(context.TODO())

	title := renderMessage(config.ConfigData.Alerts.General.Title, event)
	title = fmt.Sprintf("**%v**\n\n", title)
	message = title + message

	// Send alert & attach snapshot if one was saved
	var msg *discord.Message
	if event.HasSnapshot {
		image := discord.NewFile("snapshot.jpg", "", snapshot)
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetImage("attachment://snapshot.jpg").SetColor(5793266).Build()
		msg, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).SetFiles(image).Build())
		log.Trace().
			Str("event_id", event.ID).
			Interface("payload", msg).
			Msg("Send Discord Alert")

	} else {
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetColor(5793266).Build()
		msg, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).Build())
		log.Trace().
			Str("event_id", event.ID).
			Interface("payload", msg).
			Msg("Send Discord Alert")
	}
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Discord").
			Err(err).
			Msg("Unable to send alert")
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Discord").
		Msg("Alert sent")
}
