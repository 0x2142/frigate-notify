package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
	"github.com/kkyr/fig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Frigate Frigate `fig:"frigate"`
	Alerts  Alerts  `fig:"alerts"`
	Monitor Monitor `fig:"monitor"`
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
	Version      int
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
	Nfty      Nfty     `fig:"nfty"`
	Ntfy      Ntfy     `fig:"ntfy"`
	Webhook   Webhook  `fig:"webhook"`
}

type General struct {
	Title         string `fig:"title" default:"Frigate Alert"`
	TimeFormat    string `fig:"timeformat" default:""`
	NoSnap        string `fig:"nosnap" default:"allow"`
	SnapBbox      bool   `fig:"snap_bbox" default:false`
	SnapTimestamp bool   `fig:"snap_timestamp" default:false`
	SnapCrop      bool   `fig:"snap_crop" default:false`
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
	Recipient string `fig:"recipient" default:""`
	Template  string `fig:"template" default:""`
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
	Priority int    `fig:"priority" default:0`
	Retry    int    `fig:"retry" default:0`
	Expire   int    `fig:"expire" default:0`
	TTL      int    `fig:"ttl" default:0`
	Template string `fig:"template" default:""`
}

// DEPRECATED: Misspelling of Ntfy
type Nfty struct {
	Enabled  bool                `fig:"enabled" default:false`
	Server   string              `fig:"server" default:""`
	Topic    string              `fig:"topic" default:""`
	Insecure bool                `fig:"ignoressl" default:false`
	Headers  []map[string]string `fig:"headers"`
	Template string              `fig:"template" default:""`
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
	validateConfig()
}

// validateConfig checks config file structure & loads info into associated packages
func validateConfig() {
	var configErrors []string
	var response []byte
	var err error
	log.Debug().Msg("Validating config file...")

	if (ConfigData.Frigate.WebAPI.Enabled && ConfigData.Frigate.MQTT.Enabled) || (!ConfigData.Frigate.WebAPI.Enabled && !ConfigData.Frigate.MQTT.Enabled) {
		configErrors = append(configErrors, "Please configure only one polling method: Frigate Web API or MQTT")
	}

	// Set default web API interval if not specified
	if ConfigData.Frigate.WebAPI.Enabled && ConfigData.Frigate.WebAPI.Interval == 0 {
		ConfigData.Frigate.WebAPI.Interval = 30
	}

	// Warn on test mode being enabled
	if ConfigData.Frigate.WebAPI.Enabled && ConfigData.Frigate.WebAPI.TestMode {
		log.Warn().Msg("~~~~~~~~~~~~~~~~~~~")
		log.Warn().Msg("WARNING: Test Mode is enabled.")
		log.Warn().Msg("This is intended for development only & will only query Frigate for the last event.")
		log.Warn().Msg("Do not enable this in production! App will not accurately check for events.")
		log.Warn().Msg("~~~~~~~~~~~~~~~~~~~")
	}

	// Check if Frigate server URL contains protocol, assume HTTP if not specified
	if !strings.Contains(ConfigData.Frigate.Server, "http://") && !strings.Contains(ConfigData.Frigate.Server, "https://") {
		log.Warn().Msg("No protocol specified on Frigate Server. Assuming http://. If this is incorrect, please adjust the config file.")
		ConfigData.Frigate.Server = fmt.Sprintf("http://%s", ConfigData.Frigate.Server)
	}

	// Test connectivity to Frigate
	log.Debug().Msg("Checking connection to Frigate server...")
	statsAPI := fmt.Sprintf("%s/api/stats", ConfigData.Frigate.Server)
	current_attempt := 1
	if ConfigData.Frigate.StartupCheck.Attempts == 0 {
		ConfigData.Frigate.StartupCheck.Attempts = 5
	}
	if ConfigData.Frigate.StartupCheck.Interval == 0 {
		ConfigData.Frigate.StartupCheck.Interval = 30
	}
	for current_attempt < ConfigData.Frigate.StartupCheck.Attempts {
		response, err = util.HTTPGet(statsAPI, ConfigData.Frigate.Insecure, ConfigData.Frigate.Headers...)
		if err != nil {
			log.Warn().
				Err(err).
				Int("attempt", current_attempt).
				Int("max_tries", ConfigData.Frigate.StartupCheck.Attempts).
				Int("interval", ConfigData.Frigate.StartupCheck.Interval).
				Msgf("Cannot reach Frigate server at %v", ConfigData.Frigate.Server)
			time.Sleep(time.Duration(ConfigData.Frigate.StartupCheck.Interval) * time.Second)
			current_attempt += 1
		} else {
			break
		}
	}
	if current_attempt == ConfigData.Frigate.StartupCheck.Attempts {
		log.Fatal().
			Err(err).
			Msgf("Max attempts reached - Cannot reach Frigate server at %v", ConfigData.Frigate.Server)
	}
	var stats models.FrigateStats
	json.Unmarshal([]byte(response), &stats)
	log.Info().Msgf("Successfully connected to %v", ConfigData.Frigate.Server)
	if stats.Service.Version != "" {
		log.Debug().Msgf("Frigate server is running version %v", stats.Service.Version)
		// Save major version number
		ConfigData.Frigate.Version, _ = strconv.Atoi(strings.Split(stats.Service.Version, ".")[1])
	}

	// Check Public / External URL if set
	if ConfigData.Frigate.PublicURL != "" {
		if !strings.Contains(ConfigData.Frigate.PublicURL, "http://") && !strings.Contains(ConfigData.Frigate.PublicURL, "https://") {
			configErrors = append(configErrors, "Public URL must include http:// or https://")
		}
	}

	// Check for camera exclusions
	if len(ConfigData.Frigate.Cameras.Exclude) > 0 {
		log.Debug().Msg("Cameras to exclude from alerting:")
		for _, c := range ConfigData.Frigate.Cameras.Exclude {
			log.Debug().Msgf(" - %v", c)
		}
	}

	// Check MQTT Config
	if ConfigData.Frigate.MQTT.Enabled {
		log.Debug().Msg("MQTT Enabled.")
		if ConfigData.Frigate.MQTT.Server == "" {
			configErrors = append(configErrors, "No MQTT server address specified!")
		}
		if ConfigData.Frigate.MQTT.Username != "" && ConfigData.Frigate.MQTT.Password == "" {
			configErrors = append(configErrors, "MQTT user provided, but no password!")
		}
		if ConfigData.Frigate.MQTT.Port == 0 {
			ConfigData.Frigate.MQTT.Port = 1883
		}
	}

	// Check quiet hours config
	if ConfigData.Alerts.Quiet.Start != "" || ConfigData.Alerts.Quiet.End != "" {
		timeformat := "15:04"
		validstart := true
		validend := true
		if _, ok := time.Parse(timeformat, ConfigData.Alerts.Quiet.Start); ok != nil {
			configErrors = append(configErrors, "Start time for quiet hours does not match format: 00:00")
			validstart = false
		}
		if _, ok := time.Parse(timeformat, ConfigData.Alerts.Quiet.End); ok != nil {
			configErrors = append(configErrors, "End time for quiet hours does not match format: 00:00")
			validend = false
		}
		if validstart && validend {
			log.Debug().Msgf("Quiet hours enabled. Start: %v, End: %v", ConfigData.Alerts.Quiet.Start, ConfigData.Alerts.Quiet.End)
		}
	}

	// Check action on no snapshot available
	if strings.ToLower(ConfigData.Alerts.General.NoSnap) != "allow" && strings.ToLower(ConfigData.Alerts.General.NoSnap) != "drop" {
		configErrors = append(configErrors, "Option for nosnap must be 'allow' or 'drop'")
	} else {
		log.Debug().Msgf("Events without a snapshot: %v", strings.ToLower(ConfigData.Alerts.General.NoSnap))
	}

	// Check Zone filtering config
	if strings.ToLower(ConfigData.Alerts.Zones.Unzoned) != "allow" && strings.ToLower(ConfigData.Alerts.Zones.Unzoned) != "drop" {
		configErrors = append(configErrors, "Option for unzoned events must be 'allow' or 'drop'")
	} else {
		log.Debug().Msgf("Events outside a zone: %v", strings.ToLower(ConfigData.Alerts.Zones.Unzoned))
	}

	if len(ConfigData.Alerts.Zones.Allow) > 0 {
		log.Debug().Msg("Zones to generate alerts for:")
		for _, c := range ConfigData.Alerts.Zones.Allow {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("All zones included in alerts")
	}
	if len(ConfigData.Alerts.Zones.Block) > 0 {
		log.Debug().Msg("Zones to exclude from alerting:")
		for _, c := range ConfigData.Alerts.Zones.Block {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("No zones excluded")
	}

	// Check Label filtering config
	if ConfigData.Alerts.Labels.MinScore > 0 {
		log.Debug().Msgf("Label required minimum score: %v", ConfigData.Alerts.Labels.MinScore)
	}
	if len(ConfigData.Alerts.Labels.Allow) > 0 {
		log.Debug().Msg("Labels to generate alerts for:")
		for _, c := range ConfigData.Alerts.Labels.Allow {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("All labels included in alerts")
	}
	if len(ConfigData.Alerts.Labels.Block) > 0 {
		log.Debug().Msg("Labels to exclude from alerting:")
		for _, c := range ConfigData.Alerts.Labels.Block {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("No labels excluded")
	}

	// Check Subabel filtering config
	if len(ConfigData.Alerts.SubLabels.Allow) > 0 {
		log.Debug().Msg("Sublabels to generate alerts for:")
		for _, c := range ConfigData.Alerts.SubLabels.Allow {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("All Sublabels included in alerts")
	}
	if len(ConfigData.Alerts.SubLabels.Block) > 0 {
		log.Debug().Msg("Sublabels to exclude from alerting:")
		for _, c := range ConfigData.Alerts.SubLabels.Block {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("No Sublabels excluded")
	}

	// Check / Load alerting configuration
	if ConfigData.Alerts.Discord.Enabled {
		log.Debug().Msg("Discord alerting enabled.")
		if ConfigData.Alerts.Discord.Webhook == "" {
			configErrors = append(configErrors, "No Discord webhook specified!")
		}
		// Check template syntax
		if msg := checkTemplate("Discord", ConfigData.Alerts.Discord.Template); msg != "" {
			configErrors = append(configErrors, msg)
		}
	}
	if ConfigData.Alerts.Gotify.Enabled {
		log.Debug().Msg("Gotify alerting enabled.")
		// Check if Gotify server URL contains protocol, assume HTTP if not specified
		if !strings.Contains(ConfigData.Alerts.Gotify.Server, "http://") && !strings.Contains(ConfigData.Alerts.Gotify.Server, "https://") {
			log.Debug().Msg("No protocol specified on Gotify Server. Assuming http://. If this is incorrect, please adjust the config file.")
			ConfigData.Alerts.Gotify.Server = fmt.Sprintf("http://%s", ConfigData.Alerts.Gotify.Server)
		}
		if ConfigData.Alerts.Gotify.Server == "" {
			configErrors = append(configErrors, "No Gotify server specified!")
		}
		if ConfigData.Alerts.Gotify.Token == "" {
			configErrors = append(configErrors, "No Gotify token specified!")
		}
		// Check template syntax
		if msg := checkTemplate("Gotify", ConfigData.Alerts.Gotify.Template); msg != "" {
			configErrors = append(configErrors, msg)
		}
	}
	if ConfigData.Alerts.SMTP.Enabled {
		log.Debug().Msg("SMTP alerting enabled.")
		if ConfigData.Alerts.SMTP.Server == "" {
			configErrors = append(configErrors, "No SMTP server specified!")
		}
		if ConfigData.Alerts.SMTP.Recipient == "" {
			configErrors = append(configErrors, "No SMTP recipients specified!")
		}
		if ConfigData.Alerts.SMTP.User != "" && ConfigData.Alerts.SMTP.Password == "" {
			configErrors = append(configErrors, "SMTP username in config, but no password provided!")
		}
		if ConfigData.Alerts.SMTP.Port == 0 {
			ConfigData.Alerts.SMTP.Port = 25
		}
		// Check template syntax
		if msg := checkTemplate("SMTP", ConfigData.Alerts.SMTP.Template); msg != "" {
			configErrors = append(configErrors, msg)
		}
	}
	if ConfigData.Alerts.Telegram.Enabled {
		log.Debug().Msg("Telegram alerting enabled.")
		if ConfigData.Alerts.Telegram.ChatID == 0 {
			configErrors = append(configErrors, "No Telegram Chat ID specified!")
		}
		if ConfigData.Alerts.Telegram.Token == "" {
			configErrors = append(configErrors, "No Telegram bot token specified!")
		}
		// Check template syntax
		if msg := checkTemplate("Telegram", ConfigData.Alerts.Telegram.Template); msg != "" {
			configErrors = append(configErrors, msg)
		}
	}
	if ConfigData.Alerts.Pushover.Enabled {
		log.Debug().Msg("Pushover alerting enabled.")
		if ConfigData.Alerts.Pushover.Token == "" {
			configErrors = append(configErrors, "No Pushover API token specified!")
		}
		if ConfigData.Alerts.Pushover.Userkey == "" {
			configErrors = append(configErrors, "No Pushover user key specified!")
		}
		if ConfigData.Alerts.Pushover.Priority < -2 || ConfigData.Alerts.Pushover.Priority > 2 {
			configErrors = append(configErrors, "Pushover priority must be between -2 and 2!")
		}
		// Priority 2 is emergency, needs a retry interval & expiration set
		if ConfigData.Alerts.Pushover.Priority == 2 {
			if ConfigData.Alerts.Pushover.Retry == 0 || ConfigData.Alerts.Pushover.Expire == 0 {
				configErrors = append(configErrors, "Pushover retry interval & expiration must be set with priority 2!")
			}
			if ConfigData.Alerts.Pushover.Retry < 30 {
				configErrors = append(configErrors, "Pushover retry cannot be less than 30 seconds!")
			}
		}
		if ConfigData.Alerts.Pushover.TTL < 0 {
			configErrors = append(configErrors, "Pushover TTL cannot be negative!")
		}
		// Check template syntax
		if msg := checkTemplate("Pushover", ConfigData.Alerts.Pushover.Template); msg != "" {
			configErrors = append(configErrors, msg)
		}
	}
	// Deprecation warning
	// TODO: Remove misspelled Ntfy config with v0.4.0 or later
	if ConfigData.Alerts.Nfty.Enabled {
		log.Warn().Msg("Config for 'nfty' will be deprecated due to misspelling. Please update config to 'ntfy'")
		// Copy data to new Ntfy struct
		ConfigData.Alerts.Ntfy.Enabled = ConfigData.Alerts.Nfty.Enabled
		ConfigData.Alerts.Ntfy.Server = ConfigData.Alerts.Nfty.Server
		ConfigData.Alerts.Ntfy.Topic = ConfigData.Alerts.Nfty.Topic
		ConfigData.Alerts.Ntfy.Insecure = ConfigData.Alerts.Nfty.Insecure
		ConfigData.Alerts.Ntfy.Headers = ConfigData.Alerts.Nfty.Headers
		ConfigData.Alerts.Ntfy.Template = ConfigData.Alerts.Nfty.Template
	}
	if ConfigData.Alerts.Ntfy.Enabled {
		log.Debug().Msg("Ntfy alerting enabled.")
		if ConfigData.Alerts.Ntfy.Server == "" {
			configErrors = append(configErrors, "No Ntfy server specified!")
		}
		if ConfigData.Alerts.Ntfy.Topic == "" {
			configErrors = append(configErrors, "No Ntfy topic specified!")
		}
		// Check template syntax
		if msg := checkTemplate("Ntfy", ConfigData.Alerts.Ntfy.Template); msg != "" {
			configErrors = append(configErrors, msg)
		}
	}
	if ConfigData.Alerts.Webhook.Enabled {
		log.Debug().Msg("Webhook alerting enabled.")
		if ConfigData.Alerts.Webhook.Server == "" {
			configErrors = append(configErrors, "No Webhook server specified!")
		}
	}

	// Validate monitoring config
	if ConfigData.Monitor.Enabled {
		log.Debug().Msg("App monitoring enabled.")
		if ConfigData.Monitor.URL == "" {
			configErrors = append(configErrors, "App monitor enabled but no URL specified!")
		}
		if ConfigData.Monitor.Interval == 0 {
			ConfigData.Monitor.Interval = 60
		}
	}

	if len(configErrors) > 0 {
		fmt.Println()
		log.Error().Msg("Config validation failed:")
		for _, msg := range configErrors {
			log.Error().Msgf(" - %v", msg)
		}
		fmt.Println()
		log.Fatal().Msg("Please fix config errors before restarting app.")
	} else {
		log.Info().Msg("Config file validated!")
	}
}

func checkTemplate(provider, alertTemplate string) string {
	var templateError string
	_, err := template.New("").Parse(alertTemplate)
	if err != nil {
		templateError = fmt.Sprintf("Failed to parse %s template: %s", provider, err.Error())
	}
	return templateError
}
