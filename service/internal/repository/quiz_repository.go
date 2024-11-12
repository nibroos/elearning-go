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

type QuizRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewQuizRepository(db *gorm.DB, sqlDB *sqlx.DB) *QuizRepository {
	return &QuizRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *QuizRepository) ListQuizes(ctx context.Context, filters map[string]string) ([]dtos.QuizListDTO, int, error) {
	quizes := []dtos.QuizListDTO{}
	var total int

	from := `FROM (
        SELECT q.id, q.name, q.description, q.threshold, q.created_at, q.updated_at

        FROM quizes q
        WHERE q.deleted_at IS NULL
    ) AS alias WHERE 1=1`

	query := `SELECT * ` + from
	countQuery := `SELECT COUNT(*) ` + from

	var args []interface{}
	i := 1
	for key, value := range filters {
		switch key {
		case "name", "description":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				countQuery += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", i, i+1)
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", i, i+1)
		args = append(args, "%"+value+"%", "%"+value+"%")
		i += 2
	}

	countArgs := append([]interface{}{}, args...)

	allowedOrderColumns := []string{"id", "name", "description", "threshold", "created_at", "updated_at"}
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
		err := r.sqlDB.SelectContext(ctx, &quizes, query, args...)
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

	return quizes, total, nil
}

func (r *QuizRepository) GetQuizByID(ctx context.Context, params *dtos.GetQuizParams) (*dtos.QuizDetailDTO, error) {
	var quiz dtos.QuizDetailDTO
	// deletedAt := params.IsDeleted

	query := `SELECT q.id, q.name, q.description, q.threshold, q.created_at, q.updated_at, q.deleted_at

	FROM quizes q
	WHERE 1=1`

	var args []interface{}

	i := 1
	query += " AND q.id = $1"
	args = append(args, params.ID)
	i++

	isDeletedQuery := ` AND q.deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND q.deleted_at IS NOT NULL"
	}

	query += isDeletedQuery

	if err := r.sqlDB.Get(&quiz, query, args...); err != nil {
		return nil, err
	}

	return &quiz, nil
}

// BeginTransaction starts a new transaction
func (r *QuizRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *QuizRepository) CreateQuiz(tx *gorm.DB, quiz *models.Quiz) error {
	if err := tx.Create(quiz).Error; err != nil {
		return err
	}
	return nil
}

func (r *QuizRepository) UpdateQuiz(tx *gorm.DB, quiz *models.Quiz) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(quiz).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *QuizRepository) DeleteQuiz(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Quiz{}, id).Error; err != nil {
		if err := tx.Delete(&models.Quiz{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *QuizRepository) RestoreQuiz(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE quizes SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
