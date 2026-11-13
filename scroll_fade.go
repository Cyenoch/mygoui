package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/fx"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
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
// MaskFade uses an effect on the viewport, replacing any existing effect.
// Without a GPU, or for a viewport above MyGo's 8192-pixel effect layer, the
// shadow stands in for the mask. Embedded web views do not show through it.
func ScrollFade(c *ui.Context, viewport ui.Element, state *ui.ScrollState, o ScrollFadeOptions) ui.Element {
	if state == nil {
		panic("mygoui.ScrollFade: nil scroll state")
	}
	if o.Mode != ShadowFade && o.Mode != SurfaceFade && o.Mode != MaskFade {
		panic("mygoui.ScrollFade: unsupported mode")
	}
	depth := o.Size
	if depth <= 0 {
		depth = theme.Effects().FadeEdge
	}
	edges := o.Edges
	if edges == 0 {
		edges = FadeTop | FadeRight | FadeBottom | FadeLeft
	}
	ink := o.Color
	if o.Mode == SurfaceFade {
		if o.Surface.A != 255 {
			panic("mygoui.ScrollFade: SurfaceFade requires an opaque surface")
		}
		ink = o.Surface
	} else if ink.A == 0 || o.Mode == MaskFade {
		ink = theme.Colors(c).Foreground.Alpha(0.18)
	}
	rtl := c.Direction() == ui.RightToLeft
	if o.Mode == MaskFade {
		mask := ui.Local(viewport, scrollMaskKey{}, newScrollMask)
		mask.state, mask.depth, mask.edges, mask.rtl = state, depth, edges, rtl
		// Where the mask does not run, the shadow stands in for it.
		scrollEdgeLayer(c, viewport, func(p *ui.Painter, r ui.Rect) {
			if !mask.runs(p, r) {
				paintScrollEdges(p, r, state, depth, edges, rtl, ink)
			}
		})
		return viewport.Effect(mask)
	}
	scrollEdgeLayer(c, viewport.Clip(), func(p *ui.Painter, r ui.Rect) {
		paintScrollEdges(p, r, state, depth, edges, rtl, ink)
	})
	return viewport
}

// scrollEdgeLayer draws over the viewport's content. Absolute children stay
// fixed in a scroll viewport and paint inside its content clip, before its
// border and native scrollbars.
func scrollEdgeLayer(c *ui.Context, viewport ui.Element, draw func(*ui.Painter, ui.Rect)) {
	viewport.Children(func() {
		ui.Box(c.Key(scrollFadeLayerKey{})).Absolute().Left(0).Top(0).Right(0).Bottom(0).
			PassThrough().Role(ui.RoleNone).Draw(draw)
	})
}

// paintScrollEdges overlays an ink gradient on each selected edge that still
// has content scrolled beyond it.
func paintScrollEdges(p *ui.Painter, r ui.Rect, state *ui.ScrollState, depth float32, edges ScrollFadeEdges, rtl bool, ink ui.Color) {
	if r.W <= 0 || r.H <= 0 {
		return
	}
	x, y := scrollOffsets(state, rtl)
	w, h := min(depth, r.W/2), min(depth, r.H/2)
	strength := func(edge ScrollFadeEdges, remaining, length float32) float32 {
		if edges&edge == 0 || remaining <= 0 {
			return 0
		}
		return min(1, remaining/length)
	}
	fx.PaintEdges(p, r, [4]float32{fx.Top: h, fx.Right: w, fx.Bottom: h, fx.Left: w}, [4]float32{
		fx.Top: strength(FadeTop, y, h), fx.Right: strength(FadeRight, state.MaxX-x, w),
		fx.Bottom: strength(FadeBottom, state.MaxY-y, h), fx.Left: strength(FadeLeft, x, w),
	}, ink)
}

// scrollOffsets returns the clamped physical offsets from the left and top.
func scrollOffsets(state *ui.ScrollState, rtl bool) (x, y float32) {
	x, y = max(0, min(state.X, state.MaxX)), max(0, min(state.Y, state.MaxY))
	if rtl {
		x = state.MaxX - x
	}
	return x, y
}

type scrollFadeLayerKey struct{}

type scrollMaskKey struct{}

type scrollMask struct {
	state  *ui.ScrollState
	shader fx.Mask
	depth  float32
	edges  ScrollFadeEdges
	rtl    bool
	// forced makes a headless frame run the shader, for tests.
	forced bool
}

func newScrollMask() scrollMask { return scrollMask{} }

// gpu reports whether the frame runs shaders: without a GPU MyGo would
// interpret the mask per pixel.
func (m *scrollMask) gpu(p *ui.Painter) bool { return m.forced || p.Accelerated() }

// runs reports whether the mask shader draws the viewport box in this
// frame. Above the layer limit MyGo paints the viewport without its effect.
func (m *scrollMask) runs(p *ui.Painter, box ui.Rect) bool {
	return m.gpu(p) && fx.Fits(box, p.Scale())
}

func (m *scrollMask) EffectMargin() ui.Point { return ui.Point{} }

// PaintEffect runs the shader only where it can change something: the
// content otherwise holds ScrollFade's shadow layer and is drawn as it is.
func (m *scrollMask) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	x, y := scrollOffsets(m.state, m.rtl)
	length := func(edge ScrollFadeEdges, remaining float32) float32 {
		if m.edges&edge == 0 {
			return 0
		}
		return max(0, min(m.depth, remaining))
	}
	top, right := length(FadeTop, y), length(FadeRight, m.state.MaxX-x)
	bottom, left := length(FadeBottom, m.state.MaxY-y), length(FadeLeft, x)
	faded := top > 0 || right > 0 || bottom > 0 || left > 0
	if faded && m.gpu(p) {
		m.paintMask(p, box, content, inset, top, right, bottom, left)
		return
	}
	fx.Plain(p, box, content, inset)
}

func (m *scrollMask) paintMask(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point, top, right, bottom, left float32) {
	m.shader.Paint(p, box, content, inset, [4]float32{fx.Top: top, fx.Right: right, fx.Bottom: bottom, fx.Left: left})
}

// ScrollAxis selects the axes of a native scroll area.
type ScrollAxis uint8

const (
	ScrollVertical ScrollAxis = iota
	ScrollHorizontal
	ScrollBoth
)

// ScrollAreaOptions configures a scroll viewport with native scrollbars.
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
