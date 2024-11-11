package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupRecordRoutes(records fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	recordRepo := repository.NewRecordRepository(gormDB, sqlDB)
	recordService := service.NewRecordService(recordRepo)
	recordController := rest.NewRecordController(recordService)

	// prefix /records

	records.Post("/index-record", recordController.ListRecords)
	records.Post("/show-record", recordController.GetRecordByID)
	records.Post("/create-record", recordController.CreateRecord)
	records.Post("/update-record", recordController.UpdateRecord)
	records.Post("/delete-record", recordController.DeleteRecord)
	records.Post("/restore-record", recordController.RestoreRecord)
	records.Post("/auth-index-record", recordController.ListRecordsByAuthUser)
	records.Post("/auth-show-record", recordController.GetRecordByIDByAuthUser)
	records.Post("/auth-create-record", recordController.CreateRecordByAuthUser)
	records.Post("/auth-update-record", recordController.UpdateRecordByAuthUser)
	records.Post("/auth-delete-record", recordController.DeleteRecordByAuthUser)
	records.Post("/auth-restore-record", recordController.RestoreRecordByAuthUser)
}
