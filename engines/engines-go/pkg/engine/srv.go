package engine

import (
	"gocv.io/x/gocv"
	"net"
)

func WindowDisplayHandler(conn net.Conn) {
	frames := make(chan []byte, 30)
	images := make(chan gocv.Mat)
	ctx := NewStreamContext() // FIXME: initialize engine context

	go ImageDecoder(ctx, frames, images)
	go ConnectionHandler(ctx, conn, frames)
	WindowDisplaySink(ctx, images)

	close(frames)
	close(images)
}
