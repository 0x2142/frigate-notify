package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/rs/zerolog/log"
)

type ReloadOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// PostReload reloads current app config
func PostReload(ctx context.Context, input *struct{}) (*ReloadOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/reload").
		Str("method", "POST").
		Msg("Received API request")

	resp := &ReloadOutput{}
	resp.Body.Message = "ok"

	go func() {
		log.Info().Msg("Received request to reload config")
		// Re-load from file & trigger reload
		config.Load()
		newconfig := config.ConfigData
		go reloadCfg(newconfig, true, true)
	}()

	log.Trace().
		Str("uri", V1_PREFIX+"/reload").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
