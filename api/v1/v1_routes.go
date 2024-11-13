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
		Description: "Retrieve health and status of Frigate-Notify",
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

	// GET /state
	huma.Register(api, huma.Operation{
		OperationID: "get-state",
		Method:      http.MethodGet,
		Path:        V1_PREFIX + "/state",
		Summary:     V1_PREFIX + "/state",
		Description: "Retrieve current state of Frigate-Notify alerting",
		Tags:        []string{"Control"},
	}, GetState)

	// POST /state
	huma.Register(api, huma.Operation{
		OperationID: "post-state",
		Method:      http.MethodPost,
		Path:        V1_PREFIX + "/state",
		Summary:     V1_PREFIX + "/state",
		Description: "Set state of Frigate-Notify alerting",
		Tags:        []string{"Control"},
	}, PostState)
}
