package config

import (
	"github.com/maghaze/gateway/internal/ports/grpc"
	"github.com/maghaze/gateway/internal/ports/http"
	"github.com/maghaze/gateway/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	HTTP   *http.Config   `koanf:"http"`
	GRPC   *grpc.Config   `koanf:"grpc"`
}
