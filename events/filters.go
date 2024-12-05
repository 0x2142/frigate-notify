package events

import (
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

// checkFilters processes incoming event through configured filters to determine if it should generate a notification
func checkFilters(event models.Event) bool {
	// Check if notifications are currently disabled
	if !config.Internal.Status.Notifications.Enabled {
		log.Info().Msg("Event dropped - Notifications currently disabled.")
		return false
	}

	// Skip excluded cameras
	if slices.Contains(config.ConfigData.Frigate.Cameras.Exclude, event.Camera) {
		log.Info().
			Str("event_id", event.ID).
			Str("camera", event.Camera).
			Msg("Event dropped - Camera Excluded")
		return false
	}

	// Drop event if no snapshot or clip is available - Event is likely being filtered on Frigate side.
	// For example, if a camera has `required_zones` set - then there may not be any clip or snap until
	// object moves into required zone
	if !event.HasClip && !event.HasSnapshot {
		log.Info().
			Str("event_id", event.ID).
			Msg("Event dropped - No snapshot or clip available")
		return false
	}

	// Check if notify_once is set & we already notified on this event
	if config.ConfigData.Alerts.General.NotifyOnce {
		// Check if cache already contains event ID
		if getCachebyID(event.ID) != nil {
			log.Info().
				Str("event_id", event.ID).
				Msg("Event dropped - Already notified & notify_once is set")
			return false
		}
	}

	// Check if already notified on zones
	if zoneAlreadyAlerted(event) {
		log.Info().
			Str("event_id", event.ID).
			Str("camera", event.Camera).
			Str("label", event.Label).
			Str("zones", strings.Join(event.CurrentZones, ",")).
			Msg("Event dropped - Already notified on this zone")
		return false
	} else {
		log.Debug().
			Str("event_id", event.ID).
			Str("camera", event.Camera).
			Str("label", event.Label).
			Str("zones", strings.Join(event.CurrentZones, ",")).
			Msg("Object entered new zone")
	}

	// Drop event if no snapshot & skip_nosnap is true
	if !event.HasSnapshot && strings.ToLower(config.ConfigData.Alerts.General.NoSnap) == "drop" {
		log.Info().
			Str("event_id", event.ID).
			Msg("Event dropped - No snapshot available")
		return false
	}

	// Check quiet hours
	if isQuietHours() {
		log.Info().
			Str("event_id", event.ID).
			Msg("Event dropped - Quiet hours.")
		return false
	}

	// Check Zone filter
	if !isAllowedZone(event.ID, event.CurrentZones) {
		return false
	}

	// Check Label filter
	if !isAllowedLabel(event.ID, event.Label, "label") {
		return false
	}

	// Check label score
	if !aboveMinScore(event.ID, event.TopScore) {
		return false
	}

	// Check Sublabel filter
	if len(event.SubLabel) == 0 {
		if !isAllowedLabel(event.ID, "", "sublabel") {
			return false
		}
	} else {
		for _, sublabel := range event.SubLabel {
			if !isAllowedLabel(event.ID, sublabel, "sublabel") {
				return false
			}
		}
	}

	// Default
	return true
}

// isQuietHours checks to see if current event time is within window for supressing notifications
func isQuietHours() bool {
	currentTime, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
	start, _ := time.Parse("15:04", config.ConfigData.Alerts.Quiet.Start)
	end, _ := time.Parse("15:04", config.ConfigData.Alerts.Quiet.End)
	log.Trace().
		Time("current_time", currentTime).
		Time("quiet_start", start).
		Time("quiet_end", end).
		Msg("Check quiet hours")
	// Check if quiet period is overnight
	if end.Before(start) {
		if currentTime.After(start) || currentTime.Before(end) {
			return true
		}
	}
	// Otherwise check if between start & end times
	if currentTime.After(start) && currentTime.Before(end) {
		return true
	}
	return false
}

// isAllowedZone verifies whether a zone should be allowed to generate a notification
func isAllowedZone(id string, zones []string) bool {
	log.Trace().
		Str("event_id", id).
		Strs("zones", zones).
		Str("allow_unzoned", config.ConfigData.Alerts.Zones.Unzoned).
		Strs("blocked", config.ConfigData.Alerts.Zones.Block).
		Strs("allowed", config.ConfigData.Alerts.Zones.Allow).
		Msg("Check allowed zone")
	// By default, send events without a zone unless specified otherwise
	if strings.ToLower(config.ConfigData.Alerts.Zones.Unzoned) == "drop" && len(zones) == 0 {
		log.Info().
			Str("event_id", id).
			Str("zones", strings.Join(zones, ",")).
			Msg("Event dropped - Outside of zone.")
		return false
	} else if len(zones) == 0 {
		return true
	}
	// Check zone block list
	for _, zone := range zones {
		if slices.Contains(config.ConfigData.Alerts.Zones.Block, zone) {
			log.Info().
				Str("event_id", id).
				Str("zones", strings.Join(zones, ",")).
				Msg("Event dropped - Zone block list.")
			return false
		}
	}
	// If no allow list, all events are permitted
	if len(config.ConfigData.Alerts.Zones.Allow) == 0 {
		return true
	}
	// Check zone allow list
	for _, zone := range zones {
		if slices.Contains(config.ConfigData.Alerts.Zones.Allow, zone) {
			return true
		}
	}
	// Default drop event
	log.Info().
		Str("event_id", id).
		Str("zones", strings.Join(zones, ",")).
		Msg("Event dropped - Not on zone allow list.")
	return false
}

// isAllowedLabel verifies whether a label or sublabel should be allowed to generate a notification
func isAllowedLabel(id string, label string, kind string) bool {
	var blocked []string
	var allowed []string
	if kind == "label" {
		blocked = config.ConfigData.Alerts.Labels.Block
		allowed = config.ConfigData.Alerts.Labels.Allow
		log.Trace().
			Str("event_id", id).
			Str("label", label).
			Strs("blocked", blocked).
			Strs("allowed", allowed).
			Msg("Check allowed label")
	}
	if kind == "sublabel" {
		blocked = config.ConfigData.Alerts.SubLabels.Block
		allowed = config.ConfigData.Alerts.SubLabels.Allow
		log.Trace().
			Str("event_id", id).
			Str("label", label).
			Strs("blocked", blocked).
			Strs("allowed", allowed).
			Msg("Check allowed sublabel")
	}
	// Check block list
	if slices.Contains(blocked, label) {
		log.Info().
			Str("event_id", id).
			Str(kind, label).
			Msgf("Event dropped - %s block list.", kind)
		return false
	}
	// If no allow list, all events are permitted
	if len(allowed) == 0 {
		return true
	}
	// Check allow list
	if slices.Contains(allowed, label) {
		return true
	}

	// Default drop event
	log.Info().
		Str("event_id", id).
		Str(kind, label).
		Msgf("Event dropped - Not on %s allow list.", kind)
	return false
}

// aboveMinScore checks if label score is above configured minimum
func aboveMinScore(id string, score float64) bool {
	minScore := config.ConfigData.Alerts.Labels.MinScore
	score = score * 100
	log.Trace().
		Str("event_id", id).
		Float64("event_score", score).
		Float64("min_score", minScore).
		Msg("Check minimum score")
	if score >= minScore {
		return true
	} else {
		log.Info().
			Str("event_id", id).
			Float64("score", score).
			Msg("Event dropped - Does not meet minimum label score.")
		return false
	}
}
