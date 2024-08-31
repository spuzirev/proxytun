package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

var (
	listenAddr string
	remoteAddr string
	proxyAddr  string
)

func registerFlags() {
	flag.StringVar(&listenAddr, "listen", "", "Listen addr, such as \":5000\"")
	flag.StringVar(&remoteAddr, "remote", "", "Remote addr to connect, such as \"127.0.0.1:5001\"")
	flag.StringVar(&proxyAddr, "proxy", "", "Proxy addr")
}

func parseFlags() {
	flag.Parse()
}

func main() {
	registerFlags()
	parseFlags()

	if listenAddr == "" && remoteAddr == "" {
		log.Fatalf("You should specify at least one of -listen|-remote")
	}

	if listenAddr != "" && remoteAddr != "" {
		log.Fatalf("You should specify only one of -listen|-remote")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if listenAddr != "" {
		runServer(ctx, listenAddr)
	}

	if remoteAddr != "" {
		runClient(ctx, remoteAddr, proxyAddr)
	}
}
