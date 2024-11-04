package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/0x2142/frigate-notify/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/zerolog/log"
)

var API_PREFIX string

func RunAPIServer() error {
	API_PREFIX = config.ConfigData.App.API.Prefix

	router := http.NewServeMux()

	// Configure API
	apiConfig := huma.DefaultConfig("Frigate-Notify", config.Internal.AppVersion)
	apiConfig.Info.License = &huma.License{Name: "MIT",
		URL: "https://github.com/0x2142/frigate-notify/blob/main/LICENSE"}
	apiConfig.Info.Contact = &huma.Contact{Name: "Matt Schmitz",
		URL: "https://github.com/0x2142/frigate-notify",
	}
	api := humago.New(router, apiConfig)

	registerRoutes(api)

	log.Debug().Msg("Starting API server...")
	listenAddr := fmt.Sprintf("0.0.0.0:%v", config.ConfigData.App.API.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	go http.Serve(listener, router)
	return nil
}

func registerRoutes(api huma.API) {

	// GET /readyz
	huma.Register(api, huma.Operation{
		OperationID: "get-readyz",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/readyz",
		Hidden:      true,
		Summary:     "Readyz",
		Description: "Retrieve Frigate-Notify ready state",
		Tags:        []string{"App"},
	}, GetReadyz)

	// GET /version
	huma.Register(api, huma.Operation{
		OperationID: "get-version",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/version",
		Summary:     "Version",
		Description: "Retrieve Frigate-Notify application version",
		Tags:        []string{"App"},
	}, GetVersion)

	// GET /status
	huma.Register(api, huma.Operation{
		OperationID: "get-status",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/status",
		Summary:     "Status",
		Description: "Retrieve health and status of Frigate-Notify",
		Tags:        []string{"App"},
	}, GetStatus)

	// GET /config
	huma.Register(api, huma.Operation{
		OperationID: "get-config",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/config",
		Summary:     "Config",
		Description: "Retrieve current running configuration",
		Tags:        []string{"App"},
	}, GetConfig)

	// GET /state
	huma.Register(api, huma.Operation{
		OperationID: "get-state",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/state",
		Summary:     "State",
		Description: "Retrieve current state of Frigate-Notify alerting",
		Tags:        []string{"Control"},
	}, GetState)

	// POST /state
	huma.Register(api, huma.Operation{
		OperationID: "post-state",
		Method:      http.MethodPost,
		Path:        API_PREFIX + "/state",
		Summary:     "State",
		Description: "Set state of Frigate-Notify alerting",
		Tags:        []string{"Control"},
	}, PostState)
}
