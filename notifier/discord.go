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
	var file *discord.File
	var err error

	title := fmt.Sprintf("**%v**\n\n", AlertTitle)
	message = title + message

	// Send alert & attach snapshot if one was saved
	if snapshot != nil {
		file = discord.NewFile("snapshot.jpg", "", snapshot)
		_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetContent(message).SetFiles(file).Build())
	} else {
		message += "\nNo snapshot saved."
		_, err = client.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetContent(message).Build())
	}
	if err != nil {
		log.Printf("Unable to send Discord Alert: %v", err)
	}
	log.Println("Discord alert sent")
}
