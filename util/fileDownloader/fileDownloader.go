package fileDownloader

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/fs"
)

// get the zero file, and split and send each chunks of it to the sernder, and delete file and the zero file.
func sendFileChunks(sender_psid string, fileNamePattern string, contentSize int64, chunkSize int64, responseMediaType string) error {
	// Open the source file
	sourceFile, err := os.Open(fmt.Sprintf(fileNamePattern, 0))
	if err != nil {
		return fmt.Errorf("sendFileChunks: %w", err)
	}
	defer sourceFile.Close()

	// Calculate the number of chunks
	numChunks := (contentSize + chunkSize - 1) / chunkSize
	// fmt.Println("numChunks:", numChunks)

	// Iterate through the chunks and copy them to separate files
	for i := int64(0); i < numChunks; i++ {
		// Calculate the offset for the current chunk
		offset := i * chunkSize

		// Seek to the beginning of the current chunk
		_, err := sourceFile.Seek(offset, io.SeekStart)
		if err != nil {
			return fmt.Errorf("sendFileChunks: %w", err)
		}

		// Determine the size of the current chunk
		size := chunkSize
		if offset+chunkSize > contentSize {
			size = contentSize - offset
		}

		// Create the output file for the current chunk
		filePath := fmt.Sprintf(fileNamePattern, i+1)
		outputFile, err := os.Create(filePath)
		if err != nil {
			fs.DeleteFile(fmt.Sprintf(fileNamePattern, i+1))
			return fmt.Errorf("sendFileChunks: %w", err)
		}
		defer outputFile.Close()

		// if file is (apk, obb) add one byte to the first of the file to make it elegible on facebook
		if (i == 0) && (strings.Contains(filePath, "_pac.apk") || strings.Contains(filePath, "_pac.obb")) {
			_, err := outputFile.Write([]byte("0"))
			if err != nil {
				fs.DeleteFile(filePath)
				return fmt.Errorf("DownloadEtireFile: %w", err)
			}
		}

		// Write the chunk to the output file
		_, err = io.CopyN(outputFile, sourceFile, size)
		if err != nil {
			fs.DeleteFile(fmt.Sprintf(fileNamePattern, i+1))
			return fmt.Errorf("sendFileChunks: %w", err)
		}

		// send the chunk file and delete it
		fileUrl := fmt.Sprintf("https://pacshare.omzor.com/%s", strings.ReplaceAll(filePath, "./public/", ""))
		response := facebook.ResponseMediaAttachment{
			Type: responseMediaType,
			Payload: facebook.WebhookBodyAttachmentPayload{
				URL:         fileUrl,
				Is_reusable: false,
			},
		}
		facebookSender.CallSendAPIWithCallback(sender_psid, response, func(err error) error {
			return fs.DeleteFile(fmt.Sprintf(fileNamePattern, i+1))
		})
	}

	// delete the zero file
	fs.DeleteFile(fmt.Sprintf(fileNamePattern, 0))

	return nil
}

// download file by chunks whithout using range request header
func DownloadAndSendFile(sender_psid string, mediaUrl string, outputPath string, fileName string, fileExtentions string, contentSize int, chunkSize int, responseMediaType string) error {
	randomNumber := rand.Intn(1000)
	numChunks := (contentSize + chunkSize - 1) / chunkSize
	fileNamePattern := outputPath + fileName + "_" + fmt.Sprint(randomNumber) + "." + fmt.Sprint(numChunks) + "_" + "%d" + fileExtentions

	if contentSize <= chunkSize {
		filePath := fmt.Sprintf(fileNamePattern, 1)
		err := DownloadEtireFile(mediaUrl, filePath)
		if err != nil {
			fs.DeleteFile(filePath)
		}

		// send the entire file and delete it
		fileUrl := fmt.Sprintf("https://pacshare.omzor.com/%s", strings.ReplaceAll(filePath, "./public/", ""))
		response := facebook.ResponseMediaAttachment{
			Type: responseMediaType,
			Payload: facebook.WebhookBodyAttachmentPayload{
				URL:         fileUrl,
				Is_reusable: false,
			},
		}
		facebookSender.CallSendAPIWithCallback(sender_psid, response, func(err error) error {
			return fs.DeleteFile(filePath)
		})
	} else {
		// download the zero file
		filePath := fmt.Sprintf(fileNamePattern, 0)
		err := DownloadEtireFile(mediaUrl, filePath)
		if err != nil {
			fs.DeleteFile(filePath)
		}

		// if the contentSize is bigger than the chunkSize
		sendFileChunks(sender_psid, fileNamePattern, int64(contentSize), int64(chunkSize), responseMediaType)
	}

	return nil
}
