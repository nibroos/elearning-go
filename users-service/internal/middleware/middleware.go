package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler middleware
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := http.StatusInternalServerError
	message := "Internal server error"

	if err == sql.ErrNoRows {
		code = http.StatusNotFound
		message = "No result found"
	} else if e, ok := err.(*fiber.Error); ok {
		// Use Fiber's default error message
		code = e.Code
		message = e.Message
	}

	// Capture the stack trace
	_, file, line, _ := runtime.Caller(1)
	stackTrace := fmt.Sprintf("%s:%d", file, line)

	// Log the error and stack trace
	log.Printf("Error: %v\nStack Trace: %s\n", err, stackTrace)

	// Return a JSON response with the error
	return ctx.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": message,
		"errors":  err.Error(),
		// "stack":   stackTrace, // Optionally include stack trace
	})
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
					if v == "" {
						requestBody[key] = nil
					} else {
						filters[key] = v
					}
				case int:
					filters[key] = strconv.Itoa(v)
				case float64:
					filters[key] = strconv.FormatFloat(v, 'f', -1, 64)
				default:
					log.Printf("Unsupported type for key %s: %T", key, v)
				}
			}

			// Marshal the modified body back to JSON
			modifiedBody, err := json.Marshal(requestBody)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process request body"})
			}

			// Replace the request body with the modified body
			ctx.Request().SetBody(modifiedBody)

			ctx.Locals("filters", filters)
		}

		return ctx.Next()
	}
}

func ConvertEmptyStringsToNull() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Parse the request body into a map
		var body map[string]interface{}
		if err := json.Unmarshal(ctx.Body(), &body); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Convert empty strings to null
		for key, value := range body {
			if str, ok := value.(string); ok && str == "" {
				body[key] = nil
			}
		}

		// Marshal the modified body back to JSON
		modifiedBody, err := json.Marshal(body)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process request body"})
		}

		// Replace the request body with the modified body
		ctx.Request().SetBody(modifiedBody)

		return ctx.Next()
	}
}
