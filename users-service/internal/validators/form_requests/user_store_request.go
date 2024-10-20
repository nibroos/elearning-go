package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// UserStoreRequest handles the validation for the RegisterRequest.
type UserStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of UserStoreRequest.
func NewUserStoreRequest() *UserStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &UserStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *UserStoreRequest) Validate(req *dtos.CreateUserRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"name":     []string{"required", "min:3"},
		"username": []string{"unique:users,username"},
		"email":    []string{"required", "email", "unique:users,email"},
		"password": []string{"required", "min:4"},
		"role_ids": []string{"required"},
	}

	messages := govalidator.MapData{
		"role_ids": []string{"required:The roles field is required."},
	}

	opts := govalidator.Options{
		Data:     req,
		Rules:    rules,
		Messages: messages,
	}

	v := govalidator.New(opts)
	mappedErrors := v.ValidateStruct()

	if len(mappedErrors) == 0 {
		return nil
	}

	errors := make(map[string]string)
	for field, err := range mappedErrors {
		errors[field] = err[0]
	}
	return errors
}
