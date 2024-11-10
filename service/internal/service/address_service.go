package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type AddressService struct {
	repo *repository.AddressRepository
}

func NewAddressService(repo *repository.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) ListAddresses(ctx context.Context, filters map[string]string) ([]dtos.AddressListDTO, int, error) {

	resultChan := make(chan dtos.ListAddressesResult, 1)

	go func() {
		addresses, total, err := s.repo.ListAddresses(ctx, filters)
		resultChan <- dtos.ListAddressesResult{Addresses: addresses, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Addresses, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *AddressService) CreateAddress(ctx context.Context, address *models.Address) (*models.Address, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create address
	if err := s.repo.CreateAddress(tx, address); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return address, nil
}

func (s *AddressService) GetAddressByID(ctx context.Context, params *dtos.GetAddressParams) (*dtos.AddressDetailDTO, error) {
	addressChan := make(chan *dtos.AddressDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		address, err := s.repo.GetAddressByID(ctx, params)
		if err != nil {
			errChan <- err
			return
		}
		addressChan <- address
	}()

	select {
	case address := <-addressChan:
		return address, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *AddressService) UpdateAddress(ctx context.Context, address *models.Address) (*models.Address, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update address
	if err := s.repo.UpdateAddress(tx, address); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return address, nil
}

func (s *AddressService) DeleteAddress(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete address
	if err := s.repo.DeleteAddress(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *AddressService) RestoreAddress(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore address
	if err := s.repo.RestoreAddress(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
