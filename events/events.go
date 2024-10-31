package events

import (
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/rs/zerolog/log"
)

// processEvent handles preparing event for alerting
func processEvent(event models.Event) {
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
	if !checkFilters(event) {
		return
	}

	// Send alert with snapshot
	notifier.SendAlert(event)
}
