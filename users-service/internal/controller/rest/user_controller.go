package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/users-service/internal/service"
)

type UserController struct {
    service service.UserService
}

func NewUserController(service service.UserService) *UserController {
    return &UserController{service: service}
}

func (c *UserController) GetUsers(ctx *fiber.Ctx) {
    searchParams := map[string]string{
        "global": ctx.Query("global"),
        "name":   ctx.Query("name"),
        "email":  ctx.Query("email"),
    }

    users, err := c.service.GetUsers(searchParams)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, fiber.Error{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, users)
}
