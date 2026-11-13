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

const (
	// fade is the crossfade between content and its indicator.
	fade = 150 * time.Millisecond
	// wait keeps a region's indicator back, so a load that ends sooner
	// shows none.
	wait = 150 * time.Millisecond
	// dim is the opacity content keeps under a refresh.
	dim = 0.55
	// LineHeight is the height of the refresh line, before scaling.
	LineHeight = 2
)

type (
	swapKey  struct{}
	showKey  struct{}
	startKey struct{}
	builtKey struct{}
)

// Swap returns how far e has gone from its content (0) to its indicator
// (1). A press is answered at once, so there is no wait.
func Swap(e ui.Element, on bool) float32 {
	return e.AnimateTo(swapKey{}, target(on), ui.Animation{Duration: fade, NonMotion: true})
}

// Show returns the opacity of a region's indicator. It waits before fading
// in and leaves at once, so content that arrives is not held back. A region
// built already loading waits too: an animation starts at its first target,
// so the first build names the hidden one.
func Show(e ui.Element, on bool) float32 {
	a := ui.Animation{Duration: fade, NonMotion: true}
	if on {
		a.Delay = wait
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
func Dimmed(shown float32) float32 { return 1 - (1-dim)*shown }

// Spinner is the arc indicator at size DIPs, before scaling, named Loading.
func Spinner(c *ui.Context, size float32, ink ui.Color) ui.Element {
	return loading.Arc(c, loading.Options{
		Size: theme.Scale(c, size), Color: ink, Duration: time.Second,
	}).Label(locale.M().Loading)
}

// PaintSpinner paints the arc indicator centered on (x, y) for a control
// that cannot hold a child, such as a native editor. It turns once a second
// and rests still under reduced motion.
func PaintSpinner(p *ui.Painter, x, y, size float32, ink ui.Color, still bool) {
	const steps = 12
	turn := 0.0
	if !still {
		turn = float64(p.Now().UnixMilli()%1000) / 1000 * 2 * math.Pi
		p.AnimationFrameRate(60)
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
		band := r.W * 0.4
		if reduced {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: band, H: r.H}, ink, r.H/2)
			return
		}
		const period = 1500 * time.Millisecond
		phase := float32(p.Now().Sub(*start)%period) / float32(period)
		// The band starts half inside the edge, so the line shows ink as
		// soon as it fades in.
		x := r.X - band/2 + ui.EaseInOut(phase)*(band*4)
		left, right := max(x, r.X), min(x+band, r.X+r.W)
		if right > left {
			p.Fill(ui.Rect{X: left, Y: r.Y, W: right - left, H: r.H}, ink, r.H/2)
		}
		p.AnimationFrameRate(60)
	})
}
