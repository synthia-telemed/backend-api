package token

import (
	"context"
	pb "github.com/synthia-telemed/backend-api/pkg/services/token/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service interface {
	GenerateToken(userID uint64, role string) (string, error)
}

type gRPCTokenService struct {
	tokenClient pb.TokenClient
}

func NewGRPCTokenService(serviceHost string) (*gRPCTokenService, error) {
	conn, err := grpc.Dial(serviceHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	tokenClient := pb.NewTokenClient(conn)

	return &gRPCTokenService{tokenClient: tokenClient}, nil
}

func (s gRPCTokenService) GenerateToken(userID uint64, role string) (string, error) {
	req := &pb.GenerateTokenRequest{
		UserID: userID,
		Role:   role,
	}
	res, err := s.tokenClient.GenerateToken(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetToken(), nil
}
