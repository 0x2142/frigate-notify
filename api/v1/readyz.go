package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

type ReadyzOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

// GetReadyz returns current app ready state
func GetReadyz(ctx context.Context, input *struct{}) (*ReadyzOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/readyz").
		Str("method", "GET").
		Msg("Received API request")

	resp := &ReadyzOutput{}

	if config.Internal.Status.Health == "ok" {
		resp.Body.Status = "ok"
		log.Trace().
			Str("uri", V1_PREFIX+"/readyz").
			Interface("response_json", resp.Body).
			Msg("Sent API response")
		return resp, nil
	} else {
		log.Trace().
			Str("uri", V1_PREFIX+"/readyz").
			Int("status_code", 500).
			Msg("Sent API response - App not ready")
		return nil, huma.Error500InternalServerError("app not ready")
	}

}
