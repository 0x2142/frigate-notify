package models

import "time"

// Store internal-use only info
type InternalConfig struct {
	AppVersion     string
	FrigateVersion int
	Status         Status
}

type Status struct {
	Health           string            `json:"health" example:"ok" doc:"Overall health of Frigate-Notify app"`
	API              string            `json:"api" example:"v0.0.0" doc:"Health of Frigate-Notify API"`
	Frigate          FrigateConnection `json:"frigate"`
	LastEvent        time.Time         `json:"last_event" example:"0001-01-01T00:00:00Z" doc:"Timestamp of last received event from Frigate"`
	LastNotification time.Time         `json:"last_notification" example:"0001-01-01T00:00:00Z" doc:"Timestamp of last sent notification"`
	Notifications    Notifiers         `json:"notifications" doc:"Status of notification providers"`
	Monitor          string            `json:"monitor" example:"ok" doc:"Health of reporting state to external health monitor app"`
}

type FrigateConnection struct {
	API  string `json:"api" example:"ok" doc:"Health of connection to Frigate via API"`
	MQTT string `json:"mqtt" example:"ok" doc:"Health of connection to Frigate via MQTT"`
}

type Notifiers struct {
	Enabled    bool             `json:"enabled" example:"true" doc:"State of whether Frigate-Notify is enabled for notifications"`
	AppriseAPI []NotifierStatus `json:"apprise-api" doc:"Status of Apprise API notifications"`
	Discord    []NotifierStatus `json:"discord" doc:"Status of Discord notifications"`
	Gotify     []NotifierStatus `json:"gotify" doc:"Status of Gotify notifications"`
	Ntfy       []NotifierStatus `json:"ntfy" doc:"Status of Ntfy notifications"`
	Pushover   []NotifierStatus `json:"pushover" doc:"Status of Pushover notifications"`
	Signal     []NotifierStatus `json:"signal" doc:"Status of Signal notifications"`
	SMTP       []NotifierStatus `json:"smtp" doc:"Status of SMTP notifications"`
	Telegram   []NotifierStatus `json:"telegram" doc:"Status of Telegram notifications"`
	Webhook    []NotifierStatus `json:"webhook" doc:"Status of Webhook notifications"`
	Mattermost []NotifierStatus `json:"mattermost" doc:"Status of Mattermost notifications"`
}

type NotifierStatus struct {
	ID          int       `json:"id"`
	Enabled     bool      `json:"enabled" doc:"Whether notification provider is enabled"`
	Status      string    `json:"status" default:"not enabled"`
	Sent        int64     `json:"sent" doc:"Number of notifications sent via this provider"`
	Failed      int64     `json:"failed" doc:"Number of errors while attempting to send notifications via this provider"`
	LastSuccess time.Time `json:"last_success" doc:"Timestamp of last successful notification sent"`
	LastFailure time.Time `json:"last_failure" doc:"Timestamp of last failure to send notification"`
	LastError   string    `json:"last_error" doc:"Error message from last failure, if applicable" default:"n/a"`
}

func (n *NotifierStatus) InitNotifStatus(id int, enabled bool) {
	n.ID = id
	n.Enabled = enabled
	n.Status = "configured, not used yet"
}

func (n *NotifierStatus) NotifSuccess() {
	n.LastSuccess = time.Now()
	n.Sent += 1
	n.Status = "ok"
}

func (n *NotifierStatus) NotifFailure(message string) {
	n.LastFailure = time.Now()
	n.LastError = message
	n.Failed += 1
	n.Status = "error"
}
