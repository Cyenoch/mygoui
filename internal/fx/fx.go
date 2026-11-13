// Package fx holds the shaders and paint helpers the effect family and the
// root components draw with. A shader is compiled the first time it paints
// and shared by every element.
package fx

import (
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// maxLayer is the largest side, in pixels, of the bitmap MyGo renders an
// effect's content into (ui/effect.go, maxEffectLayer). A larger element is
// painted as it is and its effect is not called.
const maxLayer = 8192

// Fits reports whether MyGo renders box, in DIPs at scale, into an effect
// layer. The effects of this package have no margin.
func Fits(box ui.Rect, scale float32) bool {
	side := func(from, size float32) float64 {
		return math.Ceil(float64((from+size)*scale)) - math.Floor(float64(from*scale))
	}
	return side(box.X, box.W) <= maxLayer && side(box.Y, box.H) <= maxLayer
}

// Plain draws an effect's content unchanged.
func Plain(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	p.ImageRegion(content, source(p, box, inset), box, ui.ImageOptions{})
}

// Faded draws an effect's content at alpha.
func Faded(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point, alpha float32) {
	p.ImageRegion(content, source(p, box, inset), box, ui.ImageOptions{Tint: ui.RGBA(255, 255, 255, alpha)})
}

func source(p *ui.Painter, box ui.Rect, inset ui.Point) ui.Rect {
	s := p.Scale()
	return ui.Rect{X: inset.X, Y: inset.Y, W: box.W * s, H: box.H * s}
}

// Phase returns how far now is through a loop of period, from 0 to 1.
func Phase(now time.Time, period time.Duration) float32 {
	return float32(now.UnixNano()%int64(period)) / float32(period)
}

// Pulse returns how much of its alpha a pulsing placeholder has lost at
// phase of the loop: none at the ends, half at the middle, easing in and
// out.
func Pulse(phase float32) float32 {
	k := 1 - float32(math.Abs(float64(2*phase-1)))
	return 0.5 * ui.EaseInOut(k)
}

// Edge indexes the physical edges in the order the mask and PaintEdges take
// them.
const (
	Top = iota
	Right
	Bottom
	Left
)

// PaintEdges overlays ink on the edges of r: edge i fades inwards over
// depth[i] DIPs from strength[i] of the ink. A non-positive depth or
// strength leaves the edge alone.
func PaintEdges(p *ui.Painter, r ui.Rect, depth, strength [4]float32, ink ui.Color) {
	paint := func(edge int, box ui.Rect, angle float32, reverse bool) {
		if depth[edge] <= 0 || strength[edge] <= 0 {
			return
		}
		from, to := ink.Alpha(strength[edge]), ink.Alpha(0)
		if reverse {
			from, to = to, from
		}
		p.FillGradient(box, ui.LinearGradient{From: from, To: to, Angle: angle}, 0)
	}
	paint(Top, ui.Rect{X: r.X, Y: r.Y, W: r.W, H: depth[Top]}, 180, false)
	paint(Bottom, ui.Rect{X: r.X, Y: r.Y + r.H - depth[Bottom], W: r.W, H: depth[Bottom]}, 180, true)
	paint(Left, ui.Rect{X: r.X, Y: r.Y, W: depth[Left], H: r.H}, 90, false)
	paint(Right, ui.Rect{X: r.X + r.W - depth[Right], Y: r.Y, W: depth[Right], H: r.H}, 90, true)
}

var maskShader = sync.OnceValue(func() *shader.Shader {
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

// Mask multiplies the content's alpha by a ramp at each edge. It belongs to
// one element: keep it with ui.Local. The zero value is ready; the shader is
// compiled when it first paints.
type Mask struct {
	values *shader.Values
	edges  shader.Uniform
	effect shader.Effect
}

// Paint draws content fading to transparency over depth DIPs at each edge;
// a zero depth leaves that edge alone.
func (m *Mask) Paint(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point, depth [4]float32) {
	if m.values == nil {
		s := maskShader()
		m.values, m.edges = s.NewValues(), s.Uniform("Edges")
		m.effect = s.Effect(shader.Options{Values: m.values})
	}
	m.values.Vec4(m.edges, depth[Top], depth[Right], depth[Bottom], depth[Left])
	m.effect.PaintEffect(p, box, content, inset)
}

var shimmerShader = sync.OnceValue(func() *shader.Shader {
	return shader.MustCompile(`package main

var Window vec4
var Motion vec3

func ColorEffect(pos vec2, color vec4) vec4 {
	if Motion.z > 0.5 {
		return color * (1.0 - Motion.x)
	}
	center := (-0.5 + 2.0 * Motion.x) * Window.z
	x := pos.x + Window.x - center
	y := pos.y + Window.y - Window.w / 2.0
	half := (2.0 * Window.z * 0.8660254 + Window.w * 0.5) / 2.0
	t := 0.5 + (0.8660254 * x + 0.5 * y) / (2.0 * half)
	highlight := clamp(1.0 - abs(t - 0.5) / 0.1, 0.0, 1.0) * Motion.y
	return vec4(mix(color.rgb, vec3(color.a), highlight), color.a)
}
`)
})

// Shimmer passes a highlight band across the content, placed in window
// coordinates so that every region in a window moves as one, or pulses the
// content's alpha. It belongs to one element: keep it with ui.Local. The
// zero value is ready; the shader is compiled when it first paints.
type Shimmer struct {
	values *shader.Values
	window shader.Uniform
	motion shader.Uniform
	effect shader.Effect
}

// Paint draws content at phase of the loop in a window of width by height
// DIPs: a band mixing towards white by up to peak, or with pulse the alpha
// falling as Pulse does.
func (s *Shimmer) Paint(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point, width, height, phase, peak float32, pulse bool) {
	if s.values == nil {
		compiled := shimmerShader()
		s.values, s.window, s.motion = compiled.NewValues(), compiled.Uniform("Window"), compiled.Uniform("Motion")
		s.effect = compiled.Effect(shader.Options{Values: s.values})
	}
	mode := float32(0)
	if pulse {
		mode, phase = 1, Pulse(phase)
	}
	s.values.Vec4(s.window, box.X, box.Y, width, height)
	s.values.Vec3(s.motion, phase, peak, mode)
	s.effect.PaintEffect(p, box, content, inset)
}

var wipeShader = sync.OnceValue(func() *shader.Shader {
	return shader.MustCompile(`package main

var Dir vec2
var Front vec2

func ColorEffect(pos vec2, color vec4) vec4 {
	size := rectSize()
	d := dot(pos, Dir)
	if Dir.x < 0.0 {
		d += size.x
	}
	if Dir.y < 0.0 {
		d += size.y
	}
	return color * clamp((Front.x - d) / Front.y, 0.0, 1.0)
}
`)
})

// Wipe shows the content up to a soft front that travels across it. It
// belongs to one element: keep it with ui.Local. The zero value is ready;
// the shader is compiled when it first paints.
type Wipe struct {
	values *shader.Values
	dir    shader.Uniform
	front  shader.Uniform
	effect shader.Effect
}

// Paint draws content revealed to progress, from 0 (nothing) to 1 (all),
// travelling along the unit axis (dx, dy) with a front soft DIPs deep.
func (w *Wipe) Paint(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point, dx, dy, progress, soft float32) {
	if w.values == nil {
		s := wipeShader()
		w.values, w.dir, w.front = s.NewValues(), s.Uniform("Dir"), s.Uniform("Front")
		w.effect = s.Effect(shader.Options{Values: w.values})
	}
	extent := box.H
	if dx != 0 {
		extent = box.W
	}
	w.values.Vec2(w.dir, dx, dy)
	w.values.Vec2(w.front, progress*(extent+soft), soft)
	w.effect.PaintEffect(p, box, content, inset)
}
