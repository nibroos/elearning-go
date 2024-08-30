package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/users-service/internal/model"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
)

// type UserService interface {
//     GetUsers(searchParams map[string]string) ([]model.User, error)
// }

// type UserService struct {
//     repo repository.UsersRepository
//     // pb.UnimplementedUserServiceServer
// }

type UserService struct {
    repo *repository.UsersRepository
}


func NewUserService(repo *repository.UsersRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUsers(searchParams map[string]string) ([]model.User, error) {
    return s.repo.GetUsers(searchParams)
}
func (s *UserService) CreateUser(ctx context.Context, tx *sqlx.Tx, name string, email string, password string, roleID int64) (*model.User, error) {
    user := &model.User{Name: name, Email: email, Password: password}

    // Creating user in a goroutine to handle concurrency
    ch := make(chan error)
    go func() {
        id, err := s.repo.CreateUser(ctx, tx, user)
        if err != nil {
            ch <- err
            return
        }
        user.ID = id
        ch <- nil
    }()

    // Wait for the user to be created
    if err := <-ch; err != nil {
        return nil, err
    }

    // Attach role to the user concurrently
    go func() {
        ch <- s.repo.AttachRoleToUser(ctx, tx, user.ID, roleID)
    }()

    if err := <-ch; err != nil {
        return nil, err
    }

    return user, nil
}
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
//         Id:             int32(user.ID),
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
//         Id:             int32(user.ID),
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
//         Id:             int32(user.ID),
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
