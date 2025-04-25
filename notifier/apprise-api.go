package notifier

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

type AppriseAPIPayload struct {
	Title      string                 `json:"title"`
	Body       string                 `json:"body"`
	Tags       string                 `json:"tags,omitempty"`
	URLs       string                 `json:"urls,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Format     string                 `json:"format,omitempty"`
	Attachment []AppriseAPIAttachment `json:"attachment,omitempty"`
}

type AppriseAPIAttachment struct {
	Filename string `json:"filename"`
	Base64   string `json:"base64"`
	Mimetype string `json:"mimetype"`
}

// SendAppriseAPI forwards alert messages to Apprise API notification server
func SendAppriseAPI(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.AppriseAPI[provider.index]
	status := &config.Internal.Status.Notifications.AppriseAPI[provider.index]

	// Build notification
	var message string
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Apprise API")
	} else {
		message = renderMessage("markdown", event, "message", "Apprise API")
	}

	payload := AppriseAPIPayload{
		Title: renderMessage(config.ConfigData.Alerts.General.Title, event, "title", "Apprise API"),
		Body:  message}

	if len(profile.URLs) != 0 {
		payload.URLs = strings.Join(profile.URLs, ",")
	}

	if len(profile.Tags) != 0 {
		payload.Tags = strings.Join(profile.Tags, ",")
	}

	if event.HasSnapshot {
		img, _ := io.ReadAll(snapshot)
		attach := AppriseAPIAttachment{
			Filename: "snapshot.jpg",
			Base64:   base64.StdEncoding.EncodeToString(img),
			Mimetype: http.DetectContentType(img),
		}
		payload.Attachment = append(payload.Attachment, attach)
	}

	header := map[string]string{"Content-Type": "application/json"}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "apprise_api").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	// Build URL
	appriseapiURL := profile.Server
	if !strings.HasSuffix(appriseapiURL, "/notify") {
		appriseapiURL += "/notify"
	}
	if profile.Token != "" {
		appriseapiURL += "/" + profile.Token
	}

	response, err := util.HTTPPost(appriseapiURL, profile.Insecure, data, "", header)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Apprise API").
			Str("response", string(response)).
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Apprise API").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
