package service

import (
	"context"
	"errors"

	"github.com/nibroos/elearning-go/users-service/internal/dtos"
	"github.com/nibroos/elearning-go/users-service/internal/models"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
	"github.com/nibroos/elearning-go/users-service/internal/utils"
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

func (s *UserService) GetUserByID(ctx context.Context, id uint32) (*dtos.UserDetailDTO, error) {
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

func (s *UserService) UpdateUser(user *models.User, roleIDs []uint32) error {
	// TODO update hash password
	// TODO update user
	// TODO update roles

	return s.repo.UpdateUser(user)
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
