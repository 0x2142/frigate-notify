package events

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
	"github.com/rs/zerolog/log"
)

// processEvent handles preparing event for alerting
func processEvent(event models.Event) {
	if config.ConfigData.Alerts.General.RecheckDelay != 0 {
		event = recheckEvent(event)
	}

	config.Internal.Status.LastEvent = time.Now()
	// For events collected via API, top-level top_score value is no longer used
	// So need to replace it with data.top_score value
	if event.TopScore == 0 {
		event.TopScore = event.Data.TopScore
	}

	// Convert to human-readable timestamp
	eventTime := time.Unix(int64(event.StartTime), 0)
	log.Info().
		Str("event_id", event.ID).
		Str("camera", event.Camera).
		Str("label", event.Label).
		Str("zones", strings.Join(event.CurrentZones, ",")).
		Msg("Processing event...")
	log.Debug().
		Str("event_id", event.ID).
		Msgf("Event start time: %s", eventTime)

	// Check that event passes configured filters
	if !checkEventFilters(event) {
		return
	}

	// Send alert with snapshot
	notifier.SendAlert(event)
}

func recheckEvent(event models.Event) models.Event {
	delay := config.ConfigData.Alerts.General.RecheckDelay
	log.Debug().
		Str("event_id", event.ID).
		Int("recheck_delay", delay).
		Msg("Waiting to re-check event details")
	time.Sleep(time.Duration(delay) * time.Second)
	log.Debug().
		Str("event_id", event.ID).
		Int("recheck_delay", delay).
		Msg("Re-checking event details")

	url := config.ConfigData.Frigate.Server + "/api/events/" + event.ID
	response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
	if err != nil {
		config.Internal.Status.Health = "frigate webapi unreachable"
		config.Internal.Status.Frigate.API = "unreachable"
		log.Error().
			Err(err).
			Msgf("Cannot get event from %s", url)
	}
	config.Internal.Status.Health = "ok"
	config.Internal.Status.Frigate.API = "ok"

	json.Unmarshal([]byte(response), &event)
	return event
}
