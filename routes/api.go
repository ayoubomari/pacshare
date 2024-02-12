package routes

import (
	"github.com/ayoubomari/pacshare/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegisterAPI(api fiber.Router) {
	registerwebhooks(api)
}

func registerwebhooks(api fiber.Router) {
	webhooks := api.Group("/webhook")

	webhooks.Get("/", controllers.FacebookGet)

	webhooks.Post("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("hello world")
	})
}
