package recorder

import (
	"cognitivexr.at/cogstream/api/engines"
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/api/messages"
	"cognitivexr.at/cogstream/pkg/engine"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"context"
	"errors"
	"gocv.io/x/gocv"
	"log"
	"sync"
	"time"
)

var Descriptor = engines.EngineDescriptor{
	Name: "record",
	Specification: engines.Specification{
		Operation:   messages.OperationRecord,
		InputFormat: format.AnyFormat,
		Attributes:  messages.NewAttributes(),
	},
}

var Factory engine.Factory = &engineFactory{}

type Engine struct {
	writer *gocv.VideoWriter

	running     bool
	initialized bool

	initError error

	buffer chan *pipeline.Frame
	mutex  sync.Mutex
}

type engineFactory struct{}

func (e *engineFactory) Descriptor() engines.EngineDescriptor {
	return Descriptor
}

func (e *engineFactory) NewEngine() pipeline.Engine {
	return NewEngine()
}

func NewEngine() *Engine {
	return &Engine{
		buffer: make(chan *pipeline.Frame),
	}
}

func (e *Engine) Process(_ context.Context, frame *pipeline.Frame, _ pipeline.EngineResultWriter) error {
	if e.running {
		return e.writer.Write(*frame.Mat)
	}

	return e.processInit(frame)
}

func (e *Engine) processInit(frame *pipeline.Frame) error {
	// determining stream parameters
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !e.initialized {
		e.initialized = true
		go e.initializeVideoWriter()
	}

	if e.initError != nil {
		return e.initError
	}

	if !e.running {
		e.buffer <- frame
		return nil
	} else {
		close(e.buffer)
		return e.writer.Write(*frame.Mat)
	}
}

func (e *Engine) initializeVideoWriter() {
	// FIXME: determine parameters from engine context

	var fileName = "/tmp/go-recorder-" + time.Now().Format("20060102-150405") + ".avi"
	var frame *pipeline.Frame
	var more bool
	fpsWindowSize := 30 // frames to capture before determining FPS

	// get image parameters
	frame, more = <-e.buffer
	if !more {
		e.initError = errors.New("could not determine engine dimensions")
		return
	}
	cols, rows := frame.Mat.Cols(), frame.Mat.Rows()
	log.Printf("determined dimensions: %d x %d\n", cols, rows)

	// guess video framerate by reading 30 frames
	then := time.Now()
	for i := 0; i < fpsWindowSize; i++ {
		_, more = <-e.buffer
		if !more {
			e.initError = errors.New("could not determine engine fps")
			return
		}
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	duration := time.Since(then)
	fps := float64(fpsWindowSize) / duration.Seconds()
	log.Printf("determined fps: %.2f\n", fps)

	writer, err := gocv.VideoWriterFile(fileName, "MJPG", fps, cols, rows, true)
	if err != nil {
		e.initError = err
		return
	}

	e.writer = writer
	e.running = true
}
