// Package busy holds the loading treatment components share, so a pending
// button, a working field and a refreshing list fade, wait and dim alike.
package busy

import (
	"math"
	"time"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// LineHeight is the height of the refresh line, before scaling.
const LineHeight = 2

// frameRate caps the indicators' redraws per second.
const frameRate = 60

// fade is the crossfade between content and its indicator: a state change
// in place.
func fade() ui.Animation {
	m := theme.Motion()
	return ui.Animation{Duration: m.State, Ease: m.Ease, NonMotion: true}
}

type (
	swapKey  struct{}
	showKey  struct{}
	startKey struct{}
	builtKey struct{}
)

// Swap returns how far e has gone from its content (0) to its indicator
// (1). A press is answered at once, so there is no wait.
func Swap(e ui.Element, on bool) float32 {
	return e.AnimateTo(swapKey{}, target(on), fade())
}

// Show returns the opacity of a region's indicator. It waits before fading
// in and leaves at once, so content that arrives is not held back. A region
// built already loading waits too: an animation starts at its first target,
// so the first build names the hidden one.
func Show(e ui.Element, on bool) float32 {
	a := fade()
	if on {
		a.Delay = theme.Effects().LoadingDelay
	}
	if built := ui.Local(e, builtKey{}, func() bool { return false }); !*built {
		*built = true
		e.AnimateTo(showKey{}, 0, a)
	}
	return e.AnimateTo(showKey{}, target(on), a)
}

func target(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// Dim fades content under a refresh by shown, the value Show returned, and
// takes its input away while on.
func Dim(content ui.Element, on bool, shown float32) ui.Element {
	return content.Opacity(Dimmed(shown)).IgnoreInput(on)
}

// Dimmed is the opacity content takes under a refresh shown this far.
func Dimmed(shown float32) float32 { return 1 - (1-theme.Effects().RefreshOpacity)*shown }

// Spinner is the arc indicator at size DIPs, before scaling, named Loading.
func Spinner(c *ui.Context, size float32, ink ui.Color) ui.Element {
	return loading.Arc(c, loading.Options{
		Size: theme.Scale(c, size), Color: ink, Duration: theme.Effects().SpinnerPeriod,
	}).Label(locale.M().Loading)
}

// PaintSpinner paints the arc indicator centered on (x, y) for a control
// that cannot hold a child, such as a native editor. It turns once a spinner
// period and rests still under reduced motion.
func PaintSpinner(p *ui.Painter, x, y, size float32, ink ui.Color, still bool) {
	const steps = 12
	turn := 0.0
	if !still {
		period := theme.Effects().SpinnerPeriod.Milliseconds()
		turn = float64(p.Now().UnixMilli()%period) / float64(period) * 2 * math.Pi
		p.AnimationFrameRate(frameRate)
	}
	width := size * 2.5 / 24
	radius := float64(size/2 - width/2)
	var points [steps + 1][2]float32
	for i := range points {
		a := turn + float64(i)/steps*1.8
		points[i] = [2]float32{x + float32(radius*math.Cos(a)), y + float32(radius*math.Sin(a))}
	}
	p.StrokePolyline(points[:], width, ink)
}

// Sweep returns where the band of an indeterminate line from x, w wide,
// starts after elapsed, and how wide it is: the theme's SweepBand of the
// width, entering from wholly before the line and leaving wholly past it
// over each SweepPeriod, easing in and out. A narrow band travels further,
// so it leaves the far edge before the next pass. Reduced motion holds it at
// the start.
func Sweep(x, w float32, elapsed time.Duration, reduced bool) (from, band float32) {
	fx := theme.Effects()
	band = w * fx.SweepBand
	if reduced {
		return x, band
	}
	phase := float32(elapsed%fx.SweepPeriod) / float32(fx.SweepPeriod)
	return x - band + ui.EaseInOut(phase)*max(band*4.5, w+band), band
}

// Line lays an indeterminate line over the top edge of the region being
// built, inset by its corner radius so it needs no clip. Call it among the
// region's children while shown, the value Show returned, is above zero.
func Line(c *ui.Context, shown, inset float32) {
	if shown <= 0 {
		return
	}
	ink := theme.Colors(c).Primary
	reduced := theme.ReduceMotion(c)
	height := theme.Scale(c, LineHeight)
	line := ui.Box(c).Absolute().Top(0).Left(inset).Right(inset).Height(height).
		PassThrough().Role(ui.RoleStatus).Label(locale.M().Loading).Opacity(shown)
	start := ui.Local(line, startKey{}, func() time.Time { return c.Now() })
	line.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 {
			return
		}
		x, band := Sweep(r.X, r.W, p.Now().Sub(*start), reduced)
		left, right := max(x, r.X), min(x+band, r.X+r.W)
		if right > left {
			p.Fill(ui.Rect{X: left, Y: r.Y, W: right - left, H: r.H}, ink, r.H/2)
		}
		if !reduced {
			p.AnimationFrameRate(frameRate)
		}
	})
}
