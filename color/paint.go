package color

import (
	"github.com/egoist/mygo/ui"
)

// Checker colors of HeroUI's transparency pattern (16 DIP tile, 8 DIP
// squares).
var (
	checkerLight = ui.RGB(0xf7, 0xf7, 0xf7)
	checkerDark  = ui.RGB(0xef, 0xef, 0xef)
)

// checkers paints the transparency checkerboard inside r. It allocates
// nothing beyond the paint operations themselves.
func checkers(p *ui.Painter, r ui.Rect, size float32) {
	if size <= 0 || r.W <= 0 || r.H <= 0 {
		return
	}
	p.Fill(r, checkerLight, 0)
	for y, row := r.Y, 0; y < r.Y+r.H; y, row = y+size, row+1 {
		for x, col := r.X, 0; x < r.X+r.W; x, col = x+size, col+1 {
			if (row+col)%2 == 0 {
				continue
			}
			p.Fill(ui.Rect{
				X: x, Y: y,
				W: min(size, r.X+r.W-x),
				H: min(size, r.Y+r.H-y),
			}, checkerDark, 0)
		}
	}
}

// paintSwatch paints color over a checkerboard, rounded by radius.
func paintSwatch(p *ui.Painter, r ui.Rect, color ui.Color, radius float32) {
	if color.A < 255 {
		p.Clip(r, radius, func() { checkers(p, r, 8) })
	}
	if color.A != 0 {
		p.Fill(r, color, radius)
	}
}

// paintSwatchBorder paints the hairline that separates a swatch from a
// similarly colored surface.
func paintSwatchBorder(p *ui.Painter, r ui.Rect, radius float32) {
	p.Stroke(r, ui.RGBA(0, 0, 0, 0.1), radius, 1)
}

// paintThumb paints a color thumb: a three DIP white ring around the chosen
// color, with a dark hairline that keeps it visible on light surfaces.
func paintThumb(p *ui.Painter, centerX, centerY, diameter float32, color ui.Color) {
	outer := diameter + 2
	p.Fill(ui.Rect{
		X: centerX - outer/2, Y: centerY - outer/2,
		W: outer, H: outer,
	}, ui.RGBA(0, 0, 0, 0.12), outer/2)
	p.Fill(ui.Rect{
		X: centerX - diameter/2, Y: centerY - diameter/2,
		W: diameter, H: diameter,
	}, ui.RGB(255, 255, 255), diameter/2)
	inner := max(4, diameter-6)
	p.Fill(ui.Rect{
		X: centerX - inner/2, Y: centerY - inner/2,
		W: inner, H: inner,
	}, color, inner/2)
}

// visibleColor returns c composited over the surface it sits on, so a border
// or checkmark keeps its contrast when the color is translucent.
func visibleColor(c ui.Color, over ui.Color) ui.Color {
	if c.A == 255 {
		return c
	}
	return c.Over(over)
}

// luminance returns the perceived brightness of a color, from 0 to 1.
func luminance(c ui.Color) float64 {
	return (0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)) / 255
}

// paintCheckmark draws the selection tick centered in r. It strokes a
// ui.Path: ui.Painter.Line paints only horizontal and vertical strokes, so
// the tick's diagonals became a bar.
func paintCheckmark(p *ui.Painter, r ui.Rect, color ui.Color) {
	s := min(r.W, r.H) / 3
	if s < 6 {
		s = 6
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	var path ui.Path
	path.MoveTo(cx-s*0.5, cy).
		LineTo(cx-s*0.1, cy+s*0.4).
		LineTo(cx+s*0.5, cy-s*0.4)
	p.StrokePath(&path, 2, color)
}
