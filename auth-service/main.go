package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// var err error
	// db, err = sql.Open("postgres", "user=postgres password=secret dbname=authdb sslmode=disable")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	app := fiber.New()

	registerRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "4001"
	}

	log.Println("Auth Service running on port " + port)
	log.Fatal(app.Listen(":" + port))
}
