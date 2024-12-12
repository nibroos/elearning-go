package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupSubscribeRoutes(subscribes fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	subscribeRepo := repository.NewSubscribeRepository(gormDB, sqlDB)
	subscribeService := service.NewSubscribeService(subscribeRepo)
	subscribeController := rest.NewSubscribeController(subscribeService)

	// prefix /subscribes

	subscribes.Post("/index-subscribe", subscribeController.GetSubscribes)
	subscribes.Post("/index-subscribe-r", subscribeController.GetSubscribesFromRedis)
	subscribes.Post("/show-subscribe", subscribeController.GetSubscribeByID)
	subscribes.Post("/create-subscribe", subscribeController.CreateSubscribe)
	subscribes.Post("/update-subscribe", subscribeController.UpdateSubscribe)
	subscribes.Post("/delete-subscribe", subscribeController.DeleteSubscribe)
	subscribes.Post("/restore-subscribe", subscribeController.RestoreSubscribe)
}
