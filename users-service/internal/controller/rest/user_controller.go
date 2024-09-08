package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/models"
	"github.com/nibroos/elearning-go/users-service/internal/service"
	"github.com/nibroos/elearning-go/users-service/internal/utils"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
	// Initialize the filters struct
	var GetUsersRequest dtos.GetUsersRequest

	// Parse the request body into the GetUsersRequest struct
	if err := ctx.BodyParser(&GetUsersRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"status":  "error",
			"err":     err.Error(),
		})
	}
	// return utils.DD(ctx, "ABC", GetUsersRequest, []int{1, 2, 3}, map[string]string{"key": "value"})

	// Convert the filters struct to a map
	filters := utils.ConvertStructToMap(GetUsersRequest)

	users, total, err := c.service.GetUsers(ctx.Context(), filters)
	if err != nil {
		// return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, utils.ErrorWithLocation(err), "error"), http.StatusInternalServerError)
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	currentPage := utils.AtoiDefault(filters["page"], 1)
	lastPage := (total + GetUsersRequest.PerPage.Value - 1) / GetUsersRequest.PerPage.Value

	paginationMeta := &utils.PaginationMeta{
		Total:       total,
		PerPage:     GetUsersRequest.PerPage.Value,
		CurrentPage: currentPage,
		LastPage:    lastPage,
	}

	response := utils.WrapResponse(users, paginationMeta, "Users fetched successfully", http.StatusOK)
	return utils.SendResponse(ctx, response, http.StatusOK)
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req dtos.CreateUserRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Address:  req.Address,
	}

	if err := c.service.CreateUser(&user, req.RoleIDs); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusInternalServerError})
	}

	createdUser, err := c.service.GetUserByID(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusInternalServerError})
	}

	paginationMeta := &utils.PaginationMeta{}

	response := utils.WrapResponse(createdUser, paginationMeta, "Users fetched successfully", http.StatusOK)
	return utils.SendResponse(ctx, response, http.StatusOK)
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	var req dtos.GetUserByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	user, err := c.service.GetUserByID(req.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Invalid request", "status": http.StatusInternalServerError})
	}

	paginationMeta := &utils.PaginationMeta{}

	response := utils.WrapResponse(user, paginationMeta, "Users fetched successfully", http.StatusOK)
	return utils.SendResponse(ctx, response, http.StatusOK)
}
