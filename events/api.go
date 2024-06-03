package frigate

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
	"golang.org/x/exp/slices"
)

const eventsURI = "/api/events"
const snapshotURI = "/snapshot.jpg"

// LastEventTime tracks the timestamp of the last event seen
var LastEventTime float64 = float64(time.Now().Unix())

// CheckForEvents queries for all detection events since last alert time
func CheckForEvents() {
	var params string
	if config.ConfigData.Frigate.WebAPI.TestMode {
		// For testing, pull 1 event immediately
		params = "?include_thumbnails=0&limit=1"
	} else {
		// Check for any events after last query time
		params = "?include_thumbnails=0&after=" + strconv.FormatFloat(LastEventTime, 'f', 6, 64)
	}

	url := config.ConfigData.Frigate.Server + eventsURI + params
	log.Debug().Msg("Checking for new events...")

	// Query events
	response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, config.ConfigData.Frigate.Headers...)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("Cannot get events from %s", url)
	}

	var events []models.Event

	json.Unmarshal([]byte(response), &events)

	log.Debug().Msgf("Found %v new events", len(events))

	for _, event := range events {
		// Convert to human-readable timestamp
		eventTime := time.Unix(int64(event.StartTime), 0)

		// Update last event check time with most recent timestamp
		if event.StartTime > LastEventTime {
			LastEventTime = event.StartTime
		}

		// Skip excluded cameras
		if slices.Contains(config.ConfigData.Frigate.Cameras.Exclude, event.Camera) {
			log.Debug().
				Str("event_id", event.ID).
				Str("camera", event.Camera).
				Msg("Skipping event from excluded camera")
			continue
		}

		log.Info().
			Str("event_id", event.ID).
			Str("camera", event.Camera).
			Str("label", event.Label).
			Str("zones", strings.Join(event.Zones, ",")).
			Msg("Event Detected")
		log.Debug().
			Str("event_id", event.ID).
			Msgf("Event start time: %s", eventTime)

		// Check that event passes the zone & label filters
		if !isAllowedZone(event.ID, event.Zones) || !isAllowedLabel(event.ID, event.Label) {
			return
		}

		// If snapshot was collected, pull down image to send with alert
		var snapshot io.Reader
		var snapshotURL string
		if event.HasSnapshot {
			snapshotURL = config.ConfigData.Frigate.Server + eventsURI + "/" + event.ID + snapshotURI
			snapshot = GetSnapshot(snapshotURL, event.ID)
		}

		// Send alert with snapshot
		notifier.SendAlert(event, snapshotURL, snapshot, event.ID)
	}

}

// GetSnapshot downloads a snapshot from Frigate
func GetSnapshot(snapshotURL, eventID string) io.Reader {
	response, err := util.HTTPGet(snapshotURL, config.ConfigData.Frigate.Insecure, config.ConfigData.Frigate.Headers...)
	if err != nil {
		log.Warn().
			Str("event_id", eventID).
			Err(err).
			Msgf("Could not access snaphot")
		return nil
	}

	return bytes.NewReader(response)
}
