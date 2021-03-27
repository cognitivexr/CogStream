package transform

import (
	"gocv.io/x/gocv"
	"image"
)

func ResizeWithScale(sz image.Point) Function {
	return func(src gocv.Mat, dst *gocv.Mat) {
		gocv.Resize(src, dst, sz, 0, 0, gocv.InterpolationDefault)
	}
}
