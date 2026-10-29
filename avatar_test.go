package mygoui_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAvatarImageSourcesAndFallback(t *testing.T) {
	pixels := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := range 4 {
		for x := range 4 {
			pixels.SetRGBA(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, pixels); err != nil {
		t.Fatal(err)
	}
	bitmap, err := ui.DecodeImage(pngData.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	svg, err := ui.DecodeImage([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="4" height="4"><path fill="#ff0000" d="M0 0h4v4H0z"/></svg>`))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name     string
		source   ui.ImageSource
		fallback bool
	}{
		{name: "bitmap", source: bitmap},
		{name: "svg", source: svg},
		{name: "absent", fallback: true},
		{name: "nil bitmap", source: (*ui.Bitmap)(nil), fallback: true},
		{name: "nil svg", source: (*ui.SVG)(nil), fallback: true},
		{name: "empty bitmap", source: new(ui.Bitmap), fallback: true},
		{name: "empty svg", source: new(ui.SVG), fallback: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				mygoui.Avatar(c, mygoui.AvatarOptions{Image: tc.source, Name: "Ada Lovelace"})
			}, 100, 80)
			if tt.HasText("AL") != tc.fallback {
				t.Fatalf("fallback visible = %v, want %v", tt.HasText("AL"), tc.fallback)
			}
			bounds := findRect(t, tt, "Ada Lovelace")
			if !tc.fallback {
				wantPixel(t, tt.Image(), int(bounds.X+bounds.W/2), int(bounds.Y+bounds.H/2), ui.RGB(255, 0, 0), "loaded portrait")
			}
		})
	}
}
