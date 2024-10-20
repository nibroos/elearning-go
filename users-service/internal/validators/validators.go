package validators

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/utils"
	"github.com/thedevsaddam/govalidator"
)

var validate *validator.Validate
var db *sqlx.DB

func InitValidator(database *sqlx.DB) {
	db = database
	validate = validator.New()

	// Register custom validation functions if needed
	validate.RegisterValidation("unique", uniqueValidator)

	// Register custom validation rules
	govalidator.AddCustomRule("unique", uniqueRule)
}

// uniqueValidator checks if a field value is unique in the database.
func uniqueValidator(fl validator.FieldLevel) bool {

	utils.DD(map[string]interface{}{
		"perPage": fl.Field().Interface(),
		"fl":      fl,
	})

	value := fl.Field().Interface()

	// Debugging: Dump the value
	utils.DD(value)

	// If the value is empty or null, pass the validation
	if value == nil || reflect.ValueOf(value).IsZero() {
		return true
	}

	// Convert value to string for query
	valueStr, ok := value.(string)
	if !ok {
		return false
	}

	param := fl.Param()
	params := strings.Split(param, ",")
	if len(params) != 2 {
		return false
	}

	table := params[0]
	column := params[1]

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", table, column)
	err := db.Get(&count, query, valueStr)
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

func ValidateUpdateUserRequest(req *dtos.UpdateUserRequest) map[string]string {
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

func ValidateRegisterRequest(req *dtos.RegisterRequest, ctx context.Context) map[string]string {
	err := validate.Struct(req)
	// utils.DD(ctx, map[string]interface{}{
	// 	"req":      req,
	// 	"testbool": true,
	// 	"err":      err,
	// })
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		return errors
	}
	return nil
}

// uniqueRule checks if a field value is unique in the database.
func uniqueRule(field string, rule string, message string, value interface{}) error {
	valueStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value type")
	}

	params := strings.Split(rule, ":")
	if len(params) != 2 {
		return fmt.Errorf("invalid rule format")
	}

	tableColumn := strings.Split(params[1], ",")
	if len(tableColumn) != 2 {
		return fmt.Errorf("invalid table and column format")
	}

	table := tableColumn[0]
	column := tableColumn[1]

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = $1", table, column)
	err := db.Get(&count, query, valueStr)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("the %s has already been taken", field)
	}

	return nil
}
