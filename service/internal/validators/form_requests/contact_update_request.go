package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// ContactUpdateRequest handles the validation for the RegisterRequest.
type ContactUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of ContactUpdateRequest.
func NewContactUpdateRequest() *ContactUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &ContactUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *ContactUpdateRequest) Validate(req *dtos.UpdateContactRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"type_identifier_id": []string{"exists:mix_values,id"},
		"ref_num":            []string{"required", fmt.Sprintf("unique_ig:identifiers,id,%d", req.ID)},
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
