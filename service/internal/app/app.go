package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	grpcapp "service1/internal/app/grpc"
	"service1/internal/utils"
)

type Discovery interface {
	Broadcast(ctx context.Context, currentAddr string, addNodeCh chan string)
	RemoveNode(addr string)
}

type Sender interface {
	StartDispatch(ctx context.Context, currentAddr, addr string, rmNodeCh chan string) error
}

type App struct {
	log       *slog.Logger
	GRPCSrv   *grpcapp.App
	Discovery Discovery
	Sender    Sender
	addr      string
}

func New(log *slog.Logger, port string, discoveryService Discovery, grpcSender Sender) *App {
	grpcApp := grpcapp.New(log, port)
	currentAddr := net.JoinHostPort(utils.GetLocalHost(), port)

	return &App{
		log:       log,
		GRPCSrv:   grpcApp,
		Discovery: discoveryService,
		Sender:    grpcSender,
		addr:      currentAddr,
	}
}

func (a *App) Start(ctx context.Context) error {
	addNodeCh := make(chan string)

	go a.Discovery.Broadcast(ctx, a.addr, addNodeCh)
	go a.nodeLoop(ctx, addNodeCh)

	return a.GRPCSrv.Start()
}
func (a *App) nodeLoop(ctx context.Context, addNodeCh chan string) {
	rmNodeCh := make(chan string)
	for {
		select {
		case <-ctx.Done():
			return
		case receiverAddr := <-addNodeCh:
			a.log.Info(fmt.Sprintf("new node: %s", receiverAddr))
			go func() {
				if err := a.Sender.StartDispatch(ctx, a.addr, receiverAddr, rmNodeCh); err != nil {
					a.log.Error(fmt.Sprintf("failed to dial to new peer: %v", err))
				}
			}()
		case addr := <-rmNodeCh:
			a.Discovery.RemoveNode(addr)
			fmt.Println("removing node: ", addr)
		}
	}
}
