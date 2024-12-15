package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type SubscribeService struct {
	repo *repository.SubscribeRepository
}

func NewSubscribeService(repo *repository.SubscribeRepository) *SubscribeService {
	return &SubscribeService{repo: repo}
}

func (s *SubscribeService) GetSubscribes(ctx context.Context, filters map[string]string) ([]dtos.SubscribeListDTO, int, error) {

	resultChan := make(chan dtos.GetSubscribesResult, 1)

	go func() {
		subscribes, total, err := s.repo.GetSubscribes(ctx, filters)
		resultChan <- dtos.GetSubscribesResult{Subscribes: subscribes, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Subscribes, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *SubscribeService) GetSubscribesFromRedis(ctx context.Context, filters map[string]string) ([]dtos.SubscribeListDTO, int, error) {
	resultChan := make(chan dtos.GetSubscribesResult, 1)

	go func() {
		subscribes, total, err := s.repo.GetSubscribesFromRedis(ctx, filters)
		resultChan <- dtos.GetSubscribesResult{Subscribes: subscribes, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Subscribes, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *SubscribeService) CreateSubscribe(ctx context.Context, subscribe *models.Subscribe) (*models.Subscribe, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create subscribe
	if err := s.repo.CreateSubscribe(tx, subscribe); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return subscribe, nil
}

func (s *SubscribeService) GetSubscribeByID(ctx context.Context, id uint) (*dtos.SubscribeDetailDTO, error) {
	subscribeChan := make(chan *dtos.SubscribeDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		subscribe, err := s.repo.GetSubscribeByID(ctx, id)
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

func (s *SubscribeService) UpdateSubscribe(ctx context.Context, subscribe *models.Subscribe) (*models.Subscribe, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update subscribe
	if err := s.repo.UpdateSubscribe(tx, subscribe); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return subscribe, nil
}

func (s *SubscribeService) DeleteSubscribe(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete subscribe
	if err := s.repo.DeleteSubscribe(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *SubscribeService) RestoreSubscribe(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore subscribe
	if err := s.repo.RestoreSubscribe(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
