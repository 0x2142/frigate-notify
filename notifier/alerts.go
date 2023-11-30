package notifier

import (
	"bytes"
	"io"
)

var DiscordEnabled = false
var GotifyEnabled = false
var SMTPEnabled = false
var AlertTitle = "Frigate Alert"

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(message, snapshotURL string, snapshot io.Reader) {
	// Create copy of snapshot for each alerting method
	snap, _ := io.ReadAll(snapshot)
	if DiscordEnabled {
		SendDiscordMessage(message, bytes.NewReader(snap))
	}
	if GotifyEnabled {
		SendGotifyPush(message, snapshotURL)
	}
	if SMTPEnabled {
		SendSMTP(message, bytes.NewReader(snap))
	}
}
