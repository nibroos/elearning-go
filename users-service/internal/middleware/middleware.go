package middleware

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler middleware
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next() // Execute the next handler
		if err != nil {
			// Capture the stack trace
			_, file, line, _ := runtime.Caller(1)
			stackTrace := fmt.Sprintf("%s:%d", file, line)

			// Log the error and stack trace
			log.Printf("Error: %v\nStack Trace: %s\n", err, stackTrace)

			// Return a JSON response with the error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
				"stack":   stackTrace, // Optionally include stack trace
			})
		}
		return nil
	}
}
