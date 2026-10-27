package color

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestColorStateClampsChannels covers the channel boundaries a caller can hit
// through a slider, a field or a direct setter.
func TestColorStateClampsChannels(t *testing.T) {
	s := NewColorState(ui.RGBA(0x33, 0x66, 0x99, 0.5))
	s.SetSaturation(2)
	if s.Saturation() != 1 {
		t.Fatalf("saturation = %v, want 1", s.Saturation())
	}
	s.SetBrightness(-1)
	if s.Brightness() != 0 {
		t.Fatalf("brightness = %v, want 0", s.Brightness())
	}
	s.SetBrightness(2)
	if s.Brightness() != 1 {
		t.Fatalf("brightness = %v, want 1", s.Brightness())
	}

	s.SetChannel(ChannelRed, 300)
	if s.Value.R != 255 {
		t.Fatalf("red = %d, want 255", s.Value.R)
	}
	s.SetChannel(ChannelRed, -10)
	if s.Value.R != 0 {
		t.Fatalf("red = %d, want 0", s.Value.R)
	}
	s.SetChannel(ChannelGreen, 300)
	if s.Value.G != 255 {
		t.Fatalf("green = %d, want 255", s.Value.G)
	}
	s.SetChannel(ChannelHue, 400)
	if s.Hue() != 360 || s.Value.R != 255 || s.Value.G != 0 {
		t.Fatalf("hue = %v color %v, want 360 and red", s.Hue(), s.Value)
	}

	s.SetAlpha(2)
	if s.Alpha() != 1 || s.Value.A != 255 {
		t.Fatalf("alpha = %v/%d, want 1/255", s.Alpha(), s.Value.A)
	}
	alpha := NewColorState(ui.Hex("#0066ff"))
	alpha.SetAlpha(-1)
	if alpha.Alpha() != 0 || FormatHex(alpha.Value) != "#0066ff00" {
		t.Fatalf("alpha = %v hex %q, want 0 and #0066ff00", alpha.Alpha(), FormatHex(alpha.Value))
	}
}

// TestColorSliderPointerClampsToEnds checks that a press outside the cap
// centers still reaches the channel's bounds.
func TestColorSliderPointerClampsToEnds(t *testing.T) {
	s := NewColorState(ui.RGB(128, 80, 40))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ColorSlider(c, ColorSliderOptions{Label: "Red", State: s, Channel: ChannelRed})
	}, 320, 100)

	tt.ClickAt(1, 34)
	if s.Value.R != 0 {
		t.Fatalf("left edge picked red %d, want 0", s.Value.R)
	}
	tt.ClickAt(319, 34)
	if s.Value.R != 255 {
		t.Fatalf("right edge picked red %d, want 255", s.Value.R)
	}
}

// TestColorSliderCapsPaintChannelEnds checks the solid end caps the control draws
// outside the gradient: the unselected red channel's extremes.
func TestColorSliderCapsPaintChannelEnds(t *testing.T) {
	s := NewColorState(ui.RGB(128, 80, 40))
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal := theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Background(pal.Background).Children(func() {
			ColorSlider(c, ColorSliderOptions{Label: "Red", State: s, Channel: ChannelRed})
		})
	}, 320, 100)

	img := tt.Image()
	low := img.RGBAAt(2, 34)
	high := img.RGBAAt(318, 34)
	wantLow, wantHigh := ui.RGB(0, 80, 40), ui.RGB(255, 80, 40)
	if !colorHoverNear(low.R, wantLow.R) || !colorHoverNear(low.G, wantLow.G) ||
		!colorHoverNear(low.B, wantLow.B) {
		t.Fatalf("low cap = %d,%d,%d, want %d,%d,%d", low.R, low.G, low.B,
			wantLow.R, wantLow.G, wantLow.B)
	}
	if !colorHoverNear(high.R, wantHigh.R) || !colorHoverNear(high.G, wantHigh.G) ||
		!colorHoverNear(high.B, wantHigh.B) {
		t.Fatalf("high cap = %d,%d,%d, want %d,%d,%d", high.R, high.G, high.B,
			wantHigh.R, wantHigh.G, wantHigh.B)
	}
}

// pickerCheckInk reports whether the item shows a selection tick.
func pickerCheckInk(tt *ui.Tester, item ui.Rect) bool {
	scan := ui.Rect{
		X: item.X + item.W*0.25, Y: item.Y + item.H*0.25,
		W: item.W * 0.5, H: item.H * 0.5,
	}
	_, _, ok := swatchCheckInk(tt.Image(), scan)
	return ok
}

// TestSwatchPickerCheckmarkAnimatesIn checks the tick's 150 ms arrival: the
// newly chosen swatch shows no ink on the selection frame and its tick once
// the pose settles.
func TestSwatchPickerCheckmarkAnimatesIn(t *testing.T) {
	red, green := ui.Hex("#e0010f"), ui.Hex("#0a64b5")
	value := red
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal := theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(16).Background(pal.Background).
			Children(func() {
				ColorSwatchPicker(c, ColorSwatchPickerOptions{
					Label: "Palette", Colors: []ui.Color{red, green},
					Names: []string{"Red", "Green"}, Value: &value,
				})
			})
	}, 240, 160)

	greenItem, ok := tt.Find("Green")
	if !ok {
		t.Fatal("the green swatch is missing")
	}
	if err := tt.Click("Green"); err != nil {
		t.Fatal(err)
	}
	if value != green {
		t.Fatalf("the picker chose %v, want green", value)
	}
	if pickerCheckInk(tt, greenItem) {
		t.Fatal("the tick is already drawn on the selection frame")
	}
	tt.Advance(200 * time.Millisecond)
	greenItem, _ = tt.Find("Green")
	if !pickerCheckInk(tt, greenItem) {
		t.Fatal("the tick did not arrive")
	}
}

// TestColorPickerPopoverFadesIn checks the panel's native arrival fade: the
// padding is transparent on the opening frame and the overlay surface once the
// transition settles.
func TestColorPickerPopoverFadesIn(t *testing.T) {
	s := NewColorState(ui.Hex("#6366f1"))
	open := false
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(16).Background(pal.Background).
			Children(func() {
				ColorPicker(c, ColorPickerOptions{
					Label: "Fill", State: s, Open: &open,
					HideSwatches: true,
				})
			})
	}, 520, 760)

	if err := tt.Click("Fill"); err != nil {
		t.Fatal(err)
	}
	area, ok := tt.Find("Color area")
	if !ok {
		t.Fatalf("the panel is missing; texts: %v", tt.Texts())
	}
	// The panel's 8 DIP padding above the area, right of the rounded corner.
	x, y := int(area.X+area.W-6), int(area.Y)-4
	first := tt.Image().RGBAAt(x, y)
	if !colorHoverNear(first.R, pal.Background.R) || !colorHoverNear(first.G, pal.Background.G) ||
		!colorHoverNear(first.B, pal.Background.B) {
		t.Fatalf("the panel is already opaque: %d,%d,%d", first.R, first.G, first.B)
	}
	tt.Advance(200 * time.Millisecond)
	settled := tt.Image().RGBAAt(x, y)
	if !colorHoverNear(settled.R, pal.Overlay.R) || !colorHoverNear(settled.G, pal.Overlay.G) ||
		!colorHoverNear(settled.B, pal.Overlay.B) {
		t.Fatalf("settled panel = %d,%d,%d, want overlay %d,%d,%d",
			settled.R, settled.G, settled.B, pal.Overlay.R, pal.Overlay.G, pal.Overlay.B)
	}
}
