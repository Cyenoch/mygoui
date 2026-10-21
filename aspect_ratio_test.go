package mygoui_test

import (
	"math"
	"testing"

	"github.com/Cyenoch/mygoui"
	myui "github.com/egoist/mygo/ui"
)

func TestAspectRatioTracksContainerResize(t *testing.T) {
	tt := myui.NewTester(func(c *myui.Context) {
		mygoui.AspectRatio(c, 16.0/9, nil).Label("Preview")
	}, 320, 400)
	before, ok := tt.Find("Preview")
	if !ok || math.Abs(float64(before.W/before.H-16.0/9)) > 0.02 {
		t.Fatalf("initial ratio bounds=%v", before)
	}
	tt.SetSize(180, 400)
	after, ok := tt.Find("Preview")
	if !ok || after.W >= before.W || math.Abs(float64(after.W/after.H-16.0/9)) > 0.02 {
		t.Fatalf("resize did not preserve ratio: before=%v after=%v", before, after)
	}
}
