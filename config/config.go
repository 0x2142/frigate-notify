package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/0x2142/frigate-notify/models"
	"github.com/kkyr/fig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	App     models.App      `fig:"app" json:"app"`
	Frigate *models.Frigate `fig:"frigate" json:"frigate" validate:"required"`
	Alerts  *models.Alerts  `fig:"alerts" json:"alerts" validate:"required"`
	Monitor models.Monitor  `fig:"monitor" json:"monitor"`
}

var ConfigData Config

// loadConfig opens & attempts to parse configuration file
func LoadConfig(configFile string) {
	// Set config file location
	if configFile == "" {
		var ok bool
		configFile, ok = os.LookupEnv("FN_CONFIGFILE")
		if !ok {
			configFile = "./config.yml"
		}
	}

	// Load Config file
	log.Debug().Msgf("Attempting to load config file: %v", configFile)

	err := fig.Load(&ConfigData, fig.File(filepath.Base(configFile)), fig.Dirs(filepath.Dir(configFile)), fig.UseEnv("FN"))
	if err != nil {
		if errors.Is(err, fig.ErrFileNotFound) {
			log.Warn().Msg("Config file could not be read, attempting to load config from environment")
			err = fig.Load(&ConfigData, fig.IgnoreFile(), fig.UseEnv("FN"))
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("Failed to load config from environment!")
			}
		} else {
			log.Fatal().
				Err(err).
				Msg("Failed to load config from file!")
		}
	}
	log.Info().Msg("Config loaded.")

	// Send config file to validation before completing
	validationErrors := ConfigData.validate()

	if len(validationErrors) > 0 {
		fmt.Println()
		log.Error().Msg("Config validation failed:")
		for _, msg := range validationErrors {
			log.Error().Msgf(" - %v", msg)
		}
		fmt.Println()
		log.Fatal().Msg("Please fix config errors before restarting app.")
	} else {
		log.Info().Msg("Config file validated!")
	}
}
