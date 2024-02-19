package facebook

import (
	"fmt"
	"testing"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	facebookModel "github.com/ayoubomari/pacshare/app/models/facebook"
)

// test seding text message
func TestSendFacebookMessage(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebookModel.ResponseMessage{
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
	response := facebookModel.ResponseMediaAttachment{
		Type: "video",
		Payload: facebookModel.WebhookBodyAttachmentPayload{
			URL:         "pacshare.omzor.com/src/videos/woody.mp4",
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
	response := facebookModel.ResponseTemplateAttachment{
		Type: "template",
		Payload: facebookModel.TemplateAttachmentPayload{
			TemplateType: "generic",
			Elements: []facebookModel.TemplateAttachmentElement{
				{
					Title:    "title",
					Subtitle: "subtitle",
					ImageURL: "https://scontent.frak2-1.fna.fbcdn.net/v/t1.15752-9/370319136_3585075875091292_546338281524030784_n.png?stp=dst-png_s2048x2048&_nc_cat=111&ccb=1-7&_nc_sid=8cd0a2&_nc_ohc=ZzFhvIpnXDEAX8NSLWQ&_nc_ht=scontent.frak2-1.fna&oh=03_AdR_XV_5NDkC5FEH3kpziPRRaDWKayFbVULxTeYO3PCsFw&oe=65F43D2E",
					Buttons: []facebookModel.TemplateAttachmentButton{
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
	response := facebookModel.ResponseTemplateAttachment{
		Type: "template",
		Payload: facebookModel.TemplateAttachmentPayload{
			TemplateType: "button",
			Text:         "button text",
			Buttons: []facebookModel.TemplateButtonButton{
				{
					Type:    "postback",
					Title:   "Show",
					Payload: "template button payload",
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
	response := "TYPING_ON"
	err := facebookSender.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

// test seding quick replay message
func TestSendFacebookQuickReplay(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebookModel.QuickReplyMessage{
		Text: "title",
		QuickReplies: []facebookModel.QuickReplyResponse{
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
