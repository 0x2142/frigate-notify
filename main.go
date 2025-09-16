package main

import (
	"embed"
	"flag"
	"io"
	"os"
	"os/signal"
	debuginfo "runtime/debug"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/0x2142/frigate-notify/api"
	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/events"
	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
)

var debug, debugenv bool
var jsonlog, jsonlogenv bool
var nocolor, nocolorenv bool
var configFile string
var logLevel string

//go:embed templates/*
var NotifTemplates embed.FS

func main() {
	config.Internal.Status.Health = "starting"

	// Parse flags
	flag.StringVar(&configFile, "c", "", "Configuration file location (default \"./config.yml\")")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging (Overrides loglevel, if also set)")
	flag.StringVar(&logLevel, "loglevel", "", "Set logging level")
	flag.BoolVar(&jsonlog, "jsonlog", false, "Enable JSON logging")
	flag.BoolVar(&nocolor, "nocolor", false, "Disable color on console logging")
	flag.Parse()

	// Set up logging
	var logwriter io.Writer = os.Stdout
	logfile := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		LocalTime:  true,
	}

	_, jsonlogenv = os.LookupEnv("FN_JSONLOG")
	if !jsonlog && !jsonlogenv {
		_, nocolorenv = os.LookupEnv("FN_NOCOLOR")
		logwriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05 -0700", NoColor: nocolorenv || nocolor}
	}

	zerolog.TimeFieldFormat = "2006/01/02 15:04:05 -0700"
	log.Logger = zerolog.New(zerolog.MultiLevelWriter(logwriter, logfile)).With().Timestamp().Logger()

	// Apply custom log level, if set
	if logLevel == "" {
		logLevel, _ = os.LookupEnv("FN_LOGLEVEL")
	}
	switch strings.ToLower(logLevel) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "debug":
		log.Debug().Msg("Debug logging enabled")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
	case "trace":
		log.Trace().Msg("Trace logging enabled")
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Logger = log.With().Caller().Logger()
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Enable debug logging if set
	_, debugenv = os.LookupEnv("FN_DEBUG")
	if debug || debugenv {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		// Add calling module to debug logs
		log.Logger = log.With().Caller().Logger()
		log.Debug().Msg("Debug logging enabled")
	}

	info, _ := debuginfo.ReadBuildInfo()
	buildinfo := make(map[string]interface{})
	for _, s := range info.Settings {
		buildinfo[s.Key] = s.Value
	}
	log.Info().Msgf("Frigate Notify - %v", config.Internal.AppVersion)
	log.Trace().Fields(buildinfo).Msg("Build Info")
	log.Info().Msg("Starting...")

	// Load & validate config
	config.ConfigFile = configFile
	config.Load()

	notifier.TemplateFiles = NotifTemplates

	// Set up monitor
	if config.ConfigData.Monitor.Enabled {
		log.Debug().Msg("App monitoring enabled.")
		go func() {
			for {
				_, err := util.HTTPGet(config.ConfigData.Monitor.URL, config.ConfigData.Monitor.Insecure, "")
				if err != nil {
					config.Internal.Status.Monitor = err.Error()
					log.Warn().
						Err(err).
						Msg("Unable to reach polling monitoring URL")
				}
				config.Internal.Status.Monitor = "ok"
				log.Debug().Msg("Completed monitoring check-in.")
				time.Sleep(time.Duration(config.ConfigData.Monitor.Interval) * time.Second)
			}
		}()
	}

	// Set up event cache
	events.InitZoneCache()
	defer events.CloseZoneCache()

	// Start API server if enabled
	if config.ConfigData.App.API.Enabled {
		err := api.RunAPIServer()
		if err != nil {
			config.Internal.Status.API = err.Error()
			log.Error().Err(err).Msg("Failed to start API server")
		} else {
			config.Internal.Status.API = "ok"
			log.Info().Msgf("API server ready on :%v", config.ConfigData.App.API.Port)
		}
	}

	// Loop & watch for events
	if config.ConfigData.Frigate.WebAPI.Enabled {
		log.Info().Msg("App ready!")
		config.Internal.Status.Health = "ok"
		for {
			events.QueryAPI()
			time.Sleep(time.Duration(config.ConfigData.Frigate.WebAPI.Interval) * time.Second)
		}
	}

	// Connect MQTT
	if config.ConfigData.Frigate.MQTT.Enabled {
		log.Debug().Msg("Connecting to MQTT Server...")
		events.SubscribeMQTT()
		defer events.DisconnectMQTT()
		log.Info().Msg("App ready!")
		config.Internal.Status.Health = "ok"
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
	}

}
