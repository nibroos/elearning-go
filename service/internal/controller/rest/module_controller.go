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

type ModuleController struct {
	service *service.ModuleService
}

func NewModuleController(service *service.ModuleService) *ModuleController {
	return &ModuleController{service: service}
}

func (c *ModuleController) GetModules(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	modules, total, err := c.service.GetModules(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, modules, paginationMeta, "Modules fetched successfully", http.StatusOK, nil, nil)
}
func (c *ModuleController) CreateModule(ctx *fiber.Ctx) error {
	var req dtos.CreateModuleRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewModuleStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	module := models.Module{
		Name:        req.Name,
		ClassID:     req.ClassID,
		Description: req.Description,
		CreatedByID: &userID,
	}

	createdModule, err := c.service.CreateModule(ctx.Context(), &module)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create module", http.StatusInternalServerError, err.Error(), nil)
	}

	getModule, err := c.service.GetModuleByID(ctx.Context(), createdModule.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getModule}, paginationMeta, "Module created successfully", http.StatusCreated, nil, nil)
}
func (c *ModuleController) GetModuleByID(ctx *fiber.Ctx) error {
	var req dtos.GetModuleByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusBadRequest, "ID is required", nil)
	}

	module, err := c.service.GetModuleByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusNotFound, err.Error(), nil)
	}

	moduleArray := []interface{}{module}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, moduleArray, paginationMeta, "Module fetched successfully", http.StatusOK, nil, nil)
}

// update module
func (c *ModuleController) UpdateModule(ctx *fiber.Ctx) error {
	var req dtos.UpdateModuleRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewModuleUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Fetch the existing module to get the current data
	existingModule, err := c.service.GetModuleByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusNotFound, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error(), "message": "Unauthorized", "status": fiber.StatusUnauthorized})
	}
	userID := uint(claims["user_id"].(float64))

	module := models.Module{
		ID:          req.ID,
		ClassID:     req.ClassID,
		Name:        req.Name,
		Description: req.Description,
		CreatedByID: &existingModule.CreatedByID,
		UpdatedByID: &userID,
	}

	updatedModule, err := c.service.UpdateModule(ctx.Context(), &module)
	if err != nil {
		if err.Error() == "module name already exists" {
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{"errors": err.Error(), "message": "Module already exists", "status": http.StatusConflict})
		}
		return utils.GetResponse(ctx, nil, nil, "Failed to update module", http.StatusInternalServerError, err.Error(), nil)
	}

	getModule, err := c.service.GetModuleByID(ctx.Context(), updatedModule.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getModule}, paginationMeta, "Module updated successfully", http.StatusOK, nil, nil)
}

// delete module
func (c *ModuleController) DeleteModule(ctx *fiber.Ctx) error {
	var req dtos.DeleteModuleRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET module by ID
	_, err := c.service.GetModuleByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteModule(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete module", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Module deleted successfully", http.StatusOK, nil, nil)
}

// restore module
func (c *ModuleController) RestoreModule(ctx *fiber.Ctx) error {
	var req dtos.DeleteModuleRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET module by ID
	_, err := c.service.GetModuleByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Module not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreModule(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore module", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Module restored successfully", http.StatusOK, nil, nil)
}
