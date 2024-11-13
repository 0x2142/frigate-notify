package config

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
	"github.com/rs/zerolog/log"
)

func (c *Config) validate() []string {
	var validationErrors []string
	log.Debug().Msg("Validating config file...")

	// Check Frigate Connectivity
	if results := c.validateFrigateConnectivity(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate App Mode
	if results := c.validateAppMode(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate API settings
	if c.App.API.Enabled {
		if results := c.validateAPI(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate Frigate polling method
	if results := c.validateFrigatePolling(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate MQTT settings
	if c.Frigate.MQTT.Enabled {
		if results := c.validateMQTT(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Check / Log info on Camera exclusions
	c.validateCameraExclusions()

	// Validate Quiet Hours Config
	if results := c.validateQuietHours(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate alert general section settings
	if results := c.validateAlertGeneral(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate Zone Filters
	if results := c.validateZoneFilters(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate Label Filters
	if results := c.validateLabelFiltering(); len(results) > 0 {
		validationErrors = append(validationErrors, results...)
	}

	// Validate Discord
	if c.Alerts.Discord.Enabled {
		if results := c.validateDiscord(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate Gotify
	if c.Alerts.Gotify.Enabled {
		if results := c.validateGotify(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate SMTP
	if c.Alerts.SMTP.Enabled {
		if results := c.validateSMTP(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate Telegram
	if c.Alerts.Telegram.Enabled {
		if results := c.validateTelegram(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate Pushover
	if c.Alerts.Pushover.Enabled {
		if results := c.validatePushover(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate Ntfy
	if c.Alerts.Ntfy.Enabled {
		if results := c.validateNtfy(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate Webhook
	if c.Alerts.Webhook.Enabled {
		if results := c.validateWebhook(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}

	// Validate that at least one alert profile is enabled
	if result := c.validateAlertingEnabled(); result != "" {
		validationErrors = append(validationErrors, result)
	}

	// Validate app health check / monitoring config
	if c.Monitor.Enabled {
		if results := c.validateAppMonitoring(); len(results) > 0 {
			validationErrors = append(validationErrors, results...)
		}
	}
	return validationErrors
}

func (c *Config) validateAppMode() []string {
	var appErrors []string
	if strings.ToLower(c.App.Mode) != "events" && strings.ToLower(c.App.Mode) != "reviews" {
		appErrors = append(appErrors, "MQTT mode must be 'events' or 'reviews'")
	}
	if Internal.FrigateVersion < 14 && strings.ToLower(c.App.Mode) == "reviews" {
		appErrors = append(appErrors, "Frigate must be version 0.14 or higher to use 'reviews' mode. Please use 'events' mode or update Frigate.")
	}
	log.Debug().Msgf("App mode: %v", c.App.Mode)
	return appErrors
}

func (c *Config) validateAPI() []string {
	var apiErrors []string

	// Set default port if needed
	if c.App.API.Port == 0 {
		c.App.API.Port = 8000
	}

	if c.App.API.Port <= 0 || c.App.API.Port > 65535 {
		apiErrors = append(apiErrors, "Invalid API port")
	}

	return apiErrors
}

func (c *Config) validateFrigatePolling() []string {
	var pollingErrors []string
	webapi := c.Frigate.WebAPI.Enabled
	mqtt := c.Frigate.MQTT.Enabled
	// Check that only one polling method is configured
	if (webapi && mqtt) || (!webapi && !mqtt) {
		pollingErrors = append(pollingErrors, "Please configure only one polling method: Frigate Web API or MQTT")
	}
	if webapi {
		log.Debug().Msgf("Event polling method: Web API")
	}
	if mqtt {
		log.Debug().Msgf("Event polling method: MQTT")
	}

	// Set default web API interval if not specified
	if c.Frigate.WebAPI.Interval == 0 {
		c.Frigate.WebAPI.Interval = 30
	}

	// Warn on test mode being enabled
	if c.Frigate.WebAPI.Enabled && c.Frigate.WebAPI.TestMode {
		log.Warn().Msg("~~~~~~~~~~~~~~~~~~~")
		log.Warn().Msg("WARNING: Test Mode is enabled.")
		log.Warn().Msg("This is intended for development only & will only query Frigate for the last event.")
		log.Warn().Msg("Do not enable this in production! App will not accurately check for events.")
		log.Warn().Msg("~~~~~~~~~~~~~~~~~~~")
	}

	return pollingErrors
}

func (c *Config) validateFrigateConnectivity() []string {
	var response []byte
	var err error
	var connectivityErrors []string

	url := c.Frigate.Server
	max_attempts := c.Frigate.StartupCheck.Attempts
	interval := c.Frigate.StartupCheck.Interval

	// Check if Frigate server URL contains protocol, assume HTTP if not specified
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		log.Warn().Msgf("No protocol specified on Frigate server URL, so we'll try http://%s. If this is incorrect, please adjust the config file.", c.Frigate.Server)
		c.Frigate.Server = fmt.Sprintf("http://%s", url)
		url = c.Frigate.Server
	}

	// Check Public / External URL if set
	if c.Frigate.PublicURL != "" {
		if !strings.Contains(c.Frigate.PublicURL, "http://") && !strings.Contains(c.Frigate.PublicURL, "https://") {
			connectivityErrors = append(connectivityErrors, "Public URL must include http:// or https://")
		}
	} else {
		// If Public URL not explicitly set, use local Frigate URL
		c.Frigate.PublicURL = c.Frigate.Server
	}

	// Check HTTP header template syntax
	if msg := validateTemplate("Frigate HTTP Headers", c.Alerts.General.Title); msg != "" {
		connectivityErrors = append(connectivityErrors, msg)
	}

	// Test connectivity to Frigate
	log.Debug().Msg("Checking connection to Frigate server...")
	statsAPI := fmt.Sprintf("%s/api/stats", url)
	current_attempt := 1
	if max_attempts == 0 {
		max_attempts = 5
	}
	if interval == 0 {
		interval = 30
	}
	for current_attempt < max_attempts {
		response, err = util.HTTPGet(statsAPI, c.Frigate.Insecure, "", c.Frigate.Headers...)
		if err != nil {
			Internal.Status.Frigate.API = "unreachable"
			log.Warn().
				Err(err).
				Int("attempt", current_attempt).
				Int("max_tries", max_attempts).
				Int("interval", interval).
				Msgf("Cannot reach Frigate server at %v", url)
			time.Sleep(time.Duration(interval) * time.Second)
			current_attempt += 1
		} else {
			break
		}
	}
	if current_attempt == max_attempts {
		Internal.Status.Frigate.API = "unreachable"
		log.Fatal().
			Err(err).
			Msgf("Max attempts reached - Cannot reach Frigate server at %v", url)
	}
	var stats models.FrigateStats
	json.Unmarshal([]byte(response), &stats)
	log.Info().Msgf("Successfully connected to %v", url)
	Internal.Status.Frigate.API = "ok"
	if stats.Service.Version != "" {
		log.Debug().Msgf("Frigate server is running version %v", stats.Service.Version)
		// Save major version number
		Internal.FrigateVersion, _ = strconv.Atoi(strings.Split(stats.Service.Version, ".")[1])
	}
	return connectivityErrors
}

func (c *Config) validateMQTT() []string {
	var configErrors []string
	// Check MQTT Config
	log.Debug().Msg("MQTT Enabled.")
	if c.Frigate.MQTT.Server == "" {
		configErrors = append(configErrors, "No MQTT server address specified")
	}
	if c.Frigate.MQTT.Username != "" && c.Frigate.MQTT.Password == "" {
		configErrors = append(configErrors, "MQTT user provided, but no password")
	}
	// Set default port if needed
	if c.Frigate.MQTT.Port == 0 {
		c.Frigate.MQTT.Port = 1883
	}
	return configErrors
}

func (c *Config) validateCameraExclusions() {
	// Check for camera exclusions
	if len(c.Frigate.Cameras.Exclude) > 0 {
		log.Debug().Msg("Cameras to exclude from alerting:")
		for _, c := range c.Frigate.Cameras.Exclude {
			log.Debug().Msgf(" - %v", c)
		}
	}
}

func (c *Config) validateQuietHours() []string {
	var quietHoursErrors []string
	// Check quiet hours config
	if c.Alerts.Quiet.Start != "" || c.Alerts.Quiet.End != "" {
		timeformat := "15:04"
		validstart := true
		validend := true
		if _, ok := time.Parse(timeformat, c.Alerts.Quiet.Start); ok != nil {
			quietHoursErrors = append(quietHoursErrors, "Start time for quiet hours does not match format: 00:00")
			validstart = false
		}
		if _, ok := time.Parse(timeformat, c.Alerts.Quiet.End); ok != nil {
			quietHoursErrors = append(quietHoursErrors, "End time for quiet hours does not match format: 00:00")
			validend = false
		}
		if validstart && validend {
			log.Debug().Msgf("Quiet hours enabled. Start: %v, End: %v", c.Alerts.Quiet.Start, c.Alerts.Quiet.End)
		}
	}
	return quietHoursErrors
}

func (c *Config) validateAlertGeneral() []string {
	var alertErrors []string
	// Check action on no snapshot available
	if strings.ToLower(c.Alerts.General.NoSnap) != "allow" && strings.ToLower(c.Alerts.General.NoSnap) != "drop" {
		alertErrors = append(alertErrors, "Option for nosnap must be 'allow' or 'drop'")
	} else {
		log.Debug().Msgf("Events without a snapshot: %v", strings.ToLower(c.Alerts.General.NoSnap))
	}

	// Notify_Once
	log.Debug().Msgf("Notify only once per event: %v", c.Alerts.General.NotifyOnce)

	// Notify_Detections
	log.Debug().Msgf("Notify on Detections: %v", c.Alerts.General.NotifyDetections)

	// Check title template syntax
	if msg := validateTemplate("Alert Title", c.Alerts.General.Title); msg != "" {
		alertErrors = append(alertErrors, msg)
	}

	return alertErrors
}

func (c *Config) validateZoneFilters() []string {
	var filterErrors []string
	// Check Zone filtering config
	if strings.ToLower(c.Alerts.Zones.Unzoned) != "allow" && strings.ToLower(c.Alerts.Zones.Unzoned) != "drop" {
		filterErrors = append(filterErrors, "Option for unzoned events must be 'allow' or 'drop'")
	} else {
		log.Debug().Msgf("Events outside a zone: %v", strings.ToLower(c.Alerts.Zones.Unzoned))
	}

	if len(c.Alerts.Zones.Allow) > 0 {
		log.Debug().Msg("Zones to generate alerts for:")
		for _, c := range c.Alerts.Zones.Allow {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("All zones included in alerts")
	}
	if len(c.Alerts.Zones.Block) > 0 {
		log.Debug().Msg("Zones to exclude from alerting:")
		for _, c := range c.Alerts.Zones.Block {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("No zones excluded")
	}
	return filterErrors
}

func (c *Config) validateLabelFiltering() []string {
	var labelErrors []string
	// Check Label filtering config
	if c.Alerts.Labels.MinScore > 0 {
		log.Debug().Msgf("Label required minimum score: %v", c.Alerts.Labels.MinScore)
	}
	if len(c.Alerts.Labels.Allow) > 0 {
		log.Debug().Msg("Labels to generate alerts for:")
		for _, c := range c.Alerts.Labels.Allow {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("All labels included in alerts")
	}
	if len(c.Alerts.Labels.Block) > 0 {
		log.Debug().Msg("Labels to exclude from alerting:")
		for _, c := range c.Alerts.Labels.Block {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("No labels excluded")
	}

	// Check Subabel filtering config
	if len(c.Alerts.SubLabels.Allow) > 0 {
		log.Debug().Msg("Sublabels to generate alerts for:")
		for _, c := range c.Alerts.SubLabels.Allow {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("All Sublabels included in alerts")
	}
	if len(c.Alerts.SubLabels.Block) > 0 {
		log.Debug().Msg("Sublabels to exclude from alerting:")
		for _, c := range c.Alerts.SubLabels.Block {
			log.Debug().Msgf(" - %v", c)
		}
	} else {
		log.Debug().Msg("No Sublabels excluded")
	}
	return labelErrors
}

func (c *Config) validateDiscord() []string {
	var discordErrors []string
	log.Debug().Msg("Discord alerting enabled.")
	if c.Alerts.Discord.Webhook == "" {
		discordErrors = append(discordErrors, "No Discord webhook specified!")
	}
	// Check template syntax
	if msg := validateTemplate("Discord", c.Alerts.Discord.Template); msg != "" {
		discordErrors = append(discordErrors, msg)
	}
	return discordErrors
}

func (c *Config) validateGotify() []string {
	var gotifyErrors []string
	log.Debug().Msg("Gotify alerting enabled.")
	if c.Alerts.Gotify.Server == "" {
		gotifyErrors = append(gotifyErrors, "No Gotify server specified!")
	}
	if c.Alerts.Gotify.Token == "" {
		gotifyErrors = append(gotifyErrors, "No Gotify token specified!")
	}
	// Check if Gotify server URL contains protocol, assume HTTP if not specified
	if !strings.Contains(c.Alerts.Gotify.Server, "http://") && !strings.Contains(c.Alerts.Gotify.Server, "https://") {
		log.Debug().Msg("No protocol specified on Gotify Server. Assuming http://. If this is incorrect, please adjust the config file.")
		c.Alerts.Gotify.Server = fmt.Sprintf("http://%s", c.Alerts.Gotify.Server)
	}
	// Check template syntax
	if msg := validateTemplate("Gotify", c.Alerts.Gotify.Template); msg != "" {
		gotifyErrors = append(gotifyErrors, msg)
	}
	return gotifyErrors
}

func (c *Config) validateSMTP() []string {
	var smtpErrors []string
	log.Debug().Msg("SMTP alerting enabled.")
	if c.Alerts.SMTP.Server == "" {
		smtpErrors = append(smtpErrors, "No SMTP server specified!")
	}
	if c.Alerts.SMTP.Recipient == "" {
		smtpErrors = append(smtpErrors, "No SMTP recipients specified!")
	}
	if c.Alerts.SMTP.User != "" && c.Alerts.SMTP.Password == "" {
		smtpErrors = append(smtpErrors, "SMTP username in config, but no password provided!")
	}
	if c.Alerts.SMTP.Port == 0 {
		c.Alerts.SMTP.Port = 25
	}
	// Copy `user` to `from` if `from` not explicitly configured
	if c.Alerts.SMTP.From == "" && c.Alerts.SMTP.User != "" {
		c.Alerts.SMTP.From = c.Alerts.SMTP.User
	}
	// Check template syntax
	if msg := validateTemplate("SMTP", c.Alerts.SMTP.Template); msg != "" {
		smtpErrors = append(smtpErrors, msg)
	}

	return smtpErrors
}

func (c *Config) validateTelegram() []string {
	var telegramErrors []string
	log.Debug().Msg("Telegram alerting enabled.")
	if c.Alerts.Telegram.ChatID == 0 {
		telegramErrors = append(telegramErrors, "No Telegram Chat ID specified!")
	}
	if c.Alerts.Telegram.Token == "" {
		telegramErrors = append(telegramErrors, "No Telegram bot token specified!")
	}
	// Check template syntax
	if msg := validateTemplate("Telegram", c.Alerts.Telegram.Template); msg != "" {
		telegramErrors = append(telegramErrors, msg)
	}
	return telegramErrors
}

func (c *Config) validatePushover() []string {
	var pushoverErrors []string
	log.Debug().Msg("Pushover alerting enabled.")
	if c.Alerts.Pushover.Token == "" {
		pushoverErrors = append(pushoverErrors, "No Pushover API token specified!")
	}
	if c.Alerts.Pushover.Userkey == "" {
		pushoverErrors = append(pushoverErrors, "No Pushover user key specified!")
	}
	if c.Alerts.Pushover.Priority < -2 || c.Alerts.Pushover.Priority > 2 {
		pushoverErrors = append(pushoverErrors, "Pushover priority must be between -2 and 2!")
	}
	// Priority 2 is emergency, needs a retry interval & expiration set
	if c.Alerts.Pushover.Priority == 2 {
		if c.Alerts.Pushover.Retry == 0 || c.Alerts.Pushover.Expire == 0 {
			pushoverErrors = append(pushoverErrors, "Pushover retry interval & expiration must be set with priority 2!")
		}
		if c.Alerts.Pushover.Retry < 30 {
			pushoverErrors = append(pushoverErrors, "Pushover retry cannot be less than 30 seconds!")
		}
	}
	if c.Alerts.Pushover.TTL < 0 {
		pushoverErrors = append(pushoverErrors, "Pushover TTL cannot be negative!")
	}

	// Check template syntax
	if msg := validateTemplate("Pushover", c.Alerts.Pushover.Template); msg != "" {
		pushoverErrors = append(pushoverErrors, msg)
	}
	return pushoverErrors
}

func (c *Config) validateNtfy() []string {
	var ntfyErrors []string
	log.Debug().Msg("Ntfy alerting enabled.")
	if c.Alerts.Ntfy.Server == "" {
		ntfyErrors = append(ntfyErrors, "No Ntfy server specified!")
	}
	if c.Alerts.Ntfy.Topic == "" {
		ntfyErrors = append(ntfyErrors, "No Ntfy topic specified!")
	}
	// Check template syntax
	if msg := validateTemplate("Ntfy", c.Alerts.Ntfy.Template); msg != "" {
		ntfyErrors = append(ntfyErrors, msg)
	}

	// Check HTTP header template syntax
	if msg := validateTemplate("Ntfy HTTP Headers", c.Alerts.General.Title); msg != "" {
		ntfyErrors = append(ntfyErrors, msg)
	}

	return ntfyErrors
}

func (c *Config) validateWebhook() []string {
	var webhookErrors []string
	log.Debug().Msg("Webhook alerting enabled.")
	if c.Alerts.Webhook.Server == "" {
		webhookErrors = append(webhookErrors, "No Webhook server specified!")
	}
	// Check HTTP header template syntax
	if msg := validateTemplate("Webhook HTTP Headers", c.Alerts.General.Title); msg != "" {
		webhookErrors = append(webhookErrors, msg)
	}

	return webhookErrors
}

func (c *Config) validateAlertingEnabled() string {
	// Check to ensure at least one alert provider is configured
	if c.Alerts.Discord.Enabled {
		return ""
	}
	if c.Alerts.Gotify.Enabled {
		return ""
	}
	if c.Alerts.SMTP.Enabled {
		return ""
	}
	if c.Alerts.Telegram.Enabled {
		return ""
	}
	if c.Alerts.Pushover.Enabled {
		return ""
	}
	if c.Alerts.Ntfy.Enabled {
		return ""
	}
	if c.Alerts.Webhook.Enabled {
		return ""
	}
	return "No alerting methods have been configured. Please check config file syntax!"
}

func (c *Config) validateAppMonitoring() []string {
	var monitoringErrors []string
	log.Debug().Msg("App monitoring enabled.")
	if c.Monitor.URL == "" {
		monitoringErrors = append(monitoringErrors, "App monitor enabled but no URL specified!")
	}
	if c.Monitor.Interval == 0 {
		c.Monitor.Interval = 60
	}
	return monitoringErrors
}

func validateTemplate(provider, customTemplate string) string {
	var templateError string
	_, err := template.New("").Parse(customTemplate)
	if err != nil {
		templateError = fmt.Sprintf("Failed to parse %s template: %s", provider, err.Error())
	}
	return templateError
}
