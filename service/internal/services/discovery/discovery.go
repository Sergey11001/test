package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

const (
	BroadcastHost = "255.255.255.255"
)

type Discovery struct {
	log             *slog.Logger
	mu              *sync.Mutex
	peers           map[string]struct{}
	localPort       string
	broadcastPort   string
	broadcastPrefix string
}

func NewDiscovery(log *slog.Logger, localPort, broadcastPort, broadcastPrefix string) *Discovery {
	return &Discovery{
		log:             log,
		mu:              &sync.Mutex{},
		peers:           make(map[string]struct{}),
		localPort:       localPort,
		broadcastPort:   broadcastPort,
		broadcastPrefix: broadcastPrefix,
	}
}

func (d *Discovery) Broadcast(ctx context.Context, currentAddr string, addNodeCh chan string) {
	pc, err := net.ListenPacket("udp4", fmt.Sprintf(":%s", d.localPort))
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(BroadcastHost, d.broadcastPort))
	if err != nil {
		d.log.Error(fmt.Sprintf("failed to resolve broadcast address: %v", err))
	}

	buf := make([]byte, 1024)
	go func() {
		for {
			n, _, err := pc.ReadFrom(buf)
			if err != nil {
				d.log.Error(fmt.Sprintf("failed to read from broadcast: %v", err))
				continue
			}
			if d.nodeContains(string(buf[:n])) {
				continue
			}

			remoteAddr := string(buf[:n])
			d.addPeer(remoteAddr)
			addNodeCh <- remoteAddr
		}
	}()

	for {
		body := d.broadcastPrefix + currentAddr
		_, err = pc.WriteTo([]byte(body), addr)
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}

func (d *Discovery) RemoveNode(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.peers, addr)
}

func (d *Discovery) addPeer(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.peers[addr] = struct{}{}
}

func (d *Discovery) nodeContains(addr string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.peers[addr]; ok {
		return true
	}
	return false
}
