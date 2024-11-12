package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/controller/rest"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/service"
	"gorm.io/gorm"
)

func SetupQuizRoutes(quizes fiber.Router, gormDB *gorm.DB, sqlDB *sqlx.DB) {
	quizRepo := repository.NewQuizRepository(gormDB, sqlDB)
	quizService := service.NewQuizService(quizRepo)
	quizController := rest.NewQuizController(quizService)

	// prefix /quizes

	quizes.Post("/index-quiz", quizController.ListQuizes)
	quizes.Post("/show-quiz", quizController.GetQuizByID)
	quizes.Post("/create-quiz", quizController.CreateQuiz)
	quizes.Post("/update-quiz", quizController.UpdateQuiz)
	quizes.Post("/delete-quiz", quizController.DeleteQuiz)
	quizes.Post("/restore-quiz", quizController.RestoreQuiz)
}
