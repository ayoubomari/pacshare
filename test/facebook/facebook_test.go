package facebook

import (
	"fmt"
	"testing"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

// test seding text message
func TestSendFacebookMessage(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebook.ResponseMessage{
		Text: "hi from the test",
	}
	err := facebookSender.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

// test seding media message
func TestSendFacebookMediaAttachment(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebook.ResponseMediaAttachment{
		Type: "file",
		Payload: facebook.WebhookBodyAttachmentPayload{
			URL:         "https://www.emse.fr/~picard/cours/1A/java/livretJava.pdf",
			Is_reusable: false,
		},
	}
	err := facebookSender.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

// test seding geniric template message
func TestSendFacebookGeniricTamplateAttachment(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebook.ResponseTemplateAttachment{
		Type: "template",
		Payload: facebook.TemplateAttachmentPayload{
			TemplateType: "generic",
			Elements: []facebook.TemplateAttachmentElement{
				{
					Title:    "title",
					Subtitle: "subtitle",
					ImageURL: "https://pacshare.omzor.com/img/backgrounds/minimalisme.jpeg",
					Buttons: []facebook.TemplateAttachmentButton{
						{
							Type:    "postback",
							Title:   "button title",
							Payload: "button payload",
						},
					},
				},
			},
		},
	}
	err := facebookSender.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

// test seding button template message
func TestSendFacebookButtonTamplateAttachment(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebook.ResponseTemplateAttachment{
		Type: "template",
		Payload: facebook.TemplateAttachmentPayload{
			TemplateType: "button",
			Text:         "button text",
			Buttons: []facebook.TemplateButtonButton{
				{
					Type:    "postback",
					Title:   "Show",
					Payload: "template button payload",
					// Url:   "https://www.messenger.com/",
				},
			},
		},
	}
	err := facebookSender.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

// test seding response action like TYPING_ON
func TestSendFacebookResponseAction(t *testing.T) {
	sender_psid := "4345084215546247"
	err := facebookSender.SendTypingOn(sender_psid)
	if err != nil {
		t.Fatal(err)
	}
}

// test seding quick replay message
func TestSendFacebookQuickReplay(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebook.QuickReplyMessage{
		Text: "title",
		QuickReplies: []facebook.QuickReplyResponse{
			{
				ContentType: "text",
				Title:       "pac1",
				Payload:     "PAYLOAD",
				ImageURL:    "https://static.wikia.nocookie.net/pacman/images/2/24/Pac-Man-0.png/revision/latest?cb=20190526005949",
			},
			{
				ContentType: "text",
				Title:       "pac2",
				Payload:     "PAYLOAD",
				ImageURL:    "https://static.wikia.nocookie.net/pacman/images/2/24/Pac-Man-0.png/revision/latest?cb=20190526005949",
			},
		},
	}
	err := facebookSender.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

// test to get facebook message info
func TestGetFacebookMessageInfo(t *testing.T) {
	mid := "m_CnahzSvUI_YxZZhSiXZ0jpFs0JtiuDOO_WqFxoAp5dZbGnBcbn38q3qYSIRJxtA3JC9TnhK6Bdyig7BMuJDTww"
	messageInfo, err := facebookSender.GetMessageInfo(mid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", messageInfo)
}

// test to get facebook sender info
func TestGetFacebookSenderInfo(t *testing.T) {
	sender_psid := "4345084215546247"
	senderInfo, err := facebookSender.GetSenderInfo(sender_psid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", senderInfo)
}
