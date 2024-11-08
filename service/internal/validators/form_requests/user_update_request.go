package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// UserdUpdateRequest handles the validation for the RegisterRequest.
type UserdUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterdUpdateRequest creates a new instance of UserdUpdateRequest.
func NewUserdUpdateRequest() *UserdUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &UserdUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *UserdUpdateRequest) Validate(req *dtos.UpdateUserRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"name":     []string{"required", "min:3"},
		"username": []string{fmt.Sprintf("unique_ig:users,username,%d", req.ID)},
		"email":    []string{"required", "email", fmt.Sprintf("unique_ig:users,email,%d", req.ID)},
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
