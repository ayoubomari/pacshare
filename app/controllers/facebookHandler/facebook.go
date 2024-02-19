package facebookHandler

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/gofiber/fiber/v2"

	"github.com/ayoubomari/pacshare/app/controllers/facebookSender"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/apk"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/pdf"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/wiki"
	"github.com/ayoubomari/pacshare/app/controllers/scrapers/yt"
)

// fiber reqest handler for GET: /webhook
func FacebookGet(c *fiber.Ctx) error {
	// Your verify token. Should be a random string.
	verifyToken := os.Getenv("FACEBOOK_PAGE_ACCESS_TOKEN")

	// Parse the query params
	mode := c.Query("hub.mode", "")
	token := c.Query("hub.verify_token", "")
	challenge := c.Query("hub.challenge", "")

	// Checks if a token and mode is in the query string of the request
	if mode != "" && token != "" {
		// Checks the mode and token sent is correct
		if mode == "subscribe" && token == verifyToken {
			// Responds with the challenge token from the request
			//console.log('WEBHOOK_VERIFIED');
			return c.SendString(challenge)
		} else {
			// Responds with '403 Forbidden' if verify tokens do not match
			return c.SendStatus(fiber.StatusForbidden)
		}
	}
	return c.SendStatus(fiber.StatusForbidden)
}

// fiber reqest handler for POST: /webhook
func FacebookPost(c *fiber.Ctx) error {
	var body facebook.FacebookWebhookBody
	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("FacebookPost: %w", err)
	}

	// Checks this is an event from a page subscription
	if body.Object == "page" {
		for _, entry := range body.Entry {
			// Gets the body of the webhook event
			webHookEvent := entry.Messaging[0]

			// Get the sender PSID
			sender_psid := webHookEvent.Sender.ID

			if webHookEvent.Message.MID != "" && webHookEvent.Message.Quick_reply.Payload != "" { // handle quick replay
				err := handleQuickReplay(sender_psid, webHookEvent.Message)
				if err != nil {
					fmt.Println(err)
					return nil
				}
			} else if webHookEvent.Message.MID != "" { // handle message
				err := handleMessage(sender_psid, webHookEvent.Message)
				if err != nil {
					fmt.Println(err)
					return nil
				}
			} else if webHookEvent.PostBack.Title != "" { // handle postback
				err := handlePostback(sender_psid, webHookEvent.PostBack)
				if err != nil {
					fmt.Println(err)
					return nil
				}
			}
		}

		// Returns a '200 OK' response to all requests
		return c.Status(200).SendString("EVENT_RECEIVED")
	}

	// forbid any thing else
	return c.SendStatus(fiber.StatusForbidden)
}

// handle quick replay message
func handleQuickReplay(sender_psid string, message facebook.Message) error {
	return nil
}

// handle message come from facebook (text message)
func handleMessage(sender_psid string, message facebook.Message) error {
	// work with lower case from now one
	trimmedMessage := strings.Trim(message.Text, " ")

	// redirect to specific regex handler

	// yt -> ^(.yt) (.+)$
	ytRegex := regexp.MustCompile("^(.yt) (.+)$")
	if match := ytRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return yt.RegexHundlerMessage(sender_psid, match[1:])
	}

	// wiki -> ^(.wiki)( |-)([a-z]{2}) (.+)$
	wikiRegex := regexp.MustCompile("^(.wiki)( |-)([a-z]{2}) (.+)$")
	if match := wikiRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return wiki.RegexHundlerMessage(match[1:])
	}

	// apk -> ^(.apk) (.+)$
	apkRegex := regexp.MustCompile("^(.apk) (.+)$")
	if match := apkRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return apk.RegexHundlerMessage(match[1:])
	}

	// pdf -> ^(.pdf) (.+)$
	pdfRegex := regexp.MustCompile("^(.pdf) (.+)$")
	if match := pdfRegex.FindStringSubmatch(trimmedMessage); match != nil {
		return pdf.RegexHundlerMessage(match[1:])
	}

	//default
	return yt.RegexHundlerMessage(
		sender_psid,
		[]string{
			".wiki-en",
			trimmedMessage,
		},
	)
}

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
		response := facebook.ResponseMediaAttachment{
			Type: "file",
			Payload: facebook.WebhookBodyAttachmentPayload{
				URL:         "pacshare.omzor.com/static_src/apks/pacshare.apk",
				Is_reusable: false,
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
		return apk.RegexHundlerPostback(match[1:])
	}

	// pdf -> ^(PDF_::_)(.+)$
	pdfRegex := regexp.MustCompile("^(PDF_::_)(.+)$")
	if match := pdfRegex.FindStringSubmatch(postback.Payload); match != nil {
		return pdf.RegexHundlerPostback(match[1:])
	}

	//default
	return nil
}
