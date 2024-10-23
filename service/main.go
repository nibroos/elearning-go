package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nibroos/elearning-go/service/internal/config"
	restController "github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/middleware"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/routes"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/nibroos/elearning-go/service/internal/validators"
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

	// Initialize the validator with the database connection
	validators.InitValidator(sqlDB)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Attach middleware
	app.Use(middleware.ConvertEmptyStringsToNull())
	app.Use(middleware.ConvertRequestToFilters())

	// Initialize repositories
	userRepo := repository.NewUserRepository(gormDB, sqlDB)

	// Initialize controllers
	userController := restController.NewUserController(service.NewUserService(userRepo))
	seederController := restController.NewSeederController(sqlDB.DB)
	// grpcUserController := grpcController.GRPCUserController(grpcServer, service.NewUserService(userRepo))

	// Setup REST routes
	routes.SetupRoutes(app, userController, seederController)

	// Protect routes with JWT middleware
	// app.Use(middleware.JWTMiddleware())

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
	// 	if err := runGRPCServer(grpcUserController); err != nil {
	// 		log.Fatalf("Failed to run gRPC server: %v", err)
	// 	}
	// }()

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
