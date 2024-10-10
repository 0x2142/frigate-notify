package models

type Review struct {
	Type   string `json:"type"`
	Before struct {
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
		} `json:"data"`
	} `json:"before"`
	After struct {
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
		} `json:"data"`
	} `json:"after"`
}
