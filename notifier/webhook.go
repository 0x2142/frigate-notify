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
func SendWebhook(event models.Event) {
	// Build notification
	var message string
	payload, err := json.Marshal(config.ConfigData.Alerts.Webhook.Template)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Webhook").
			Err(err).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Webhook[0].NotifFailure(err.Error())
		return
	}
	if string(payload) != "null" {
		message = renderMessage(string(payload), event, "message", "Webhook")
	} else {
		message = renderMessage("json", event, "message", "Webhook")
	}

	headers := renderHTTPKV(config.ConfigData.Alerts.Webhook.Headers, event, "headers", "Webhook")
	params := renderHTTPKV(config.ConfigData.Alerts.Webhook.Params, event, "params", "Webhook")
	paramString := util.BuildHTTPParams(params...)
	if strings.ToUpper(config.ConfigData.Alerts.Webhook.Method) == "GET" {
		_, err = util.HTTPGet(config.ConfigData.Alerts.Webhook.Server, config.ConfigData.Alerts.Webhook.Insecure, paramString, headers...)

	} else {
		_, err = util.HTTPPost(config.ConfigData.Alerts.Webhook.Server, config.ConfigData.Alerts.Webhook.Insecure, []byte(message), paramString, headers...)
	}

	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Webhook").
			Err(err).
			Msg("Unable to send alert")
		config.Internal.Status.Notifications.Webhook[0].NotifFailure(err.Error())
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Webhook").
		Msg("Alert sent")
	config.Internal.Status.Notifications.Webhook[0].NotifSuccess()
}
