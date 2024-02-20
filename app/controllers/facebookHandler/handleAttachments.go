package facebookHandler

import (
	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

// handle Attachments message
func handleAttachments(sender_psid string, message facebook.Message) error {
	response := facebook.ResponseMessage{
		Text: "🤔",
	}
	facebookSender.CallSendAPI(sender_psid, response)

	return nil
}
