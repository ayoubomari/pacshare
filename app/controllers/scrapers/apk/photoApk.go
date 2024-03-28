package apk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/request"
)

var SomethingWasWrong = facebook.ResponseMessage{
	Text: "Something wrong try another time 🙁.",
}

func photoApk(sender_psid string, arguments []string) error {
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	apkId := arguments[0]

	appInfo, err := GetApkInfoWS2(apkId)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("photoApk: %w", err)
	}

	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://%s.en.aptoide.com/app", appInfo.Nodes.Meta.Data.Uname),
		nil,
		nil,
	)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("photoApk: %w", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("photoApk: %w", err)
	}

	// fmt.Println("getting the html page...")
	imgNum := 0
	// send screenshots
	doc.Find(".app-view__SlideBundlerContainer-sc-oiuh9w-2 img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")

		if exists {
			imgNum++

			// fmt.Printf("Image %d: %s\n", i+1, src)
			response := facebook.ResponseMediaAttachment{
				Type: "image",
				Payload: facebook.WebhookBodyAttachmentPayload{
					URL:         strings.Split(src, "?")[0] + "?w=720",
					Is_reusable: false,
				},
			}
			go facebookSender.CallSendAPI(sender_psid, response)
		}
	})

	// if there is no screenshot found
	if imgNum == 0 {
		response := facebook.ResponseMessage{
			Text: "No screenshots were found. 🤷‍♂️",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	return nil
}
