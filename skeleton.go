package mygoui

import (
	"math"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SkeletonAnimation selects HeroUI's placeholder motion.
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
		Background(t.SurfaceTertiary.Alpha(0.7))
	if o.Width > 0 {
		e.Width(o.Width)
	} else {
		e.FillWidth()
	}
	if o.Animation == SkeletonNone || theme.ReduceMotion(c) || c.Preferences().HighContrast {
		return e
	}
	start := ui.Local(e, skeletonStartKey{}, func() time.Time { return c.Now() })
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		now := p.Now()
		switch o.Animation {
		case SkeletonPulse:
			phase := float32(now.Sub(*start)%(1500*time.Millisecond)) / float32(1500*time.Millisecond)
			alpha := 0.35 * (0.5 - 0.5*float32(math.Cos(float64(phase)*2*math.Pi)))
			p.Fill(r, t.SurfaceTertiary.Alpha(alpha), radius)
			p.AnimationFrame()
		default:
			phase := float32(now.Sub(*start)%(2000*time.Millisecond)) / float32(2000*time.Millisecond)
			band := min(float32(160), r.W*0.6)
			x := r.X - band + phase*(r.W+2*band)
			p.FillGradient(ui.Rect{X: x, Y: r.Y, W: band / 2, H: r.H}, ui.LinearGradient{
				From: ui.Transparent, To: t.SurfaceTertiary.Alpha(0.85), Angle: 90,
			}, 0)
			p.FillGradient(ui.Rect{X: x + band/2, Y: r.Y, W: band / 2, H: r.H}, ui.LinearGradient{
				From: t.SurfaceTertiary.Alpha(0.85), To: ui.Transparent, Angle: 90,
			}, 0)
			p.AnimationFrame()
		}
	})
	return e
}

type skeletonStartKey struct{}

// SkeletonGroupOptions coordinates placeholder animation.
type SkeletonGroupOptions struct {
	Animation SkeletonAnimation
}

// SkeletonGroup lays out placeholders and paints one synchronized shimmer or
// pulse across the whole group. Skeletons inside a group should use
// SkeletonNone so the group owns the animation; their shapes remain.
func SkeletonGroup(c *ui.Context, o SkeletonGroupOptions, body func()) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).Gap(theme.Space8).FillWidth().Clip()
	e.Children(body)
	if o.Animation == SkeletonNone || theme.ReduceMotion(c) || c.Preferences().HighContrast {
		return e
	}
	start := ui.Local(e, skeletonGroupStartKey{}, func() time.Time { return c.Now() })
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if o.Animation == SkeletonPulse {
			phase := float32(p.Now().Sub(*start)%(1500*time.Millisecond)) / float32(1500*time.Millisecond)
			alpha := 0.22 * (0.5 - 0.5*float32(math.Cos(float64(phase)*2*math.Pi)))
			p.Fill(r, t.SurfaceTertiary.Alpha(alpha), 0)
			p.AnimationFrame()
			return
		}
		phase := float32(p.Now().Sub(*start)%(2000*time.Millisecond)) / float32(2000*time.Millisecond)
		band := min(float32(220), r.W*0.7)
		x := r.X - band + phase*(r.W+2*band)
		p.FillGradient(ui.Rect{X: x, Y: r.Y, W: band / 2, H: r.H}, ui.LinearGradient{
			From: ui.Transparent, To: t.SurfaceTertiary.Alpha(0.5), Angle: 90,
		}, 0)
		p.FillGradient(ui.Rect{X: x + band/2, Y: r.Y, W: band / 2, H: r.H}, ui.LinearGradient{
			From: t.SurfaceTertiary.Alpha(0.5), To: ui.Transparent, Angle: 90,
		}, 0)
		p.AnimationFrame()
	})
	return e
}

type skeletonGroupStartKey struct{}
