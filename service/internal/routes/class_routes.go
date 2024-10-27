package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupClassRoutes(classes fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	classRepo := repository.NewClassRepository(gormDB, sqlDB)
	classService := service.NewClassService(classRepo)
	classController := rest.NewClassController(classService)

	// prefix /classes

	classes.Post("/index-class", classController.GetClasss)
	classes.Post("/show-class", classController.GetClassByID)
	classes.Post("/create-class", classController.CreateClass)
	classes.Post("/update-class", classController.UpdateClass)
	classes.Post("/delete-class", classController.DeleteClass)
	classes.Post("/restore-class", classController.RestoreClass)
}
