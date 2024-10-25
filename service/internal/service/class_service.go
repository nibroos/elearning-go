package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type ClassService struct {
	repo *repository.ClassRepository
}

func NewClassService(repo *repository.ClassRepository) *ClassService {
	return &ClassService{repo: repo}
}

func (s *ClassService) GetClasses(ctx context.Context, filters map[string]string) ([]dtos.ClassListDTO, int, error) {

	resultChan := make(chan dtos.GetClassesResult, 1)

	go func() {
		subscribes, total, err := s.repo.GetClasses(ctx, filters)
		resultChan <- dtos.GetClassesResult{Classes: subscribes, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Classes, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *ClassService) CreateClass(ctx context.Context, subscribe *models.Class) (*models.Class, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create subscribe
	if err := s.repo.CreateClass(tx, subscribe); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return subscribe, nil
}

func (s *ClassService) GetClassByID(ctx context.Context, id uint) (*dtos.ClassDetailDTO, error) {
	subscribeChan := make(chan *dtos.ClassDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		subscribe, err := s.repo.GetClassByID(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		subscribeChan <- subscribe
	}()

	select {
	case subscribe := <-subscribeChan:
		return subscribe, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *ClassService) UpdateClass(ctx context.Context, subscribe *models.Class) (*models.Class, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update subscribe
	if err := s.repo.UpdateClass(tx, subscribe); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return subscribe, nil
}

func (s *ClassService) DeleteClass(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete subscribe
	if err := s.repo.DeleteClass(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *ClassService) RestoreClass(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore subscribe
	if err := s.repo.RestoreClass(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
