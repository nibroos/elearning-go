package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// ModuleStoreRequest handles the validation for the RegisterRequest.
type ModuleStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of ModuleStoreRequest.
func NewModuleStoreRequest() *ModuleStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &ModuleStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *ModuleStoreRequest) Validate(req *dtos.CreateModuleRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"class_id":    []string{"required", "exists:classes,id"},
		"name":        []string{"required", "unique:modules,name"},
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
