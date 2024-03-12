package app

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
)

type App struct {
	log *slog.Logger
}

func New(log *slog.Logger) *App {
	return &App{
		log: log,
	}
}

func (a *App) Start(ctx context.Context, broadcastPort int, broadcastIP, prefixIP string) error {
	addr := &net.UDPAddr{
		Port: broadcastPort,
		IP:   net.ParseIP(broadcastIP),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		a.log.Error("failed to listen udp: " + err.Error())
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	remoteCons := new(sync.Map)

	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		body := buf[:n]
		if !bytes.HasPrefix(body, []byte(prefixIP)) {
			continue
		}

		body = bytes.TrimPrefix(body, []byte(prefixIP))

		if _, ok := remoteCons.Load(remoteAddr.String()); !ok {
			remoteCons.Store(remoteAddr.String(), &remoteAddr)

			go func() {
				remoteCons.Range(func(key, value interface{}) bool {
					keyAddr, ok := key.(string)
					if !ok {
						return true
					}

					if keyAddr == remoteAddr.String() {
						return true
					}

					if _, err := conn.WriteTo(body, *value.(*net.Addr)); err != nil {
						remoteCons.Delete(key)
						return true
					}

					if _, err := conn.WriteTo([]byte(keyAddr), remoteAddr); err != nil {
						return true
					}

					return true
				})
			}()
		}

	}
}
