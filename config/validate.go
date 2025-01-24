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

func (c *Config) Validate() []string {
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
	Internal.Status.Notifications.Discord = make([]models.NotifierStatus, len(c.Alerts.Discord))
	for id, profile := range c.Alerts.Discord {
		Internal.Status.Notifications.Discord[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validateDiscord(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
		}
	}

	// Validate Gotify
	Internal.Status.Notifications.Gotify = make([]models.NotifierStatus, len(c.Alerts.Gotify))
	for id, profile := range c.Alerts.Gotify {
		Internal.Status.Notifications.Gotify[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validateGotify(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
		}
	}

	// Validate SMTP
	Internal.Status.Notifications.SMTP = make([]models.NotifierStatus, len(c.Alerts.SMTP))
	for id, profile := range c.Alerts.SMTP {
		Internal.Status.Notifications.SMTP[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validateSMTP(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
		}
	}
	// Validate Telegram
	Internal.Status.Notifications.Telegram = make([]models.NotifierStatus, len(c.Alerts.Telegram))
	for id, profile := range c.Alerts.Telegram {
		Internal.Status.Notifications.Telegram[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validateTelegram(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
		}
	}

	// Validate Pushover
	Internal.Status.Notifications.Pushover = make([]models.NotifierStatus, len(c.Alerts.Pushover))
	for id, profile := range c.Alerts.Pushover {
		Internal.Status.Notifications.Pushover[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validatePushover(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
		}
	}

	// Validate Ntfy
	Internal.Status.Notifications.Ntfy = make([]models.NotifierStatus, len(c.Alerts.Ntfy))
	for id, profile := range c.Alerts.Ntfy {
		Internal.Status.Notifications.Ntfy[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validateNtfy(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
		}
	}

	// Validate Webhook
	Internal.Status.Notifications.Webhook = make([]models.NotifierStatus, len(c.Alerts.Webhook))
	for id, profile := range c.Alerts.Webhook {
		Internal.Status.Notifications.Webhook[id].InitNotifStatus(id, profile.Enabled)
		if profile.Enabled {
			if results := c.validateWebhook(id); len(results) > 0 {
				validationErrors = append(validationErrors, results...)
			}
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
		log.Error().
			Err(err).
			Msgf("Max attempts reached - Cannot reach Frigate server at %v", url)
		connectivityErrors = append(connectivityErrors, "Max attempts reached - Cannot reach Frigate server at "+url)
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

	// Check action on audio-only events
	if strings.ToLower(c.Alerts.General.AudioOnly) != "allow" && strings.ToLower(c.Alerts.General.AudioOnly) != "drop" {
		alertErrors = append(alertErrors, "Option for audio_only must be 'allow' or 'drop'")
	} else {
		log.Debug().Msgf("Audio-only events: %v", strings.ToLower(c.Alerts.General.AudioOnly))
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

func (c *Config) validateDiscord(id int) []string {
	var discordErrors []string
	log.Debug().Msgf("Alerting enabled for Discord profile ID %v", id)
	if c.Alerts.Discord[id].Webhook == "" {
		discordErrors = append(discordErrors, fmt.Sprintf("No Discord webhook specified! Profile ID %v", id))
	}
	// Check template syntax
	if msg := validateTemplate("Discord", c.Alerts.Discord[id].Template); msg != "" {
		discordErrors = append(discordErrors, msg+fmt.Sprintf(" Profile ID %v", id))
	}
	return discordErrors
}

func (c *Config) validateGotify(id int) []string {
	var gotifyErrors []string
	log.Debug().Msgf("Alerting enabled for Gotify profile ID %v", id)
	if c.Alerts.Gotify[id].Server == "" {
		gotifyErrors = append(gotifyErrors, fmt.Sprintf("No Gotify server specified! Profile ID %v", id))
	}
	if c.Alerts.Gotify[id].Token == "" {
		gotifyErrors = append(gotifyErrors, fmt.Sprintf("No Gotify token specified! Profile ID %v", id))
	}
	// Check if Gotify server URL contains protocol, assume HTTP if not specified
	if !strings.Contains(c.Alerts.Gotify[id].Server, "http://") && !strings.Contains(c.Alerts.Gotify[id].Server, "https://") {
		log.Debug().Msgf("No protocol specified on Gotify Server. Assuming http://. If this is incorrect, please adjust the config file. Profile ID %v", id)
		c.Alerts.Gotify[id].Server = fmt.Sprintf("http://%s", c.Alerts.Gotify[id].Server)
	}
	// Check template syntax
	if msg := validateTemplate("Gotify", c.Alerts.Gotify[id].Template); msg != "" {
		gotifyErrors = append(gotifyErrors, msg+fmt.Sprintf(" Profile ID %v", id))
	}
	return gotifyErrors
}

func (c *Config) validateSMTP(id int) []string {
	var smtpErrors []string
	log.Debug().Msgf("Alerting enabled for SMTP profile ID %v", id)
	if c.Alerts.SMTP[id].Server == "" {
		smtpErrors = append(smtpErrors, fmt.Sprintf("No SMTP server specified! Profile ID %v", id))
	}
	if c.Alerts.SMTP[id].Recipient == "" {
		smtpErrors = append(smtpErrors, fmt.Sprintf("No SMTP recipients specified! Profile ID %v", id))
	}
	if c.Alerts.SMTP[id].User != "" && c.Alerts.SMTP[id].Password == "" {
		smtpErrors = append(smtpErrors, fmt.Sprintf("SMTP username in config, but no password provided! Profile ID %v", id))
	}
	if c.Alerts.SMTP[id].Port == 0 {
		c.Alerts.SMTP[id].Port = 25
	}
	// Copy `user` to `from` if `from` not explicitly configured
	if c.Alerts.SMTP[id].From == "" && c.Alerts.SMTP[id].User != "" {
		c.Alerts.SMTP[id].From = c.Alerts.SMTP[id].User
	}
	// Check template syntax
	if msg := validateTemplate("SMTP", c.Alerts.SMTP[id].Template); msg != "" {
		smtpErrors = append(smtpErrors, msg+fmt.Sprintf(" Profile ID %v", id))
	}

	return smtpErrors
}

func (c *Config) validateTelegram(id int) []string {
	var telegramErrors []string
	log.Debug().Msgf("Alerting enabled for Telegram profile ID %v", id)
	if c.Alerts.Telegram[id].ChatID == 0 {
		telegramErrors = append(telegramErrors, fmt.Sprintf("No Telegram Chat ID specified! Profile ID %v", id))
	}
	if c.Alerts.Telegram[id].Token == "" {
		telegramErrors = append(telegramErrors, fmt.Sprintf("No Telegram bot token specified! Profile ID %v", id))
	}
	// Check template syntax
	if msg := validateTemplate("Telegram", c.Alerts.Telegram[id].Template); msg != "" {
		telegramErrors = append(telegramErrors, msg+fmt.Sprintf(" Profile ID %v", id))
	}
	return telegramErrors
}

func (c *Config) validatePushover(id int) []string {
	var pushoverErrors []string
	log.Debug().Msgf("Alerting enabled for Pushover profile ID %v", id)
	if c.Alerts.Pushover[id].Token == "" {
		pushoverErrors = append(pushoverErrors, fmt.Sprintf("No Pushover API token specified! Profile ID %v", id))
	}
	if c.Alerts.Pushover[id].Userkey == "" {
		pushoverErrors = append(pushoverErrors, fmt.Sprintf("No Pushover user key specified! Profile ID %v", id))
	}
	if c.Alerts.Pushover[id].Priority < -2 || c.Alerts.Pushover[id].Priority > 2 {
		pushoverErrors = append(pushoverErrors, fmt.Sprintf("Pushover priority must be between -2 and 2! Profile ID %v", id))
	}
	// Priority 2 is emergency, needs a retry interval & expiration set
	if c.Alerts.Pushover[id].Priority == 2 {
		if c.Alerts.Pushover[id].Retry == 0 || c.Alerts.Pushover[id].Expire == 0 {
			pushoverErrors = append(pushoverErrors, fmt.Sprintf("Pushover retry interval & expiration must be set with priority 2! Profile ID %v", id))
		}
		if c.Alerts.Pushover[id].Retry < 30 {
			pushoverErrors = append(pushoverErrors, fmt.Sprintf("Pushover retry cannot be less than 30 seconds! Profile ID %v", id))
		}
	}
	if c.Alerts.Pushover[id].TTL < 0 {
		pushoverErrors = append(pushoverErrors, fmt.Sprintf("Pushover TTL cannot be negative! Profile ID %v", id))
	}

	// Check template syntax
	if msg := validateTemplate("Pushover", c.Alerts.Pushover[id].Template); msg != "" {
		pushoverErrors = append(pushoverErrors, msg+fmt.Sprintf("Profile ID %v", id))
	}
	return pushoverErrors
}

func (c *Config) validateNtfy(id int) []string {
	var ntfyErrors []string
	log.Debug().Msgf("Alerting enabled for Ntfy profile ID %v", id)
	if c.Alerts.Ntfy[id].Server == "" {
		ntfyErrors = append(ntfyErrors, fmt.Sprintf("No Ntfy server specified! Profile ID %v", id))
	}
	if c.Alerts.Ntfy[id].Topic == "" {
		ntfyErrors = append(ntfyErrors, fmt.Sprintf("No Ntfy topic specified! Profile ID %v", id))
	}
	// Check template syntax
	if msg := validateTemplate("Ntfy", c.Alerts.Ntfy[id].Template); msg != "" {
		ntfyErrors = append(ntfyErrors, msg+fmt.Sprintf("Profile ID %v", id))
	}

	// Check HTTP header template syntax
	if msg := validateTemplate("Ntfy HTTP Headers", c.Alerts.General.Title); msg != "" {
		ntfyErrors = append(ntfyErrors, msg+fmt.Sprintf("Profile ID %v", id))
	}

	return ntfyErrors
}

func (c *Config) validateWebhook(id int) []string {
	var webhookErrors []string
	log.Debug().Msgf("Alerting enabled for Webhook profile ID %v", id)
	if c.Alerts.Webhook[id].Server == "" {
		webhookErrors = append(webhookErrors, fmt.Sprintf("No Webhook server specified! Profile ID %v", id))
	}
	// Check HTTP header template syntax
	if msg := validateTemplate("Webhook HTTP Headers", c.Alerts.General.Title); msg != "" {
		webhookErrors = append(webhookErrors, msg+fmt.Sprintf("Profile ID %v", id))
	}

	return webhookErrors
}

func (c *Config) validateAlertingEnabled() string {
	// Check to ensure at least one alert provider is configured
	for _, profile := range c.Alerts.Discord {
		if profile.Enabled {
			return ""
		}
	}
	for _, profile := range c.Alerts.Gotify {
		if profile.Enabled {
			return ""
		}
	}
	for _, profile := range c.Alerts.SMTP {
		if profile.Enabled {
			return ""
		}
	}
	for _, profile := range c.Alerts.Telegram {
		if profile.Enabled {
			return ""
		}
	}
	for _, profile := range c.Alerts.Pushover {
		if profile.Enabled {
			return ""
		}
	}
	for _, profile := range c.Alerts.Ntfy {
		if profile.Enabled {
			return ""
		}
	}
	for _, profile := range c.Alerts.Webhook {
		if profile.Enabled {
			return ""
		}
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
