package gallery

import (
	"fmt"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/chart"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const bklitSourceRoot = "https://github.com/bklit/bklit-ui/blob/ec5a9705ee973f95fd0ef2d974e9d00ae3fda011/apps/web/content/docs/"

var chartSubjects = []struct{ id, title, slug string }{
	{"areachart", "Area Chart", "area-chart"},
	{"barchart", "Bar Chart", "bar-chart"},
	{"candlestickchart", "Candlestick Chart", "candlestick-chart"},
	{"choroplethchart", "Choropleth Chart", "choropleth-chart"},
	{"composedchart", "Composed Chart", "composed-chart"},
	{"funnelchart", "Funnel Chart", "funnel-chart"},
	{"gaugechart", "Gauge", "gauge-chart"},
	{"heatmapchart", "Heatmap Chart", "heatmap-chart"},
	{"linechart", "Line Chart", "line-chart"},
	{"profitlossline", "Profit/Loss Line", "profit-loss-line"},
	{"livelinechart", "Live Line Chart", "live-line-chart"},
	{"piechart", "Pie Chart", "pie-chart"},
	{"radarchart", "Radar Chart", "radar-chart"},
	{"radialchart", "Radial Chart", "radial-chart"},
	{"ringchart", "Ring Chart", "ring-chart"},
	{"scatterchart", "Scatter Chart", "scatter-chart"},
	{"sankeychart", "Sankey Chart", "sankey-chart"},
	{"sparkline", "Sparkline", "sparkline"},
	{"sunburstchart", "Sunburst Chart", "sunburst-chart"},
}

var chartUtilities = []struct{ id, title, slug string }{
	{"chartlegend", "Chart Legend", "legend"},
	{"chartgrid", "Chart Grid", "grid"},
	{"chartbackground", "Chart Background", "background"},
	{"referencearea", "Reference Area", "reference-area"},
	{"projectionline", "Projection Line", "projection-line"},
	{"chartband", "Band", "band"},
	{"chartweekends", "Weekends", "weekends"},
	{"chartdeadline", "Deadline", "deadline"},
	{"trendline", "Trend Line", "trend-line"},
	{"chartmarkers", "Chart Markers", "markers"},
	{"charttooltip", "Chart Tooltip", "tooltip"},
	{"chartbrush", "Chart Brush", "brush"},
	{"chartxaxis", "X Axis", "axis/x-axis"},
	{"chartyaxis", "Y Axis", "axis/y-axis"},
	{"chartcontext", "Chart Context", "use-chart"},
	{"customindicator", "Custom Indicator", "custom-indicator"},
}

func chartSource(id string) string {
	for _, subject := range chartSubjects {
		if id == subject.id {
			return bklitSourceRoot + "components/" + subject.slug + ".mdx"
		}
	}
	for _, subject := range chartUtilities {
		if id == subject.id {
			return bklitSourceRoot + "utility/" + subject.slug + ".mdx"
		}
	}
	return ""
}

type chartDemoState struct {
	Selection chart.State
	Drill     chart.SunburstState
	Map       chart.MapState
	Brush     chart.BrushState
	Offset    int
}

func chartStateOf(s *State, id string) *chartDemoState {
	st := StateOf[*chartDemoState](s, "chart/"+id)
	if *st == nil {
		*st = &chartDemoState{
			Selection: chart.State{Selected: -1},
			Drill:     chart.SunburstState{State: chart.State{Selected: -1}},
			Brush:     chart.BrushState{From: 0, To: 5},
		}
	}
	return *st
}

func init() {
	for i, subject := range chartSubjects {
		id := subject.id
		RegisterComponent(ComponentPage{
			ID: id, Title: subject.title, Group: "Charts", Order: 1100 + i*10,
			Description: "Native data visualization with coss chart tokens and caller-owned inspection.",
			Examples:    []ComponentExample{{Title: "Sample data", Demo: func(c *ui.Context, s *State) { chartDemo(c, s, id) }}},
			Code:        chartSnippet(id), Import: "github.com/Cyenoch/mygoui/chart",
			Params: chartParams(id),
			Tokens: []string{"Chart1", "Chart2", "Chart3", "Chart4", "Chart5", "Border", "MutedForeground", "Ring"},
		})
	}
	for i, subject := range chartUtilities {
		id := subject.id
		RegisterComponent(ComponentPage{
			ID: id, Title: subject.title, Group: "Chart Utilities", Order: 1400 + i*10,
			Description: "Composable native chart geometry, inspection and controls.",
			Examples:    []ComponentExample{{Title: "Sample data", Demo: func(c *ui.Context, s *State) { chartDemo(c, s, id) }}},
			Code:        chartSnippet(id), Import: "github.com/Cyenoch/mygoui/chart",
			Params: chartParams(id),
			Tokens: []string{"Chart1", "Border", "MutedForeground", "Ring", "Popover", "PopoverForeground"},
		})
	}
}

func demoSeries(offset int) []chart.Series {
	points := []chart.Point{}
	for i, v := range []float64{24, 42, 30, 65, 48, 80} {
		points = append(points, chart.Point{X: float64(i + offset), Y: v, Label: fmt.Sprintf("Day %d", i+offset+1)})
	}
	return []chart.Series{{ID: "sales", Label: "Sales", Points: points}}
}

func demoValues() []chart.Value {
	return []chart.Value{{ID: "direct", Label: "Direct", Value: 72}, {ID: "search", Label: "Search", Value: 48}, {ID: "social", Label: "Social", Value: 24}}
}

func demoHierarchy() *chart.HierarchyNode {
	return &chart.HierarchyNode{ID: "all", Label: "All", Children: []*chart.HierarchyNode{
		{ID: "product", Label: "Product", Children: []*chart.HierarchyNode{{ID: "web", Label: "Web", Value: 40}, {ID: "native", Label: "Native", Value: 30}}},
		{ID: "services", Label: "Services", Children: []*chart.HierarchyNode{{ID: "support", Label: "Support", Value: 20}, {ID: "other", Label: "Other", Value: 10}}},
	}}
}

func demoRegions() []chart.Region {
	regions := []chart.Region{}
	for i, item := range []struct {
		id       string
		lon, lat float64
	}{{"West", -8, 40}, {"Central", 0, 42}, {"East", 8, 40}} {
		regions = append(regions, chart.Region{ID: item.id, Label: item.id, Value: float64(30 + i*25), Polygons: []chart.Polygon{{Rings: [][]chart.GeoPoint{{
			{Longitude: item.lon, Latitude: item.lat}, {Longitude: item.lon + 7, Latitude: item.lat},
			{Longitude: item.lon + 7, Latitude: item.lat + 7}, {Longitude: item.lon, Latitude: item.lat + 7},
		}}}}})
	}
	return regions
}

func chartDemo(c *ui.Context, s *State, id string) {
	st := chartStateOf(s, id)
	o := chart.Options{Label: "Chart specimen", State: &st.Selection, Height: 220}
	series := demoSeries(st.Offset)
	cart := chart.CartesianOptions{Options: o, Series: series}
	var err error
	switch id {
	case "areachart":
		_, err = chart.Area(c, cart)
	case "barchart":
		_, err = chart.Bar(c, cart)
	case "linechart":
		_, err = chart.Line(c, cart)
	case "scatterchart":
		_, err = chart.Scatter(c, cart)
	case "profitlossline":
		for i := range series[0].Points {
			series[0].Points[i].Y -= 45
		}
		_, err = chart.ProfitLoss(c, cart)
	case "livelinechart":
		mygoui.Button(c, mygoui.ButtonOptions{Label: "Next sample", Variant: mygoui.ButtonOutline, OnClick: func() { st.Offset++ }})
		_, err = chart.LiveLine(c, cart, 4)
	case "composedchart":
		series[0].Kind = chart.BarSeries
		second := demoSeries(st.Offset)[0]
		second.ID, second.Label, second.Kind = "target", "Target", chart.LineSeries
		for i := range second.Points {
			second.Points[i].Y += 12
		}
		cart.Series = append(series, second)
		_, err = chart.Composed(c, cart)
	case "sparkline":
		for _, kind := range []chart.SeriesKind{chart.LineSeries, chart.AreaSeries, chart.BarSeries} {
			cart.State = &chartStateOf(s, fmt.Sprintf("sparkline/%d", kind)).Selection
			cart.Series[0].Kind = kind
			cart.Height = 48
			_, err = chart.Sparkline(c.Key(kind), cart)
			if err != nil {
				break
			}
		}
	case "candlestickchart":
		candles := []chart.Candle{}
		for i := range 8 {
			open := float64(40 + i*3)
			close := open + 6
			if i%3 == 0 {
				close = open - 4
			}
			candles = append(candles, chart.Candle{X: float64(i), Open: open, Close: close, Low: min(open, close) - 4, High: max(open, close) + 6, Label: fmt.Sprint("Day ", i+1)})
		}
		_, err = chart.Candlestick(c, chart.CandlestickOptions{Options: o, Candles: candles})
	case "piechart":
		items := []chart.LegendItem{}
		for _, v := range demoValues() {
			items = append(items, chart.LegendItem{ID: v.ID, Label: v.Label})
		}
		chart.Legend(c, &st.Selection, items, false)
		_, err = chart.Pie(c, chart.PieOptions{Options: o, Values: demoValues(), InnerRadius: .65})
	case "ringchart":
		_, err = chart.Ring(c, chart.RingOptions{Options: o, Values: demoValues()})
	case "funnelchart":
		_, err = chart.Funnel(c, chart.FunnelOptions{Options: o, Values: demoValues()})
	case "gaugechart":
		_, err = chart.Gauge(c, chart.GaugeOptions{Options: o, Value: 72})
	case "radarchart":
		chart.Legend(c, &st.Selection, []chart.LegendItem{{ID: "current", Label: "Current"}, {ID: "prior", Label: "Prior"}}, false)
		_, err = chart.Radar(c, chart.RadarOptions{Options: o, Metrics: []string{"Speed", "Quality", "Reach", "Value", "Trust"}, Max: 100,
			Series: []chart.RadarSeries{{ID: "current", Label: "Current", Values: []float64{80, 95, 65, 85, 90}}, {ID: "prior", Label: "Prior", Values: []float64{60, 70, 85, 65, 75}}}})
	case "radialchart":
		_, err = chart.Radial(c, chart.RadialOptions{Options: o, Values: demoValues()})
	case "heatmapchart":
		ui.Row(c).FillWidth().Wrap().Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
			for i, label := range []string{"0", "1–4", "5–8", "9–12", "13–18"} {
				ui.Row(c).Gap(theme.Space4).AlignItems(ui.Center).Children(func() {
					fill := theme.Colors(c).MutedBackground
					if i > 0 {
						fill = theme.Colors(c).MutedBackground.Mix(chart.SeriesColor(c, 0), float32(i)/4)
					}
					ui.Box(c).Size(10, 10).Radius(2).Background(fill)
					mygoui.Typography(c, label, mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
				})
			}
		})
		cells := []chart.Cell{}
		for column := range 12 {
			for row := range 7 {
				cells = append(cells, chart.Cell{Column: column, Row: row, Value: float64((column*3 + row*5) % 19)})
			}
		}
		_, err = chart.Heatmap(c, chart.HeatmapOptions{Options: o, Thresholds: []float64{0, 4, 8, 12}, Columns: []string{"Jan", "", "", "Apr", "", "", "Jul", "", "", "Oct", "", ""}, Rows: []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, Cells: cells})
	case "choroplethchart":
		_, err = chart.Choropleth(c, chart.ChoroplethOptions{Options: o, Regions: demoRegions(), Viewport: &st.Map})
	case "sankeychart":
		_, err = chart.Sankey(c, chart.SankeyOptions{Options: o, Nodes: []chart.FlowNode{{ID: "visits", Label: "Visits"}, {ID: "trial", Label: "Trial"}, {ID: "paid", Label: "Paid"}, {ID: "exit", Label: "Exit"}},
			Links: []chart.FlowLink{{Source: "visits", Target: "trial", Value: 60}, {Source: "visits", Target: "exit", Value: 40}, {Source: "trial", Target: "paid", Value: 45}, {Source: "trial", Target: "exit", Value: 15}}})
	case "sunburstchart":
		_, err = chart.Sunburst(c, chart.SunburstOptions{Options: o, Root: demoHierarchy(), Drill: &st.Drill})
	default:
		err = chartUtilityDemo(c, st, id, cart)
	}
	if err != nil {
		mygoui.Typography(c, err.Error(), mygoui.TypographyOptions{Type: theme.TypeBodySM})
	}
}

func chartUtilityDemo(c *ui.Context, st *chartDemoState, id string, cart chart.CartesianOptions) error {
	points := cart.Series[0].Points
	switch id {
	case "chartlegend":
		chart.Legend(c, &st.Selection, []chart.LegendItem{{ID: "sales", Label: "Sales"}}, false)
	case "chartbrush":
		if _, err := chart.Brush(c, &st.Brush, 0, 5, false); err != nil {
			return err
		}
		cart.Brush = &st.Brush
	case "charttooltip":
		st.Selection.Selected = 2
	case "referencearea":
		cart.Annotations = []chart.Annotation{{Kind: chart.ReferenceArea, X1: 1, X2: 3, Y1: 20, Y2: 60}}
	case "projectionline":
		cart.XDomain = &chart.Domain{Min: 0, Max: 7}
		cart.Annotations = []chart.Annotation{{Kind: chart.ProjectionLine, X1: 5, Y1: 80, X2: 7, Y2: 90}}
	case "chartband":
		upper := append([]chart.Point{}, points...)
		for i := range upper {
			upper[i].Y += 12
		}
		cart.Annotations = []chart.Annotation{{Kind: chart.Band, Points: points, Upper: upper}}
	case "chartweekends":
		for i := range points {
			points[i].X = 1792790400 + float64(i)*86400
		}
		cart.XAxis.Time = true
		cart.Annotations = []chart.Annotation{{Kind: chart.Weekends}}
	case "chartdeadline":
		cart.Annotations = []chart.Annotation{{Kind: chart.Deadline, X1: 4, Y1: 60, Label: "Due"}}
	case "trendline":
		cart.Annotations = []chart.Annotation{{Kind: chart.TrendLine, Points: points}}
	case "chartmarkers":
		cart.Annotations = []chart.Annotation{{Kind: chart.Marker, X1: 3, Y1: 65, Label: "Release"}}
	case "chartxaxis":
		cart.XAxis.Format = func(v float64) string { return fmt.Sprintf("Day %.0f", v+1) }
	case "chartyaxis":
		cart.YAxis.Format = func(v float64) string { return fmt.Sprintf("$%.0f", v) }
	case "chartcontext":
		cart.Build = func(c *ui.Context, ctx chart.PlotContext) {
			ui.Box(c).Absolute().LeftPercent(ctx.X.At(3)*100).TopPercent((1-ctx.Y.At(65))*100).Size(8, 8).Background(theme.Colors(c).Chart2).Radius(4)
		}
	case "customindicator":
		st.Selection.Selected = 3
		cart.Indicator = func(p *ui.Painter, r ui.Rect, i int) {
			x := r.X + float32(i)/5*r.W
			p.Fill(ui.Rect{X: x - 4, Y: r.Y, W: 8, H: r.H}, theme.Colors(c).Chart2.Alpha(.2), 2)
		}
	case "chartbackground":
		cart.Background = chart.BackgroundDots
	}
	_, err := chart.Line(c, cart)
	return err
}

func chartParams(id string) []ComponentParam {
	params := []ComponentParam{
		{"Options", "chart.Options", "zero", "Label, description, height, caller state, loading, disabled and callbacks."},
		{"State", "*chart.State", "nil", "Window-owned selection, hidden series and native focus handle."},
	}
	name, typ := "Series", "[]chart.Series"
	switch id {
	case "candlestickchart":
		name, typ = "Candles", "[]chart.Candle"
	case "choroplethchart":
		name, typ = "Regions", "[]chart.Region"
	case "heatmapchart":
		name, typ = "Cells", "[]chart.Cell"
	case "piechart", "ringchart", "radialchart", "funnelchart":
		name, typ = "Values", "[]chart.Value"
	case "radarchart":
		name, typ = "Series", "[]chart.RadarSeries"
	case "sankeychart":
		name, typ = "Nodes / Links", "[]chart.FlowNode / []chart.FlowLink"
	case "sunburstchart":
		name, typ = "Root / Drill", "*chart.HierarchyNode / *chart.SunburstState"
	case "gaugechart":
		name, typ = "Value / Min / Max", "float64"
	}
	return append(params, ComponentParam{name, typ, "empty", "Caller-owned data; constructors return contextual validation errors."})
}

func chartSnippet(id string) string {
	const prefix = `state := chart.State{Selected: -1} // retain per window
opts := chart.Options{Label: "Measurements", State: &state}
`
	const values = `values := []chart.Value{
    {ID: "direct", Label: "Direct", Value: 72},
    {ID: "search", Label: "Search", Value: 48},
}
`
	const series = `series := []chart.Series{{ID: "sales", Label: "Sales", Points: []chart.Point{
    {X: 0, Y: 24}, {X: 1, Y: 42}, {X: 2, Y: 30},
}}}
plot := chart.CartesianOptions{Options: opts, Series: series}
`
	const finish = `
if err != nil { ui.Text(c, err.Error()) }`
	functions := map[string]string{"areachart": "Area", "barchart": "Bar", "linechart": "Line", "scatterchart": "Scatter", "profitlossline": "ProfitLoss", "composedchart": "Composed", "sparkline": "Sparkline"}
	if fn := functions[id]; fn != "" {
		return prefix + series + "_, err := chart." + fn + "(c, plot)" + finish
	}
	switch id {
	case "livelinechart":
		return prefix + series + "_, err := chart.LiveLine(c, plot, 120)" + finish
	case "piechart":
		return prefix + values + "_, err := chart.Pie(c, chart.PieOptions{Options: opts, Values: values, InnerRadius: .65})" + finish
	case "ringchart":
		return prefix + values + "_, err := chart.Ring(c, chart.RingOptions{Options: opts, Values: values})" + finish
	case "funnelchart":
		return prefix + values + "_, err := chart.Funnel(c, chart.FunnelOptions{Options: opts, Values: values})" + finish
	case "radialchart":
		return prefix + values + "_, err := chart.Radial(c, chart.RadialOptions{Options: opts, Values: values, Kind: chart.RadialDots})" + finish
	case "gaugechart":
		return prefix + "_, err := chart.Gauge(c, chart.GaugeOptions{Options: opts, Value: 72})" + finish
	case "radarchart":
		return prefix + `_, err := chart.Radar(c, chart.RadarOptions{Options: opts,
    Metrics: []string{"Speed", "Quality", "Reach"}, Max: 100,
    Series: []chart.RadarSeries{{ID: "current", Label: "Current", Values: []float64{80, 95, 65}}},
})` + finish
	case "candlestickchart":
		return prefix + `_, err := chart.Candlestick(c, chart.CandlestickOptions{Options: opts,
    Candles: []chart.Candle{{X: 1, Open: 40, High: 52, Low: 38, Close: 48, Label: "Today"}},
})` + finish
	case "heatmapchart":
		return prefix + `_, err := chart.Heatmap(c, chart.HeatmapOptions{Options: opts,
    Columns: []string{"Week 1", "Week 2"}, Rows: []string{"Mon", "Tue"},
    Cells: []chart.Cell{{Column: 0, Row: 0, Value: 4}, {Column: 1, Row: 1, Value: 8}},
})` + finish
	case "choroplethchart":
		return prefix + `viewport := chart.MapState{} // retain per window
_, err := chart.Choropleth(c, chart.ChoroplethOptions{Options: opts, Viewport: &viewport,
    Regions: []chart.Region{{ID: "west", Label: "West", Value: 42,
        Polygons: []chart.Polygon{{Rings: [][]chart.GeoPoint{{
            {Longitude: 0, Latitude: 40}, {Longitude: 8, Latitude: 40},
            {Longitude: 4, Latitude: 48},
        }}}},
    }},
})` + finish
	case "sankeychart":
		return prefix + `_, err := chart.Sankey(c, chart.SankeyOptions{Options: opts,
    Nodes: []chart.FlowNode{{ID: "visit", Label: "Visit"}, {ID: "paid", Label: "Paid"}},
    Links: []chart.FlowLink{{Source: "visit", Target: "paid", Value: 42}},
})` + finish
	case "sunburstchart":
		return prefix + `drill := chart.SunburstState{State: chart.State{Selected: -1}} // retain per window
_, err := chart.Sunburst(c, chart.SunburstOptions{Options: opts, Drill: &drill,
    Root: &chart.HierarchyNode{ID: "all", Label: "All", Children: []*chart.HierarchyNode{
        {ID: "web", Label: "Web", Value: 42}, {ID: "native", Label: "Native", Value: 58},
    }},
})` + finish
	}
	extra := ""
	switch id {
	case "chartlegend":
		extra = `chart.Legend(c, &state, []chart.LegendItem{{ID: "sales", Label: "Sales"}}, false)`
	case "chartbrush":
		extra = `brush := chart.BrushState{From: 0, To: 2} // retain per window
if _, err := chart.Brush(c, &brush, 0, 2, false); err != nil { ui.Text(c, err.Error()) }
plot.Brush = &brush`
	case "charttooltip":
		extra = `state.Selected = 1`
	case "chartxaxis":
		extra = `plot.XAxis = chart.Axis{Ticks: 3, Format: func(v float64) string { return fmt.Sprintf("Day %.0f", v+1) }}`
	case "chartyaxis":
		extra = `plot.YAxis = chart.Axis{Ticks: 4, Format: func(v float64) string { return fmt.Sprintf("$%.0f", v) }}`
	case "chartbackground":
		extra = `plot.Background = chart.BackgroundDots`
	case "chartgrid":
		extra = `plot.HideGrid = false`
	case "chartcontext":
		extra = `plot.Build = func(c *ui.Context, ctx chart.PlotContext) {
    ui.Box(c).Absolute().LeftPercent(ctx.X.At(1)*100).TopPercent((1-ctx.Y.At(42))*100).
        Size(8, 8).Radius(4).Background(chart.SeriesColor(c, 1))
}`
	case "customindicator":
		extra = `plot.Indicator = func(p *ui.Painter, r ui.Rect, selected int) {
    p.Fill(ui.Rect{X: r.X+float32(selected)*r.W/2, Y: r.Y, W: 3, H: r.H}, chart.SeriesColor(c, 1), 0)
}`
	default:
		kind := map[string]string{"referencearea": "ReferenceArea", "projectionline": "ProjectionLine", "chartband": "Band", "chartweekends": "Weekends", "chartdeadline": "Deadline", "trendline": "TrendLine", "chartmarkers": "Marker"}[id]
		extra = "plot.Annotations = []chart.Annotation{{Kind: chart." + kind + ", X1: 1, X2: 2, Y1: 10, Y2: 40}}"
		if id == "chartband" {
			extra = `plot.Annotations = []chart.Annotation{{Kind: chart.Band,
    Points: []chart.Point{{X: 0, Y: 10}, {X: 2, Y: 20}},
    Upper: []chart.Point{{X: 0, Y: 30}, {X: 2, Y: 45}},
}}`
		}
		if id == "trendline" {
			extra = `plot.Annotations = []chart.Annotation{{Kind: chart.TrendLine, Points: series[0].Points}}`
		}
	}
	return prefix + series + extra + "\n_, err := chart.Line(c, plot)" + finish
}
