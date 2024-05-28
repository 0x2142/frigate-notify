package frigate

import (
	"slices"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
)

// isAllowedZone verifies whether a zone should be allowed to generate a notification
func isAllowedZone(id string, zones []string) bool {
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

// isAllowedLabel verifies whether a label should be allowed to generate a notification
func isAllowedLabel(id string, label string) bool {
	// Check label block list
	if slices.Contains(config.ConfigData.Alerts.Labels.Block, label) {
		log.Info().
			Str("event_id", id).
			Str("label", label).
			Msg("Event dropped - Label block list.")
		return false
	}
	// If no allow list, all events are permitted
	if len(config.ConfigData.Alerts.Labels.Allow) == 0 {
		return true
	}
	// Check label allow list
	if slices.Contains(config.ConfigData.Alerts.Labels.Allow, label) {
		return true
	}

	// Default drop event
	log.Info().
		Str("event_id", id).
		Str("label", label).
		Msg("Event dropped - Not on label allow list.")
	return false
}
