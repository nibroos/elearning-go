package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/middleware"
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

	getUser, err := c.service.GetUserByID(ctx.Context(), createdUser.ID)
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

	user, err := c.service.GetUserByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusNotFound, err.Error())
	}

	userArray := []interface{}{user}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, userArray, paginationMeta, "User fetched successfully", http.StatusOK)
}

// update user
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var req dtos.UpdateUserRequest

	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	validationErrors := validators.ValidateUpdateUserRequest(&req)
	if validationErrors != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": validationErrors, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	// Fetch the existing user to get the current password if needed
	existingUser, err := c.service.GetUserByID(ctx.Context(), req.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusNotFound, err.Error())
	}

	user := models.User{
		ID:       req.ID,
		Name:     req.Name,
		Username: req.Username.Value,
		Email:    req.Email,
		Password: *existingUser.Password.Value, // Default to existing password
		Address:  req.Address.Value,
	}

	// Update password only if a new one is provided
	if req.Password.Value != nil {
		user.Password = *req.Password.Value
	}

	updatedUser, err := c.service.UpdateUser(ctx.Context(), &user, req.RoleIDs)
	if err != nil {
		if err.Error() == "username already exists" {
			return ctx.Status(http.StatusConflict).JSON(fiber.Map{"errors": err.Error(), "message": "Username already exists", "status": http.StatusConflict})
		}
		return utils.GetResponse(ctx, nil, nil, "Failed to update user", http.StatusInternalServerError, err.Error())
	}

	getUser, err := c.service.GetUserByID(ctx.Context(), updatedUser.ID)
	if err != nil {
		return utils.GetResponse(ctx, nil, nil, "User not found", http.StatusNotFound, err.Error())
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, getUser, paginationMeta, "User updated successfully", http.StatusOK)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req dtos.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request", "status": "error", "err": err.Error()})
	}

	user, err := c.service.Authenticate(ctx.Context(), req.Email, req.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials", "status": "error", "err": err.Error()})
	}

	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token", "status": "error", "err": err.Error()})
	}

	return ctx.JSON(fiber.Map{"token": token})
}
func (c *UserController) Register(ctx *fiber.Ctx) error {
	var req dtos.RegisterRequest

	// Use the utility function to parse the request body
	if err := utils.BodyParserWithNull(ctx, &req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": err.Error(), "message": "Invalid request", "status": http.StatusBadRequest})
	}

	// Validate the request
	validationErrors := validators.ValidateRegisterRequest(&req)
	if validationErrors != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": validationErrors, "message": "Validation failed", "status": http.StatusBadRequest})
	}

	utils.DD(ctx.Context(), map[string]interface{}{
		"req":  req,
		"test": 'a',
		// "valid": validationErrors,
	})

	user := models.User{
		Name:     req.Name,
		Username: req.Username.Value,
		Email:    req.Email,
		Password: req.Password,
	}

	roleIDS := []uint32{utils.RoleStudent}

	createdUser, err := c.service.CreateUser(ctx.Context(), &user, roleIDS)
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

	token, err := middleware.GenerateJWT(createdUser.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token", "status": "error", "err": err.Error()})
	}

	filters := ctx.Locals("filters").(map[string]string)
	paginationMeta := utils.CreatePaginationMeta(filters, 1)

	return utils.GetResponse(ctx, getUser, paginationMeta, "User registered successfully", http.StatusCreated, fiber.Map{"token": token})
}
