package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	pb "github.com/nibroos/elearning-go/users-service/internal/proto"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
	"github.com/nibroos/elearning-go/users-service/internal/service"
)

// UserController holds the methods for handling user-related HTTP requests.
type UserController struct {
    userService service.UserService
	repo repository.UserRepository
    db *sqlx.DB
}

// func RunGRPCServer(repo repository.UserRepository) error {
//     lis, err := net.Listen("tcp", ":50051")
//     if err != nil {
//         return err
//     }

//     grpcServer := grpc.NewServer()
//     userService := service.NewUserService(repo)
//     pb.RegisterUserServiceServer(grpcServer, userService)
//     log.Println("gRPC server is running on port 50051")

//     return grpcServer.Serve(lis)
// }

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

    users, err := c.userService.GetUsers(searchParams)
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

// func (uc *UserController) CreateUser(c *fiber.Ctx) error {
//     // repo := repository.NewUserRepository()
//     // userService := service.NewUserService(repo)

//     // userService.CreateUser(c)
//     return c.SendString("Create User")
// }
func (c *UserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    tx, err := c.db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    user, err := c.userService.CreateUser(ctx, tx, req.GetName(), req.GetEmail(), req.GetRoleId())
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return &pb.CreateUserResponse{
        Id:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
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