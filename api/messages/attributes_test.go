package messages

import (
	"cognitivexr.at/cogstream/api/format"
	"reflect"
	"testing"
)

func TestAttributes_getInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		a       Attributes
		args    args
		wantN   int
		wantOk  bool
		wantErr bool
	}{
		{name: "correct", a: NewAttributes().Set("foo", "1"), args: args{"foo"}, wantN: 1, wantOk: true, wantErr: false},
		{name: "missingKey", a: NewAttributes().Set("foo", "1"), args: args{"bar"}, wantN: 0, wantOk: false, wantErr: false},
		{name: "notAnInt", a: NewAttributes().Set("foo", "a"), args: args{"foo"}, wantN: 0, wantOk: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, gotOk, err := tt.a.getInt(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("getInt() gotN = %v, want %v", gotN, tt.wantN)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getInt() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestFormatFromAttributes(t *testing.T) {
	type args struct {
		a Attributes
	}
	tests := []struct {
		name    string
		args    args
		wantF   format.Format
		wantErr bool
	}{
		{
			name: "default",
			args: args{NewAttributes().
				Set("format.width", "640").
				Set("format.height", "480").
				Set("format.colorMode", "1").
				Set("format.orientation", "1")},
			wantF:   format.Format{640, 480, format.RGB, format.TopLeft},
			wantErr: false,
		},
		{
			name: "stringColorMode",
			args: args{NewAttributes().
				Set("format.width", "640").
				Set("format.height", "480").
				Set("format.colorMode", "RGB").
				Set("format.orientation", "1")},
			wantF:   format.Format{640, 480, format.RGB, format.TopLeft},
			wantErr: false,
		},
		{
			name: "stringColorModeLower",
			args: args{NewAttributes().
				Set("format.width", "640").
				Set("format.height", "480").
				Set("format.colorMode", "rgba").
				Set("format.orientation", "1")},
			wantF:   format.Format{640, 480, format.RGBA, format.TopLeft},
			wantErr: false,
		},
		{
			name: "missingAttribute",
			args: args{NewAttributes().
				Set("format.width", "640").
				Set("format.colorMode", "1").
				Set("format.orientation", "1")},
			wantF:   format.Format{640, 0, format.RGB, format.TopLeft},
			wantErr: false,
		},
		{
			name: "formatError",
			args: args{NewAttributes().
				Set("format.width", "640").
				Set("format.height", "abcd").
				Set("format.colorMode", "1").
				Set("format.orientation", "1")},
			wantF:   format.Format{640, 0, 0, 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, err := FormatFromAttributes(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatFromAttributes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("FormatFromAttributes() gotF = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestFormatToAttributes(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		a := NewAttributes()
		f := format.Format{640, 480, format.RGB, format.TopLeft}

		want := NewAttributes().
			Set("format.width", "640").
			Set("format.height", "480").
			Set("format.colorMode", "1").
			Set("format.orientation", "1")

		FormatToAttributes(f, a)

		if !reflect.DeepEqual(a, want) {
			t.Errorf("FormatFromAttributes() gotF = %v, want %v", a, want)
		}
	})

	t.Run("retainsExistingKeys", func(t *testing.T) {
		a := NewAttributes()
		a.Set("foo", "bar")
		f := format.Format{640, 480, format.RGB, format.TopLeft}

		want := NewAttributes().
			Set("foo", "bar").
			Set("format.width", "640").
			Set("format.height", "480").
			Set("format.colorMode", "1").
			Set("format.orientation", "1")

		FormatToAttributes(f, a)

		if !reflect.DeepEqual(a, want) {
			t.Errorf("FormatFromAttributes() gotF = %v, want %v", a, want)
		}
	})
}
