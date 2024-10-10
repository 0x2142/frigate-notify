package frigate

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

var reviewURI = "/review?id="

// processReview handles incoming /reviews MQTT messages & pulls out relevant info for alerting
func processReview(client mqtt.Client, msg mqtt.Message) {
	// Parse incoming MQTT message
	var review models.Review
	json.Unmarshal(msg.Payload(), &review)

	log.Trace().
		RawJSON("payload", msg.Payload()).
		Msg("MQTT event received")

	if review.Type == "new" {
		// Convert to human-readable timestamp
		reviewTime := time.Unix(int64(review.After.StartTime), 0)
		log.Info().
			Str("review_id", review.After.ID).
			Str("camera", review.After.Camera).
			Int("num_detections", len(review.After.Data.Detections)).
			Str("objects", strings.Join(review.After.Data.Objects, ",")).
			Str("zones", strings.Join(review.After.Data.Zones, ",")).
			Str("severity", review.After.Severity).
			Msg("New review received")
		log.Debug().
			Str("review_id", review.After.ID).
			Msgf("Review start time: %s", reviewTime)

		if !config.ConfigData.Alerts.General.NotifyDetections && review.After.Severity == "detection" {
			log.Info().
				Str("review_id", review.After.ID).
				Msg("Review dropped - Event is detection only, not alert")
			return
		}

		// Retrieve detailed detection information
		reviewFiltered := false
		var firstDetection models.Event
		for _, id := range review.After.Data.Detections {
			url := fmt.Sprintf("%s%s/%s", config.ConfigData.Frigate.Server, eventsURI, id)

			response, err := util.HTTPGet(url, config.ConfigData.Frigate.Insecure, "")
			if err != nil {
				log.Error().
					Err(err).
					Str("review_id", review.After.ID).
					Str("detection_id", id).
					Msgf("Unable to retrieve detection information")
				continue
			}

			var detection models.Event
			json.Unmarshal(response, &detection)

			// Store first detection for this review
			if firstDetection.ID == "" {
				firstDetection = detection
			}
			// Check that event passes configured filters
			detection.CurrentZones = detection.Zones
			if !checkEventFilters(detection) {
				reviewFiltered = true
				break
			}
		}
		// If any detection would be filtered, skip notifying on this review
		if reviewFiltered {
			log.Info().
				Str("review_id", review.After.ID).
				Msgf("Review dropped - One or more detections are filtered")
			return
		}

		// Check if already notified on zones
		//if zoneAlreadyAlerted(event.After.Event) {
		//	log.Info().
		//		Str("event_id", event.After.ID).
		//		Str("camera", event.After.Camera).
		//		Str("label", event.After.Label).
		//		Str("zones", strings.Join(event.After.CurrentZones, ",")).
		//		Msg("Event dropped - Already notified on this zone")
		//	return
		//} else {
		//	log.Debug().
		//		Str("event_id", event.After.ID).
		//		Str("camera", event.After.Camera).
		//		Str("label", event.After.Label).
		//		Str("zones", strings.Join(event.After.CurrentZones, ",")).
		//		Msg("Object entered new zone")
		//}

		// Build snapshot & alert based on first detection
		// If snapshot was collected, pull down image to send with alert
		var snapshot io.Reader
		var snapshotURL string
		if firstDetection.HasSnapshot {
			snapshotURL = config.ConfigData.Frigate.Server + eventsURI + "/" + firstDetection.ID + snapshotURI
			snapshot = GetSnapshot(snapshotURL, firstDetection.ID)
		}

		// Set Review & Event links
		firstDetection.Extra.ReviewLink = config.ConfigData.Frigate.PublicURL + reviewURI + review.After.ID
		firstDetection.Extra.EventLink = config.ConfigData.Frigate.PublicURL + eventsURI + "/" + firstDetection.ID + "/clip.mp4"

		// Send alert with snapshot
		notifier.SendAlert(firstDetection, snapshot, firstDetection.ID)
	}

	// Clear event cache entry when event ends
	if review.Type == "end" {
		log.Debug().
			Str("review_id", review.After.ID).
			Msg("Review ended")
		return
	}
}
