package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// SubscribeUpdateRequest handles the validation for the RegisterRequest.
type SubscribeUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of SubscribeUpdateRequest.
func NewSubscribeUpdateRequest() *SubscribeUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &SubscribeUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *SubscribeUpdateRequest) Validate(req *dtos.UpdateSubscribeRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"name":        []string{"required", "min:3", fmt.Sprintf("unique_ig:users,email,%d", req.ID)},
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
