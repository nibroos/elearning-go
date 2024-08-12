package main

import (
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App) {
	app.Post("/login", loginHandler)
	app.Post("/register", registerHandler)
}
