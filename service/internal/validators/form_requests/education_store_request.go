package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// EducationStoreRequest handles the validation for the RegisterRequest.
type EducationStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of EducationStoreRequest.
func NewEducationStoreRequest() *EducationStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &EducationStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *EducationStoreRequest) Validate(req *dtos.CreateEducationRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"module_id":       []string{"required", "exists:modules,id"},
		"no_urut":         []string{},
		"name":            []string{"required", "unique:educations,name"},
		"description":     []string{"required"},
		"text_materi":     []string{},
		"attachment_urls": []string{"array"},
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
