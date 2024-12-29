package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nibroos/elearning-go/service/internal/config"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/middleware"
	"github.com/nibroos/elearning-go/service/internal/routes"
	"github.com/nibroos/elearning-go/service/internal/validators"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Determine the environment (production or test)
	env := os.Getenv("APP_ENV")
	var dbURL string
	if env == "test" {
		dbURL = config.GetTestDatabaseURL()
	} else {
		dbURL = config.GetDatabaseURL()
	}

	// Initialize the SQLx database connection
	sqlDB, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the SQL database: %v", err)
	}

	// Configure SQLx connection pool
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDB.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	// Initialize the Gorm database connection
	gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the Gorm database: %v", err)
	}

	// Configure Gorm connection pool
	sqlDBGorm, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get Gorm DB instance: %v", err)
	}
	sqlDBGorm.SetMaxOpenConns(100)          // Maximum number of open connections
	sqlDBGorm.SetMaxIdleConns(10)           // Maximum number of idle connections
	sqlDBGorm.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	// Initialize the Redis client
	if env == "test" {
		config.InitRedisClientTest()
	} else {
		config.InitRedisClient()
	}

	// Fetch needed data from the database and cache it in Redis
	config.FetchCachedData(context.Background(), sqlDB)

	// Initialize the validator with the database connection
	validators.InitValidator(sqlDB)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Attach middleware
	app.Use(middleware.ConvertEmptyStringsToNull())
	app.Use(middleware.ConvertRequestToFilters())

	// Setup REST routes
	routes.SetupRoutes(app, gormDB, sqlDB)

	// Protect routes with JWT middleware
	// app.Use(middleware.JWTMiddleware())

	var wg sync.WaitGroup

	// Check if the service type is "scheduler"
	if os.Getenv("SERVICE_TYPE") == "scheduler" {
		// Start the scheduler in a separate goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Initialize the cron scheduler
			cron := cron.New()
			schedulerController := rest.NewSchedulerController(cron, gormDB, sqlDB)

			// Reload schedules from the database
			if err := schedulerController.ReloadSchedules(); err != nil {
				log.Printf("Failed to reload schedules: %v", err)
				return
			}

			// Start the cron scheduler
			cron.Start()
			log.Println("Scheduler started successfully")
		}()
	} else {
		// Start REST server
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := app.Listen(":4001"); err != nil {
				log.Fatalf("Failed to start REST server: %v", err)
			}

			println("Server started on :4001")
		}()

		// Start gRPC server
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		// 	if err := runGRPCServer(grpcUserController); err != nil {
		// 		log.Fatalf("Failed to run gRPC server: %v", err)
		// 	}
		// }()
	}

	// Wait for all servers to exit
	wg.Wait()
}

// func runGRPCServer(grpcUserController grpcController.GRPCUserController) error {
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		return err
// 	}

// 	server := grpc.NewServer(
// 		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor()),
// 	)

// 	grpcController.RegisterUserServiceServer(server, grpcUserController)

// 	log.Printf("gRPC server listening on %v", lis.Addr())
// 	return server.Serve(lis)
// }
