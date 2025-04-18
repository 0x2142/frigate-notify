package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/rs/zerolog/log"
)

type VersionOutput struct {
	Body struct {
		Version string `json:"version" example:"v0.0.0" doc:"Current version of Frigate-Notify"`
	}
}

// GetVersion returns the current running app version
func GetVersion(ctx context.Context, input *struct{}) (*VersionOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/version").
		Str("method", "GET").
		Msg("Received API request")

	resp := &VersionOutput{}
	resp.Body.Version = config.Internal.AppVersion

	log.Trace().
		Str("uri", V1_PREFIX+"/version").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
