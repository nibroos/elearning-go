package utils

import (
	"fmt"
	"reflect"
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
	Status  string      `json:"status"`
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
func WrapResponse(data interface{}, pagination *PaginationMeta, message string, status string) Response {
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
