package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// SubscribeStoreRequest handles the validation for the RegisterRequest.
type SubscribeStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of SubscribeStoreRequest.
func NewSubscribeStoreRequest() *SubscribeStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &SubscribeStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *SubscribeStoreRequest) Validate(req *dtos.CreateSubscribeRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"name":        []string{"required", "min:3", "unique:subscribes,name"},
		"description": []string{"required"},
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
