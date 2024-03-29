package notifier

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/gomarkdown/markdown"
	"github.com/wneessen/go-mail"
)

// SendSMTP forwards alert data via email
func SendSMTP(message string, snapshot io.Reader) {
	// Set up email alert
	m := mail.NewMsg()
	m.From(config.ConfigData.Alerts.SMTP.User)
	m.To(ParseSMTPRecipients()...)
	m.Subject(config.ConfigData.Alerts.General.Title)
	// Attach snapshot if one exists
	if snapshot != nil {
		m.AttachReader("snapshot.jpg", snapshot)
	} else {
		message += "\n\nNo snapshot saved."
	}
	// Convert message body to HTML
	htmlMessage := markdown.ToHTML([]byte(message), nil, nil)
	m.SetBodyString(mail.TypeTextHTML, string(htmlMessage))

	time.Sleep(5 * time.Second)

	// Set up SMTP Connection
	c, err := mail.NewClient(config.ConfigData.Alerts.SMTP.Server, mail.WithPort(config.ConfigData.Alerts.SMTP.Port))
	// Add authentication params if needed
	if config.ConfigData.Alerts.SMTP.User != "" && config.ConfigData.Alerts.SMTP.Password != "" {
		c.SetSMTPAuth(mail.SMTPAuthPlain)
		c.SetUsername(config.ConfigData.Alerts.SMTP.User)
		c.SetPassword(config.ConfigData.Alerts.SMTP.Password)
	}
	// Mandatory TLS is enabled by default, so disable TLS if config flag is set
	if !config.ConfigData.Alerts.SMTP.TLS {
		c.SetTLSPolicy(mail.NoTLS)
	}

	if err != nil {
		log.Print("Failed to connect to SMTP Server: ", err)
	}

	// Send message
	if err := c.DialAndSend(m); err != nil {
		log.Print("Failed to send SMTP message: ", err)
		return
	}
	log.Println("SMTP alert sent")

}

// ParseSMTPRecipients splits individual email addresses from config file
func ParseSMTPRecipients() []string {
	var recipients []string
	list := strings.Split(config.ConfigData.Alerts.SMTP.Recipient, ",")
	for _, addr := range list {
		recipients = append(recipients, strings.TrimSpace(addr))
	}
	return recipients
}
