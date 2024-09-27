package notifier

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

var TemplateFiles embed.FS

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(event models.Event, snapshotURL string, snapshot io.Reader, eventid string) {
	// Add Frigate Major version metadata
	event.Extra.FrigateMajorVersion = config.ConfigData.Frigate.Version
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

// setExtras adds additional data into the event model to be used for templates
func setExtras(event models.Event) models.Event {
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

	return event
}

// Build notification based on template
func renderMessage(sourceTemplate string, event models.Event) string {
	event = setExtras(event)

	// Render template
	var tmpl *template.Template
	var err error
	if sourceTemplate == "markdown" || sourceTemplate == "plaintext" || sourceTemplate == "html" || sourceTemplate == "json" {
		tmpl = template.Must(template.ParseFS(TemplateFiles, "templates/"+sourceTemplate+".template"))
	} else {
		tmpl, err = template.New("custom").Funcs(template.FuncMap{"env": includeenv}).Parse(sourceTemplate)
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

// Build HTTP headers or params based on template
func renderHTTPKV(list []map[string]string, event models.Event, kvtype string) []map[string]string {
	event = setExtras(event)

	var renderedList []map[string]string

	for _, item := range list {
		for k, v := range item {
			// Render
			tmpl, err := template.New("custom").Funcs(template.FuncMap{"env": includeenv}).Parse(v)
			if err != nil {
				log.Warn().Err(err).Msgf("Failed to render HTTP %s", kvtype)
			}

			var renderedTemplate bytes.Buffer
			err = tmpl.Execute(&renderedTemplate, event)
			if err != nil {
				log.Fatal().
					Err(err).
					Msgf("Failed to render HTTP %s", kvtype)
			}

			v = renderedTemplate.String()
			renderedList = append(renderedList, map[string]string{k: v})
		}
	}

	return renderedList
}

// includeenv retrieves environment variables for use within templates
func includeenv(env string) string {
	if strings.HasPrefix(env, "FN_") {
		value, ok := os.LookupEnv(env)
		if !ok {
			log.Warn().
				Msgf("Could not find matching env: %v", env)
			return ""
		}
		return value
	} else {
		log.Warn().
			Msg("Env vars used in templates must contain FN_ prefix")
		return ""

	}
}
