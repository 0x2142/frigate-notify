package notifier

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/gomarkdown/markdown"
	"github.com/gregdel/pushover"
)

// SendPushoverMessage sends alert message through Pushover service
func SendPushoverMessage(message string, snapshot io.Reader, eventid string) {
	push := pushover.New(config.ConfigData.Alerts.Pushover.Token)
	recipient := pushover.NewRecipient(config.ConfigData.Alerts.Pushover.Userkey)

	// Convert message to HTML & strip newline characters
	htmlMessage := string(markdown.ToHTML([]byte(message), nil, nil))
	htmlMessage = strings.Replace(htmlMessage, "\n", "", -1)

	// Create new message
	notif := &pushover.Message{
		Message:  htmlMessage,
		Title:    config.ConfigData.Alerts.General.Title,
		Priority: config.ConfigData.Alerts.Pushover.Priority,
		HTML:     true,
		TTL:      time.Duration(config.ConfigData.Alerts.Pushover.TTL) * time.Second,
	}

	// If emergency priority, set retry / expiration
	if notif.Priority == 2 {
		notif.Retry = time.Duration(config.ConfigData.Alerts.Pushover.Retry) * time.Second
		notif.Expire = time.Duration(config.ConfigData.Alerts.Pushover.Expire) * time.Second
		fmt.Print(notif.Retry, notif.Expire)
	}

	// Add target devices if specified
	if config.ConfigData.Alerts.Pushover.Devices != "" {
		devices := strings.ReplaceAll(config.ConfigData.Alerts.Pushover.Devices, " ", "")
		notif.DeviceName = devices
	}

	// Send notification
	if snapshot != nil {
		notif.AddAttachment(snapshot)
		if _, err := push.SendMessage(notif, recipient); err != nil {
			log.Print("Event ID %v - Error sending Pushover notification:", eventid, err)
			return
		}
	} else {
		if _, err := push.SendMessage(notif, recipient); err != nil {
			log.Print("Event ID %v - Error sending Pushover notification:", eventid, err)
			return
		}
	}

	log.Printf("Event ID %v - Pushover alert sent", eventid)
}
