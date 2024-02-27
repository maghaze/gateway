package grpc

import (
	"context"
	"errors"

	pb "github.com/CafeKetab/PBs/golang/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthClient interface {
	Authenticate(ctx context.Context, token string) (uint64, error)
}

type authClient struct {
	logger *zap.Logger
	api    pb.AuthClient
}

func NewAuthClient(cfg *Config, lg *zap.Logger) *authClient {
	client := &authClient{logger: lg}

	connection, err := grpc.Dial(cfg.Targets.Auth, grpc.WithInsecure())
	if err != nil {
		lg.Panic("error while instantiating auth grpc client", zap.Error(err))
	}
	client.api = pb.NewAuthClient(connection)

	return client
}

func (c *authClient) Authenticate(ctx context.Context, token string) (uint64, error) {
	pbId, err := c.api.GetIdFromToken(ctx, &pb.Token{Value: token})
	if err != nil {
		errString := "Error getting id from token"
		c.logger.Error(errString, zap.String("token", token), zap.Error(err))
		return 0, errors.New(errString)
	}
	return pbId.Value, nil
}
