package facebookHandler

import (
	"fmt"
	"regexp"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/apk"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/pdf"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/wiki"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/yt"
	"github.com/ayoubomari/pacshare/app/models/facebook"
)

// handle attachment message come from facebook
func handlePostback(sender_psid string, postback facebook.PostBack) error {
	// global postbacks
	switch postback.Payload {
	case "GET_STARTED_PAYLOAD":

		responseMessage := facebook.ResponseMessage{
			Text: "Hi! Welcome on PacShare 💖\n" +
				"It's a messenger chat bot 🤖\n" +
				"For watching Youtube on Facebook messenger 🎬\n" +
				"We hope you like It 😁",
		}
		err := facebookSender.CallSendAPI(sender_psid, responseMessage)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}

		responseTemplateAttachment := facebook.ResponseTemplateAttachment{
			Type: "template",
			Payload: facebook.TemplateAttachmentPayload{
				TemplateType: "button",
				Text:         "🚨For using this service you have to Like our page.🚨\n\n🚨لستخدام هذه الخدمة ، عليك أن تقوم بلإعجاب بالصفحة🚨\n\n☟\nfb.com/PacShare1",
				Buttons: []facebook.TemplateButtonButton{
					{
						Type:  "web_url",
						Url:   "https://fb.com/PacShare1",
						Title: "Like page",
					},
				},
			},
		}
		err = facebookSender.CallSendAPI(sender_psid, responseTemplateAttachment)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}

		responseMediaAttachment := facebook.ResponseMediaAttachment{
			Type: "image",
			Payload: facebook.WebhookBodyAttachmentPayload{
				URL:         "https://pacshare.omzor.com/static_src/imgs/run.gif",
				Is_reusable: true,
			},
		}
		err = facebookSender.CallSendAPI(sender_psid, responseMediaAttachment)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}

		return nil

	// for menu
	case "DOWNLOAD_PACSHARE_APP":
		response := facebook.ResponseTemplateAttachment{
			Type: "template",
			Payload: facebook.TemplateAttachmentPayload{
				TemplateType: "button",
				Text:         "Select your OS 📲",
				Buttons: []facebook.TemplateButtonButton{
					{
						Type:  "web_url",
						Title: "Android 🤖",
						Url:   "https://www.facebook.com/groups/1759083970948072/permalink/2340195572836906/",
					},
					// {
					// 	Type:  "web_url",
					// 	Title: "IOS 🍏",
					// 	Url:   "https://www.facebook.com/groups/1759083970948072/permalink/2325309854325478/",
					// },
				},
			},
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "SUPPORT_US":
		response := facebook.ResponseMessage{
			Text: "⚙️ Behind the scenes, our team dedicates significant time and resources to keep the bot running smoothly and to introduce new features that enhance your experience.\n" +
				"\n" +
				"🌟 Your contribution goes a long way in helping us cover server costs and invest in further development.\n" +
				"\n" +
				"☕ If you appreciate our efforts and would like to support us, consider buying us a coffee!\n" +
				"\n" +
				"https://buymeacoffee.com/pacshare",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "REVIEW":
		response := facebook.ResponseMessage{
			Text: "Help us to improve our service by review us 🙏, Write your opinions and your criticisms ✍️, We will be very happy to read your review 😊.\n" +
				"\n\n" +
				"To review us click here 👇\n" +
				"https://www.facebook.com/pacshare1/reviews/",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "HELP":
		response := facebook.ResponseMessage{
			Text: "If you have any questions or need any help about this service 📺, You can write It in the comments section of this post 💬, we will answer your questions as soon as possible 👍.\n" +
				"\n" +
				"https://fb.com/pacshare1/photos/127230276413440",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	}

	// redirect to specific regex handler

	// yt -> ^(YT_::_)(.+)$
	ytRegex := regexp.MustCompile("^(YT_::_)(.+)$")
	if match := ytRegex.FindStringSubmatch(postback.Payload); match != nil {
		return yt.RegexHundlerPostback(sender_psid, match[1:])
	}

	// wiki -> ^(WIKI_::_)(.+)$
	wikiRegex := regexp.MustCompile("^(WIKI_::_)(.+)$")
	if match := wikiRegex.FindStringSubmatch(postback.Payload); match != nil {
		return wiki.RegexHundlerPostback(match[1:])
	}

	// apk -> ^(APK_::_)(.+)$
	apkRegex := regexp.MustCompile("^(APK_::_)(.+)$")
	if match := apkRegex.FindStringSubmatch(postback.Payload); match != nil {
		return apk.RegexHundlerPostback(sender_psid, match[1:])
	}

	// pdf -> ^(PDF_::_)(.+)$
	pdfRegex := regexp.MustCompile("^(PDF_::_)(.+)$")
	if match := pdfRegex.FindStringSubmatch(postback.Payload); match != nil {
		return pdf.RegexHundlerPostback(sender_psid, match[1:])
	}

	//default
	return nil
}
