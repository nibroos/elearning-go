package form_requests

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// AddressStoreRequest handles the validation for the RegisterRequest.
type AddressStoreRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterStoreRequest creates a new instance of AddressStoreRequest.
func NewAddressStoreRequest() *AddressStoreRequest {
	v := govalidator.New(govalidator.Options{})
	return &AddressStoreRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *AddressStoreRequest) Validate(req *dtos.CreateAddressRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"type_address_id": []string{"required", "exists:mix_values,id"},
		"user_id":         []string{"required", "exists:users,id"},
		"ref_num":         []string{"required"},
		"status":          []string{"required"},
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
