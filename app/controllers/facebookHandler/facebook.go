package facebookHandler

import (
	"fmt"
	"os"

	"github.com/ayoubomari/pacshare/app/models/facebook"
	"github.com/gofiber/fiber/v2"
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
			} else if webHookEvent.Message.MID != "" && len(webHookEvent.Message.Attachments) > 0 { // handle attachments
				err := handleAttachments(sender_psid, webHookEvent.Message)
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
