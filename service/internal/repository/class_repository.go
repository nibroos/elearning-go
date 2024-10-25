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

type ClassRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewClassRepository(db *gorm.DB, sqlDB *sqlx.DB) *ClassRepository {
	return &ClassRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

// subcribe_id INT REFERENCES subcribes(id),
// incharge_id INT REFERENCES users(id),
// name VARCHAR(255),
// description TEXT,
// banner_url TEXT,
// logo_url TEXT,
// video_url TEXT,
// created_by_id INT REFERENCES users(id),
// updated_by_id INT REFERENCES users(id),
// created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
// updated_at timestamp with time zone,
// deleted_at timestamp with time zone
func (r *ClassRepository) GetClasses(ctx context.Context, filters map[string]string) ([]dtos.ClassListDTO, int, error) {
	classes := []dtos.ClassListDTO{}
	var total int
	query := `SELECT 
        alias.id, 
        alias.name, 
        alias.description, 
        alias.banner_url, 
        alias.logo_url, 
        alias.video_url,
        alias.created_by_name,
        alias.updated_by_name,
        alias.incharge_name,
        alias.subcribe_name,
        alias.created_at,
        alias.updated_at
    FROM (
        SELECT 
            c.id AS id, 
            c.name AS name, 
            c.description AS description, 
            c.banner_url AS banner_url, 
            c.logo_url AS logo_url, 
            c.video_url AS video_url,
            cb.name AS created_by_name,
            ub.name AS updated_by_name,
            ic.name AS incharge_name,
            s.name AS subcribe_name,
            c.created_at AS created_at,
            c.updated_at AS updated_at,
            c.subcribe_id AS subcribe_id,
            c.incharge_id AS incharge_id
        FROM classes AS c
        JOIN users AS cb ON c.created_by_id = cb.id
        JOIN users AS ub ON c.updated_by_id = ub.id
        LEFT JOIN users AS ic ON c.incharge_id = ic.id
        JOIN subcribes AS s ON c.subcribe_id = s.id
    ) AS alias WHERE 1=1`
	var args []interface{}

	i := 1

	if value, ok := filters["name"]; ok {
		query += fmt.Sprintf(" AND name ILIKE $%d", i)
		args = append(args, "%"+value+"%")
		i++
	}

	if value, ok := filters["subcribe_id"]; ok {
		query += fmt.Sprintf(" AND subcribe_id = $%d", i)
		args = append(args, value)
		i++
	}

	if value, ok := filters["incharge_id"]; ok {
		query += fmt.Sprintf(" AND incharge_id = $%d", i)
		args = append(args, value)
		i++
	}

	// search by global multiple column using OR ILIKE
	if value, ok := filters["global"]; ok {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR created_by_name ILIKE $%d OR updated_by_name ILIKE $%d OR incharge_name ILIKE $%d OR subcribe_name ILIKE $%d)", i, i+1, i+2, i+3, i+4, i+5)
		args = append(args, "%"+value+"%", "%"+value+"%", "%"+value+"%", "%"+value+"%", "%"+value+"%", "%"+value+"%")
		i += 6
	}

	countQuery := `SELECT COUNT(*) FROM (
        SELECT 
            c.id AS id, 
            c.name AS name, 
            c.description AS description, 
            c.banner_url AS banner_url, 
            c.logo_url AS logo_url, 
            c.video_url AS video_url,
            cb.name AS created_by_name,
            ub.name AS updated_by_name,
            ic.name AS incharge_name,
            s.name AS subcribe_name,
            c.created_at AS created_at,
            c.updated_at AS updated_at,
            c.subcribe_id AS subcribe_id,
            c.incharge_id AS incharge_id
        FROM classes AS c
        JOIN users AS cb ON c.created_by_id = cb.id
        JOIN users AS ub ON c.updated_by_id = ub.id
        LEFT JOIN users AS ic ON c.incharge_id = ic.id
        JOIN subcribes AS s ON c.subcribe_id = s.id
    ) AS alias WHERE 1=1`
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
func (r *ClassRepository) GetClassByID(ctx context.Context, id uint) (*dtos.ClassDetailDTO, error) {
	var class dtos.ClassDetailDTO

	query := `SELECT 
        id, 
        name, 
        description, 
        banner_url, 
        logo_url, 
        video_url,
        created_by_name,
        updated_by_name,
        incharge_name,
        subcribe_name,
        created_at,
        updated_at,
        deleted_at
    FROM (
        SELECT 
            c.id AS id, 
            c.name AS name, 
            c.description AS description, 
            c.banner_url AS banner_url, 
            c.logo_url AS logo_url, 
            c.video_url AS video_url,
            cb.name AS created_by_name,
            ub.name AS updated_by_name,
            ic.name AS incharge_name,
            s.name AS subcribe_name,
            c.created_at AS created_at,
            c.updated_at AS updated_at,
            c.deleted_at AS deleted_at,
            c.subcribe_id AS subcribe_id,
            c.incharge_id AS incharge_id
        FROM classes AS c
        JOIN users AS cb ON c.created_by_id = cb.id
        JOIN users AS ub ON c.updated_by_id = ub.id
        LEFT JOIN users AS ic ON c.incharge_id = ic.id
        JOIN subcribes AS s ON c.subcribe_id = s.id
    ) AS alias WHERE id = $1 AND deleted_at IS NULL`
	if err := r.sqlDB.Get(&class, query, id); err != nil {
		return nil, err
	}

	return &class, nil
}

// BeginTransaction starts a new transaction
func (r *ClassRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *ClassRepository) CreateClass(tx *gorm.DB, class *models.Class) error {
	if err := tx.Create(class).Error; err != nil {
		return err
	}
	return nil
}

func (r *ClassRepository) UpdateClass(tx *gorm.DB, class *models.Class) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(class).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *ClassRepository) DeleteClass(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Class{}, id).Error; err != nil {
		if err := tx.Delete(&models.Class{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ClassRepository) RestoreClass(tx *gorm.DB, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var class models.Class
		if err := tx.Unscoped().First(&class, id).Error; err != nil {
			return err
		}
		return tx.Model(&class).Update("deleted_at", nil).Error
	})
}
