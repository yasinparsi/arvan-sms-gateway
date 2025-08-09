package grpc

import (
	"billing-service/internal/storage"
	pb "billing-service/proto"
	"context"
)

type Server struct {
	pb.UnimplementedBillingServiceServer
	redisClient *storage.RedisClient
}

func NewServer(redisClient *storage.RedisClient) *Server {
	return &Server{redisClient: redisClient}
}

func (s *Server) CheckBalance(ctx context.Context, req *pb.BillingRequest) (*pb.BillingResponse, error) {
	allowed, _, err := s.redisClient.CheckAmountAtomically(ctx, req.UserId, req.Cost)
	if err != nil {
		return &pb.BillingResponse{
			Allowed: false,
			Error:   err.Error(),
		}, nil
	}

	if !allowed {
		return &pb.BillingResponse{
			Allowed: false,
			Error:   "insufficient balance",
		}, nil
	}

	return &pb.BillingResponse{
		Allowed: true,
		Error:   "",
	}, nil
}
