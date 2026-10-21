// Package chart composes maintained MyGo plots, scales, marks and selection.
// Polar extensions use its public Polar/Sector and native Path primitives.
package chart

import (
	"math"
	"slices"
	"strings"

	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
	native "github.com/egoist/mygo/ui/chart"
)

type ChartKind uint8

const (
	Bar ChartKind = iota
	Line
	Area
	Pie
	Radar
	Radial
)

type SeriesID string
type SeriesColorToken uint8

const (
	Chart1 SeriesColorToken = iota
	Chart2
	Chart3
	Chart4
	Chart5
)

type Value struct {
	Number  float64
	Missing bool
}
type Series struct {
	ID     SeriesID
	Label  string
	Color  SeriesColorToken
	Values []Value
}
type Data struct {
	Labels []string
	Series []Series
}
type StackMode uint8

const (
	Unstacked StackMode = iota
	Stacked
	Normalized
	Centered
)

// State is caller-owned presentation state, one instance per chart/window.
// Selected is the category index, or for pie/radial the visible series index.
// Set Selected to -1 for no initial reading. Geometry is native window-local state.
type State struct {
	Selected int
	Hidden   map[SeriesID]bool
}
type Options struct {
	Kind                              ChartKind
	Data                              Data
	Label, Description                string
	Height                            float32
	Stack                             StackMode
	HideAxes, HideLegend, HideTooltip bool
	OnSelection                       func(int)
	OnVisibility                      func(SeriesID, bool)
}

type run struct {
	from, to   int
	line, area ui.Path
}
type preparedSeries struct {
	source    int
	color     ui.Color
	low, high []float32
	runs      []run
	radar     ui.Path
	fill      bool
}
type cache struct {
	data                  Data
	hidden                []bool
	kind                  ChartKind
	stack                 StackMode
	ready                 bool
	x                     native.Band
	y                     native.Linear
	xs                    []float32
	series                []preparedSeries
	values                [][]float64
	segments              [][]native.Segment
	slices                []native.Slice
	polarValues           []float64
	polarXS               []float32
	ticks, bottomTicks    []native.Tick
	polarPaths, tracks    []ui.Path
	grid                  ui.Path
	rect                  ui.Rect
	geometryReady         bool
	message               string
	locale                string
	readings              []string
	hasValues, hasNonzero bool
}

func present(v Value) bool { return !v.Missing && !math.IsNaN(v.Number) && !math.IsInf(v.Number, 0) }
func valueAt(s Series, i int) (float64, bool) {
	if i < 0 || i >= len(s.Values) || !present(s.Values[i]) {
		return 0, false
	}
	return s.Values[i].Number, true
}
func (k *cache) unchanged(s *State, o Options) bool {
	if !k.ready || k.kind != o.Kind || k.stack != o.Stack || !slices.Equal(k.data.Labels, o.Data.Labels) || len(k.data.Series) != len(o.Data.Series) {
		return false
	}
	for i, series := range o.Data.Series {
		old := k.data.Series[i]
		if series.ID != old.ID || series.Label != old.Label || series.Color != old.Color || k.hidden[i] != s.Hidden[series.ID] || !slices.EqualFunc(series.Values, old.Values, func(a, b Value) bool {
			return a.Missing == b.Missing && (a.Number == b.Number || (math.IsNaN(a.Number) && math.IsNaN(b.Number)))
		}) {
			return false
		}
	}
	return true
}
func (k *cache) update(s *State, o Options, locale string) {
	if k.unchanged(s, o) && k.locale == locale {
		return
	}
	k.ready = true
	k.locale = locale
	k.kind = o.Kind
	k.stack = o.Stack
	k.geometryReady = false
	k.message = ""
	k.hasValues = false
	k.hasNonzero = false
	k.data.Labels = append(k.data.Labels[:0], o.Data.Labels...)
	if cap(k.data.Series) < len(o.Data.Series) {
		k.data.Series = make([]Series, len(o.Data.Series))
	} else {
		k.data.Series = k.data.Series[:len(o.Data.Series)]
	}
	k.hidden = k.hidden[:0]
	k.series = k.series[:0]
	k.values = k.values[:0]
	k.polarValues = k.polarValues[:0]
	for i, series := range o.Data.Series {
		oldValues := k.data.Series[i].Values
		k.data.Series[i] = series
		k.data.Series[i].Values = append(oldValues[:0], series.Values...)
		hidden := s.Hidden[series.ID]
		k.hidden = append(k.hidden, hidden)
		if hidden {
			continue
		}
		k.series = append(k.series, preparedSeries{source: i})
		numbers := make([]float64, len(o.Data.Labels))
		for j := range numbers {
			if n, ok := valueAt(series, j); ok {
				numbers[j] = n
				if (o.Kind != Pie && o.Kind != Radial) || j == 0 {
					k.hasValues = true
					k.hasNonzero = k.hasNonzero || n != 0
				}
				if o.Kind >= Pie && n < 0 {
					k.message = locales.ComponentPolarNonnegative()
				}
			}
		}
		k.values = append(k.values, numbers)
		v, _ := valueAt(series, 0)
		k.polarValues = append(k.polarValues, max(0, v))
	}
	k.x = native.NewBand(o.Data.Labels...).Padding(0.25, 0.15)
	k.xs = make([]float32, len(o.Data.Labels))
	for i := range k.xs {
		k.xs[i] = k.x.At(i)
	}
	k.polarXS = make([]float32, len(k.series))
	for i := range k.polarXS {
		k.polarXS[i] = (float32(i) + 0.5) / float32(len(k.series))
	}
	k.segments = nil
	if o.Stack != Unstacked && (o.Kind == Bar || o.Kind == Area) {
		how := native.StackStandard
		if o.Stack == Normalized {
			how = native.StackNormalized
		}
		if o.Stack == Centered {
			how = native.StackCenter
		}
		k.segments = native.Stack(k.values, how)
	}
	lo, hi := 0.0, 0.0
	for si, series := range k.series {
		for i := range k.xs {
			n, ok := valueAt(o.Data.Series[series.source], i)
			if !ok {
				continue
			}
			if k.segments != nil {
				seg := k.segments[si][i]
				lo = min(lo, seg.Lo, seg.Hi)
				hi = max(hi, seg.Lo, seg.Hi)
			} else {
				lo = min(lo, n)
				hi = max(hi, n)
			}
		}
	}
	// A zero domain needs a finite scale, not fabricated samples or NaN endpoints.
	if lo == hi {
		hi = lo + 1
	}
	k.y = native.NewLinear(lo, hi).Nice(4)
	for si := range k.series {
		ps := &k.series[si]
		ps.low = make([]float32, len(k.xs))
		ps.high = make([]float32, len(k.xs))
		start := -1
		for i := 0; i <= len(k.xs); i++ {
			var n float64
			ok := false
			if i < len(k.xs) {
				n, ok = valueAt(o.Data.Series[ps.source], i)
			}
			if ok {
				low, high := 0.0, n
				if k.segments != nil {
					low, high = k.segments[si][i].Lo, k.segments[si][i].Hi
				}
				ps.low[i], ps.high[i] = k.y.At(low), k.y.At(high)
				if start < 0 {
					start = i
				}
			} else if start >= 0 {
				ps.runs = append(ps.runs, run{from: start, to: i})
				start = -1
			}
		}
	}
	k.slices = native.Pie(k.polarValues, 0)
	k.ticks = k.y.Ticks(4)
	for i := range k.ticks {
		k.ticks[i].Label = i18n.FormatNumber(locale, k.ticks[i].Value, i18n.NumberOptions{MaxFractionDigits: 100})
	}
	k.bottomTicks = k.x.Ticks()
	if o.Kind == Radar && len(k.xs) < 3 {
		k.message = locales.ComponentRadarMinimumCategories()
	}
	if len(k.series) == 0 && len(o.Data.Series) > 0 {
		k.message = locales.ComponentNoVisibleSeries()
	} else if !k.hasValues {
		k.message = locales.ComponentNoData()
	} else if !k.hasNonzero && k.message == "" {
		k.message = locales.ComponentAllValuesZero()
	}
	k.polarPaths = make([]ui.Path, len(k.series))
	k.tracks = make([]ui.Path, len(k.series))
	count := len(o.Data.Labels)
	if o.Kind == Pie || o.Kind == Radial {
		count = len(k.series)
	}
	k.readings = make([]string, count)
	for i := range k.readings {
		k.readings[i] = describe(k, o, i)
	}
}
func seriesColor(c *ui.Context, token SeriesColorToken) ui.Color {
	p := theme.Colors(c)
	switch token {
	case Chart2:
		return p.Chart2
	case Chart3:
		return p.Chart3
	case Chart4:
		return p.Chart4
	case Chart5:
		return p.Chart5
	default:
		return p.Chart1
	}
}
func reading(k *cache, o Options, index int) string {
	if index < 0 || index >= len(k.readings) {
		return ""
	}
	return k.readings[index]
}
func describe(k *cache, o Options, index int) string {
	if o.Kind == Pie || o.Kind == Radial {
		if index < 0 || index >= len(k.series) {
			return ""
		}
		series := o.Data.Series[k.series[index].source]
		n, ok := valueAt(series, 0)
		if !ok {
			return locales.ComponentSeriesMissing(series.Label)
		}
		return series.Label + ": " + i18n.FormatNumber(k.locale, n, i18n.NumberOptions{MaxFractionDigits: 100})
	}
	if index < 0 || index >= len(o.Data.Labels) {
		return ""
	}
	var b strings.Builder
	b.WriteString(o.Data.Labels[index])
	for _, ps := range k.series {
		series := o.Data.Series[ps.source]
		b.WriteString(" · ")
		b.WriteString(series.Label)
		b.WriteString(": ")
		if n, ok := valueAt(series, index); ok {
			b.WriteString(i18n.FormatNumber(k.locale, n, i18n.NumberOptions{MaxFractionDigits: 100}))
		} else {
			b.WriteString(locales.ComponentMissing())
		}
	}
	return b.String()
}

// Chart keeps its selection description, tooltip, legend and plotted values on
// one caller data model. Missing samples break paths without generated endpoints.
func Chart(c *ui.Context, s *State, o Options) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	root := ui.Column(c).MinWidth(0).Gap(theme.Space8).Label(o.Label).Description(o.Description)
	k := ui.Local(root, "chart.geometry", func() *cache { return &cache{} })
	(*k).update(s, o, c.Region().FormatLocale(c.Locale()))
	cached := *k
	for i := range cached.series {
		cached.series[i].color = seriesColor(c, o.Data.Series[cached.series[i].source].Color)
	}
	p := theme.Colors(c)
	height := o.Height
	if height <= 0 {
		height = 220
	}
	root.Children(func() {
		axes := native.Axes{}
		if !o.HideAxes && o.Kind <= Area {
			axes.Left = cached.ticks
			axes.Bottom = cached.bottomTicks
			axes.Label = func(_ native.Side, tick native.Tick) {
				ui.Text(c, tick.Label).NoWrap().FontSize(theme.Scale(c, 11)).TextColor(p.MutedForeground)
			}
		}
		native.Plot(c.Key("plot"), axes, func(area ui.Element) {
			area.Label(o.Label).Description(o.Description).Clip()
			if o.Kind <= Area {
				cartesian(c, area, cached, s, o)
			} else {
				polar(c, area, cached, s, o)
			}
			if !o.HideTooltip && s.Selected >= 0 {
				overlay.Tooltip(c, area, overlay.TooltipOptions{Text: reading(cached, o, s.Selected)})
			}
		}).Height(height).FillWidth().MinWidth(0).Gap(theme.Space8)
		if !o.HideLegend {
			ui.Row(c.Key("legend")).Wrap().Gap(theme.Space8).Children(func() {
				for _, series := range o.Data.Series {
					visible := !s.Hidden[series.ID]
					e := mygoui.Toggle(c.Key(series.ID), &visible, mygoui.ToggleOptions{Label: series.Label, Size: mygoui.ToggleSM, OnChange: func(on bool) {
						if s.Hidden == nil {
							s.Hidden = make(map[SeriesID]bool)
						}
						s.Hidden[series.ID] = !on
						s.Selected = -1
						if o.OnVisibility != nil {
							o.OnVisibility(series.ID, on)
						}
					}})
					e.Border(2, seriesColor(c, series.Color))
				}
			})
		}
		// A stable reading line is also the keyboard tooltip and handles narrow plots.
		text := cached.message
		if !o.HideTooltip && s.Selected >= 0 {
			if selected := reading(cached, o, s.Selected); selected != "" {
				text = selected
			}
		}
		ui.Text(c, text).MinHeight(theme.Scale(c, 20)).FontSize(theme.Scale(c, 12)).TextColor(p.MutedForeground).Role(ui.RoleStatus).Selectable()
	})
	return root
}
func cartesian(c *ui.Context, area ui.Element, k *cache, s *State, o Options) {
	p := theme.Colors(c)
	area.Draw(func(painter *ui.Painter, r ui.Rect) {
		if !o.HideAxes {
			for _, tick := range k.ticks {
				painter.Line(r.X, native.Y(r, tick.At), r.X+r.W, native.Y(r, tick.At), 1, p.Border)
			}
		}
		if o.Kind == Bar {
			return
		}
		if !k.geometryReady || r != k.rect {
			for si := range k.series {
				ps := &k.series[si]
				for ri := range ps.runs {
					run := &ps.runs[ri]
					native.LineInto(&run.line, r, k.xs[run.from:run.to], ps.high[run.from:run.to], native.CurveLinear)
					if o.Kind == Area {
						run.area = *native.Area(r, k.xs[run.from:run.to], ps.low[run.from:run.to], ps.high[run.from:run.to], native.CurveLinear)
					}
				}
			}
			k.rect = r
			k.geometryReady = true
		}
		for si := range k.series {
			ps := &k.series[si]
			for ri := range ps.runs {
				run := &ps.runs[ri]
				if o.Kind == Area {
					painter.FillPath(&run.area, ps.color.Alpha(0.2))
				}
				painter.StrokePath(&run.line, 2, ps.color)
			}
		}
	})
	if o.Kind == Bar {
		for si, ps := range k.series {
			series := o.Data.Series[ps.source]
			for i := range k.xs {
				n, ok := valueAt(series, i)
				if !ok {
					continue
				}
				x := k.x.Band(i)
				if k.segments == nil {
					width := (x.Hi - x.Lo) / float32(max(1, len(k.series)))
					x.Lo += float32(si) * width
					x.Hi = x.Lo + width
				}
				span := k.y.Span(0, n)
				if k.segments != nil {
					seg := k.segments[si][i]
					span = k.y.Span(seg.Lo, seg.Hi)
				}
				native.Rect(c.Key([2]int{ps.source, i}), x, span).Radius(theme.Radii().SM).Background(seriesColor(c, series.Color)).Label(locales.ComponentChartValue(o.Data.Labels[i], series.Label, i18n.FormatNumber(k.locale, n, i18n.NumberOptions{MaxFractionDigits: 100})))
			}
		}
	} else {
		for _, ps := range k.series {
			series := o.Data.Series[ps.source]
			for i := range k.xs {
				if _, ok := valueAt(series, i); ok {
					native.Point(c.Key([2]int{ps.source, i}), k.xs[i], ps.high[i], 5, 5).Radius(3).Background(seriesColor(c, series.Color))
				}
			}
		}
	}
	if native.Select(c, area, &s.Selected, k.xs, func(i int) string { return reading(k, o, i) }) && o.OnSelection != nil {
		o.OnSelection(s.Selected)
	}
	if s.Selected >= 0 && s.Selected < len(k.xs) {
		native.RuleX(c, k.xs[s.Selected], 1).Background(p.Ring)
	}
}

type polarPointer struct {
	x, y float32
	over bool
}

func polar(c *ui.Context, area ui.Element, k *cache, s *State, o Options) {
	xs := k.xs
	if o.Kind == Pie || o.Kind == Radial {
		xs = k.polarXS
	}
	pointer := ui.Local(area, "polar.pointer", func() polarPointer { return polarPointer{} })
	old := s.Selected
	native.Select(c, area, &s.Selected, xs, func(i int) string { return reading(k, o, i) })
	// A resting pointer must not override native keyboard selection.
	if x, y, over := area.PointerPosition(); over {
		moved := !pointer.over || x != pointer.x || y != pointer.y
		pointer.x, pointer.y, pointer.over = x, y, true
		if moved {
			r := area.Bounds()
			cx, cy, radius := r.W/2, r.H/2, min(r.W, r.H)/2-12
			switch o.Kind {
			case Pie:
				if i, ok := native.SliceAt(k.slices, x-cx, y-cy, 0, radius); ok && k.polarValues[i] > 0 {
					s.Selected = i
				} else {
					s.Selected = -1
				}
			case Radial:
				distance := float32(math.Hypot(float64(x-cx), float64(y-cy)))
				s.Selected = -1
				n := len(k.series)
				if n > 0 {
					band := radius / float32(n+1)
					i := int((radius - distance) / band)
					if i >= 0 && i < n && distance >= radius-float32(i+1)*band+2 && distance <= radius-float32(i)*band-2 {
						a := math.Atan2(float64(x-cx), -float64(y-cy))
						if a < 0 {
							a += 2 * math.Pi
						}
						maximum := slices.Max(k.polarValues)
						if maximum > 0 && a <= 2*math.Pi*k.polarValues[i]/maximum {
							s.Selected = i
						}
					}
				}
			case Radar:
				if len(xs) > 0 {
					a := math.Atan2(float64(x-cx), -float64(y-cy))
					if a < 0 {
						a += 2 * math.Pi
					}
					s.Selected = int(math.Round(a/(2*math.Pi)*float64(len(xs)))) % len(xs)
				}
			}
		}
	} else {
		pointer.over = false
	}
	area.Value(reading(k, o, s.Selected))
	if o.Kind == Radar && !o.HideAxes && len(k.xs) >= 3 {
		for i, label := range o.Data.Labels {
			a := 2 * math.Pi * float64(i) / float64(len(k.xs))
			x, y := native.Polar(0.5, 0.5, 0.42, a)
			native.Anchor(c, x, 1-y, ui.Center, ui.Center, func() {
				ui.Text(c, label).NoWrap().FontSize(theme.Scale(c, 11)).TextColor(theme.Colors(c).MutedForeground)
			})
		}
	}
	if old != s.Selected && o.OnSelection != nil {
		o.OnSelection(s.Selected)
	}
	colors := theme.Colors(c)
	area.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, radius := r.X+r.W/2, r.Y+r.H/2, max(float32(0), min(r.W, r.H)/2-12)
		if o.Kind == Radar {
			radarPaint(p, r, k, s, o, cx, cy, radius, colors.Border)
			return
		}
		if !k.geometryReady || k.rect != r {
			if o.Kind == Pie {
				for i, slice := range k.slices {
					k.polarPaths[i] = *native.Sector(cx, cy, 0, radius, slice)
				}
			} else if len(k.series) > 0 {
				maximum := slices.Max(k.polarValues)
				band := radius / float32(len(k.series)+1)
				for i, value := range k.polarValues {
					outer := radius - float32(i)*band
					inner := outer - band + 4
					k.tracks[i] = *native.Sector(cx, cy, inner, outer, native.Slice{Start: 0, End: 2 * math.Pi})
					k.polarPaths[i].Reset()
					if maximum > 0 && value > 0 {
						k.polarPaths[i] = *native.Sector(cx, cy, inner, outer, native.Slice{Start: 0, End: 2 * math.Pi * value / maximum})
					}
				}
			}
			k.geometryReady = true
			k.rect = r
		}
		for i, value := range k.polarValues {
			if o.Kind == Radial {
				p.FillPath(&k.tracks[i], colors.Muted)
			}
			if value <= 0 {
				continue
			}
			p.FillPath(&k.polarPaths[i], k.series[i].color)
			if i == s.Selected {
				p.StrokePath(&k.polarPaths[i], 2, colors.Ring)
			}
		}
	})
}
func radarPaint(p *ui.Painter, r ui.Rect, k *cache, s *State, o Options, cx, cy, radius float32, border ui.Color) {
	n := len(k.xs)
	if n < 3 {
		return
	}
	maximum := 0.0
	for _, series := range k.values {
		for _, v := range series {
			maximum = max(maximum, v)
		}
	}
	if maximum <= 0 {
		maximum = 1
	}
	if !k.geometryReady || r != k.rect {
		k.grid.Reset()
		for level := 1; level <= 4; level++ {
			for i := range n {
				x, y := native.Polar(cx, cy, radius*float32(level)/4, 2*math.Pi*float64(i)/float64(n))
				if i == 0 {
					k.grid.MoveTo(x, y)
				} else {
					k.grid.LineTo(x, y)
				}
			}
			k.grid.Close()
		}
		for i := range n {
			x, y := native.Polar(cx, cy, radius, 2*math.Pi*float64(i)/float64(n))
			k.grid.MoveTo(cx, cy).LineTo(x, y)
		}
		for si := range k.series {
			ps := &k.series[si]
			ps.radar.Reset()
			ps.fill = true
			penDown := false
			for i := range n {
				v, ok := valueAt(o.Data.Series[ps.source], i)
				if !ok || v < 0 {
					ps.fill = false
					penDown = false
					continue
				}
				x, y := native.Polar(cx, cy, radius*float32(v/maximum), 2*math.Pi*float64(i)/float64(n))
				if !penDown {
					ps.radar.MoveTo(x, y)
				} else {
					ps.radar.LineTo(x, y)
				}
				penDown = true
			}
			if ps.fill {
				ps.radar.Close()
			}
		}
		k.rect = r
		k.geometryReady = true
	}
	p.StrokePath(&k.grid, 1, border)
	for si := range k.series {
		ps := &k.series[si]
		if ps.fill {
			p.FillPath(&ps.radar, ps.color.Alpha(0.15))
		}
		p.StrokePath(&ps.radar, 2, ps.color)
		for i := range n {
			v, ok := valueAt(o.Data.Series[ps.source], i)
			if !ok || v < 0 {
				continue
			}
			x, y := native.Polar(cx, cy, radius*float32(v/maximum), 2*math.Pi*float64(i)/float64(n))
			size := float32(4)
			if i == s.Selected {
				size = 8
			}
			p.Fill(ui.Rect{X: x - size/2, Y: y - size/2, W: size, H: size}, ps.color, size/2)
		}
	}
}
