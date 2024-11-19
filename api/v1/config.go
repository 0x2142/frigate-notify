package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/events"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

type ConfigOutput struct {
	Body struct {
		Config config.Config `json:"config"`
	}
}

type PutConfigInput struct {
	Body struct {
		Config     config.Config `json:"config"`
		SkipSave   bool          `json:"skipsave,omitempty" doc:"Skip writing new config to file" default:"false"`
		SkipBackup bool          `json:"skipbackup,omitempty" doc:"Skip creating config file backup" default:"false"`
	}
}

type PutConfigOutput struct {
	Body struct {
		Status string   `json:"status"`
		Errors []string `json:"errors,omitempty"`
	}
}

// GetConfig returns the current running configuration
func GetConfig(ctx context.Context, input *struct{}) (*ConfigOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/config").
		Str("method", "GET").
		Msg("Received API request")

	resp := &ConfigOutput{}
	resp.Body.Config = config.ConfigData

	log.Trace().
		Str("uri", V1_PREFIX+"/config").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}

// PutConfig replaces the current running configuration
func PutConfig(ctx context.Context, input *PutConfigInput) (*PutConfigOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/config").
		Str("method", "PUT").
		Msg("Received API request")

	resp := &PutConfigOutput{}

	newConfig := input.Body.Config
	validationErrors := newConfig.Validate()

	if len(validationErrors) == 0 {
		resp.Body.Status = "ok"
		go reloadCfg(newConfig, input.Body.SkipSave, input.Body.SkipBackup)

		log.Trace().
			Str("uri", V1_PREFIX+"/config").
			Interface("response_json", resp.Body).
			Msg("Sent API response")
		return resp, nil
	} else {
		resp.Body.Status = "validation error"
		resp.Body.Errors = validationErrors

		log.Trace().
			Str("uri", V1_PREFIX+"/config").
			Interface("response_json", resp.Body).
			Msg("Sent API response")

		return resp, huma.Error422UnprocessableEntity("config validation failed")
	}
}

func reloadCfg(newconfig config.Config, skipSave bool, skipBackup bool) {
	log.Info().Msg("Reloading app config...")
	log.Trace().
		Bool("skipSave", skipSave).
		Bool("skipBackup", skipBackup).
		Msg("Config reload via API")
	if config.ConfigData.Frigate.MQTT.Enabled {
		events.DisconnectMQTT()
	}

	config.ConfigData = newconfig
	if !skipSave {
		config.Save(skipBackup)
	}

	if config.ConfigData.Frigate.MQTT.Enabled {
		events.SubscribeMQTT()
	}
	log.Info().Msg("Config reload completed")
}
