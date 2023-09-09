package main

import (

	"practice-api/configs"
	router "practice-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Connect to the Database
	configs.ConnectDB()

	// Start a new fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PUT,DELETE",
	}))

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup the router
	router.SetupRoutes(app)

	// Listen on PORT 3000
	app.Listen(":3000")
}