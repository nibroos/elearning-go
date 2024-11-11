package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"gorm.io/gorm"
)

type RecordRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewRecordRepository(db *gorm.DB, sqlDB *sqlx.DB) *RecordRepository {
	return &RecordRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *RecordRepository) ListRecords(ctx context.Context, filters map[string]string) ([]dtos.RecordListDTO, int, error) {
	records := []dtos.RecordListDTO{}
	var total int

	from := `FROM (
        SELECT r.id, r.user_id, r.education_id, r.time_spent, r.created_at, r.updated_at,
        u.name as user_name,
        e.name as education_name

        FROM records r
        JOIN users u ON r.user_id = u.id
        JOIN educations e ON r.education_id = e.id
        WHERE r.deleted_at IS NULL
    ) AS alias WHERE 1=1`

	query := `SELECT * ` + from
	countQuery := `SELECT COUNT(*) ` + from

	var args []interface{}
	i := 1
	for key, value := range filters {
		switch key {
		case "education_name", "user_name":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				countQuery += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	if value, ok := filters["user_id"]; ok && value != "" {
		query += fmt.Sprintf(" AND user_id = $%d", i)
		countQuery += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, value)
		i++
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (user_name ILIKE $%d OR education_name ILIKE $%d)", i, i+1)
		countQuery += fmt.Sprintf(" AND (user_name ILIKE $%d OR education_name ILIKE $%d)", i, i+1)
		args = append(args, "%"+value+"%", "%"+value+"%", "%"+value+"%")
		i += 3
	}

	countArgs := append([]interface{}{}, args...)

	allowedOrderColumns := []string{"id", "ref_num", "user_name", "education_name"}
	orderColumn := utils.GetStringOrDefaultFromArray(filters["order_column"], allowedOrderColumns, "id")
	orderDirection := utils.GetStringOrDefault(filters["order_direction"], "asc")
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, (currentPage-1)*perPage)

	// Channels for concurrent execution
	countChan := make(chan error)
	selectChan := make(chan error)

	// Goroutine for count query
	go func() {
		err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
		countChan <- err
	}()

	// Goroutine for select query
	go func() {
		err := r.sqlDB.SelectContext(ctx, &records, query, args...)
		selectChan <- err
	}()

	// Wait for both goroutines to finish
	countErr := <-countChan
	selectErr := <-selectChan

	if countErr != nil {
		return nil, 0, countErr
	}

	if selectErr != nil {
		return nil, 0, selectErr
	}

	return records, total, nil
}

func (r *RecordRepository) GetRecordByID(ctx context.Context, params *dtos.GetRecordParams) (*dtos.RecordDetailDTO, error) {
	var record dtos.RecordDetailDTO
	// deletedAt := params.IsDeleted

	query := `SELECT r.id, r.user_id, r.education_id, r.time_spent, r.created_at, r.updated_at, r.deleted_at,
	u.name as user_name,
	e.name as education_name

	FROM records r
	JOIN users u ON r.user_id = u.id
	JOIN educations e ON r.education_id = e.id
	WHERE 1=1`

	var args []interface{}

	i := 1
	query += " AND r.id = $1"
	args = append(args, params.ID)
	i++

	isDeletedQuery := ` AND r.deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND r.deleted_at IS NOT NULL"
	}

	if params.UserID != 0 {
		query += fmt.Sprintf(" AND r.user_id = $%d", i)
		args = append(args, params.UserID)
		i++
	}

	query += isDeletedQuery

	if err := r.sqlDB.Get(&record, query, args...); err != nil {
		return nil, err
	}

	return &record, nil
}

// BeginTransaction starts a new transaction
func (r *RecordRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *RecordRepository) CreateRecord(tx *gorm.DB, record *models.Record) error {
	if err := tx.Create(record).Error; err != nil {
		return err
	}
	return nil
}

func (r *RecordRepository) UpdateRecord(tx *gorm.DB, record *models.Record) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *RecordRepository) DeleteRecord(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Record{}, id).Error; err != nil {
		if err := tx.Delete(&models.Record{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *RecordRepository) RestoreRecord(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE records SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
