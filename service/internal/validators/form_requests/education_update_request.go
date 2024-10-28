package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// EducationUpdateRequest handles the validation for the RegisterRequest.
type EducationUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of EducationUpdateRequest.
func NewEducationUpdateRequest() *EducationUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &EducationUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *EducationUpdateRequest) Validate(req *dtos.UpdateEducationRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"module_id":       []string{"required", "exists:modules,id"},
		"no_urut":         []string{},
		"name":            []string{"required", fmt.Sprintf("unique_ig:modules,name,%d", req.ID)},
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
