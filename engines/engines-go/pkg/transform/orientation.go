package transform

import (
	"cognitivexr.at/cogstream/api/format"
	"gocv.io/x/gocv"
)

func Flip(flipCode int) Function {
	return func(src gocv.Mat, dst *gocv.Mat) { gocv.Flip(src, dst, flipCode) }
}

func Rotate(code gocv.RotateFlag) Function {
	return func(src gocv.Mat, dst *gocv.Mat) { gocv.Rotate(src, dst, code) }
}

// clockwiseRotations returns the number of 90 degree clockwise rotations that that are necessary to transform
// the orientation `from int` to `to int`.
func getRotateFlag(from int, to int) gocv.RotateFlag {
	angle := to - from

	if angle < 0 {
		angle += 360
	}

	rotations := (angle % 360) / 90

	switch rotations {
	case 1:
		return gocv.Rotate90Clockwise
	case 2:
		return gocv.Rotate180Clockwise
	case 3:
		return gocv.Rotate90CounterClockwise
	}

	return -1
}

func GetOrientationTransformation(source format.Orientation, target format.Orientation) (rotate Function, flip Function) {
	rotate, flip = NoTransform, NoTransform

	if source == target {
		return
	}

	// TODO: could optimize, for example rotate 180+flip vertical = flip horizontal

	flag := getRotateFlag(source.Angle(), target.Angle())
	if flag != -1 {
		rotate = Rotate(flag)
	}

	if source.Mirrored() != target.Mirrored() {
		a := ((target.Angle() % 360) / 90) % 2
		if a == 0 { // 0 or 180 degrees -> flip vertically
			flip = Flip(1)
		} else { // 90 or 270 degrees -> flip horizontally
			flip = Flip(0)
		}
	}

	return
}
