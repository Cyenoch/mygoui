package color

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// greenSum adds the green channel of every pixel, a measure of how much of
// the chosen colour the picker shows against the light page.
func greenSum(img *image.RGBA) int {
	n := 0
	for i := 1; i < len(img.Pix); i += 4 {
		n += int(img.Pix[i])
	}
	return n
}

// TestSwatchPickerChoiceBorderFades pins the chosen swatch's border
// arriving through a fade, not at once.
func TestSwatchPickerChoiceBorderFades(t *testing.T) {
	value := ui.Hex("#0a64b5")
	colors := []ui.Color{ui.Hex("#e0010f"), ui.Hex("#0a64b5")}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			ColorSwatchPicker(c, ColorSwatchPickerOptions{
				Label: "Palette", Colors: colors, Names: []string{"Crimson", "Ocean"}, Value: &value,
			})
		})
	}, 240, 120)
	tt.Frame()
	value = colors[0]
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	mid := greenSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	end := greenSum(tt.Image())
	tt2 := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			ColorSwatchPicker(c, ColorSwatchPickerOptions{
				Label: "Palette", Colors: colors, Names: []string{"Crimson", "Ocean"}, Value: &value,
			})
		})
	}, 240, 120)
	tt2.Frame()
	if built := greenSum(tt2.Image()); built != end {
		t.Fatalf("a palette built with a choice shows %d, want %d at once", built, end)
	}
	if mid == end {
		t.Fatal("the border reached its final colour at once")
	}
}

// TestSwatchTickDrawsInWhenSelected pins the marked swatch's tick drawing in.
func TestSwatchTickDrawsInWhenSelected(t *testing.T) {
	selected := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			ColorSwatch(c, ColorSwatchOptions{Color: ui.Hex("#e0010f"), Label: "chip", Selected: selected, Size: SwatchXL})
		})
	}, 160, 160)
	tt.Frame()
	none := greenSum(tt.Image())
	selected = true
	tt.Frame()
	tt.Advance(40 * time.Millisecond)
	mid := greenSum(tt.Image())
	tt.Advance(300 * time.Millisecond)
	end := greenSum(tt.Image())
	if mid <= none || mid >= end {
		t.Fatalf("tick ink: none %d, mid %d, settled %d", none, mid, end)
	}
}
