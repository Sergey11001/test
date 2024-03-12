package server

import (
	"fmt"
	chatv1 "github.com/Sergey11001/protocol/gen/go"
	"google.golang.org/grpc"
	"io"
	"service1/internal/utils"
)

type server struct {
	chatv1.UnimplementedChatServer
}

func RegisterServer(gRPC *grpc.Server) {
	chatv1.RegisterChatServer(gRPC, &server{})
}

func (s *server) SendAndGetMessage(stream chatv1.Chat_SendAndGetMessageServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		fmt.Printf("received from %s: %s\n", req.Sender, req.Message)

		res := &chatv1.Message{
			Message:  utils.GenerateString(5),
			Sender:   req.Receiver,
			Receiver: req.Sender,
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}
