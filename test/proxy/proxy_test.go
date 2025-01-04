package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ayoubomari/pacshare/config"
)

var proxies = config.Proxies
var proxy = proxies[0]

func TestProxy(t *testing.T) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		t.Fatalf("Invalid proxy URL %s: %v\n", proxy, err)
	}

	// Set up the proxy client
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	// Make a request using the proxy client
	resp, err := client.Get("http://icanhazip.com")
	if err != nil {
		t.Fatalf("Failed to make request through proxy %s: %v\n", proxy, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v\n", err)
	}

	fmt.Printf("Response from proxy %s: %s\n", proxy, string(body))
}
