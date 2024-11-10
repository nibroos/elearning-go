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

type ContactRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewContactRepository(db *gorm.DB, sqlDB *sqlx.DB) *ContactRepository {
	return &ContactRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *ContactRepository) ListContacts(ctx context.Context, filters map[string]string) ([]dtos.ContactListDTO, int, error) {
	contacts := []dtos.ContactListDTO{}
	var total int

	from := `FROM (
        SELECT c.id, c.user_id, c.type_contact_id, c.ref_num, c.status, c.created_at, c.updated_at,
        u.name as user_name,
        ti.name as type_contact_name

        FROM contacts c
        JOIN users u ON c.user_id = u.id
        JOIN mix_values ti ON c.type_contact_id = ti.id
        WHERE c.deleted_at IS NULL
    ) AS alias WHERE 1=1`

	query := `SELECT * ` + from
	countQuery := `SELECT COUNT(*) ` + from

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

	if value, ok := filters["user_id"]; ok && value != "" {
		query += fmt.Sprintf(" AND user_id = $%d", i)
		countQuery += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, value)
		i++
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (ref_num ILIKE $%d OR user_name ILIKE $%d OR type_contact_name ILIKE $%d)", i, i+1, i+2)
		countQuery += fmt.Sprintf(" AND (ref_num ILIKE $%d OR user_name ILIKE $%d OR type_contact_name ILIKE $%d)", i, i+1, i+2)
		args = append(args, "%"+value+"%", "%"+value+"%", "%"+value+"%")
		i += 3
	}

	countArgs := append([]interface{}{}, args...)

	allowedOrderColumns := []string{"id", "ref_num", "user_name", "type_contact_name"}
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
		err := r.sqlDB.SelectContext(ctx, &contacts, query, args...)
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

	return contacts, total, nil
}

func (r *ContactRepository) GetContactByID(ctx context.Context, params *dtos.GetContactParams) (*dtos.ContactDetailDTO, error) {
	var contact dtos.ContactDetailDTO
	// deletedAt := params.IsDeleted

	query := `SELECT c.id, c.user_id, c.type_contact_id, c.ref_num, c.status, c.created_at, c.updated_at, c.deleted_at,
	u.name as user_name,
	ti.name as type_contact_name

	FROM contacts c
	JOIN users u ON c.user_id = u.id
	JOIN mix_values ti ON c.type_contact_id = ti.id
	WHERE 1=1`

	var args []interface{}

	i := 1
	query += " AND c.id = $1"
	args = append(args, params.ID)
	i++

	isDeletedQuery := ` AND c.deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND c.deleted_at IS NOT NULL"
	}

	if params.UserID != 0 {
		query += fmt.Sprintf(" AND c.user_id = $%d", i)
		args = append(args, params.UserID)
		i++
	}

	query += isDeletedQuery

	if err := r.sqlDB.Get(&contact, query, args...); err != nil {
		return nil, err
	}

	return &contact, nil
}

// BeginTransaction starts a new transaction
func (r *ContactRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *ContactRepository) CreateContact(tx *gorm.DB, contact *models.Contact) error {
	if err := tx.Create(contact).Error; err != nil {
		return err
	}
	return nil
}

func (r *ContactRepository) UpdateContact(tx *gorm.DB, contact *models.Contact) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(contact).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *ContactRepository) DeleteContact(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Contact{}, id).Error; err != nil {
		if err := tx.Delete(&models.Contact{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ContactRepository) RestoreContact(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE contacts SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
