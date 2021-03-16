package messages

import (
	"cognitivexr.at/cogstream/api/format"
	"fmt"
	"strconv"
)

func (a Attributes) getInt(key string) (int, bool, error) {
	val := a.Get(key)
	if val == "" {
		return 0, false, nil
	}
	n, err := strconv.Atoi(val)
	return n, true, err
}

func FormatFromAttributes(a Attributes) (f format.Format, err error) {
	var intVal int
	var ok bool

	f.Width, _, err = a.getInt("format.width")
	if err != nil {
		return
	}

	f.Height, _, err = a.getInt("format.height")
	if err != nil {
		return
	}

	colorModeVal := a.Get("format.colorMode")
	if colorModeVal != "" {
		if intVal, err = strconv.Atoi(colorModeVal); err == nil {
			f.ColorMode = format.ColorMode(intVal)
		} else {
			f.ColorMode = format.StringToColorMode(colorModeVal)

			if f.ColorMode == 0 {
				err = fmt.Errorf("unknown color mode: %s", colorModeVal)
				return
			}
		}
	}

	intVal, ok, err = a.getInt("format.orientation")
	if err != nil {
		return
	}
	if ok {
		f.Orientation = format.Orientation(intVal)
	}

	return
}

func FormatToAttributes(f format.Format, a Attributes) {
	a.Set("format.width", strconv.Itoa(f.Width))
	a.Set("format.height", strconv.Itoa(f.Height))
	a.Set("format.colorMode", strconv.Itoa(int(f.ColorMode)))
	a.Set("format.orientation", strconv.Itoa(int(f.Orientation)))
}
