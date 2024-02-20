package yt

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

// download, send, and delete the youtube audio
// params: sender_psid="sender_id" string
// params: arguments="the arguments of the subcommand" []string[videoID, duration_per_seconds]
func listenYt(sender_psid string, arguments []string) error {
	if len(arguments) < 2 {
		return errors.New("arguments length is lower then 2")
	}

	// get audio url .mp3
	videoFormatsAndDetails, err := getVideoFormatsUrls(arguments[0])
	if errors.Is(err, ErrVideoWayWasDeleted) {
		response := facebook.ResponseMessage{
			Text: "The video may have been removed from YouTube.",
		}
		facebookSender.CallSendAPI(sender_psid, response)
		return nil
	} else if err != nil {
		return fmt.Errorf("listenYt: could not get the audio formats: %w", err)
	}
	audioUrl := videoFormatsAndDetails.FormatsUrls[1]

	// get content size
	fileSize, err := request.GetContentLengthFromResponseHeader(audioUrl)
	if err != nil {
		return fmt.Errorf("listenYt: couldn't get the content size: %w", err)
	}

	// call DownloadFileByRangeWithCallBack to download the file and send It and delete It.
	err = fileDownloader.DownloadFileByRangeWithCallBack(sender_psid, audioUrl, "./public/src/audios/", formats.ToFileNameString(videoFormatsAndDetails.Title), "_pac.mp4", fileSize, config.AudioChunksMaxSize, "file")
	if err != nil {
		return fmt.Errorf("listenYt: couldn't download file by chunks: %w", err)
	}

	// send recommended videos button
	response := facebook.ResponseTemplateAttachment{
		Type: "template",
		Payload: facebook.TemplateAttachmentPayload{
			TemplateType: "button",
			Text:         "Recommended videos",
			Buttons: []facebook.TemplateButtonButton{
				{
					Type:    "postback",
					Title:   "Show",
					Payload: fmt.Sprintf("YT_::_RECOMMENDED_VIDEO_::_%s", arguments[0]),
				},
			},
		},
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
