package apiv1

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var V1_PREFIX string

func Registerv1Routes(api huma.API) {
	V1_PREFIX = "/api/v1"

	// GET /readyz
	huma.Register(api, huma.Operation{
		OperationID: "get-readyz",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/readyz",
		Hidden:      true,
		Summary:     V1_PREFIX + "/readyz",
		Description: "Retrieve Frigate-Notify ready state",
		Tags:        []string{"Status"},
	}, GetReadyz)

	// GET /healthz
	huma.Register(api, huma.Operation{
		OperationID: "get-healthz",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/healthz",
		Hidden:      true,
		Summary:     V1_PREFIX + "/healthz",
		Description: "Retrieve Frigate-Notify liveness state",
		Tags:        []string{"Status"},
	}, GetHealthz)

	// GET /version
	huma.Register(api, huma.Operation{
		OperationID: "get-version",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/version",
		Summary:     V1_PREFIX + "/version",
		Description: "Retrieve Frigate-Notify application version",
		Tags:        []string{"Status"},
	}, GetVersion)

	// GET /status
	huma.Register(api, huma.Operation{
		OperationID: "get-status",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/status",
		Summary:     V1_PREFIX + "/status",
		Description: "Retrieve detailed health and status of Frigate-Notify",
		Tags:        []string{"Status"},
	}, GetStatus)

	// GET /config
	huma.Register(api, huma.Operation{
		OperationID: "get-config",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/config",
		Summary:     V1_PREFIX + "/config",
		Description: "Retrieve current running configuration",
		Tags:        []string{"Config"},
	}, GetConfig)

	// PUT /config
	huma.Register(api, huma.Operation{
		OperationID:   "put-config",
		Method:        http.MethodPut,
		Path:          V1_PREFIX + "/config",
		Summary:       V1_PREFIX + "/config",
		Description:   "Set current running configuration",
		Tags:          []string{"Config"},
		DefaultStatus: http.StatusAccepted,
	}, PutConfig)

	// POST /reload
	huma.Register(api, huma.Operation{
		OperationID:   "post-reload",
		Method:        http.MethodPost,
		Path:          V1_PREFIX + "/reload",
		Summary:       V1_PREFIX + "/reload",
		Description:   "Reload config from file & restart app",
		Tags:          []string{"Control"},
		DefaultStatus: http.StatusAccepted,
	}, PostReload)

	// GET /notif_state
	huma.Register(api, huma.Operation{
		OperationID: "get-notif-state",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/notif_state",
		Summary:     V1_PREFIX + "/notif_state",
		Description: "Retrieve current state of Frigate-Notify notifications",
		Tags:        []string{"Control"},
	}, GetNotifState)

	// POST /notif_state
	huma.Register(api, huma.Operation{
		OperationID:   "post-notif_state",
		Method:        http.MethodPost,
		Path:          V1_PREFIX + "/notif_state",
		Summary:       V1_PREFIX + "/notif_state",
		Description:   "Set state of Frigate-Notify alerting",
		Tags:          []string{"Control"},
		DefaultStatus: http.StatusAccepted,
	}, PostNotifState)

	// POST /notiftest
	huma.Register(api, huma.Operation{
		OperationID:   "post-notiftest",
		Method:        http.MethodPost,
		Path:          V1_PREFIX + "/notif_test",
		Summary:       V1_PREFIX + "/notif_test",
		Description:   "Send test notification via configured providers",
		Tags:          []string{"Control"},
		DefaultStatus: http.StatusAccepted,
	}, PostNotifTest)
}
