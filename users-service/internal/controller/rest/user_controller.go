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
		})
	}

	if GetUsersRequest.PerPage == "" {
		GetUsersRequest.PerPage = "10"
	}
	if GetUsersRequest.Page == "" {
		GetUsersRequest.Page = "1"
	}
	if GetUsersRequest.OrderColumn == "" {
		GetUsersRequest.OrderColumn = "users.name"
	}
	if GetUsersRequest.OrderDirection == "" {
		GetUsersRequest.OrderDirection = "desc"
	}

	// Convert the filters struct to a map
	filters := utils.ConvertStructToMap(GetUsersRequest)

	users, total, err := c.service.GetUsers(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Error fetching users", "error"), http.StatusInternalServerError)
	}

	perPage := utils.AtoiDefault(filters["per_page"], 10)
	currentPage := utils.AtoiDefault(filters["page"], 1)
	lastPage := (total + perPage - 1) / perPage

	paginationMeta := &utils.PaginationMeta{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		LastPage:    lastPage,
	}

	response := utils.WrapResponse(users, paginationMeta, "Users fetched successfully", string(rune(http.StatusOK)))
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

	response := utils.WrapResponse(createdUser, paginationMeta, "Users fetched successfully", string(rune(http.StatusOK)))
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

	response := utils.WrapResponse(user, paginationMeta, "Users fetched successfully", string(rune(http.StatusOK)))
	return utils.SendResponse(ctx, response, http.StatusOK)
}
