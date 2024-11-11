package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// RecordStoreRequest handles the validation for the RegisterRequest.
type RecordStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of RecordStoreRequest.
func NewRecordStoreRequest() *RecordStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &RecordStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *RecordStoreRequest) Validate(req *dtos.CreateRecordRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"education_id": []string{"required", "exists:educations,id"},
		"user_id":      []string{"required", "exists:users,id"},
		"time_spent":   []string{"required"},
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
