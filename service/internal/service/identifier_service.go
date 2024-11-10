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

func (s *IdentifierService) ListIdentifiers(ctx context.Context, filters map[string]string) ([]dtos.IdentifierListDTO, int, error) {

	resultChan := make(chan dtos.ListIdentifiersResult, 1)

	go func() {
		identifiers, total, err := s.repo.ListIdentifiers(ctx, filters)
		resultChan <- dtos.ListIdentifiersResult{Identifiers: identifiers, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Identifiers, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *IdentifierService) CreateIdentifier(ctx context.Context, identifier *models.Identifier) (*models.Identifier, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create identifier
	if err := s.repo.CreateIdentifier(tx, identifier); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return identifier, nil
}

func (s *IdentifierService) GetIdentifierByID(ctx context.Context, params *dtos.GetIdentifierParams) (*dtos.IdentifierDetailDTO, error) {
	identifierChan := make(chan *dtos.IdentifierDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		identifier, err := s.repo.GetIdentifierByID(ctx, params)
		if err != nil {
			errChan <- err
			return
		}
		identifierChan <- identifier
	}()

	select {
	case identifier := <-identifierChan:
		return identifier, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *IdentifierService) UpdateIdentifier(ctx context.Context, identifier *models.Identifier) (*models.Identifier, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update identifier
	if err := s.repo.UpdateIdentifier(tx, identifier); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return identifier, nil
}

func (s *IdentifierService) DeleteIdentifier(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete identifier
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

	// Restore identifier
	if err := s.repo.RestoreIdentifier(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *IdentifierService) ListIdentifiersByAuthUser(ctx context.Context, filters map[string]string) ([]dtos.IdentifierListDTO, int, error) {
	resultChan := make(chan dtos.ListIdentifiersResult, 1)

	go func() {
		identifiers, total, err := s.repo.ListIdentifiers(ctx, filters)
		resultChan <- dtos.ListIdentifiersResult{Identifiers: identifiers, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Identifiers, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}
