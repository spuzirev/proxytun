package main

import (
	"context"
	"log"
	"net"
	"net/url"

	"github.com/songgao/water"
	"golang.org/x/net/proxy"
)

func runClient(ctx context.Context, remoteAddr, proxyAddr string) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatalf("Failed to create device: %v", err)
	}
	log.Printf("Created tun interface: %s", iface.Name())

	var dialer proxy.Dialer = &net.Dialer{}
	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			log.Fatalf("Failed to parse proxyAddr: %v", err)
		}

		proxyDialer, err := proxy.FromURL(proxyURL, dialer)
		if err != nil {
			log.Fatalf("Failed to configure proxy dialer")
		}
		dialer = proxyDialer
	}

	go client(ctx, iface, dialer, remoteAddr)

	<-ctx.Done()
}

func client(ctx context.Context, iface *water.Interface, dialer proxy.Dialer, remoteAddr string) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down client")
			return
		default:
		}

		var (
			conn net.Conn
			err  error
		)
		if ctxDialer, ok := dialer.(proxy.ContextDialer); ok {
			conn, err = ctxDialer.DialContext(ctx, "tcp", remoteAddr)
		} else {
			conn, err = dialer.Dial("tcp", remoteAddr)
		}
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			continue
		}
		log.Printf("Connected to %s", conn.RemoteAddr())

		err = communicate(ctx, iface, conn)
		if err != nil {
			conn.Close()
			log.Printf("Communication error: %v", err)
		}
	}
}
