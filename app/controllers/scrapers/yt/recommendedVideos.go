package yt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

// response body json struct
type recommendedVideosResponseBody struct {
	RecommendedVideos []struct {
		VideoID       string `json:"videoID"`
		Title         string `json:"title"`
		LengthSeconds int    `json:"lengthSeconds"`
	} `json:"recommendedVideos"`
}

func recommendedVideos(sender_psid string, arguments []string) error {
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	videoID := arguments[0]

	// use invidios to get the video formats
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("%s/api/v1/videos/%s?fields=recommendedVideos", config.InvidiousEndpoint, videoID),
		nil,
		make(map[string]string),
	)
	if err != nil {
		return fmt.Errorf("VetvideoFormatsUrls: json requeset err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("VetvideoFormatsUrls: fail to read res.body %w", err)
	}
	var bodyJson recommendedVideosResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("VetvideoFormatsUrls: failt to unmarshale the body %w", err)
	}
	RecommendedVideos := bodyJson.RecommendedVideos

	comul := 0
	n := config.MaxReturnedVideo
	for i := 0; i < len(RecommendedVideos); i++ {
		if RecommendedVideos[i].VideoID != "" {
			durationInSeconds := RecommendedVideos[i].LengthSeconds
			if durationInSeconds <= config.MaxDurationInSeconds && durationInSeconds > 0 {
				n--
				comul++
				if n+1 == 0 {
					break
				}

				response := facebook.ResponseTemplateAttachment{
					Type: "template",
					Payload: facebook.TemplateAttachmentPayload{
						TemplateType: "generic",
						Elements: []facebook.TemplateAttachmentElement{
							{
								Title:    fmt.Sprintf("%d# %s", comul, RecommendedVideos[i].Title),
								Subtitle: formats.DisplaySecends(durationInSeconds),
								ImageURL: fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", RecommendedVideos[i].VideoID),
								Buttons: []facebook.TemplateAttachmentButton{
									{
										Type:    "postback",
										Title:   "Watch now",
										Payload: fmt.Sprintf("YT_::_WATCH_::_%s_::_%d", RecommendedVideos[i].VideoID, durationInSeconds),
									},
									{
										Type:    "postback",
										Title:   "Listen now",
										Payload: fmt.Sprintf("YT_::_LISTEN_::_%s_::_%d", RecommendedVideos[i].VideoID, durationInSeconds),
									},
									{
										Type:    "postback",
										Title:   "Description",
										Payload: fmt.Sprintf("YT_::_DESCRIPTION_::_%s", RecommendedVideos[i].VideoID),
									},
								},
							},
						},
					},
				}

				go facebookSender.CallSendAPI(sender_psid, response)
			}
		}
	}
	if comul == 0 {
		response := facebook.ResponseMessage{
			Text: "All these videos are long 😥.\n" +
				"Try different keywords.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
	}
	return nil
}
