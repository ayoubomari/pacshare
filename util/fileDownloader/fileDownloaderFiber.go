package fileDownloader

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/fs"
)

// download file by chunks whithout using range request header
func DownloadAndSendFileWithFiber(sender_psid string, mediaUrl string, outputPath string, fileName string, fileExtentions string, chunkSize int, responseMediaType string) error {
	contentSize := 1
	randomNumber := rand.Intn(1000)
	numChunks := (contentSize + chunkSize - 1) / chunkSize
	fileNamePattern := outputPath + fileName + "_" + fmt.Sprint(randomNumber) + "." + fmt.Sprint(numChunks) + "_" + "%d" + fileExtentions

	filePath := fmt.Sprintf(fileNamePattern, 0)
	err := DownloadEtireFileWithFiber(mediaUrl, filePath)
	if err != nil {
		fs.DeleteFile(filePath)
		return nil
	}

	// get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fs.DeleteFile(filePath)
		return nil
	}
	contentSize = int(fileInfo.Size())
	randomNumber = rand.Intn(1000)
	numChunks = (contentSize + chunkSize - 1) / chunkSize
	fileNamePattern = outputPath + fileName + "_" + fmt.Sprint(randomNumber) + "." + fmt.Sprint(numChunks) + "_" + "%d" + fileExtentions

	fmt.Println("contentSize:", contentSize)

	if contentSize <= chunkSize {
		newFilePath := fmt.Sprintf(fileNamePattern, 1)
		err := os.Rename(filePath, newFilePath)
		if err != nil {
			fmt.Println("error renaming file:", err)
			fs.DeleteFile(filePath)
			return nil
		}

		// send the entire file and delete it
		fileUrl := fmt.Sprintf("https://pacshare.omzor.com/%s", strings.ReplaceAll(newFilePath, "./public/", ""))
		fmt.Println("fileUrl", fileUrl)
		response := facebook.ResponseMediaAttachment{
			Type: responseMediaType,
			Payload: facebook.WebhookBodyAttachmentPayload{
				URL:         fileUrl,
				Is_reusable: false,
			},
		}
		facebookSender.CallSendAPIWithCallback(sender_psid, response, func(err error) error {
			return fs.DeleteFile(newFilePath)
		})
	} else {
		newFilePath := fmt.Sprintf(fileNamePattern, 0)
		err := os.Rename(filePath, newFilePath)
		if err != nil {
			fmt.Println("error renaming file:", err)
			fs.DeleteFile(filePath)
			return nil
		}

		// if the contentSize is bigger than the chunkSize
		sendFileChunks(sender_psid, fileNamePattern, int64(contentSize), int64(chunkSize), responseMediaType)
	}

	return nil
}
