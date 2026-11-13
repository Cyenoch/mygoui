package mygoui

import (
	"image/color"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Above MyGo's effect layer limit a GPU frame paints the viewport without
// its mask, so the shadow has to stand in there as it does without a GPU.
// The forced flag makes the headless frame behave as a GPU frame.
func TestScrollMaskFadeShadesAboveTheLayerLimit(t *testing.T) {
	content, background := ui.RGB(30, 70, 230), ui.RGB(230, 40, 70)
	for _, tc := range []struct {
		name   string
		width  float32
		shaded bool
	}{{"fits", 120, false}, {"oversized", 9000, true}} {
		t.Run(tc.name, func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Box(c).Fill().Background(background).Children(func() {
					viewport := ui.ScrollBoth(c).Absolute().Left(0).Top(0).Size(tc.width, 100).TrackScroll(&state)
					viewport.Children(func() {
						ui.Box(c).Size(tc.width, 300).Shrink(0).Background(content)
					})
					ScrollFade(c, viewport, &state, ScrollFadeOptions{Mode: MaskFade, Edges: FadeBottom})
					ui.Local(viewport, scrollMaskKey{}, newScrollMask).forced = true
				})
			}, 120, 100)
			tt.Frame()
			img := tt.Image()
			want := color.RGBA{R: content.R, G: content.G, B: content.B, A: 255}
			if got := img.RGBAAt(40, 40); got != want {
				t.Fatalf("content away from the edge = %v, want it unchanged", got)
			}
			edge := img.RGBAAt(40, 99)
			if got := edge.R < 100 && edge != want; got != tc.shaded {
				t.Fatalf("bottom edge = %v, shaded content = %v, want %v", edge, got, tc.shaded)
			}
		})
	}
}
