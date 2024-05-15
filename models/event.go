package models

// MQTTEvent stores incoming MQTT payloads from Frigate
type MQTTEvent struct {
	Before struct {
		Event
	} `json:"before,omitempty"`
	After struct {
		Event
	} `json:"after,omitempty"`
	Type string `json:"type"`
}

// Event stores Frigate alert attributes
type Event struct {
	Area   interface{} `json:"area"`
	Box    interface{} `json:"box"`
	Camera string      `json:"camera"`
	Data   struct {
		Attributes []interface{} `json:"attributes"`
		Box        []float64     `json:"box"`
		Region     []float64     `json:"region"`
		Score      float64       `json:"score"`
		TopScore   float64       `json:"top_score"`
		Type       string        `json:"type"`
	} `json:"data"`
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
	EnteredZones       []string    `json:"entered_zones"`
	Extra              ExtraFields
}

// Additional custom fields
type ExtraFields struct {
	FormattedTime   string
	TopScorePercent string
	ZoneList        string
	LocalURL        string
	PublicURL       string
}
