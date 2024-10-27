package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/middleware"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"github.com/nibroos/elearning-go/service/internal/validators/form_requests"
)

type SubscribeController struct {
	service *service.SubscribeService
}

func NewSubscribeController(service *service.SubscribeService) *SubscribeController {
	return &SubscribeController{service: service}
}

func (c *SubscribeController) GetSubscribes(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	subscribes, total, err := c.service.GetSubscribes(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, subscribes, paginationMeta, "Subscribes fetched successfully", http.StatusOK, nil, nil)
}
func (c *SubscribeController) CreateSubscribe(ctx *fiber.Ctx) error {
	var req dtos.CreateSubscribeRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewSubscribeStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	subscribe := models.Subscribe{
		Name:        req.Name,
		Description: req.Description,
		CreatedByID: &userID,
	}

	createdSubscribe, err := c.service.CreateSubscribe(ctx.Context(), &subscribe)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create subscribe", http.StatusInternalServerError, err.Error(), nil)
	}

	getSubscribe, err := c.service.GetSubscribeByID(ctx.Context(), createdSubscribe.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getSubscribe}, paginationMeta, "Subscribe created successfully", http.StatusCreated, nil, nil)
}
func (c *SubscribeController) GetSubscribeByID(ctx *fiber.Ctx) error {
	var req dtos.GetSubscribeByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusBadRequest, "ID is required", nil)
	}

	subscribe, err := c.service.GetSubscribeByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusNotFound, err.Error(), nil)
	}

	subscribeArray := []interface{}{subscribe}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, subscribeArray, paginationMeta, "Subscribe fetched successfully", http.StatusOK, nil, nil)
}

// update subscribe
func (c *SubscribeController) UpdateSubscribe(ctx *fiber.Ctx) error {
	var req dtos.UpdateSubscribeRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewSubscribeUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Fetch the existing subscribe to get the current data
	existingSubscribe, err := c.service.GetSubscribeByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusNotFound, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error(), "message": "Unauthorized", "status": fiber.StatusUnauthorized})
	}
	userID := uint(claims["user_id"].(float64))

	subscribe := models.Subscribe{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		CreatedByID: &existingSubscribe.CreatedByID,
		UpdatedByID: &userID,
	}

	updatedSubscribe, err := c.service.UpdateSubscribe(ctx.Context(), &subscribe)
	if err != nil {
		if err.Error() == "subscribe name already exists" {
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{"errors": err.Error(), "message": "Subscribe already exists", "status": http.StatusConflict})
		}
		return utils.GetResponse(ctx, nil, nil, "Failed to update subscribe", http.StatusInternalServerError, err.Error(), nil)
	}

	getSubscribe, err := c.service.GetSubscribeByID(ctx.Context(), updatedSubscribe.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getSubscribe}, paginationMeta, "Subscribe updated successfully", http.StatusOK, nil, nil)
}

// delete subscribe
func (c *SubscribeController) DeleteSubscribe(ctx *fiber.Ctx) error {
	var req dtos.DeleteSubscribeRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET subscribe by ID
	_, err := c.service.GetSubscribeByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteSubscribe(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete subscribe", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Subscribe deleted successfully", http.StatusOK, nil, nil)
}

// restore subscribe
func (c *SubscribeController) RestoreSubscribe(ctx *fiber.Ctx) error {
	var req dtos.DeleteSubscribeRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET subscribe by ID
	_, err := c.service.GetSubscribeByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Subscribe not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreSubscribe(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore subscribe", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Subscribe restored successfully", http.StatusOK, nil, nil)
}
