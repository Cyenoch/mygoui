package mygoui

import (
	"math"
	"sync"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// SkeletonAnimation selects placeholder motion; shimmer is the COSS default.
type SkeletonAnimation uint8

const (
	SkeletonShimmer SkeletonAnimation = iota
	SkeletonPulse
	SkeletonNone
)

// SkeletonOptions reserves the shape of loading content.
type SkeletonOptions struct {
	Width, Height float32
	Radius        float32
	Circle        bool
	Animation     SkeletonAnimation
}

// Skeleton paints a muted placeholder with a paint-only shimmer or pulse.
// Reduced motion and high contrast keep it still.
func Skeleton(c *ui.Context, o SkeletonOptions) ui.Element {
	t := theme.Colors(c)
	height := o.Height
	if height <= 0 {
		height = theme.Scale(c, theme.IconMD)
	}
	radius := o.Radius
	if radius <= 0 {
		radius = theme.RadiusScale().SM
	}
	if o.Circle {
		radius = height / 2
	}
	e := ui.Box(c).Height(height).Shrink(0).Radius(radius).Clip().Role(ui.RoleNone).
		Background(t.MutedBackground)
	if o.Width > 0 {
		e.Width(o.Width)
	} else {
		e.FillWidth()
	}
	if o.Animation == SkeletonNone || theme.ReduceMotion(c) || c.Preferences().HighContrast {
		return e
	}
	windowWidth, windowHeight := c.Size()
	peak := skeletonHighlight(c)
	stops := []ui.GradientStop{
		{Color: ui.Transparent, At: 0.4},
		{Color: ui.RGBA(255, 255, 255, peak), At: 0.5},
		{Color: ui.Transparent, At: 0.6},
	}
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		phase := skeletonPhase(p.Now())
		switch o.Animation {
		case SkeletonPulse:
			k := 1 - float32(math.Abs(float64(2*phase-1)))
			p.Fill(r, t.Surface.Alpha(0.5*ui.EaseInOut(k)), radius)
		default:
			box := skeletonGradientBox(windowWidth, windowHeight, phase)
			p.Clip(r, radius, func() {
				p.FillGradient(box, ui.LinearGradient{Angle: 120, Stops: stops}, 0)
			})
		}
		p.AnimationFrame()
	})
	return e
}

func skeletonPhase(now time.Time) float32 {
	const period = 2 * time.Second
	return float32(now.UnixNano()%int64(period)) / float32(period)
}

func skeletonHighlight(c *ui.Context) float32 {
	if c.Theme().Dark {
		return 0.04
	}
	return 0.64
}

func skeletonGradientBox(width, height, phase float32) ui.Rect {
	return ui.Rect{X: (-1.5 + 2*phase) * width, W: 2 * width, H: height}
}

// SkeletonGroupOptions coordinates placeholder animation.
type SkeletonGroupOptions struct {
	Animation SkeletonAnimation
}

// SkeletonGroup lays out placeholders and animates the entire rendered subtree
// together. Skeletons inside a group should use SkeletonNone. While animated,
// MyGo excludes embedded web views and skips the effect if either layer side
// exceeds 8192 pixels.
func SkeletonGroup(c *ui.Context, o SkeletonGroupOptions, body func()) ui.Element {
	e := ui.Column(c).Gap(theme.Space8).FillWidth().Clip()
	e.Children(body)
	if o.Animation == SkeletonNone || theme.ReduceMotion(c) || c.Preferences().HighContrast {
		return e
	}
	fx := ui.Local(e, skeletonGroupKey{}, newSkeletonGroupEffect)
	fx.width, fx.height = c.Size()
	fx.peak, fx.pulse = skeletonHighlight(c), o.Animation == SkeletonPulse
	return e.Effect(fx)
}

type skeletonGroupKey struct{}

var skeletonGroupShader = sync.OnceValue(func() *shader.Shader {
	return shader.MustCompile(`package main

var Window vec4
var Motion vec3

func ColorEffect(pos vec2, color vec4) vec4 {
	if Motion.z > 0.5 {
		return color * (1.0 - 0.5 * (1.0 - abs(2.0 * Motion.x - 1.0)))
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

type skeletonGroupEffect struct {
	values        *shader.Values
	window        shader.Uniform
	motion        shader.Uniform
	effect        shader.Effect
	width, height float32
	peak          float32
	pulse         bool
}

func newSkeletonGroupEffect() skeletonGroupEffect {
	s := skeletonGroupShader()
	values := s.NewValues()
	return skeletonGroupEffect{
		values: values, window: s.Uniform("Window"), motion: s.Uniform("Motion"),
		effect: s.Effect(shader.Options{Values: values}),
	}
}

func (fx *skeletonGroupEffect) EffectMargin() ui.Point { return ui.Point{} }

func (fx *skeletonGroupEffect) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	mode := float32(0)
	if fx.pulse {
		mode = 1
	}
	fx.values.Vec4(fx.window, box.X, box.Y, fx.width, fx.height)
	fx.values.Vec3(fx.motion, skeletonPhase(p.Now()), fx.peak, mode)
	fx.effect.PaintEffect(p, box, content, inset)
	p.AnimationFrame()
}
