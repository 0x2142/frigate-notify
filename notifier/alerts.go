package notifier

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/url"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

var TemplateFiles embed.FS

type notifMeta struct {
	name  string
	index int
}

// SendAlert forwards alert information to all enabled alerting methods
func SendAlert(event models.Event) {
	config.Internal.Status.LastNotification = time.Now()
	// Collect snapshot, if available
	var snapshot io.Reader
	if event.HasSnapshot {
		snapshot = GetSnapshot(event.ID)
	}

	// Set Event link
	event.Extra.EventLink = config.ConfigData.Frigate.PublicURL + "/api/events/" + event.ID + "/clip.mp4"

	// Add Frigate Major version metadata
	event.Extra.FrigateMajorVersion = config.Internal.FrigateVersion

	// Create copy of snapshot for each alerting method
	var snap []byte
	if snapshot != nil {
		snap, _ = io.ReadAll(snapshot)
	}

	// Send Alerts
	// Discord
	for id, profile := range config.ConfigData.Alerts.Discord {
		if profile.Enabled {
			provider := notifMeta{name: "discord", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendDiscordMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Gotify
	for id, profile := range config.ConfigData.Alerts.Gotify {
		if profile.Enabled {
			provider := notifMeta{name: "gotify", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendGotifyPush(event, provider)
			}
		}
	}
	// SMTP
	for id, profile := range config.ConfigData.Alerts.SMTP {
		if profile.Enabled {
			provider := notifMeta{name: "smtp", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendSMTP(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Telegram
	for id, profile := range config.ConfigData.Alerts.Telegram {
		if profile.Enabled {
			provider := notifMeta{name: "telegram", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendTelegramMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Pushover
	for id, profile := range config.ConfigData.Alerts.Pushover {
		if profile.Enabled {
			provider := notifMeta{name: "pushover", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendPushoverMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Ntfy
	for id, profile := range config.ConfigData.Alerts.Ntfy {
		if profile.Enabled {
			provider := notifMeta{name: "ntfy", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendNtfyPush(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Webhook
	for id, profile := range config.ConfigData.Alerts.Webhook {
		if profile.Enabled {
			provider := notifMeta{name: "webhook", index: id}
			if checkAlertFilters(event, profile.Filters, provider) {
				go SendWebhook(event, provider)
			}
		}
	}
}

// GetSnapshot downloads a snapshot from Frigate
func GetSnapshot(eventID string) io.Reader {
	// Add optional snapshot modifiers
	url, _ := url.Parse(config.ConfigData.Frigate.Server + "/api/events/" + eventID + "/snapshot.jpg")
	q := url.Query()
	if config.ConfigData.Alerts.General.SnapBbox {
		q.Add("bbox", "1")
	}
	if config.ConfigData.Alerts.General.SnapTimestamp {
		q.Add("timestamp", "1")
	}
	if config.ConfigData.Alerts.General.SnapCrop {
		q.Add("crop", "1")
	}
	url.RawQuery = q.Encode()
	response, err := util.HTTPGet(url.String(), config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
	if err != nil {
		log.Warn().
			Str("event_id", eventID).
			Err(err).
			Msgf("Could not access snaphot")
		return nil
	}

	return bytes.NewReader(response)
}

// setExtras adds additional data into the event model to be used for templates
func setExtras(event models.Event) models.Event {
	// Transform camera names, example: "test_camera" to "Test Camera"
	caser := cases.Title(language.Und)
	event.Extra.CameraName = caser.String(strings.ReplaceAll(event.Camera, "_", " "))

	// Assign Frigate URL to extra event fields
	event.Extra.LocalURL = config.ConfigData.Frigate.Server
	event.Extra.PublicURL = config.ConfigData.Frigate.PublicURL

	// MQTT uses CurrentZones, Web API uses Zones
	// Combine into one object to use regardless of connection method
	event.Zones = append(event.Zones, event.CurrentZones...)
	// Remove duplicates
	slices.Sort(event.Zones)
	event.Zones = slices.Compact(event.Zones)
	// Join zones into plain comma-separated string
	event.Extra.ZoneList = strings.Join(event.Zones, ", ")

	// If certain time format is provided, re-format date / time string
	eventTime := time.Unix(int64(event.StartTime), 0)
	event.Extra.FormattedTime = eventTime.String()
	if config.ConfigData.Alerts.General.TimeFormat != "" {
		event.Extra.FormattedTime = eventTime.Format(config.ConfigData.Alerts.General.TimeFormat)
	}

	// Calc TopScore percentage
	event.Extra.TopScorePercent = fmt.Sprintf("%v%%", int((event.TopScore * 100)))

	return event
}

// Build notification based on template
func renderMessage(sourceTemplate string, event models.Event, mtype string, provider string) string {
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

	log.Debug().
		Str("event_id", event.ID).
		Str("provider", provider).
		Str("rendered_template", renderedTemplate.String()).
		Msgf("Rendered %s template", mtype)

	return renderedTemplate.String()

}

// Build HTTP headers or params based on template
func renderHTTPKV(list []map[string]string, event models.Event, kvtype string, provider string) []map[string]string {
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

	log.Debug().
		Str("event_id", event.ID).
		Str("provider", provider).
		Interface("rendered_template", renderedList).
		Msgf("Rendered HTTP %s template", kvtype)

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
