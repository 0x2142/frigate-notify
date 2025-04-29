package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	envconfig "github.com/0x2142/frigate-notify/config/providers/env"
	secretsconfig "github.com/0x2142/frigate-notify/config/providers/secrets"
	"github.com/0x2142/frigate-notify/models"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	yml "gopkg.in/yaml.v3"
)

type Config struct {
	App     models.App     `koanf:"app" json:"app" required:"false"`
	Frigate models.Frigate `koanf:"frigate" json:"frigate" required:"true"`
	Alerts  models.Alerts  `koanf:"alerts" json:"alerts" required:"true"`
	Monitor models.Monitor `koanf:"monitor" json:"monitor" required:"false"`
}

var ConfigData Config
var ConfigFile string
var k = koanf.New(".")

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

	// Set config defaults
	k.Load(structs.Provider(DefaultConfig, "koanf"), nil)

	// Attempt to load config from file
	log.Debug().Msgf("Loading config from file: %v", ConfigFile)
	if err := k.Load(file.Provider(ConfigFile), yaml.Parser()); err != nil {
		log.Warn().Msg("Unable to load config from file")
	}

	// Attempt to load config from env var
	log.Debug().Msg("Checking for environment variables")
	if err := k.Load(envconfig.ProviderWithValue("FN_", "__", processENV), json.Parser()); err != nil {
		log.Debug().
			Err(err).
			Msg("Unable to load environment variables")
	}

	// Attempt to load config from docker secrets
	log.Debug().Msg("Checking for docker secrets")
	if err := k.Load(secretsconfig.ProviderWithValue("FN_", "__", processENV), json.Parser()); err != nil {
		log.Debug().
			Err(err).
			Msg("Unable to load docker secrets")
	}

	k.Unmarshal("", &ConfigData)

	log.Info().Msg("Config loaded")

	log.Trace().
		Interface("config", ConfigData).
		Msg("Config loaded")

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

	data, err := yml.Marshal(&ConfigData)
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

func processENV(s, v string) (string, interface{}) {
	key := strings.ToLower(strings.TrimPrefix(s, "FN_"))

	// Split multiple values separated by semicolon into slice
	// for example FN_FRIGATE__CAMERAS__EXCLUDE="camera1;camera2"
	if strings.Contains(v, ";") {
		if strings.Contains(strings.ToLower(key), "headers") {
			var headers []map[string]string
			for _, header := range strings.Split(v, ";") {
				split := strings.Split(header, ":")
				if len(split) != 2 {
					continue
				}
				newHeader := map[string]string{split[0]: split[1]}
				headers = append(headers, newHeader)
			}
			return key, headers
		}
		return key, strings.Split(v, ";")
	}
	return key, v
}
