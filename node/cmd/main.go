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
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())

	disc := discovery.NewDiscovery(log, cfg.GRPCConfig.Port, cfg.BroadcastPort, cfg.BroadcastPrefix)
	sender := grpc_sender.NewGRPCSender(log)

	application := app.New(log, cfg.GRPCConfig.Port, disc, sender)

	stopCh := make(chan os.Signal)
	errorCh := make(chan error)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Start(ctx); err != nil {
			errorCh <- err
		}
	}()

	select {
	case <-stopCh:
		log.Info("shutting down...")
		cancel()
		application.GRPCSrv.Stop()
		return
	case err := <-errorCh:
		log.Error(fmt.Sprintf("application error: %v", err))
		return
	}
}
