package color

import (
	"image"
	"testing"

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

// TestColorSwatchPickerHoverGrowsSwatch checks the picker item's hover
// snapshot: HeroUI scales the swatch to 1.1, so the rest fill becomes the
// color under the pointer and shrinks back when it leaves. The single
// ColorSwatch preview is a RoleImage with no pointer state.
func TestColorSwatchPickerHoverGrowsSwatch(t *testing.T) {
	var colors theme.Palette
	red := ui.RGB(255, 0, 0)
	var chosen ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			ColorSwatchPicker(c, ColorSwatchPickerOptions{
				Label: "Palette", Size: SwatchMD, Value: &chosen,
				Colors: []ui.Color{red, ui.RGB(0, 255, 0)},
				Names:  []string{"Red", "Green"},
			})
		})
	}, 400, 140)

	item, ok := tt.Find("Red")
	if !ok {
		t.Fatal("red swatch missing")
	}
	// One pixel inside the resting circle, exposed when it grows 1.1×.
	x, y := int(item.X)-1, int(item.Y+item.H/2)
	wantColorHoverPixel(t, tt.Image(), x, y, colors.Background, "swatch rest fill")

	tt.Move(item.X+item.W/2, item.Y+item.H/2)
	tt.Frame()
	wantColorHoverPixel(t, tt.Image(), x, y, red, "swatch hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantColorHoverPixel(t, tt.Image(), x, y, colors.Background, "swatch fill after leaving")
}

// TestColorAreaThumbPressed checks the area's pressed snapshot: a press grows
// the thumb from 16 to 20 DIP, exposing its white border. HeroUI defines no
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

	tt.Press(cx, cy)
	tt.Frame()
	wantColorHoverPixel(t, tt.Image(), x, y, white, "thumb pressed border")
	tt.Release(cx, cy)
}
