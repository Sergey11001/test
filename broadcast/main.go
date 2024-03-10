package main

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

const (
	BroadcastIP   = "255.255.255.255"
	BroadcastPort = 12345
	PrefixIP      = "qwe"
)

func main() {
	fmt.Println("Broadcasting...")

	addr := &net.UDPAddr{
		Port: BroadcastPort,
		IP:   net.ParseIP(BroadcastIP),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	remoteConns := new(sync.Map)

	for {
		buf := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		body := buf[:n]
		if !bytes.HasPrefix(body, []byte(PrefixIP)) {
			continue
		}

		body = bytes.TrimPrefix(body, []byte(PrefixIP))

		if _, ok := remoteConns.Load(remoteAddr.String()); !ok {
			remoteConns.Store(remoteAddr.String(), &remoteAddr)
		}

		go func() {
			remoteConns.Range(func(key, value interface{}) bool {
				if key == remoteAddr.String() {
					return true
				}

				if _, err := conn.WriteTo(body, *value.(*net.Addr)); err != nil {
					remoteConns.Delete(key)
					return true
				}
				return true
			})
		}()
	}
}
