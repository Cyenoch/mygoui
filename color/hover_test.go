package color

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func colorHoverPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

func colorHoverNear(got, want uint8) bool {
	d := int(got) - int(want)
	return d >= -3 && d <= 3
}

func wantColorHoverPixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := colorHoverPixel(img, x, y)
	if !colorHoverNear(r, c.R) || !colorHoverNear(g, c.G) || !colorHoverNear(b, c.B) {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

// TestColorAreaThumbPressed checks the area's pressed snapshot: a press grows
// the thumb from 16 to 20 DIP, exposing its white border. The control defines no
// hover face for the area or its thumb.
func TestColorAreaThumbPressed(t *testing.T) {
	var colors theme.Palette
	state := &ColorState{}
	state.Set(ui.RGB(128, 64, 32))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			ColorArea(c, ColorAreaOptions{Label: "Fill", State: state, Size: 200})
		})
	}, 300, 300)

	area, ok := tt.Find("Fill")
	if !ok {
		t.Fatal("color area missing")
	}
	cx := area.X + area.W*float32(state.Saturation())
	cy := area.Y + area.H*float32(1-state.Brightness())
	// 8 DIP from the centre is outside the 16 DIP resting thumb and inside
	// the 20 DIP pressed one.
	x, y := int(cx)+8, int(cy)
	beforeR, beforeG, beforeB := colorHoverPixel(tt.Image(), x, y)
	white := ui.RGB(255, 255, 255)
	if colorHoverNear(beforeR, white.R) && colorHoverNear(beforeG, white.G) &&
		colorHoverNear(beforeB, white.B) {
		t.Fatal("the resting thumb already covers the pressed sample")
	}

	notWhite := func(what string) {
		t.Helper()
		r, g, b := colorHoverPixel(tt.Image(), x, y)
		if colorHoverNear(r, white.R) && colorHoverNear(g, white.G) && colorHoverNear(b, white.B) {
			t.Fatalf("%s: the thumb already covers the pressed sample", what)
		}
	}
	tt.Press(cx, cy)
	tt.Frame()
	// The pose starts from the resting size and grows in paint-only frames.
	notWhite("before the grow")
	tt.Advance(160 * time.Millisecond)
	wantColorHoverPixel(t, tt.Image(), x, y, white, "thumb pressed border")
	tt.Release(cx, cy)
	tt.Advance(160 * time.Millisecond)
	notWhite("after release")
}
