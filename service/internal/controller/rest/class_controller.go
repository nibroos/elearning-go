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

type ClassController struct {
	service *service.ClassService
}

func NewClassController(service *service.ClassService) *ClassController {
	return &ClassController{service: service}
}

func (c *ClassController) GetClasss(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	classes, total, err := c.service.GetClasses(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, classes, paginationMeta, "Classs fetched successfully", http.StatusOK, nil, nil)
}
func (c *ClassController) CreateClass(ctx *fiber.Ctx) error {
	var req dtos.CreateClassRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewClassStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error(), "message": "Unauthorized", "status": fiber.StatusUnauthorized})
	}
	userID := uint(claims["user_id"].(float64))

	bannerURL := ""
	logoURL := ""
	videoURL := ""

	class := models.Class{
		Name:        req.Name,
		Description: req.Description,
		BannerURL:   &bannerURL,
		LogoURL:     &logoURL,
		VideoURL:    &videoURL,
		SubscribeID: req.SubscribeID,
		InchargeID:  req.InchargeID,
		CreatedByID: &userID,
	}

	createdClass, err := c.service.CreateClass(ctx.Context(), &class)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create class", http.StatusInternalServerError, err.Error(), nil)
	}

	getClass, err := c.service.GetClassByID(ctx.Context(), createdClass.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getClass}, paginationMeta, "Class created successfully", http.StatusCreated, nil, nil)
}
func (c *ClassController) GetClassByID(ctx *fiber.Ctx) error {
	var req dtos.GetClassByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusBadRequest, "ID is required", nil)
	}

	class, err := c.service.GetClassByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusNotFound, err.Error(), nil)
	}

	classArray := []interface{}{class}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, classArray, paginationMeta, "Class fetched successfully", http.StatusOK, nil, nil)
}

// update class
func (c *ClassController) UpdateClass(ctx *fiber.Ctx) error {
	var req dtos.UpdateClassRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewClassUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Fetch the existing subscribe to get the current data
	existingClass, err := c.service.GetClassByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusNotFound, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error(), "message": "Unauthorized", "status": fiber.StatusUnauthorized})
	}
	userID := uint(claims["user_id"].(float64))

	bannerURL := ""
	logoURL := ""
	videoURL := ""

	class := models.Class{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		BannerURL:   &bannerURL,
		LogoURL:     &logoURL,
		VideoURL:    &videoURL,
		SubscribeID: req.SubscribeID,
		InchargeID:  req.InchargeID,
		CreatedByID: &existingClass.CreatedByID,
		UpdatedByID: &userID,
		CreatedAt:   *existingClass.CreatedAt,
	}

	updatedClass, err := c.service.UpdateClass(ctx.Context(), &class)
	if err != nil {
		if err.Error() == "classname already exists" {
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{"errors": err.Error(), "message": "Classname already exists", "status": http.StatusConflict})
		}
		return utils.GetResponse(ctx, nil, nil, "Failed to update class", http.StatusInternalServerError, err.Error(), nil)
	}

	getClass, err := c.service.GetClassByID(ctx.Context(), updatedClass.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getClass}, paginationMeta, "Class updated successfully", http.StatusOK, nil, nil)
}

// delete class
func (c *ClassController) DeleteClass(ctx *fiber.Ctx) error {
	var req dtos.DeleteClassRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET class by ID
	_, err := c.service.GetClassByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteClass(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete class", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Class deleted successfully", http.StatusOK, nil, nil)
}

// restore class
func (c *ClassController) RestoreClass(ctx *fiber.Ctx) error {
	var req dtos.DeleteClassRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET class by ID
	_, err := c.service.GetClassByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Class not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreClass(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore class", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Class restored successfully", http.StatusOK, nil, nil)
}
