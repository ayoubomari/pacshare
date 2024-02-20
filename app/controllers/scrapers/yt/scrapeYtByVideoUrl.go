package yt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
)

func scrapeYtByVideoUrl(sender_psid string, url string) error {
	var videoID string
	// extract the videoID
	if strings.Contains(url, "watch?v=") {
		url2 := strings.ReplaceAll(url, "&", "watch?v=")
		videoID = strings.Split(url2, "watch?v=")[1]
	} else if strings.Contains(url, "youtu.be/") {
		videoID = strings.Split(url, "youtu.be/")[1]
		videoID = strings.Split(videoID, "?")[0]
	} else if strings.Contains(url, "shorts/") {
		videoID = strings.Split(url, "shorts/")[1]
		videoID = strings.Split(videoID, "?")[0]
	}

	if len(videoID) != 11 {
		return errors.New("couldn't extract the videoID")
	}

	videoDetails, err := getVideoDetails(videoID)
	if err != nil && errors.Is(err, ErrVideoDetailsIsNil) {
		response := facebook.ResponseMessage{
			Text: "Anvalide Youtube Url 💔.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return nil
	} else if err != nil {
		return fmt.Errorf("scrapeYtByVideoUrl : %w", err)
	}

	if videoDetails.DurationInSeconds == 0 {
		response := facebook.ResponseMessage{
			Text: "You can't watch live 😞.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return nil
	}
	if videoDetails.DurationInSeconds > config.MaxDurationInSeconds {
		response := facebook.ResponseMessage{
			Text: "This video is too long 🤷‍♂️.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return nil
	}

	response := facebook.ResponseTemplateAttachment{
		Type: "template",
		Payload: facebook.TemplateAttachmentPayload{
			TemplateType: "generic",
			Elements: []facebook.TemplateAttachmentElement{
				{
					Title:    videoDetails.Title,
					Subtitle: formats.DisplaySecends(videoDetails.DurationInSeconds),
					ImageURL: fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", videoID),
					Buttons: []facebook.TemplateAttachmentButton{
						{
							Type:    "postback",
							Title:   "Watch now",
							Payload: fmt.Sprintf("YT_::_WATCH_::_%s_::_%d", videoID, videoDetails.DurationInSeconds),
						},
						{
							Type:    "postback",
							Title:   "Listen now",
							Payload: fmt.Sprintf("YT_::_LISTEN_::_%s_::_%d", videoID, videoDetails.DurationInSeconds),
						},
						{
							Type:    "postback",
							Title:   "Description",
							Payload: fmt.Sprintf("YT_::_DESCRIPTION_::_%s", videoID),
						},
					},
				},
			},
		},
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
