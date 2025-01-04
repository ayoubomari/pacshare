package request

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ayoubomari/pacshare/config"
)

// createHTTPClient creates an HTTP client that uses the next proxy from the list
func createHTTPClient() (*http.Client, error) {
	proxyStr := config.GetNextProxy()
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy URL: %v", err)
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}, nil
}

// JSONRequestWithQuery sends a request with JSON body and query strings using a proxy
func JSONRequestWithQuery(method string, url string, jsonBytes []byte, headers map[string]string, queryParams map[string]string, useProxy bool) (*http.Response, error) {
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
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create the HTTP client with proxy
	client := &http.Client{}
	if useProxy {
		client, err = createHTTPClient()
		if err != nil {
			return nil, err
		}
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}

// JSONReqest sends a request with JSON body using a proxy
func JSONReqest(method string, url string, jsonBytes []byte, headers map[string]string, useProxy bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create the HTTP client with proxy
	client := &http.Client{}
	if useProxy {
		client, err = createHTTPClient()
		if err != nil {
			return nil, err
		}
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}

// GetContentLengthFromResponseHeader gets the content length from the response header using a HEAD request and a proxy
func GetContentLengthFromResponseHeader(url string, headers map[string]string, useProxy bool) (int, error) {
	// Create the HTTP client with proxy
	client := &http.Client{}
	var err error
	if useProxy {
		client, err = createHTTPClient()
		if err != nil {
			return 0, err
		}
	}

	// Create a new HEAD request
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	fmt.Printf("int(resp.ContentLength): %d\n", int(resp.ContentLength))
	return int(resp.ContentLength), nil
}

// GetContentLengthWithGetReq gets the content length using a GET request and a proxy
func GetContentLengthWithGetReq(url string, useProxy bool) (int, error) {
	// Create the HTTP client with proxy
	client := &http.Client{}
	var err error
	if useProxy {
		client, err = createHTTPClient()
		if err != nil {
			return 0, err
		}
	}

	// Make the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("non-OK status code: %d", resp.StatusCode)
	}

	return int(resp.ContentLength), nil
}

// GetRedirectLocation gets the redirect location using a proxy
func GetRedirectLocation(url string, useProxy bool) (string, error) {
	// Create the HTTP client with proxy
	client := &http.Client{}
	var err error
	if useProxy {
		client, err = createHTTPClient()
		if err != nil {
			return "", err
		}
	}

	// Make the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		redirectURL, err := resp.Location()
		if err != nil {
			return "", err
		}
		return redirectURL.String(), nil
	}

	return "", nil
}
