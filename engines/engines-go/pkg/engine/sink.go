package engine

import (
	"gocv.io/x/gocv"
	"log"
)

// Sink to display the frames in a gocv Window
func WindowDisplaySink(ctx StreamContext, src <-chan gocv.Mat) {
	window := gocv.NewWindow("window")
	defer window.Close()

	doneCh := make(chan bool)
	defer func() {
		close(doneCh)
		log.Println("window display sink returning")
	}()

	go func() {
		for {
			if window.WaitKey(1000) >= 0 {
				doneCh <- true
				return
			}
		}
	}()

	for {
		select {
		case img := <-src:
			window.IMShow(img)
		case done := <-doneCh:
			if done {
				return
			}
		}
	}

}
