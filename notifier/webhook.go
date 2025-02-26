package notifier

import (
	"fmt"
	"strings"

	"github.com/disgoorg/json"
	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

type WebhookPayload struct {
	Time         string   `json:"time"`
	ID           string   `json:"id"`
	Camera       string   `json:"camera"`
	Label        string   `json:"label"`
	SubLabel     string   `json:"sublabel"`
	Score        string   `json:"score"`
	Audio        string   `json:"audio"`
	CurrentZones []string `json:"current_zones"`
	EnteredZones []string `json:"entered_zones"`
	HasClip      bool     `json:"has_clip"`
	HasSnap      bool     `json:"has_snapshot"`
	Links        struct {
		Camera string `json:"camera"`
		Clip   string `json:"clip,omitempty"`
		Review string `json:"review,omitempty"`
		Snap   string `json:"snapshot,omitempty"`
	} `json:"links"`
}

// SendWebhook sends alert through HTTP POST to target webhook
func SendWebhook(event models.Event, provider notifMeta) {
	profile := config.ConfigData.Alerts.Webhook[provider.index]
	status := &config.Internal.Status.Notifications.Webhook[provider.index]

	// Build notification
	var message string
	payload, err := json.Marshal(profile.Template)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Webhook").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	if string(payload) != "null" {
		message = renderMessage(string(payload), event, "message", "Webhook")
	} else {
		defaultTemplate := WebhookPayload{
			Time:         event.Extra.FormattedTime,
			ID:           event.ID,
			Camera:       event.Extra.CameraName,
			Label:        event.Label,
			SubLabel:     event.SubLabel,
			Score:        event.Extra.TopScorePercent,
			Audio:        event.Extra.Audio,
			CurrentZones: event.CurrentZones,
			EnteredZones: event.EnteredZones,
			HasClip:      event.HasClip,
			HasSnap:      event.HasSnapshot,
		}
		if event.Extra.FrigateMajorVersion >= 14 {
			defaultTemplate.Links.Camera = fmt.Sprintf("%s/#%s", event.Extra.PublicURL, event.Camera)
		} else {
			defaultTemplate.Links.Camera = fmt.Sprintf("%s/cameras/%s", event.Extra.PublicURL, event.Camera)
		}
		if event.HasClip {
			defaultTemplate.Links.Clip = event.Extra.EventLink
		}
		if event.HasSnapshot {
			defaultTemplate.Links.Snap = fmt.Sprintf("%s/api/events/%s/snapshot.jpg", event.Extra.PublicURL, event.ID)
		}
		if event.Extra.ReviewLink != "" {
			defaultTemplate.Links.Review = event.Extra.ReviewLink
		}
		payload, _ = json.Marshal(defaultTemplate)
		message = string(payload)
	}

	headers := renderHTTPKV(profile.Headers, event, "headers", "Webhook")
	params := renderHTTPKV(profile.Params, event, "params", "Webhook")
	paramString := util.BuildHTTPParams(params...)
	if strings.ToUpper(profile.Method) == "GET" {
		_, err = util.HTTPGet(profile.Server, profile.Insecure, paramString, headers...)

	} else {
		_, err = util.HTTPPost(profile.Server, profile.Insecure, []byte(message), paramString, headers...)
	}

	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Webhook").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Webhook").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
