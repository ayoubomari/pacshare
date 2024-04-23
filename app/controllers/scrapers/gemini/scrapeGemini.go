package gemini

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/fileDownloader"
	"github.com/ayoubomari/pacshare/util/formats"
	"github.com/ayoubomari/pacshare/util/fs"
	"github.com/ayoubomari/pacshare/util/request"
)

var SomethingWasWrong = facebook.ResponseMessage{
	Text: "Oops! Something went wrong.\n" +
		"Try again in a bit?  ⏱️",
}

type getGeminiAnserResponseBody struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
}

// for scrape the gemini text content as markdown
// @params: sender_psid="facebook sender ID"
// @params: arguments=[]string{"user question"}
func scrapeGemini(sender_psid string, arguments []string, mid string) error {
	question := arguments[0]

	var anser string

	//chack if the image is exist on the replay
	if mid != "" {
		noImageFoundResponse := facebook.ResponseMessage{
			Text: "You can ask a question about the image by replying directly to it! ➡️",
		}

		ConversationMessage, err := facebookSender.GetMessageInfo(mid)
		if err != nil || len(ConversationMessage.Attachments.Data) < 1 || !strings.Contains(ConversationMessage.Attachments.Data[0].MimeType, "image/") {
			facebookSender.CallSendAPI(sender_psid, noImageFoundResponse)
			return fmt.Errorf("scrapeGemini: %w", err)
		}

		imageContentBase64, err := fileDownloader.FetchFileContentAsString(ConversationMessage.Attachments.Data[0].ImageData.URL)
		if err != nil {
			facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
			return fmt.Errorf("scrapeGemini: %w", err)
		}

		anser, err = getAnserGeminiProVision(question, ConversationMessage.Attachments.Data[0].MimeType, imageContentBase64)
		if err != nil {
			facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
			return fmt.Errorf("scrapeGemini: %w", err)
		}
	} else {
		var err error
		anser, err = getAnserGeminiPro(question)
		if err != nil {
			facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
			return fmt.Errorf("scrapeGemini: %w", err)
		}
	}

	//send the markdown as chunks
	//write the markdown content in a file
	var questionSlug string
	if len(question) > 50 {
		questionSlug = question[0:50]
	} else {
		questionSlug = question
	}
	questionSlug = strings.ReplaceAll(questionSlug, "_", "-")
	randomNumber := rand.Intn(1000)
	fileName := formats.ToFileNameString(questionSlug) + "_" + fmt.Sprint(randomNumber) + ".1_1_pac.md"
	err := fs.WriteFile(fmt.Sprintf("./public/src/geminis/%s", fileName), &anser)
	if err != nil {
		facebookSender.CallSendAPI(sender_psid, SomethingWasWrong)
		return fmt.Errorf("scrapeGemini: %w", err)
	}

	response := facebook.ResponseMediaAttachment{
		Type: "file",
		Payload: facebook.WebhookBodyAttachmentPayload{
			URL: fmt.Sprintf("https://pacshare.omzor.com/src/geminis/%s", fileName),
		},
	}
	go facebookSender.CallSendAPIWithCallback(sender_psid, response, func(err error) error {
		return fs.DeleteFile(fmt.Sprintf("./public/src/geminis/%s", fileName))
	})

	return nil
}

// todo: write a function that return the anser using gemini pro api
func getAnserGeminiPro(question string) (string, error) {
	jsonBytes := []byte(fmt.Sprintf(`{"contents":[{"parts":[{"text":"%s"}]}]}`, question))

	// send http request
	res, err := request.JSONReqest(
		"POST",
		"http://129.151.173.72/geminipro",
		jsonBytes,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("getAnserGeminiPro: json requeset err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("getAnserGeminiPro: fail to read res.body %w", err)
	}
	var bodyJson getGeminiAnserResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return "", fmt.Errorf("getAnserGeminiPro: failt to unmarshale the body %w", err)
	}

	if len(bodyJson.Candidates) < 1 || len(bodyJson.Candidates[0].Content.Parts) < 1 {
		return "", fmt.Errorf("getAnserGeminiPro: no ansert found %w", err)
	}

	// anser, err := formats.Utf8ToBase64(bodyJson.Candidates[0].Content.Parts[0].Text)
	// if err != nil {
	// 	return "", fmt.Errorf("error converting utf8 to base64: %w", err)
	// }
	return bodyJson.Candidates[0].Content.Parts[0].Text, nil
}

// todo: write a function that return the anser using gemini pro vesion api (use it only if there is an image)
func getAnserGeminiProVision(question string, mimeType string, imageContentBase64 string) (string, error) {
	jsonBytesVision := []byte(fmt.Sprintf(`{"contents":[{"parts":[{"text":"%s"},{"inlineData":{"mimeType": "%s","data":"%s"}}]}]}`, question, mimeType, imageContentBase64))

	// send http request
	res, err := request.JSONReqest(
		"POST",
		"http://129.151.173.72/geminivesionpro",
		jsonBytesVision,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("getAnserGeminiProVision: json requeset err: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("getAnserGeminiProVision: fail to read res.body %w", err)
	}
	var bodyJson getGeminiAnserResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return "", fmt.Errorf("getAnserGeminiProVision: failt to unmarshale the body %w", err)
	}

	if len(bodyJson.Candidates) < 1 || len(bodyJson.Candidates[0].Content.Parts) < 1 {
		return "", fmt.Errorf("getAnserGeminiProVision: no ansert found %w", err)
	}

	// anser, err := formats.Utf8ToBase64(bodyJson.Candidates[0].Content.Parts[0].Text)
	// if err != nil {
	// 	return "", fmt.Errorf("error converting utf8 to base64: %w", err)
	// }
	return bodyJson.Candidates[0].Content.Parts[0].Text, nil
}
