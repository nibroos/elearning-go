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

type EducationRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewEducationRepository(db *gorm.DB, sqlDB *sqlx.DB) *EducationRepository {
	return &EducationRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *EducationRepository) GetEducations(ctx context.Context, filters map[string]string) ([]dtos.EducationListDTO, int, error) {
	educations := []dtos.EducationListDTO{}
	var total int

	query := `SELECT *
	FROM ( 
		SELECT s.id, s.name, s.description, s.class_id, s.created_at, s.updated_at, s.deleted_at,
		c.name as class_name,
		cu.name as created_by_name,
		uu.name as updated_by_name

		FROM educations s
		JOIN users cu ON s.created_by_id = cu.id
		LEFT JOIN users uu ON s.updated_by_id = uu.id
		JOIN classes c ON s.class_id = c.id
	) AS alias WHERE 1=1 AND deleted_at IS NULL`

	countQuery := `SELECT COUNT(*) FROM (
		SELECT s.id, s.name, s.description, s.class_id, s.created_at, s.updated_at, s.deleted_at,
		c.name as class_name,
		cu.name as created_by_name,
		uu.name as updated_by_name

		FROM educations s
		JOIN users cu ON s.created_by_id = cu.id
		LEFT JOIN users uu ON s.updated_by_id = uu.id
		JOIN classes c ON s.class_id = c.id
	) AS alias WHERE 1=1 AND deleted_at IS NULL`

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
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR class_name ILIKE $%d)", i, i+1, i+2)
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR class_name ILIKE $%d)", i, i+1, i+2)
		args = append(args, "%"+value+"%", "%"+value+"%", "%"+value+"%")
		i += 3
	}

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

	err = r.sqlDB.SelectContext(ctx, &educations, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return educations, total, nil
}

func (r *EducationRepository) GetEducationByID(ctx context.Context, id uint) (*dtos.EducationDetailDTO, error) {
	var education dtos.EducationDetailDTO

	query := `SELECT s.*,
	c.name as class_name,
	cu.name as created_by_name,
	uu.name as updated_by_name

	FROM educations s
	JOIN users cu ON s.created_by_id = cu.id
	LEFT JOIN users uu ON s.updated_by_id = uu.id
	JOIN classes c ON s.class_id = c.id
	WHERE s.id = $1 AND s.deleted_at IS NULL`
	if err := r.sqlDB.Get(&education, query, id); err != nil {
		return nil, err
	}

	return &education, nil
}

// BeginTransaction starts a new transaction
func (r *EducationRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *EducationRepository) CreateEducation(tx *gorm.DB, education *models.Education) error {
	if err := tx.Create(education).Error; err != nil {
		return err
	}
	return nil
}

func (r *EducationRepository) UpdateEducation(tx *gorm.DB, education *models.Education) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(education).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *EducationRepository) DeleteEducation(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Education{}, id).Error; err != nil {
		if err := tx.Delete(&models.Education{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *EducationRepository) RestoreEducation(tx *gorm.DB, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var education models.Education
		if err := tx.Unscoped().First(&education, id).Error; err != nil {
			return err
		}
		return tx.Model(&education).Update("deleted_at", nil).Error
	})
}
