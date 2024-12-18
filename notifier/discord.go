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
func SendDiscordMessage(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Discord[provider.index]

	var err error
	var message string
	// Build notification
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Discord")
	} else {
		message = renderMessage("markdown", event, "message", "Discord")
	}

	// Connect to Discord
	client, err := webhook.NewWithURL(profile.Webhook)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Discord").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Discord[0].NotifFailure(err.Error())
	}
	defer client.Close(context.TODO())

	title := renderMessage(config.ConfigData.Alerts.General.Title, event, "title", "Discord")
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
			Int("provider_id", provider.index).
			Interface("payload", msg).
			Msg("Send Discord Alert")

	} else {
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetColor(5793266).Build()
		msg, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).Build())
		log.Trace().
			Str("event_id", event.ID).
			Int("provider_id", provider.index).
			Interface("payload", msg).
			Msg("Send Discord Alert")
	}
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Discord").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Discord[0].NotifFailure(err.Error())
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Discord").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	config.Internal.Status.Notifications.Discord[0].NotifSuccess()
}
