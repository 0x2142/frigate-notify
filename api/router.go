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

	huma.Register(api, huma.Operation{
		OperationID: "get-version",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/version",
		Summary:     "Version",
		Description: "Retrieve Frigate-Notify application version",
		Tags:        []string{"App"},
	}, GetVersion)

	huma.Register(api, huma.Operation{
		OperationID: "get-status",
		Method:      http.MethodGet,
		Path:        API_PREFIX + "/status",
		Summary:     "Status",
		Description: "Retrieve health and status of Frigate-Notify",
		Tags:        []string{"App"},
	}, GetStatus)

}
