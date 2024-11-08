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

	query := `SELECT id, username, name, email FROM users WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	var args []interface{}

	i := 1
	for key, value := range filters {
		switch key {
		case "username", "name", "email":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				countQuery += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (username ILIKE $%d OR name ILIKE $%d OR email ILIKE $%d)", i, i+1, i+2)
		countQuery += fmt.Sprintf(" AND (username ILIKE $%d OR name ILIKE $%d OR email ILIKE $%d)", i, i+1, i+2)
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

	err = r.sqlDB.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
func (r *UserRepository) GetUserByID(ctx context.Context, params *dtos.GetUserByIDParams) (*dtos.UserDetailDTO, error) {
	var user dtos.UserDetailDTO

	query := `SELECT id, username, name, email, address, password FROM users WHERE id = $1`

	var args []interface{}

	i := 1
	query += " AND id = $1"
	args = append(args, params.ID)
	i++

	isDeletedQuery := ` AND deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND deleted_at IS NOT NULL"
		i++
	}

	query += isDeletedQuery

	if err := r.sqlDB.Get(&user, query, args...); err != nil {
		return nil, err
	}

	// utils.DD(ctx, map[string]interface{}{
	// 	"user": user,
	// })

	// Fetch roles
	var roleNames []string
	roleQuery := `
        SELECT mv.name 
        FROM pools p
        JOIN mix_values mv ON p.mv2_id = mv.id
        JOIN groups g1 ON p.group1_id = g1.id
        JOIN groups g2 ON p.group2_id = g2.id
        WHERE g1.name = 'users' AND g2.name = 'roles' 
				AND p.deleted_at IS NULL
				AND p.mv1_id = $1
    `
	if err := r.sqlDB.Select(&roleNames, roleQuery, params.ID); err != nil {
		return nil, err
	}
	user.Roles = roleNames

	// Fetch permissions
	var permissionNames []string
	permissionQuery := `
        SELECT mv.name 
        FROM pools p
        JOIN mix_values mv ON p.mv2_id = mv.id
        JOIN groups g1 ON p.group1_id = g1.id
        JOIN groups g2 ON p.group2_id = g2.id
        WHERE g1.name = 'roles' AND g2.name = 'permissions' AND p.mv1_id IN (
            SELECT mv.id 
            FROM pools p
            JOIN mix_values mv ON p.mv2_id = mv.id
            JOIN groups g1 ON p.group1_id = g1.id
            JOIN groups g2 ON p.group2_id = g2.id
            WHERE g1.name = 'users' AND g2.name = 'roles' AND p.mv1_id = $1
        )
    `
	if err := r.sqlDB.Select(&permissionNames, permissionQuery, params.ID); err != nil {
		return nil, err
	}
	user.Permissions = permissionNames

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*dtos.UserDetailDTO, error) {
	var user dtos.UserDetailDTO

	query := `SELECT id, username, name, email, password, address FROM users WHERE deleted_at IS NULL AND email = $1 OR username = $1`
	if err := r.sqlDB.Get(&user, query, email); err != nil {
		return nil, err
	}

	id := user.ID

	// Fetch roles
	var roleNames []string
	roleQuery := `
				SELECT mv.name
				FROM pools p
				JOIN mix_values mv ON p.mv2_id = mv.id
				JOIN groups g1 ON p.group1_id = g1.id
				JOIN groups g2 ON p.group2_id = g2.id
				WHERE g1.name = 'users' AND g2.name = 'roles' AND p.mv1_id = $1
			`
	if err := r.sqlDB.Select(&roleNames, roleQuery, id); err != nil {
		return nil, err
	}

	user.Roles = roleNames

	// Fetch permissions
	var permissionNames []string
	permissionQuery := `
				SELECT mv.name
				FROM pools p
				JOIN mix_values mv ON p.mv2_id = mv.id
				JOIN groups g1 ON p.group1_id = g1.id
				JOIN groups g2 ON p.group2_id = g2.id
				WHERE g1.name = 'roles' AND g2.name = 'permissions' AND p.mv1_id IN (
					SELECT mv.id
					FROM pools p
					JOIN mix_values mv ON p.mv2_id = mv.id
					JOIN groups g1 ON p.group1_id = g1.id
					JOIN groups g2 ON p.group2_id = g2.id
					WHERE g1.name = 'users' AND g2.name = 'roles' AND p.mv1_id = $1
				)
			`
	if err := r.sqlDB.Select(&permissionNames, permissionQuery, id); err != nil {
		return nil, err
	}

	user.Permissions = permissionNames

	return &user, nil
}

// func (r *UserRepository) GetUserByID(ctx context.Context, id uint32) (*dtos.UserDetailDTO, error) {
// 	var user dtos.UserDetailDTO

// 	query := `SELECT * FROM users WHERE id = $1`
// 	if err := r.sqlDB.Get(&user, query, id); err != nil {
// 		return nil, err
// 	}

// 	utils.DD(ctx, map[string]interface{}{
// 		"user": user,
// 	})

// 	// Fetch roles
// 	var roleIDs []uint32
// 	roleQuery := `SELECT role_id FROM role_user WHERE user_id = $1`
// 	if err := r.sqlDB.Select(&roleIDs, roleQuery, id); err != nil {
// 		return nil, err
// 	}
// 	user.RoleIDs = roleIDs

// 	// Fetch permissions
// 	var permissionIDs []uint32
// 	permissionQuery := `SELECT permission_id FROM permission_user WHERE user_id = $1`
// 	if err := r.sqlDB.Select(&permissionIDs, permissionQuery, id); err != nil {
// 		return nil, err
// 	}
// 	user.PermissionIDs = permissionIDs

// 	return &user, nil
// }

// BeginTransaction starts a new transaction
func (r *UserRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

//	func (r *UserRepository) CreateUser(user *model.User) error {
//	    _, err := r.db.NamedExec(`INSERT INTO users (name, email, password) VALUES (:name, :email, :password)`, user)
//	    return err
//	}

func (r *UserRepository) AttachRoles(tx *gorm.DB, user *models.User, roleIDs []uint32) error {
	// Prepare batch insert for new role_user relationships
	var pools []models.Pool
	for _, roleID := range roleIDs {
		pool := models.Pool{
			Group1ID: utils.GroupIDUsers, // users
			Group2ID: utils.GroupIDRoles, // roles
			Mv1ID:    uint32(user.ID),
			Mv2ID:    roleID,
		}
		pools = append(pools, pool)
	}

	// delete existing roles
	if err := r.DeleteRolesByUserID(tx, user.ID); err != nil {
		return err
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

func (r *UserRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	if err := tx.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUser(tx *gorm.DB, user *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *UserRepository) DeleteUser(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		if err := tx.Delete(&models.User{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteRolesByUserID
func (r *UserRepository) DeleteRolesByUserID(tx *gorm.DB, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			UPDATE pools SET deleted_at = NOW() 
			WHERE group1_id = ? AND mv1_id = ?
			AND group2_id = ?
		`, utils.GroupIDUsers, userID, utils.GroupIDRoles).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *UserRepository) RestoreUser(tx *gorm.DB, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
