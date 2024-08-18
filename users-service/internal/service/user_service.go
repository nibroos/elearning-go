package service

import (
	"context"

	"github.com/your-username/users-service/internal/model"
	"github.com/your-username/users-service/internal/repository"
	pb "github.com/your-username/users-service/proto"
)

type UserService struct {
    repo repository.UserRepository
    pb.UnimplementedUserServiceServer
}

func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
    user := &model.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }

    err := s.repo.CreateUser(user)
    if err != nil {
        return nil, err
    }

    return &pb.UserResponse{User: &pb.User{
        Id:       int32(user.ID),
        Name:     user.Name,
        Email:    user.Email,
        Password: user.Password,
    }}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    user, err := s.repo.GetUserByID(int(req.Id))
    if err != nil {
        return nil, err
    }

    return &pb.UserResponse{User: &pb.User{
        Id:       int32(user.ID),
        Name:     user.Name,
        Email:    user.Email,
        Password: user.Password,
    }}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
    user := &model.User{
        ID:       int(req.Id),
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }

    err := s.repo.UpdateUser(user)
    if err != nil {
        return nil, err
    }

    return &pb.UserResponse{User: &pb.User{
        Id:       int32(user.ID),
        Name:     user.Name,
        Email:    user.Email,
        Password: user.Password,
    }}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    err := s.repo.DeleteUser(int(req.Id))
    if err != nil {
        return nil, err
    }

    return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
