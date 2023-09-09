package inventoryRoutes

import (
	"practice-api/controllers/practice"
	practiceHandler "practice-api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupInvRoutes(router fiber.Router) {
	app := router.Group("/practice")
	// Register a user
	app.Post("/register", practiceHandler.RegisterUser)
	app.Post("/login", practiceHandler.LoginUser)

	app.Use(practiceHandler.Authenticate)
	app.Get("/getme", practiceHandler.GetMe)
	app.Get("/", practice.Index)
	app.Get("/:id", practice.Show)
	app.Post("/", practice.Create)
	app.Put("/update/:id", practice.Update)
	app.Delete("/delete/:id", practice.Delete)
	app.Post("/logout", practiceHandler.LogoutUser)
}