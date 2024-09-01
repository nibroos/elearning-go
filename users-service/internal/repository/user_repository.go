package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/models"
	"gorm.io/gorm"
)

// type UserRepository interface {
//     CreateUser(user *model.User) error
//     GetUserByID(id int) (*model.User, error)
//     UpdateUser(user *model.User) error
//     DeleteUser(id int) error
//     GetUsers(searchParams map[string]string) ([]models.User, error)
// }

type UserRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewUserRepository(db *gorm.DB, sqlDB *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}
func (r *UserRepository) GetUsers(ctx context.Context, filters map[string]string) ([]dtos.UserListDTO, int, error) {
	users := []dtos.UserListDTO{}
	var total int

	query := `SELECT * FROM users WHERE 1=1`
	var args []interface{}

	i := 1
	for key, value := range filters {
		switch key {
		case "username", "name", "email", "global":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	countArgs := append([]interface{}{}, args...)
	err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	orderColumn := filters["order_column"]
	orderDirection := filters["order_direction"]
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage, _ := strconv.Atoi(filters["per_page"])
	currentPage, _ := strconv.Atoi(filters["page"])
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, (currentPage-1)*perPage)

	err = r.sqlDB.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetUserByID(id uint) (*dtos.UserDetailDTO, error) {
	var user dtos.UserDetailDTO
	query := `SELECT * FROM users WHERE id = $1`
	if err := r.sqlDB.Get(&user, query, id); err != nil {
		return nil, err
	}
	// TODO: add roles
	// TODO: add permissions

	return &user, nil
}

// BeginTransaction starts a new transaction
func (r *UserRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

//	func (r *UserRepository) CreateUser(user *model.User) error {
//	    _, err := r.db.NamedExec(`INSERT INTO users (name, email, password) VALUES (:name, :email, :password)`, user)
//	    return err
//	}

func (r *UserRepository) AttachRoles(tx *gorm.DB, user *models.User, roleIDs []uint) error {
	// Prepare batch insert for new role_user relationships
	var pools []models.Pool
	for _, roleID := range roleIDs {
		pool := models.Pool{
			Name:  "role_user",
			Tb1ID: user.ID,
			Tb2ID: roleID,
		}
		pools = append(pools, pool)
	}

	// Insert all role_user relationships in a single query
	if len(pools) > 0 {
		if err := tx.Create(&pools).Error; err != nil {
			return err
		}
	}

	return nil
}

// func (r *UserRepository) AttachRoles(user *models.User, roleIDs []int32) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		if len(roleIDs) > 0 {
// 			roles := []models.Role{}
// 			if err := tx.Where("id IN (?)", roleIDs).Find(&roles).Error; err != nil {
// 				return err
// 			}
// 			if err := tx.Model(user).Association("Roles").Replace(roles); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

// func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
//     var user model.User
//     err := r.db.Get(&user, `SELECT * FROM users WHERE id = $1`, id)
//     return &user, err
// }

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *UserRepository) DeleteUser(user *models.User) error {
	return r.db.Delete(user).Error
}

// func (r *UserRepository) GetUsers(searchParams map[string]string) ([]dtos.UserListDTO, error) {
//     // var users []models.User
//     var users []dtos.UserListDTO

//     query := `SELECT * FROM users WHERE 1=1`
//     args := []interface{}{}

//     if global, exists := searchParams["global"]; exists && global != "" {
//         query += ` AND (username ILIKE $1 OR email ILIKE $2)`
//         args = append(args, "%"+global+"%", "%"+global+"%")
//     }

//     if name, exists := searchParams["name"]; exists && name != "" {
//         query += ` AND username ILIKE $3`
//         query += ` AND name ILIKE $4`
//         args = append(args, "%"+name+"%")
//         args = append(args, "%"+name+"%")
//     }

//     if email, exists := searchParams["email"]; exists && email != "" {
//         query += ` AND email ILIKE $5`
//         args = append(args, "%"+email+"%")
//     }

//     query = strings.Replace(query, "$1", fmt.Sprintf("$%d", len(args)-2), -1)
//     query = strings.Replace(query, "$2", fmt.Sprintf("$%d", len(args)-1), -1)
//     query = strings.Replace(query, "$3", fmt.Sprintf("$%d", len(args)), -1)
//     query = strings.Replace(query, "$4", fmt.Sprintf("$%d", len(args)+1), -1)
//     query = strings.Replace(query, "$5", fmt.Sprintf("$%d", len(args)+2), -1)

//     err := r.sqlDB.Select(&users, query, args...)
//     if err != nil {
//         return nil, err
//     }

//     return users, nil
// }
// func (r *UserRepository) UpdateUser(user *model.User) error {
//     _, err := r.db.NamedExec(`UPDATE users SET name=:name, email=:email, password=:password WHERE id=:id`, user)
//     return err
// }

// func (r *UserRepository) DeleteUser(id int) error {
//     _, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
//     return err
// }
