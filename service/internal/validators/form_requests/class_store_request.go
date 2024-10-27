package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// ClassStoreRequest handles the validation for the RegisterRequest.
type ClassStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of ClassStoreRequest.
func NewClassStoreRequest() *ClassStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &ClassStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *ClassStoreRequest) Validate(req *dtos.CreateClassRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"name":         []string{"required", "unique:classes,name"},
		"description":  []string{},
		"subscribe_id": []string{"required", "numeric", "exists:subscribes,id"},
		"incharge_id":  []string{"required", "numeric", "exists:users,id"},
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
