package controller

import (
	"log"
	"net"

	pb "github.com/nibroos/elearning-go/users-service/internal/proto"
	"github.com/nibroos/elearning-go/users-service/internal/repository"
	"github.com/nibroos/elearning-go/users-service/internal/service"

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
    log.Println("gRPC server is running on port 50051")

    return grpcServer.Serve(lis)
}
