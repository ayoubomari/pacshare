package yt

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

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
		false,
	)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("VetvideoFormatsUrls: json requeset err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	// fmt.Println(string(bodyBytes))
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

// XML structures
type MPD struct {
	XMLName xml.Name `xml:"MPD"`
	Period  Period   `xml:"Period"`
}

type Period struct {
	AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
	MimeType        string           `xml:"mimeType,attr"`
	Representations []Representation `xml:"Representation"`
}

type Representation struct {
	ID        string `xml:"id,attr"`
	Codecs    string `xml:"codecs,attr"`
	BaseURL   string `xml:"BaseURL"`
	Width     string `xml:"width,attr,omitempty"`
	Height    string `xml:"height,attr,omitempty"`
	Bandwidth string `xml:"bandwidth,attr"`
}

func getVideoFormatsUrlsByXML(videoID string) (VideoFormatsUrlsAndDetails, error) {
	var videoFormatsUrlsAndDetails VideoFormatsUrlsAndDetails

	// Make request to get the DASH manifest
	fmt.Printf("%s/api/manifest/dash/id/%s?local=true&unique_res=1", config.InvidiousEndpoint, videoID)
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("%s/api/manifest/dash/id/%s?local=true&unique_res=1", config.InvidiousEndpoint, videoID),
		nil,
		nil,
		false,
	)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("GetVideoFormatsUrlsByXML: request err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("GetVideoFormatsUrlsByXML: fail to read res.body %w", err)
	}

	var mpd MPD
	err = xml.Unmarshal(bodyBytes, &mpd)
	if err != nil {
		return videoFormatsUrlsAndDetails, fmt.Errorf("GetVideoFormatsUrlsByXML: fail to unmarshal XML %w", err)
	}

	// Find best quality video and audio URLs
	var videoURL, audioURL string
	var maxVideoBandwidth, maxAudioBandwidth int64

	for _, adaptationSet := range mpd.Period.AdaptationSets {
		switch {
		case strings.HasPrefix(adaptationSet.MimeType, "video/"):
			// Find highest quality video
			for _, rep := range adaptationSet.Representations {
				bandwidth, _ := strconv.ParseInt(rep.Bandwidth, 10, 64)
				if bandwidth > maxVideoBandwidth {
					maxVideoBandwidth = bandwidth
					videoURL = rep.BaseURL
				}
			}
		case strings.HasPrefix(adaptationSet.MimeType, "audio/"):
			// Find highest quality audio
			for _, rep := range adaptationSet.Representations {
				bandwidth, _ := strconv.ParseInt(rep.Bandwidth, 10, 64)
				if bandwidth > maxAudioBandwidth {
					maxAudioBandwidth = bandwidth
					audioURL = rep.BaseURL
				}
			}
		}
	}

	if videoURL == "" || audioURL == "" {
		return videoFormatsUrlsAndDetails, ErrVideoWayWasDeleted
	}

	videoFormatsUrlsAndDetails.Title = videoID // Note: Title isn't available in DASH manifest
	videoFormatsUrlsAndDetails.FormatsUrls[0] = videoURL
	videoFormatsUrlsAndDetails.FormatsUrls[1] = audioURL

	return videoFormatsUrlsAndDetails, nil
}
