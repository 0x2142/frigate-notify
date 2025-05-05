package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	AuthEnabled     bool
	FrigateServer   string
	FrigateInsecure bool
	FrigateUser     string
	FrigatePass     string
	cookies         *cookiejar.Jar
)

type FrigateAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func init() {
	cookies, _ = cookiejar.New(nil)
}

func GetFrigateVersion(headers []map[string]string) (int, error) {
	url := fmt.Sprintf("%s/api/version", FrigateServer)
	response, err := HTTPGet(url, FrigateInsecure, "", headers...)
	if err != nil {
		return 0, err
	}

	version, _ := strconv.Atoi(strings.Split(string(response), ".")[1])
	return version, nil
}

func checkFrigateAuth() error {
	log.Trace().Msg("Checking Frigate auth token...")
	url := fmt.Sprintf("%s/api/profile", FrigateServer)
	if _, err := HTTPGet(url, FrigateInsecure, ""); err != nil {
		log.Trace().Msg("Frigate auth token expired or not obtained yet")
		if err := getFrigateAuthToken(); err != nil {
			return err
		}
		return nil
	}
	log.Trace().Msg("Frigate auth token still valid")
	return nil
}

func getFrigateAuthToken() error {
	log.Debug().Msg("Authenticating to Frigate...")
	authurl := fmt.Sprintf("%s/api/login", FrigateServer)

	frigate_auth := FrigateAuth{User: FrigateUser, Password: FrigatePass}
	auth_payload, _ := json.Marshal(frigate_auth)

	client := http.Client{Timeout: 10 * time.Second}

	if FrigateInsecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	auth, _ := http.NewRequest(http.MethodPost, authurl, bytes.NewBuffer(auth_payload))

	log.Trace().
		Str("url", authurl).
		Bool("insecure", FrigateInsecure).
		Msg("Attempting authentication")
	response, err := client.Do(auth)
	if err != nil {
		return err
	}
	if response.StatusCode == 401 {
		return errors.New("frigate authentication failed - unauthorized")
	}
	if response.StatusCode == 200 {
		// Save cookies
		log.Debug().Msg("Successfully authenticated to Frigate")
		u, _ := url.Parse(FrigateServer)
		cookies.SetCookies(u, response.Cookies())
		log.Trace().
			Interface("cookies", cookies.Cookies(u)).
			Msg("Saved Frigate cookies")
	}
	return nil
}
