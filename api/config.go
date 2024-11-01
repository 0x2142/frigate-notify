package api

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/rs/zerolog/log"
)

type ConfigOutput struct {
	Body struct {
		Config config.Config `json:"config"`
	}
}

// GetConfig returns the current running configuratio
func GetConfig(ctx context.Context, input *struct{}) (*ConfigOutput, error) {
	log.Trace().
		Str("uri", API_PREFIX+"/config").
		Str("method", "GET").
		Msg("Received API request")

	resp := &ConfigOutput{}
	resp.Body.Config = config.ConfigData

	log.Trace().
		Str("uri", API_PREFIX+"/config").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
