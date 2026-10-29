package chart

import (
	"fmt"
	"math"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

// Value is a labeled quantity. Max is used by progress rings (default 100).
type Value struct {
	ID, Label  string
	Value, Max float64
	Color      ui.Color
}

type PieOptions struct {
	Options
	Values []Value
	// InnerRadius is a fraction of the outer radius, 0 for pie, e.g. .65 for donut.
	InnerRadius float32
}

func validateValues(values []Value, progress bool) error {
	ids := map[string]bool{}
	for i, v := range values {
		if v.ID == "" || ids[v.ID] || !finite(v.Value) || v.Value < 0 {
			return fmt.Errorf("chart: value %d needs a unique ID and finite nonnegative quantity", i)
		}
		ids[v.ID] = true
		if progress && (!finite(v.Max) || v.Max < 0 || (v.Max > 0 && v.Value > v.Max) || (v.Max == 0 && v.Value > 100)) {
			return fmt.Errorf("chart: value %q must be within its progress range", v.ID)
		}
	}
	return nil
}

func visibleValues(c *ui.Context, values []Value, state *State) []Value {
	result := []Value{}
	for i, v := range values {
		if state == nil || !state.Hidden[v.ID] {
			if v.Color == (ui.Color{}) {
				v.Color = SeriesColor(c, i)
			}
			result = append(result, v)
		}
	}
	return result
}

func valueColor(c *ui.Context, v Value, i int) ui.Color {
	if v.Color != (ui.Color{}) {
		return v.Color
	}
	return SeriesColor(c, i)
}

func Pie(c *ui.Context, o PieOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	if err := validateValues(o.Values, false); err != nil {
		return ui.Element{}, err
	}
	if !finite(float64(o.InnerRadius)) || o.InnerRadius < 0 || o.InnerRadius >= 1 {
		return ui.Element{}, fmt.Errorf("pie: inner radius must be within [0,1)")
	}
	values := visibleValues(c, o.Values, o.State)
	numbers := []float64{}
	total := 0.0
	for _, v := range values {
		numbers = append(numbers, v.Value)
		total += v.Value
	}
	if !finite(total) {
		return ui.Element{}, fmt.Errorf("pie: total value overflows")
	}
	slices := native.Pie(numbers, 0.02)
	count := len(values)
	if total == 0 {
		count = 0
	}
	describe := func(i int) string {
		return fmt.Sprintf("%s: %s (%.1f%%)", values[i].Label, number(values[i].Value), values[i].Value/total*100)
	}
	return surface(c, o.Options, count, describe, func(p *ui.Painter, r ui.Rect, selected int) {
		cx, cy, radius := circle(r)
		for i, s := range slices {
			outer := radius
			if selected == i {
				outer += 3
			}
			p.FillPath(native.Sector(cx, cy, radius*o.InnerRadius, outer, s), valueColor(c, values[i], i))
		}
	}, func(x, y float32) int {
		// Preserve circular hit geometry in rectangular plots.
		w, h := o.State.Handle.Bounds(c).W, o.State.Handle.Bounds(c).H
		radius := min(w, h)/2 - 8
		i, ok := native.SliceAt(slices, (x-.5)*w, (y-.5)*h, radius*o.InnerRadius, radius+3)
		if !ok {
			return -1
		}
		return i
	}), nil
}

func circle(r ui.Rect) (float32, float32, float32) {
	return r.X + r.W/2, r.Y + r.H/2, max(0, min(r.W, r.H)/2-8)
}

type GaugeOptions struct {
	Options
	Value, Min, Max float64
	Linear          bool
	Notches         int
}

// Gauge uses a notched track. Zero Min/Max means 0..100.
func Gauge(c *ui.Context, o GaugeOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	if o.Min == 0 && o.Max == 0 {
		o.Max = 100
	}
	if !finite(o.Value) || !finite(o.Min) || !finite(o.Max) || o.Max <= o.Min || o.Value < o.Min || o.Value > o.Max {
		return ui.Element{}, fmt.Errorf("gauge: value must be within a finite increasing domain")
	}
	if o.Notches == 0 {
		o.Notches = 40
	}
	if o.Notches < 2 || o.Notches > 1000 {
		return ui.Element{}, fmt.Errorf("gauge: notches must be between 2 and 1000")
	}
	t := theme.Colors(c)
	fraction := (o.Value - o.Min) / (o.Max - o.Min)
	return surface(c, o.Options, 1, func(int) string { return o.Label + ": " + number(o.Value) },
		func(p *ui.Painter, r ui.Rect, _ int) {
			cx, cy, radius := circle(r)
			for i := range o.Notches {
				fill := t.Input
				if float64(i)/float64(o.Notches) < fraction {
					fill = t.Chart1
				}
				if o.Linear {
					width := r.W / float32(o.Notches)
					p.Fill(ui.Rect{X: r.X + float32(i)*width, Y: cy - 12, W: max(1, width-2), H: 24}, fill, 2)
					continue
				}
				start := -math.Pi*.75 + float64(i)/float64(o.Notches)*math.Pi*1.5
				end := start + math.Pi*1.5/float64(o.Notches)*.65
				p.FillPath(native.Sector(cx, cy, radius*.8, radius, native.Slice{Start: start, End: end}), fill)
			}
		}, func(float32, float32) int { return 0 }), nil
}

type RingOptions struct {
	Options
	Values []Value
}

func Ring(c *ui.Context, o RingOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	if err := validateValues(o.Values, true); err != nil {
		return ui.Element{}, err
	}
	values := visibleValues(c, o.Values, o.State)
	t := theme.Colors(c)
	return surface(c, o.Options, len(values), func(i int) string {
		return values[i].Label + ": " + number(values[i].Value)
	}, func(p *ui.Painter, r ui.Rect, selected int) {
		cx, cy, radius := circle(r)
		step := radius * .75 / float32(max(1, len(values)))
		for i, v := range values {
			outer := radius - float32(i)*step
			inner := outer - step*.7
			p.FillPath(native.Sector(cx, cy, inner, outer, native.Slice{Start: 0, End: 2 * math.Pi}), t.Input)
			maxValue := v.Max
			if maxValue == 0 {
				maxValue = 100
			}
			fill := valueColor(c, v, i)
			if selected >= 0 && selected != i {
				fill = fill.Alpha(.4)
			}
			p.FillPath(native.Sector(cx, cy, inner, outer, native.Slice{Start: 0, End: (v.Value / maxValue) * (2 * math.Pi)}), fill)
		}
	}, func(x, y float32) int {
		r := o.State.Handle.Bounds(c)
		rad := min(r.W, r.H)/2 - 8
		distance := float32(math.Hypot(float64((x-.5)*r.W), float64((y-.5)*r.H)))
		i := int((rad - distance) / (rad * .75 / float32(max(1, len(values)))))
		if distance > rad || i < 0 || i >= len(values) {
			return -1
		}
		return i
	}), nil
}

type RadarSeries struct {
	ID, Label string
	Values    []float64
	Color     ui.Color
}

type RadarOptions struct {
	Options
	Metrics []string
	Series  []RadarSeries
	Max     float64
}

func Radar(c *ui.Context, o RadarOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	if len(o.Metrics) > 0 && len(o.Metrics) < 3 {
		return ui.Element{}, fmt.Errorf("radar: at least three metrics are required")
	}
	maxValue := o.Max
	if !finite(maxValue) || maxValue < 0 {
		return ui.Element{}, fmt.Errorf("radar: maximum must be finite and nonnegative")
	}
	for si, s := range o.Series {
		if len(s.Values) != len(o.Metrics) {
			return ui.Element{}, fmt.Errorf("radar: series %d must match metric count", si)
		}
		for _, v := range s.Values {
			if !finite(v) || v < 0 || (o.Max > 0 && v > o.Max) {
				return ui.Element{}, fmt.Errorf("radar: invalid value in series %d", si)
			}
			maxValue = max(maxValue, v)
		}
	}
	maxValue = max(1, maxValue)
	t := theme.Colors(c)
	metricLabels := make([]*ui.TextLayout, len(o.Metrics))
	clearance := theme.Scale(c, 24)
	for i, text := range o.Metrics {
		metricLabels[i] = ui.ShapeText(text, ui.Font{Family: theme.FontFamily(c), Size: theme.Scale(c, 12)}, 0)
		w, h := metricLabels[i].Size()
		clearance = max(clearance, w/2+theme.Scale(c, 12), h+theme.Scale(c, 8))
	}
	return surface(c, o.Options, len(o.Metrics), func(i int) string {
		text := o.Metrics[i]
		for _, s := range o.Series {
			text += " · " + s.Label + ": " + number(s.Values[i])
		}
		return text
	}, func(p *ui.Painter, r ui.Rect, selected int) {
		cx, cy, radius := circle(r)
		radius = max(0, radius-clearance)
		for level := 1; level <= 4; level++ {
			path := &ui.Path{}
			for i := range o.Metrics {
				x, y := native.Polar(cx, cy, radius*float32(level)/4, float64(i)/float64(len(o.Metrics))*2*math.Pi)
				if i == 0 {
					path.MoveTo(x, y)
				} else {
					path.LineTo(x, y)
				}
			}
			p.StrokePath(path.Close(), 1, t.Border)
		}
		for i := range o.Metrics {
			angle := float64(i) / float64(len(o.Metrics)) * 2 * math.Pi
			x, y := native.Polar(cx, cy, radius, angle)
			line(p, ui.Point{X: cx, Y: cy}, ui.Point{X: x, Y: y}, 1, t.Border)
			lx, ly := native.Polar(cx, cy, radius+theme.Scale(c, 12), angle)
			w, h := metricLabels[i].Size()
			if math.Sin(angle) > .2 {
				lx += w / 2
			} else if math.Sin(angle) < -.2 {
				lx -= w / 2
			}
			if math.Cos(angle) > .2 {
				ly -= h / 2
			} else if math.Cos(angle) < -.2 {
				ly += h / 2
			}
			p.TextLayout(metricLabels[i], max(r.X, min(r.X+r.W-w, lx-w/2)), max(r.Y, min(r.Y+r.H-h, ly-h/2)), t.MutedForeground)
		}
		for si, s := range o.Series {
			if o.State != nil && o.State.Hidden[s.ID] {
				continue
			}
			path := &ui.Path{}
			fill := s.Color
			if fill == (ui.Color{}) {
				fill = SeriesColor(c, si)
			}
			for i, v := range s.Values {
				x, y := native.Polar(cx, cy, radius*float32(v/maxValue), float64(i)/float64(len(o.Metrics))*2*math.Pi)
				if i == 0 {
					path.MoveTo(x, y)
				} else {
					path.LineTo(x, y)
				}
				if selected == i {
					p.Fill(ui.Rect{X: x - 4, Y: y - 4, W: 8, H: 8}, fill, 4)
				}
			}
			path.Close()
			p.FillPath(path, fill.Alpha(.15))
			p.StrokePath(path, 2, fill)
		}
	}, angularHit(c, o.Options, len(o.Metrics), true)), nil
}

type RadialKind uint8

const (
	RadialBar RadialKind = iota
	RadialLine
	RadialArea
	RadialDots
)

type RadialOptions struct {
	Options
	Values []Value
	Kind   RadialKind
}

func Radial(c *ui.Context, o RadialOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	if err := validateValues(o.Values, false); err != nil {
		return ui.Element{}, err
	}
	if o.Kind > RadialDots {
		return ui.Element{}, fmt.Errorf("radial: unsupported layer kind")
	}
	values := visibleValues(c, o.Values, o.State)
	maxValue := 1.0
	for _, v := range values {
		maxValue = max(maxValue, v.Value)
	}
	return surface(c, o.Options, len(values), func(i int) string {
		return values[i].Label + ": " + number(values[i].Value)
	}, func(p *ui.Painter, r ui.Rect, selected int) {
		cx, cy, radius := circle(r)
		path := &ui.Path{}
		for i, v := range values {
			a := float64(i) / float64(len(values)) * 2 * math.Pi
			rad := radius * (.2 + .8*float32(v.Value/maxValue))
			color := valueColor(c, v, i)
			if selected >= 0 && selected != i {
				color = color.Alpha(.45)
			}
			x, y := native.Polar(cx, cy, rad, a)
			switch o.Kind {
			case RadialBar:
				p.FillPath(native.Sector(cx, cy, radius*.2, rad, native.Slice{Start: a, End: a + 2*math.Pi/float64(len(values)), Pad: .04}), color)
			case RadialDots:
				for step := radius * .2; step <= rad; step += 9 {
					dx, dy := native.Polar(cx, cy, step, a)
					p.Fill(ui.Rect{X: dx - 3, Y: dy - 3, W: 6, H: 6}, color, 3)
				}
			default:
				if i == 0 {
					path.MoveTo(x, y)
				} else {
					path.LineTo(x, y)
				}
			}
		}
		if o.Kind == RadialLine || o.Kind == RadialArea {
			path.Close()
			if o.Kind == RadialArea {
				p.FillPath(path, SeriesColor(c, 0).Alpha(.15))
			}
			p.StrokePath(path, 2, SeriesColor(c, 0))
		}
	}, angularHit(c, o.Options, len(values), o.Kind != RadialBar)), nil
}

func angularHit(c *ui.Context, o Options, count int, nearest bool) func(float32, float32) int {
	return func(x, y float32) int {
		r := o.State.Handle.Bounds(c)
		a := math.Atan2(float64((x-.5)*r.W), float64((.5-y)*r.H))
		if a < 0 {
			a += 2 * math.Pi
		}
		position := a / (2 * math.Pi) * float64(count)
		if nearest {
			return int(math.Round(position)) % max(1, count)
		}
		return min(count-1, int(position))
	}
}

type FunnelOptions struct {
	Options
	Values []Value
}

func Funnel(c *ui.Context, o FunnelOptions) (ui.Element, error) {
	if err := checkOptions(o.Options); err != nil {
		return ui.Element{}, err
	}
	if err := validateValues(o.Values, false); err != nil {
		return ui.Element{}, err
	}
	values := visibleValues(c, o.Values, o.State)
	maxValue := 0.0
	for _, v := range values {
		maxValue = max(maxValue, v.Value)
	}
	count := len(values)
	if maxValue == 0 {
		count = 0
	}
	return surface(c, o.Options, count, func(i int) string {
		return values[i].Label + ": " + number(values[i].Value)
	}, func(p *ui.Painter, r ui.Rect, selected int) {
		step := r.H / float32(max(1, count))
		for i, v := range values {
			width := r.W * float32(v.Value/maxValue) * .9
			bottom := width
			if i+1 < len(values) {
				bottom = r.W * float32(values[i+1].Value/maxValue) * .9
			}
			y := r.Y + float32(i)*step
			cx := r.X + r.W/2
			path := (&ui.Path{}).MoveTo(cx-width/2, y+2).LineTo(cx+width/2, y+2).
				LineTo(cx+bottom/2, y+step-2).LineTo(cx-bottom/2, y+step-2).Close()
			fill := valueColor(c, v, i)
			if selected >= 0 && selected != i {
				fill = fill.Alpha(.45)
			}
			p.FillPath(path, fill)
		}
	}, func(_, y float32) int { return min(count-1, int(y*float32(count))) }), nil
}
