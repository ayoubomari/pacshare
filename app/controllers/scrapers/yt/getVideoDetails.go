package yt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

var (
	ErrVideoDetailsIsNil = errors.New("video details is nil")
	ErrInvalidResponse   = errors.New("invalid response from YouTube API")
)

type VideoDetails struct {
	Title             string `json:"title"`
	DurationInSeconds int    `json:"durationInSeconds"`
	Description       string `json:"description"`
	Thumbnail         string `json:"Thumbnail"`
	UploadDate        string `json:"uploadDate"`
	Author            string `json:"author"`
	ViewCount         string `json:"viewCount,omitempty"`
}

type getVideoDetailsResponseBody struct {
	VideoDetails struct {
		Author           string `json:"author,omitempty"`
		ViewCount        string `json:"viewCount,omitempty"`
		ShortDescription string `json:"shortDescription,omitempty"`
		LengthSeconds    string `json:"lengthSeconds,omitempty"`
		Title            string `json:"title,omitempty"`
	} `json:"videoDetails"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			UploadDate string `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
}

func getVideoDetails(videoID string) (VideoDetails, error) {
	var videoDetails VideoDetails

	if videoID == "" {
		return videoDetails, errors.New("video ID is empty")
	}

	jsonBytes := []byte(`{"context":{"client":{"hl":"en","gl":"MA","remoteHost":"41.141.106.74","deviceMake":"","deviceModel":"","visitorData":"CgtJNUVzR3M1NHR6SSijk_WLBg%3D%3D","userAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36,gzip(gfe)","clientName":"WEB","clientVersion":"2.20211029.00.00","osName":"Windows","osVersion":"10.0","originalUrl":"https://www.youtube.com/watch?v=` + videoID + `","platform":"DESKTOP","clientFormFactor":"UNKNOWN_FORM_FACTOR","configInfo":{"appInstallData":"CKOT9YsGEN_WrQUQsNStBRC3y60FEJLXrQUQktWtBRDz460FENi-rQUQkfj8Eg%3D%3D"},"browserName":"Chrome","browserVersion":"95.0.4638.54","screenWidthPoints":1440,"screenHeightPoints":241,"screenPixelDensity":1,"screenDensityFloat":1,"utcOffsetMinutes":60,"userInterfaceTheme":"USER_INTERFACE_THEME_LIGHT","clientScreen":"WATCH","mainAppWebInfo":{"graftUrl":"/watch?v=` + videoID + `","webDisplayMode":"WEB_DISPLAY_MODE_BROWSER","isWebNativeShareAvailable":true},"timeZone":"Africa/Casablanca"},"user":{"lockedSafetyMode":false},"request":{"useSsl":true,"internalExperimentFlags":[],"consistencyTokenJars":[]},"clickTracking":{"clickTrackingParams":"COsBENwwGAAiEwi5pdX8ofLzAhVtQk8EHSd0C2kyBnNlYXJjaFISaXRhbGlhIHVsdHJhIGJlYXRzmgEDEPQk"},"adSignalsInfo":{"params":[{"key":"dt","value":"1635600806944"},{"key":"flash","value":"0"},{"key":"frm","value":"0"},{"key":"u_tz","value":"60"},{"key":"u_his","value":"7"},{"key":"u_h","value":"900"},{"key":"u_w","value":"1440"},{"key":"u_ah","value":"860"},{"key":"u_aw","value":"1440"},{"key":"u_cd","value":"24"},{"key":"bc","value":"31"},{"key":"bih","value":"241"},{"key":"biw","value":"1424"},{"key":"brdim","value":"0,0,0,0,1440,0,1440,860,1440,241"},{"key":"vis","value":"1"},{"key":"wgl","value":"true"},{"key":"ca_type","value":"image"}]}},"videoId":"` + videoID + `","playbackContext":{"contentPlaybackContext":{"currentUrl":"/watch?v=` + videoID + `","vis":0,"splay":false,"autoCaptionsDefaultOn":false,"autonavState":"STATE_NONE","html5Preference":"HTML5_PREF_WANTS","signatureTimestamp":18927,"referer":"https://www.youtube.com/results?search_query=morocco","lactMilliseconds":"-1"}},"racyCheckOk":false,"contentCheckOk":false}`)

	res, err := request.JSONReqest(
		"POST",
		"https://www.youtube.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false",
		jsonBytes,
		nil,
		false,
	)
	if err != nil {
		return videoDetails, fmt.Errorf("getVideoDetails: JSON request error: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return videoDetails, fmt.Errorf("getVideoDetails: failed to read response body: %w", err)
	}

	var bodyJson getVideoDetailsResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return videoDetails, fmt.Errorf("getVideoDetails: failed to unmarshal the body: %w", err)
	}

	if bodyJson.VideoDetails.Title == "" || bodyJson.Microformat.PlayerMicroformatRenderer.UploadDate == "" {
		return videoDetails, ErrInvalidResponse
	}

	durationInSeconds, err := strconv.Atoi(bodyJson.VideoDetails.LengthSeconds)
	if err != nil {
		durationInSeconds = 0
	}

	description := bodyJson.VideoDetails.ShortDescription
	if len(description) > 2000 {
		description = description[:2000]
	}

	viewCount, err := strconv.ParseInt(bodyJson.VideoDetails.ViewCount, 10, 64)
	if err != nil {
		viewCount = 0
	}

	uploadDateParts := strings.Split(bodyJson.Microformat.PlayerMicroformatRenderer.UploadDate, "T")
	uploadDate := ""
	if len(uploadDateParts) > 0 {
		uploadDate = uploadDateParts[0]
	}

	videoDetails = VideoDetails{
		Title:             bodyJson.VideoDetails.Title,
		DurationInSeconds: durationInSeconds,
		Description:       description,
		Thumbnail:         fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", videoID),
		UploadDate:        uploadDate,
		Author:            bodyJson.VideoDetails.Author,
		ViewCount:         formats.FormatNumberWithCommas(viewCount),
	}

	return videoDetails, nil
}