package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/rs/zerolog/log"
)

type HealthzOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

// GetHealthz returns current app liveness state
func GetHealthz(ctx context.Context, input *struct{}) (*HealthzOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/healthz").
		Str("method", "GET").
		Msg("Received API request")

	resp := &HealthzOutput{}

	if config.Internal.Status.Health == "ok" {
		resp.Body.Status = "ok"
		log.Trace().
			Str("uri", V1_PREFIX+"/healthz").
			Interface("response_json", resp.Body).
			Msg("Sent API response")
		return resp, nil
	} else {
		log.Trace().
			Str("uri", V1_PREFIX+"/healthz").
			Int("status_code", 500).
			Msg("Sent API response")
		return nil, huma.Error500InternalServerError("app not ready")
	}

}
