package filedownloader

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/fs"
)

// download file by chunks usign request range header
func DownloadFileByRangeWithCallBack(sender_psid string, mediaUrl string, outputPath string, fileName string, fileExtentions string, contentSize int, chunkSize int, responseMediaType string) error {
	randomNumber := rand.Intn(1000)
	totalFileNumber := contentSize / chunkSize
	if contentSize%chunkSize > 0 {
		totalFileNumber += 1
	}
	fileNamePattern := outputPath + fileName + "_" + fmt.Sprint(randomNumber) + "." + fmt.Sprint(totalFileNumber) + "_" + "%d" + fileExtentions

	// Initialize HTTP client
	client := &http.Client{}

	// waitgroup
	var wg sync.WaitGroup

	// Download the file in chunks
	var offset int
	for i := 0; i < totalFileNumber; i++ {
		wg.Add(1)
		go func(goOffset int, goFileNumber int) {
			// Create HTTP GET request
			req, err := http.NewRequest("GET", mediaUrl, nil)
			if err != nil {
				fmt.Println("err:", err)
				wg.Done()
				return
			}

			outputFileFullPathName := fmt.Sprintf(fileNamePattern, goFileNumber)

			// Create an output file (you can use any io.Writer)
			outputFile, err := os.Create(outputFileFullPathName)
			if err != nil {
				fmt.Println("err:", err)
				wg.Done()
				return
			}
			defer outputFile.Close()

			endRange := goOffset + chunkSize - 1
			if endRange >= contentSize {
				endRange = contentSize - 1
			}

			// add the range to the request header
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", goOffset, endRange))
			// Perform the HTTP request
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("err:", err)
				fs.DeleteFile(outputFileFullPathName)
				wg.Done()
				return
			}
			defer resp.Body.Close()

			// Write the response body to the output file
			_, err = io.Copy(outputFile, resp.Body)
			if err != nil {
				wg.Done()
				fs.DeleteFile(outputFileFullPathName)
				wg.Done()
				return
			}

			fileUrl := fmt.Sprintf("https://pacshare.omzor.com/%s", strings.ReplaceAll(outputFileFullPathName, "./public/", ""))
			// send the file to the client and remove it
			response := facebook.ResponseMediaAttachment{
				Type: responseMediaType,
				Payload: facebook.WebhookBodyAttachmentPayload{
					URL:         fileUrl,
					Is_reusable: false,
				},
			}
			facebookSender.CallSendAPI(sender_psid, response)
			fs.DeleteFile(outputFileFullPathName)

			wg.Done()
		}(offset, i+1)

		// Move the offset to the next chunk
		offset += chunkSize
	}

	wg.Wait()

	//send complition response message
	response := facebook.ResponseMessage{
		Text: "All files have been sent. ✅",
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
