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

type IdentifierController struct {
	service *service.IdentifierService
}

func NewIdentifierController(service *service.IdentifierService) *IdentifierController {
	return &IdentifierController{service: service}
}

func (c *IdentifierController) ListIdentifiers(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	identifiers, total, err := c.service.ListIdentifiers(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, identifiers, paginationMeta, "Identifiers fetched successfully", http.StatusOK, nil, nil)
}
func (c *IdentifierController) CreateIdentifier(ctx *fiber.Ctx) error {
	var req dtos.CreateIdentifierRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewIdentifierStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	createdAt := time.Now()

	identifier := models.Identifier{
		TypeIdentifierID: req.TypeIdentifierID,
		UserID:           req.UserID,
		RefNum:           req.RefNum,
		Status:           req.Status,
		CreatedAt:        &createdAt,
		OptionsJSON:      nil,
	}

	createdIdentifier, err := c.service.CreateIdentifier(ctx.Context(), &identifier)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetIdentifierParams{ID: createdIdentifier.ID}
	getIdentifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getIdentifier}, paginationMeta, "Identifier created successfully", http.StatusCreated, nil, nil)
}

func (c *IdentifierController) GetIdentifierByID(ctx *fiber.Ctx) error {
	var req dtos.GetIdentifierByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetIdentifierParams{ID: req.ID}
	identifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	identifierArray := []interface{}{identifier}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, identifierArray, paginationMeta, "Identifier fetched successfully", http.StatusOK, nil, nil)
}

// update identifier
func (c *IdentifierController) UpdateIdentifier(ctx *fiber.Ctx) error {
	var req dtos.UpdateIdentifierRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewIdentifierUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	params := &dtos.GetIdentifierParams{ID: req.ID}
	// Fetch the existing identifier to get the current data
	existingIdentifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	identifier := models.Identifier{
		ID:               req.ID,
		TypeIdentifierID: existingIdentifier.TypeIdentifierID,
		UserID:           req.UserID,
		RefNum:           req.RefNum,
		Status:           req.Status,
		CreatedAt:        existingIdentifier.CreatedAt,
	}

	if req.TypeIdentifierID != nil {
		identifier.TypeIdentifierID = *req.TypeIdentifierID
	}

	updatedIdentifier, err := c.service.UpdateIdentifier(ctx.Context(), &identifier)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetIdentifierParams{ID: updatedIdentifier.ID}
	getIdentifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getIdentifier}, paginationMeta, "Identifier updated successfully", http.StatusOK, nil, nil)
}

// delete identifier
func (c *IdentifierController) DeleteIdentifier(ctx *fiber.Ctx) error {
	var req dtos.DeleteIdentifierRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetIdentifierParams{ID: req.ID}
	// GET identifier by ID
	_, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteIdentifier(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Identifier deleted successfully", http.StatusOK, nil, nil)
}

// restore identifier
func (c *IdentifierController) RestoreIdentifier(ctx *fiber.Ctx) error {
	var req dtos.DeleteIdentifierRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, "ID is required", nil)
	}

	isDeleted := 1
	params := &dtos.GetIdentifierParams{ID: req.ID, IsDeleted: &isDeleted}
	// GET identifier by ID
	_, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreIdentifier(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Identifier restored successfully", http.StatusOK, nil, nil)
}

func (c *IdentifierController) ListIdentifiersByAuthUser(ctx *fiber.Ctx) error {
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

	identifiers, total, err := c.service.ListIdentifiersByAuthUser(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, identifiers, paginationMeta, "Identifiers fetched successfully", http.StatusOK, nil, nil)
}

func (c *IdentifierController) GetIdentifierByAuthUser(ctx *fiber.Ctx) error {
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	var req dtos.GetIdentifierByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetIdentifierParams{ID: req.ID, UserID: userID}
	identifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	identifierArray := []interface{}{identifier}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, identifierArray, paginationMeta, "Identifier fetched successfully", http.StatusOK, nil, nil)
}

func (c *IdentifierController) CreateIdentifierByAuthUser(ctx *fiber.Ctx) error {
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	var req dtos.CreateIdentifierRequest

	req.UserID = userID

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewIdentifierStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	createdAt := time.Now()

	identifier := models.Identifier{
		TypeIdentifierID: req.TypeIdentifierID,
		UserID:           userID,
		RefNum:           req.RefNum,
		Status:           req.Status,
		CreatedAt:        &createdAt,
		OptionsJSON:      nil,
	}

	createdIdentifier, err := c.service.CreateIdentifier(ctx.Context(), &identifier)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetIdentifierParams{ID: createdIdentifier.ID, UserID: userID}
	getIdentifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getIdentifier}, paginationMeta, "Identifier created successfully", http.StatusCreated, nil, nil)
}

func (c *IdentifierController) UpdateIdentifierByAuthUser(ctx *fiber.Ctx) error {
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	var req dtos.UpdateIdentifierRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewIdentifierUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	params := &dtos.GetIdentifierParams{ID: req.ID, UserID: userID}
	// Fetch the existing identifier to get the current data
	existingIdentifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	identifier := models.Identifier{
		ID:               req.ID,
		TypeIdentifierID: existingIdentifier.TypeIdentifierID,
		UserID:           userID,
		RefNum:           req.RefNum,
		Status:           req.Status,
		CreatedAt:        existingIdentifier.CreatedAt,
	}

	if req.TypeIdentifierID != nil {
		identifier.TypeIdentifierID = *req.TypeIdentifierID
	}

	updatedIdentifier, err := c.service.UpdateIdentifier(ctx.Context(), &identifier)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetIdentifierParams{ID: updatedIdentifier.ID, UserID: userID}
	getIdentifier, err := c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getIdentifier}, paginationMeta, "Identifier updated successfully", http.StatusOK, nil, nil)
}

func (c *IdentifierController) DeleteIdentifierByAuthUser(ctx *fiber.Ctx) error {
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	var req dtos.DeleteIdentifierRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetIdentifierParams{ID: req.ID, UserID: userID}
	// GET identifier by ID
	_, err = c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteIdentifier(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Identifier deleted successfully", http.StatusOK, nil, nil)
}

func (c *IdentifierController) RestoreIdentifierByAuthUser(ctx *fiber.Ctx) error {
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	var req dtos.DeleteIdentifierRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusBadRequest, "ID is required", nil)
	}

	isDeleted := 1
	params := &dtos.GetIdentifierParams{ID: req.ID, UserID: userID, IsDeleted: &isDeleted}
	// GET identifier by ID
	_, err = c.service.GetIdentifierByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Identifier not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreIdentifier(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore identifier", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Identifier restored successfully", http.StatusOK, nil, nil)
}
