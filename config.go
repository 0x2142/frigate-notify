package main

import (
	"log"
	"os"

	frigate "github.com/0x2142/frigate-notify/events"
	"github.com/0x2142/frigate-notify/notifier"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Frigate Frigate `yaml:"frigate"`
	Alerts  Alerts  `yaml:"alerts"`
	Monitor Monitor `yaml:"monitor"`
}

type Frigate struct {
	Server   string  `yaml:"server"`
	Insecure bool    `yaml:"ignoressl"`
	WebAPI   WebAPI  `yaml:"webapi"`
	MQTT     MQTT    `yaml:"mqtt"`
	Cameras  Cameras `yaml:"cameras"`
}

type WebAPI struct {
	Enabled  bool `yaml:"enabled"`
	Interval int  `yaml:"interval"`
}

type MQTT struct {
	Enabled  bool   `yaml:"enabled"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Cameras struct {
	Exclude []string `yaml:"exclude"`
}

type Alerts struct {
	General General `yaml:"general"`
	Discord Discord `yaml:"discord"`
	Gotify  Gotify  `yaml:"gotify"`
	SMTP    SMTP    `yaml:"smtp"`
}

type General struct {
	Title string `yaml:"title"`
}

type Discord struct {
	Enabled bool   `yaml:"enabled"`
	Webhook string `yaml:"webhook"`
}

type Gotify struct {
	Enabled  bool   `yaml:"enabled"`
	Server   string `yaml:"server"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"ignoressl"`
}

type SMTP struct {
	Enabled   bool   `yaml:"enabled"`
	Server    string `yaml:"server"`
	Port      int    `yaml:"port"`
	TLS       bool   `yaml:"tls"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Recipient string `yaml:"recipient"`
}

type Monitor struct {
	Enabled  bool   `yaml:"enabled"`
	URL      string `yaml:"url"`
	Interval int    `yaml:"interval"`
	Insecure bool   `yaml:"ignoressl"`
}

var ConfigData Config

// loadConfig opens & attempts to parse configuration file
func loadConfig(configFile string) {
	if configFile == "" {
		var ok bool
		configFile, ok = os.LookupEnv("FN_CONFIGFILE")
		if !ok {
			configFile = "./config.yml"
		}
	}

	log.Print("Loading config file: ", configFile)
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&ConfigData)
	if err != nil {
		log.Fatal("Error loading config file: ", err)
	}

	// Send config file to validation before completing
	validateConfig()

	log.Print("Config file loaded.")
}

// validateConfig checks config file structure & loads info into associated packages
func validateConfig() {

	if (ConfigData.Frigate.WebAPI.Enabled && ConfigData.Frigate.MQTT.Enabled) || (!ConfigData.Frigate.WebAPI.Enabled && !ConfigData.Frigate.MQTT.Enabled) {
		log.Fatal("Please configure only one polling method: Frigate Web API or MQTT")
	}

	// Check for Frigate server config
	if ConfigData.Frigate.Server != "" {
		frigate.FrigateServerURL = ConfigData.Frigate.Server
		if ConfigData.Frigate.Insecure {
			frigate.FrigateInsecure = true
		}
	} else {
		log.Fatal("Please configure Frigate server URL")
	}

	if ConfigData.Frigate.WebAPI.Enabled {
		// Check for run interval. If none, set at 30 seconds
		if ConfigData.Frigate.WebAPI.Interval == 0 {
			log.Println("No interval specified, using default: 30 seconds")
			ConfigData.Frigate.WebAPI.Interval = 30
		}
	}

	if ConfigData.Frigate.MQTT.Enabled {
		if ConfigData.Frigate.MQTT.Port == 0 {
			ConfigData.Frigate.MQTT.Port = 1883
		}
		// Set MQTT config
		frigate.MQTTServer = ConfigData.Frigate.MQTT.Server
		frigate.MQTTPort = ConfigData.Frigate.MQTT.Port
		frigate.MQTTUser = ConfigData.Frigate.MQTT.Username
		frigate.MQTTPass = ConfigData.Frigate.MQTT.Password
	}

	// Check for camera exclusions
	if len(ConfigData.Frigate.Cameras.Exclude) > 0 {
		log.Println("Cameras to exclude from alerting:")
		for _, c := range ConfigData.Frigate.Cameras.Exclude {
			log.Println(" -", c)
		}
		frigate.ExcludeCameras = ConfigData.Frigate.Cameras.Exclude
	}

	// Load Alert general settings
	if ConfigData.Alerts.General.Title != "" {
		notifier.AlertTitle = ConfigData.Alerts.General.Title
	}

	// Check / Load alerting configuration
	if ConfigData.Alerts.Discord.Enabled {
		log.Print("Discord alerting enabled.")
		notifier.DiscordEnabled = true
		notifier.DiscordWebhookURL = ConfigData.Alerts.Discord.Webhook
		notifier.SetupDiscord()
	}
	if ConfigData.Alerts.Gotify.Enabled {
		log.Print("Gotify alerting enabled.")
		notifier.GotifyEnabled = true
		notifier.GotifyServerURL = ConfigData.Alerts.Gotify.Server
		notifier.GotifyToken = ConfigData.Alerts.Gotify.Token
		if ConfigData.Alerts.Gotify.Insecure {
			notifier.GotifyInsecure = true
		}
	}
	if ConfigData.Alerts.SMTP.Enabled {
		log.Print("SMTP alerting enabled.")
		notifier.SMTPEnabled = true
		notifier.SMTPServer = ConfigData.Alerts.SMTP.Server
		notifier.SMTPPort = ConfigData.Alerts.SMTP.Port
		notifier.SMTPTLS = ConfigData.Alerts.SMTP.TLS
		notifier.SMTPUser = ConfigData.Alerts.SMTP.User
		notifier.SMTPPassword = ConfigData.Alerts.SMTP.Password
		notifier.ParseSMTPRecipients(ConfigData.Alerts.SMTP.Recipient)
	}
}
