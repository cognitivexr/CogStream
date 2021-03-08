package engine

import (
	"gocv.io/x/gocv"
	"log"
	"net"
)

// WindowDisplayHandler pipes the decoded frames into a WindowDisplaySink.
func WindowDisplayHandler(conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)

	defer func() {
		close(frames)
		close(images)
	}()

	ctx, err := InitStreamContext(conn)
	if err != nil {
		log.Println("error initializing stream context", err)
		conn.Close()
		return
	}

	go ImageDecoder(ctx, frames, images)
	go ConnectionHandler(ctx, conn, frames)
	WindowDisplaySink(ctx, images)

}
