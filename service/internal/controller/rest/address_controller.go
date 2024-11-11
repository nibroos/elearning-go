package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/middleware"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"github.com/nibroos/elearning-go/service/internal/validators/form_requests"
)

type AddressController struct {
	service *service.AddressService
}

func NewAddressController(service *service.AddressService) *AddressController {
	return &AddressController{service: service}
}

func (c *AddressController) ListAddresses(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	addresses, total, err := c.service.ListAddresses(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, addresses, paginationMeta, "Addresses fetched successfully", http.StatusOK, nil, nil)
}
func (c *AddressController) CreateAddress(ctx *fiber.Ctx) error {
	var req dtos.CreateAddressRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewAddressStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	createdAt := time.Now()

	address := models.Address{
		TypeAddressID: req.TypeAddressID,
		UserID:        req.UserID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     &createdAt,
		OptionsJSON:   nil,
	}

	createdAddress, err := c.service.CreateAddress(ctx.Context(), &address)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create address", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetAddressParams{ID: createdAddress.ID}
	getAddress, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getAddress}, paginationMeta, "Address created successfully", http.StatusCreated, nil, nil)
}

func (c *AddressController) GetAddressByID(ctx *fiber.Ctx) error {
	var req dtos.GetAddressByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetAddressParams{ID: req.ID}
	address, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	addressArray := []interface{}{address}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, addressArray, paginationMeta, "Address fetched successfully", http.StatusOK, nil, nil)
}

// update address
func (c *AddressController) UpdateAddress(ctx *fiber.Ctx) error {
	var req dtos.UpdateAddressRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewAddressUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	params := &dtos.GetAddressParams{ID: req.ID}
	// Fetch the existing address to get the current data
	existingAddress, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	address := models.Address{
		ID:            req.ID,
		TypeAddressID: existingAddress.TypeAddressID,
		UserID:        req.UserID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     existingAddress.CreatedAt,
	}

	if req.TypeAddressID != nil {
		address.TypeAddressID = *req.TypeAddressID
	}

	updatedAddress, err := c.service.UpdateAddress(ctx.Context(), &address)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update address", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetAddressParams{ID: updatedAddress.ID}
	getAddress, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getAddress}, paginationMeta, "Address updated successfully", http.StatusOK, nil, nil)
}

// delete address
func (c *AddressController) DeleteAddress(ctx *fiber.Ctx) error {
	var req dtos.DeleteAddressRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetAddressParams{ID: req.ID}
	// GET address by ID
	_, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteAddress(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete address", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Address deleted successfully", http.StatusOK, nil, nil)
}

// restore address
func (c *AddressController) RestoreAddress(ctx *fiber.Ctx) error {
	var req dtos.DeleteAddressRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, "ID is required", nil)
	}

	isDeleted := 1
	params := &dtos.GetAddressParams{ID: req.ID, IsDeleted: &isDeleted}
	// GET address by ID
	_, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreAddress(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore address", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Address restored successfully", http.StatusOK, nil, nil)
}

func (c *AddressController) ListAddressesByAuthUser(ctx *fiber.Ctx) error {
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	filters, ok := ctx.Locals("filters").(map[string]string)
	filters["user_id"] = fmt.Sprint(userID)

	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	addresses, total, err := c.service.ListAddresses(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, addresses, paginationMeta, "Addresses fetched successfully", http.StatusOK, nil, nil)
}

// make auth create address
func (c *AddressController) CreateAddressByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.CreateAddressRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))
	req.UserID = userID

	// Validate the request
	reqValidator := form_requests.NewAddressStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return utils.GetResponse(ctx, nil, nil, "Validation failed", http.StatusBadRequest, reqValidator, nil)
	}

	createdAt := time.Now()

	address := models.Address{
		TypeAddressID: req.TypeAddressID,
		UserID:        userID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     &createdAt,
		OptionsJSON:   nil,
	}

	createdAddress, err := c.service.CreateAddress(ctx.Context(), &address)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create address", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetAddressParams{ID: createdAddress.ID}
	getAddress, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getAddress}, paginationMeta, "Address created successfully", http.StatusCreated, nil, nil)
}

func (c *AddressController) GetAddressByIDByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.GetAddressByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetAddressParams{ID: req.ID}
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	address, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	addressArray := []interface{}{address}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, addressArray, paginationMeta, "Address fetched successfully", http.StatusOK, nil, nil)
}

// update address
func (c *AddressController) UpdateAddressByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.UpdateAddressRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	// Validate the request
	reqValidator := form_requests.NewAddressUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return utils.GetResponse(ctx, nil, nil, "Validation failed", http.StatusBadRequest, reqValidator, nil)
	}

	params := &dtos.GetAddressParams{ID: req.ID}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	// Fetch the existing address to get the current data
	existingAddress, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	address := models.Address{
		ID:            req.ID,
		TypeAddressID: existingAddress.TypeAddressID,
		UserID:        userID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     existingAddress.CreatedAt,
	}

	if req.TypeAddressID != nil {
		address.TypeAddressID = *req.TypeAddressID
	}

	updatedAddress, err := c.service.UpdateAddress(ctx.Context(), &address)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update address", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetAddressParams{ID: updatedAddress.ID}
	getAddress, err := c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getAddress}, paginationMeta, "Address updated successfully", http.StatusOK, nil, nil)
}

// delete address
func (c *AddressController) DeleteAddressByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.DeleteAddressRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetAddressParams{ID: req.ID}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	// GET address by ID
	_, err = c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteAddress(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete address", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Address deleted successfully", http.StatusOK, nil, nil)
}

// restore address
func (c *AddressController) RestoreAddressByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.DeleteAddressRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusBadRequest, "ID is required", nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	isDeleted := 1

	params := &dtos.GetAddressParams{ID: req.ID, IsDeleted: &isDeleted, UserID: userID}

	// GET address by ID
	_, err = c.service.GetAddressByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Address not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreAddress(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore address", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Address restored successfully", http.StatusOK, nil, nil)
}
