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

type IdentifierRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewIdentifierRepository(db *gorm.DB, sqlDB *sqlx.DB) *IdentifierRepository {
	return &IdentifierRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *IdentifierRepository) GetIdentifiers(ctx context.Context, filters map[string]string) ([]dtos.IdentifierListDTO, int, error) {
	identifers := []dtos.IdentifierListDTO{}
	var total int

	query := `SELECT *
	FROM ( 
		SELECT i.id, i.ref_num, i.status, i.created_at, i.updated_at, i.deleted_at,
		u.name as user_name,
		ti.name as type_identifier_name

		FROM identifers i
		JOIN users u ON i.user_id = u.id
		JOIN type_identifiers ti ON i.type_identifier_id = ti.id
	) AS alias WHERE 1=1 AND deleted_at IS NULL`

	countQuery := `SELECT COUNT(*) FROM (
		SELECT i.id, i.ref_num, i.status, i.created_at, i.updated_at, i.deleted_at,
		u.name as user_name,
		ti.name as type_identifier_name

		FROM identifers i
		JOIN users u ON i.user_id = u.id
		JOIN type_identifiers ti ON i.type_identifier_id = ti.id
	) AS alias WHERE 1=1 AND deleted_at IS NULL`

	var args []interface{}

	i := 1
	for key, value := range filters {
		switch key {
		case "ref_num":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				countQuery += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (ref_num ILIKE $%d OR user_name ILIKE $%d OR type_identifier_name ILIKE $%d)", i, i+1, i+2)
		countQuery += fmt.Sprintf(" AND (ref_num ILIKE $%d OR user_name ILIKE $%d OR type_identifier_name ILIKE $%d)", i, i+1, i+2)
		args = append(args, "%"+value+"%", "%"+value+"%", "%"+value+"%")
		i += 3
	}

	countArgs := append([]interface{}{}, args...)
	err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// orderColumn := utils.GetStringOrDefault(filters["order_column"], "id")
	allowedOrderColumns := []string{"id", "ref_num", "user_name", "type_identifier_name"}
	orderColumn := utils.GetStringOrDefaultFromArray(filters["order_column"], allowedOrderColumns, "id")
	orderDirection := utils.GetStringOrDefault(filters["order_direction"], "asc")
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, (currentPage-1)*perPage)

	err = r.sqlDB.SelectContext(ctx, &identifers, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return identifers, total, nil
}

func (r *IdentifierRepository) GetIdentifierByID(ctx context.Context, params *dtos.GetIdentifierParams) (*dtos.IdentifierDetailDTO, error) {
	var identifer dtos.IdentifierDetailDTO
	// deletedAt := params.IsDeleted

	query := `SELECT i.id, i.ref_num, i.status, i.created_at, i.updated_at, i.deleted_at,
	u.name as user_name,
	ti.name as type_identifier_name

	FROM identifers i
	JOIN users u ON i.user_id = u.id
	JOIN type_identifiers ti ON i.type_identifier_id = ti.id
	WHERE 1=1`

	var args []interface{}

	i := 1
	query += " AND e.id = $1"
	args = append(args, params.ID)
	i++

	isDeletedQuery := ` AND i.deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND i.deleted_at IS NOT NULL"
		i++
	}

	query += isDeletedQuery

	if err := r.sqlDB.Get(&identifer, query, args...); err != nil {
		return nil, err
	}

	return &identifer, nil
}

// BeginTransaction starts a new transaction
func (r *IdentifierRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *IdentifierRepository) CreateIdentifier(tx *gorm.DB, identifer *models.Identifier) error {
	if err := tx.Create(identifer).Error; err != nil {
		return err
	}
	return nil
}

func (r *IdentifierRepository) UpdateIdentifier(tx *gorm.DB, identifer *models.Identifier) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(identifer).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *IdentifierRepository) DeleteIdentifier(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Identifier{}, id).Error; err != nil {
		if err := tx.Delete(&models.Identifier{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *IdentifierRepository) RestoreIdentifier(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE identifers SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
