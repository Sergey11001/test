package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test_project/broadcast/internal/app"
	"test_project/broadcast/internal/config"
)

func main() {
	cfg := config.MustLoad()
	ctx, cancel := context.WithCancel(context.Background())

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("starting broadcast service")

	application := app.New(log)

	stopCh := make(chan os.Signal)
	errorCh := make(chan error)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Start(ctx, cfg.BroadcastPort, cfg.BroadcastIP, cfg.PrefixIP); err != nil {
			errorCh <- err
		}
	}()

	select {
	case <-stopCh:
		cancel()
		return
	case err := <-errorCh:
		log.Error(fmt.Sprintf("application error: %v", err))
		return
	}
}
