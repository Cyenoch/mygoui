package mygoui

import (
	"github.com/Cyenoch/mygoui/effect"
	"github.com/Cyenoch/mygoui/internal/fx"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
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
	e := ui.Box(c).Height(height).Shrink(0).Radius(radius).Clip().Role(ui.RoleNone)
	if o.Width > 0 {
		e.Width(o.Width)
	} else {
		e.FillWidth()
	}
	tokens := theme.Effects()
	period := tokens.ShimmerPeriod
	fill := t.MutedBackground
	if o.Animation == SkeletonNone || period <= 0 || theme.ReduceMotion(c) || c.Preferences().HighContrast {
		return e.Background(fill)
	}
	if o.Animation == SkeletonPulse {
		// The fill loses alpha as a SkeletonGroup's content does.
		return e.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, fill.Alpha(1-fx.Pulse(fx.Phase(p.Now(), period))), radius)
			p.AnimationFrameRate(60)
		})
	}
	e.Background(fill)
	windowWidth, windowHeight := c.Size()
	peak := tokens.ShimmerPeak
	if c.Theme().Dark {
		peak = tokens.ShimmerPeakDark
	}
	stops := []ui.GradientStop{
		{Color: ui.Transparent, At: 0.4},
		{Color: ui.RGBA(255, 255, 255, peak), At: 0.5},
		{Color: ui.Transparent, At: 0.6},
	}
	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		box := skeletonGradientBox(windowWidth, windowHeight, fx.Phase(p.Now(), period))
		p.Clip(r, radius, func() {
			p.FillGradient(box, ui.LinearGradient{Angle: 120, Stops: stops}, 0)
		})
		p.AnimationFrameRate(60)
	})
	return e
}

func skeletonGradientBox(width, height, phase float32) ui.Rect {
	return ui.Rect{X: (-1.5 + 2*phase) * width, W: 2 * width, H: height}
}

// SkeletonGroupOptions coordinates placeholder animation.
type SkeletonGroupOptions struct {
	Animation SkeletonAnimation
}

// SkeletonGroup lays out placeholders, as a busy group named by the locale's
// Loading, and animates the entire rendered subtree
// together through effect.Shimmer. Skeletons inside a group should use
// SkeletonNone. Without a GPU, above MyGo's 8192-pixel effect layer, under
// reduced motion and under high contrast the group rests; embedded web
// views do not show through it while it animates.
func SkeletonGroup(c *ui.Context, o SkeletonGroupOptions, body func()) ui.Element {
	e := ui.Column(c).Gap(theme.Space(2)).FillWidth().Clip().
		Role(ui.RoleGroup).Label(locale.M().Loading).Busy(true)
	e.Children(body)
	if o.Animation == SkeletonNone {
		return e
	}
	return effect.Shimmer(c, e, effect.ShimmerOptions{Pulse: o.Animation == SkeletonPulse})
}
