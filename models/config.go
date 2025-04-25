package models

type App struct {
	Mode string `koanf:"mode" json:"mode" enum:"events,reviews" doc:"Type of polling method used when connecting to Frigate" default:"reviews"`
	API  API    `koanf:"api" json:"api" doc:"Frigate-Notify API settings"`
}

type API struct {
	Enabled bool `koanf:"enabled" json:"enabled" doc:"Enable Frigate-Notify API server" enum:"true,false" default:"false"`
	Port    int  `koanf:"port" json:"port,omitempty" doc:"API server port" minimum:"1" maximum:"65535" default:"8000"`
}

type Frigate struct {
	Server       string              `koanf:"server" json:"server" validate:"required" example:"http://192.0.2.10:5000" doc:"Server hostname, IP address, or URL for Frigate"`
	Insecure     bool                `koanf:"ignoressl" json:"ignoressl,omitempty" enum:"true,false" doc:"Ignore TLS/SSL errors" default:"false"`
	PublicURL    string              `koanf:"public_url" json:"public_url,omitempty" example:"https://frigate.test" doc:"Public/External-reachable URL for Frigate" default:""`
	Headers      []map[string]string `koanf:"headers" json:"headers,omitempty" doc:"HTTP headers to include with requests to Frigate"`
	StartupCheck StartupCheck        `koanf:"startup_check" json:"startup_check,omitempty" doc:"Check connectivity to Frigate at app startup"`
	WebAPI       WebAPI              `koanf:"webapi" json:"webapi,omitempty" doc:"Event collection via Frigate API"`
	MQTT         MQTT                `koanf:"mqtt" json:"mqtt,omitempty" doc:"Event collection via MQTT`
	Cameras      Cameras             `koanf:"cameras" json:"cameras,omitempty" doc:"Camera settings"`
}

type StartupCheck struct {
	Attempts int `koanf:"attempts" json:"attempts,omitempty" doc:"Maximum attempts for connecting to Frigate" minimum:"1" maximum:"10000000" default:"5"`
	Interval int `koanf:"interval" json:"interval,omitempty" doc:"Interval for connection attempts" minimum:"1" maximum:"10000000" default:"30"`
}

type WebAPI struct {
	Enabled  bool `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable event collection via Frigate API" default:"false"`
	Interval int  `koanf:"interval" json:"interval,omitempty" doc:"Interval of API event collection from Frigate" minimum:"1" maximum:"65535" default:"30"`
	TestMode bool `koanf:"testmode" json:"testmode,omitempty" enum:"true,false" doc:"Used for testing only" hidden:"true" default:"false"`
}

type MQTT struct {
	Enabled     bool   `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable event collection via MQTT" default:"false"`
	Server      string `koanf:"server" json:"server,omitempty" doc:"MQTT server address" default:""`
	Port        int    `koanf:"port" json:"port,omitempty" doc:"MQTT port" minimum:"1" maximum:"65535" default:"1883"`
	ClientID    string `koanf:"clientid" json:"clientid,omitempty" doc:"MQTT client ID" default:"frigate-notify"`
	Username    string `koanf:"username" json:"username,omitempty" doc:"MQTT username" default:""`
	Password    string `koanf:"password" json:"password,omitempty" doc:"MQTT password" default:""`
	TopicPrefix string `koanf:"topic_prefix" json:"topic_prefix,omitempty" doc:"MQTT topic prefix" default:"frigate"`
}

type Cameras struct {
	Exclude []string `koanf:"exclude" json:"exclude" uniqueItems:"true" doc:"Exclude cameras from alerting"`
}

type Alerts struct {
	General      General      `koanf:"general" json:"general,omitempty" doc:"Common alert settings"`
	Quiet        Quiet        `koanf:"quiet" json:"quiet,omitempty" doc:"Alert quiet periods"`
	Zones        Zones        `koanf:"zones" json:"zones,omitempty" doc:"Allow/Block zones from alerting"`
	Labels       Labels       `koanf:"labels" json:"labels,omitempty" doc:"Allow/Block labels from alerting"`
	SubLabels    Labels       `koanf:"sublabels" json:"sublabels,omitempty" doc:"Allow/Block sublabels from alerting"`
	LicensePlate LicensePlate `koanf:"license_plate" json:"license_plate,omitempty" doc:"License plate recognition settings"`
	AppriseAPI   []AppriseAPI `koanf:"apprise_api" json:"apprise_api,omitempty" doc:"Apprise API notification settings"`
	Discord      []Discord    `koanf:"discord" json:"discord,omitempty" doc:"Discord notification settings"`
	Gotify       []Gotify     `koanf:"gotify" json:"gotify,omitempty" doc:"Gotify notification settings"`
	Matrix       []Matrix     `koanf:"matrix" json:"matrix,omitempty" doc:"Matrix notification settings"`
	Mattermost   []Mattermost `koanf:"mattermost" json:"mattermost,omitempty" doc:"Mattermost notification settings"`
	Ntfy         []Ntfy       `koanf:"ntfy" json:"ntfy,omitempty" doc:"Ntfy notification settings"`
	Pushover     []Pushover   `koanf:"pushover" json:"pushover,omitempty" doc:"Pushover notification settings"`
	Signal       []Signal     `koanf:"signal" json:"signal,omitempty" doc:"Signal notification settings"`
	SMTP         []SMTP       `koanf:"smtp" json:"smtp,omitempty" doc:"SMTP notification settings"`
	Telegram     []Telegram   `koanf:"telegram" json:"telegram,omitempty" doc:"Telegram notification settings"`
	Webhook      []Webhook    `koanf:"webhook" json:"webhook,omitempty" doc:"Webhook notification settings"`
}

type General struct {
	Title            string `koanf:"title" json:"title,omitempty" doc:"Notification title" default:"Frigate Alert"`
	TimeFormat       string `koanf:"timeformat" json:"timeformat,omitempty" doc:"Time format used in notifications" default:""`
	NoSnap           string `koanf:"nosnap,omitempty" json:"nosnap" enum:"allow,drop" doc:"Allow/Drop events if they do not have a snapshot" default:"allow"`
	SnapBbox         bool   `koanf:"snap_bbox,omitempty" json:"snap_bbox" enum:"true,false" doc:"Include bounding box on snapshots" default:"false"`
	SnapTimestamp    bool   `koanf:"snap_timestamp,omitempty" json:"snap_timestamp" enum:"true,false" doc:"Include timestamp on snapshots" default:"false"`
	SnapCrop         bool   `koanf:"snap_crop,omitempty"  json:"snap_crop" enum:"true,false" doc:"Crop snapshots" default:"false"`
	MaxSnapRetry     int    `koanf:"max_snap_retry,omitempty" json:"max_snap_retry" doc:"Maximum number of retry attempts when snapshot is not ready yet" default:"10"`
	NotifyOnce       bool   `koanf:"notify_once,omitempty"  json:"notify_once" enum:"true,false" doc:"Only notify once per event (For app mode: events)" default:"false"`
	NotifyDetections bool   `koanf:"notify_detections,omitempty" json:"notify_detections" enum:"true,false" doc:"Enable notifications on detection (For app mode: reviews)" default:"false"`
	RecheckDelay     int    `koanf:"recheck_delay" json:"recheck_delay,omitempty" default:"0" doc:"Delay before re-checking event details from Frigate"`
	AudioOnly        string `koanf:"audio_only" json:"audio_only,omitempty" enum:"allow,drop" doc:"Allow/Drop events that only contain audio detections" default:"allow"`
}

type LicensePlate struct {
	Enabled bool     `koanf:"enabled" json:"enabled,omitempty" enum:"true,false" doc:"Enable waiting for license plate recognition when car & license plate are detected" default:"false"`
	Allow   []string `koanf:"allow" json:"allow,omitempty" doc:"List of license plates to allow alerts from"`
	Block   []string `koanf:"block" json:"block,omitempty" doc:"List of license plates to always block"`
}

type Quiet struct {
	Start string `koanf:"start" json:"start,omitempty" example:"02:30" pattern:"(\d)?\d:\d\d" doc:"Start time for quiet hours" default:""`
	End   string `koanf:"end" json:"end,omitempty" example:"05:45" pattern:"(\d)?\d:\d\d" doc:"End time for quiet hours" default:""`
}

type Zones struct {
	Unzoned string   `koanf:"unzoned" json:"unzoned,omitempty" enum:"allow,drop" doc:"Allow/Drop events when object is outside a zone" default:"allow"`
	Allow   []string `koanf:"allow" json:"allow,omitempty" doc:"List of zones to allow alerts from"`
	Block   []string `koanf:"block" json:"block,omitempty" doc:"List of zones to always block"`
}

type Labels struct {
	MinScore float64  `koanf:"min_score" json:"min_score" minimum:"0" maximum:"100" doc:"Set minimum score before events will notify" default:"0"`
	Allow    []string `koanf:"allow" json:"allow,omitempty" doc:"List of labels to allow alerts from" `
	Block    []string `koanf:"block" json:"block,omitempty" doc:"List of labels to always block"`
}

type AlertFilter struct {
	Cameras   []string `koanf:"cameras" json:"cameras,omitempty" doc:"List of cameras that will use this alert provider`
	Zones     []string `koanf:"zones" json:"zones,omitempty" doc:"List of zones that will use this alert provider`
	Quiet     Quiet    `koanf:"quiet" json:"quiet,omitempty" doc:"Quiet period for this alert provider"`
	Labels    []string `koanf:"labels" json:"labels,omitempty" doc:"List of labels that will use this alert provider"`
	Sublabels []string `koanf:"sublabels" json:"sublabels,omitempty" doc:"List of sublabels that will use this alert provider"`
}

type AppriseAPI struct {
	Enabled  bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Apprise API" default:"false"`
	Server   string      `koanf:"server" json:"server,omitempty" doc:"Apprise API URL to send alerts" default:""`
	Token    string      `koanf:"token" json:"token,omitempty" doc:"Apprise API config token"`
	Tags     []string    `koanf:"tags" json:"tags,omitempty" doc:"Notification group tags to receive alert"`
	URLs     []string    `koanf:"urls" json:"urls,omitempty" doc:"Apprise notification URLs to send alerts to"`
	Insecure bool        `koanf:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:"false"`
	Template string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Discord struct {
	Enabled      bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Discord" default:"false"`
	Webhook      string      `koanf:"webhook" json:"webhook,omitempty" doc:"Discord webhook URL to send alerts" default:""`
	DisableEmbed bool        `koanf:"disable_embed" json:"disable_embed,omitempty" doc:"Disable sending notification as Discord embedded message"`
	Template     string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters      AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Gotify struct {
	Enabled  bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Gotify" default:"false"`
	Server   string      `koanf:"server" json:"server,omitempty" doc:"Gotify server URL" default:""`
	Token    string      `koanf:"token" json:"token,omitempty" doc:"Gotify app token" default:""`
	Insecure bool        `koanf:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:"false"`
	Template string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Matrix struct {
	Enabled  bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Matrix" default:"false"`
	Server   string      `koanf:"server" json:"server,omitempty" doc:"Matrix Homeserver address" default:""`
	Username string      `koanf:"username" json:"username,omitempty" doc:"Matrix username" default:""`
	Password string      `koanf:"password" json:"password,omitempty" doc:"Matrix password" default:""`
	RoomID   string      `koanf:"roomid" json:"roomid,omitempty" doc:"Room ID to send notifications to" default:""`
	Insecure bool        `koanf:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:"false"`
	Template string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Mattermost struct {
	Enabled  bool                `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Mattermost" default:"false"`
	Webhook  string              `koanf:"webhook" json:"webhook,omitempty" doc:"Mattermost webhook URL" default:""`
	Channel  string              `koanf:"channel" json:"channel,omitempty" doc:"Mattermost channel" default:""`
	Priority string              `koanf:"priority" json:"priority,omitempty" enum:"standard,important,urgent" doc:"Mattermost message priority" default:"standard"`
	Username string              `koanf:"username" json:"username,omitempty" doc:"Override Mattermost username for messages"`
	Insecure bool                `koanf:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:"false"`
	Headers  []map[string]string `koanf:"headers" json:"headers,omitempty" doc:"HTTP headers to include with Mattermost notifications" `
	Template string              `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter         `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Ntfy struct {
	Enabled  bool                `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Ntfy" default:"false"`
	Server   string              `koanf:"server" json:"server,omitempty" doc:"Ntfy Server address" default:""`
	Topic    string              `koanf:"topic" json:"topic,omitempty" doc:"Ntfy topic" default:""`
	Insecure bool                `koanf:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:"false"`
	Headers  []map[string]string `koanf:"headers" json:"headers,omitempty" doc:"HTTP headers to include with Ntfy notifications"`
	Template string              `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter         `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Pushover struct {
	Enabled  bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Pushover" default:"false"`
	Token    string      `koanf:"token" json:"token,omitempty" doc:"Pushover app token" default:""`
	Userkey  string      `koanf:"userkey" json:"userkey,omitempty" doc:"Pushover user key" default:""`
	Devices  string      `koanf:"devices" json:"devices,omitempty" doc:"Pushover devices to target for notification" default:""`
	Sound    string      `koanf:"sound" json:"sound,omitempty" doc:"Pushover notification sound" default:"pushover"`
	Priority int         `koanf:"priority" json:"priority,omitempty" minimum:"-2" maximum:"2" doc:"Pushover message priority" default:"0"`
	Retry    int         `koanf:"retry" json:"retry,omitempty" doc:"Retry interval for emergency notifications (Priority 2)" default:"0"`
	Expire   int         `koanf:"expire" json:"expire,omitempty" doc:"Expiration timer for emergency notifications (Priority 2)" default:"0"`
	TTL      int         `koanf:"ttl" json:"ttl,omitempty" doc:"Time to Live for notification messages" default:"0"`
	Template string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Signal struct {
	Enabled    bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Signal" default:"false"`
	Server     string      `koanf:"server" json:"server,omitempty" doc:"Signal REST API server hostname or IP address" default:""`
	Account    string      `koanf:"account" json:"account,omitempty" doc:"Number of account used to send messages" default:""`
	Recipients []string    `koanf:"recipients" json:"recipients,omitempty" doc:"List of recipients to receive messages"`
	Insecure   bool        `koanf:"ignoressl" enum:"true,false" json:"ignoressl,omitempty" default:"false"`
	Template   string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters    AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type SMTP struct {
	Enabled   bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via SMTP" default:"false"`
	Server    string      `koanf:"server" json:"server,omitempty" doc:"SMTP server hostname or IP address" default:""`
	Port      int         `koanf:"port" json:"port,omitempty" minimum:"1" maximum:"65535" doc:"SMTP server port" default:25`
	TLS       bool        `koanf:"tls" json:"tls,omitempty" enum:"true,false" doc:"Enable/Disable TLS connection" default:"false"`
	User      string      `koanf:"user" json:"user,omitempty" doc:"SMTP user for authentication" default:""`
	Password  string      `koanf:"password" json:"password,omitempty" doc:"SMTP password for authentication" default:""`
	From      string      `koanf:"from" json:"from,omitempty" format:"email" doc:"SMTP sender" default:""`
	Recipient string      `koanf:"recipient" json:"recipient,omitempty" format:"email" doc:"SMTP recipient" default:""`
	Template  string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Insecure  bool        `koanf:"ignoressl" enum:"true,false" json:"ignoressl,omitempty" default:"false"`
	Filters   AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Telegram struct {
	Enabled         bool        `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Telegram" default:"false"`
	ChatID          int64       `koanf:"chatid" json:"chatid,omitempty" minimum:"1" doc:"Telegram chat ID" default:"0"`
	MessageThreadID int         `koanf:"message_thread_id" json:"message_thread_id,omitempty" doc:"Send message to thread by ID" default:"0`
	Token           string      `koanf:"token" json:"token,omitempty" doc:"Telegram bot token" default:""`
	Template        string      `koanf:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters         AlertFilter `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Webhook struct {
	Enabled  bool                `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Webhook" default:"false"`
	Server   string              `koanf:"server" json:"server,omitempty" doc:"Webhook address" default:""`
	Insecure bool                `koanf:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:"false"`
	Method   string              `koanf:"method" json:"method,omitempty" enum:"GET,POST" doc:"HTTP method for webhook notifications" default:"POST"`
	Params   []map[string]string `fix:"params" json:"params,omitempty" doc:"URL parameters for webhook notifications"`
	Headers  []map[string]string `koanf:"headers" json:"headers,omitempty" doc:"HTTP headers for webhook notifications"`
	Template interface{}         `koanf:"template" json:"template,omitempty" doc:"Custom message template"`
	Filters  AlertFilter         `koanf:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Monitor struct {
	Enabled  bool   `koanf:"enabled" json:"enabled" enum:"true,false" doc:"Enable monitoring via external uptime application" default:"false"`
	URL      string `koanf:"url" json:"url,omitempty" doc:"Address of monitoring server" default:""`
	Interval int    `koanf:"interval" json:"interval,omitempty" minimum:"1" maximum:"10000000" doc:"Interval between check-in messages" default:"60"`
	Insecure bool   `koanf:"ignoressl" json:"ignoressl,omitempty" enum:"true,false" doc:"Ignore TLS/SSL errors" default:"false"`
}
