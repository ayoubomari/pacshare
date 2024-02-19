package yt

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	filedownloader "github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

// download, send, and delete the video
// params: sender_psid="sender_id" string
// params: arguments="the arguments of the subcommand" []string[videoID, duration_per_seconds]
func watchYt(sender_psid string, arguments []string) error {
	if len(arguments) < 2 {
		return errors.New("arguments length is lower then 2")
	}

	// get video url .mp4
	videoFormatsAndDetails, err := getVideoFormatsUrls(arguments[0])
	if err != nil {
		return fmt.Errorf("watchYt: could not get the video formats: %w", err)
	}
	videoUrl := videoFormatsAndDetails.FormatsUrls[0]

	// get content size
	fileSize, err := request.GetContentLengthFromResponseHeader(videoUrl)
	if err != nil {
		return fmt.Errorf("watchYt: couldn't get the content size: %w", err)
	}

	// call DownloadFileByRangeWithCallBack to download the file and send It and delete It.
	err = filedownloader.DownloadFileByRangeWithCallBack(sender_psid, videoUrl, "./public/src/videos/", formats.ToFileNameString(videoFormatsAndDetails.Title), "_pac.mp4", fileSize, config.VideoChunksMaxSize, "file")
	if err != nil {
		return fmt.Errorf("watchYt: couldn't download file by chunks: %w", err)
	}

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
	go facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
