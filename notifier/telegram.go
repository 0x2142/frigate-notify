package notifier

import (
	"io"
	"log"
	"strings"

	"github.com/0x2142/frigate-notify/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gomarkdown/markdown"
)

// SendTelegramMessage sends alert through Telegram to individual users
func SendTelegramMessage(message string, snapshot io.Reader, eventid string) {
	bot, err := tgbotapi.NewBotAPI(config.ConfigData.Alerts.Telegram.Token)
	if err != nil {
		log.Print("Event ID %v - Failed to connect to Telegram:", eventid, err)
		return
	}

	// Convert message to HTML & remove tags not permitted by Telegram
	htmlMessage := string(markdown.ToHTML([]byte(message), nil, nil))
	htmlMessage = strings.Replace(htmlMessage, "<p>", "", -1)
	htmlMessage = strings.Replace(htmlMessage, "</p>", "", -1)
	htmlMessage = strings.Replace(htmlMessage, "<br>", "", -1)

	if snapshot != nil {
		// Attach & send snapshot
		photo := tgbotapi.NewPhoto(config.ConfigData.Alerts.Telegram.ChatID, tgbotapi.FileReader{Name: "Snapshot", Reader: snapshot})
		photo.Caption = htmlMessage
		photo.ParseMode = "HTML"
		if _, err := bot.Send(photo); err != nil {
			log.Print("Event ID %v - Failed to send alert via Telegram:", eventid, err)
			return
		}
	} else {
		// Send plain text message if no snapshot available
		htmlMessage += "No snapshot saved."
		msg := tgbotapi.NewMessage(config.ConfigData.Alerts.Telegram.ChatID, htmlMessage)
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Print("Event ID %v - Failed to send alert via Telegram:", eventid, err)
			return
		}
	}
	log.Printf("Event ID %v - Telegram alert sent", eventid)
}
