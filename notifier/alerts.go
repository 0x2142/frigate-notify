package notifier

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(event models.Event, snapshotURL string, snapshot io.Reader, eventid string) {
	// Add Frigate Major version metadata
	event.Extra.FrigateMajorVersion = config.ConfigData.Frigate.Version
	// Drop event if no snapshot or clip is available - Event is likely being filtered on Frigate side.
	// For example, if a camera has `required_zones` set - then there may not be any clip or snap until
	// object moves into required zone
	if !event.HasClip && !event.HasSnapshot {
		log.Info().
			Str("event_id", event.ID).
			Msg("Event dropped - No snapshot or clip available")
		return
	}
	// Drop event if no snapshot & skip_nosnap is true
	if !event.HasSnapshot && strings.ToLower(config.ConfigData.Alerts.General.NoSnap) == "drop" {
		log.Info().
			Str("event_id", event.ID).
			Msg("Event dropped - No snapshot available")
		return
	}
	// Create copy of snapshot for each alerting method
	var snap []byte
	if snapshot != nil {
		snap, _ = io.ReadAll(snapshot)
	}
	if config.ConfigData.Alerts.Discord.Enabled {
		go SendDiscordMessage(event, bytes.NewReader(snap))
	}
	if config.ConfigData.Alerts.Gotify.Enabled {
		go SendGotifyPush(event, snapshotURL)
	}
	if config.ConfigData.Alerts.SMTP.Enabled {
		go SendSMTP(event, bytes.NewReader(snap))
	}
	if config.ConfigData.Alerts.Telegram.Enabled {
		go SendTelegramMessage(event, bytes.NewReader(snap))
	}
	if config.ConfigData.Alerts.Pushover.Enabled {
		go SendPushoverMessage(event, bytes.NewReader(snap))
	}
	if config.ConfigData.Alerts.Ntfy.Enabled {
		go SendNtfyPush(event, bytes.NewReader(snap))
	}
	if config.ConfigData.Alerts.Webhook.Enabled {
		go SendWebhook(event)
	}
}

// Build notification based on template
func renderMessage(sourceTemplate string, event models.Event) string {
	// Assign Frigate URL to extra event fields
	event.Extra.LocalURL = config.ConfigData.Frigate.Server
	event.Extra.PublicURL = config.ConfigData.Frigate.PublicURL

	// MQTT uses CurrentZones, Web API uses Zones
	// Combine into one object to use regardless of connection method
	event.Zones = append(event.Zones, event.CurrentZones...)
	// Join zones into plain comma-separated string
	event.Extra.ZoneList = strings.Join(event.Zones, ", ")

	// If certain time format is provided, re-format date / time string
	eventTime := time.Unix(int64(event.StartTime), 0)
	event.Extra.FormattedTime = eventTime.String()
	if config.ConfigData.Alerts.General.TimeFormat != "" {
		event.Extra.FormattedTime = eventTime.Format(config.ConfigData.Alerts.General.TimeFormat)
	}

	// For Web API query, top-level top_score value is no longer used
	// So need to replace it with data.top_score value
	if event.TopScore == 0 {
		event.TopScore = event.Data.TopScore
	}
	// Calc TopScore percentage
	event.Extra.TopScorePercent = fmt.Sprintf("%v%%", int((event.TopScore * 100)))

	// Render template
	var tmpl *template.Template
	var err error
	if sourceTemplate == "markdown" || sourceTemplate == "plaintext" || sourceTemplate == "html" || sourceTemplate == "json" {
		var templateFile = "./templates/" + sourceTemplate + ".template"
		tmpl = template.Must(template.ParseFiles(templateFile))
	} else {
		tmpl, err = template.New("custom").Parse(sourceTemplate)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to render event message")
		}
	}

	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, event)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Failed to render event message")
	}

	return renderedTemplate.String()

}
