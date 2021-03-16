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
	// see https://sirv.sirv.com/website/exif-orientation-values.jpg

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

type Format struct {
	Width       int         `json:"width"`
	Height      int         `json:"height"`
	ColorMode   ColorMode   `json:"colorMode"`
	Orientation Orientation `json:"orientation"`
}

// AnyFormat indicates that the any format is supported or accepted.
var AnyFormat = Format{}
