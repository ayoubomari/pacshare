package controllers

import (
	"fmt"
	"os"

	"github.com/ayoubomari/pacshare/app/models"
	"github.com/gofiber/fiber/v2"
)

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

func FacebookPost(c *fiber.Ctx) error {
	var body models.FacebookWebhookBody

	// Checks this is an event from a page subscription
	if body.Object == "page" {
		for _, entry := range body.Entry {
			// Gets the body of the webhook event
			webHookEvent := entry.Messaging[0]

			// Get the sender PSID
			sender_psid := webHookEvent.Sender.ID

			if webHookEvent.Message.MID != "" {
				handleMessage(sender_psid, webHookEvent.Message)
			} else if webHookEvent.PostBack.Title != "" {
				handlePostback(sender_psid, webHookEvent.PostBack)
			}
		}

		// Returns a '200 OK' response to all requests
		c.Status(200).SendString("EVENT_RECEIVED")
	}

	// forbid any thing else
	return c.SendStatus(fiber.StatusForbidden)
}

func handleMessage(sender_psid string, message models.Message) {
	fmt.Println("from post message")
}

func handlePostback(sender_psid string, message models.PostBack) {
	fmt.Println("from post back")
}
