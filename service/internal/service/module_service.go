package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type ModuleService struct {
	repo *repository.ModuleRepository
}

func NewModuleService(repo *repository.ModuleRepository) *ModuleService {
	return &ModuleService{repo: repo}
}

func (s *ModuleService) GetModules(ctx context.Context, filters map[string]string) ([]dtos.ModuleListDTO, int, error) {

	resultChan := make(chan dtos.GetModulesResult, 1)

	go func() {
		modules, total, err := s.repo.GetModules(ctx, filters)
		resultChan <- dtos.GetModulesResult{Modules: modules, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Modules, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *ModuleService) CreateModule(ctx context.Context, module *models.Module) (*models.Module, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create module
	if err := s.repo.CreateModule(tx, module); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return module, nil
}

func (s *ModuleService) GetModuleByID(ctx context.Context, id uint) (*dtos.ModuleDetailDTO, error) {
	moduleChan := make(chan *dtos.ModuleDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		module, err := s.repo.GetModuleByID(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		moduleChan <- module
	}()

	select {
	case module := <-moduleChan:
		return module, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *ModuleService) UpdateModule(ctx context.Context, module *models.Module) (*models.Module, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update module
	if err := s.repo.UpdateModule(tx, module); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return module, nil
}

func (s *ModuleService) DeleteModule(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete module
	if err := s.repo.DeleteModule(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *ModuleService) RestoreModule(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore module
	if err := s.repo.RestoreModule(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
