package form_requests

import (
	"context"
	"fmt"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/thedevsaddam/govalidator"
)

// AddressUpdateRequest handles the validation for the RegisterRequest.
type AddressUpdateRequest struct {
	Validator *govalidator.Validator
}

// NewRegisterUpdateRequest creates a new instance of AddressUpdateRequest.
func NewAddressUpdateRequest() *AddressUpdateRequest {
	v := govalidator.New(govalidator.Options{})
	return &AddressUpdateRequest{Validator: v}
}

// Validate validates the RegisterRequest.
func (r *AddressUpdateRequest) Validate(req *dtos.UpdateAddressRequest, ctx context.Context) map[string]string {
	rules := govalidator.MapData{
		"type_address_id": []string{"exists:mix_values,id"},
		"user_id":         []string{"required", "exists:users,id"},
		"ref_num":         []string{"required", fmt.Sprintf("unique_ig:addresses,id,%d", req.ID)},
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
