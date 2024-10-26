package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/middleware"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

// SetupRoutes sets up the REST routes for the user service.
func SetupRoutes(app *fiber.App, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	// Public routes
	app.Get("/api/v1/users/test", func(c *fiber.Ctx) error {
		return c.SendString("REST Users Service!")
	})

	version := app.Group("/api/v1")
	auth := version.Group("/auth")

	version.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Service is running",
		})
	})

	// Setup auth routes
	auth.Post("/login", rest.NewUserController(service.NewUserService(repository.NewUserRepository(gormDB, sqlDB))).Login)
	auth.Post("/register", rest.NewUserController(service.NewUserService(repository.NewUserRepository(gormDB, sqlDB))).Register)

	// Protected routes
	app.Use(middleware.JWTMiddleware())

	// Grouped routes
	users := version.Group("/users")
	SetupUserRoutes(users, gormDB, sqlDB)

	subscribes := version.Group("/subscribes")
	SetupSubscribeRoutes(subscribes, gormDB, sqlDB)

	// Seeder route
	version.Post("/seeders/run", rest.NewSeederController(sqlDB.DB).RunSeeders)
}