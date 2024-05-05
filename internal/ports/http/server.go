package http

import (
	"encoding/json"
	"fmt"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/maghaze/gateway/internal/ports/grpc"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	logger *zap.Logger
	auth   grpc.AuthClient

	managementApp *fiber.App
	clientApp     *fiber.App
}

func New(cfg *Config, log *zap.Logger, auth grpc.AuthClient) *Server {
	server := &Server{config: cfg, logger: log, auth: auth}

	fiberConfig := fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal, DisableStartupMessage: true}
	server.managementApp, server.clientApp = fiber.New(fiberConfig), fiber.New(fiberConfig)

	prometheus := fiberprometheus.New("auth")
	prometheus.RegisterAt(server.managementApp, "/metrics")
	server.managementApp.Use(prometheus.Middleware)

	server.managementApp.Get("/healthz/liveness", server.liveness)
	server.managementApp.Get("/healthz/readiness", server.readiness)

	v1 := server.clientApp.Group("/v1")
	v1.Group("/users", server.optionalAuthentication, server.proxy)
	v1.Group("/books", server.requiredAuthentication, server.proxy)

	return server
}

func (server *Server) Serve(managementPort, clientPort int) {
	go func() {
		server.logger.Info("HTTP management server starts listening on", zap.Int("port", managementPort))
		if err := server.managementApp.Listen(fmt.Sprintf(":%d", managementPort)); err != nil {
			server.logger.Fatal("error resolving HTTP server", zap.Error(err))
		}
	}()

	go func() {
		server.logger.Info("HTTP client server starts listening on", zap.Int("port", clientPort))
		if err := server.clientApp.Listen(fmt.Sprintf(":%d", clientPort)); err != nil {
			server.logger.Fatal("error resolving HTTP server", zap.Error(err))
		}
	}()
}
