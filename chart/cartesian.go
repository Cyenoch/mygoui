package chart

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

// Point holds a numeric or Unix-time x coordinate and one measurement.
// A NaN Y is an explicit missing sample and breaks a line or area.
type Point struct {
	X, Y  float64
	Label string
}

type SeriesKind uint8

const (
	LineSeries SeriesKind = iota
	AreaSeries
	BarSeries
	ScatterSeries
	ProfitLossSeries
)

// Series belongs to the caller; ID is the stable legend identity.
type Series struct {
	ID, Label string
	Points    []Point
	Kind      SeriesKind
	Color     ui.Color
}

// Domain is an explicit increasing numeric range.
type Domain struct{ Min, Max float64 }

// Axis configures native tick labels. Time interprets X as Unix seconds.
type Axis struct {
	Hidden bool
	Ticks  int
	Format func(float64) string
	Time   bool
}

// PlotContext is build-scoped scale information for custom marks. Do not retain
// it or its Context/Element inputs. X and Y return plot fractions.
type PlotContext struct {
	X, Y     native.Linear
	Selected int
}

// CartesianOptions composes lines, areas, bars and scatter on one shared axis.
// Stacked requires aligned x coordinates for the bar/area series. Horizontal
// is supported for bars. Brush filters by data X before domains are computed.
type CartesianOptions struct {
	Options
	Series              []Series
	XAxis, YAxis        Axis
	XDomain, YDomain    *Domain
	Brush               *BrushState
	Stacked, Horizontal bool
	HideGrid            bool
	Background          BackgroundPattern
	Curve               native.Curve
	Annotations         []Annotation
	Build               func(*ui.Context, PlotContext)
}

func Line(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	return cartesian(c, o, LineSeries, true)
}

func Area(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	return cartesian(c, o, AreaSeries, true)
}

func Bar(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	return cartesian(c, o, BarSeries, true)
}

func Scatter(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	return cartesian(c, o, ScatterSeries, true)
}

func ProfitLoss(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	return cartesian(c, o, ProfitLossSeries, true)
}

func Composed(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	return cartesian(c, o, LineSeries, false)
}

// LiveLine shows the latest count samples from each caller-updated series.
// No goroutine or clock is created. Zero count shows all supplied samples.
func LiveLine(c *ui.Context, o CartesianOptions, count int) (ui.Element, error) {
	if count < 0 {
		return ui.Element{}, fmt.Errorf("live line: sample count must be nonnegative")
	}
	o.Series = slices.Clone(o.Series)
	for i := range o.Series {
		if count > 0 && len(o.Series[i].Points) > count {
			o.Series[i].Points = o.Series[i].Points[len(o.Series[i].Points)-count:]
		}
	}
	return Line(c, o)
}

// Sparkline omits axes and grid and supports line, area or bar layers.
func Sparkline(c *ui.Context, o CartesianOptions) (ui.Element, error) {
	o.XAxis.Hidden, o.YAxis.Hidden, o.HideGrid = true, true, true
	if o.Height == 0 {
		o.Height = 48
	}
	return Composed(c, o)
}

func cartesian(c *ui.Context, o CartesianOptions, kind SeriesKind, force bool) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	o.Series = slices.Clone(o.Series)
	for i := range o.Series {
		if o.Series[i].Color == (ui.Color{}) {
			o.Series[i].Color = SeriesColor(c, i)
		}
	}
	series, err := prepareSeries(o, kind, force)
	if err != nil {
		return ui.Element{}, err
	}
	if err := validateAnnotations(o.Annotations); err != nil {
		return ui.Element{}, err
	}
	if o.Curve < native.CurveLinear || o.Curve > native.CurveStep || o.Background > BackgroundLines {
		return ui.Element{}, fmt.Errorf("chart %q: unsupported curve or background", o.Label)
	}
	x, y, stacks, err := cartesianScales(o, series)
	if err != nil {
		return ui.Element{}, err
	}
	for _, annotation := range o.Annotations {
		if annotation.Kind == Weekends {
			if err := validateWeekendDomain(x); err != nil {
				return ui.Element{}, err
			}
		}
	}
	ctx := PlotContext{X: x, Y: y}
	t := theme.Colors(c)
	textSize := theme.Scale(c, 12)
	xs := seriesXs(series)
	describe := func(i int) string {
		parts := []string{}
		for _, s := range series {
			for _, pt := range s.Points {
				if pt.X == xs[i] && !math.IsNaN(pt.Y) {
					label := pt.Label
					if label == "" {
						label = axisFormat(o.XAxis, pt.X)
					}
					parts = append(parts, label+" · "+s.Label+": "+axisFormat(o.YAxis, pt.Y))
				}
			}
		}
		return strings.Join(parts, "; ")
	}
	axes := Axes(c, ctx, o.XAxis, o.YAxis)
	if o.Horizontal {
		axes = Axes(c, PlotContext{X: y, Y: x}, o.YAxis, o.XAxis)
	}
	barCount := 0
	for _, s := range series {
		if s.Kind == BarSeries {
			barCount++
		}
	}
	return plotFrame(c, o.Options, axes, len(xs), describe, func(area ui.Element, inner Options) {
		area.Children(func() {
			surface(c, inner, len(xs), describe, func(p *ui.Painter, r ui.Rect, selected int) {
				if o.Background != BackgroundSolid {
					Background(p, r, ui.Transparent, t.Border, o.Background)
				}
				if !o.HideGrid {
					gridCtx := ctx
					if o.Horizontal {
						gridCtx.X, gridCtx.Y = y, x
					}
					Grid(p, r, gridCtx, t.Border)
				}
				for _, annotation := range o.Annotations {
					drawAnnotation(p, r, ctx, annotation, t, textSize)
				}
				barIndex := 0
				for i, s := range series {
					color := s.Color
					if color == (ui.Color{}) {
						color = SeriesColor(c, i)
					}
					drawSeries(p, r, ctx, s, barIndex, barCount, o, stacks[i], color, t)
					if s.Kind == BarSeries {
						barIndex++
					}
					if selected >= 0 && o.Indicator == nil && !o.Horizontal {
						for _, pt := range s.Points {
							if pt.X != xs[selected] || math.IsNaN(pt.Y) {
								continue
							}
							a := point(r, x.At(pt.X), y.At(pt.Y))
							p.Fill(ui.Rect{X: a.X - 4, Y: a.Y - 4, W: 8, H: 8}, color, 4)
						}
					}
				}
				if selected >= 0 && o.Indicator == nil {
					v := x.At(xs[selected])
					if o.Horizontal {
						p.Line(r.X, native.Y(r, v), r.X+r.W, native.Y(r, v), 1, t.Ring)
					} else {
						p.Line(native.X(r, v), r.Y, native.X(r, v), r.Y+r.H, 1, t.Ring)
					}
				}
			}, func(px, py float32) int {
				if kind == ScatterSeries && force {
					best, distance := -1, math.Inf(1)
					r := o.State.Handle.Bounds(c)
					for _, s := range series {
						for _, pt := range s.Points {
							if math.IsNaN(pt.Y) {
								continue
							}
							dx, dy := float64((x.At(pt.X)-px)*r.W), float64((1-y.At(pt.Y)-py)*r.H)
							if d := dx*dx + dy*dy; d < distance {
								best, _ = slices.BinarySearch(xs, pt.X)
								distance = d
							}
						}
					}
					return best
				}
				if o.Horizontal {
					px = 1 - py
				}
				v := x.Invert(px)
				i, _ := slices.BinarySearch(xs, v)
				if i > 0 && (i == len(xs) || v-xs[i-1] < xs[i]-v) {
					i--
				}
				return min(len(xs)-1, i)
			})
			if o.Build != nil {
				ctx.Selected = -1
				if o.State != nil {
					ctx.Selected = o.State.Selected
				}
				o.Build(c, ctx)
			}
		})
	}), nil
}

// seriesXs returns the distinct X coordinates of series in increasing order.
func seriesXs(series []Series) []float64 {
	total := 0
	for _, s := range series {
		total += len(s.Points)
	}
	xs := make([]float64, 0, total)
	for _, s := range series {
		for _, pt := range s.Points {
			xs = append(xs, pt.X)
		}
	}
	if len(series) > 1 {
		// Each series is already increasing.
		slices.Sort(xs)
		xs = slices.Compact(xs)
	}
	return xs
}

func prepareSeries(o CartesianOptions, kind SeriesKind, force bool) ([]Series, error) {
	if o.Brush != nil && (!finite(o.Brush.From) || !finite(o.Brush.To) || o.Brush.To < o.Brush.From) {
		return nil, fmt.Errorf("chart %q: invalid brush domain", o.Label)
	}
	result := make([]Series, 0, len(o.Series))
	ids := make(map[string]bool, len(o.Series))
	for si, input := range o.Series {
		if input.ID == "" || ids[input.ID] {
			return nil, fmt.Errorf("chart %q: series %d needs a unique nonempty ID", o.Label, si)
		}
		ids[input.ID] = true
		s := input
		if force {
			s.Kind = kind
		}
		if s.Kind > ProfitLossSeries || (o.Horizontal && s.Kind != BarSeries) {
			return nil, fmt.Errorf("chart %q: unsupported series kind/orientation for %q", o.Label, s.ID)
		}
		// The caller's points are read, never copied: a build runs per frame.
		for i, pt := range input.Points {
			if !finite(pt.X) || math.IsInf(pt.Y, 0) || (i > 0 && pt.X <= input.Points[i-1].X) {
				return nil, fmt.Errorf("chart %q: series %q point %d needs finite increasing X and finite or missing Y", o.Label, s.ID, i)
			}
		}
		if o.Brush != nil {
			from, _ := slices.BinarySearchFunc(s.Points, o.Brush.From, func(pt Point, x float64) int { return cmp.Compare(pt.X, x) })
			to, found := slices.BinarySearchFunc(s.Points, o.Brush.To, func(pt Point, x float64) int { return cmp.Compare(pt.X, x) })
			if found {
				to++
			}
			s.Points = s.Points[from:max(from, to)]
		}
		if o.State == nil || !o.State.Hidden[s.ID] {
			result = append(result, s)
		}
	}
	return result, nil
}

// extent is the running range of the values added to it, skipping NaN.
type extent struct{ lo, hi float64 }

func newExtent() extent { return extent{math.Inf(1), math.Inf(-1)} }

func (e *extent) add(values ...float64) {
	for _, v := range values {
		if !math.IsNaN(v) {
			e.lo, e.hi = min(e.lo, v), max(e.hi, v)
		}
	}
}

func cartesianScales(o CartesianOptions, series []Series) (native.Linear, native.Linear, [][]native.Segment, error) {
	xe, ye := newExtent(), newExtent()
	ye.add(0)
	var values [][]float64
	if o.Stacked {
		values = make([][]float64, len(series))
	}
	for i, s := range series {
		if o.Stacked {
			if i > 0 && len(s.Points) != len(series[0].Points) {
				return native.Linear{}, native.Linear{}, nil, fmt.Errorf("chart %q: stacked series must have aligned X coordinates", o.Label)
			}
			if s.Kind != BarSeries && s.Kind != AreaSeries {
				return native.Linear{}, native.Linear{}, nil, fmt.Errorf("chart %q: stacking requires bar or area series", o.Label)
			}
			values[i] = make([]float64, len(s.Points))
		}
		for j, pt := range s.Points {
			if o.Stacked {
				if i > 0 && pt.X != series[0].Points[j].X {
					return native.Linear{}, native.Linear{}, nil, fmt.Errorf("chart %q: stacked series must have aligned X coordinates", o.Label)
				}
				values[i][j] = pt.Y
			}
			xe.add(pt.X)
			ye.add(pt.Y)
		}
	}
	for _, a := range o.Annotations {
		if a.Kind == ProjectionLine {
			xe.add(a.X1, a.X2)
			ye.add(a.Y1, a.Y2)
		}
		if a.Kind == Band {
			for i, pt := range a.Points {
				xe.add(pt.X)
				ye.add(pt.Y, a.Upper[i].Y)
			}
		}
	}
	stacks := make([][]native.Segment, len(series))
	if o.Stacked {
		stacks = native.Stack(values, native.StackStandard)
		for _, stack := range stacks {
			for _, v := range stack {
				if !finite(v.Lo) || !finite(v.Hi) {
					return native.Linear{}, native.Linear{}, nil, fmt.Errorf("chart %q: stacked values overflow", o.Label)
				}
				ye.add(v.Lo, v.Hi)
			}
		}
	}
	x, err := scaleExtent(xe, o.XDomain)
	if err != nil {
		return x, native.Linear{}, nil, err
	}
	if o.XDomain == nil && xe.lo < xe.hi {
		step := math.Inf(1)
		for _, s := range series {
			if s.Kind == BarSeries {
				for i := 1; i < len(s.Points); i++ {
					step = min(step, s.Points[i].X-s.Points[i-1].X)
				}
			}
		}
		if finite(step) {
			x.Lo -= step / 2
			x.Hi += step / 2
		}
	}
	y, err := scaleExtent(ye, o.YDomain)
	return x, y, stacks, err
}

func scale(values []float64, domain *Domain) (native.Linear, error) {
	e := newExtent()
	e.add(values...)
	return scaleExtent(e, domain)
}

func scaleExtent(e extent, domain *Domain) (native.Linear, error) {
	lo, hi := e.lo, e.hi
	if lo > hi {
		lo, hi = 0, 1
	}
	if lo == hi {
		lo, hi = lo-0.5, hi+0.5
	}
	if domain != nil {
		lo, hi = domain.Min, domain.Max
		if !finite(lo) || !finite(hi) || hi <= lo {
			return native.Linear{}, fmt.Errorf("chart: domain must be finite and increasing")
		}
	}
	if !finite(lo) || !finite(hi) || !finite(hi-lo) || hi <= lo {
		return native.Linear{}, fmt.Errorf("chart: numeric extent is not representable")
	}
	return native.NewLinear(lo, hi), nil
}

func axisFormat(a Axis, v float64) string {
	if a.Format != nil {
		return a.Format(v)
	}
	if a.Time {
		return time.Unix(int64(v), 0).UTC().Format("Jan 02")
	}
	return number(v)
}

// Axes returns styled native axes for a custom Plot as well as built-in charts.
func Axes(c *ui.Context, ctx PlotContext, x, y Axis) native.Axes {
	ticks := func(s native.Linear, a Axis) []native.Tick {
		if a.Hidden {
			return nil
		}
		n := a.Ticks
		if n <= 0 {
			n = 4
		}
		s.Format = func(v float64) string { return axisFormat(a, v) }
		return s.Ticks(n)
	}
	return native.Axes{Left: ticks(ctx.Y, y), Bottom: ticks(ctx.X, x),
		Label: func(_ native.Side, tick native.Tick) {
			style.Text(ui.Text(c, tick.Label), c, theme.XS()).NoWrap().TextColor(theme.Colors(c).MutedForeground)
		}}
}

func drawSeries(p *ui.Painter, r ui.Rect, ctx PlotContext, s Series, index, count int,
	o CartesianOptions, stack []native.Segment, color ui.Color, t theme.Palette,
) {
	xs, ys, low := make([]float32, len(s.Points)), make([]float32, len(s.Points)), make([]float32, len(s.Points))
	for i, pt := range s.Points {
		xs[i], ys[i], low[i] = ctx.X.At(pt.X), ctx.Y.At(pt.Y), ctx.Y.At(0)
		if len(stack) > 0 {
			ys[i], low[i] = ctx.Y.At(stack[i].Hi), ctx.Y.At(stack[i].Lo)
			if math.IsNaN(pt.Y) {
				ys[i], low[i] = float32(math.NaN()), float32(math.NaN())
			}
		}
	}
	switch s.Kind {
	case AreaSeries:
		if len(stack) == 0 && len(xs) > int(r.W)*4 {
			// Unstacked areas share one baseline, so the outline can thin.
			xs, ys = native.MinMax(nil, nil, xs, ys, max(1, int(r.W)))
			low = low[:len(xs)]
		}
		p.FillPath(native.Area(r, xs, low, ys, o.Curve), color.Alpha(0.18))
		p.StrokePath(native.Line(r, xs, ys, o.Curve), 2, color)
	case LineSeries:
		drawX, drawY := xs, ys
		if len(xs) > int(r.W)*4 {
			drawX, drawY = native.MinMax(nil, nil, xs, ys, max(1, int(r.W)))
		}
		p.StrokePath(native.Line(r, drawX, drawY, o.Curve), 2, color)
		if len(xs) == 1 {
			a := point(r, xs[0], ys[0])
			p.Fill(ui.Rect{X: a.X - 3, Y: a.Y - 3, W: 6, H: 6}, color, 3)
		}
	case ProfitLossSeries:
		for i := 1; i < len(xs); i++ {
			a, b := s.Points[i-1], s.Points[i]
			if math.IsNaN(a.Y) || math.IsNaN(b.Y) {
				continue
			}
			fill := t.Success
			if a.Y < 0 {
				fill = t.Destructive
			}
			pa, pb := point(r, xs[i-1], ys[i-1]), point(r, xs[i], ys[i])
			if (a.Y < 0) != (b.Y < 0) {
				zero := a.X + (b.X-a.X)*(-a.Y)/(b.Y-a.Y)
				mid := point(r, ctx.X.At(zero), ctx.Y.At(0))
				line(p, pa, mid, 2, fill)
				fill = t.Success
				if b.Y < 0 {
					fill = t.Destructive
				}
				line(p, mid, pb, 2, fill)
				continue
			}
			line(p, pa, pb, 2, fill)
		}
	case ScatterSeries:
		for i := range xs {
			if !math.IsNaN(s.Points[i].Y) {
				a := point(r, xs[i], ys[i])
				p.Fill(ui.Rect{X: a.X - 4, Y: a.Y - 4, W: 8, H: 8}, color, 4)
			}
		}
	case BarSeries:
		step := float32(1)
		for i := 1; i < len(xs); i++ {
			step = min(step, xs[i]-xs[i-1])
		}
		width, offset := step*0.7, float32(0)
		if !o.Stacked {
			width /= float32(max(1, count))
			offset = width * (float32(index) - float32(count-1)/2)
		}
		for i := range xs {
			if math.IsNaN(s.Points[i].Y) {
				continue
			}
			a, b := point(r, xs[i]+offset-width/2, low[i]), point(r, xs[i]+offset+width/2, ys[i])
			if o.Horizontal {
				a, b = point(r, low[i], xs[i]+offset-width/2), point(r, ys[i], xs[i]+offset+width/2)
			}
			p.Fill(ui.Rect{X: min(a.X, b.X), Y: min(a.Y, b.Y), W: max(1, abs(a.X-b.X)), H: max(1, abs(a.Y-b.Y))}, color, 2)
		}
	}
}

func abs(v float32) float32 { return float32(math.Abs(float64(v))) }
