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
// a bar across the indicator.
func drawChevron(p *ui.Painter, r ui.Rect, angle float32, width float32, color ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	d := min(r.W, r.H) / 4
	sin, cos := sinCos(angle)
	var path ui.Path
	path.MoveTo(cx-d*cos+2*d*sin, cy-d*sin-2*d*cos).
		LineTo(cx+d*cos, cy+d*sin).
		LineTo(cx-d*cos-2*d*sin, cy-d*sin+2*d*cos)
	p.StrokePath(&path, width, color)
}

// drawChevronDown paints a chevron pointing down, turned clockwise by angle
// degrees.
func drawChevronDown(p *ui.Painter, r ui.Rect, angle float32, width float32, color ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	d := min(r.W, r.H) / 4
	sin, cos := sinCos(angle)
	var path ui.Path
	path.MoveTo(cx-d*cos+d/2*sin, cy-d*sin-d/2*cos).
		LineTo(cx-d/2*sin, cy+d/2*cos).
		LineTo(cx+d*cos+d/2*sin, cy+d*sin-d/2*cos)
	p.StrokePath(&path, width, color)
}

// drawCheck paints HeroUI's checkmark in r.
func drawCheck(p *ui.Painter, r ui.Rect, color ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	s := min(r.W, r.H)
	var path ui.Path
	path.MoveTo(cx-s*0.30, cy).
		LineTo(cx-s*0.08, cy+s*0.22).
		LineTo(cx+s*0.32, cy-s*0.24)
	p.StrokePath(&path, 1.5, color)
}

// sinCos returns the sine and cosine of angle degrees, for the rotated
// points of the glyphs.
func sinCos(angle float32) (sin, cos float32) {
	s, c := math.Sincos(float64(angle) * math.Pi / 180)
	return float32(s), float32(c)
}

// drawDot paints HeroUI's selected-radio dot in r.
func drawDot(p *ui.Painter, r ui.Rect, color ui.Color) {
	s := min(r.W, r.H) / 2
	p.Fill(ui.Rect{X: r.X + r.W/2 - s, Y: r.Y + r.H/2 - s, W: 2 * s, H: 2 * s}, color, s)
}

// overlayPanel applies HeroUI's floating-panel surface: overlay background,
// overlay foreground, 24 DIP corners and the overlay shadow (an inset
// hairline in dark mode, where upstream disables the drop shadow).
func overlayPanel(c *ui.Context, e ui.Element, radius float32) ui.Element {
	t := theme.Colors(c)
	e = e.Background(t.Overlay).TextColor(t.OverlayForeground).Radius(radius).Clip()
	return style.Shadow(e, c, theme.OverlayShadows(c.Theme().Dark))
}

// scale converts an unscaled HeroUI DIP metric.
func scale(c *ui.Context, dip float32) float32 { return theme.Scale(c, dip) }

// hairline is one device-independent pixel.
const hairline = theme.Hairline
