package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

func loginHandler(c *fiber.Ctx) error {
	// Implement login logic
	return c.SendString("Login")
}

func registerHandler(c *fiber.Ctx) error {
	var user User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	go func() {
		_, err := DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
		if err != nil {
			log.Println("Error inserting user:", err)
		}
	}()

	return c.SendString("User registered")

	// return c.SendString("Register")
}
