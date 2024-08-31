package main

import (
	"context"
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/songgao/water"
	"golang.org/x/sync/errgroup"
)

func connToTunLoop(ctx context.Context, iface *water.Interface, conn net.Conn) error {
	header := make([]byte, 2)
	buf := make([]byte, 8192)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := io.ReadFull(conn, header)
		if err != nil {
			return err
		}
		packetLen := binary.BigEndian.Uint16(header)
		log.Printf("Expecting %d bytes from conn", packetLen)

		n, err := io.ReadFull(conn, buf[:packetLen])
		if err != nil {
			return err
		}
		log.Printf("Read %d bytes from conn", n)

		n, err = iface.Write(buf[:packetLen])
		if err != nil {
			return err
		}
		log.Printf("Wrote %d bytes to tun", n)
	}
}

func tunToConnLoop(ctx context.Context, iface *water.Interface, conn net.Conn) error {
	buf := make([]byte, 8192)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := iface.Read(buf)
		if err != nil {
			return err
		}
		log.Printf("Read %d bytes from tun", n)

		err = binary.Write(conn, binary.BigEndian, uint16(n))
		if err != nil {
			return err
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			return err
		}
		log.Printf("Wrote %d bytes to tun", n)
	}
}

func communicate(ctx context.Context, iface *water.Interface, conn net.Conn) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return connToTunLoop(ctx, iface, conn) })
	eg.Go(func() error { return tunToConnLoop(ctx, iface, conn) })
	return eg.Wait()
}
