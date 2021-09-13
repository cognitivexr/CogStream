package stats

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/engine"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"context"
	"fmt"
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

type statsSink struct {
	firstFrameTime time.Time
	lastFrameTime  time.Time
	frameCount     int
}

func NewEngine() pipeline.Engine {
	return &statsSink{
		time.Time{},
		time.Time{},
		0,
	}
}

func (w *statsSink) Process(_ context.Context, _ *pipeline.Frame, _ pipeline.EngineResultWriter) error {
	w.frameCount += 1
	currentFrameTime := time.Now()
	if w.firstFrameTime.IsZero() {
		w.firstFrameTime = currentFrameTime
		w.lastFrameTime = currentFrameTime
		return nil
	}
	timeBetween := currentFrameTime.Sub(w.lastFrameTime)
	fps := 1 / timeBetween.Seconds()
	cumulativeAverage := currentFrameTime.Sub(w.firstFrameTime).Seconds() / float64(w.frameCount)
	fpsAvg := 1 / cumulativeAverage
	fmt.Printf("FRAME %d\ncurr. fps: %v\nframetime: %v\ncumu. ftm: %v\navrg. fps: %v\n", w.frameCount, fps, timeBetween, cumulativeAverage, fpsAvg)

	w.lastFrameTime = currentFrameTime
	return nil
}
