package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Meta struct {
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta,omitempty"`
	Message string      `json:"message"`
	Status  int16       `json:"status"`
}

type StringOrInt struct {
	Value int
}

// JSONError formats and returns an error response
func JSONError(ctx *fiber.Ctx, status int, err error) error {
	return ctx.Status(status).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
func WrapResponse(data interface{}, pagination *PaginationMeta, message string, status int16) Response {
	meta := Meta{}
	if pagination != nil {
		meta.Pagination = *pagination
	}

	return Response{
		Data:    data,
		Meta:    meta,
		Message: message,
		Status:  status,
	}
}

func SendResponse(ctx *fiber.Ctx, response Response, statusCode int) error {
	return ctx.Status(statusCode).JSON(response)
}

func AtoiDefault(str string, def int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return value
}

// ConvertStructToMap function converts a struct to a map
func ConvertStructToMap(filters interface{}) map[string]string {
	result := make(map[string]string)

	v := reflect.ValueOf(filters)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		key := typeOfS.Field(i).Tag.Get("json")
		value := v.Field(i).Interface()

		// Convert value to string
		switch v := value.(type) {
		case int:
			result[key] = strconv.Itoa(v)
		case string:
			result[key] = v
		default:
			result[key] = fmt.Sprintf("%v", v)
		}
	}

	return result
}

// GenerateIndexName creates a standardized index name based on table and columns
func GenerateIndexName(table string, columns ...string) string {
	return fmt.Sprintf("idx_%s_%s", table, strings.Join(columns, "_"))
}

func (s *StringOrInt) UnmarshalJSON(data []byte) error {
	var strValue string
	if err := json.Unmarshal(data, &strValue); err == nil {
		val, err := strconv.Atoi(strValue)
		if err != nil {
			return err
		}
		s.Value = val
		return nil
	}

	var intValue int
	if err := json.Unmarshal(data, &intValue); err != nil {
		return err
	}
	s.Value = intValue
	return nil
}

func DefaultString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func DefaultInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

// DD takes multiple values, creates a JSON response, and stops execution.
func DD(c *fiber.Ctx, values ...interface{}) error {
	// Create a map to hold the values
	response := fiber.Map{}

	// Dynamically add the passed values to the response
	for i, value := range values {
		// The key will be "value_0", "value_1", etc.
		key := fmt.Sprintf("value_%d", i)
		response[key] = value
	}

	// Return a JSON response with status 200 and stop further execution
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "debug",
		"message": "Debugging Output",
		"data":    response,
	})
}

// ErrorWithLocation returns an error message with file and line number information.
func ErrorWithLocation(err error) string {
	// Retrieve the program counter, file, and line number
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("error: %v at %s:%d", err, file, line)
}

// GetStringOrDefault retrieves a string value from a variable or map-based on the provided default value.
// It accepts both direct values and map-based values.
// If the value is empty or the key does not exist, it returns the provided default value.
func GetStringOrDefault(value interface{}, defaultValue string) string {
	// Check if value is a string directly
	if str, ok := value.(string); ok {
		if str != "" {
			return str
		}
		return defaultValue
	}

	// Check if value is a map
	if reflect.TypeOf(value).Kind() == reflect.Map {
		// Ensure the value is a map of strings to interfaces
		if m, ok := value.(map[string]interface{}); ok {
			// Try to retrieve value from map and check if it is a string
			if v, exists := m["order_column"]; exists {
				if str, ok := v.(string); ok && str != "" {
					return str
				}
			}
		}
	}

	return defaultValue
}
