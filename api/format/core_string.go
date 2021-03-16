package format

import (
	"fmt"
	"strings"
)

func (c ColorMode) String() string {
	switch c {
	case RGB:
		return "RGB"
	case RGBA:
		return "RGBA"
	case Gray:
		return "Gray"
	case BGR:
		return "BGR"
	case BGRA:
		return "BGRA"
	case HLS:
		return "HLS"
	case Lab:
		return "Lab"
	case Luv:
		return "Luv"
	}
	return "unknown"
}

func StringToColorMode(str string) ColorMode {
	lower := strings.ToLower(str)

	for _, colorMode := range ColorModes {
		if lower == strings.ToLower(colorMode.String()) {
			return colorMode
		}
	}

	return 0
}

func (f Format) String() string {
	return fmt.Sprintf("(%d, %d, %s, %d)", f.Width, f.Height, f.ColorMode, f.Orientation)
}
