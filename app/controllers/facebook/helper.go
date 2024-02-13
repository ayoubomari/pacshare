package facebook

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
	res, err := request.PostJSONReqest(
		fmt.Sprintf("https://graph.facebook.com/v%s.0/me/messages?access_token=%s", os.Getenv("GRAPHQL_V"), os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")),
		requestBodyBytes,
		make(map[string]string),
	)
	if err != nil {
		return fmt.Errorf("CallSendAPI: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("CallSendAPI: %w", err)
	}
	var bodyJson facebook.CallSendAPIResonse
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return fmt.Errorf("CallSendAPI: %w", err)
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
			fmt.Println("err:", err)
			return fmt.Errorf("CallSendAPI: %w", err)
		}
	} else if responseAttachment, ok := response.(facebook.ResponseAttachment); ok {
		requestBody := facebook.ResponseWithAttachment{
			Recipient:  facebook.ResponseRecipient{ID: sender_psid},
			Attachment: responseAttachment,
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
		return errors.New("the response formart doesn't math any of the available formats.")
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
