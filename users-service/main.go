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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db, err := sqlx.Connect("postgres", config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
    userController := controller.NewUserController(userRepo)

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := controller.RunGRPCServer(userRepo); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

    wg.Wait()
}
