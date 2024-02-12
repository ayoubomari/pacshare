package main

import (
	"flag"
	"log"

	"github.com/ayoubomari/pacshare/config"
	"github.com/ayoubomari/pacshare/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":5000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// api handler
	routes.RegisterAPI(app)

	// Setup static files
	app.Static("/", "./public")

	// Listen on port 5000
	log.Fatal(app.Listen(":" + config.GetServerPort())) // go run app.go -port=:5000
}
