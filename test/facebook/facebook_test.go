package facebook

import (
	"fmt"
	"testing"

	"github.com/ayoubomari/pacshare/app/controllers/facebook"
	facebookModel "github.com/ayoubomari/pacshare/app/models/facebook"
)

func TestSendFacebookMessage(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebookModel.ResponseMessage{
		Text: "hi from the test",
	}
	err := facebook.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendFacebookMediaAttachment(t *testing.T) {
	sender_psid := "4345084215546247"
	response := facebookModel.ResponseMediaAttachment{
		Type: "video",
		Payload: facebookModel.WebhookBodyAttachmentPayload{
			URL:         "https://rr1---sn-5hnednss.googlevideo.com/videoplayback?expire=1707937379&ei=A7rMZa-eC6rY6dsPlYGHyAU&ip=2a0a%3A4cc0%3A1%3A11a2%3Af5ea%3A9331%3A64ee%3A6d9b&id=o-AG41sUa_0O0pFfdSGTjs8inidG2xCPNCMenGbham4dkF&itag=18&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&vprv=1&svpuc=1&mime=video%2Fmp4&cnr=14&ratebypass=yes&dur=25.147&lmt=1688629803719031&fexp=24007246&c=ANDROID&txp=1438434&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cvprv%2Csvpuc%2Cmime%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AJfQdSswRAIgRCSqhYlHTUllcAogkalogkVFSAJAIPgIIm79hqxYkc4CICjGzoz40y-4RIU7FnmgS6isIf1-2qSr29hck2VxJn0Z&host=rr2---sn-5oxmp55u-8pxe.googlevideo.com&rm=sn-5oxmp55u-8pxe7l,sn-4g5erz7z&req_id=a13b44c114e636e2&ipbypass=yes&redirect_counter=3&cm2rm=sn-apns7s&cms_redirect=yes&cmsv=e&mh=u1&mip=105.158.228.125&mm=34&mn=sn-5hnednss&ms=ltu&mt=1707916147&mv=m&mvi=1&pl=21&lsparams=ipbypass,mh,mip,mm,mn,ms,mv,mvi,pl&lsig=APTiJQcwRgIhAPclYdozr88uJYuwUMPpFeMnI2htRm5YqZBIL3kvXXcWAiEA_CWZCFNKjFJwtVBUmN4aN_Xeg17AyUbiL0YWjwvvXvk%3D",
			Is_reusable: false,
		},
	}
	err := facebook.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

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
	err := facebook.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

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
	err := facebook.CallSendAPI(sender_psid, response)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFacebookMessageInfo(t *testing.T) {
	mid := "m_CnahzSvUI_YxZZhSiXZ0jpFs0JtiuDOO_WqFxoAp5dZbGnBcbn38q3qYSIRJxtA3JC9TnhK6Bdyig7BMuJDTww"
	messageInfo, err := facebook.GetMessageInfo(mid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", messageInfo)
}
