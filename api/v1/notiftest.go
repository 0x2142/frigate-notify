package apiv1

import (
	"context"
	"encoding/json"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
	"github.com/rs/zerolog/log"
)

type NotifTestOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// PostNotifTest collects most the recent event from Frigate & sends a test notification to configured providers
func PostNotifTest(ctx context.Context, input *struct{}) (*NotifTestOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/notif_test").
		Str("method", "POST").
		Msg("Received API request")

	resp := &NotifTestOutput{}
	resp.Body.Message = "ok"

	go func() {
		log.Info().Msg("Received request to test notifications")

		// Query frigate API for most recent event
		var events []models.Event
		uri := "/api/events"
		params := "?include_thumbnails=0&limit=1"
		url := config.ConfigData.Frigate.Server + uri + params

		response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("Cannot get event from %s", url)
		}
		json.Unmarshal([]byte(response), &events)

		// Send test notification
		notifier.SendAlert(events)
	}()

	log.Trace().
		Str("uri", V1_PREFIX+"/notif_test").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
