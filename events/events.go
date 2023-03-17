package frigate

import (
	"fmt"
	"time"
)

// Event stores Frigate alert attributes
type Event struct {
	Area               interface{}   `json:"area"`
	Box                interface{}   `json:"box"`
	Camera             string        `json:"camera"`
	EndTime            interface{}   `json:"end_time"`
	FalsePositive      interface{}   `json:"false_positive"`
	HasClip            bool          `json:"has_clip"`
	HasSnapshot        bool          `json:"has_snapshot"`
	ID                 string        `json:"id"`
	Label              string        `json:"label"`
	PlusID             interface{}   `json:"plus_id"`
	Ratio              interface{}   `json:"ratio"`
	Region             interface{}   `json:"region"`
	RetainIndefinitely bool          `json:"retain_indefinitely"`
	StartTime          float64       `json:"start_time"`
	SubLabel           interface{}   `json:"sub_label"`
	Thumbnail          string        `json:"thumbnail"`
	TopScore           float64       `json:"top_score"`
	Zones              []interface{} `json:"zones"`
}

var FrigateServerURL string
var FrigateInsecure = false

// buildMessage constructs message payload for all alerting methods
func buildMessage(time time.Time, event Event) string {
	// Build alert message payload
	message := fmt.Sprintf("Detection at %v.", time)
	message += fmt.Sprintf("\nCamera: %s", event.Camera)
	// Attach detection label & caculate score percentage
	message += fmt.Sprintf("\nLabel: %v (%v%%)", event.Label, int((event.TopScore * 100)))
	// If zones configured / detected, include details
	if len(event.Zones) >= 1 {
		message += fmt.Sprintf("\nZone(s): %v", event.Zones)
	}
	message += fmt.Sprintf("\n\n[Link to Frigate](%s)", FrigateServerURL)

	return message
}
