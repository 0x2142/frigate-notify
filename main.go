package main

import (
	"flag"
	"os"
	"os/signal"
	"time"
	_ "time/tzdata"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	frigate "github.com/0x2142/frigate-notify/events"
	"github.com/0x2142/frigate-notify/util"
)

var APP_VER = "v0.2.8"
var debug, debugenv bool
var jsonlog, jsonlogenv bool
var configFile string

func main() {
	// Parse flags
	flag.StringVar(&configFile, "c", "", "Configuration file location (default \"./config.yml\")")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&jsonlog, "jsonlog", false, "Enable JSON logging")
	flag.Parse()

	// Set up logging
	_, jsonlogenv = os.LookupEnv("FN_JSONLOG")
	if jsonlog || jsonlogenv {
		zerolog.TimeFieldFormat = "2006/01/02 15:04:05"
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05"})
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Enable debug logging if set
	_, debugenv = os.LookupEnv("FN_DEBUG")
	if debug || debugenv {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	}

	log.Info().Msgf("Frigate Notify - %v", APP_VER)
	log.Info().Msg("Starting...")

	// Load & validate config
	config.LoadConfig(configFile)

	// Set up monitor
	if config.ConfigData.Monitor.Enabled {
		log.Debug().Msg("App monitoring enabled.")
		go func() {
			for {
				_, err := util.HTTPGet(config.ConfigData.Monitor.URL, config.ConfigData.Monitor.Insecure)
				if err != nil {
					log.Warn().
						Err(err).
						Msg("Unable to reach polling monitoring URL")
				}
				log.Debug().Msg("Completed monitoring check-in.")
				time.Sleep(time.Duration(config.ConfigData.Monitor.Interval) * time.Second)
			}
		}()
	}

	// Loop & watch for events
	if config.ConfigData.Frigate.WebAPI.Enabled {
		log.Info().Msg("App running. Press Ctrl-C to quit.")
		for {
			frigate.CheckForEvents()
			time.Sleep(time.Duration(config.ConfigData.Frigate.WebAPI.Interval) * time.Second)
		}
	}
	// Connect MQTT
	if config.ConfigData.Frigate.MQTT.Enabled {
		log.Debug().Msg("Connecting to MQTT Server...")
		frigate.SubscribeMQTT()
		log.Info().Msg("App running. Press Ctrl-C to quit.")
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
	}

}
