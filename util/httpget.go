package util

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

// HTTPGet is a simple HTTP client function to return page body
func HTTPGet(url string, insecure bool) ([]byte, error) {

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

	return body, nil
}
