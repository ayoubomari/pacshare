package pdf

import (
	"errors"
	"fmt"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

func descriptionPdf(sender_psid string, arguments []string) error {
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

	responseText := ""
	if pdfInfo.Name != "" {
		responseText += fmt.Sprintf("title: %s\n\n", pdfInfo.Name)
	}
	if pdfInfo.Author != "" {
		responseText += fmt.Sprintf("✒️: %s\n\n", pdfInfo.Author)
	}
	if pdfInfo.FileSize != "" {
		responseText += fmt.Sprintf("⬇️: %s\n\n", pdfInfo.FileSize)
	}
	if pdfInfo.PagesNum != "" {
		responseText += fmt.Sprintf("📜: %s\n\n", pdfInfo.PagesNum)
	}
	if pdfInfo.Language != "" {
		responseText += fmt.Sprintf("🌐: %s\n\n", pdfInfo.Language)
	}

	// final result
	if responseText == "" {
		response := facebook.ResponseMessage{
			Text: "No description was found. 🤷‍♂️",
		}
		return facebookSender.CallSendAPI(sender_psid, response)
	}

	response := facebook.ResponseMessage{
		Text: responseText,
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
