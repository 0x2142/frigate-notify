// Custom docker secrets provider
package secretsconfig

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
	// Read in docker secrets
	files, err := os.ReadDir("/run/secrets")
	if err != nil {
		return nil, err
	}
	for _, k := range files {
		if e.prefix != "" {
			if strings.HasPrefix(k.Name(), e.prefix) {
				keys = append(keys, k.Name())
			}
		} else {
			keys = append(keys, k.Name())
		}
	}

	if len(keys) == 0 {
		return nil, errors.New("no docker secrets found")
	}
	log.Debug().Msgf("Found %v docker secret(s)", len(keys))

	for _, k := range keys {
		v, err := os.ReadFile("/run/secrets/" + k)

		if err != nil {
			continue
		}

		var (
			key   string
			value interface{}
		)

		sv := strings.TrimSuffix(string(v), "\n")

		if e.cb != nil {
			key, value = e.cb(k, sv)
			if key == "" {
				continue
			}
		} else {
			key = k
			value = sv
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
