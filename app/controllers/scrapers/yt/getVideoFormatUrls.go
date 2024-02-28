package yt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/request"
)

var ErrVideoWayWasDeleted = errors.New("the video may have been removed from YouTube")

// the return type of getVideoFormatsUrls
// string title
// VideoFormatsUrls[videoURL, audioUrl]
type VideoFormatsUrlsAndDetails struct {
	Title       string    `json:"title"`
	FormatsUrls [2]string `json:"formatsUrls"`
}

// response body json struct
type getVideoFormatsUrlsResponseBody struct {
	Title         string `json:"title"`
	FormatStreams []struct {
		Url string `json:"url"`
	} `json:"formatStreams"`
	AdaptiveFormats []struct {
		Url string `json:"url"`
	} `json:"adaptiveFormats"`
}

func getVideoFormatsUrls(videoID string) (VideoFormatsUrlsAndDetails, error) {
	var videoFormatsUrlsAndDetails VideoFormatsUrlsAndDetails

	// use invidios to get the video formats
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("%s/api/v1/videos/%s?fields=formatStreams,adaptiveFormats,title", config.InvidiousEndpoint, videoID),
		nil,
		nil,
	)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("VetvideoFormatsUrls: json requeset err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("VetvideoFormatsUrls: fail to read res.body %w", err)
	}
	var bodyJson getVideoFormatsUrlsResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("VetvideoFormatsUrls: failt to unmarshale the body %w", err)
	}
	if bodyJson.FormatStreams == nil || bodyJson.FormatStreams[0].Url == "" || bodyJson.AdaptiveFormats == nil || bodyJson.AdaptiveFormats[0].Url == "" {
		return videoFormatsUrlsAndDetails, ErrVideoWayWasDeleted
	}

	videoFormatsUrlsAndDetails.Title = bodyJson.Title
	videoFormatsUrlsAndDetails.FormatsUrls[0] = bodyJson.FormatStreams[0].Url   // video/mp4
	videoFormatsUrlsAndDetails.FormatsUrls[1] = bodyJson.AdaptiveFormats[0].Url // // audio/mp4
	return videoFormatsUrlsAndDetails, nil
}
