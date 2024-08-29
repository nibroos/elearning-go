package controller

import (
	"context"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	pb "github.com/nibroos/elearning-go/users-service/internal/proto"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
	"github.com/nibroos/elearning-go/users-service/internal/service"

	"google.golang.org/grpc"
)

// UserController holds the methods for handling user-related HTTP requests.
type UserController struct {
    service service.UserService
	repo repository.UserRepository
}

func RunGRPCServer(repo repository.UserRepository) error {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        return err
    }

    grpcServer := grpc.NewServer()
    userService := service.NewUserService(repo)
    pb.RegisterUserServiceServer(grpcServer, userService)
    log.Println("gRPC server is running on port 50051")

    return grpcServer.Serve(lis)
}

// NewUserController creates a new instance of UserController.
func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (c *UserController) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
    searchParams := map[string]string{
        "global": req.GetGlobal(),
        "name":   req.GetName(),
        "email":  req.GetEmail(),
    }

    users, err := c.service.GetUsers(searchParams)
    if err != nil {
        return nil, err
    }

    // Mapping the response
    var pbUsers []*pb.User
    for _, user := range users {
        pbUsers = append(pbUsers, &pb.User{
            Id:        user.ID,
            Name:      user.Name,
            Username:  user.Username,
            Email:     user.Email,
        })
    }

    return &pb.GetUsersResponse{Users: pbUsers}, nil
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
    // repo := repository.NewUserRepository()
    // userService := service.NewUserService(repo)

    // userService.CreateUser(c)
    return c.SendString("Create User")
}

func (uc *UserController) GetUser(c *fiber.Ctx) error {
	// Handle the request here
    return c.SendString("Get User")
}

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	// Handle the request here
    return c.SendString("Update User")
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
    // Handle the request here
    return c.SendString("Delete User")
}