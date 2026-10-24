package mygoui

import (
	"math"
	"strconv"
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// IndicatorColor is the color of a progress or meter indicator. Its zero
// value is HeroUI's default accent, unlike the status colors used by chips.
type IndicatorColor uint8

const (
	IndicatorAccent IndicatorColor = iota
	IndicatorDefault
	IndicatorSuccess
	IndicatorWarning
	IndicatorDanger
)

// ProgressSize is HeroUI's meter/progress scale.
type ProgressSize uint8

const (
	ProgressMedium ProgressSize = iota
	ProgressSmall
	ProgressLarge
)

// ProgressOptions configures a HeroUI bar indicator.
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
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4).
		Role(ui.RoleProgress).Label(o.Label)
	if o.Indeterminate {
		e.Range(0, 1, -1)
	} else {
		e.Range(lo, hi, o.Value)
	}
	e.Children(func() {
		progressHeading(c, o.Label, o.ShowOutput, progressText(o))
		height, radius := progressMetrics(o.Size)
		track := ui.Box(c).FillWidth().Height(height).Radius(radius).Clip().
			Background(t.Default)
		if o.Indeterminate {
			start := ui.Local(track, progressStartKey{}, func() time.Time { return c.Now() })
			track.DrawOver(func(p *ui.Painter, r ui.Rect) {
				band := r.W * 0.4
				if theme.ReduceMotion(c) {
					p.Fill(ui.Rect{X: r.X + r.W*0.3, Y: r.Y, W: band, H: r.H},
						indicatorFill(c, o.Color), radius)
					return
				}
				const period = 1500 * time.Millisecond
				phase := float32(p.Now().Sub(*start)%period) / float32(period)
				x := r.X - band + phase*(r.W+band)
				p.Fill(ui.Rect{X: x, Y: r.Y, W: band, H: r.H}, indicatorFill(c, o.Color), radius)
				p.AnimationFrame()
			})
		} else {
			pct := progressFraction(o.Value, lo, hi)
			track.Children(func() {
				ui.Box(c).Height(height).WidthPercent(pct * 100).
					Background(indicatorFill(c, o.Color)).Radius(radius)
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

// Meter renders a HeroUI meter with its label and output.
func Meter(c *ui.Context, o MeterOptions) ui.Element {
	t := theme.Colors(c)
	lo, hi := progressBounds(o.Min, o.Max)
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4).
		Role(ui.RoleMeter).Label(o.Label).Range(lo, hi, o.Value)
	e.Children(func() {
		progressHeading(c, o.Label, o.ShowOutput, meterText(o))
		height, radius := progressMetrics(o.Size)
		pct := progressFraction(o.Value, lo, hi)
		ui.Box(c).FillWidth().Height(height).Radius(radius).Clip().Background(t.Default).
			Children(func() {
				ui.Box(c).Height(height).WidthPercent(pct * 100).
					Background(indicatorFill(c, o.Color)).Radius(radius)
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
			style.Text(ui.Text(c, output), c, ts).TextColor(theme.Colors(c).Muted)
		}
	})
}

func progressMetrics(size ProgressSize) (height, radius float32) {
	height, radius = theme.Space8, theme.RadiusScale().SM
	switch size {
	case ProgressSmall:
		height, radius = theme.Space4, theme.RadiusScale().XS
	case ProgressLarge:
		height, radius = theme.Space12, theme.RadiusScale().MD
	}
	return height, radius
}

type progressStartKey struct{}

func progressText(o ProgressOptions) string {
	if o.OutputText != "" {
		return o.OutputText
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
	return t.Accent
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

// ProgressCircle draws a determinate ring from native paint segments or an
// indeterminate native spinner. The value is exposed as accessible text.
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
		label = percentText(o.Value, o.Min, o.Max)
	}
	if o.Indeterminate {
		return ui.Box(c).Size(size, size).Shrink(0).Children(func() {
			loading.Ring(c, loading.Options{Size: size, Color: indicatorFill(c, o.Color)}).Label(label)
		})
	}
	progress := progressFraction(o.Value, lo, hi)
	e := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleProgress).
		Range(lo, hi, o.Value).Label(label)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		drawRing(p, r, progress, indicatorFill(c, o.Color), t.Default)
	})
	if o.ShowOutput {
		e.Children(func() {
			ui.Box(c).Size(size, size).Center().Children(func() {
				style.Text(ui.Text(c, percentText(o.Value, o.Min, o.Max)), c, theme.XS()).
					TextColor(t.Muted)
			})
		})
	}
	return e
}

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
