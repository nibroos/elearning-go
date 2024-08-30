package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/users-service/internal/service"
)

type UserController struct {
    service service.UserService
    tx *sqlx.Tx
}

func NewUserController(service service.UserService) *UserController {
    return &UserController{service: service}
}

// func (c *UserController) GetUsers(ctx *fiber.Ctx) {
//     searchParams := map[string]string{
//         "global": ctx.Query("global"),
//         "name":   ctx.Query("name"),
//         "email":  ctx.Query("email"),
//     }

//     users, err := c.service.GetUsers(searchParams)
//     if err != nil {
//         ctx.JSON(http.StatusInternalServerError, fiber.Error{"error": err.Error()})
//         return
//     }

//     ctx.JSON(http.StatusOK, users)
// }

func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
    // Extract query parameters for search
    searchParams := map[string]string{
        "global": ctx.Query("global"),
        "name":   ctx.Query("name"),
        "email":  ctx.Query("email"),
    }

    // Fetch users from the service layer
    users, err := c.service.GetUsers(searchParams)
    if err != nil {
        // Return an error response
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Return the list of users as a JSON response
    return ctx.Status(fiber.StatusOK).JSON(users)
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
    var body struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
        RoleID   int64  `json:"role_id"`
    }

    if err := ctx.BodyParser(&body); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }

    user, err := c.service.CreateUser(ctx.Context(), c.tx, body.Name, body.Email, body.Password, body.RoleID)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return ctx.JSON(fiber.Map{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
    })
}