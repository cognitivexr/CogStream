package transform

import (
	"gocv.io/x/gocv"
	"testing"
)

func TestIsNoTransform(t *testing.T) {
	ref := NoTransform

	type args struct {
		fn Function
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "equal", args: args{NoTransform}, want: true},
		{name: "withReference", args: args{ref}, want: true},
		{name: "notEqual", args: args{Flip(0)}, want: false},
		{name: "emptyNotEqual", args: args{func(src gocv.Mat, dst *gocv.Mat) {}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNoTransform(tt.args.fn); got != tt.want {
				t.Errorf("IsNoTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}
