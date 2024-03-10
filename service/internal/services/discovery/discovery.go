package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"service1/internal/utils"
	"sync"
	"time"
)

const (
	PrefixIP      = "qwe"
	BroadcastHost = "255.255.255.255"
)

type Discovery struct {
	log       *slog.Logger
	mu        *sync.Mutex
	peers     map[string]struct{}
	AddPeerCh chan string
}

func NewDiscovery(log *slog.Logger, addCh chan string) *Discovery {
	return &Discovery{
		log:       log,
		mu:        &sync.Mutex{},
		peers:     make(map[string]struct{}),
		AddPeerCh: addCh,
	}
}

func (d *Discovery) AddPeer(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.peers[addr] = struct{}{}
}

func (d *Discovery) Broadcast(ctx context.Context, broadcastPort, localPort string) {
	pc, err := net.ListenPacket("udp4", fmt.Sprintf(":%s", localPort))
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	localHost := utils.GetLocalHost()

	addr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(BroadcastHost, broadcastPort))
	if err != nil {
		d.log.Error(fmt.Sprintf("failed to resolve broadcast address: %v", err))
	}

	buf := make([]byte, 1024)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, _, err := pc.ReadFrom(buf)
				if err != nil {
					d.log.Error(fmt.Sprintf("failed to read from broadcast: %v", err))
					continue
				}
				if d.PeerContains(string(buf[:n])) {
					continue
				}

				remoteAddr := string(buf[:n])
				d.AddPeer(remoteAddr)
				d.AddPeerCh <- remoteAddr
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			body := PrefixIP + net.JoinHostPort(localHost, localPort)
			_, err = pc.WriteTo([]byte(body), addr)
			if err != nil {
				panic(err)
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func (d *Discovery) PeerContains(addr string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.peers[addr]; ok {
		return true
	}
	return false
}

func (d *Discovery) RemovePeer(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.peers, addr)
}
