package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/maghaze/gateway/internal/config"
	"github.com/maghaze/gateway/internal/ports/grpc"
	"github.com/maghaze/gateway/internal/ports/http"
	"github.com/maghaze/gateway/pkg/logger"
)

type Server struct {
	managementPort int
	clientPort     int
	grpcPort       int
}

func (server Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		server.main(config.Load(true), trap)
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "run api-gateway server",
		Run:   run,
	}

	cmd.Flags().IntVar(&server.managementPort, "management-port", 8080, "The port the metrics and probe endpoints binds to")
	cmd.Flags().IntVar(&server.clientPort, "client-port", 8081, "The port the api-gateway http server endpoints binds to")
	cmd.Flags().IntVar(&server.grpcPort, "grpc-port", 9090, "The port the grpc endpoint listens on")

	return cmd
}

func (server *Server) main(cfg *config.Config, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	authGrpcClient := grpc.NewAuthClient(cfg.GRPC, logger)

	httpServer := http.New(cfg.HTTP, logger, authGrpcClient)
	go httpServer.Serve(server.managementPort, server.clientPort)

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	logger.Info("exiting by receiving a unix signal", field)
}
