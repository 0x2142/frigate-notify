package models

type FrigateStats struct {
	Service struct {
		LastUpdated   int    `json:"last_updated"`
		LatestVersion string `json:"latest_version"`
		Uptime        int    `json:"uptime"`
		Version       string `json:"version"`
	} `json:"service"`
}
