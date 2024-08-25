package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/users-service/internal/controller"
)

// SetupRoutes sets up the REST routes for the user service.
func SetupRoutes(app *fiber.App, userController *controller.UserController) {
	app.Get("/api/v1/users/test", func(c *fiber.Ctx) error {
		return c.SendString("REST Users Service!")
	})

	users := app.Group("/api/v1/users")

	users.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Users-Service is running",
		})
	})

	// Define other routes here
	users.Get("/index-user", userController.CreateUser)
	users.Get("/show-user", userController.GetUser)
	users.Get("/create-user", userController.CreateUser)
	users.Get("/update-user", userController.UpdateUser)
	users.Get("/delete-user", userController.DeleteUser)
	
}
