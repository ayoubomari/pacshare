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
		response := facebook.ResponseMessage{
			Text: "Hi, how can I help you?",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "Greeting":
		response := facebook.ResponseMessage{
			Text: "Hi, how can I help you?",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil

	// for menu
	case "YOUTUBE":
		response := facebook.ResponseMessage{
			Text: "for youtube",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "REVIEW":
		response := facebook.ResponseMessage{
			Text: "for review",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "HELP":
		response := facebook.ResponseMessage{
			Text: "for help",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "SUPPORT_US":
		response := facebook.ResponseMessage{
			Text: "for SUPPORT_US",
		}
		err := facebookSender.CallSendAPI(sender_psid, response)
		if err != nil {
			return fmt.Errorf("handleMessage: %w", err)
		}
		return nil
	case "DOWNLOAD_PACSHARE_APP":
		response := facebook.QuickReplyMessage{
			Text: "Download PacShare App",
			QuickReplies: []facebook.QuickReplyResponse{
				{
					ContentType: "text",
					Title:       "Android",
					Payload:     "DOWNLOAD_PACSHARE_APP_ANDROID",
					ImageURL:    "https://pacshare.omzor.com/static_src/imgs/android_logo.png",
				},
				{
					ContentType: "text",
					Title:       "IOS",
					Payload:     "DOWNLOAD_PACSHARE_APP_IOS",
					ImageURL:    "https://pacshare.omzor.com/static_src/imgs/ios_logo.png",
				},
			},
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
