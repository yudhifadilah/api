package router

import (
	routes "practice-api/route/practice"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	practiceapi := app.Group("/pract", logger.New())

	// Setup the Node Routes
	routes.SetupInvRoutes(practiceapi)
}