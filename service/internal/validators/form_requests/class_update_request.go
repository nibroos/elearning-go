package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// ClassUpdateRequest handles the validation for the RegisterRequest.
type ClassUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of ClassUpdateRequest.
func NewClassUpdateRequest() *ClassUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &ClassUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *ClassUpdateRequest) Validate(req *dtos.UpdateClassRequest, ctx context.Context) map[string]string {
	// utils.DD(req)
	rules := govalidator.MapData{
		"name":         []string{"required", fmt.Sprintf("unique_ig:classes,name,%d", req.ID)},
		"description":  []string{},
		"subscribe_id": []string{"required", "numeric", "exists:subscribes,id"},
		"incharge_id":  []string{"required", "numeric", "exists:users,id"},
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
