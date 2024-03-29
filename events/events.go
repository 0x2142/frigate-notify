package frigate

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/config"
)

// Event stores Frigate alert attributes
type Event struct {
	Area               interface{} `json:"area"`
	Box                interface{} `json:"box"`
	Camera             string      `json:"camera"`
	EndTime            interface{} `json:"end_time"`
	FalsePositive      interface{} `json:"false_positive"`
	HasClip            bool        `json:"has_clip"`
	HasSnapshot        bool        `json:"has_snapshot"`
	ID                 string      `json:"id"`
	Label              string      `json:"label"`
	PlusID             interface{} `json:"plus_id"`
	Ratio              interface{} `json:"ratio"`
	Region             interface{} `json:"region"`
	RetainIndefinitely bool        `json:"retain_indefinitely"`
	StartTime          float64     `json:"start_time"`
	SubLabel           interface{} `json:"sub_label"`
	Thumbnail          string      `json:"thumbnail"`
	TopScore           float64     `json:"top_score"`
	Zones              []string    `json:"zones"`
	CurrentZones       []string    `json:"current_zones"`
}

// buildMessage constructs message payload for all alerting methods
func buildMessage(time time.Time, event Event) string {
	// Build alert message payload, include two spaces at end to force markdown newline
	message := fmt.Sprintf("Detection at %v  ", time)
	message += fmt.Sprintf("\nCamera: %s  ", event.Camera)
	// Attach detection label & caculate score percentage
	message += fmt.Sprintf("\nLabel: %v (%v%%)  ", event.Label, int((event.TopScore * 100)))
	// If zones configured / detected, include details
	var zones []string
	zones = append(zones, event.Zones...)
	zones = append(zones, event.CurrentZones...)
	if len(zones) >= 1 {
		message += fmt.Sprintf("\nZone(s): %v  ", strings.Join(zones, ", "))
	}
	message += fmt.Sprintf("\n\n[Link to Frigate](%s)", config.ConfigData.Frigate.Server)

	return message
}

// isAllowedZone verifies whether a zone should be allowed to generate a notification
func isAllowedZone(id string, zones []string) bool {
	// By default, send events without a zone unless specified otherwise
	if strings.ToLower(config.ConfigData.Alerts.Zones.Unzoned) == "drop" && len(zones) == 0 {
		log.Printf("Event ID %v - Dropped since it was outside a zone.", id)
		return false
	} else if len(zones) == 0 {
		return true
	}
	// Check zone block list
	for _, zone := range zones {
		if slices.Contains(config.ConfigData.Alerts.Zones.Block, zone) {
			log.Printf("Event ID %v - Dropped by zone block list.", id)
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
	log.Printf("Event ID %v - Dropped. Not on allow list.", id)
	return false
}
