package notifier

import (
	"fmt"
	"io"
	"log"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
)

var client webhook.Client

var DiscordWebhookURL string

// SetupDiscord creates a Discord webhook client
func SetupDiscord() {
	// Connect to Discord
	var err error
	client, err = webhook.NewWithURL(DiscordWebhookURL)
	if err != nil {
		log.Printf("Unable to send Discord Alert: %v", err)
	}
}

// SendDiscordMessage pushes alert message to Discord via webhook
func SendDiscordMessage(message string, snapshot io.Reader) {
	var err error

	title := fmt.Sprintf("**%v**\n\n", AlertTitle)
	message = title + message

	// Send alert & attach snapshot if one was saved
	if snapshot != nil {
		image := discord.NewFile("snapshot.jpg", "", snapshot)
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetImage("attachment://snapshot.jpg").SetColor(5793266).Build()
		_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).SetFiles(image).Build())
	} else {
		message += "\nNo snapshot saved."
		embed := discord.NewEmbedBuilder().SetDescription(message).SetTitle(title).SetColor(5793266).Build()
		_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetEmbeds(embed).Build())

	}
	if err != nil {
		log.Printf("Unable to send Discord Alert: %v", err)
	}
	log.Println("Discord alert sent")
}
