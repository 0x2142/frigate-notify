package notifier

import (
	"strings"

	"github.com/disgoorg/json"
	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

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
		message = renderMessage("json", event, "message", "Webhook")
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
