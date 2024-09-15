package middleware

import (
	"fmt"
	"log"
	"runtime"
	"strconv"

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

func ConvertRequestToFilters() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Check if the content type is JSON
		if ctx.Get("Content-Type") == "application/json" {
			var requestBody map[string]interface{}
			if err := ctx.BodyParser(&requestBody); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Failed to parse request body",
					"status":  "error",
					"err":     err.Error(),
				})
			}

			filters := make(map[string]string)
			for key, value := range requestBody {
				switch v := value.(type) {
				case string:
					filters[key] = v
				case int:
					filters[key] = strconv.Itoa(v)
				case float64:
					filters[key] = strconv.FormatFloat(v, 'f', -1, 64)
				default:
					log.Printf("Unsupported type for key %s: %T", key, v)
				}
			}

			ctx.Locals("filters", filters)
		}

		return ctx.Next()
	}
}
