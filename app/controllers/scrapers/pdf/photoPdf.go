package pdf

import (
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/fs"
)

// photoPdf handles the process of fetching and sending a PDF's cover image
func photoPdf(sender_psid string, arguments []string) error {
	// Check if we have sufficient arguments
	if len(arguments) < 1 {
		return errors.New("arguments length is lower than 1")
	}

	pdfLink := arguments[0]

	// Fetch PDF information
	pdfInfo, err := GetPdfInfo(pdfLink)
	if err != nil {
		// If there's an error, send a generic error message to the user
		go facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("photoPdf: %w", err)
	}

	// Check if a cover image is available
	if pdfInfo.Cover == "" {
		response := facebook.ResponseMessage{
			Text: "No screenshots were found. 🤷‍♂️",
		}
		go facebookSender.CallSendAPI(sender_psid, response)
		return nil
	}

	// Safely extract file extension
	fileExt := filepath.Ext(pdfInfo.Cover)
	if fileExt == "" {
		fileExt = ".jpg" // Default to .jpg if no extension found
	}

	// Generate a unique filename for the cover image
	coverPath := fmt.Sprintf("./public/src/images/%s_%d%s", 
		formats.ToFileNameString(pdfInfo.Name), 
		rand.Intn(1000), 
		fileExt)

	// Download the cover image
	err = fileDownloader.DownloadEtireFile(pdfInfo.Cover, coverPath)
	if err != nil {
		go facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("photoPdf: %w", err)
	}

	// Prepare the response with the image
	response := facebook.ResponseMediaAttachment{
		Type: "image",
		Payload: facebook.WebhookBodyAttachmentPayload{
			URL:         fmt.Sprintf("https://pacshare.omzor.com/%s", strings.TrimPrefix(coverPath, "./public/")),
			Is_reusable: false,
		},
	}

	// Send the response asynchronously and delete the file after sending
	go func() {
		facebookSender.CallSendAPI(sender_psid, response)
		// Delete the file after sending
		fs.DeleteFile(coverPath)
	}()

	return nil
}