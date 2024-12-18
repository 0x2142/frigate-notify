package config

import (
	"testing"

	"github.com/0x2142/frigate-notify/models"
)

func TestValidateAppMode(t *testing.T) {
	config := Config{Frigate: &models.Frigate{}}

	// Check good config
	config.App.Mode = "reviews"
	Internal.FrigateVersion = 14
	result := config.validateAppMode()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
	// Check bad config
	config.App.Mode = "asdf"
	result = config.validateAppMode()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: error, Got: %v", result)
	}

	// Check incompatible version
	config.App.Mode = "reviews"
	Internal.FrigateVersion = 13
	result = config.validateAppMode()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: error, Got: %v", result)
	}
}

func TestValidateAPI(t *testing.T) {
	config := Config{App: models.App{}}

	config.App.API.Enabled = true

	// Validate default port set
	config.validateAPI()
	if config.App.API.Port != 8000 {
		t.Errorf("Expected: 80, Got: %v", config.App.API.Port)

	}

	// Check good config
	config.App.API.Port = 8080
	result := config.validateAPI()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Check bad config
	config.App.API.Port = 65540
	result = config.validateAPI()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

}

func TestValidateFrigatePolling(t *testing.T) {
	config := Config{Frigate: &models.Frigate{}}

	// Test one method configured
	config.Frigate.MQTT.Enabled = true
	result := config.validateFrigatePolling()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test both methods configured
	config.Frigate.WebAPI.Enabled = true
	result = config.validateFrigatePolling()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: error, Got: %v", result)
	}

	// Test no methods configured
	config.Frigate.WebAPI.Enabled = false
	config.Frigate.MQTT.Enabled = false
	result = config.validateFrigatePolling()
	if len(result) != expected {
		t.Errorf("Expected: error, Got: %v", result)
	}
}

func TestValidateMQTT(t *testing.T) {
	config := Config{Frigate: &models.Frigate{}}

	// Test correct config
	config.Frigate.MQTT.Enabled = true
	config.Frigate.MQTT.Server = "192.0.2.10"
	config.Frigate.MQTT.Username = "test"
	config.Frigate.MQTT.Password = "testddd"
	result := config.validateMQTT()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing server address
	config.Frigate.MQTT.Server = ""
	result = config.validateMQTT()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: err, Got: %v", result)
	}

	// Test user but no password
	config.Frigate.MQTT.Server = "192.0.2.10"
	config.Frigate.MQTT.Password = ""
	result = config.validateMQTT()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: err, Got: %v", result)
	}
}

func TestValidateQuietHours(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}

	// Test valid config
	config.Alerts.Quiet.Start = "03:50"
	config.Alerts.Quiet.End = "04:00"
	result := config.validateQuietHours()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test bad start time
	config.Alerts.Quiet.Start = "03"
	result = config.validateQuietHours()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test bad start & end time
	config.Alerts.Quiet.End = "abc"
	result = config.validateQuietHours()
	expected = 2
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateAlertGeneral(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}

	// Test valid config
	config.Alerts.General.NoSnap = "allow"
	result := config.validateAlertGeneral()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test invalid config
	config.Alerts.General.NoSnap = "something else"
	result = config.validateAlertGeneral()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateDiscord(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Discord = make([]models.Discord, 1)

	// Test valid config
	config.Alerts.Discord[0].Webhook = "https://something.test"
	result := config.validateDiscord(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing webhook config
	config.Alerts.Discord[0].Webhook = ""
	result = config.validateDiscord(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateGotify(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Gotify = make([]models.Gotify, 1)

	// Test valid config
	config.Alerts.Gotify[0].Server = "https://something.test"
	config.Alerts.Gotify[0].Token = "abcdefg"
	result := config.validateGotify(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing server config
	config.Alerts.Gotify[0].Server = ""
	result = config.validateGotify(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing token config
	config.Alerts.Gotify[0].Token = ""
	result = config.validateGotify(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateSMTP(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.SMTP = make([]models.SMTP, 1)

	// Test valid config
	config.Alerts.SMTP[0].Server = "192.0.2.10"
	config.Alerts.SMTP[0].Recipient = "someone@none.test"
	config.Alerts.SMTP[0].User = "someuser"
	config.Alerts.SMTP[0].Password = "abcd"
	result := config.validateSMTP(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Check Default port set
	if config.Alerts.SMTP[0].Port != 25 {
		t.Errorf("Expected: port 25 , Got: %v", config.Alerts.SMTP[0].Port)
	}

	// Check SMTP From is copied
	if config.Alerts.SMTP[0].User != config.Alerts.SMTP[0].From {
		t.Errorf("Expected: %v, Got: %v", config.Alerts.SMTP[0].User, config.Alerts.SMTP[0].From)
	}

	// Test missing server
	config.Alerts.SMTP[0].Server = ""
	result = config.validateSMTP(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing recipient
	config.Alerts.SMTP[0].Recipient = ""
	result = config.validateSMTP(0)
	expected = 2
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test invalid auth config
	config.Alerts.SMTP[0].Password = ""
	result = config.validateSMTP(0)
	expected = 3
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateTelegram(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Telegram = make([]models.Telegram, 1)

	// Test valid config
	config.Alerts.Telegram[0].ChatID = 1234
	config.Alerts.Telegram[0].Token = "abcd"
	result := config.validateTelegram(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing Chat ID
	config.Alerts.Telegram[0].ChatID = 0
	result = config.validateTelegram(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing Token
	config.Alerts.Telegram[0].Token = ""
	result = config.validateTelegram(0)
	expected = 2
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidatePushover(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Pushover = make([]models.Pushover, 1)

	// Test valid config
	config.Alerts.Pushover[0].Token = "abcd"
	config.Alerts.Pushover[0].Userkey = "abcd"
	config.Alerts.Pushover[0].Priority = 1
	result := config.validatePushover(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing token
	config.Alerts.Pushover[0].Token = ""
	result = config.validatePushover(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing Userkey
	config.Alerts.Pushover[0].Userkey = ""
	result = config.validatePushover(0)
	expected = 2
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test priority 2 missing retry / expiration config
	config.Alerts.Pushover[0].Priority = 2
	result = config.validatePushover(0)
	expected = 4
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test priority 2 with low retry interval
	config.Alerts.Pushover[0].Retry = 2
	config.Alerts.Pushover[0].Expire = 10
	result = config.validatePushover(0)
	expected = 3
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test negative TTL
	config.Alerts.Pushover[0].TTL = -2
	result = config.validatePushover(0)
	expected = 4
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateNtfy(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Ntfy = make([]models.Ntfy, 1)

	// Test valid config
	config.Alerts.Ntfy[0].Server = "https://ntfy.test"
	config.Alerts.Ntfy[0].Topic = "frigate"
	result := config.validateNtfy(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing server config
	config.Alerts.Ntfy[0].Server = ""
	result = config.validateNtfy(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing topic config
	config.Alerts.Ntfy[0].Topic = ""
	result = config.validateNtfy(0)
	expected = 2
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}
func TestValidateWebhook(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Webhook = make([]models.Webhook, 1)

	// Test valid config
	config.Alerts.Webhook[0].Server = "https://webhook.test"
	result := config.validateWebhook(0)
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Test missing server config
	config.Alerts.Webhook[0].Server = ""
	result = config.validateWebhook(0)
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

}

func TestValidateAlertingEnabled(t *testing.T) {
	config := Config{Alerts: &models.Alerts{}}
	config.Alerts.Discord = make([]models.Discord, 1)

	// Test valid config
	config.Alerts.Discord[0].Enabled = true
	result := config.validateAlertingEnabled()
	expected := ""
	if result != expected {
		t.Errorf("Expected: '', Got: %v", result)
	}

	// Test missing server config
	config.Alerts.Discord[0].Enabled = false
	result = config.validateAlertingEnabled()
	if result == "" {
		t.Errorf("Expected: error message, Got: %v", result)
	}
}

func TestValidateppMonitoring(t *testing.T) {
	config := Config{}

	// Test valid config
	config.Monitor.URL = "https://something.test"
	result := config.validateAppMonitoring()
	expected := 0
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}

	// Validate default interval
	if config.Monitor.Interval != 60 {
		t.Errorf("Expected: 60, Got: %v", config.Monitor.Interval)
	}

	// Test missing URL
	config.Monitor.URL = ""
	result = config.validateAppMonitoring()
	expected = 1
	if len(result) != expected {
		t.Errorf("Expected: %v error(s), Got: %v", expected, result)
	}
}

func TestValidateTemplate(t *testing.T) {

	// Test valid template
	result := validateTemplate("discord", "{{ .Camera }} detected {{ .Label }}")
	expected := ""
	if result != expected {
		t.Errorf("Expected: '', Got: %v", result)
	}

	// Test invalid template
	result = validateTemplate("discord", "{{ Camera }} detected {{ .Label }}")
	if result == "" {
		t.Errorf("Expected: error message, Got: %v", result)
	}
}
