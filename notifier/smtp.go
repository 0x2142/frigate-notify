package notifier

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/wneessen/go-mail"
)

var SMTPUser string
var SMTPRecipients []string
var SMTPServer string
var SMTPPassword string
var SMTPTLS bool
var SMTPPort int

// SendSMTP forwards alert data via email
func SendSMTP(message string, snapshot io.Reader) {
	// Set up email alert
	m := mail.NewMsg()
	m.From(SMTPUser)
	m.To(SMTPRecipients...)
	m.Subject(AlertTitle)
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

	// Connect to SMTP server
	c, err := mail.NewClient(SMTPServer, mail.WithPort(SMTPPort), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(SMTPUser), mail.WithPassword(SMTPPassword))
	if err != nil {
		log.Print("Failed to connect to SMTP Server: ", err)
	}

	// Send message
	if err := c.DialAndSend(m); err != nil {
		log.Print("Failed to send SMTP message: ", err)
	}
	log.Println("SMTP alert sent")

}

// ParseSMTPRecipients splits individual email addresses from config file
func ParseSMTPRecipients(r string) {
	list := strings.Split(r, ",")
	for _, addr := range list {
		SMTPRecipients = append(SMTPRecipients, strings.TrimSpace(addr))
	}
}
