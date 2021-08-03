package fps

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/engine"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"context"
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

var Factory engine.Factory = &engineFactory{}
var Descriptor = engines.EngineDescriptor{
	Name: "fps",
	Specification: engines.Specification{
		Operation:   messages.OperationRecord,
		InputFormat: format.AnyFormat,
		Attributes:  messages.NewAttributes(),
	},
}

type engineFactory struct{}

func (e *engineFactory) Descriptor() engines.EngineDescriptor {
	return Descriptor
}

func (e *engineFactory) NewEngine() pipeline.Engine {
	return NewEngine()
}

type fpsDisplaySink struct {
	window        *gocv.Window
	lastFrameTime time.Time
}

func NewEngine() pipeline.Engine {
	return &fpsDisplaySink{
		gocv.NewWindow("FPS display"),
		time.Now(),
	}
}

func (w *fpsDisplaySink) Process(_ context.Context, frame *pipeline.Frame, _ pipeline.EngineResultWriter) error {
	currentFrameTime := time.Now()
	timeBetween := currentFrameTime.Sub(w.lastFrameTime)
	fps := 1 / timeBetween.Seconds()
	fmt.Printf("fps: %v", fps)
	w.window.IMShow(*frame.Mat)

	if w.window.WaitKey(1) >= 0 {
		return pipeline.Stop
	}

	w.lastFrameTime = currentFrameTime
	return nil
}
