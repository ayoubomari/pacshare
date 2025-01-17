package facebookSender

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/util/request"
)

// function to send request to facebook user
func facebookSendRequest(sender_psid string, requestBodyBytes []byte, Errornotify bool) error {
	res, err := request.JSONReqest(
		"POST",
		fmt.Sprintf("https://graph.facebook.com/v%s.0/%s/messages?access_token=%s", os.Getenv("GRAPHQL_V"), os.Getenv("FACEBOOK_PAGE_ID"), os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")),
		requestBodyBytes,
		nil,
		false,
	)
	if err != nil {
		return fmt.Errorf("facebookSendRequest: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	// fmt.Println("bodyBytes: ", string(bodyBytes))
	if err != nil {
		return fmt.Errorf("facebookSendRequest: %w", err)
	}
	// fmt.Println(string(bodyBytes))
	var bodyJson facebook.CallSendAPIResonse
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("facebookSendRequest: %w", err)
	}
	if Errornotify && bodyJson.Error.Message != "" {
		fmt.Println("facebook was returned an error", bodyJson.Error.Message)
		if bodyJson.Error.Message == "(#100) Upload failed" { // if it's an attachment upload error
			response := facebook.ResponseMessage{
				Text: "Failed to upload an attachment 🙁.",
			}
			CallSendAPI(sender_psid, response)
		} else {
			response := facebook.ResponseMessage{
				Text: "Something wrong try another time 🙁.",
			}
			CallSendAPI(sender_psid, response)
		}
	}

	return nil
}

// Sends response messages via the Send API
func CallSendAPI(sender_psid string, response interface{}) error {
	if responseMessage, ok := response.(facebook.ResponseMessage); ok {
		requestBody := facebook.ResponseWithMessage{
			Recipient: facebook.ResponseRecipient{ID: sender_psid},
			Message:   responseMessage,
		}
		bodyJsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}

		err = facebookSendRequest(sender_psid, bodyJsonBytes, false)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else if responseAttachment, ok := response.(facebook.ResponseMediaAttachment); ok {
		requestBody := facebook.ResponseWithMediaAttachment{
			Recipient: facebook.ResponseRecipient{ID: sender_psid},
			Message: facebook.MediaMessage{
				Attachment: responseAttachment,
			},
		}
		bodyJsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}

		err = facebookSendRequest(sender_psid, bodyJsonBytes, true)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else if responseAttachment, ok := response.(facebook.ResponseTemplateAttachment); ok {
		requestBody := facebook.ResponseWithTemplateAttachment{
			Recipient: facebook.ResponseRecipient{ID: sender_psid},
			Message: facebook.TemplateMessage{
				Attachment: responseAttachment,
			},
		}
		bodyJsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}

		err = facebookSendRequest(sender_psid, bodyJsonBytes, true)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else if responseQuickReplay, ok := response.(facebook.QuickReplyMessage); ok {
		requestBody := facebook.ResponseWithQuickReplay{
			Recipient: facebook.ResponseRecipient{ID: sender_psid},
			Message:   responseQuickReplay,
		}
		bodyJsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}

		err = facebookSendRequest(sender_psid, bodyJsonBytes, true)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else if responseResponseAction, ok := response.(string); ok {
		requestBody := facebook.ResponseWithResponseAction{
			Recipient:     facebook.ResponseRecipient{ID: sender_psid},
			Sender_action: responseResponseAction,
		}
		bodyJsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}

		err = facebookSendRequest(sender_psid, bodyJsonBytes, false)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else {
		return errors.New("the response formart doesn't math any of the available formats")
	}
	return nil
}

// send typeing action response
func SendTypingOn(sender_psid string) error {
	return CallSendAPI(sender_psid, "TYPING_ON")
}

// Sends response messages via the Send API, with a callbach functions
func CallSendAPIWithCallback(sender_psid string, response interface{}, cb func(err error) error) error {
	err := CallSendAPI(sender_psid, response)
	if err != nil {
		cb(err)
		return fmt.Errorf("CallSendAPIWithCallback: %w", err)
	}

	err = cb(nil)
	if err != nil {
		return fmt.Errorf("CallSendAPIWithCallback => cb: %w", err)
	}

	return nil
}

// SendMessageByChunks sends a message in chunks, respecting Facebook's message length limit
func SendMessageByChunks(sender_psid string, message string) error {
    totalMessages := (len(message) + config.MaxMessageLength - 1) / config.MaxMessageLength

    for i := 0; i < totalMessages; i++ {
        start := i * config.MaxMessageLength
        end := (i + 1) * config.MaxMessageLength
        if end > len(message) {
            end = len(message)
        }

        response := facebook.ResponseMessage{
            Text: message[start:end],
        }

        err := CallSendAPI(sender_psid, response)
        if err != nil {
            return fmt.Errorf("failed to send chunk %d: %w", i+1, err)
        }

        // Wait for 1.5 seconds before sending the next chunk
        if i < totalMessages-1 {
            time.Sleep(1500 * time.Millisecond)
        }
    }

    return nil
}

func GetMessageInfo(mid string) (facebook.ConversationMessage, error) {
	var conversationMessage facebook.ConversationMessage
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://graph.facebook.com/v%s.0/%s?fields=from,message,attachments&access_token=%s", os.Getenv("GRAPHQL_V"), mid, os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")),
		nil,
		nil,
		false,
	)
	if err != nil {
		return conversationMessage, fmt.Errorf("GetMessageInfo: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return conversationMessage, fmt.Errorf("GetMessageInfo: %w", err)
	}
	err = json.Unmarshal(bodyBytes, &conversationMessage)
	if err != nil {
		return conversationMessage, fmt.Errorf("GetMessageInfo: %w", err)
	}

	return conversationMessage, nil
}

func GetSenderInfo(sender_psid string) (facebook.SenderInfo, error) {
	var senderInfo facebook.SenderInfo
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://graph.facebook.com/%s?fields=first_name,last_name,profile_pic&access_token=%s", sender_psid, os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")),
		nil,
		nil,
		false,
	)
	if err != nil {
		return senderInfo, fmt.Errorf("GetSenderInfo: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return senderInfo, fmt.Errorf("GetSenderInfo: %w", err)
	}
	err = json.Unmarshal(bodyBytes, &senderInfo)
	if err != nil {
		return senderInfo, fmt.Errorf("GetSenderInfo: %w", err)
	}

	return senderInfo, nil
}
