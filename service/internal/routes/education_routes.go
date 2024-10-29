package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupEducationRoutes(educations fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	educationRepo := repository.NewEducationRepository(gormDB, sqlDB)
	educationService := service.NewEducationService(educationRepo)
	educationController := rest.NewEducationController(educationService)

	// prefix /educations

	educations.Post("/index-education", educationController.GetEducations)
	educations.Post("/show-education", educationController.GetEducationByID)
	educations.Post("/create-education", educationController.CreateEducation)
	educations.Post("/update-education", educationController.UpdateEducation)
	educations.Post("/delete-education", educationController.DeleteEducation)
	educations.Post("/restore-education", educationController.RestoreEducation)
}
