package api

import (
	"fmt"
	"net"
	"net/http"

	apiv1 "github.com/0x2142/frigate-notify/api/v1"
	"github.com/0x2142/frigate-notify/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/zerolog/log"
)

func RunAPIServer() error {

	router := http.NewServeMux()

	// Configure API
	apiConfig := huma.DefaultConfig("Frigate-Notify", config.Internal.AppVersion)
	apiConfig.Info.License = &huma.License{Name: "MIT",
		URL: "https://github.com/0x2142/frigate-notify/blob/main/LICENSE"}
	apiConfig.Info.Contact = &huma.Contact{Name: "Matt Schmitz",
		URL: "https://github.com/0x2142/frigate-notify",
	}
	api := humago.New(router, apiConfig)

	apiv1.Registerv1Routes(api)

	log.Debug().Msg("Starting API server...")
	listenAddr := fmt.Sprintf("0.0.0.0:%v", config.ConfigData.App.API.Port)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	go http.Serve(listener, router)
	return nil
}
