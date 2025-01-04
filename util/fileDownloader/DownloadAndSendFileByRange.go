package fileDownloader

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/fs"
)

// DownloadAndSendFileByRange downloads a file in chunks using the Range header and sends it
func DownloadAndSendFileByRange(sender_psid string, mediaUrl string, outputPath string, fileName string, fileExtentions string, contentSize int, chunkSize int, responseMediaType string, headers map[string]string, useProxy bool) error {
	randomNumber := rand.Intn(1000)
	numChunks := (contentSize + chunkSize - 1) / chunkSize
	fileNamePattern := outputPath + fileName + "_" + fmt.Sprint(randomNumber) + "." + fmt.Sprint(numChunks) + "_" + "%d" + fileExtentions

	// waitgroup
	var wg sync.WaitGroup

	// Download the file in chunks
	var offset int
	for i := 0; i < 1; /*numChunks*/ i++ {
		wg.Add(1)
		go func(goOffset int, goFileNumber int) {
			// Get the next proxy
			proxyURL, err := url.Parse(config.GetNextProxy())
			if err != nil {
				fmt.Println("Invalid proxy URL:", err)
				return
			}

			// Set up the proxy client
			client := &http.Client{}
			if useProxy {
				client = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyURL),
					},
					Timeout: 10 * time.Second,
				}
			}

			// Create HTTP GET request
			req, err := http.NewRequest("GET", mediaUrl, nil)
			if err != nil {
				fmt.Println("err:", err)
				wg.Done()
				return
			}

			// Add the custom headers to the request
			for key, value := range headers {
				req.Header.Set(key, value)
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

			// Add the range to the request header
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
				fs.DeleteFile(outputFileFullPathName)
				wg.Done()
				return
			}

			fileUrl := fmt.Sprintf("https://pacshare.omzor.com/%s", strings.ReplaceAll(outputFileFullPathName, "./public/", ""))
			// Send the file to the client and remove it
			response := facebook.ResponseMediaAttachment{
				Type: responseMediaType,
				Payload: facebook.WebhookBodyAttachmentPayload{
					URL:         fileUrl,
					Is_reusable: false,
				},
			}
			facebookSender.CallSendAPI(sender_psid, response)
			// fs.DeleteFile(outputFileFullPathName)

			wg.Done()
		}(offset, i+1)

		// Move the offset to the next chunk
		offset += chunkSize
	}

	wg.Wait()

	// Send completion response message
	response := facebook.ResponseMessage{
		Text: "All files have been sent. ✅",
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
