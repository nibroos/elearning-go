package rest

import (
	"database/sql"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/nibroos/elearning-go/users-service/internal/utils"
)

type SeederController struct {
	db *sql.DB
}

func NewSeederController(db *sql.DB) *SeederController {
	return &SeederController{db: db}
}

func (c *SeederController) RunSeeders(ctx *fiber.Ctx) error {
	seedFiles := []string{
		"001_create_roles.sql",
		"002_create_permissions.sql",
		"003_create_roles_values.sql",
		"004_create_permissions_values.sql",
	}

	// Prepend the directory path to each seed file
	for i, file := range seedFiles {
		seedFiles[i] = filepath.Join("internal", "database", "seeders", file)
	}

	err := utils.ExecuteSeeders(c.db, seedFiles)
	if err != nil {
		return utils.JSONError(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(fiber.Map{
		"message": "Seeders executed successfully",
	})
}
