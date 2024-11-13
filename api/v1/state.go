package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/rs/zerolog/log"
)

type StateInput struct {
	Body struct {
		Enabled bool `json:"enabled" enum:"true,false" doc:"Set state of alerting" required:"true"`
	}
}

type StateOutput struct {
	Body struct {
		Enabled bool `json:"enabled" enum:"true,false" doc:"Frigate-Notify enabled for alerting" default:"true"`
	}
}

// GetState returns whether app is enabled for sending notifications or not
func GetState(ctx context.Context, input *struct{}) (*StateOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/state").
		Str("method", "GET").
		Msg("Received API request")

	resp := &StateOutput{}
	resp.Body.Enabled = config.Internal.Status.Enabled

	log.Trace().
		Str("uri", V1_PREFIX+"/state").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}

// PostState updates state to enable or disable app notifications
func PostState(ctx context.Context, input *StateInput) (*StateOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/state").
		Str("method", "POST").
		Msg("Received API request")

	config.Internal.Status.Enabled = input.Body.Enabled

	log.Debug().
		Bool("state", input.Body.Enabled).
		Msg("App state changed via API")

	resp := &StateOutput{}
	resp.Body.Enabled = config.Internal.Status.Enabled

	log.Trace().
		Str("uri", V1_PREFIX+"/state").
		//Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
