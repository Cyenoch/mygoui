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

func TestCardPaintsSurfaceAndRadius(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Padding(20).Background(colors.Background).Children(func() {
			mygoui.Card(c, mygoui.CardOptions{}, func() {
				mygoui.CardHeader(c, func() {
					mygoui.CardTitle(c, "Surface")
				})
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

func TestCardSectionsOwnTheirPadding(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(20).Children(func() {
			mygoui.Card(c, mygoui.CardOptions{Label: "Card"}, func() {
				mygoui.CardHeader(c, func() {
					mygoui.CardTitle(c, "Title")
					mygoui.CardDescription(c, "Description")
				}).Label("Header")
				mygoui.CardContent(c, func() {
					ui.Text(c, "Body")
				}).Label("Content")
				mygoui.CardFooter(c, func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Action"})
				}).Label("Footer")
			})
		})
	}, 360, 320)
	card := findRect(t, tt, "Card")
	header := findRect(t, tt, "Header")
	content := findRect(t, tt, "Content")
	footer := findRect(t, tt, "Footer")
	title := findRect(t, tt, "Title")
	description := findRect(t, tt, "Description")
	body := findRect(t, tt, "Body")
	action := findRect(t, tt, "Action")
	if header.X-card.X != 1 || header.Y-card.Y != 1 || header.W != card.W-2 {
		t.Fatalf("card adds padding outside its sections: card=%+v header=%+v", card, header)
	}
	if title.X-header.X != 24 || title.Y-header.Y != 24 ||
		header.Y+header.H-description.Y-description.H != 24 {
		t.Fatalf("header padding: header=%+v title=%+v description=%+v", header, title, description)
	}
	if description.Y-title.Y-title.H != 6 || description.H != 20 {
		t.Fatalf("card text geometry: title=%+v description=%+v", title, description)
	}
	if content.Y != header.Y+header.H || footer.Y != content.Y+content.H {
		t.Fatalf("sections have unexpected outer gaps: %+v %+v %+v", header, content, footer)
	}
	if body.X-content.X != 24 || body.Y != content.Y || body.H != content.H {
		t.Fatalf("content horizontal padding: content=%+v body=%+v", content, body)
	}
	if action.X-footer.X != 24 || action.Y-footer.Y != 24 ||
		footer.Y+footer.H-action.Y-action.H != 24 {
		t.Fatalf("footer padding: footer=%+v action=%+v", footer, action)
	}
}
