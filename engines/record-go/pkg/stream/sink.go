package stream

import (
	"errors"
	"gocv.io/x/gocv"
	"log"
	"time"
)

// getVideoWriter creates a new gocv VideoWriter for the given StreamContext.
func getVideoWriter(ctx StreamContext, src <-chan gocv.Mat) (*gocv.VideoWriter, error) {
	// FIXME: determine parameters from stream context

	var fileName = "/tmp/go-record-" + time.Now().Format("20060102-150405") + ".avi"
	var img gocv.Mat
	var more bool

	// get image parameters
	img, more = <-src
	if !more {
		return nil, errors.New("could not determine stream dimensions")
	}
	cols, rows := img.Cols(), img.Rows()
	log.Printf("determined dimensions: %d x %d\n", cols, rows)

	// guess video framerate by reading 30 frames
	then := time.Now()
	n := 30
	for i := 0; i < n; i++ {
		_, more = <-src
		if !more {
			return nil, errors.New("could not determine stream fps")
		}
	}

	duration := time.Since(then)
	fps := float64(n) / duration.Seconds()
	log.Printf("determined fps: %.2f\n", fps)

	return gocv.VideoWriterFile(fileName, "MJPG", fps, cols, rows, true)
}

func SaveVideo(ctx StreamContext, src <-chan gocv.Mat) {
	writer, err := getVideoWriter(ctx, src)
	defer writer.Close()

	for {
		img, more := <-src
		if !more {
			return
		}
		err = writer.Write(img)

		if err != nil {
			log.Println("error while writing video", err)
			return
		}
	}
}

// Sink to display the frames in a gocv Window
func WindowDisplay(ctx StreamContext, src <-chan gocv.Mat) {
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
