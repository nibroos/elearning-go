package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/users-service/internal/model"
	"github.com/pkg/errors"
)

// type UsersRepository interface {
//     CreateUser(user *model.User) error
//     GetUserByID(id int) (*model.User, error)
//     UpdateUser(user *model.User) error
//     DeleteUser(id int) error
//     GetUsers(searchParams map[string]string) ([]model.User, error)
// }

type UsersRepository struct {
    db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
    return &UsersRepository{db: db}
}

func (r *UsersRepository) GetUsers(searchParams map[string]string) ([]model.User, error) {
    var users []model.User
    query := `SELECT * FROM users WHERE 1=1`
    args := []interface{}{}

    if global, exists := searchParams["global"]; exists && global != "" {
        query += ` AND (username ILIKE $1 OR email ILIKE $2)`
        args = append(args, "%"+global+"%", "%"+global+"%")
    }

    if name, exists := searchParams["name"]; exists && name != "" {
        query += ` AND username ILIKE $3`
        query += ` AND name ILIKE $4`
        args = append(args, "%"+name+"%")
        args = append(args, "%"+name+"%")
    }

    if email, exists := searchParams["email"]; exists && email != "" {
        query += ` AND email ILIKE $5`
        args = append(args, "%"+email+"%")
    }

    query = strings.Replace(query, "$1", fmt.Sprintf("$%d", len(args)-2), -1)
    query = strings.Replace(query, "$2", fmt.Sprintf("$%d", len(args)-1), -1)
    query = strings.Replace(query, "$3", fmt.Sprintf("$%d", len(args)), -1)
    query = strings.Replace(query, "$4", fmt.Sprintf("$%d", len(args)+1), -1)
    query = strings.Replace(query, "$5", fmt.Sprintf("$%d", len(args)+2), -1)

    err := r.db.Select(&users, query, args...)
    if err != nil {
        return nil, err
    }

    return users, nil
}

// func (r *UsersRepository) CreateUser(user *model.User) error {
//     _, err := r.db.NamedExec(`INSERT INTO users (name, email, password) VALUES (:name, :email, :password)`, user)
//     return err
// }

func (r *UsersRepository) CreateUser(ctx context.Context, tx *sqlx.Tx, user *model.User) (int64, error) {
    var id int64
    err := tx.QueryRowContext(ctx, `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`, user.Name, user.Email).Scan(&id)
    if err != nil {
        return 0, errors.Wrap(err, "failed to insert user")
    }
    return id, nil
}

func (r *UsersRepository) AttachRoleToUser(ctx context.Context, tx *sqlx.Tx, userID, roleID int64) error {
    _, err := tx.ExecContext(ctx, `INSERT INTO role_user (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
    if err != nil {
        return errors.Wrap(err, "failed to attach role to user")
    }
    return nil
}

// func (r *UsersRepository) GetUserByID(id int) (*model.User, error) {
//     var user model.User
//     err := r.db.Get(&user, `SELECT * FROM users WHERE id = $1`, id)
//     return &user, err
// }

// func (r *UsersRepository) UpdateUser(user *model.User) error {
//     _, err := r.db.NamedExec(`UPDATE users SET name=:name, email=:email, password=:password WHERE id=:id`, user)
//     return err
// }

// func (r *UsersRepository) DeleteUser(id int) error {
//     _, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
//     return err
// }
