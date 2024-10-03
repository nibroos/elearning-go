package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/users-service/internal/dtos"
)

var validate *validator.Validate
var db *sqlx.DB

func init() {
	validate = validator.New()

	// Register custom validation functions if needed
	validate.RegisterValidation("unique", uniqueValidator)
}

// uniqueValidator checks if a field value is unique in the database.
func uniqueValidator(fl validator.FieldLevel) bool {
	param := fl.Param()
	params := strings.Split(param, ",")
	if len(params) != 2 {
		return false
	}

	table := params[0]
	column := params[1]
	value := fl.Field().String()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
	err := db.Get(&count, query, value)
	if err != nil {
		return false
	}

	return count == 0
}

// ValidateCreateUserRequest validates the CreateUserRequest struct.
func ValidateCreateUserRequest(req *dtos.CreateUserRequest) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	errors := make(map[string]string)
	for _, err := range validationErrors {
		errors[err.Field()] = err.Tag()
	}
	return errors
}
