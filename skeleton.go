package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type SkeletonOptions struct {
	Width, Height float32
	Round         bool
}

// Skeleton reserves content geometry with a muted, decorative placeholder.
// Its clipped shimmer uses the native paint clock without rebuilding the view;
// reduced motion and high contrast keep the placeholder still.
// Source shape/color: https://ui.shadcn.com/r/styles/base-rhea/skeleton.json
func Skeleton(c *myui.Context, o SkeletonOptions) myui.Element {
	height := o.Height
	if height <= 0 {
		height = theme.Scale(c, theme.IconMedium)
	}
	e := myui.Box(c).Height(height).Shrink(0).Radius(theme.Radii().XXL).
		Background(theme.Colors(c).Muted).Role(myui.RoleNone)
	if o.Width > 0 {
		e.Width(o.Width)
	} else {
		e.FillWidth()
	}
	if o.Round {
		e.Radius(height)
	}
	if !theme.ReduceMotion(c) && !c.Preferences().HighContrast {
		started := myui.Local(e, "skeleton-start", func() time.Time { return c.Now() })
		colors := theme.Colors(c)
		shine := colors.Background.Alpha(0.85)
		if c.Theme().Dark {
			shine = colors.Foreground.Alpha(0.12)
		}
		radius := theme.Radii().XXL
		if o.Round {
			radius = height
		}
		rtl := c.Direction() == myui.RightToLeft
		e.DrawOver(func(p *myui.Painter, r myui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			const period = 1600 * time.Millisecond
			phase := float32(p.Now().Sub(*started)%period) / float32(period)
			if rtl {
				phase = 1 - phase
			}
			band := min(float32(120), r.W*0.6)
			x := r.X - band + phase*(r.W+2*band)
			p.Clip(r, radius, func() {
				p.FillGradient(myui.Rect{X: x, Y: r.Y, W: band / 2, H: r.H}, myui.LinearGradient{
					From: shine.Alpha(0), To: shine, Angle: 90,
				}, 0)
				p.FillGradient(myui.Rect{X: x + band/2, Y: r.Y, W: band / 2, H: r.H}, myui.LinearGradient{
					From: shine, To: shine.Alpha(0), Angle: 90,
				}, 0)
			})
			p.AnimationFrame()
		})
	}
	return e
}
