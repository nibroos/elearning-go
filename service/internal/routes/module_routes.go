package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupModuleRoutes(modules fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	moduleRepo := repository.NewModuleRepository(gormDB, sqlDB)
	moduleService := service.NewModuleService(moduleRepo)
	moduleController := rest.NewModuleController(moduleService)

	// prefix /modules

	modules.Post("/index-module", moduleController.GetModules)
	modules.Post("/show-module", moduleController.GetModuleByID)
	modules.Post("/create-module", moduleController.CreateModule)
	modules.Post("/update-module", moduleController.UpdateModule)
	modules.Post("/delete-module", moduleController.DeleteModule)
	modules.Post("/restore-module", moduleController.RestoreModule)
}
