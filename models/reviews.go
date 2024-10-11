package models

type MQTTReview struct {
	Before struct {
		Review
	} `json:"before,omitempty"`
	After struct {
		Review
	} `json:"after,omitempty"`
	Type string `json:"type"`
}

type Review struct {
	ID        string  `json:"id"`
	Camera    string  `json:"camera"`
	StartTime float64 `json:"start_time"`
	EndTime   any     `json:"end_time"`
	Severity  string  `json:"severity"`
	ThumbPath string  `json:"thumb_path"`
	Data      struct {
		Detections []string `json:"detections"`
		Objects    []string `json:"objects"`
		SubLabels  []any    `json:"sub_labels"`
		Zones      []string `json:"zones"`
		Audio      []any    `json:"audio"`
	}
}
