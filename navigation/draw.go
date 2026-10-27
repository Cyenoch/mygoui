package navigation

import (
	"math"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// drawChevron paints a two-segment chevron pointing right at angle degrees,
// clockwise, around the center of r. It strokes a ui.Path: ui.Painter.Line
// paints only horizontal and vertical strokes, so a diagonal one came out as
// a bar across the indicator. path holds the strokes between paints, so a
// repainting frame rebuilds no command slice.
func drawChevron(p *ui.Painter, path *ui.Path, r ui.Rect, angle float32, width float32, color ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	d := min(r.W, r.H) / 4
	sin, cos := sinCos(angle)
	path.Reset().
		MoveTo(cx-d*cos+2*d*sin, cy-d*sin-2*d*cos).
		LineTo(cx+d*cos, cy+d*sin).
		LineTo(cx-d*cos-2*d*sin, cy-d*sin+2*d*cos)
	p.StrokePath(path, width, color)
}

// drawChevronDown paints a chevron pointing down, turned clockwise by angle
// degrees. It reuses path as drawChevron does.
func drawChevronDown(p *ui.Painter, path *ui.Path, r ui.Rect, angle float32, width float32, color ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	d := min(r.W, r.H) / 4
	sin, cos := sinCos(angle)
	path.Reset().
		MoveTo(cx-d*cos+d/2*sin, cy-d*sin-d/2*cos).
		LineTo(cx-d/2*sin, cy+d/2*cos).
		LineTo(cx+d*cos+d/2*sin, cy+d*sin-d/2*cos)
	p.StrokePath(path, width, color)
}

// drawCheck paints a checkmark in r up to progress of its length.
func drawCheck(p *ui.Painter, path *ui.Path, r ui.Rect, progress float32, color ui.Color) {
	if progress <= 0 {
		return
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	s := min(r.W, r.H)
	path.Reset().
		MoveTo(cx-s*0.30, cy).
		LineTo(cx-s*0.08, cy+s*0.22).
		LineTo(cx+s*0.32, cy-s*0.24)
	if progress < 1 {
		path.Trim(0, progress)
	}
	p.StrokePath(path, 1.5, color)
}

// sinCos returns the sine and cosine of angle degrees, for the rotated
// points of the glyphs.
func sinCos(angle float32) (sin, cos float32) {
	s, c := math.Sincos(float64(angle) * math.Pi / 180)
	return float32(s), float32(c)
}

// overlayPanel applies the popover surface, border, radius, and shadow.
func overlayPanel(c *ui.Context, e ui.Element, radius float32) ui.Element {
	t := theme.Colors(c)
	e = e.Background(t.Popover).TextColor(t.PopoverForeground).Radius(radius).Border(theme.Hairline, t.Border).Clip()
	return style.Shadow(e, c, theme.OverlayShadows(c.Theme().Dark))
}

// scale applies the host's text scale to a DIP metric.
func scale(c *ui.Context, dip float32) float32 { return theme.Scale(c, dip) }

func layoutChanged(a, b ui.Rect) bool {
	for _, delta := range [...]float32{a.X - b.X, a.Y - b.Y, a.W - b.W, a.H - b.H} {
		if delta < -0.5 || delta > 0.5 {
			return true
		}
	}
	return false
}

// hairline is one device-independent pixel.
const hairline = theme.Hairline
