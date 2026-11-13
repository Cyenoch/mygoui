package navigation

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// minTextContrast is the contrast WCAG asks of text.
const minTextContrast = 4.5

// legibleInk returns ink, drawn toward strong just far enough to read at
// minTextContrast on face, which is opaque. The theme's muted ink alone falls
// short on its own muted face in the dark appearance.
func legibleInk(ink, strong, face ui.Color) ui.Color {
	for t := float32(0); t < 1; t += 0.1 {
		if mixed := ink.Mix(strong, t); contrastRatio(mixed, face) >= minTextContrast {
			return mixed
		}
	}
	return strong
}

// contrastRatio is the WCAG contrast of two opaque colors.
func contrastRatio(a, b ui.Color) float64 {
	la, lb := luminance(a), luminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

func luminance(c ui.Color) float64 {
	channel := func(v uint8) float64 {
		x := float64(v) / 255
		if x <= 0.03928 {
			return x / 12.92
		}
		return math.Pow((x+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
}

// readableOn is white or black, whichever reads better on the opaque face of a
// solid status color, whose own foreground is made for its soft tint.
func readableOn(face ui.Color) ui.Color {
	if white := ui.RGB(255, 255, 255); contrastRatio(white, face) >= minTextContrast {
		return white
	}
	return ui.RGB(0, 0, 0)
}
