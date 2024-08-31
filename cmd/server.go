package main

import (
	"context"
	"log"
	"net"

	"github.com/songgao/water"
)

func runServer(ctx context.Context, listenAddr string) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatalf("Failed to open tun device: %v", err)
	}
	log.Printf("Created tun interface: %s", iface.Name())
	defer iface.Close()

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	go server(ctx, iface, listener)

	<-ctx.Done()
	listener.Close()
}

func server(ctx context.Context, iface *water.Interface, listener net.Listener) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down server")
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		err = communicate(ctx, iface, conn)
		if err != nil {
			conn.Close()
			log.Printf("Communication error: %v", err)
		}
	}
}
