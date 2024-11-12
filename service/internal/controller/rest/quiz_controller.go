package rest

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/service"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"github.com/nibroos/elearning-go/service/internal/validators/form_requests"
)

type QuizController struct {
	service *service.QuizService
}

func NewQuizController(service *service.QuizService) *QuizController {
	return &QuizController{service: service}
}

// TODO : test all the functions

func (c *QuizController) ListQuizes(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	contacts, total, err := c.service.ListQuizes(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, contacts, paginationMeta, "Quizes fetched successfully", http.StatusOK, nil, nil)
}
func (c *QuizController) CreateQuiz(ctx *fiber.Ctx) error {
	var req dtos.CreateQuizRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewQuizStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	createdAt := time.Now()

	contact := models.Quiz{
		Name:        req.Name,
		Description: req.Description,
		Threshold:   req.Threshold,
		CreatedAt:   &createdAt,
	}

	createdQuiz, err := c.service.CreateQuiz(ctx.Context(), &contact)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create contact", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetQuizParams{ID: createdQuiz.ID}
	getQuiz, err := c.service.GetQuizByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getQuiz}, paginationMeta, "Quiz created successfully", http.StatusCreated, nil, nil)
}

func (c *QuizController) GetQuizByID(ctx *fiber.Ctx) error {
	var req dtos.GetQuizByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetQuizParams{ID: req.ID}
	contact, err := c.service.GetQuizByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusNotFound, err.Error(), nil)
	}

	contactArray := []interface{}{contact}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, contactArray, paginationMeta, "Quiz fetched successfully", http.StatusOK, nil, nil)
}

// update contact
func (c *QuizController) UpdateQuiz(ctx *fiber.Ctx) error {
	var req dtos.UpdateQuizRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewQuizUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	params := &dtos.GetQuizParams{ID: req.ID}
	// Fetch the existing contact to get the current data
	existingQuiz, err := c.service.GetQuizByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusNotFound, err.Error(), nil)
	}

	contact := models.Quiz{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Threshold:   req.Threshold,
		CreatedAt:   existingQuiz.CreatedAt,
	}

	updatedQuiz, err := c.service.UpdateQuiz(ctx.Context(), &contact)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update contact", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetQuizParams{ID: updatedQuiz.ID}
	getQuiz, err := c.service.GetQuizByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getQuiz}, paginationMeta, "Quiz updated successfully", http.StatusOK, nil, nil)
}

// delete contact
func (c *QuizController) DeleteQuiz(ctx *fiber.Ctx) error {
	var req dtos.DeleteQuizRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetQuizParams{ID: req.ID}
	// GET contact by ID
	_, err := c.service.GetQuizByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteQuiz(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete contact", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Quiz deleted successfully", http.StatusOK, nil, nil)
}

// restore contact
func (c *QuizController) RestoreQuiz(ctx *fiber.Ctx) error {
	var req dtos.DeleteQuizRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusBadRequest, "ID is required", nil)
	}

	isDeleted := 1
	params := &dtos.GetQuizParams{ID: req.ID, IsDeleted: &isDeleted}
	// GET contact by ID
	_, err := c.service.GetQuizByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Quiz not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreQuiz(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore contact", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Quiz restored successfully", http.StatusOK, nil, nil)
}

// TODO : check accessible quiz on student
