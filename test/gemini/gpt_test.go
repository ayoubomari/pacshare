package gemini

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/request"
)

type getGeminiAnserResponseBody struct {
	Candidates []struct {
		Content *struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
}

func TestGeminiUsinCurl(t *testing.T) {
	// Convert the request body to JSON
	question := "what is the minimum wage in Morocco"
	jsonBytes := []byte(fmt.Sprintf(`{"contents":[{"parts":[{"text":"%s"}]}]}`, question))

	// send http request
	res, err := request.JSONReqest(
		"POST",
		fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", config.GeminiApiKey),
		jsonBytes,
		nil,
	)
	if err != nil {
		fmt.Printf("getVideoDetails: json requeset err: %s", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("getVideoDetails: fail to read res.body %s", err)
	}
	var bodyJson getGeminiAnserResponseBody
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		fmt.Printf("getVideoDetails: failt to unmarshale the body %s", err)
	}

	if len(bodyJson.Candidates) < 1 || len(bodyJson.Candidates[0].Content.Parts) < 1 {
		fmt.Printf("getVideoDetails: no ansert found %s", err)
	}

	fmt.Println(bodyJson.Candidates[0].Content.Parts[0].Text)
}
