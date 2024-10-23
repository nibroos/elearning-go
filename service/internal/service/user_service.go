package service

import (
	"context"
	"errors"

	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/repository"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// type UserService interface {
//     GetUsers(searchParams map[string]string) ([]models.User, error)
// }

// type UserService struct {
//     repo repository.UserRepository
//     // pb.UnimplementedUserServiceServer
// }

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(ctx context.Context, filters map[string]string) ([]dtos.UserListDTO, int, error) {

	resultChan := make(chan dtos.GetUsersResult, 1)

	go func() {
		users, total, err := s.repo.GetUsers(ctx, filters)
		resultChan <- dtos.GetUsersResult{Users: users, Total: total, Err: err}
	}()

	select {
	case res := <-resultChan:
		return res.Users, res.Total, res.Err
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
}
func (s *UserService) CreateUser(ctx context.Context, user *models.User, roleIDs []uint32) (*models.User, error) {
	// Hash password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Create user
	if err := s.repo.CreateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Attach roles
	if err := s.repo.AttachRoles(tx, user, roleIDs); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*dtos.UserDetailDTO, error) {
	userChan := make(chan *dtos.UserDetailDTO, 1)
	errChan := make(chan error, 1)

	go func() {
		user, err := s.repo.GetUserByID(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		userChan <- user
	}()

	select {
	case user := <-userChan:
		return user, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User, roleIDs []uint32) (*models.User, error) {
	// TODO update hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// TODO update user
	user.Password = hashedPassword

	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// Update user
	if err := s.repo.UpdateUser(tx, user); err != nil {
		tx.Rollback()
		return nil, err
	}

	// TODO update roles
	// Attach roles
	if err := s.repo.AttachRoles(tx, user, roleIDs); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*dtos.UserDetailDTO, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password.Value), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// delete roles
	if err := s.repo.DeleteRolesByUserID(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	// Delete user
	if err := s.repo.DeleteUser(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) RestoreUser(ctx context.Context, id uint) error {
	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Restore user
	if err := s.repo.RestoreUser(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
