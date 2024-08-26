package notifier

import (
	"crypto/tls"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/wneessen/go-mail"
)

// SendSMTP forwards alert data via email
func SendSMTP(event models.Event, snapshot io.Reader) {
	// Build notification
	var message string
	if config.ConfigData.Alerts.SMTP.Template != "" {
		message = renderMessage(config.ConfigData.Alerts.SMTP.Template, event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "SMTP").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("html", event)
	}

	// Set up email alert
	m := mail.NewMsg()
	m.From(config.ConfigData.Alerts.SMTP.From)
	m.To(ParseSMTPRecipients()...)
	m.Subject(config.ConfigData.Alerts.General.Title)
	// Attach snapshot if one exists
	if event.HasSnapshot {
		m.AttachReader("snapshot.jpg", snapshot)
	}

	// Convert message body to HTML
	m.SetBodyString(mail.TypeTextHTML, message)

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
	// Disable certificate verification if needed
	if config.ConfigData.Alerts.SMTP.Insecure {
		c.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "SMTP").
			Err(err).
			Msg("Unable to send alert")
	}

	// Send message
	if err := c.DialAndSend(m); err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "SMTP").
			Err(err).
			Msg("Unable to send alert")
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "SMTP").
		Msg("Alert sent")

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
