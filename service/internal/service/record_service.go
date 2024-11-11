package service

import (
	"context"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
)

type RecordService struct {
	repo *repository.RecordRepository
}

func NewRecordService(repo *repository.RecordRepository) *RecordService {
	return &RecordService{repo: repo}
}

func (s *RecordService) ListRecords(ctx context.Context, filters map[string]string) ([]dtos.RecordListDTO, int, error) {

	resultChan := make(chan dtos.ListRecordsResult, 1)

	go func() {
		records, total, err := s.repo.ListRecords(ctx, filters)
		resultChan <- dtos.ListRecordsResult{Records: records, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Records, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}

func (s *RecordService) CreateRecord(ctx context.Context, record *models.Record) (*models.Record, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create record
	if err := s.repo.CreateRecord(tx, record); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (s *RecordService) GetRecordByID(ctx context.Context, params *dtos.GetRecordParams) (*dtos.RecordDetailDTO, error) {
	recordChan := make(chan *dtos.RecordDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		record, err := s.repo.GetRecordByID(ctx, params)
		if err != nil {
			errChan <- err
			return
		}
		recordChan <- record
	}()

	select {
	case record := <-recordChan:
		return record, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *RecordService) UpdateRecord(ctx context.Context, record *models.Record) (*models.Record, error) {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update record
	if err := s.repo.UpdateRecord(tx, record); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return record, nil
}

func (s *RecordService) DeleteRecord(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Delete record
	if err := s.repo.DeleteRecord(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *RecordService) RestoreRecord(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore record
	if err := s.repo.RestoreRecord(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
