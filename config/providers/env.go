// Custom environment variable provider
// Needed in order to support multiple notification provider profiles config via env vars
// Thanks to https://github.com/knadh/koanf/issues/74#issuecomment-1778885434
package envconfig

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/sjson"
)

type Env struct {
	prefix string
	delim  string
	cb     func(key string, value string) (string, interface{})
	out    string
}

func Provider(prefix, delim string, cb func(s string) string) *Env {
	e := &Env{
		prefix: prefix,
		delim:  delim,
		out:    "{}",
	}
	if cb != nil {
		e.cb = func(key string, value string) (string, interface{}) {
			return cb(key), value
		}
	}
	return e
}

func ProviderWithValue(prefix, delim string, cb func(key string, value string) (string, interface{})) *Env {
	return &Env{
		prefix: prefix,
		delim:  delim,
		cb:     cb,
	}
}

func (e *Env) Read() (map[string]interface{}, error) {
	return nil, errors.New("provider does not support this method")
}

func (e *Env) ReadBytes() ([]byte, error) {
	var keys []string
	for _, k := range os.Environ() {
		if e.prefix != "" {
			if strings.HasPrefix(k, e.prefix) {
				keys = append(keys, k)
			}
		} else {
			keys = append(keys, k)
		}
	}

	if len(keys) == 0 {
		return nil, errors.New("no app environment variables found")
	}
	log.Debug().Msgf("Found %v environment variable(s)", len(keys))

	for _, k := range keys {
		parts := strings.SplitN(k, "=", 2)

		var (
			key   string
			value interface{}
		)

		if e.cb != nil {
			key, value = e.cb(parts[0], parts[1])
			if key == "" {
				continue
			}
		} else {
			key = parts[0]
			value = parts[1]
		}

		if err := e.set(key, value); err != nil {
			return []byte{}, err
		}
	}

	return []byte(e.out), nil
}

func (e *Env) set(key string, value interface{}) error {
	out, err := sjson.Set(e.out, strings.Replace(key, e.delim, ".", -1), value)
	if err != nil {
		return err
	}
	e.out = out

	return nil
}
