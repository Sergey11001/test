package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"service1/internal/app"
	"service1/internal/config"
	"service1/internal/services/discovery"
	"service1/internal/services/grpc_sender"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())

	disc := discovery.NewDiscovery(log, cfg.GRPCConfig.Port, cfg.BroadcastPort)
	sender := grpc_sender.NewGRPCSender(log)

	application := app.New(log, cfg.GRPCConfig.Port, disc, sender)

	go func() {
		if err := application.Start(ctx); err != nil {
			log.Error(fmt.Sprintf("failed to start application: %v", err))
			return
		}
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	cancel()
	application.GRPCSrv.Stop()
}
