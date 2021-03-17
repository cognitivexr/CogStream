package engine

import (
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/engines/pkg/transform"
	"context"
	"gocv.io/x/gocv"
	"log"
)

// TODO: separate transport (JPEG channel, bitmap channel, ...) and stream (source format, engine format)

func ImageDecoder(ctx context.Context, source <-chan []byte, dest chan<- gocv.Mat) {
	metadata, ok := GetStreamMetadata(ctx)
	if !ok {
		log.Println("could not get stream metadata from context")
		metadata = NewStreamMetadata()
	}

	// create transformer
	var tf transform.Function

	if metadata.EngineFormat == format.AnyFormat {
		log.Println("engine format = any format")
		tf = transform.NoTransform
	} else if metadata.ClientFormat == format.AnyFormat {
		log.Println("clientFormat format = any format")
		// FIXME: determine format from first frame
		tf = transform.NoTransform
	}  else {
		tf1, err := transform.BuildTransformer(metadata.ClientFormat, metadata.EngineFormat)
		tf = tf1
		if err != nil {
			log.Printf("error building transformer: %v\n", err)
			tf = transform.NoTransform
		}
	}


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
		tf(img, &img)

		dest <- img
	}

	log.Println("decoder returning")
}
