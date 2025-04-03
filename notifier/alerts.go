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
func SendAlert(events []models.Event) {
	config.Internal.Status.LastNotification = time.Now()

	// Collect snapshot, if available
	var snapshot io.Reader
	for _, event := range events {
		if event.HasSnapshot {
			snapshot = GetSnapshot(event.ID)
			break
		}
	}

	// Set extra event details & get event used for notifications
	event := setExtras(events)

	// Create copy of snapshot for each alerting method
	var snap []byte
	if snapshot != nil {
		snap, _ = io.ReadAll(snapshot)
	} else {
		event.HasSnapshot = false
	}

	// Send Alerts
	// Apprise API
	for id, profile := range config.ConfigData.Alerts.AppriseAPI {
		if profile.Enabled {
			provider := notifMeta{name: "apprise-api", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendAppriseAPI(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Discord
	for id, profile := range config.ConfigData.Alerts.Discord {
		if profile.Enabled {
			provider := notifMeta{name: "discord", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendDiscordMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Gotify
	for id, profile := range config.ConfigData.Alerts.Gotify {
		if profile.Enabled {
			provider := notifMeta{name: "gotify", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendGotifyPush(event, provider)
			}
		}
	}
	// Matrix
	for id, profile := range config.ConfigData.Alerts.Matrix {
		if profile.Enabled {
			provider := notifMeta{name: "matrix", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendMatrix(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Mattermost
	for id, profile := range config.ConfigData.Alerts.Mattermost {
		if profile.Enabled {
			provider := notifMeta{name: "mattermost", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendMattermost(event, provider)
			}
		}
	}
	// Ntfy
	for id, profile := range config.ConfigData.Alerts.Ntfy {
		if profile.Enabled {
			provider := notifMeta{name: "ntfy", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendNtfyPush(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Pushover
	for id, profile := range config.ConfigData.Alerts.Pushover {
		if profile.Enabled {
			provider := notifMeta{name: "pushover", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendPushoverMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Signal
	for id, profile := range config.ConfigData.Alerts.Signal {
		if profile.Enabled {
			provider := notifMeta{name: "signal", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendSignalMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// SMTP
	for id, profile := range config.ConfigData.Alerts.SMTP {
		if profile.Enabled {
			provider := notifMeta{name: "smtp", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendSMTP(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Telegram
	for id, profile := range config.ConfigData.Alerts.Telegram {
		if profile.Enabled {
			provider := notifMeta{name: "telegram", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
				go SendTelegramMessage(event, bytes.NewReader(snap), provider)
			}
		}
	}
	// Webhook
	for id, profile := range config.ConfigData.Alerts.Webhook {
		if profile.Enabled {
			provider := notifMeta{name: "webhook", index: id}
			if checkAlertFilters(events, profile.Filters, provider) {
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
	var response []byte

	attempts := 0
	max_attempts := config.ConfigData.Alerts.General.MaxSnapRetry
	for attempts < max_attempts {
		var err error
		response, err = util.HTTPGet(url.String(), config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
		if err != nil {
			attempts += 1
			if err.Error() == "404" {
				time.Sleep(2 * time.Second)
				log.Info().
					Str("event_id", eventID).
					Int("attempt", attempts).
					Int("max_attempts", max_attempts).
					Msgf("Waiting for snapshot to be available")
				continue
			} else {
				log.Warn().
					Str("event_id", eventID).
					Err(err).
					Msgf("Could not access snapshot")
				return nil
			}
		} else {
			break
		}
	}
	if attempts == max_attempts {
		return nil
	}
	return bytes.NewReader(response)
}

// setExtras adds additional data into the event model to be used for templates
func setExtras(events []models.Event) models.Event {
	// Pull first event, which will be used to store info relevant to notifications
	key := events[0]

	// Set Event link
	key.Extra.EventLink = config.ConfigData.Frigate.PublicURL + "/api/events/" + key.ID + "/clip.mp4"

	// Add Frigate Major version metadata
	key.Extra.FrigateMajorVersion = config.Internal.FrigateVersion

	// Transform camera names, example: "test_camera" to "Test Camera"
	caser := cases.Title(language.Und)
	key.Extra.CameraName = caser.String(strings.ReplaceAll(key.Camera, "_", " "))

	// Assign Frigate URL to extra event fields
	key.Extra.LocalURL = config.ConfigData.Frigate.Server
	key.Extra.PublicURL = config.ConfigData.Frigate.PublicURL

	// Create list of all detected object (mostly applicable to /reviews)
	var labelList []string
	var sublabelList []string
	for _, event := range events {
		if !slices.Contains(labelList, event.Label) {
			labelList = append(labelList, event.Label)
		}
		if !slices.Contains(sublabelList, event.SubLabel) {
			sublabelList = append(sublabelList, event.SubLabel)
		}
	}
	key.Extra.LabelList = strings.Join(labelList, ", ")
	key.Extra.SubLabelList = strings.Join(sublabelList, ", ")

	// MQTT uses CurrentZones, Web API uses Zones
	// Combine into one object to use regardless of connection method
	for _, event := range events {
		key.Zones = append(key.Zones, event.CurrentZones...)
	}
	// Remove duplicates
	slices.Sort(key.Zones)
	key.Zones = slices.Compact(key.Zones)
	// Join zones into plain comma-separated string
	key.Extra.ZoneList = strings.Join(key.Zones, ", ")

	// If certain time format is provided, re-format date / time string
	eventTime := time.Unix(int64(key.StartTime), 0)
	key.Extra.FormattedTime = eventTime.String()
	if config.ConfigData.Alerts.General.TimeFormat != "" {
		key.Extra.FormattedTime = eventTime.Format(config.ConfigData.Alerts.General.TimeFormat)
	}

	// Calc TopScore percentage
	key.Extra.TopScorePercent = fmt.Sprintf("%v%%", int((key.TopScore * 100)))

	return key
}

// Build notification based on template
func renderMessage(sourceTemplate string, event models.Event, mtype string, provider string) string {
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
