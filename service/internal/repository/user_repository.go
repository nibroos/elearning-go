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

type UserRepository interface {
	GetUsers(ctx context.Context, filters map[string]string) ([]dtos.UserListDTO, int, error)
	GetUserByID(ctx context.Context, params *dtos.GetUserByIDParams) (*dtos.UserDetailDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*dtos.UserDetailDTO, error)
	BeginTransaction() *gorm.DB
	AttachRoles(tx *gorm.DB, user *models.User, roleIDs []uint32) error
	CreateUser(tx *gorm.DB, user *models.User) error
	UpdateUser(tx *gorm.DB, user *models.User) error
	DeleteUser(tx *gorm.DB, id uint) error
	DeleteRolesByUserID(tx *gorm.DB, userID uint) error
	RestoreUser(tx *gorm.DB, id uint) error
	Commit(tx *gorm.DB) error
}

type userRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewUserRepository(db *gorm.DB, sqlDB *sqlx.DB) *userRepository {
	return &userRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *userRepository) GetUsers(ctx context.Context, filters map[string]string) ([]dtos.UserListDTO, int, error) {
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

	orderColumn := utils.GetStringOrDefault(filters["order_column"], "id")
	orderDirection := utils.GetStringOrDefault(filters["order_direction"], "asc")
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)

	countArgs := append([]interface{}{}, args...)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)

	args = append(args, perPage, (currentPage-1)*perPage)

	countChan := make(chan error)
	selectChan := make(chan error)

	// Goroutine for count query
	go func() {
		err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
		countChan <- err
	}()

	// Goroutine for select query
	go func() {
		err := r.sqlDB.SelectContext(ctx, &users, query, args...)
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

	return users, total, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, params *dtos.GetUserByIDParams) (*dtos.UserDetailDTO, error) {
	var user dtos.UserDetailDTO

	query := `SELECT id, username, name, email, address, password FROM users WHERE id = $1`

	var args []interface{}
	args = append(args, params.ID)

	isDeletedQuery := ` AND deleted_at IS NULL`
	if params.IsDeleted != nil && *params.IsDeleted == 1 {
		isDeletedQuery = " AND deleted_at IS NOT NULL"
	}

	query += isDeletedQuery

	// Channels for concurrent execution
	userChan := make(chan error)
	roleChan := make(chan error)
	permissionChan := make(chan error)

	// Goroutine for user query
	go func() {
		err := r.sqlDB.GetContext(ctx, &user, query, args...)
		userChan <- err
	}()

	// Goroutine for role query
	go func() {
		var roleNames []string
		roleQuery := `
            SELECT mv.name 
            FROM pools p
            JOIN mix_values mv ON p.mv2_id = mv.id
            JOIN groups g1 ON p.group1_id = g1.id
            JOIN groups g2 ON p.group2_id = g2.id
            WHERE p.deleted_at IS NULL AND
						g1.name = 'users' AND g2.name = 'roles' 
            AND p.deleted_at IS NULL
            AND p.mv1_id = $1
        `
		err := r.sqlDB.SelectContext(ctx, &roleNames, roleQuery, params.ID)
		if err == nil {
			user.Roles = roleNames
		}
		roleChan <- err
	}()

	// Goroutine for permission query
	go func() {
		var permissionNames []string
		permissionQuery := `
            SELECT mv.name 
            FROM pools p
            JOIN mix_values mv ON p.mv2_id = mv.id
            JOIN groups g1 ON p.group1_id = g1.id
            JOIN groups g2 ON p.group2_id = g2.id
            WHERE p.deleted_at IS NULL AND
						g1.name = 'roles' AND g2.name = 'permissions' AND p.mv1_id IN (
                SELECT mv.id 
                FROM pools p
                JOIN mix_values mv ON p.mv2_id = mv.id
                JOIN groups g1 ON p.group1_id = g1.id
                JOIN groups g2 ON p.group2_id = g2.id
                WHERE g1.name = 'users' AND g2.name = 'roles' AND p.mv1_id = $1
            )
        `
		err := r.sqlDB.SelectContext(ctx, &permissionNames, permissionQuery, params.ID)
		if err == nil {
			user.Permissions = permissionNames
		}
		permissionChan <- err
	}()

	// Wait for all goroutines to finish
	userErr := <-userChan
	roleErr := <-roleChan
	permissionErr := <-permissionChan

	if userErr != nil {
		return nil, userErr
	}

	if roleErr != nil {
		return nil, roleErr
	}

	if permissionErr != nil {
		return nil, permissionErr
	}

	return &user, nil
}
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*dtos.UserDetailDTO, error) {
	var user dtos.UserDetailDTO

	query := `SELECT id, username, name, email, password, address FROM users WHERE deleted_at IS NULL AND (email = $1 OR username = $1)`
	if err := r.sqlDB.GetContext(ctx, &user, query, email); err != nil {
		return nil, err
	}

	id := user.ID

	// Channels for concurrent execution
	roleChan := make(chan error)
	permissionChan := make(chan error)

	// Goroutine for role query
	go func() {
		var roleNames []string
		roleQuery := `
            SELECT mv.name
            FROM pools p
            JOIN mix_values mv ON p.mv2_id = mv.id
            JOIN groups g1 ON p.group1_id = g1.id
            JOIN groups g2 ON p.group2_id = g2.id
            WHERE p.deleted_at IS NULL AND
						g1.name = 'users' AND g2.name = 'roles' AND p.mv1_id = $1
        `
		err := r.sqlDB.SelectContext(ctx, &roleNames, roleQuery, id)
		if err == nil {
			user.Roles = roleNames
		}
		roleChan <- err
	}()

	// Goroutine for permission query
	go func() {
		var permissionNames []string
		permissionQuery := `
            SELECT mv.name
            FROM pools p
            JOIN mix_values mv ON p.mv2_id = mv.id
            JOIN groups g1 ON p.group1_id = g1.id
            JOIN groups g2 ON p.group2_id = g2.id
            WHERE p.deleted_at IS NULL AND
						g1.name = 'roles' AND g2.name = 'permissions' AND p.mv1_id IN (
                SELECT mv.id
                FROM pools p
                JOIN mix_values mv ON p.mv2_id = mv.id
                JOIN groups g1 ON p.group1_id = g1.id
                JOIN groups g2 ON p.group2_id = g2.id
                WHERE g1.name = 'users' AND g2.name = 'roles' AND p.mv1_id = $1
            )
        `
		err := r.sqlDB.SelectContext(ctx, &permissionNames, permissionQuery, id)
		if err == nil {
			user.Permissions = permissionNames
		}
		permissionChan <- err
	}()

	// Wait for both goroutines to finish
	roleErr := <-roleChan
	permissionErr := <-permissionChan

	if roleErr != nil {
		return nil, roleErr
	}

	if permissionErr != nil {
		return nil, permissionErr
	}

	return &user, nil
}

// BeginTransaction starts a new transaction
func (r *userRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *userRepository) AttachRoles(tx *gorm.DB, user *models.User, roleIDs []uint32) error {
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

func (r *userRepository) CreateUser(tx *gorm.DB, user *models.User) error {
	if err := tx.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateUser(tx *gorm.DB, user *models.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *userRepository) DeleteUser(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		if err := tx.Delete(&models.User{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteRolesByUserID
func (r *userRepository) DeleteRolesByUserID(tx *gorm.DB, userID uint) error {
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

func (s *userRepository) RestoreUser(tx *gorm.DB, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET deleted_at = NULL WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

// commit or rollback
func (r *userRepository) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}
