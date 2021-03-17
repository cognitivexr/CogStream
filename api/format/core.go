package format

type ColorMode int

const (
	RGB ColorMode = iota + 1
	RGBA
	Gray
	BGR
	BGRA
	HLS
	Lab
	Luv
	Bayer
)

var ColorModes []ColorMode = []ColorMode{RGB, RGBA, Gray, BGR, BGRA, HLS, Lab, Luv, Bayer}

// Orientation encodes the frame oriantation in an EXIF orientation code
type Orientation int

const (
	// see https://sirv.com/help/articles/rotate-photos-to-be-upright/

	TopLeft     Orientation = iota + 1 // = 0 degrees: the correct orientation, no adjustment is required.
	TopRight                           // = 0 degrees, mirrored: image has been flipped back-to-front.
	BottomRight                        // = 180 degrees: image is upside down.
	BottomLeft                         // = 180 degrees, mirrored: image has been flipped back-to-front and is upside down.
	LeftTop                            // = 90 degrees: image has been flipped back-to-front and is on its side.
	RightTop                           // = 90 degrees, mirrored: image is on its side.
	RightBottom                        // = 270 degrees: image has been flipped back-to-front and is on its far side.
	LeftBottom                         // = 270 degrees, mirrored: image is on its far side.
)

var Orientations []Orientation = []Orientation{TopLeft, TopRight, BottomRight, BottomLeft, LeftTop, RightTop, RightBottom, LeftBottom}

func (o Orientation) Angle() int {
	switch o {
	case TopLeft:
		return 0
	case TopRight:
		return 0
	case BottomRight:
		return 180
	case BottomLeft:
		return 180
	case LeftTop:
		return 90
	case RightTop:
		return 90
	case RightBottom:
		return 270
	case LeftBottom:
		return 270
	}
	return 0
}

func (o Orientation) Mirrored() bool {
	switch o {
	case TopLeft:
		return false
	case TopRight:
		return true
	case BottomRight:
		return false
	case BottomLeft:
		return true
	case LeftTop:
		return true
	case RightTop:
		return false
	case RightBottom:
		return true
	case LeftBottom:
		return false
	}
	return false
}

type Format struct {
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	ColorMode   ColorMode   `json:"colorMode"`
	Orientation Orientation `json:"orientation"`
}

// AnyFormat indicates that the any format is supported or accepted.
var AnyFormat = Format{}
