package service

import (
	"context"
	"sync"

	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/models"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
	"github.com/nibroos/elearning-go/users-service/internal/utils"
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
	var wg sync.WaitGroup
	var users []dtos.UserListDTO
	var total int
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		users, total, err = s.repo.GetUsers(ctx, filters)
	}()

	wg.Wait()
	users, total, err = s.repo.GetUsers(ctx, filters)

	return users, total, err
}

func (s *UserService) CreateUser(user *models.User, roleIDs []uint) error {
	// Hash password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Transaction handling
	tx := s.repo.BeginTransaction()
	if err := tx.Error; err != nil {
		return err
	}

	// Create user
	if err := s.repo.CreateUser(user); err != nil {
		tx.Rollback()
		return err
	}

	// Attach roles
	if err := s.repo.AttachRoles(tx, user, roleIDs); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *UserService) GetUserByID(id uint) (*dtos.UserDetailDTO, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUser(user *models.User, roleIDs []uint) error {
	// TODO update hash password
	// TODO update user
	// TODO update roles

	return s.repo.UpdateUser(user)
}

// func (s *UserService) CreateUser(ctx context.Context, tx *sqlx.Tx, name string, email string, password string, roleID uint) (*model.User, error) {
//     user := &model.User{Name: name, Email: email, Password: password}

//     // Creating user in a goroutine to handle concurrency
//     ch := make(chan error)
//     go func() {
//         id, err := s.repo.CreateUser(ctx, tx, user)
//         if err != nil {
//             ch <- err
//             return
//         }
//         user.ID = id
//         ch <- nil
//     }()

//     // Wait for the user to be created
//     if err := <-ch; err != nil {
//         return nil, err
//     }

//     // Attach role to the user concurrently
//     go func() {
//         ch <- s.repo.AttachRoleToUser(ctx, tx, user.ID, roleID)
//     }()

//     if err := <-ch; err != nil {
//         return nil, err
//     }

//     return user, nil
// }
// func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
//     user := &model.User{
//         Name:           req.Name,
//         Email:          req.Email,
//         Password:       req.Password,
//     }

//     err := s.repo.CreateUser(user)
//     if err != nil {
//         return nil, err
//     }

//     return &pb.UserResponse{User: &pb.User{
//         Id:             uint(user.ID),
//         Name:           user.Name,
//         Email:          user.Email,
//         Password:       user.Password,
//     }}, nil
// }

// func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
//     user, err := s.repo.GetUserByID(int(req.Id))
//     if err != nil {
//         return nil, err
//     }

//     return &pb.UserResponse{User: &pb.User{
//         Id:             uint(user.ID),
//         Name:           user.Name,
//         Email:          user.Email,
//         Password:       user.Password,
//     }}, nil
// }

// func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
//     user := &model.User{
//         ID:             int(req.Id),
//         Name:           req.Name,
//         Email:          req.Email,
//         Password:       req.Password,
//     }

//     err := s.repo.UpdateUser(user)
//     if err != nil {
//         return nil, err
//     }

//     return &pb.UserResponse{User: &pb.User{
//         Id:             uint(user.ID),
//         Name:           user.Name,
//         Email:          user.Email,
//         Password:       user.Password,
//     }}, nil
// }

// func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
//     err := s.repo.DeleteUser(int(req.Id))
//     if err != nil {
//         return nil, err
//     }

//     return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
// }
