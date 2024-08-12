package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=secret dbname=customerdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	registerRoutes(e)

	log.Println("Customer Service running on port 8081")
	log.Fatal(e.Start(":8081"))
}
