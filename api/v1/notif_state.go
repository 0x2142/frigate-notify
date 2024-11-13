package apiv1

import (
	"context"

	"github.com/0x2142/frigate-notify/config"
	"github.com/rs/zerolog/log"
)

type NotifStateInput struct {
	Body struct {
		Enabled bool `json:"enabled" enum:"true,false" doc:"Set state of notifications" required:"true"`
	}
}

type NotifStateOutput struct {
	Body struct {
		Enabled bool `json:"enabled" enum:"true,false" doc:"Frigate-Notify enabled for notifications" default:"true"`
	}
}

// GetNotifState returns whether app is enabled for sending notifications or not
func GetNotifState(ctx context.Context, input *struct{}) (*NotifStateOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/notif_state").
		Str("method", "GET").
		Msg("Received API request")

	resp := &NotifStateOutput{}
	resp.Body.Enabled = config.Internal.Status.Notifications.Enabled

	log.Trace().
		Str("uri", V1_PREFIX+"/notif_state").
		Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}

// PostNotifState updates state to enable or disable app notifications
func PostNotifState(ctx context.Context, input *NotifStateInput) (*NotifStateOutput, error) {
	log.Trace().
		Str("uri", V1_PREFIX+"/notif_state").
		Str("method", "POST").
		Msg("Received API request")

	config.Internal.Status.Notifications.Enabled = input.Body.Enabled

	log.Debug().
		Bool("state", input.Body.Enabled).
		Msg("App state changed via API")

	resp := &NotifStateOutput{}
	resp.Body.Enabled = config.Internal.Status.Notifications.Enabled

	log.Trace().
		Str("uri", V1_PREFIX+"/notif_state").
		//Interface("response_json", resp.Body).
		Msg("Sent API response")

	return resp, nil
}
