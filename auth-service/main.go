package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("postgres", "user=postgres password=secret dbname=authdb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    app := fiber.New()

    registerRoutes(app)

    log.Println("Auth Service running on port 8080")
    log.Fatal(app.Listen(":8080"))
}