package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"service1/internal/app"
	"service1/internal/config"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := context.WithCancel(context.Background())

	application := app.New(log, cfg.GRPCConfig.Port)

	go application.MustRun()
	go application.Start(ctx, cfg.BroadcastPort, cfg.GRPCConfig.Port)

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	cancel()
	application.GRPCSrv.Stop()
}
