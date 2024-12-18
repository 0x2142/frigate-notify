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
	config.Internal.Status.LastEvent = time.Now()

	// Convert to human-readable timestamp
	reviewTime := time.Unix(int64(review.StartTime), 0)
	log.Info().
		Str("review_id", review.ID).
		Str("camera", review.Camera).
		Int("num_detections", len(review.Data.Detections)).
		Str("objects", strings.Join(review.Data.Objects, ",")).
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

	// Retrieve detailed detection information
	reviewFiltered := false
	var firstDetection models.Event
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

		// Check that event passes configured filters
		detection.CurrentZones = detection.Zones
		if !checkEventFilters(detection) {
			reviewFiltered = true
			break
		}

		// Store first detection for this review, alerts will be based on this event's data
		if firstDetection.ID == "" {
			firstDetection = detection
		}
	}
	// If any detection would be filtered, skip notifying on this review
	if reviewFiltered {
		log.Info().
			Str("review_id", review.ID).
			Msgf("Review dropped - One or more detections are filtered")
		return
	}

	// Add special link to review page
	firstDetection.Extra.ReviewLink = config.ConfigData.Frigate.PublicURL + "/review?id=" + review.ID

	// Send alert with snapshot
	notifier.SendAlert(firstDetection)
}
