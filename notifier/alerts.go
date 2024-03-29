package notifier

import (
	"bytes"
	"io"

	"github.com/0x2142/frigate-notify/config"
)

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(message, snapshotURL string, snapshot io.Reader) {
	// Create copy of snapshot for each alerting method
	var snap []byte
	if snapshot != nil {
		snap, _ = io.ReadAll(snapshot)
	}
	if config.ConfigData.Alerts.Discord.Enabled {
		SendDiscordMessage(message, bytes.NewReader(snap))
	}
	if config.ConfigData.Alerts.Gotify.Enabled {
		SendGotifyPush(message, snapshotURL)
	}
	if config.ConfigData.Alerts.SMTP.Enabled {
		SendSMTP(message, bytes.NewReader(snap))
	}
}
