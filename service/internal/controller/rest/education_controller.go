package rest

import (
	"encoding/json"
	"log"
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

type EducationController struct {
	service *service.EducationService
}

func NewEducationController(service *service.EducationService) *EducationController {
	return &EducationController{service: service}
}

func (c *EducationController) GetEducations(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	educations, total, err := c.service.GetEducations(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, educations, paginationMeta, "Educations fetched successfully", http.StatusOK, nil, nil)
}
func (c *EducationController) CreateEducation(ctx *fiber.Ctx) error {
	var req dtos.CreateEducationRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewEducationStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	thumbnailURL := ""
	videoURL := ""

	attachmentUrls := []string{}
	// Convert attachmentUrls to JSON
	// attachmentUrlsJSON, err := json.Marshal(req.AttachmentUrls)
	attachmentUrlsJSON, err := json.Marshal(attachmentUrls)
	if err != nil {
		log.Println("Error converting attachment URLs to JSON:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// utils.DD(map[string]interface{}{
	// 	// "attachmentUrls":     req.AttachmentUrls,
	// 	"attachmentUrlsJSON": string(attachmentUrlsJSON),
	// })

	createdAt := time.Now()

	education := models.Education{
		ModuleID:      req.ModuleID,
		NoUrut:        req.NoUrut,
		Name:          req.Name,
		Description:   req.Description,
		TextMateri:    req.TextMateri,
		ThumbnailURL:  thumbnailURL,
		VideoURL:      videoURL,
		AttachmentURL: string(attachmentUrlsJSON),
		CreatedByID:   &userID,
		CreatedAt:     &createdAt,
	}

	createdEducation, err := c.service.CreateEducation(ctx.Context(), &education)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create education", http.StatusInternalServerError, err.Error(), nil)
	}

	getEducation, err := c.service.GetEducationByID(ctx.Context(), createdEducation.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getEducation}, paginationMeta, "Education created successfully", http.StatusCreated, nil, nil)
}

func (c *EducationController) GetEducationByID(ctx *fiber.Ctx) error {
	var req dtos.GetEducationByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusBadRequest, "ID is required", nil)
	}

	education, err := c.service.GetEducationByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusNotFound, err.Error(), nil)
	}

	educationArray := []interface{}{education}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, educationArray, paginationMeta, "Education fetched successfully", http.StatusOK, nil, nil)
}

// TODO : Implement the UpdateEducation method
// update education
func (c *EducationController) UpdateEducation(ctx *fiber.Ctx) error {
	var req dtos.UpdateEducationRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewEducationUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Fetch the existing education to get the current data
	existingEducation, err := c.service.GetEducationByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusNotFound, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error(), "message": "Unauthorized", "status": fiber.StatusUnauthorized})
	}
	userID := uint(claims["user_id"].(float64))

	thumbnailURL := ""
	videoURL := ""

	attachmentUrls := []string{}
	// Convert attachmentUrls to JSON
	// attachmentUrlsJSON, err := json.Marshal(req.AttachmentUrls)
	attachmentUrlsJSON, err := json.Marshal(attachmentUrls)
	if err != nil {
		log.Println("Error converting attachment URLs to JSON:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	education := models.Education{
		ID:            req.ID,
		ModuleID:      req.ModuleID,
		NoUrut:        req.NoUrut,
		Name:          req.Name,
		Description:   req.Description,
		TextMateri:    req.TextMateri,
		ThumbnailURL:  thumbnailURL,
		VideoURL:      videoURL,
		AttachmentURL: string(attachmentUrlsJSON),
		CreatedByID:   &existingEducation.CreatedByID,
		UpdatedByID:   &userID,
		CreatedAt:     existingEducation.CreatedAt,
	}

	updatedEducation, err := c.service.UpdateEducation(ctx.Context(), &education)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update education", http.StatusInternalServerError, err.Error(), nil)
	}

	getEducation, err := c.service.GetEducationByID(ctx.Context(), updatedEducation.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getEducation}, paginationMeta, "Education updated successfully", http.StatusOK, nil, nil)
}

// delete education
func (c *EducationController) DeleteEducation(ctx *fiber.Ctx) error {
	var req dtos.DeleteEducationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET education by ID
	_, err := c.service.GetEducationByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteEducation(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete education", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Education deleted successfully", http.StatusOK, nil, nil)
}

// restore education
func (c *EducationController) RestoreEducation(ctx *fiber.Ctx) error {
	var req dtos.DeleteEducationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusBadRequest, "ID is required", nil)
	}

	// GET education by ID
	_, err := c.service.GetEducationByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Education not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreEducation(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore education", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Education restored successfully", http.StatusOK, nil, nil)
}
