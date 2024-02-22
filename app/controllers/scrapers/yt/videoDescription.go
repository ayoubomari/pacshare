package yt

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
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

	// send description box by chunks
	facebookSender.SendMessageByChunks(sender_psid, videoDetails.Description)

	return nil
}
