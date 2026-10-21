package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ScrollEdgesOptions describes the opaque surface behind a scrolling viewport.
// Depth is measured in DIPs; zero uses the shared spacing scale. Axes must match
// the native scroll constructor; zero means ui.Vertical.
type ScrollEdgesOptions struct {
	Surface  ui.Color
	Depth    float32
	Axes     ui.Orientation
	Disabled bool
}

// ScrollEdges reserves native scrollbar gutters and softens hidden content at
// the inner clip edges. Apply it before children to a tracked, square Scroll
// with no existing border or scrollbar insets, inside the containing frame.
// Disabled and high contrast omit only the fade, preserving the content clip.
// The surface must be opaque: this is a surface-colored fade, not an alpha mask.
// It installs DrawOver; do not install a second callback on the same viewport.
func ScrollEdges(c *ui.Context, viewport ui.Element, state *ui.ScrollState, o ScrollEdgesOptions) {
	if state == nil || o.Surface.A != 255 {
		return
	}
	axes := o.Axes
	if axes == 0 {
		axes = ui.Vertical
	}
	barWidth := c.Theme().ScrollbarWidth
	if barWidth <= 0 {
		barWidth = theme.Space6
	}
	// MyGo overlays its bars. Padding scrolls with the content and cannot
	// exclude that content from a bar lane. Transparent borders reserve both
	// layout and the native child clip, while leaving scrollbar paint intact.
	// An opaque border would paint AFTER the native bars and erase them.
	lane := barWidth + theme.Space6
	right, bottom, left := float32(0), float32(0), float32(0)
	rtl := false
	viewport.Children(func() { rtl = c.Direction() == ui.RightToLeft })
	if axes&ui.Vertical != 0 {
		if rtl {
			left = lane
		} else {
			right = lane
		}
	}
	if axes&ui.Horizontal != 0 {
		bottom = lane
	}
	viewport.Background(o.Surface).BorderWidth(0, right, bottom, left).BorderColor(ui.Color{})
	if o.Disabled || c.Preferences().HighContrast {
		return
	}
	depth := o.Depth
	if depth <= 0 {
		depth = theme.Space24
	}
	viewport.DrawOver(func(p *ui.Painter, r ui.Rect) {
		r.X, r.W = r.X+left, max(0, r.W-left-right)
		r.H = max(0, r.H-bottom)
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if h := min(depth, r.H/2); axes&ui.Vertical != 0 && state.MaxY > 0 && h > 0 {
			at := max(0, min(state.Y, state.MaxY))
			if strength := min(float32(1), at/h); strength > 0 {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: h}, ui.LinearGradient{
					From: o.Surface.Alpha(strength), To: o.Surface.Alpha(0), Angle: 180,
				}, 0)
			}
			if strength := min(float32(1), (state.MaxY-at)/h); strength > 0 {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, ui.LinearGradient{
					From: o.Surface.Alpha(0), To: o.Surface.Alpha(strength), Angle: 180,
				}, 0)
			}
		}
		if w := min(depth, r.W/2); axes&ui.Horizontal != 0 && state.MaxX > 0 && w > 0 {
			at := max(0, min(state.X, state.MaxX))
			start, end := min(float32(1), at/w), min(float32(1), (state.MaxX-at)/w)
			if rtl {
				start, end = end, start
			}
			if start > 0 {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}, ui.LinearGradient{
					From: o.Surface.Alpha(start), To: o.Surface.Alpha(0), Angle: 90,
				}, 0)
			}
			if end > 0 {
				p.FillGradient(ui.Rect{X: r.X + r.W - w, Y: r.Y, W: w, H: r.H}, ui.LinearGradient{
					From: o.Surface.Alpha(0), To: o.Surface.Alpha(end), Angle: 90,
				}, 0)
			}
		}
	})
}
