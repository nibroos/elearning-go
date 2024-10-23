package controller

import (
	"context"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/models"
	pb "github.com/nibroos/elearning-go/service/internal/proto"
	"github.com/nibroos/elearning-go/service/internal/service"
	"google.golang.org/grpc"
)

// UserController holds the methods for handling user-related HTTP and gRPC requests.
type UserController struct {
	userService *service.UserService
	db          *sqlx.DB
}

// NewUserController creates a new instance of UserController.
func NewUserController(userService *service.UserService, db *sqlx.DB) *UserController {
	return &UserController{userService: userService, db: db}
}

// RunGRPCServer starts the gRPC server.
func RunGRPCServer(userService *service.UserService) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &GRPCUserController{userService: userService})
	log.Println("gRPC server is running on port 50051")

	return grpcServer.Serve(lis)
}

// GRPCUserController implements the gRPC user service.
type GRPCUserController struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

// CreateUser handles the gRPC request to create a user.
func (c *GRPCUserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Address:  req.Address,
	}

	createdUser, err := c.userService.CreateUser(ctx, user, req.GetRoleIds())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id:    uint32(createdUser.ID),
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}

// GetUsers handles the gRPC request to get users.
func (c *GRPCUserController) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	searchParams := map[string]string{
		"global": req.GetGlobal(),
		"name":   req.GetName(),
		"email":  req.GetEmail(),
	}

	users, total, err := c.userService.GetUsers(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:       uint32(user.ID),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return &pb.GetUsersResponse{Users: pbUsers, Total: int32(total)}, nil
}

// GetUserByID handles the gRPC request to get a user by ID.
func (c *GRPCUserController) GetUserByID(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := c.userService.GetUserByID(ctx, uint32(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:            uint32(user.ID),
			Name:          user.Name,
			Email:         user.Email,
			Address:       user.Address,
			RoleIds:       user.RoleIDs,
			PermissionIds: user.PermissionIDs,
		},
	}, nil
}
