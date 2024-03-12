package grpc_sender

import (
	"context"
	"errors"
	"fmt"
	chatv1 "github.com/Sergey11001/protocol/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log/slog"
	"service1/internal/utils"
	"time"
)

type GRPCSender struct {
	log *slog.Logger
}

func NewGRPCSender(log *slog.Logger) *GRPCSender {
	return &GRPCSender{
		log: log,
	}
}

func (g *GRPCSender) StartDispatch(ctx context.Context, currentAddr, addr string, rmNodeCh chan string) error {
	client, err := newGRPCClient(addr)
	if err != nil {
		g.log.Error(fmt.Sprintf("failed to create grpc client: %v", err))
		return err
	}

	stream, err := client.SendAndGetMessage(ctx)
	if err != nil {
		g.log.Error(fmt.Sprintf("failed to create grpc stream: %v", err))
		return err
	}

	closeStream := make(chan struct{})

	go func() {
		for {
			_, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || status.Code(err) == codes.Canceled {
					g.log.Error(fmt.Sprintf("eof: %v", err.Error()))
				} else {
					g.log.Error(fmt.Sprintf("error while receiving stream response: %v", err.Error()))
				}
				closeStream <- struct{}{}
				return
			}
		}
	}()

	for {
		select {
		case <-time.After(time.Second * 1):
			msg := utils.GenerateString(5)
			err := stream.Send(&chatv1.Message{
				Message:  msg,
				Sender:   currentAddr,
				Receiver: addr,
			})

			if err != nil {
				if errors.Is(err, io.EOF) || status.Code(err) == codes.Canceled {
					g.log.Error(fmt.Sprintf("failed to send: %v", err))
					continue
				}
				continue
			}
			fmt.Println(fmt.Sprintf("sent message to %s: %s", addr, msg))
		case <-closeStream:
			stream.CloseSend()
			g.log.Info(fmt.Sprintf("stream closed: %s", addr))
			rmNodeCh <- addr
			return nil
		}
	}
}

func newGRPCClient(addr string) (chatv1.ChatClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s:%w", "failed to dial", err)
	}

	client := chatv1.NewChatClient(conn)

	return client, nil
}
