package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupHealthRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Users-Service is running",
		})
	})
}
