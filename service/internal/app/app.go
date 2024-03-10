package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	grpcapp "service1/internal/app/grpc"
	"service1/internal/services/discovery"
	"service1/internal/services/grpc_sender"
	"service1/internal/utils"
)

type App struct {
	log        *slog.Logger
	GRPCSrv    *grpcapp.App
	Discovery  *discovery.Discovery
	GRPCSender *grpc_sender.GRPCSender
	addr       string
}

func New(log *slog.Logger, port string) *App {
	grpcApp := grpcapp.New(log, port)

	removeCh := make(chan string, 20)
	addCh := make(chan string)

	currentAddr := net.JoinHostPort(utils.GetLocalHost(), port)

	discoveryService := discovery.NewDiscovery(log, addCh)
	grpcSender := grpc_sender.NewGRPCSender(log, removeCh)

	return &App{
		log:        log,
		GRPCSrv:    grpcApp,
		Discovery:  discoveryService,
		GRPCSender: grpcSender,
		addr:       currentAddr,
	}
}

func (a *App) MustRun() {
	err := a.GRPCSrv.Start()
	panic(err)
}

func (a *App) Start(ctx context.Context, broadcastPort string, localPort string) {
	go a.Discovery.Broadcast(ctx, broadcastPort, localPort)
	a.peerLoop(ctx)
}
func (a *App) peerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case receiverAddr := <-a.Discovery.AddPeerCh:
			a.log.Info(fmt.Sprintf("new peer: %s", receiverAddr))
			go func() {
				if err := a.GRPCSender.BindStream(ctx, a.addr, receiverAddr); err != nil {
					a.log.Error(fmt.Sprintf("failed to dial to new peer: %v", err))
				}
			}()
		case addr := <-a.GRPCSender.RemovePeerCh:
			a.Discovery.RemovePeer(addr)
			fmt.Println("removing node: ", addr)
		}
	}
}
