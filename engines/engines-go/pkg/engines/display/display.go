package display

import (
	"cognitivexr.at/cogstream/pkg/pipeline"
	"context"
	"gocv.io/x/gocv"
)

type windowDisplaySink struct {
	window *gocv.Window
}

func NewEngine() pipeline.Engine {
	return &windowDisplaySink{
		gocv.NewWindow("window"),
	}
}

func (w *windowDisplaySink) Process(_ context.Context, frame *pipeline.Frame, _ pipeline.EngineResultWriter) error {
	w.window.IMShow(*frame.Mat)

	if w.window.WaitKey(1) >= 0 {
		return pipeline.Stop
	}

	return nil
}
