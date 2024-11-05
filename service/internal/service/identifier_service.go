package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type IdentifierService struct {
	repo *repository.IdentifierRepository
}

func NewIdentifierService(repo *repository.IdentifierRepository) *IdentifierService {
	return &IdentifierService{repo: repo}
}

func (s *IdentifierService) GetIdentifiers(ctx context.Context, filters map[string]string) ([]dtos.IdentifierListDTO, int, error) {

	resultChan := make(chan dtos.GetIdentifiersResult, 1)

	go func() {
		identifers, total, err := s.repo.GetIdentifiers(ctx, filters)
		resultChan <- dtos.GetIdentifiersResult{Identifiers: identifers, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Identifiers, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *IdentifierService) CreateIdentifier(ctx context.Context, identifer *models.Identifier) (*models.Identifier, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create identifer
	if err := s.repo.CreateIdentifier(tx, identifer); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return identifer, nil
}

func (s *IdentifierService) GetIdentifierByID(ctx context.Context, params *dtos.GetIdentifierParams) (*dtos.IdentifierDetailDTO, error) {
	identiferChan := make(chan *dtos.IdentifierDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		identifer, err := s.repo.GetIdentifierByID(ctx, params)
		if err != nil {
			errChan <- err
			return
		}
		identiferChan <- identifer
	}()

	select {
	case identifer := <-identiferChan:
		return identifer, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *IdentifierService) UpdateIdentifier(ctx context.Context, identifer *models.Identifier) (*models.Identifier, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update identifer
	if err := s.repo.UpdateIdentifier(tx, identifer); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return identifer, nil
}

func (s *IdentifierService) DeleteIdentifier(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete identifer
	if err := s.repo.DeleteIdentifier(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *IdentifierService) RestoreIdentifier(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore identifer
	if err := s.repo.RestoreIdentifier(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
