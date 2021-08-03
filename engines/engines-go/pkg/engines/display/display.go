package display

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/engine"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"context"
	"gocv.io/x/gocv"
)

var Factory engine.Factory = &engineFactory{}
var Descriptor = engines.EngineDescriptor{
	Name: "display",
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
