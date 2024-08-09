package notifier

import (
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
		return
	}
	if string(payload) != "null" {
		message = renderMessage(string(payload), event)
		log.Debug().
			Str("event_id", event.ID).
			Str("provider", "Webhook").
			Str("rendered_template", message).
			Msg("Custom message template used")
	} else {
		message = renderMessage("json", event)
	}

	headers := renderHeaders(config.ConfigData.Alerts.Webhook.Headers, event)
	_, err = util.HTTPPost(config.ConfigData.Alerts.Webhook.Server, config.ConfigData.Alerts.Webhook.Insecure, []byte(message), headers...)
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Webhook").
			Err(err).
			Msg("Unable to send alert")
		return
	}
	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Webhook").
		Msg("Alert sent")
}
