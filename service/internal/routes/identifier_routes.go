package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupIdentifierRoutes(identifiers fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	identifierRepo := repository.NewIdentifierRepository(gormDB, sqlDB)
	identifierService := service.NewIdentifierService(identifierRepo)
	identifierController := rest.NewIdentifierController(identifierService)

	// prefix /identifiers

	identifiers.Post("/index-identifier", identifierController.GetIdentifiers)
	identifiers.Post("/show-identifier", identifierController.GetIdentifierByID)
	identifiers.Post("/create-identifier", identifierController.CreateIdentifier)
	identifiers.Post("/update-identifier", identifierController.UpdateIdentifier)
	identifiers.Post("/delete-identifier", identifierController.DeleteIdentifier)
	identifiers.Post("/restore-identifier", identifierController.RestoreIdentifier)
}
