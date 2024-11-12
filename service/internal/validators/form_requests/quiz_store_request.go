package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// QuizStoreRequest handles the validation for the RegisterRequest.
type QuizStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of QuizStoreRequest.
func NewQuizStoreRequest() *QuizStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &QuizStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *QuizStoreRequest) Validate(req *dtos.CreateQuizRequest, ctx context.Context) map[string]string {
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
