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

type SubscribeRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewSubscribeRepository(db *gorm.DB, sqlDB *sqlx.DB) *SubscribeRepository {
	return &SubscribeRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *SubscribeRepository) GetSubscribes(ctx context.Context, filters map[string]string) ([]dtos.SubscribeListDTO, int, error) {
	classes := []dtos.SubscribeListDTO{}
	var total int

	query := `SELECT id, name, description FROM classes WHERE 1=1`
	var args []interface{}

	i := 1
	for key, value := range filters {
		switch key {
		case "name", "global":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	countQuery := `SELECT COUNT(*) FROM classes WHERE 1=1`
	countArgs := append([]interface{}{}, args...)
	err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	orderColumn := utils.GetStringOrDefault(filters["order_column"], "id")
	orderDirection := utils.GetStringOrDefault(filters["order_direction"], "asc")
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, (currentPage-1)*perPage)

	err = r.sqlDB.SelectContext(ctx, &classes, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}
func (r *SubscribeRepository) GetSubscribeByID(ctx context.Context, id uint) (*dtos.SubscribeDetailDTO, error) {
	var class dtos.SubscribeDetailDTO

	query := `SELECT id, name, description, FROM classes WHERE id = $1 AND deleted_at IS NULL`
	if err := r.sqlDB.Get(&class, query, id); err != nil {
		return nil, err
	}

	return &class, nil
}

// BeginTransaction starts a new transaction
func (r *SubscribeRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *SubscribeRepository) CreateSubscribe(tx *gorm.DB, class *models.Subscribe) error {
	if err := tx.Create(class).Error; err != nil {
		return err
	}
	return nil
}

func (r *SubscribeRepository) UpdateSubscribe(tx *gorm.DB, class *models.Subscribe) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(class).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *SubscribeRepository) DeleteSubscribe(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Subscribe{}, id).Error; err != nil {
		if err := tx.Delete(&models.Subscribe{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *SubscribeRepository) RestoreSubscribe(tx *gorm.DB, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var class models.Subscribe
		if err := tx.Unscoped().First(&class, id).Error; err != nil {
			return err
		}
		return tx.Model(&class).Update("deleted_at", nil).Error
	})
}
