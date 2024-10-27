package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// ModuleUpdateRequest handles the validation for the RegisterRequest.
type ModuleUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of ModuleUpdateRequest.
func NewModuleUpdateRequest() *ModuleUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &ModuleUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *ModuleUpdateRequest) Validate(req *dtos.UpdateModuleRequest, ctx context.Context) map[string]string {
	// utils.DD(req)

	// TODO Fix the unique_ig rule
	rules := govalidator.MapData{
		"name":        []string{"required", fmt.Sprintf("unique_ig:modules,name,%d", req.ID)},
		"class_id":    []string{"required", fmt.Sprintf("unique_ig:modules,class_id,%d", req.ID)},
		"description": []string{},
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
