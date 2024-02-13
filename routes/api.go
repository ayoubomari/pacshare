package routes

import (
	"github.com/ayoubomari/pacshare/app/controllers/facebook"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(api fiber.Router) {
	registerwebhooks(api)
}

func registerwebhooks(api fiber.Router) {
	webhooks := api.Group("/webhook")

	webhooks.Get("/", facebook.FacebookGet)

	webhooks.Post("/", facebook.FacebookPost)
}
