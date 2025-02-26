package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/0x2142/frigate-notify/models"
	"github.com/kkyr/fig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App     models.App     `fig:"app" json:"app" required:"false"`
	Frigate models.Frigate `fig:"frigate" json:"frigate" required:"true"`
	Alerts  models.Alerts  `fig:"alerts" json:"alerts" required:"true"`
	Monitor models.Monitor `fig:"monitor" json:"monitor" required:"false"`
}

var ConfigData Config
var ConfigFile string

// Load opens & attempts to parse configuration file
func Load() {
	// Set config file location
	if ConfigFile == "" {
		var ok bool
		ConfigFile, ok = os.LookupEnv("FN_CONFIGFILE")
		if !ok {
			ConfigFile = "./config.yml"
		}
	}

	// Load Config file
	log.Debug().Msgf("Attempting to load config file: %v", ConfigFile)

	err := fig.Load(&ConfigData, fig.File(filepath.Base(ConfigFile)), fig.Dirs(filepath.Dir(ConfigFile)), fig.UseEnv("FN"))
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
	validationErrors := ConfigData.Validate()

	log.Trace().
		Interface("config", ConfigData).
		Msg("Config file loaded & validation completed")

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

func Save(skipBackup bool) {
	log.Debug().Msg("Writing new config file")

	data, err := yaml.Marshal(&ConfigData)
	if err != nil {
		log.Error().Err(err).Msg("Unable to save config")
		return
	}

	// Store backup of original config, if requested
	if !skipBackup {
		original, err := os.Open(ConfigFile)
		if err != nil {
			log.Error().Err(err).Msg("Unable to create config backup")
		}
		defer original.Close()

		newFile := fmt.Sprintf("%s-%s.bak", ConfigFile, time.Now().Format("20060102150405"))
		copy, err := os.Create(newFile)
		if err != nil {
			log.Error().Err(err).Msg("Unable to create config backup")
		}
		defer copy.Close()

		io.Copy(copy, original)
		log.Info().Msgf("Created config file backup: %v", newFile)

	}

	err = os.WriteFile(ConfigFile, data, 0644)
	if err != nil {
		log.Error().Err(err).Msg("Unable to save config")
		return
	}

	log.Info().Msg("Config file saved")
}
