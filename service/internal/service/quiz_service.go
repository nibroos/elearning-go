package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type QuizService struct {
	repo *repository.QuizRepository
}

func NewQuizService(repo *repository.QuizRepository) *QuizService {
	return &QuizService{repo: repo}
}

func (s *QuizService) ListQuizes(ctx context.Context, filters map[string]string) ([]dtos.QuizListDTO, int, error) {

	resultChan := make(chan dtos.ListQuizesResult, 1)

	go func() {
		quizes, total, err := s.repo.ListQuizes(ctx, filters)
		resultChan <- dtos.ListQuizesResult{Quizes: quizes, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Quizes, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *QuizService) CreateQuiz(ctx context.Context, quiz *models.Quiz) (*models.Quiz, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create quiz
	if err := s.repo.CreateQuiz(tx, quiz); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func (s *QuizService) GetQuizByID(ctx context.Context, params *dtos.GetQuizParams) (*dtos.QuizDetailDTO, error) {
	quizChan := make(chan *dtos.QuizDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		quiz, err := s.repo.GetQuizByID(ctx, params)
		if err != nil {
			errChan <- err
			return
		}
		quizChan <- quiz
	}()

	select {
	case quiz := <-quizChan:
		return quiz, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *QuizService) UpdateQuiz(ctx context.Context, quiz *models.Quiz) (*models.Quiz, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update quiz
	if err := s.repo.UpdateQuiz(tx, quiz); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return quiz, nil
}

func (s *QuizService) DeleteQuiz(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete quiz
	if err := s.repo.DeleteQuiz(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *QuizService) RestoreQuiz(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore quiz
	if err := s.repo.RestoreQuiz(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
