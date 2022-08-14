package token

import (
	"context"
	"github.com/synthia-telemed/backend-api/pkg/token/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service interface {
	GenerateToken(userID uint64, role string) (string, error)
}

type Config struct {
	Endpoint string `env:"TOKEN_SERVICE_ENDPOINT,required"`
}

type GRPCTokenService struct {
	tokenClient proto.TokenClient
}

func NewGRPCTokenService(config *Config) (*GRPCTokenService, error) {
	conn, err := grpc.Dial(config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	tokenClient := proto.NewTokenClient(conn)

	return NewGRPCTokenServiceWithClient(tokenClient), nil
}

func NewGRPCTokenServiceWithClient(tokenClient proto.TokenClient) *GRPCTokenService {
	return &GRPCTokenService{tokenClient: tokenClient}
}

func (s GRPCTokenService) GenerateToken(userID uint64, role string) (string, error) {
	req := &proto.GenerateTokenRequest{
		UserID: userID,
		Role:   role,
	}
	res, err := s.tokenClient.GenerateToken(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetToken(), nil
}
