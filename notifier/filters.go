package notifier

import (
	"slices"
	"time"

	"github.com/0x2142/frigate-notify/models"
	"github.com/rs/zerolog/log"
)

// checkAlertFilters will determine which notification provider is able to send this alert
func checkAlertFilters(event models.Event, filters models.AlertFilter, provider notifMeta) bool {
	log.Trace().
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Msg("Checking alert filters")

	// Check against quiet hours
	currentTime, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
	start, _ := time.Parse("15:04", filters.Quiet.Start)
	end, _ := time.Parse("15:04", filters.Quiet.End)
	log.Trace().
		Time("current_time", currentTime).
		Time("quiet_start", start).
		Time("quiet_end", end).
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Msg("Check quiet hours")
	// Check if quiet period is overnight
	if end.Before(start) {
		if currentTime.After(start) || currentTime.Before(end) {
			log.Debug().
				Str("provider", provider.name).
				Int("provider_id", provider.index).
				Msg("Notification droppped - Quiet hours")
			return false
		}
	}
	// Otherwise check if between start & end times
	if currentTime.After(start) && currentTime.Before(end) {
		log.Debug().
			Str("provider", provider.name).
			Int("provider_id", provider.index).
			Msg("Notification droppped - Quiet hours")
		return false
	}

	// Check filtered cameras
	log.Trace().
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Str("camera", event.Camera).
		Strs("allowed", filters.Cameras).
		Msg("Check allowed cameras")
	if len(filters.Cameras) >= 1 {
		if !slices.Contains(filters.Cameras, event.Camera) {
			log.Debug().
				Str("provider", provider.name).
				Int("provider_id", provider.index).
				Msg("Notification droppped - Camera not on filter list")
			return false
		}
	}

	// Check filtered zones
	log.Trace().
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Strs("zones", event.CurrentZones).
		Strs("allowed", filters.Zones).
		Msg("Check allowed zone")
	if len(filters.Zones) >= 1 {
		matchzone := false
		for _, zone := range event.CurrentZones {
			if slices.Contains(filters.Zones, zone) {
				matchzone = true
			}
		}
		if !matchzone {
			log.Debug().
				Str("provider", provider.name).
				Int("provider_id", provider.index).
				Msg("Notification droppped - Zone not on filter list")
			return false
		}
	}

	// Check filtered Labels
	log.Trace().
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Str("label", event.Label).
		Strs("allowed", filters.Labels).
		Msg("Check allowed label")
	if len(filters.Labels) >= 1 {
		if !slices.Contains(filters.Labels, event.Label) {
			log.Debug().
				Str("provider", provider.name).
				Int("provider_id", provider.index).
				Msg("Notification droppped - Label not on filter list")
			return false
		}
	}

	// Check filtered Sublabels
	log.Trace().
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Strs("label", event.SubLabel).
		Strs("allowed", filters.Sublabels).
		Msg("Check allowed sublabel")
	if len(filters.Sublabels) >= 1 {
		matchsublabel := false
		for _, sublabel := range event.SubLabel {
			if slices.Contains(filters.Sublabels, sublabel) {
				matchsublabel = true
			}
		}
		if !matchsublabel {
			log.Debug().
				Str("provider", provider.name).
				Int("provider_id", provider.index).
				Msg("Notification droppped - Sublabel not on filter list")
			return false
		}
	}

	// Alert permitted if all conditions pass
	log.Trace().
		Str("provider", provider.name).
		Int("provider_id", provider.index).
		Msg("Alert filters passed!")
	return true
}
