package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
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

	resp.Body.Status = config.Internal.Status.Health

	log.Trace().
		Str("uri", V1_PREFIX+"/readyz").
		Interface("response_json", resp.Body).
		Msg("Sent API response")
	return resp, nil

}
