package mygoui

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestProgressBoundsAndUnknownValues(t *testing.T) {
	for _, test := range []struct {
		name    string
		options ProgressOptions
		want    float64
	}{
		{"default percent", ProgressOptions{Value: 25}, .25},
		{"actual range", ProgressOptions{Value: 30, Min: 20, Max: 60}, .25},
		{"below start", ProgressOptions{Value: -10, Max: 10}, 0},
		{"above end", ProgressOptions{Value: 20, Max: 10}, 1},
		{"unknown", ProgressOptions{Indeterminate: true}, -1},
		{"invalid bounds", ProgressOptions{Min: 20, Max: 10}, -1},
		{"not a number", ProgressOptions{Value: math.NaN()}, -1},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := progressValue(test.options); got != test.want {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestIndeterminateProgressRespectsReducedMotionWithoutFakingCompletion(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Progress(c, ProgressOptions{Label: "Connecting", Indeterminate: true}).Width(240)
	}, 260, 60)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	before := bytes.Clone(tt.Image().Pix)
	tt.Advance(time.Second)
	if !bytes.Equal(before, tt.Image().Pix) || tt.FramePending() {
		t.Fatal("unknown progress continued animating with reduced motion")
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	tt.Advance(200 * time.Millisecond)
	if bytes.Equal(before, tt.Image().Pix) {
		t.Fatal("unknown progress did not animate when motion enabled")
	}
}
