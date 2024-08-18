package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/users-service/internal/model"
)

type UserRepository interface {
    CreateUser(user *model.User) error
    GetUserByID(id int) (*model.User, error)
    UpdateUser(user *model.User) error
    DeleteUser(id int) error
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
    _, err := r.db.NamedExec(`INSERT INTO users (name, email, password) VALUES (:name, :email, :password)`, user)
    return err
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
    var user model.User
    err := r.db.Get(&user, `SELECT * FROM users WHERE id = $1`, id)
    return &user, err
}

func (r *userRepository) UpdateUser(user *model.User) error {
    _, err := r.db.NamedExec(`UPDATE users SET name=:name, email=:email, password=:password WHERE id=:id`, user)
    return err
}

func (r *userRepository) DeleteUser(id int) error {
    _, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
    return err
}
