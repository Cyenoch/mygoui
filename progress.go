package mygoui

import (
	"math"
	"strconv"
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// IndicatorColor is the color of a progress or meter indicator. Its zero
// value uses the primary color, unlike the status colors used by chips.
type IndicatorColor uint8

const (
	IndicatorAccent IndicatorColor = iota
	IndicatorDefault
	IndicatorSuccess
	IndicatorWarning
	IndicatorDanger
)

// ProgressSize selects a native meter or progress track height.
type ProgressSize uint8

const (
	ProgressMedium ProgressSize = iota
	ProgressSmall
	ProgressLarge
)

// progressDuration is coss ui's 500 ms indicator transition.
const progressDuration = 500 * time.Millisecond

// indicatorTransition animates a value change on the indicator. MyGo lays the
// fill out at each width on the way, so the track never jumps and asks for
// frames only while the value moves.
var indicatorTransition = ui.ElementTransition{
	Duration: progressDuration, Ease: ui.EaseOut, Size: true,
}

// indicatorAnimation is the same cadence for paint that is not a box, as the
// determinate ring.
var indicatorAnimation = ui.Animation{Duration: progressDuration, Ease: ui.EaseOut}

// progressFill lays the bar's fill in its track and eases a value change,
// unless the desktop asks for less motion.
func progressFill(c *ui.Context, pct, height, radius float32, fill ui.Color) {
	bar := ui.Box(c).Height(height).WidthPercent(pct * 100).
		Background(fill).Radius(radius)
	if theme.ReduceMotion(c) {
		return
	}
	bar.Transition(indicatorTransition)
}

// ProgressOptions configures a coss ui bar indicator.
type ProgressOptions struct {
	Label           string
	Value, Min, Max float64
	Indeterminate   bool
	Color           IndicatorColor
	Size            ProgressSize
	ShowOutput      bool
	OutputText      string
}

// ProgressBar draws a linear determinate or indeterminate indicator, with an
// optional label and value output above the track.
func ProgressBar(c *ui.Context, o ProgressOptions) ui.Element {
	t := theme.Colors(c)
	lo, hi := progressBounds(o.Min, o.Max)
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).
		Role(ui.RoleProgress).Label(o.Label)
	if o.Indeterminate {
		e.Range(0, 1, -1)
	} else {
		e.Range(lo, hi, o.Value)
	}
	e.Children(func() {
		output := ""
		if o.ShowOutput {
			output = progressText(o)
		}
		progressHeading(c, o.Label, o.ShowOutput, output)
		height, radius := progressMetrics(o.Size)
		fill := indicatorFill(c, o.Color)
		track := ui.Box(c).FillWidth().Height(height).Radius(radius).Clip().
			Background(t.Input)
		if o.Indeterminate {
			start := ui.Local(track, progressStartKey{}, func() time.Time { return c.Now() })
			reduced := theme.ReduceMotion(c)
			track.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if r.W <= 0 || r.H <= 0 {
					return
				}
				band := r.W * 0.4
				if reduced {
					// Reduced motion holds the 2/5 fill at the start.
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: band, H: r.H}, fill, radius)
					return
				}
				// Sweep from -100% to 350% of the band's own width.
				const period = 1500 * time.Millisecond
				phase := float32(p.Now().Sub(*start)%period) / float32(period)
				x := r.X - band + ui.EaseInOut(phase)*(band*4.5)
				p.Fill(ui.Rect{X: x, Y: r.Y, W: band, H: r.H}, fill, radius)
				p.AnimationFrame()
			})
		} else {
			pct := progressFraction(o.Value, lo, hi)
			track.Children(func() {
				progressFill(c, pct, height, radius, fill)
			})
		}
	})
	return e
}

// MeterOptions configures a quantity within a known range.
type MeterOptions struct {
	Label           string
	Value, Min, Max float64
	Color           IndicatorColor
	Size            ProgressSize
	ShowOutput      bool
	OutputText      string
}

// Meter renders a coss ui meter with its label and output.
func Meter(c *ui.Context, o MeterOptions) ui.Element {
	t := theme.Colors(c)
	lo, hi := progressBounds(o.Min, o.Max)
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).
		Role(ui.RoleMeter).Label(o.Label).Range(lo, hi, o.Value)
	e.Children(func() {
		output := ""
		if o.ShowOutput {
			output = meterText(o)
		}
		progressHeading(c, o.Label, o.ShowOutput, output)
		height, _ := progressMetrics(o.Size)
		if o.Size == ProgressMedium {
			height = theme.Space8
		}
		radius := float32(0)
		pct := progressFraction(o.Value, lo, hi)
		fill := indicatorFill(c, o.Color)
		ui.Box(c).FillWidth().Height(height).Radius(radius).Clip().Background(t.Input).
			Children(func() {
				progressFill(c, pct, height, radius, fill)
			})
	})
	return e
}

func progressHeading(c *ui.Context, label string, showOutput bool, output string) {
	if label == "" && !showOutput {
		return
	}
	ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
		if label != "" {
			ts := theme.SM()
			ts.Weight = theme.WeightMedium
			style.Text(ui.Text(c, label), c, ts).TextColor(theme.Colors(c).Foreground)
		}
		if showOutput {
			ui.Spacer(c)
			ts := theme.SM()
			ts.Weight = theme.WeightMedium
			style.Text(ui.Text(c, output), c, ts).TextColor(theme.Colors(c).Foreground).
				FontFeatures("tnum")
		}
	})
}

func progressMetrics(size ProgressSize) (height, radius float32) {
	height, radius = theme.Space6, 999
	switch size {
	case ProgressSmall:
		height, radius = theme.Space4, 999
	case ProgressLarge:
		height, radius = theme.Space8, 999
	}
	return height, radius
}

type progressStartKey struct{}

func progressText(o ProgressOptions) string {
	if o.OutputText != "" {
		return o.OutputText
	}
	if o.Indeterminate {
		return locale.M().Loading
	}
	return percentText(o.Value, o.Min, o.Max)
}

func meterText(o MeterOptions) string {
	if o.OutputText != "" {
		return o.OutputText
	}
	return percentText(o.Value, o.Min, o.Max)
}

func progressBounds(min, max float64) (float64, float64) {
	if max == 0 && min == 0 {
		return 0, 100
	}
	if max <= min {
		return min, min + 1
	}
	return min, max
}

func progressFraction(value, min, max float64) float32 {
	if max <= min {
		return 0
	}
	return clamp01((value - min) / (max - min))
}

func percentText(value, min, max float64) string {
	lo, hi := progressBounds(min, max)
	return strconv.Itoa(int(progressFraction(value, lo, hi)*100+0.5)) + "%"
}

func indicatorFill(c *ui.Context, color IndicatorColor) ui.Color {
	t := theme.Colors(c)
	switch color {
	case IndicatorDefault:
		return t.DefaultForeground
	case IndicatorSuccess:
		return t.Success
	case IndicatorWarning:
		return t.Warning
	case IndicatorDanger:
		return t.Danger
	}
	return t.Primary
}

func clamp01(v float64) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return float32(v)
}

// ProgressCircleOptions configures a circular indicator.
type ProgressCircleOptions struct {
	Label           string
	Value, Min, Max float64
	Indeterminate   bool
	Color           IndicatorColor
	Size            ProgressSize
	ShowOutput      bool
}

// ProgressCircle draws a determinate ring or an indeterminate native spinner
// and exposes its progress through native accessibility semantics.
func ProgressCircle(c *ui.Context, o ProgressCircleOptions) ui.Element {
	t := theme.Colors(c)
	size := float32(28)
	switch o.Size {
	case ProgressSmall:
		size = 20
	case ProgressLarge:
		size = 36
	}
	size = theme.Scale(c, size)
	lo, hi := progressBounds(o.Min, o.Max)
	label := o.Label
	if label == "" {
		if o.Indeterminate {
			label = locale.M().Loading
		} else {
			label = percentText(o.Value, o.Min, o.Max)
		}
	}
	if o.Indeterminate {
		return ui.Box(c).Size(size, size).Shrink(0).Children(func() {
			loading.Ring(c, loading.Options{Size: size, Color: indicatorFill(c, o.Color)}).Label(label)
		})
	}
	progress := progressFraction(o.Value, lo, hi)
	e := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleProgress).
		Range(lo, hi, o.Value).Label(label)
	// A ring has no box that could carry the value, so the native channel
	// eases the painted arc instead: frames are asked for only while it moves.
	progress = e.AnimateTo(progressRingKey{}, progress, indicatorAnimation)
	fill := indicatorFill(c, o.Color)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		drawRing(p, r, progress, fill, t.Default)
	})
	if o.ShowOutput {
		e.Children(func() {
			ui.Box(c).Size(size, size).Center().Children(func() {
				style.Text(ui.Text(c, percentText(o.Value, o.Min, o.Max)), c, theme.XS()).
					TextColor(t.Muted).FontFeatures("tnum")
			})
		})
	}
	return e
}

// progressRingKey is the animated channel of one determinate ring.
type progressRingKey struct{}

// ringDirections holds the unit-circle offsets of the progress ring's dots.
// They depend only on the segment count, not on a frame or a context.
var ringDirections = func() (dirs [90][2]float32) {
	const segments = 90
	for i := range dirs {
		angle := -1.5707963 + float32(i)*(6.2831853/segments)
		dirs[i] = [2]float32{float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))}
	}
	return dirs
}()

func drawRing(p *ui.Painter, r ui.Rect, progress float32, fill, track ui.Color) {
	const segments = 90
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	radius := min(r.W, r.H)/2 - 1
	if radius <= 0 {
		return
	}
	dot := radius * (6.2831853 / segments) * 1.15
	if dot < 2 {
		dot = 2
	}
	for i := range segments {
		dir := ringDirections[i]
		x := cx + radius*dir[0]
		y := cy + radius*dir[1]
		color := track
		if float32(i)/segments < progress {
			color = fill
		}
		p.Fill(ui.Rect{X: x - dot/2, Y: y - dot/2, W: dot, H: dot}, color, dot/2)
	}
}
