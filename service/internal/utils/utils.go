package utils

import (
	"math/rand"
	"net"
	"strings"
)

func GenerateString(length int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(charset[rand.Intn(len(charset))])
	}

	return b.String()
}

func GetLocalIP() (string, error) {
	var localIP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return localIP, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	return localIP, nil
}
