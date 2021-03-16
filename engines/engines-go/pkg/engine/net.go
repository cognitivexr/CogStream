package engine

import (
	"context"
	"log"
	"net"
)

func ConnectionHandler(ctx context.Context, conn net.Conn, frames chan<- []byte) {
	remoteAddr := conn.RemoteAddr()
	log.Printf("[%s] accepted connection\n", remoteAddr)
	defer func() {
		log.Printf("[%s] engine handler returning\n", remoteAddr)
		conn.Close()
	}()

	scanner := NewFrameScanner(conn)

	for scanner.Next() {
		if scanner.Err() != nil {
			log.Printf("[%s] error while reading packet engine: %s\n", remoteAddr, scanner.Err())
			break
		}

		frame := scanner.frame
		frames <- frame
	}
}
