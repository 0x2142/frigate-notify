package models

type App struct {
	Mode string `fig:"mode" json:"mode" enum:"events,reviews" doc:"Type of polling method used when connecting to Frigate" default:"reviews"`
	API  API    `fig:"api" json:"api" doc:"Frigate-Notify API settings"`
}

type API struct {
	Enabled bool `fig:"enabled" json:"enabled" doc:"Enable Frigate-Notify API server" enum:"true,false" default:false`
	Port    int  `fig:"port" json:"port,omitempty" doc:"API server port" minimum:"1" maximum:"65535" default:"8000"`
}

type Frigate struct {
	Server       string              `fig:"server" json:"server" validate:"required" example:"http://192.0.2.10:5000" doc:"Server hostname, IP address, or URL for Frigate"`
	Insecure     bool                `fig:"ignoressl" json:"ignoressl,omitempty" enum:"true,false" doc:"Ignore TLS/SSL errors" default:false`
	PublicURL    string              `fig:"public_url" json:"public_url,omitempty" example:"https://frigate.test" doc:"Public/External-reachable URL for Frigate" default:""`
	Headers      []map[string]string `fig:"headers" json:"headers,omitempty" doc:"HTTP headers to include with requests to Frigate"`
	StartupCheck StartupCheck        `fig:"startup_check" json:"startup_check,omitempty" doc:"Check connectivity to Frigate at app startup"`
	WebAPI       WebAPI              `fig:"webapi" json:"webapi,omitempty" doc:"Event collection via Frigate API"`
	MQTT         MQTT                `fig:"mqtt" json:"mqtt,omitempty" doc:"Event collection via MQTT`
	Cameras      Cameras             `fig:"cameras" json:"cameras,omitempty" doc:"Camera settings"`
}

type StartupCheck struct {
	Attempts int `fig:"attempts" json:"attempts,omitempty" doc:"Maximum attempts for connecting to Frigate" minimum:"1" maximum:"10000000" default:"5"`
	Interval int `fig:"interval" json:"interval,omitempty" doc:"Interval for connection attempts" minimum:"1" maximum:"10000000" default:"30"`
}

type WebAPI struct {
	Enabled  bool `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable event collection via Frigate API" default:false`
	Interval int  `fig:"interval" json:"interval,omitempty" doc:"Interval of API event collection from Frigate" minimum:"1" maximum:"65535" default:"30"`
	TestMode bool `fig:"testmode" json:"testmode,omitempty" enum:"true,false" doc:"Used for testing only" hidden:"true" default:false`
}

type MQTT struct {
	Enabled     bool   `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable event collection via MQTT" default:false`
	Server      string `fig:"server" json:"server,omitempty" doc:"MQTT server address" default:""`
	Port        int    `fig:"port" json:"port,omitempty" doc:"MQTT port" minimum:"1" maximum:"65535" default:"1883"`
	ClientID    string `fig:"clientid" json:"clientid,omitempty" doc:"MQTT client ID" default:"frigate-notify"`
	Username    string `fig:"username" json:"username,omitempty" doc:"MQTT username" default:""`
	Password    string `fig:"password" json:"password,omitempty" doc:"MQTT password" default:""`
	TopicPrefix string `fig:"topic_prefix" json:"topic_prefix,omitempty" doc:"MQTT topic prefix" default:"frigate"`
}

type Cameras struct {
	Exclude []string `fig:"exclude" json:"exclude" uniqueItems:"true" doc:"Exclude cameras from alerting" default:[]`
}

type Alerts struct {
	General   General    `fig:"general" json:"general,omitempty" doc:"Common alert settings"`
	Quiet     Quiet      `fig:"quiet" json:"quiet,omitempty" doc:"Alert quiet periods"`
	Zones     Zones      `fig:"zones" json:"zones,omitempty" doc:"Allow/Block zones from alerting"`
	Labels    Labels     `fig:"labels" json:"labels,omitempty" doc:"Allow/Block labels from alerting"`
	SubLabels Labels     `fig:"sublabels" json:"sublabels,omitempty" doc:"Allow/Block sublabels from alerting"`
	Discord   []Discord  `fig:"discord" json:"discord,omitempty" doc:"Discord notification settings"`
	Gotify    []Gotify   `fig:"gotify" json:"gotify,omitempty" doc:"Gotify notification settings"`
	SMTP      []SMTP     `fig:"smtp" json:"smtp,omitempty" doc:"SMTP notification settings"`
	Telegram  []Telegram `fig:"telegram" json:"telegram,omitempty" doc:"Telegram notification settings"`
	Pushover  []Pushover `fig:"pushover" json:"pushover,omitempty" doc:"Pushover notification settings"`
	Ntfy      []Ntfy     `fig:"ntfy" json:"ntfy,omitempty" doc:"Ntfy notification settings"`
	Webhook   []Webhook  `fig:"webhook" json:"webhook,omitempty" doc:"Webhook notification settings"`
}

type General struct {
	Title            string `fig:"title" json:"title,omitempty" doc:"Notification title" default:"Frigate Alert"`
	TimeFormat       string `fig:"timeformat" json:"timeformat,omitempty" doc:"Time format used in notifications" default:""`
	NoSnap           string `fig:"nosnap,omitempty" json:"nosnap" enum:"allow,drop" doc:"Allow/Drop events if they do not have a snapshot" default:"allow"`
	SnapBbox         bool   `fig:"snap_bbox,omitempty" json:"snap_bbox" enum:"true,false" doc:"Include bounding box on snapshots" default:false`
	SnapTimestamp    bool   `fig:"snap_timestamp,omitempty" json:"snap_timestamp" enum:"true,false" doc:"Include timestamp on snapshots" default:false`
	SnapCrop         bool   `fig:"snap_crop,omitempty"  json:"snap_crop" enum:"true,false" doc:"Crop snapshots" default:false`
	NotifyOnce       bool   `fig:"notify_once,omitempty"  json:"notify_once" enum:"true,false" doc:"Only notify once per event (For app mode: events)" default:false`
	NotifyDetections bool   `fig:"notify_detections,omitempty" json:"notify_detections" enum:"true,false" doc:"Enable notifications on detection (For app mode: reviews)" default:false`
	RecheckDelay     int    `fig:"recheck_delay" json:"recheck_delay" default:"0" doc:"Delay before re-checking event details from Frigate"`
	AudioOnly        string `fig:"audio_only" json:"audio_only" enum:"allow,drop" doc:"Allow/Drop events that only contain audio detections" default:"allow"`
}

type Quiet struct {
	Start string `fig:"start" json:"start,omitempty" example:"02:30" pattern:"(\d)?\d:\d\d" doc:"Start time for quiet hours" default:""`
	End   string `fig:"end" json:"end,omitempty" example:"05:45" pattern:"(\d)?\d:\d\d" doc:"End time for quiet hours" default:""`
}

type Zones struct {
	Unzoned string   `fig:"unzoned" json:"unzoned,omitempty" enum:"allow,drop" doc:"Allow/Drop events when object is outside a zone" default:"allow"`
	Allow   []string `fig:"allow" json:"allow,omitempty" doc:"List of zones to allow alerts from" default:[]`
	Block   []string `fig:"block" json:"block,omitempty" doc:"List of zones to always block" default:[]`
}

type Labels struct {
	MinScore float64  `fig:"min_score" json:"min_score" minimum:"0" maximum:"100" doc:"Set minimum score before events will notify" default:"0"`
	Allow    []string `fig:"allow" json:"allow,omitempty" doc:"List of labels to allow alerts from" default:[]`
	Block    []string `fig:"block" json:"block,omitempty" doc:"List of labels to always block" default:[]`
}

type AlertFilter struct {
	Cameras   []string `fig:"cameras" json:"cameras,omitempty" doc:"List of cameras that will use this alert provider`
	Zones     []string `fig:"zones" json:"zones,omitempty" doc:"List of zones that will use this alert provider`
	Quiet     Quiet    `fig:"quiet" json:"quiet,omitempty" doc:"Quiet period for this alert provider"`
	Labels    []string `fig:"labels" json:"labels,omitempty" doc:"List of labels that will use this alert provider"`
	Sublabels []string `fig:"sublabels" json:"sublabels,omitempty" doc:"List of sublabels that will use this alert provider"`
}

type Discord struct {
	Enabled  bool        `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Discord" default:false`
	Webhook  string      `fig:"webhook" json:"webhook,omitempty" doc:"Discord webhook URL to send alerts" default:""`
	Template string      `fig:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Gotify struct {
	Enabled  bool        `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Gotify" default:false`
	Server   string      `fig:"server" json:"server,omitempty" doc:"Gotify server URL" default:""`
	Token    string      `fig:"token" json:"token,omitempty" doc:"Gotify app token" default:""`
	Insecure bool        `fig:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:false`
	Template string      `fig:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type SMTP struct {
	Enabled   bool        `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via SMTP" default:false`
	Server    string      `fig:"server" json:"server,omitempty" doc:"SMTP server hostname or IP address" default:""`
	Port      int         `fig:"port" json:"port,omitempty" minimum:"1" maximum:"65535" doc:"SMTP server port" default:25`
	TLS       bool        `fig:"tls" json:"tls,omitempty" enum:"true,false" doc:"Enable/Disable TLS connection" default:false`
	User      string      `fig:"user" json:"user,omitempty" doc:"SMTP user for authentication" default:""`
	Password  string      `fig:"password" json:"password,omitempty" doc:"SMTP password for authentication" default:""`
	From      string      `fig:"from" json:"from,omitempty" format:"email" doc:"SMTP sender" default:""`
	Recipient string      `fig:"recipient" json:"recipient,omitempty" format:"email" doc:"SMTP recipient" default:""`
	Template  string      `fig:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Insecure  bool        `fig:"ignoressl" enum:"true,false" json:"ignoressl,omitempty" default:false`
	Filters   AlertFilter `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Telegram struct {
	Enabled  bool        `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Telegram" default:false`
	ChatID   int64       `fig:"chatid" json:"chatid,omitempty" minimum:"1" doc:"Telegram chat ID" default:"0"`
	Token    string      `fig:"token" json:"token,omitempty" doc:"Telegram bot token" default:""`
	Template string      `fig:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Pushover struct {
	Enabled  bool        `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Pushover" default:false`
	Token    string      `fig:"token" json:"token,omitempty" doc:"Pushover app token" default:""`
	Userkey  string      `fig:"userkey" json:"userkey,omitempty" doc:"Pushover user key" default:""`
	Devices  string      `fig:"devices" json:"devices,omitempty" doc:"Pushover devices to target for notification" default:""`
	Sound    string      `fig:"sound" json:"sound,omitempty" doc:"Pushover notification sound" default:"pushover"`
	Priority int         `fig:"priority" json:"priority,omitempty" minimum:"-2" maximum:"2" doc:"Pushover message priority" default:"0"`
	Retry    int         `fig:"retry" json:"retry,omitempty" doc:"Retry interval for emergency notifications (Priority 2)" default:"0"`
	Expire   int         `fig:"expire" json:"expire,omitempty" doc:"Expiration timer for emergency notifications (Priority 2)" default:"0"`
	TTL      int         `fig:"ttl" json:"ttl,omitempty" doc:"Time to Live for notification messages" default:"0"`
	Template string      `fig:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Ntfy struct {
	Enabled  bool                `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Ntfy" default:false`
	Server   string              `fig:"server" json:"server,omitempty" doc:"Ntfy Server address" default:""`
	Topic    string              `fig:"topic" json:"topic,omitempty" doc:"Ntfy topic" default:""`
	Insecure bool                `fig:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:false`
	Headers  []map[string]string `fig:"headers" json:"headers,omitempty" doc:"HTTP headers to include with Ntfy notifications" default:[]`
	Template string              `fig:"template" json:"template,omitempty" doc:"Custom message template" default:""`
	Filters  AlertFilter         `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Webhook struct {
	Enabled  bool                `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable notifications via Webhook" default:false`
	Server   string              `fig:"server" json:"server,omitempty" doc:"Webhook address" default:""`
	Insecure bool                `fig:"ignoressl" json:"ignoressl,omitempty" doc:"Ignore TLS/SSL errors" default:false`
	Method   string              `fig:"method" json:"method,omitempty" enum:"GET,POST" doc:"HTTP method for webhook notifications" default:"POST"`
	Params   []map[string]string `fix:"params" json:"params,omitempty" doc:"URL parameters for webhook notifications"`
	Headers  []map[string]string `fig:"headers" json:"headers,omitempty" doc:"HTTP headers for webhook notifications"`
	Template interface{}         `fig:"template" json:"template,omitempty" doc:"Custom message template"`
	Filters  AlertFilter         `fig:"filters" json:"filters,omitempty" doc:"Filter notifications sent via this provider"`
}

type Monitor struct {
	Enabled  bool   `fig:"enabled" json:"enabled" enum:"true,false" doc:"Enable monitoring via external uptime application" default:false`
	URL      string `fig:"url" json:"url,omitempty" doc:"Address of monitoring server" default:""`
	Interval int    `fig:"interval" json:"interval,omitempty" minimum:"1" maximum:"10000000" doc:"Interval between check-in messages" default:"60"`
	Insecure bool   `fig:"ignoressl" json:"ignoressl,omitempty" enum:"true,false" doc:"Ignore TLS/SSL errors" default:false`
}
