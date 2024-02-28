package pdf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
)

var SomethingWasWrong = facebook.ResponseMessage{
	Text: "Something wrong try another time 🙁.",
}

func downloadPdf(sender_psid string, arguments []string) error {
	fmt.Println("from downloadPdf")
	if len(arguments) < 1 {
		return errors.New("arguments length is lower then 1")
	}
	pdfLink := arguments[0]
	fmt.Println("pdfLink:", pdfLink)

	pdfInfo, err := GetPdfInfo(pdfLink)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("downloadPdf: %w", err)

	}
	if len(strings.Split(pdfInfo.SessionID, "_")) < 2 {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("downloadPdf: %w", err)

	}

	// get pdf path link
	pdfFilePath := "https://z45st682fx.zlib-cdn.com/dl2.php?id=" + strings.Split(pdfInfo.SessionID, "_")[0] + "&h=" + strings.Split(pdfInfo.SessionID, "_")[1] + "&u=cache&ext=pdf&n=Living%20in%20the%20light%20a%20guide%20to%20personal%20transformation"
	fmt.Println("pdf link Path:", pdfFilePath)

	// download and send pdf
	fileDownloader.DownloadAndSendFileWithFiber(sender_psid, pdfFilePath, "./public/src/pdfs/", formats.ToFileNameString(pdfInfo.Name), "_pac.pdf", config.PdfChunksMaxSize, "file")
	//send pdf complition response message
	response := facebook.ResponseMessage{
		Text: "All pdf files have been sent. ✅",
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
