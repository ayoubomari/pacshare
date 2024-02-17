package facebookSender

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/ayoubomari/pacshare/util/request"
)

// function to send request to facebook user
func facebookSendRequest(sender_psid string, requestBodyBytes []byte) error {
	res, err := request.JSONReqest(
		"POST",
		fmt.Sprintf("https://graph.facebook.com/v%s.0/%s/messages?access_token=%s", os.Getenv("GRAPHQL_V"), os.Getenv("FACEBOOK_PAGE_ID"), os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")),
		requestBodyBytes,
		make(map[string]string),
	)
	if err != nil {
		return fmt.Errorf("facebookSendRequest: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("facebookSendRequest: %w", err)
	}
	var bodyJson facebook.CallSendAPIResonse
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("facebookSendRequest: %w", err)
	}
	if bodyJson.Error.Message != "" {
		fmt.Println("facebook was returned an error", bodyJson.Error.Message)
		response := facebook.ResponseMessage{
			Text: "Something wrong try another time 🙁.",
		}
		CallSendAPI(sender_psid, response)
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

		err = facebookSendRequest(sender_psid, bodyJsonBytes)
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

		err = facebookSendRequest(sender_psid, bodyJsonBytes)
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

		err = facebookSendRequest(sender_psid, bodyJsonBytes)
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

		err = facebookSendRequest(sender_psid, bodyJsonBytes)
		if err != nil {
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else {
		return errors.New("the response formart doesn't math any of the available formats")
	}
	return nil
}

// Sends response messages via the Send API, with a callbach functions
func CallSendAPIWithCallback(sender_psid string, response interface{}, cb func() error) error {
	err := CallSendAPI(sender_psid, response)
	if err != nil {
		cb()
		return fmt.Errorf("CallSendAPIWithCallback: %w", err)
	}

	err = cb()
	if err != nil {
		return fmt.Errorf("CallSendAPIWithCallback => cb: %w", err)
	}

	return nil
}

func GetMessageInfo(mid string) (facebook.ConversationMessage, error) {
	var conversationMessage facebook.ConversationMessage
	res, err := request.JSONReqest(
		"GET",
		fmt.Sprintf("https://graph.facebook.com/v%s.0/%s?fields=from,message,attachments&access_token=%s", os.Getenv("GRAPHQL_V"), mid, os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")),
		make([]byte, 0),
		make(map[string]string),
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
		make([]byte, 0),
		make(map[string]string),
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
