package routes

import (
	"github.com/gofiber/fiber/v2"
	rest "github.com/nibroos/elearning-go/users-service/internal/controller/rest"
	"github.com/nibroos/elearning-go/users-service/internal/middleware"
)

// SetupRoutes sets up the REST routes for the user service.
func SetupRoutes(app *fiber.App, userController *rest.UserController, seederController *rest.SeederController) {

	// Public routes
	app.Get("/api/v1/users/test", func(c *fiber.Ctx) error {
		return c.SendString("REST Users Service!")
	})

	version := app.Group("/api/v1")
	auth := version.Group("/auth")

	version.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Users-Service is running",
		})
	})

	auth.Post("/login", userController.Login)
	auth.Post("/register", userController.Register)

	// Protected routes
	app.Use(middleware.JWTMiddleware())

	users := version.Group("/users")
	users.Post("/index-user", userController.GetUsers)
	users.Post("/show-user", userController.GetUserByID)
	users.Post("/create-user", userController.CreateUser)
	users.Post("/update-user", userController.UpdateUser)
	users.Post("/delete-user", userController.DeleteUser)
	users.Post("/restore-user", userController.RestoreUser)

	// Seeder route
	version.Post("/seeders/run", seederController.RunSeeders)
}
