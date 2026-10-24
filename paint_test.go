package mygoui_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func pixelAt(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

func near(got uint8, want uint8) bool {
	d := int(got) - int(want)
	return d >= -3 && d <= 3
}

func TestButtonPaintsVariantAndHover(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Padding(20).Background(colors.Background).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Paint", OnClick: func() {}})
		})
	}, 240, 120)
	img := tt.Image()
	rect, ok := tt.Find("Paint")
	if !ok {
		t.Fatal("button label missing")
	}
	x, y := int(rect.X)+10, int(rect.Y+rect.H/2)
	r, g, b := pixelAt(img, x, y)
	if !near(r, colors.Accent.R) || !near(g, colors.Accent.G) || !near(b, colors.Accent.B) {
		t.Fatalf("button fill = %d,%d,%d, want accent %d,%d,%d", r, g, b,
			colors.Accent.R, colors.Accent.G, colors.Accent.B)
	}
	tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Frame()
	img = tt.Image()
	r, g, b = pixelAt(img, x, y)
	if !near(r, colors.AccentHover.R) || !near(g, colors.AccentHover.G) || !near(b, colors.AccentHover.B) {
		t.Fatalf("hover fill = %d,%d,%d, want accent hover %d,%d,%d", r, g, b,
			colors.AccentHover.R, colors.AccentHover.G, colors.AccentHover.B)
	}
}

func TestCardPaintsSurfaceAndRadius(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Padding(20).Background(colors.Background).Children(func() {
			mygoui.Card(c, mygoui.CardOptions{}, func() {
				mygoui.CardTitle(c, "Surface")
			})
		})
	}, 320, 200)
	img := tt.Image()
	rect, ok := tt.Find("Surface")
	if !ok {
		t.Fatal("card title missing")
	}
	x, y := int(rect.X)-12, int(rect.Y+rect.H/2)
	r, g, b := pixelAt(img, x, y)
	if !near(r, colors.Surface.R) || !near(g, colors.Surface.G) || !near(b, colors.Surface.B) {
		t.Fatalf("card fill = %d,%d,%d, want surface %d,%d,%d", r, g, b,
			colors.Surface.R, colors.Surface.G, colors.Surface.B)
	}
	if r, g, b := pixelAt(img, 22, 22); !near(r, colors.Background.R) ||
		!near(g, colors.Background.G) || !near(b, colors.Background.B) {
		t.Fatalf("corner = %d,%d,%d, want background %d,%d,%d", r, g, b,
			colors.Background.R, colors.Background.G, colors.Background.B)
	}
}
