package notifier

import (
	"io"
)

var DiscordEnabled = false
var GotifyEnabled = false
var SMTPEnabled = false
var AlertTitle = "Frigate Alert"

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(message, snapshotURL string, snapshot io.Reader) {
	if DiscordEnabled {
		SendDiscordMessage(message, snapshot)
	}
	if GotifyEnabled {
		SendGotifyPush(message, snapshotURL)
	}
	if SMTPEnabled {
		SendSMTP(message, snapshot)
	}
}
