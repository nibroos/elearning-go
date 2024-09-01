package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nibroos/elearning-go/users-service/internal/config"
	controller "github.com/nibroos/elearning-go/users-service/internal/controller/rest"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
	"github.com/nibroos/elearning-go/users-service/internal/routes"
	"github.com/nibroos/elearning-go/users-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize the SQLx database connection
	sqlDB, err := sqlx.Connect("postgres", config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to the SQL database: %v", err)
	}

	// Initialize the Gorm database connection
	gormDB, err := gorm.Open(postgres.Open(config.GetDatabaseURL()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the Gorm database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(gormDB, sqlDB)

	// Initialize controllers
	userController := controller.NewUserController(service.NewUserService(userRepo))

	// Initialize Fiber app
	app := fiber.New()

	// Setup REST routes
	routes.SetupRoutes(app, userController)

	var wg sync.WaitGroup

	// Start REST server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Listen(":4001"); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	// Start gRPC server
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if err := controller.RunGRPCServer(userRepo); err != nil {
	// 		log.Fatalf("Failed to run gRPC server: %v", err)
	// 	}
	// }()

	// Wait for all servers to exit
	wg.Wait()
}
