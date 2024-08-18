package controller

import (
	"net"

	"github.com/nibroos/users-service/internal/repository"
	"github.com/nibroos/users-service/internal/service"
	pb "github.com/nibroos/users-service/proto"
	"google.golang.org/grpc"
)

func RunGRPCServer(repo repository.UserRepository) error {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        return err
    }

    grpcServer := grpc.NewServer()
    userService := service.NewUserService(repo)
    pb.RegisterUserServiceServer(grpcServer, userService)

    return grpcServer.Serve(lis)
}
