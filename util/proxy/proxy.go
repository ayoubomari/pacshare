package proxy

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/ayoubomari/pacshare/config"
)



var proxiesListUrl = "https://api.proxyscrape.com/v3/free-proxy-list/get?request=displayproxies&protocol=http&proxy_format=protocolipport&format=text&timeout=20000"
func UpdateProxiesFromURL() error {
	resp, err := http.Get(proxiesListUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch proxies: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch proxies: status code %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	var proxies []string

	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		fmt.Println("proxy: ", proxy)
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	config.SetProxies(proxies)
	return nil
}