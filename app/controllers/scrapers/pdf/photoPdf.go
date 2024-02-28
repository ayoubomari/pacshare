package pdf

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/fs"
)

func photoPdf(sender_psid string, arguments []string) error {
	fmt.Println("from downloadPdf")
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	pdfLink := arguments[0]
	fmt.Println("pdfLink:", pdfLink)

	pdfInfo, err := GetPdfInfo(pdfLink)
	if err != nil {
		return facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
	}
	if pdfInfo.Cover == "" {
		response := facebook.ResponseMessage{
			Text: "No screenshots were found. 🤷‍♂️",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	coverPath := fmt.Sprintf("./public/src/images/%s_%d.%s", formats.ToFileNameString(pdfInfo.Name), rand.Intn(1000), formats.ToFileNameString(strings.Split(pdfInfo.Cover, ".")[3]))
	err = fileDownloader.DownloadEtireFile(pdfInfo.Cover, coverPath)
	if err != nil {
		return facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
	}

	response := facebook.ResponseMediaAttachment{
		Type: "image",
		Payload: facebook.WebhookBodyAttachmentPayload{
			URL:         fmt.Sprintf("https://pacshare.omzor.com/%s", strings.ReplaceAll(coverPath, "./public/", "")),
			Is_reusable: false,
		},
	}
	facebookSender.CallSendAPI(sender_psid, response)

	fs.DeleteFile(coverPath)

	return nil
}
