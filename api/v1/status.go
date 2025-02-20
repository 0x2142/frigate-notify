package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/rs/zerolog/log"
)

type StatusOutput struct {
	Body struct {
		Status models.Status `json:"status"`
	}
}

// GetStatus returns current app component statuses
func GetStatus(ctx context.Context, input *struct{}) (*StatusOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/status").
		Str("method", "GET").
		Msg("Received API request")

	resp := &StatusOutput{}
	resp.Body.Status = config.Internal.Status

	log.Trace().
		Str("uri", V1_PREFIX+"/status").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
