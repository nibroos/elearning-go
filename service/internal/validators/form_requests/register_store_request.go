package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// RegisterStoreRequest handles the validation for the RegisterRequest.
type RegisterStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of RegisterStoreRequest.
func NewRegisterStoreRequest() *RegisterStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &RegisterStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *RegisterStoreRequest) Validate(req *dtos.RegisterRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"name":     []string{"required", "min:3"},
		"username": []string{"unique:users,username"},
		"email":    []string{"required", "email", "unique:users,email"},
		"password": []string{"required", "min:4"},
	}

	opts := govalidator.Options{
		Data:  req,
		Rules: rules,
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
