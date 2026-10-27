package mygoui

import (
	"sync"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// ScrollFadeMode selects a native edge treatment.
type ScrollFadeMode uint8

const (
	ShadowFade ScrollFadeMode = iota
	SurfaceFade
	// MaskFade fades the viewport and its descendants to transparency.
	MaskFade
)

// ScrollFadeEdges selects physical viewport edges; zero enables every edge.
type ScrollFadeEdges uint8

const (
	FadeTop ScrollFadeEdges = 1 << iota
	FadeRight
	FadeBottom
	FadeLeft
)

// ScrollFadeOptions customizes an edge fade independently of scrolling.
// SurfaceFade requires the actual opaque surface. Color customizes a shadow.
// MaskFade uses a subtree effect and needs neither Surface nor Color.
type ScrollFadeOptions struct {
	Mode    ScrollFadeMode
	Edges   ScrollFadeEdges
	Size    float32
	Surface ui.Color
	Color   ui.Color
}

// ScrollFade decorates any native viewport tracking state, including existing
// application scroll views. It preserves native hit testing and adds no timer.
// MaskFade uses an effect on the viewport, replacing any existing effect;
// MyGo draws oversized effect layers normally and excludes embedded web views.
func ScrollFade(c *ui.Context, viewport ui.Element, state *ui.ScrollState, o ScrollFadeOptions) ui.Element {
	if state == nil {
		panic("mygoui.ScrollFade: nil scroll state")
	}
	if o.Mode != ShadowFade && o.Mode != SurfaceFade && o.Mode != MaskFade {
		panic("mygoui.ScrollFade: unsupported mode")
	}
	depth := o.Size
	if depth <= 0 {
		depth = theme.Space24
	}
	edges := o.Edges
	if edges == 0 {
		edges = FadeTop | FadeRight | FadeBottom | FadeLeft
	}
	if o.Mode == MaskFade {
		fx := ui.Local(viewport, scrollMaskKey{}, newScrollMask)
		fx.state, fx.depth, fx.edges = state, depth, edges
		fx.rtl = c.Direction() == ui.RightToLeft
		return viewport.Effect(fx)
	}
	ink := o.Color
	if o.Mode == SurfaceFade {
		if o.Surface.A != 255 {
			panic("mygoui.ScrollFade: SurfaceFade requires an opaque surface")
		}
		ink = o.Surface
	} else if ink.A == 0 {
		ink = theme.Colors(c).Foreground.Alpha(0.18)
	}
	rtl := c.Direction() == ui.RightToLeft
	// Absolute children stay fixed in a scroll viewport and paint inside its
	// rounded content clip, before its border and native scrollbars.
	viewport.Clip().Children(func() {
		ui.Box(c.Key(scrollFadeLayerKey{})).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			PassThrough().Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			x, y := max(0, min(state.X, state.MaxX)), max(0, min(state.Y, state.MaxY))
			if rtl {
				x = state.MaxX - x
			}
			paint := func(edge ScrollFadeEdges, remaining float32, box ui.Rect, angle float32, reverse bool) {
				if edges&edge == 0 || remaining <= 0 {
					return
				}
				length := box.H
				if edge == FadeLeft || edge == FadeRight {
					length = box.W
				}
				strength := min(float32(1), remaining/length)
				from, to := ink.Alpha(strength), ink.Alpha(0)
				if reverse {
					from, to = to, from
				}
				p.FillGradient(box, ui.LinearGradient{From: from, To: to, Angle: angle}, 0)
			}
			w, h := min(depth, r.W/2), min(depth, r.H/2)
			paint(FadeTop, y, ui.Rect{X: r.X, Y: r.Y, W: r.W, H: h}, 180, false)
			paint(FadeBottom, state.MaxY-y, ui.Rect{X: r.X, Y: r.Y + r.H - h, W: r.W, H: h}, 180, true)
			paint(FadeLeft, x, ui.Rect{X: r.X, Y: r.Y, W: w, H: r.H}, 90, false)
			paint(FadeRight, state.MaxX-x, ui.Rect{X: r.X + r.W - w, Y: r.Y, W: w, H: r.H}, 90, true)
		})
	})
	return viewport
}

type scrollFadeLayerKey struct{}

type scrollMaskKey struct{}

var scrollMaskShader = sync.OnceValue(func() *shader.Shader {
	return shader.MustCompile(`package main

var Edges vec4

func ColorEffect(pos vec2, color vec4) vec4 {
	size := rectSize()
	alpha := 1.0
	if Edges.x > 0.0 {
		alpha *= clamp(pos.y / Edges.x, 0.0, 1.0)
	}
	if Edges.y > 0.0 {
		alpha *= clamp((size.x - pos.x) / Edges.y, 0.0, 1.0)
	}
	if Edges.z > 0.0 {
		alpha *= clamp((size.y - pos.y) / Edges.z, 0.0, 1.0)
	}
	if Edges.w > 0.0 {
		alpha *= clamp(pos.x / Edges.w, 0.0, 1.0)
	}
	return color * alpha
}
`)
})

type scrollMask struct {
	state   *ui.ScrollState
	values  *shader.Values
	uniform shader.Uniform
	effect  shader.Effect
	depth   float32
	edges   ScrollFadeEdges
	rtl     bool
}

func newScrollMask() scrollMask {
	s := scrollMaskShader()
	values := s.NewValues()
	return scrollMask{
		values: values, uniform: s.Uniform("Edges"),
		effect: s.Effect(shader.Options{Values: values}),
	}
}

func (fx *scrollMask) EffectMargin() ui.Point { return ui.Point{} }

func (fx *scrollMask) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	x := max(0, min(fx.state.X, fx.state.MaxX))
	y := max(0, min(fx.state.Y, fx.state.MaxY))
	if fx.rtl {
		x = fx.state.MaxX - x
	}
	length := func(edge ScrollFadeEdges, remaining float32) float32 {
		if fx.edges&edge == 0 {
			return 0
		}
		return max(0, min(fx.depth, remaining))
	}
	fx.values.Vec4(fx.uniform,
		length(FadeTop, y), length(FadeRight, fx.state.MaxX-x),
		length(FadeBottom, fx.state.MaxY-y), length(FadeLeft, x))
	fx.effect.PaintEffect(p, box, content, inset)
}

// ScrollAxis selects the axes of a native scroll area.
type ScrollAxis uint8

const (
	ScrollVertical ScrollAxis = iota
	ScrollHorizontal
	ScrollBoth
)

// ScrollAreaOptions configures coss's scroll viewport with native scrollbars.
type ScrollAreaOptions struct {
	Axis   ScrollAxis
	Height float32
	Label  string
	Fade   *ScrollFadeOptions
}

// ScrollArea retains native scrolling, focus, wheel and scrollbar behavior.
func ScrollArea(c *ui.Context, state *ui.ScrollState, o ScrollAreaOptions, body func()) ui.Element {
	if state == nil {
		panic("mygoui.ScrollArea: nil scroll state")
	}
	var e ui.Element
	switch o.Axis {
	case ScrollHorizontal:
		e = ui.ScrollHorizontal(c)
	case ScrollBoth:
		e = ui.ScrollBoth(c)
	case ScrollVertical:
		e = ui.Scroll(c)
	default:
		panic("mygoui.ScrollArea: unsupported scroll axis")
	}
	e.TrackScroll(state).MinWidth(0).MinHeight(0).Label(o.Label).FocusRing(true)
	if o.Height > 0 {
		e.Height(o.Height)
	}
	e.Children(body)
	if o.Fade != nil {
		e = ScrollFade(c, e, state, *o.Fade)
	}
	return e
}
