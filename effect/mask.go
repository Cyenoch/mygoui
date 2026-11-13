package effect

import (
	"github.com/Cyenoch/mygoui/internal/fx"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Edges selects physical edges of an element; zero selects every edge.
type Edges uint8

const (
	EdgeTop Edges = 1 << iota
	EdgeRight
	EdgeBottom
	EdgeLeft
)

// maskShadow is the share of the foreground the fallback shades an edge
// with, as ScrollFade's shadow does.
const maskShadow = 0.18

// MaskOptions configures Mask. The zero value fades every edge over the
// theme's fade depth.
type MaskOptions struct {
	// Edges are the edges that fade; zero fades all four.
	Edges Edges
	// Size is the fade depth in DIPs; zero or less takes
	// theme.Effects().FadeEdge. A fade reaches at most half of the element.
	Size float32
	// Surface is the opaque color behind the element. The fallback blends
	// the edges into it; the zero color shades them with a shadow instead.
	Surface ui.Color
	// Radius is the element's corner radius, which the fallback follows.
	Radius float32
	// Fallback draws the fallback where the shader could run.
	Fallback bool
}

// Mask fades the edges of e and its children to transparency, so that what
// lies behind e shows through: use it where content continues past an edge
// and is cut there. Apply it after e's children and styling; for a scroll
// viewport use mygoui.ScrollFade, which follows the scroll position.
//
// Without a GPU, above the layer limit or with Fallback, the content is
// drawn whole and each edge is shaded instead: with a shadow, or blended
// into Surface. Either way e clips its content, as the shader's layer does.
// Mask takes e's effect and its DrawOver callback: MyGo keeps one DrawOver
// per element and gives no way to read it back to chain it. It panics for a
// translucent Surface.
func Mask(c *ui.Context, e ui.Element, o MaskOptions) ui.Element {
	if o.Surface.A != 0 && o.Surface.A != 255 {
		panic("effect.Mask: Surface must be opaque, or zero for a shadow")
	}
	size := o.Size
	if size <= 0 {
		size = theme.Effects().FadeEdge
	}
	edges := o.Edges
	if edges == 0 {
		edges = EdgeTop | EdgeRight | EdgeBottom | EdgeLeft
	}
	ink := o.Surface
	if ink.A == 0 {
		ink = theme.Colors(c).Foreground.Alpha(maskShadow)
	}
	m := ui.Local(e, maskKey{}, newMask)
	m.fallback = o.Fallback
	for i, edge := range [4]Edges{fx.Top: EdgeTop, fx.Right: EdgeRight, fx.Bottom: EdgeBottom, fx.Left: EdgeLeft} {
		m.depth[i] = 0
		if edges&edge != 0 {
			m.depth[i] = size
		}
	}
	radius := o.Radius
	e.Clip()
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 || !m.fallback && shaded(p, r, m.forced) {
			return
		}
		p.Clip(r, radius, func() {
			fx.PaintEdges(p, r, m.capped(r), [4]float32{1, 1, 1, 1}, ink)
		})
	})
	if o.Fallback {
		return e
	}
	return e.Effect(m)
}

type maskKey struct{}

type mask struct {
	shader   fx.Mask
	depth    [4]float32
	fallback bool
	// forced makes a headless frame run the shader, for tests.
	forced bool
}

func newMask() mask { return mask{} }

// capped returns the depths limited to half of r, so that opposite fades
// meet in the middle of a small element instead of dimming all of it.
func (m *mask) capped(r ui.Rect) [4]float32 {
	depth := m.depth
	for i := range depth {
		half := r.H / 2
		if i == fx.Left || i == fx.Right {
			half = r.W / 2
		}
		depth[i] = min(depth[i], half)
	}
	return depth
}

func (m *mask) EffectMargin() ui.Point { return ui.Point{} }

func (m *mask) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	if !gpu(p, m.forced) {
		// The content holds the shading DrawOver painted.
		fx.Plain(p, box, content, inset)
		return
	}
	m.shader.Paint(p, box, content, inset, m.capped(box))
}
