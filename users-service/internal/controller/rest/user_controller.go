package rest

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
	filters, ok := ctx.Locals("filters").(map[string]string)
	if !ok {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, "Invalid filters", http.StatusBadRequest), http.StatusBadRequest)
	}

	users, total, err := c.service.GetUsers(ctx.Context(), filters)
	if err != nil {
		return utils.SendResponse(ctx, utils.WrapResponse(nil, nil, err.Error(), http.StatusInternalServerError), http.StatusInternalServerError)
	}

	paginationMeta := utils.CreatePaginationMeta(filters, total)
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

	createdUser, err := c.service.CreateUser(ctx.Context(), &user, req.RoleIDs)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusInternalServerError})
	}

	getUser, err := c.service.GetUserByID(ctx.Context(), uint32(createdUser.ID))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusInternalServerError})
	}

	paginationMeta := &utils.Meta{}

	response := utils.WrapResponse(getUser, paginationMeta, "Users fetched successfully", http.StatusOK)
	return utils.SendResponse(ctx, response, http.StatusOK)
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	var req dtos.GetUserByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	user, err := c.service.GetUserByID(ctx.Context(), uint32(req.ID))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "message": "Invalid request", "status": http.StatusInternalServerError})
	}

	response := utils.WrapResponse(user, nil, "User fetched successfully", http.StatusOK)
	return utils.SendResponse(ctx, response, http.StatusOK)
}
