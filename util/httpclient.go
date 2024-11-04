package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// buildParams creates an escaped param string from a slice
func BuildHTTPParams(params ...map[string]string) string {
	var paramList string
	if len(params) > 0 {
		paramList = "?"
		for _, h := range params {
			for k, v := range h {
				k = url.QueryEscape(k)
				v = url.QueryEscape(v)
				paramList = fmt.Sprintf("%s&%s=%s", paramList, k, v)
			}

		}
	}

	return paramList
}

// HTTPGet is a simple HTTP client function to return page body
func HTTPGet(url string, insecure bool, params string, headers ...map[string]string) ([]byte, error) {
	// Append HTTP params if any
	if len(params) > 0 {
		url = url + params
	}

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

	// Remove authorization header value for logging
	for i, h := range headers {
		for k := range h {
			if strings.ToLower(k) == "authorization" {
				headers[i][k] = "--secret removed--"
			}
		}
	}
	// Send HTTP GET
	log.Trace().
		Str("url", url).
		Interface("headers", headers).
		Bool("insecure", insecure).
		Msg("HTTP GET")
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

	// Skip logging contents of snapshot image
	if strings.Contains(url, "snapshot.jpg") {
		log.Trace().
			Int64("content_length", response.ContentLength).
			Int("status_code", response.StatusCode).
			Msg("HTTP Response")
	} else {
		log.Trace().
			RawJSON("body", body).
			Int("status_code", response.StatusCode).
			Msg("HTTP Response")
	}

	if response.StatusCode != 200 {
		return nil, errors.New(strconv.Itoa(response.StatusCode))
	}

	return body, nil
}

// HTTPPost performs an HTTP POST to the target URL
// and includes auth parameters, ignoring certificates, etc
func HTTPPost(url string, insecure bool, payload []byte, params string, headers ...map[string]string) ([]byte, error) {
	// Append HTTP params if any
	if len(params) > 0 {
		url = url + params
	}

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

	// Remove authorization header value for logging
	for i, h := range headers {
		for k := range h {
			if strings.ToLower(k) == "authorization" {
				headers[i][k] = "--secret removed--"
			}
		}
	}

	// Send HTTP POST
	if json.Valid(payload) {
		log.Trace().
			Str("url", url).
			Interface("headers", headers).
			RawJSON("body", payload).
			Bool("insecure", insecure).
			Msg("HTTP POST")
	} else {
		log.Trace().
			Str("url", url).
			Interface("headers", headers).
			Interface("body", payload).
			Bool("insecure", insecure).
			Msg("HTTP POST")
	}

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

	log.Trace().
		RawJSON("body", body).
		Int("status_code", response.StatusCode).
		Msg("HTTP Response")

	// Check status codes
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("failed to send request, got status code %v", response.StatusCode)
	}

	return body, nil
}
