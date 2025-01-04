package yt

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/request"
)

// listenYt downloads, sends, and deletes the YouTube audio
func listenYt(sender_psid string, arguments []string) (err error) {
	// Defer a recovery function to handle panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered in listenYt: %v", r)
			sendErrorMessage(sender_psid, "An unexpected error occurred. Please try again later.")
		}
	}()

	if len(arguments) < 2 {
		return errors.New("arguments length is lower than 2")
	}

	videoFormatsAndDetails, err := getVideoFormatsUrls(arguments[0])
	if err != nil {
		if errors.Is(err, ErrVideoWayWasDeleted) {
			sendErrorMessage(sender_psid, "Something went wrong. Please try again.")
			return nil
		}
		return fmt.Errorf("listenYt: could not get the audio formats: %w", err)
	}

	if len(videoFormatsAndDetails.FormatsUrls) < 2 {
		sendErrorMessage(sender_psid, "Audio format not available. Please try another video.")
		return nil
	}

	audioUrl := videoFormatsAndDetails.FormatsUrls[1]

	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "en,fr-FR;q=0.9,fr;q=0.8,en-US;q=0.7",
		"Cache-Control":             "no-cache",
		"Connection":                "keep-alive",
		"Pragma":                    "no-cache",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		"X-Browser-Channel":         "stable",
		"X-Browser-Copyright":       "Copyright 2024 Google LLC. All rights reserved.",
		"X-Browser-Validation":      "3gQbjS+guBpGZLzijx6RZ1VZHAA=",
		"X-Browser-Year":            "2024",
		"X-Client-Data":             "CIS2yQEIprbJAQipncoBCICXywEIk6HLAQiHoM0BCLKezgEI/aXOAQj7vs4BCIHDzgEIo8bOAQioyM4BCInJzgEI/srOAQiYy84BCMbMzgEY9cnNARicsc4BGIDKzgE=",
		"sec-ch-ua":                 "\"Chromium\";v=\"130\", \"Google Chrome\";v=\"130\", \"Not?A_Brand\";v=\"99\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Linux\"",
	}

	fileSize, err := request.GetContentLengthFromResponseHeader(audioUrl, headers, false)
	if err != nil {
		return fmt.Errorf("listenYt: couldn't get the content size: %w", err)
	}

	err = fileDownloader.DownloadAndSendFileByRange(sender_psid, audioUrl, "./public/src/audios/", formats.ToFileNameString(videoFormatsAndDetails.Title), "_pac.mp4", fileSize, config.AudioChunksMaxSize, "file", headers, false)
	if err != nil {
		return fmt.Errorf("listenYt: couldn't download file by chunks: %w", err)
	}

	sendRecommendedVideosButton(sender_psid, arguments[0])

	return nil
}
