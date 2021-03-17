package transform

import (
	"gocv.io/x/gocv"
	"testing"
)

func Test_getRotateFlag(t *testing.T) {
	type args struct {
		from int
		to   int
	}
	tests := []struct {
		name string
		args args
		want gocv.RotateFlag
	}{
		{"0,90->90", args{0, 90}, gocv.Rotate90Clockwise},
		{"360,90->90", args{360, 90}, gocv.Rotate90Clockwise},
		{"90,180->90", args{90, 180}, gocv.Rotate90Clockwise},
		{"180,270->90", args{180, 270}, gocv.Rotate90Clockwise},
		{"270,0->90", args{270, 0}, gocv.Rotate90Clockwise},

		{"0,180->180", args{0, 180}, gocv.Rotate180Clockwise},
		{"90,270->180", args{90, 270}, gocv.Rotate180Clockwise},
		{"180,0->180", args{180, 0}, gocv.Rotate180Clockwise},
		{"270,90->180", args{270, 90}, gocv.Rotate180Clockwise},

		{"0,270->270", args{0, 270}, gocv.Rotate90CounterClockwise},
		{"360,270->270", args{360, 270}, gocv.Rotate90CounterClockwise},
		{"90,0->270", args{90, 0}, gocv.Rotate90CounterClockwise},
		{"180,90->270", args{180, 90}, gocv.Rotate90CounterClockwise},
		{"270,180->270", args{270, 180}, gocv.Rotate90CounterClockwise},

		{"error", args{1, 2}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRotateFlag(tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("getRotateFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}


