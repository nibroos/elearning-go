package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Meta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
	Message string      `json:"message"`
	Status  int16       `json:"status"`
	Errors  interface{} `json:"errors"`
}

// Nullable is a generic type that can handle null values for different data types.
type Nullable[T any] struct {
	Value *T
}

// ContextKey is a type for context keys used in this package
type ContextKey string

const (
	// ResponseWriterKey is the context key for the http.ResponseWriter
	ResponseWriterKey ContextKey = "ResponseWriter"
)

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
func WrapResponse(data interface{}, pagination *Meta, message string, status int16, errors ...interface{}) Response {
	meta := Meta{}
	if pagination != nil {
		meta = *pagination
	}

	return Response{
		Data:    data,
		Meta:    &meta,
		Message: message,
		Status:  status,
		Errors:  errors,
	}
}

func GetResponse(ctx *fiber.Ctx, data interface{}, pagination *Meta, message string, status int16, errors ...interface{}) error {
	meta := Meta{}
	if pagination != nil {
		meta = *pagination
	}

	response := Response{
		Data:    data,
		Meta:    &meta,
		Message: message,
		Status:  status,
		Errors:  errors,
	}

	return SendResponse(ctx, response, int(status))
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

// // DD takes multiple values, creates a JSON response, and stops execution.
// func DD(c *fiber.Ctx, values ...interface{}) error {
// 	// Create a map to hold the values
// 	response := fiber.Map{}

// 	// Dynamically add the passed values to the response
// 	for i, value := range values {
// 		// The key will be "value_0", "value_1", etc.
// 		key := fmt.Sprintf("value_%d", i)
// 		response[key] = value
// 	}

// 	// Return a JSON response with status 200 and stop further execution
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status":  "debug",
// 		"message": "Debugging Output",
// 		"data":    response,
// 	})
// }

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

// GetIntOrDefault retrieves an int value from a variable or map-based on the provided default value.
// It accepts both direct int values and string values (which are converted to int).
// If the value is empty, invalid, or the key does not exist, it returns the provided default value.
func GetIntOrDefault(value interface{}, defaultValue int) int {
	// Check if value is an int directly
	if intValue, ok := value.(int); ok {
		return intValue
	}

	// Check if value is a string and convert it to int
	if str, ok := value.(string); ok {
		if intValue, err := strconv.Atoi(str); err == nil {
			return intValue
		}
	}

	// Check if value is a map
	if reflect.TypeOf(value).Kind() == reflect.Map {
		// Ensure the value is a map of strings to interfaces
		if m, ok := value.(map[string]interface{}); ok {
			// Try to retrieve value from map and check if it is a string or int
			for _, v := range m {
				if intValue, ok := v.(int); ok {
					return intValue
				}
				if str, ok := v.(string); ok {
					if intValue, err := strconv.Atoi(str); err == nil {
						return intValue
					}
				}
			}
		}
	}

	return defaultValue
}

// DD is a helper function to dump the value of a variable, stop the process, and optionally send a response to the client.
func DD(ctx context.Context, value interface{}) {
	// Print the value to the console
	fmt.Printf("%+v\n", value)

	// Convert the value to JSON
	jsonValue, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal value:", err)
		os.Exit(1)
	}

	// Print the JSON value to the console
	fmt.Println(string(jsonValue))

	// Stop the process
	os.Exit(0)
}

func CreatePaginationMeta(filters map[string]string, total int) *Meta {
	currentPage := GetIntOrDefault(filters["page"], 1)
	perPage := GetIntOrDefault(filters["per_page"], 10)
	lastPage := (total + perPage - 1) / perPage

	return &Meta{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		LastPage:    lastPage,
	}
}

func ExecuteSeeders(db *sql.DB, seedFiles []string) error {
	for _, file := range seedFiles {
		err := executeSQLFile(db, file)
		if err != nil {
			return fmt.Errorf("error executing %s: %v", file, err)
		}
		fmt.Printf("Executed %s successfully\n", file)
	}
	return nil
}

func executeSQLFile(db *sql.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sqlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlBytes))
	return err
}

func generateBcryptHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// BodyParserWithNull converts empty strings to null and parses the request body into the provided struct.
func BodyParserWithNull(ctx *fiber.Ctx, out interface{}) error {
	// Parse the request body into a map
	var body map[string]interface{}
	if err := json.Unmarshal(ctx.Body(), &body); err != nil {
		return err
	}

	// Convert empty strings to null in the map
	for key, value := range body {
		if str, ok := value.(string); ok && str == "" {
			body[key] = nil
		}
	}

	// Marshal the modified body back to JSON
	modifiedBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Unmarshal the modified body into the provided struct
	if err := json.Unmarshal(modifiedBody, out); err != nil {
		return err
	}

	// Convert empty strings to null in the struct fields
	convertEmptyStringsToNull(out)

	return nil
}

// convertEmptyStringsToNull uses reflection to convert empty strings to null in struct fields.
func convertEmptyStringsToNull(out interface{}) {
	v := reflect.ValueOf(out).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			field.Set(reflect.Zero(field.Type()))
		}
	}
}

// StringPointerToString converts a string pointer to a string, returning an empty string if the pointer is nil.
func StringPointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if reflect.ValueOf(value).IsZero() {
		n.Value = nil
	} else {
		n.Value = &value
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.Value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(*n.Value)
}

// Scan implements the sql.Scanner interface.
func (n *Nullable[T]) Scan(value interface{}) error {
	if value == nil {
		n.Value = nil
		return nil
	}
	val, ok := value.(T)
	if !ok {
		return errors.New("type assertion failed")
	}
	n.Value = &val
	return nil
}
