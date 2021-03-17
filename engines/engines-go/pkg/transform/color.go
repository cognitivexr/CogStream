package transform

import (
	"cognitivexr.at/cogstream/api/format"
	"fmt"
	"gocv.io/x/gocv"
)

func CvtCol(code gocv.ColorConversionCode) Function {
	return func(src gocv.Mat, dst *gocv.Mat) { gocv.CvtColor(src, dst, code) }
}

func GetColorTransform(from format.ColorMode, to format.ColorMode) (Function, error) {
	if from == to {
		return NoTransform, nil
	}
	if from <= 0 || to <= 0 {
		return NoTransform, nil
	}

	cvtColCode := GetColorConversionCode(from, to)
	if cvtColCode < 0 {
		// FIXME: could do an intermediary transform to RGB (that almost always works) and back
		return nil, fmt.Errorf("no direct color conversion from %s to %s", from, to)
	}
	return CvtCol(cvtColCode), nil
}

func GetColorConversionCode(from format.ColorMode, to format.ColorMode) gocv.ColorConversionCode {
	if from == format.BGR && to == format.BGRA {
		return gocv.ColorBGRToBGRA
	}
	if from == format.BGRA && to == format.BGR {
		return gocv.ColorBGRAToBGR
	}
	if from == format.BGR && to == format.RGBA {
		return gocv.ColorBGRToRGBA
	}
	if from == format.RGBA && to == format.BGR {
		return gocv.ColorRGBAToBGR
	}
	if from == format.BGR && to == format.RGB {
		return gocv.ColorBGRToRGB
	}
	if from == format.BGRA && to == format.RGBA {
		return gocv.ColorBGRAToRGBA
	}
	if from == format.BGR && to == format.Gray {
		return gocv.ColorBGRToGray
	}
	if from == format.RGB && to == format.Gray {
		return gocv.ColorRGBToGray
	}
	if from == format.Gray && to == format.BGR {
		return gocv.ColorGrayToBGR
	}
	if from == format.Gray && to == format.BGRA {
		return gocv.ColorGrayToBGRA
	}
	if from == format.BGRA && to == format.Gray {
		return gocv.ColorBGRAToGray
	}
	if from == format.RGBA && to == format.Gray {
		return gocv.ColorRGBAToGray
	}
	if from == format.BGR && to == format.Lab {
		return gocv.ColorBGRToLab
	}
	if from == format.RGB && to == format.Lab {
		return gocv.ColorRGBToLab
	}
	if from == format.BGR && to == format.Luv {
		return gocv.ColorBGRToLuv
	}
	if from == format.RGB && to == format.Luv {
		return gocv.ColorRGBToLuv
	}
	if from == format.BGR && to == format.HLS {
		return gocv.ColorBGRToHLS
	}
	if from == format.RGB && to == format.HLS {
		return gocv.ColorRGBToHLS
	}
	if from == format.Lab && to == format.BGR {
		return gocv.ColorLabToBGR
	}
	if from == format.Lab && to == format.RGB {
		return gocv.ColorLabToRGB
	}
	if from == format.Luv && to == format.BGR {
		return gocv.ColorLuvToBGR
	}
	if from == format.Luv && to == format.RGB {
		return gocv.ColorLuvToRGB
	}
	if from == format.HLS && to == format.BGR {
		return gocv.ColorHLSToBGR
	}
	if from == format.HLS && to == format.RGB {
		return gocv.ColorHLSToRGB
	}

	return -1
}
