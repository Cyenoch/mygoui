package mygoui

import (
	"math"
	"time"

	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusDotOptions configures a status indicator. The default status is a
// muted dot, as for something offline. Pulse radiates a ring from the dot,
// for something live; Breathe slowly dims and restores the dot, for
// something in progress such as connecting. Both are still under reduced
// motion.
type StatusDotOptions struct {
	Status  StatusColor
	Label   string
	Pulse   bool
	Breathe bool
}

type statusDotInkKey struct{}

const (
	// statusDotPeriod is one ring of the pulse.
	statusDotPeriod = 2 * time.Second
	// statusDotBreath is one dimming and restoring of the dot.
	statusDotBreath = 1600 * time.Millisecond
)

// StatusDot shows a small colored dot before an optional label. The dot is
// decoration: without a visible Label, name the returned element with Label
// so the status does not rest on color alone. The dot takes its own size
// only: a pulse's ring grows past it, over what is around. The motion only
// repaints the dot; a still dot asks for no frames.
func StatusDot(c *ui.Context, o StatusDotOptions) ui.Element {
	t := theme.Colors(c)
	ink := t.MutedForeground
	if o.Status != ColorDefault {
		ink = statusForeground(t, o.Status)
	}
	still := theme.ReduceMotion(c)
	pulse, breathe := o.Pulse && !still, o.Breathe && !still
	root := ui.Row(c).MinWidth(0).Gap(theme.Space8).AlignItems(ui.Center)
	return root.Children(func() {
		side := theme.Scale(c, theme.Space8)
		dot := ui.Box(c).Size(side, side).Shrink(0).Role(ui.RoleNone)
		ink := motion.Color(dot, statusDotInkKey{}, ink)
		dot.Draw(func(p *ui.Painter, r ui.Rect) {
			body := ink
			switch {
			case pulse:
				phase := float32(p.Now().UnixNano()%int64(statusDotPeriod)) / float32(statusDotPeriod)
				grow := ui.EaseOut(phase)
				p.Fill(centeredSquare(r, r.W*(1+grow)), ink.Alpha(0.4*(1-grow)), r.W*2)
				p.AnimationFrameRate(30)
			case breathe:
				phase := float64(p.Now().UnixNano()%int64(statusDotBreath)) / float64(statusDotBreath)
				body = ink.Alpha(0.65 + 0.35*float32(math.Cos(2*math.Pi*phase)+1)/2)
				p.AnimationFrameRate(30)
			}
			p.Fill(r, body, r.W)
		})
		if o.Label != "" {
			style.Text(ui.Text(c, o.Label), c, theme.SM()).TextColor(t.Foreground).MinWidth(0)
		}
	})
}

func centeredSquare(r ui.Rect, side float32) ui.Rect {
	return ui.Rect{X: r.X + (r.W-side)/2, Y: r.Y + (r.H-side)/2, W: side, H: side}
}
