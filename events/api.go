package events

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

// LastQueryTime tracks the timestamp of the last event seen
var LastQueryTime float64 = float64(time.Now().Unix())

func QueryAPI() {
	appmode := strings.ToLower(config.ConfigData.App.Mode)
	var params string
	if config.ConfigData.Frigate.WebAPI.TestMode {
		// For testing, pull 1 event immediately
		params = "?include_thumbnails=0&limit=1"
	} else {
		// Check for any events after last query time
		params = "?include_thumbnails=0&after=" + strconv.FormatFloat(LastQueryTime, 'f', 6, 64)
	}

	var uri string
	if appmode == "reviews" {
		uri = "/api/review"
	} else {
		uri = "/api/events"
	}

	url := config.ConfigData.Frigate.Server + uri + params
	log.Debug().Msgf("Checking for new %s...", appmode)

	// Query API for reviews or events
	response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
	if err != nil {
		config.Internal.Status.Health = "frigate webapi unreachable"
		config.Internal.Status.Frigate.API = "unreachable"
		log.Error().
			Err(err).
			Msgf("Cannot get %s from %s", appmode, url)
	}
	config.Internal.Status.Health = "ok"
	config.Internal.Status.Frigate.API = "ok"

	switch appmode {
	case "reviews":
		var reviews []models.Review
		json.Unmarshal([]byte(response), &reviews)
		log.Debug().Msgf("Found %v new reviews", len(reviews))

		for _, review := range reviews {
			// Update last event check time with most recent timestamp
			if review.StartTime > LastQueryTime {
				LastQueryTime = review.StartTime
			}
			processReview(review)
		}
	case "events":
		var events []models.Event
		json.Unmarshal([]byte(response), &events)
		log.Debug().Msgf("Found %v new events", len(events))
		for _, event := range events {
			// Copy zones to CurrentZones, which is used for filters
			event.CurrentZones = event.Zones
			// Update last event check time with most recent timestamp
			if event.StartTime > LastQueryTime {
				LastQueryTime = event.StartTime
			}
			processEvent(event)
		}
	}

}

// Recheck Frigate event & wait for license plate recognition data
func waitforLPR(event *models.Event) models.Event {
	if event.Label == "car" {
		for _, attr := range event.Data.Attributes {
			if attr.Label == "license_plate" {
				log.Debug().
					Str("event_id", event.ID).
					Msg("Detected car & license plate - Waiting for license plate recognition...")
				max := 5
				current := 0
				for current < max {
					time.Sleep(2 * time.Second)
					current += 1

					log.Debug().
						Str("event_id", event.ID).
						Int("max_attempts", max).
						Int("current_attempts", current).
						Msg("Re-checking event details")

					url := config.ConfigData.Frigate.Server + "/api/events/" + event.ID
					response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
					if err != nil {
						config.Internal.Status.Health = "frigate webapi unreachable"
						config.Internal.Status.Frigate.API = "unreachable"
						log.Error().
							Err(err).
							Msgf("Cannot get event from %s", url)
						return *event
					}
					config.Internal.Status.Health = "ok"
					config.Internal.Status.Frigate.API = "ok"

					json.Unmarshal([]byte(response), &event)

					if event.Data.RecognizedLicensePlate != "" {
						log.Debug().
							Str("event_id", event.ID).
							Msg("License plate data received")
						return *event
					} else {
						log.Debug().
							Str("event_id", event.ID).
							Msg("No license plate data yet")
						continue
					}
				}
				log.Debug().
					Str("event_id", event.ID).
					Msg("No license plate data yet & out of attempts")
				return *event
			}
		}
	}
	return *event
}
