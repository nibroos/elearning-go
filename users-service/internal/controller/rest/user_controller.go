package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/models"
	"github.com/nibroos/elearning-go/users-service/internal/service"
	"github.com/nibroos/elearning-go/users-service/internal/utils"
	"github.com/nibroos/elearning-go/users-service/internal/validators"
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

	return utils.GetResponse(ctx, users, paginationMeta, "Users fetched successfully", http.StatusOK)
}
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req dtos.CreateUserRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	validationErrors := validators.ValidateCreateUserRequest(&req)
	if validationErrors != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": validationErrors, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username.Value,
		Email:    req.Email,
		Password: req.Password,
		Address:  req.Address.Value,
	}

	createdUser, err := c.service.CreateUser(ctx.Context(), &user, req.RoleIDs)
	if err != nil {
		if err.Error() == "username already exists" {
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{"errors": err.Error(), "message": "Username already exists", "status": http.StatusConflict})
		}
		return utils.GetResponse(ctx, nil, nil, "Failed to create user", http.StatusInternalServerError, err.Error())
	}

	getUser, err := c.service.GetUserByID(ctx.Context(), uint32(createdUser.ID))
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusNotFound, err.Error())
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, getUser, paginationMeta, "User created successfully", http.StatusCreated)
}
func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	var req dtos.GetUserByIDRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusBadRequest, err.Error())
	}

	if req.ID == 0 {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusBadRequest, "ID is required")
	}

	user, err := c.service.GetUserByID(ctx.Context(), uint32(req.ID))
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusNotFound, err.Error())
	}

	userArray := []interface{}{user}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, userArray, paginationMeta, "User fetched successfully", http.StatusOK)
}
