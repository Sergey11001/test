package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"service1/internal/grpc/server"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, port string) *App {
	gRPCServer := grpc.NewServer()
	server.RegisterServer(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", "failed to listen", err)
	}

	a.log.Info(fmt.Sprintf("gRPC server listening on port %s", a.port))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s:%w", "failed to serve", err)
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.Stop()
}
