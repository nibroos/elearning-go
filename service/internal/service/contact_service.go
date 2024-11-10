package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type ContactService struct {
	repo *repository.ContactRepository
}

func NewContactService(repo *repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) ListContacts(ctx context.Context, filters map[string]string) ([]dtos.ContactListDTO, int, error) {

	resultChan := make(chan dtos.ListContactsResult, 1)

	go func() {
		contacts, total, err := s.repo.ListContacts(ctx, filters)
		resultChan <- dtos.ListContactsResult{Contacts: contacts, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Contacts, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *ContactService) CreateContact(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create contact
	if err := s.repo.CreateContact(tx, contact); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactService) GetContactByID(ctx context.Context, params *dtos.GetContactParams) (*dtos.ContactDetailDTO, error) {
	contactChan := make(chan *dtos.ContactDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		contact, err := s.repo.GetContactByID(ctx, params)
		if err != nil {
			errChan <- err
			return
		}
		contactChan <- contact
	}()

	select {
	case contact := <-contactChan:
		return contact, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *ContactService) UpdateContact(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update contact
	if err := s.repo.UpdateContact(tx, contact); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactService) DeleteContact(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete contact
	if err := s.repo.DeleteContact(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *ContactService) RestoreContact(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore contact
	if err := s.repo.RestoreContact(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *ContactService) ListContactsByAuthUser(ctx context.Context, filters map[string]string) ([]dtos.ContactListDTO, int, error) {
	resultChan := make(chan dtos.ListContactsResult, 1)

	go func() {
		contacts, total, err := s.repo.ListContacts(ctx, filters)
		resultChan <- dtos.ListContactsResult{Contacts: contacts, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Contacts, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}
