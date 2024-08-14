package util

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// HTTPGet is a simple HTTP client function to return page body
func HTTPGet(url string, insecure bool, headers ...map[string]string) ([]byte, error) {

	// New HTTP Client
	client := http.Client{Timeout: 10 * time.Second}
	// Ignore SSL verification if set
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	// Setup new HTTP Request
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	// Add headers
	if len(headers) > 0 {
		for _, h := range headers {
			for k, v := range h {
				req.Header.Add(k, v)
			}

		}
	}

	// Send HTTP GET
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(strconv.Itoa(response.StatusCode))
	}

	return body, nil
}

// HTTPPost performs an HTTP POST to the target URL
// and includes auth parameters, ignoring certificates, etc
func HTTPPost(url string, insecure bool, payload []byte, headers ...map[string]string) ([]byte, error) {
	// New HTTP Client
	client := http.Client{Timeout: 10 * time.Second}

	// Ignore SSL verification if set
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	// Setup new HTTP Request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Add headers
	if len(headers) > 0 {
		for _, h := range headers {
			for k, v := range h {
				req.Header.Add(k, v)
			}

		}
	}

	// Send HTTP POST
	var response *http.Response
	retry := 1
	for retry <= 6 {
		response, err = client.Do(req)
		if err == nil {
			break
		} else {
			if retry == 6 {
				log.Warn().
					Int("attempt", retry).
					Int("max_tries", 6).
					Err(err).
					Msg("HTTP Request failed, retries exceeded")
				break
			}
			log.Warn().
				Int("attempt", retry).
				Int("max_tries", 6).
				Err(err).
				Msg("HTTP Request failed, retrying in 10 seconds...")
			retry += 1
			time.Sleep(1 * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
