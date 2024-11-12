package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// QuizUpdateRequest handles the validation for the RegisterRequest.
type QuizUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of QuizUpdateRequest.
func NewQuizUpdateRequest() *QuizUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &QuizUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *QuizUpdateRequest) Validate(req *dtos.UpdateQuizRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"name":        []string{"required"},
		"description": []string{},
		"threshold":   []string{"required"},
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
