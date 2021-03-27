package transform

import (
	"cognitivexr.at/cogstream/api/format"
	"cognitivexr.at/cogstream/pkg/pipeline"
	"context"
	"gocv.io/x/gocv"
	"image"
	"log"
	"reflect"
)

type Function func(src gocv.Mat, dst *gocv.Mat)

var NoTransform Function = func(src gocv.Mat, dst *gocv.Mat) {}

// Transform is a binding between Function and pipeline.Transformer.
func (f Function) Transform(ctx context.Context, src *pipeline.Frame, dest pipeline.FrameWriter) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	f(*src.Mat, src.Mat)
	return dest.WriteFrame(src)
}

func funcEqual(a, b interface{}) bool {
	av := reflect.ValueOf(&a).Elem()
	bv := reflect.ValueOf(&b).Elem()
	return av.InterfaceData() == bv.InterfaceData()
}

func IsNoTransform(fn Function) bool {
	return funcEqual(NoTransform, fn)
}

func Pipeline(fns ...Function) Function {
	if fns == nil {
		return NoTransform
	}

	// filter NoTransform functions
	pipeline := make([]Function, 0)
	for _, fn := range fns {
		if IsNoTransform(fn) {
			continue
		}
		pipeline = append(pipeline, fn)
	}

	if len(pipeline) == 0 {
		return NoTransform
	}
	if len(pipeline) == 1 {
		return pipeline[0]
	}

	return func(src gocv.Mat, dst *gocv.Mat) {
		src.CopyTo(dst)

		for _, fn := range fns {
			fn(*dst, dst)
		}
	}
}

// TODO: specify scaling type if necessary (letterbox, scale, ...)
func BuildTransformer(source format.Format, target format.Format) (Function, error) {
	log.Printf("building transformer for format %v -> %v", source, target)

	if target == format.AnyFormat {
		return NoTransform, nil
	}

	fns := make([]Function, 0)

	// TODO: could optimize by first checking whether target size < source size, and if so do the scaling first

	// color conversion
	col, err := GetColorTransform(source.ColorMode, target.ColorMode)
	if err != nil {
		return nil, err
	}
	fns = append(fns, col)

	// orientation
	rotate, flip := GetOrientationTransformation(source.Orientation, target.Orientation)
	fns = append(fns, rotate, flip)

	// resize
	if !(source.Width == target.Width && source.Height == target.Height) {
		scale := ResizeWithScale(image.Pt(target.Width, target.Height))
		fns = append(fns, scale)
	}

	return Pipeline(fns...), nil
}
