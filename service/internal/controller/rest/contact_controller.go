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

type ContactController struct {
	service *service.ContactService
}

func NewContactController(service *service.ContactService) *ContactController {
	return &ContactController{service: service}
}

// TODO : Fix all method below

func (c *ContactController) ListContacts(ctx *fiber.Ctx) error {
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	contacts, total, err := c.service.ListContacts(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, contacts, paginationMeta, "Contacts fetched successfully", http.StatusOK, nil, nil)
}
func (c *ContactController) CreateContact(ctx *fiber.Ctx) error {
	var req dtos.CreateContactRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewContactStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	createdAt := time.Now()

	contact := models.Contact{
		TypeContactID: req.TypeContactID,
		UserID:        &req.UserID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     &createdAt,
		OptionsJSON:   nil,
	}

	createdContact, err := c.service.CreateContact(ctx.Context(), &contact)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create contact", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetContactParams{ID: createdContact.ID}
	getContact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getContact}, paginationMeta, "Contact created successfully", http.StatusCreated, nil, nil)
}

func (c *ContactController) GetContactByID(ctx *fiber.Ctx) error {
	var req dtos.GetContactByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetContactParams{ID: req.ID}
	contact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	contactArray := []interface{}{contact}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, contactArray, paginationMeta, "Contact fetched successfully", http.StatusOK, nil, nil)
}

// update contact
func (c *ContactController) UpdateContact(ctx *fiber.Ctx) error {
	var req dtos.UpdateContactRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	reqValidator := form_requests.NewContactUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": reqValidator, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	params := &dtos.GetContactParams{ID: req.ID}
	// Fetch the existing contact to get the current data
	existingContact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	contact := models.Contact{
		ID:            req.ID,
		TypeContactID: existingContact.TypeContactID,
		UserID:        &req.UserID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     existingContact.CreatedAt,
	}

	if req.TypeContactID != nil {
		contact.TypeContactID = *req.TypeContactID
	}

	updatedContact, err := c.service.UpdateContact(ctx.Context(), &contact)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update contact", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetContactParams{ID: updatedContact.ID}
	getContact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getContact}, paginationMeta, "Contact updated successfully", http.StatusOK, nil, nil)
}

// delete contact
func (c *ContactController) DeleteContact(ctx *fiber.Ctx) error {
	var req dtos.DeleteContactRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetContactParams{ID: req.ID}
	// GET contact by ID
	_, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteContact(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete contact", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Contact deleted successfully", http.StatusOK, nil, nil)
}

// restore contact
func (c *ContactController) RestoreContact(ctx *fiber.Ctx) error {
	var req dtos.DeleteContactRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, "ID is required", nil)
	}

	isDeleted := 1
	params := &dtos.GetContactParams{ID: req.ID, IsDeleted: &isDeleted}
	// GET contact by ID
	_, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreContact(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore contact", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Contact restored successfully", http.StatusOK, nil, nil)
}

func (c *ContactController) ListContactsByAuthUser(ctx *fiber.Ctx) error {
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

	contacts, total, err := c.service.ListContactsByAuthUser(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)

	return utils.GetResponse(ctx, contacts, paginationMeta, "Contacts fetched successfully", http.StatusOK, nil, nil)
}

// make auth create contact
func (c *ContactController) CreateContactByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.CreateContactRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))
	req.UserID = userID

	// Validate the request
	reqValidator := form_requests.NewContactStoreRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return utils.GetResponse(ctx, nil, nil, "Validation failed", http.StatusBadRequest, reqValidator, nil)
	}

	createdAt := time.Now()

	contact := models.Contact{
		TypeContactID: req.TypeContactID,
		UserID:        &userID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     &createdAt,
		OptionsJSON:   nil,
	}

	createdContact, err := c.service.CreateContact(ctx.Context(), &contact)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to create contact", http.StatusInternalServerError, err.Error(), nil)
	}

	params := &dtos.GetContactParams{ID: createdContact.ID}
	getContact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getContact}, paginationMeta, "Contact created successfully", http.StatusCreated, nil, nil)
}

func (c *ContactController) GetContactByIDByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.GetContactByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetContactParams{ID: req.ID}
	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	contact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	contactArray := []interface{}{contact}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, contactArray, paginationMeta, "Contact fetched successfully", http.StatusOK, nil, nil)
}

// update contact
func (c *ContactController) UpdateContactByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.UpdateContactRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	// Validate the request
	reqValidator := form_requests.NewContactUpdateRequest().Validate(&req, ctx.Context())
	if reqValidator != nil {
		return utils.GetResponse(ctx, nil, nil, "Validation failed", http.StatusBadRequest, reqValidator, nil)
	}

	params := &dtos.GetContactParams{ID: req.ID}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	// Fetch the existing contact to get the current data
	existingContact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	contact := models.Contact{
		ID:            req.ID,
		TypeContactID: existingContact.TypeContactID,
		UserID:        &userID,
		RefNum:        req.RefNum,
		Status:        req.Status,
		CreatedAt:     existingContact.CreatedAt,
	}

	if req.TypeContactID != nil {
		contact.TypeContactID = *req.TypeContactID
	}

	updatedContact, err := c.service.UpdateContact(ctx.Context(), &contact)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to update contact", http.StatusInternalServerError, err.Error(), nil)
	}

	params = &dtos.GetContactParams{ID: updatedContact.ID}
	getContact, err := c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, []interface{}{getContact}, paginationMeta, "Contact updated successfully", http.StatusOK, nil, nil)
}

// delete contact
func (c *ContactController) DeleteContactByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.DeleteContactRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, "ID is required", nil)
	}

	params := &dtos.GetContactParams{ID: req.ID}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	params.UserID = userID

	// GET contact by ID
	_, err = c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.DeleteContact(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to delete contact", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Contact deleted successfully", http.StatusOK, nil, nil)
}

// restore contact
func (c *ContactController) RestoreContactByAuthUser(ctx *fiber.Ctx) error {
	var req dtos.DeleteContactRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, err.Error(), nil)
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusBadRequest, "ID is required", nil)
	}

	// Extract user ID from JWT
	claims, err := middleware.GetAuthUser(ctx)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Unauthorized", http.StatusUnauthorized, err.Error(), nil)
	}
	userID := uint(claims["user_id"].(float64))

	isDeleted := 1

	params := &dtos.GetContactParams{ID: req.ID, IsDeleted: &isDeleted, UserID: userID}

	// GET contact by ID
	_, err = c.service.GetContactByID(ctx.Context(), params)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Contact not found", http.StatusNotFound, err.Error(), nil)
	}

	err = c.service.RestoreContact(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "Failed to restore contact", http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.GetResponse(ctx, nil, nil, "Contact restored successfully", http.StatusOK, nil, nil)
}
