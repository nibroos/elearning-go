package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// IdentifierStoreRequest handles the validation for the RegisterRequest.
type IdentifierStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of IdentifierStoreRequest.
func NewIdentifierStoreRequest() *IdentifierStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &IdentifierStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *IdentifierStoreRequest) Validate(req *dtos.CreateIdentifierRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"type_identifier_id": []string{"required", "exists:modules,id"},
		"ref_num":            []string{"required"},
		"status":             []string{"required"},
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
