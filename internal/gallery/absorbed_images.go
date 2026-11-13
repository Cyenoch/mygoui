package gallery

import (
	"image"
	"image/color"
	"sync"

	"github.com/egoist/mygo/ui"
)

// demoPicture is one generated bitmap the image demos share.
type demoPicture struct {
	Label  string
	Source ui.ImageSource
}

// demoPictures are three gradients drawn once, so the image demos need no
// files and no network.
var demoPictures = sync.OnceValue(func() []demoPicture {
	return []demoPicture{
		{"Dawn", gradientBitmap(color.RGBA{255, 190, 120, 255}, color.RGBA{214, 80, 118, 255})},
		{"Lagoon", gradientBitmap(color.RGBA{96, 208, 196, 255}, color.RGBA{40, 92, 190, 255})},
		{"Dusk", gradientBitmap(color.RGBA{150, 120, 232, 255}, color.RGBA{28, 30, 72, 255})},
	}
})

// gradientBitmap blends from one corner color to the opposite one.
func gradientBitmap(from, to color.RGBA) *ui.Bitmap {
	const w, h = 480, 300
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	mix := func(a, b uint8, k float64) uint8 { return uint8(float64(a) + (float64(b)-float64(a))*k) }
	for y := range h {
		for x := range w {
			k := (float64(x)/w + float64(y)/h) / 2
			img.SetRGBA(x, y, color.RGBA{mix(from.R, to.R, k), mix(from.G, to.G, k), mix(from.B, to.B, k), 255})
		}
	}
	return ui.NewBitmap(img)
}
