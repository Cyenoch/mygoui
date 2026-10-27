package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestRangeSliderPaintAndDragWithOffset(t *testing.T) {
	for _, tc := range []struct {
		name string
		rtl  bool
	}{
		{name: "ltr"},
		{name: "rtl", rtl: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			low, high := 0.0, 50.0
			var colors theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).Fill().Padding(40).Background(colors.Background).Children(func() {
					Slider(c, SliderOptions{
						State: SliderState{Low: &low, High: &high}, Min: 0, Max: 100, HideOutput: true,
					})
				})
			}, 400, 160)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if tc.rtl {
				tt.SetDirection(ui.RightToLeft)
			}

			const left, width, pad, y = float32(40), float32(320), float32(10), float32(50)
			at := func(f float32) float32 {
				if tc.rtl {
					f = 1 - f
				}
				return left + pad + (width-2*pad)*f
			}
			checkPaint := func() {
				t.Helper()
				img := tt.Image()
				wantFormHoverPixel(t, img, int(at(float32(low/100))), int(y),
					ui.RGB(255, 255, 255), "low thumb")
				wantFormHoverPixel(t, img, int(at(float32(high/100))), int(y),
					ui.RGB(255, 255, 255), "high thumb")
				wantFormHoverPixel(t, img, int(at(float32((low+high)/200))), int(y),
					colors.Primary, "selected range")
				if low > 0 {
					wantFormHoverPixel(t, img, int(at(float32(low/200))), int(y),
						colors.Input.Over(colors.Border.Over(colors.Background)), "track before range")
				}
				if high < 100 {
					wantFormHoverPixel(t, img, int(at(float32((high+100)/200))), int(y),
						colors.Input.Over(colors.Border.Over(colors.Background)), "track after range")
				}
			}
			checkPaint()

			// Drag each real native thumb to the opposite half of the track.
			tt.Press(at(0.5), y)
			tt.Move(at(1), y)
			tt.Release(at(1), y)
			if high < 99.9 || low != 0 {
				t.Fatalf("dragging high to its endpoint gave %v..%v", low, high)
			}
			tt.Press(at(0), y)
			tt.Move(at(0.5), y)
			tt.Release(at(0.5), y)
			if low < 49.9 || low > 50.1 || high < 99.9 {
				t.Fatalf("dragging low to the midpoint gave %v..%v", low, high)
			}
			checkPaint()
		})
	}
}
