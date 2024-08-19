package service

import (
	"context"

	pb "github.com/nibroos/elearning-go/users-service/users-service/internal/proto"
)

type HealthService struct {
	pb.UnimplementedHealthServiceServer
}

func (s *HealthService) CheckHealth(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Message: "Users-Service is running"}, nil
}
