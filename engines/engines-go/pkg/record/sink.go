package record

import (
	"context"
	"errors"
	"gocv.io/x/gocv"
	"log"
	"time"
)

// getVideoWriter creates a new gocv VideoWriter for the given StreamContext.
func getVideoWriter(ctx context.Context, src <-chan gocv.Mat) (*gocv.VideoWriter, error) {
	// FIXME: determine parameters from engine context

	var fileName = "/tmp/go-record-" + time.Now().Format("20060102-150405") + ".avi"
	var img gocv.Mat
	var more bool

	// get image parameters
	img, more = <-src
	if !more {
		return nil, errors.New("could not determine engine dimensions")
	}
	cols, rows := img.Cols(), img.Rows()
	log.Printf("determined dimensions: %d x %d\n", cols, rows)

	// guess video framerate by reading 30 frames
	then := time.Now()
	n := 30
	for i := 0; i < n; i++ {
		_, more = <-src
		if !more {
			return nil, errors.New("could not determine engine fps")
		}
	}

	duration := time.Since(then)
	fps := float64(n) / duration.Seconds()
	log.Printf("determined fps: %.2f\n", fps)

	return gocv.VideoWriterFile(fileName, "MJPG", fps, cols, rows, true)
}

func SaveVideoSink(ctx context.Context, src <-chan gocv.Mat) {
	writer, err := getVideoWriter(ctx, src)
	if err != nil {
		log.Println("error creating video writer", err)
		return
	}

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
