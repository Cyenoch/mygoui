package color

import (
	"sync"

	"github.com/egoist/mygo/ui"
)

// Checker colors of the transparency pattern (16 DIP tile, 8 DIP
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

// bandPathPool reuses the path buffer of a fill with corners of unequal radii.
var bandPathPool = sync.Pool{New: func() any { return &ui.Path{} }}

// fillCorners paints a rectangle with independent corner radii, as CSS's
// border-radius does; Painter.Fill takes one radius for all four corners.
func fillCorners(p *ui.Painter, r ui.Rect, c ui.Color, tl, tr, br, bl float32) {
	if r.W <= 0 || r.H <= 0 {
		return
	}
	limit := min(r.W, r.H) / 2
	tl, tr, br, bl = min(tl, limit), min(tr, limit), min(br, limit), min(bl, limit)
	if tl == tr && tr == br && br == bl {
		p.Fill(r, c, tl)
		return
	}
	const k = 0.5523 // control point distance of a quarter circle
	path := bandPathPool.Get().(*ui.Path)
	path.Reset()
	path.MoveTo(r.X+tl, r.Y)
	path.LineTo(r.X+r.W-tr, r.Y)
	if tr > 0 {
		path.CubeTo(r.X+r.W-tr+tr*k, r.Y, r.X+r.W, r.Y+tr-tr*k, r.X+r.W, r.Y+tr)
	}
	path.LineTo(r.X+r.W, r.Y+r.H-br)
	if br > 0 {
		path.CubeTo(r.X+r.W, r.Y+r.H-br+br*k, r.X+r.W-br+br*k, r.Y+r.H, r.X+r.W-br, r.Y+r.H)
	}
	path.LineTo(r.X+bl, r.Y+r.H)
	if bl > 0 {
		path.CubeTo(r.X+bl-bl*k, r.Y+r.H, r.X, r.Y+r.H-bl+bl*k, r.X, r.Y+r.H-bl)
	}
	path.LineTo(r.X, r.Y+tl)
	if tl > 0 {
		path.CubeTo(r.X, r.Y+tl-tl*k, r.X+tl-tl*k, r.Y, r.X+tl, r.Y)
	}
	path.Close()
	p.FillPath(path, c)
	bandPathPool.Put(path)
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

// paintThumb paints a color thumb: a 1 DIP outer shadow, a three DIP white
// ring around the chosen color and a hairline inside that ring, as the control's
// border and stacked shadows draw it.
func paintThumb(p *ui.Painter, centerX, centerY, diameter float32, color ui.Color) {
	outer := diameter + 2
	p.Fill(ui.Rect{
		X: centerX - outer/2, Y: centerY - outer/2,
		W: outer, H: outer,
	}, ui.RGBA(0, 0, 0, 0.1), outer/2)
	p.Fill(ui.Rect{
		X: centerX - diameter/2, Y: centerY - diameter/2,
		W: diameter, H: diameter,
	}, ui.RGB(255, 255, 255), diameter/2)
	inner := max(4, diameter-6)
	innerRect := ui.Rect{
		X: centerX - inner/2, Y: centerY - inner/2,
		W: inner, H: inner,
	}
	p.Fill(innerRect, color, inner/2)
	p.Stroke(innerRect, ui.RGBA(0, 0, 0, 0.1), inner/2, 1)
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
