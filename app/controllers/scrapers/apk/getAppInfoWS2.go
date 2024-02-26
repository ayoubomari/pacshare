package apk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ayoubomari/pacshare/app/models/apkModels"
)

// the return struct of GetApkInfoWS2 function
type GetApkInfoWS2ResponseBody struct {
	Nodes struct {
		Meta struct {
			Data apkModels.ApkInfo `json:"data,omitempty"`
		} `json:"meta,omitempty"`
	} `json:"nodes,omitempty"`
}

func GetApkInfoWS2(appID string) (GetApkInfoWS2ResponseBody, error) {
	var apkInfo GetApkInfoWS2ResponseBody

	url := fmt.Sprintf("http://ws2.aptoide.com/api/7/getApp/app_id/%s", appID)

	resp, err := http.Get(url)
	if err != nil {
		return apkInfo, fmt.Errorf("GetApkInfoWS2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apkInfo, fmt.Errorf("GetApkInfoWS2: HTTP request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apkInfo, fmt.Errorf("GetApkInfoWS2: %w", err)
	}

	err = json.Unmarshal(body, &apkInfo)
	if err != nil {
		return apkInfo, fmt.Errorf("GetApkInfoWS2: %w", err)
	}

	return apkInfo, nil
}
