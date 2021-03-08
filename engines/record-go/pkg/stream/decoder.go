package stream

import (
	"gocv.io/x/gocv"
	"log"
)

func ImageDecoder(ctx StreamContext, source <-chan []byte, dest chan<- gocv.Mat) {
	flags := gocv.IMReadColor // FIXME determine flags from StreamContext

	for {
		frame, more := <-source

		if !more {
			break
		}
		if len(source) >= 10 {
			// TODO: back pressure
			log.Printf("queue filling up. source: %d, dest: %d\n", len(source), len(dest))
		}

		//then := time.Now()
		img, err := gocv.IMDecode(frame, flags)
		if err != nil {
			log.Println("error while decoding image", err)
			break
		}

		// FIXME determine transformations from StreamContext
		gocv.Flip(img, &img, 0)
		gocv.CvtColor(img, &img, gocv.ColorBGRToRGB)
		//log.Printf("transformation took %v\n", time.Now().Sub(then))

		dest <- img
	}

	log.Println("decoder returning")
}
