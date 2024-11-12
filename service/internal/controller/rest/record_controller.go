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

type RecordController struct {
	service *service.RecordService
}

func NewRecordController(service *service.RecordService) *RecordController {
	return &RecordController{service: service}
}

// TODO : test all the functions

func (c *RecordController) ListRecords(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	records, total, err := c.service.ListRecords(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, records, paginationMeta, "Records fetched successfully", http.StatusOK, nil, nil)
}
func (c *RecordController) CreateRecord(ctx *fiber.Ctx) error {
	var req dtos.CreateRecordRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewRecordStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	createdAt := time.Now()

	record := models.Record{
		EducationID: req.EducationID,
		UserID:      req.UserID,
		TimeSpent:   req.TimeSpent,
		CreatedAt:   &createdAt,
	}

	createdRecord, err := c.service.CreateRecord(ctx.Context(), &record)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create record", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetRecordParams{ID: createdRecord.ID}
	getRecord, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getRecord}, paginationMeta, "Record created successfully", http.StatusCreated, nil, nil)
}

func (c *RecordController) GetRecordByID(ctx *fiber.Ctx) error {
	var req dtos.GetRecordByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetRecordParams{ID: req.ID}
	record, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	recordArray := []interface{}{record}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, recordArray, paginationMeta, "Record fetched successfully", http.StatusOK, nil, nil)
}

// update record
func (c *RecordController) UpdateRecord(ctx *fiber.Ctx) error {
	var req dtos.UpdateRecordRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewRecordUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	params := &dtos.GetRecordParams{ID: req.ID}
	// Fetch the existing record to get the current data
	existingRecord, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	record := models.Record{
		ID:          req.ID,
		EducationID: existingRecord.EducationID,
		UserID:      req.UserID,
		TimeSpent:   req.TimeSpent,
		CreatedAt:   existingRecord.CreatedAt,
	}

	if req.EducationID != nil {
		record.EducationID = *req.EducationID
	}

	updatedRecord, err := c.service.UpdateRecord(ctx.Context(), &record)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update record", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetRecordParams{ID: updatedRecord.ID}
	getRecord, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getRecord}, paginationMeta, "Record updated successfully", http.StatusOK, nil, nil)
}

// delete record
func (c *RecordController) DeleteRecord(ctx *fiber.Ctx) error {
	var req dtos.DeleteRecordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetRecordParams{ID: req.ID}
	// GET record by ID
	_, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteRecord(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete record", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Record deleted successfully", http.StatusOK, nil, nil)
}

// restore record
func (c *RecordController) RestoreRecord(ctx *fiber.Ctx) error {
	var req dtos.DeleteRecordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, "ID is required", nil)
	}

	isDeleted := 1
	params := &dtos.GetRecordParams{ID: req.ID, IsDeleted: &isDeleted}
	// GET record by ID
	_, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreRecord(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore record", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Record restored successfully", http.StatusOK, nil, nil)
}

func (c *RecordController) ListRecordsByAuthUser(ctx *fiber.Ctx) error {
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

	records, total, err := c.service.ListRecords(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, records, paginationMeta, "Records fetched successfully", http.StatusOK, nil, nil)
}

// make auth create record
func (c *RecordController) CreateRecordByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.CreateRecordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))
	req.UserID = userID

	// Validate the request
	reqValidator := form_requests.NewRecordStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return utils.GetResponse(ctx, nil, nil, "Validation failed", http.StatusBadRequest, reqValidator, nil)
	}

	createdAt := time.Now()

	record := models.Record{
		EducationID: req.EducationID,
		UserID:      userID,
		TimeSpent:   req.TimeSpent,
		CreatedAt:   &createdAt,
	}

	createdRecord, err := c.service.CreateRecord(ctx.Context(), &record)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create record", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetRecordParams{ID: createdRecord.ID}
	getRecord, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getRecord}, paginationMeta, "Record created successfully", http.StatusCreated, nil, nil)
}

func (c *RecordController) GetRecordByIDByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.GetRecordByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetRecordParams{ID: req.ID}
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	record, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	recordArray := []interface{}{record}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, recordArray, paginationMeta, "Record fetched successfully", http.StatusOK, nil, nil)
}

// update record
func (c *RecordController) UpdateRecordByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.UpdateRecordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	// Validate the request
	reqValidator := form_requests.NewRecordUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return utils.GetResponse(ctx, nil, nil, "Validation failed", http.StatusBadRequest, reqValidator, nil)
	}

	params := &dtos.GetRecordParams{ID: req.ID}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	// Fetch the existing record to get the current data
	existingRecord, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	record := models.Record{
		ID:          req.ID,
		EducationID: existingRecord.EducationID,
		UserID:      userID,
		TimeSpent:   req.TimeSpent,
		CreatedAt:   existingRecord.CreatedAt,
	}

	if req.EducationID != nil {
		record.EducationID = *req.EducationID
	}

	updatedRecord, err := c.service.UpdateRecord(ctx.Context(), &record)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update record", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetRecordParams{ID: updatedRecord.ID}
	getRecord, err := c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getRecord}, paginationMeta, "Record updated successfully", http.StatusOK, nil, nil)
}

// delete record
func (c *RecordController) DeleteRecordByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.DeleteRecordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetRecordParams{ID: req.ID}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	// GET record by ID
	_, err = c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteRecord(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete record", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Record deleted successfully", http.StatusOK, nil, nil)
}

// restore record
func (c *RecordController) RestoreRecordByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.DeleteRecordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusBadRequest, "ID is required", nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	isDeleted := 1

	params := &dtos.GetRecordParams{ID: req.ID, IsDeleted: &isDeleted, UserID: userID}

	// GET record by ID
	_, err = c.service.GetRecordByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Record not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreRecord(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore record", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Record restored successfully", http.StatusOK, nil, nil)
}
