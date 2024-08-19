package server

import (
	"log"
	"net"

	pb "github.com/nibroos/elearning-go/users-service/internal/proto"
	"github.com/nibroos/elearning-go/users-service/internal/service"

	"google.golang.org/grpc"
)

func RunGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &service.HealthService{})

	log.Println("gRPC server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
