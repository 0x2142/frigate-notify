package models

import "time"

// Store internal-use only info
type InternalConfig struct {
	AppVersion     string
	FrigateVersion int
	Status         Status
}

type Status struct {
	Health    string            `json:"health" example:"ok" doc:"Overall health of Frigate-Notify app"`
	Enabled   bool              `json:"enabled" example:"true" doc:"State of whether Frigate-Notify is enabled for alerting"`
	API       string            `json:"api" example:"v0.0.0" doc:"Health of Frigate-Notify API"`
	Frigate   FrigateConnection `json:"frigate"`
	LastEvent time.Time         `json:"last_event" example:"0001-01-01T00:00:00Z" doc:"Timestamp of last received event from Frigate"`
	LastAlert time.Time         `json:"last_alert" example:"0001-01-01T00:00:00Z" doc:"Timestamp of last sent alert"`
	Monitor   string            `json:"monitor" example:"ok" doc:"Health of reporting state to external health monitor app"`
}

type FrigateConnection struct {
	API  string `json:"api" example:"ok" doc:"Health of connection to Frigate via API"`
	MQTT string `json:"mqtt" example:"ok" doc:"Health of connection to Frigate via MQTT"`
}
