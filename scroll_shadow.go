package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ScrollShadowOrientation is the scroll axis; vertical is the HeroUI default.
type ScrollShadowOrientation uint8

const (
	ScrollShadowVertical ScrollShadowOrientation = iota
	ScrollShadowHorizontal
)

// ScrollShadowOptions describes an opaque scrolling surface. HeroUI fades the
// content with a CSS mask; MyGo has no alpha mask, so this paints a surface
// colored edge fade instead. Surface should be the actual opaque background.
type ScrollShadowOptions struct {
	Orientation ScrollShadowOrientation
	Surface     ui.Color
	Size        float32
	Height      float32
	Label       string
}

// ScrollShadow wraps caller-owned content in a native scroll viewport and
// fades both edges in proportion to the remaining scroll distance.
func ScrollShadow(c *ui.Context, s *ui.ScrollState, o ScrollShadowOptions, body func()) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	surface := o.Surface
	if surface.A == 0 {
		surface = theme.Colors(c).Background
	}
	depth := o.Size
	if depth <= 0 {
		depth = theme.Space40
	}
	var e ui.Element
	if o.Orientation == ScrollShadowHorizontal {
		e = ui.ScrollHorizontal(c).TrackScroll(s)
	} else {
		e = ui.Scroll(c).TrackScroll(s)
	}
	if o.Height > 0 {
		e.Height(o.Height)
	}
	e.MinHeight(0).Label(o.Label).Background(surface)
	e.Children(body)
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if o.Orientation == ScrollShadowHorizontal {
			if s.MaxX <= 0 {
				return
			}
			at := max(0, min(s.X, s.MaxX))
			w := min(depth, r.W/2)
			if start := min(float32(1), at/w); start > 0 {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}, ui.LinearGradient{
					From: surface.Alpha(start), To: surface.Alpha(0), Angle: 90,
				}, 0)
			}
			if end := min(float32(1), (s.MaxX-at)/w); end > 0 {
				p.FillGradient(ui.Rect{X: r.X + r.W - w, Y: r.Y, W: w, H: r.H}, ui.LinearGradient{
					From: surface.Alpha(0), To: surface.Alpha(end), Angle: 90,
				}, 0)
			}
			return
		}
		if s.MaxY <= 0 {
			return
		}
		at := max(0, min(s.Y, s.MaxY))
		h := min(depth, r.H/2)
		if start := min(float32(1), at/h); start > 0 {
			p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: h}, ui.LinearGradient{
				From: surface.Alpha(start), To: surface.Alpha(0), Angle: 180,
			}, 0)
		}
		if end := min(float32(1), (s.MaxY-at)/h); end > 0 {
			p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, ui.LinearGradient{
				From: surface.Alpha(0), To: surface.Alpha(end), Angle: 180,
			}, 0)
		}
	})
	return e
}
