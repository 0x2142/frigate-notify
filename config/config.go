package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kkyr/fig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	App     App      `fig:"app"`
	Frigate *Frigate `fig:"frigate" validate:"required"`
	Alerts  *Alerts  `fig:"alerts" validate:"required"`
	Monitor *Monitor `fig:"monitor"`
}

type App struct {
	Mode string `fig:"mode" default:"events"`
}

type Frigate struct {
	Server       string              `fig:"server" validate:"required"`
	Insecure     bool                `fig:"ignoressl" default:false`
	PublicURL    string              `fig:"public_url" default:""`
	Headers      []map[string]string `fig:"headers"`
	StartupCheck StartupCheck        `fig:"startup_check"`
	WebAPI       WebAPI              `fig:"webapi"`
	MQTT         MQTT                `fig:"mqtt"`
	Cameras      Cameras             `fig:"cameras"`
	Version      int                 // Internal use only
}

type StartupCheck struct {
	Attempts int `fig:"attempts" default:5`
	Interval int `fig:"interval" default:30`
}

type WebAPI struct {
	Enabled  bool `fig:"enabled" default:false`
	Interval int  `fig:"interval" default:30`
	TestMode bool `fig:"testmode" default:false`
}

type MQTT struct {
	Enabled     bool   `fig:"enabled" default:false`
	Server      string `fig:"server" default:""`
	Port        int    `fig:"port" default:1883`
	ClientID    string `fig:"clientid" default:"frigate-notify"`
	Username    string `fig:"username" default:""`
	Password    string `fig:"password" default:""`
	TopicPrefix string `fig:"topic_prefix" default:"frigate"`
}

type Cameras struct {
	Exclude []string `fig:"exclude" default:[]`
}

type Alerts struct {
	General   General  `fig:"general"`
	Quiet     Quiet    `fig:"quiet"`
	Zones     Zones    `fig:"zones"`
	Labels    Labels   `fig:"labels"`
	SubLabels Labels   `fig:"sublabels"`
	Discord   Discord  `fig:"discord"`
	Gotify    Gotify   `fig:"gotify"`
	SMTP      SMTP     `fig:"smtp"`
	Telegram  Telegram `fig:"telegram"`
	Pushover  Pushover `fig:"pushover"`
	Ntfy      Ntfy     `fig:"ntfy"`
	Webhook   Webhook  `fig:"webhook"`
}

type General struct {
	Title            string `fig:"title" default:"Frigate Alert"`
	TimeFormat       string `fig:"timeformat" default:""`
	NoSnap           string `fig:"nosnap" default:"allow"`
	SnapBbox         bool   `fig:"snap_bbox" default:false`
	SnapTimestamp    bool   `fig:"snap_timestamp" default:false`
	SnapCrop         bool   `fig:"snap_crop" default:false`
	NotifyOnce       bool   `fig:"notify_once" default:false`
	NotifyDetections bool   `fig:"notify_detections" default:false`
}

type Quiet struct {
	Start string `fig:"start" default:""`
	End   string `fig:"end" default:""`
}

type Zones struct {
	Unzoned string   `fig:"unzoned" default:"allow"`
	Allow   []string `fig:"allow" default:[]`
	Block   []string `fig:"block" default:[]`
}

type Labels struct {
	MinScore float64  `fig:"min_score" default:0`
	Allow    []string `fig:"allow" default:[]`
	Block    []string `fig:"block" default:[]`
}

type Discord struct {
	Enabled  bool   `fig:"enabled" default:false`
	Webhook  string `fig:"webhook" default:""`
	Template string `fig:"template" default:""`
}

type Gotify struct {
	Enabled  bool   `fig:"enabled" default:false`
	Server   string `fig:"server" default:""`
	Token    string `fig:"token" default:""`
	Insecure bool   `fig:"ignoressl" default:false`
	Template string `fig:"template" default:""`
}

type SMTP struct {
	Enabled   bool   `fig:"enabled" default:false`
	Server    string `fig:"server" default:""`
	Port      int    `fig:"port" default:25`
	TLS       bool   `fig:"tls" default:false`
	User      string `fig:"user" default:""`
	Password  string `fig:"password" default:""`
	From      string `fig:"from" default:""`
	Recipient string `fig:"recipient" default:""`
	Template  string `fig:"template" default:""`
	Insecure  bool   `fig:"ignoressl" default:false`
}

type Telegram struct {
	Enabled  bool   `fig:"enabled" default:false`
	ChatID   int64  `fig:"chatid" default:0`
	Token    string `fig:"token" default:""`
	Template string `fig:"template" default:""`
}

type Pushover struct {
	Enabled  bool   `fig:"enabled" default:false`
	Token    string `fig:"token" default:""`
	Userkey  string `fig:"userkey" default:""`
	Devices  string `fig:"devices" default:""`
	Sound    string `fig:"sound" default:""`
	Priority int    `fig:"priority" default:0`
	Retry    int    `fig:"retry" default:0`
	Expire   int    `fig:"expire" default:0`
	TTL      int    `fig:"ttl" default:0`
	Template string `fig:"template" default:""`
}

type Ntfy struct {
	Enabled  bool                `fig:"enabled" default:false`
	Server   string              `fig:"server" default:""`
	Topic    string              `fig:"topic" default:""`
	Insecure bool                `fig:"ignoressl" default:false`
	Headers  []map[string]string `fig:"headers" default:[]`
	Template string              `fig:"template" default:""`
}

type Webhook struct {
	Enabled  bool                   `fig:"enabled" default:false`
	Server   string                 `fig:"server" default:""`
	Insecure bool                   `fig:"ignoressl" default:false`
	Method   string                 `fig:"method" default:"POST"`
	Params   []map[string]string    `fix:"params"`
	Headers  []map[string]string    `fig:"headers"`
	Template map[string]interface{} `fig:"template"`
}

type Monitor struct {
	Enabled  bool   `fig:"enabled" default:false`
	URL      string `fig:"url" default:""`
	Interval int    `fig:"interval" default:60`
	Insecure bool   `fig:"ignoressl" default:false`
}

var ConfigData Config

// loadConfig opens & attempts to parse configuration file
func LoadConfig(configFile string) {
	// Set config file location
	if configFile == "" {
		var ok bool
		configFile, ok = os.LookupEnv("FN_CONFIGFILE")
		if !ok {
			configFile = "./config.yml"
		}
	}

	// Load Config file
	log.Debug().Msgf("Attempting to load config file: %v", configFile)

	err := fig.Load(&ConfigData, fig.File(filepath.Base(configFile)), fig.Dirs(filepath.Dir(configFile)), fig.UseEnv("FN"))
	if err != nil {
		if errors.Is(err, fig.ErrFileNotFound) {
			log.Warn().Msg("Config file could not be read, attempting to load config from environment")
			err = fig.Load(&ConfigData, fig.IgnoreFile(), fig.UseEnv("FN"))
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("Failed to load config from environment!")
			}
		} else {
			log.Fatal().
				Err(err).
				Msg("Failed to load config from file!")
		}
	}
	log.Info().Msg("Config loaded.")

	// Send config file to validation before completing
	//validateConfig()
	validationErrors := ConfigData.validate()
	if len(validationErrors) > 0 {
		fmt.Println()
		log.Error().Msg("Config validation failed:")
		for _, msg := range validationErrors {
			log.Error().Msgf(" - %v", msg)
		}
		fmt.Println()
		log.Fatal().Msg("Please fix config errors before restarting app.")
	} else {
		log.Info().Msg("Config file validated!")
	}
}
