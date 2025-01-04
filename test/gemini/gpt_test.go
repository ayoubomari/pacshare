package gemini

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/request"
)

type getGeminiAnserResponseBody struct {
    Candidates []struct {
        Content struct {
            Parts []struct {
                Text string `json:"text,omitempty"`
            } `json:"parts,omitempty"`
            Role string `json:"role"`
        } `json:"content,omitempty"`
        FinishReason string `json:"finishReason"`
    } `json:"candidates,omitempty"`
}

func TestGeminiUsinCurl(t *testing.T) {
	startTime := time.Now()
	var wg sync.WaitGroup
	questions := []string{
		"What is the capital of Japan?",
		"How does photosynthesis work?",
		"Who wrote 'Pride and Prejudice'?",
		"What are the main causes of climate change?",
		"How do vaccines work?",
		"What is the theory of relativity?",
		"Who painted the Mona Lisa?",
		"What is the difference between HTML and CSS?",
		"How do black holes form?",
		"What are the main principles of economics?",
		"How does the human immune system function?",
		"What is artificial intelligence?",
		"Who was the first person to walk on the moon?",
		"What are the key features of democracy?",
		"How do earthquakes occur?",
		"What is the process of evolution?",
		"Who invented the telephone?",
		"What are the main components of DNA?",
		"How does a combustion engine work?",
		"What is the significance of the Magna Carta?",
	}

	for _, question := range questions {
		wg.Add(1)
		go func(question string) {
			defer wg.Done()

			jsonBytes := []byte(fmt.Sprintf(`{"contents":[{"parts":[{"text":"%s"}]}]}`, question))
			
			res, err := request.JSONReqest(
				"POST",
				fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", config.GeminiApiKey),
				jsonBytes,
				nil,
				false,
			)
			if err != nil {
				fmt.Printf("JSON request error: %s\n", err)
				return
			}
			if res == nil {
				fmt.Println("Received nil response")
				return
			}
			defer res.Body.Close()

			bodyBytes, err := io.ReadAll(res.Body)
			fmt.Println(string(bodyBytes))
			if err != nil {
				fmt.Printf("Failed to read response body: %s\n", err)
				return
			}

			var bodyJson getGeminiAnserResponseBody
			err = json.Unmarshal(bodyBytes, &bodyJson)
			if err != nil {
				fmt.Printf("Failed to unmarshal the body: %s\n", err)
				return
			}

			if len(bodyJson.Candidates) == 0 {
				fmt.Println("No candidates found")
				return
			}

			if len(bodyJson.Candidates[0].Content.Parts) == 0 {
				fmt.Println("No parts found in the first candidate")
				return
			}

			if bodyJson.Candidates[0].Content.Parts[0].Text == "" {
				fmt.Println("Empty text in the first part of the first candidate")
				return
			}

			// fmt.Println("answer: ", bodyJson.Candidates[0].Content.Parts[0].Text)
		}(question)

		
		time.Sleep(1 * time.Second)
	}

	wg.Wait()
	fmt.Printf("It took %s\n", time.Since(startTime))
}