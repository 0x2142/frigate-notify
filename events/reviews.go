package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
	"github.com/rs/zerolog/log"
)

// processReview handles querying detections under a review & preparing for sending an alert
func processReview(review models.Review) {
	if config.ConfigData.Alerts.General.RecheckDelay != 0 {
		review = recheckReview(review)
	}

	config.Internal.Status.LastEvent = time.Now()

	// Convert to human-readable timestamp
	reviewTime := time.Unix(int64(review.StartTime), 0)
	log.Info().
		Str("review_id", review.ID).
		Str("camera", review.Camera).
		Int("num_detections", len(review.Data.Detections)).
		Str("objects", strings.Join(review.Data.Objects, ",")).
		Str("audio", strings.Join(review.Data.Audio, ",")).
		Str("zones", strings.Join(review.Data.Zones, ",")).
		Str("severity", review.Severity).
		Msg("Processing review...")
	log.Debug().
		Str("review_id", review.ID).
		Msgf("Review start time: %s", reviewTime)

	if !config.ConfigData.Alerts.General.NotifyDetections && review.Severity == "detection" {
		log.Info().
			Str("review_id", review.ID).
			Msg("Review dropped - Event is detection only, not alert")
		return
	}

	// Check if audio-only event
	if len(review.Data.Detections) == 0 && len(review.Data.Audio) != 0 {
		if config.ConfigData.Alerts.General.AudioOnly == "allow" {
			// Assemble some info via Review item, since there is no detection event to look up
			var audioEvent models.Event
			audioEvent.StartTime = review.StartTime
			audioEvent.Extra.Audio = strings.Join(review.Data.Audio, ",")
			audioEvent.Camera = review.Camera
			audioEvent.Extra.ReviewLink = config.ConfigData.Frigate.PublicURL + "/review?id=" + review.ID
			notifier.SendAlert([]models.Event{audioEvent})
			return
		} else {
			log.Info().
				Str("review_id", review.ID).
				Msg("Review dropped - Audio only event")
			return
		}
	}

	// Retrieve detailed detection information
	reviewFiltered := false
	var detections []models.Event
	for _, id := range review.Data.Detections {
		url := fmt.Sprintf("%s/api/events/%s", config.ConfigData.Frigate.Server, id)

		response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "")
		if err != nil {
			config.Internal.Status.Frigate.API = "unreachable"
			log.Error().
				Err(err).
				Str("review_id", review.ID).
				Str("detection_id", id).
				Msgf("Unable to retrieve detection information")
			continue
		}
		config.Internal.Status.Frigate.API = "ok"

		var detection models.Event
		json.Unmarshal(response, &detection)

		// For events collected via API, top-level top_score value is no longer used
		// So need to replace it with data.top_score value
		if detection.TopScore == 0 {
			detection.TopScore = detection.Data.TopScore
		}

		// Wait for license plate data before notifying, if set
		if config.ConfigData.Alerts.LicensePlate.Enabled {
			waitforLPR(&detection)
		}

		// Check that event passes configured filters
		detection.CurrentZones = detection.Zones
		if !checkEventFilters(detection) {
			reviewFiltered = true
			break
		}

		// Add special link to review page
		detection.Extra.ReviewLink = config.ConfigData.Frigate.PublicURL + "/review?id=" + review.ID

		detections = append(detections, detection)
	}

	// Check to make sure at least 1 detection passed filters
	if len(detections) == 0 {
		log.Info().
			Str("review_id", review.ID).
			Msgf("Review dropped - No events eligable for notification")
		return
	}

	// If any detection would be filtered, skip notifying on this review
	if reviewFiltered {
		log.Info().
			Str("review_id", review.ID).
			Msgf("Review dropped - One or more detections are filtered")
		return
	}

	// Send alert with snapshot
	notifier.SendAlert(detections)
}

func recheckReview(review models.Review) models.Review {
	delay := config.ConfigData.Alerts.General.RecheckDelay
	log.Debug().
		Str("review_id", review.ID).
		Int("recheck_delay", delay).
		Msg("Waiting to re-check review details")
	time.Sleep(time.Duration(delay) * time.Second)
	log.Debug().
		Str("review_id", review.ID).
		Int("recheck_delay", delay).
		Msg("Re-checking review details")

	url := config.ConfigData.Frigate.Server + "/api/review/" + review.ID
	response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "", config.ConfigData.Frigate.Headers...)
	if err != nil {
		config.Internal.Status.Health = "frigate webapi unreachable"
		config.Internal.Status.Frigate.API = "unreachable"
		log.Error().
			Err(err).
			Msgf("Cannot get event from %s", url)
		return review
	}
	config.Internal.Status.Health = "ok"
	config.Internal.Status.Frigate.API = "ok"

	json.Unmarshal([]byte(response), &review)
	return review
}
