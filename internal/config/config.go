package config

import (
	"github.com/maghaze/api-gateway/internal/ports/grpc"
	"github.com/maghaze/api-gateway/internal/ports/http"
	"github.com/maghaze/api-gateway/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	HTTP   *http.Config   `koanf:"http"`
	GRPC   *grpc.Config   `koanf:"grpc"`
}
