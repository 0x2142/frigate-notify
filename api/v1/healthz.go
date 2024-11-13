package apiv1

import (
	"context"

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

	resp.Body.Status = "ok"

	log.Trace().
		Str("uri", V1_PREFIX+"/healthz").
		Interface("response_json", resp.Body).
		Msg("Sent API response")
	return resp, nil

}
