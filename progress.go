package mygoui

import (
	"math"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type ProgressOptions struct {
	Value, Min, Max float64
	Indeterminate   bool
	Label           string
}

// Progress normalizes actual bounds. Zero bounds mean 0..100; an invalid
// range or non-finite value is unknown, never a fabricated completion value.
func Progress(c *ui.Context, o ProgressOptions) ui.Element {
	value := progressValue(o)
	t := theme.Colors(c)
	bar := ui.Progress(c, value).Label(o.Label).Height(theme.Space8).
		Radius(theme.Radii().XXL).Background(t.Secondary)
	// Retain the native progress role and unknown range while replacing only
	// paint: pinned MyGo's indeterminate painter ignores ReduceMotion.
	reduced := c.Preferences().ReduceMotion
	rtl := c.Direction() == ui.RightToLeft
	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		w, x := r.W*float32(max(0, value)), r.X
		if value < 0 {
			w = r.W * 0.3
			x = r.X + (r.W-w)/2
			if !reduced {
				p.AnimationFrame()
				phase := float32(p.Now().UnixMilli()%1400) / 1400
				x = r.X - w + (r.W+w)*phase
			}
		}
		if rtl {
			x = r.X + r.W - (x - r.X) - w
		}
		p.Clip(r, theme.Radii().XXL, func() {
			p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, t.Primary, theme.Radii().XXL)
		})
	})
	return bar
}

func progressValue(o ProgressOptions) float64 {
	if o.Indeterminate || math.IsNaN(o.Value) || math.IsInf(o.Value, 0) {
		return -1
	}
	if o.Min == 0 && o.Max == 0 {
		o.Max = 100
	}
	if math.IsNaN(o.Min) || math.IsNaN(o.Max) || math.IsInf(o.Min, 0) || math.IsInf(o.Max, 0) || o.Max <= o.Min {
		return -1
	}
	return min(1, max(0, (o.Value-o.Min)/(o.Max-o.Min)))
}
