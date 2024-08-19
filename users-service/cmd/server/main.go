package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nibroos/elearning-go/users-service/internal/config"
	"github.com/nibroos/elearning-go/users-service/internal/controller"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
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

    if err := controller.RunGRPCServer(userRepo); err != nil {
        log.Fatalf("Failed to run gRPC server: %v", err)
    }
}
