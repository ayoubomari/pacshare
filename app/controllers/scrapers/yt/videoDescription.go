package yt

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/formats"
)

// send the video description in a messge request
// params: sender_psid="sender_id" string
// params: arguments="the arguments of the subcommand" []string[videoID]
func VideoDescription(sender_psid string, arguments []string) error {
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	videoID := arguments[0]

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

	// send video details
	reponse := facebook.ResponseMessage{
		Text: "🎥Channel: " + videoDetails.Author + "\n" +
			"🕔Duration: " + formats.DisplaySecends(videoDetails.DurationInSeconds) + "\n" +
			"👁️Views: " + videoDetails.ViewCount + "\n" +
			"📅Date: " + videoDetails.UploadDate + "\n",
	}
	facebookSender.CallSendAPI(sender_psid, reponse)

	// send description box
	totalMessages := len(videoDetails.Description) / config.MaxMessageLength
	if len(videoDetails.Description)%config.MaxMessageLength > 0 {
		totalMessages += 1
	}
	fmt.Println("len(videoDetails.Description):", len(videoDetails.Description))
	for i := 0; i < totalMessages; i++ {
		start := i * config.MaxMessageLength
		var end int
		if i == totalMessages-1 {
			end = len(videoDetails.Description)
		} else {
			end = i * config.MaxMessageLength
		}
		fmt.Println("start:", start)
		fmt.Println("end:", end)
		DescriptionResponse := facebook.ResponseMessage{
			Text: videoDetails.Description[start:end],
		}
		go facebookSender.CallSendAPI(sender_psid, DescriptionResponse)
	}

	return nil
}
