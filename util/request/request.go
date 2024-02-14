package request

import (
	"bytes"
	"fmt"
	"net/http"
)

// send a request with json body bytes body and query strings
func JSONRequestWithQuery(method string, url string, jsonBytes []byte, headers map[string]string, queryParams map[string]string) (*http.Response, error) {
	// Create URL with query parameters
	fullURL := url
	if len(queryParams) > 0 {
		query := ""
		for key, value := range queryParams {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", key, value)
		}
		fullURL += "?" + query
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the Content-Type header to application/json
	for k, p := range headers {
		req.Header.Set(k, p)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}

	return resp, nil
}

// send a request with json body bytes
// Note: you have to include query strings into url string
func JSONReqest(method string, url string, jsonBytes []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Set the Content-Type header to application/json
	for k, p := range headers {
		req.Header.Set(k, p)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}

	return resp, nil
}
