package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type EducationService struct {
	repo *repository.EducationRepository
}

func NewEducationService(repo *repository.EducationRepository) *EducationService {
	return &EducationService{repo: repo}
}

func (s *EducationService) GetEducations(ctx context.Context, filters map[string]string) ([]dtos.EducationListDTO, int, error) {

	resultChan := make(chan dtos.GetEducationsResult, 1)

	go func() {
		educations, total, err := s.repo.GetEducations(ctx, filters)
		resultChan <- dtos.GetEducationsResult{Educations: educations, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Educations, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *EducationService) CreateEducation(ctx context.Context, education *models.Education) (*models.Education, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create education
	if err := s.repo.CreateEducation(tx, education); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return education, nil
}

func (s *EducationService) GetEducationByID(ctx context.Context, id uint) (*dtos.EducationDetailDTO, error) {
	educationChan := make(chan *dtos.EducationDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		education, err := s.repo.GetEducationByID(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		educationChan <- education
	}()

	select {
	case education := <-educationChan:
		return education, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *EducationService) UpdateEducation(ctx context.Context, education *models.Education) (*models.Education, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update education
	if err := s.repo.UpdateEducation(tx, education); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return education, nil
}

func (s *EducationService) DeleteEducation(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete education
	if err := s.repo.DeleteEducation(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *EducationService) RestoreEducation(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore education
	if err := s.repo.RestoreEducation(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
