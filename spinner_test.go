package mygoui_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	myui "github.com/egoist/mygo/ui"
)

func TestSpinnerPaintCadenceDoesNotRebuildView(t *testing.T) {
	builds := 0
	tt := myui.NewTester(func(c *myui.Context) {
		builds++
		mygoui.Spinner(c, mygoui.SpinnerOptions{Label: "Syncing"})
	}, 120, 100)
	before := builds
	pixels := bytes.Clone(tt.Image().Pix)
	tt.Advance(time.Second)
	if bytes.Equal(pixels, tt.Image().Pix) {
		t.Fatal("visible spinner did not animate")
	}
	if builds != before {
		t.Fatalf("spinner repaint rebuilt view: before=%d after=%d", before, builds)
	}
}

func TestReducedMotionLoadingIndicatorsSettle(t *testing.T) {
	tt := myui.NewTester(func(c *myui.Context) {
		myui.Column(c).Children(func() {
			mygoui.Spinner(c, mygoui.SpinnerOptions{Label: "Syncing"})
			mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 100, Height: 16})
		})
	}, 160, 100)
	tt.SetPreferences(myui.Preferences{TextScale: 1, ReduceMotion: true})
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Fatal("reduced-motion loading indicators request continuous frames")
	}
}
