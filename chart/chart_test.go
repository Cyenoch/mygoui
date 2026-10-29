package chart

import (
	"fmt"
	"image/color"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func samples() []Series {
	return []Series{{ID: "a", Label: "Sales", Points: []Point{{X: 0, Y: -4}, {X: 1, Y: 8}, {X: 2, Y: 12}}}}
}

func TestInspectionKeyboardDisabledAndDataReplacement(t *testing.T) {
	state := State{Selected: -1}
	data := samples()
	disabled, calls := false, 0
	var buildErr error
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		_, buildErr = Line(c, CartesianOptions{Options: Options{Label: "Revenue", State: &state, Disabled: disabled, OnSelect: func(int) { calls++ }}, Series: data})
	}, 600, 400)
	tt.Frame()
	if buildErr != nil {
		t.Fatal(buildErr)
	}
	if err := tt.Click("Revenue"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	if state.Selected != 0 || !tt.HasText("0 · Sales: -4") {
		t.Fatalf("first measurement not inspected: %d, %v", state.Selected, tt.Texts())
	}
	tt.Key(0, ui.KeyEnd)
	if state.Selected != 2 || !tt.HasText("2 · Sales: 12") {
		t.Fatal("End did not select last measurement")
	}
	tt.Key(0, ui.KeyEscape)
	if state.Selected != -1 {
		t.Fatal("Escape did not clear selection")
	}
	tt.Key(0, ui.KeyHome)
	disabled = true
	tt.Frame()
	before := calls
	tt.Key(0, ui.KeyRight)
	if state.Selected != 0 || calls != before {
		t.Fatal("disabled chart handled selection")
	}
	disabled = false
	state.Selected = 2
	data[0].Points = data[0].Points[:1]
	tt.Frame()
	if state.Selected != -1 {
		t.Fatal("data replacement left an out-of-range selection")
	}
	data = nil
	tt.Frame()
	if !tt.HasText("No data") {
		t.Fatal("empty data fabricated a chart")
	}
}

func TestLegendAndBrushDriveVisibleData(t *testing.T) {
	state := State{Selected: -1}
	brush := BrushState{From: 0, To: 2}
	var buildErr error
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).FillWidth().Children(func() {
			Legend(c, &state, []LegendItem{{ID: "a", Label: "Sales"}}, false)
			_, buildErr = Brush(c, &brush, 0, 2, false)
			if buildErr != nil {
				return
			}
			_, buildErr = Line(c, CartesianOptions{Options: Options{Label: "Filtered", State: &state, Height: 160}, Series: samples(), Brush: &brush})
		})
	}, 600, 450)
	tt.Frame()
	if buildErr != nil {
		t.Fatal(buildErr)
	}
	if err := tt.Click("Sales"); err != nil {
		t.Fatal(err)
	}
	if !state.Hidden["a"] || !tt.HasText("No data") {
		t.Fatal("legend did not hide series")
	}
	if err := tt.Click("Sales"); err != nil {
		t.Fatal(err)
	}
	if state.Hidden["a"] {
		t.Fatal("legend did not restore series")
	}
	brush.From, brush.To = 1, 2
	tt.Frame()
	if err := tt.Click("Filtered"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("1 · Sales: 8") {
		t.Fatalf("brush did not filter X domain: %v", tt.Texts())
	}
}

func TestConstructorsRejectMalformedData(t *testing.T) {
	root := &HierarchyNode{ID: "cycle"}
	root.Children = []*HierarchyNode{root}
	cases := []struct {
		name    string
		build   func(*ui.Context) (ui.Element, error)
		message string
	}{
		{"infinite", func(c *ui.Context) (ui.Element, error) {
			s := samples()
			s[0].Points[0].Y = math.Inf(1)
			return Line(c, CartesianOptions{Series: s})
		}, "finite"},
		{"unordered", func(c *ui.Context) (ui.Element, error) {
			s := samples()
			s[0].Points[1].X = 0
			return Line(c, CartesianOptions{Series: s})
		}, "increasing"},
		{"domain", func(c *ui.Context) (ui.Element, error) {
			return Line(c, CartesianOptions{Series: samples(), YDomain: &Domain{Min: 2, Max: 1}})
		}, "domain"},
		{"ohlc", func(c *ui.Context) (ui.Element, error) {
			return Candlestick(c, CandlestickOptions{Candles: []Candle{{Open: 4, Close: 5, High: 3, Low: 2}}})
		}, "OHLC"},
		{"pie", func(c *ui.Context) (ui.Element, error) {
			return Pie(c, PieOptions{Values: []Value{{ID: "negative", Value: -1}}})
		}, "nonnegative"},
		{"ring", func(c *ui.Context) (ui.Element, error) {
			return Ring(c, RingOptions{Values: []Value{{ID: "over", Value: 101}}})
		}, "range"},
		{"gauge", func(c *ui.Context) (ui.Element, error) { return Gauge(c, GaugeOptions{Value: math.NaN()}) }, "domain"},
		{"heatmap", func(c *ui.Context) (ui.Element, error) {
			return Heatmap(c, HeatmapOptions{Cells: []Cell{{Column: 0, Row: 0, Value: 2}}})
		}, "out-of-range"},
		{"map", func(c *ui.Context) (ui.Element, error) {
			return Choropleth(c, ChoroplethOptions{Regions: []Region{{ID: "bad", Polygons: []Polygon{{Rings: [][]GeoPoint{{{Latitude: 90}, {}, {Longitude: 2}}}}}}}})
		}, "coordinates"},
		{"flow", func(c *ui.Context) (ui.Element, error) {
			return Sankey(c, SankeyOptions{Nodes: []FlowNode{{ID: "a"}, {ID: "b"}}, Links: []FlowLink{{Source: "a", Target: "b", Value: 1}, {Source: "b", Target: "a", Value: 1}}})
		}, "cycle"},
		{"hierarchy", func(c *ui.Context) (ui.Element, error) { return Sunburst(c, SunburstOptions{Root: root}) }, "acyclic"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			tt := ui.NewTester(func(c *ui.Context) { _, err = tc.build(c) }, 400, 300)
			tt.Frame()
			if err == nil || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("got %v, want %q", err, tc.message)
			}
		})
	}
}

func TestSignedStackAndCallerDataRemainAccurate(t *testing.T) {
	series := []Series{
		{ID: "a", Kind: BarSeries, Points: []Point{{X: 0, Y: -4}, {X: 1, Y: 8}}},
		{ID: "b", Kind: BarSeries, Points: []Point{{X: 0, Y: -6}, {X: 1, Y: 12}}},
	}
	_, y, stacks, err := cartesianScales(CartesianOptions{Stacked: true}, series)
	if err != nil {
		t.Fatal(err)
	}
	if y.Lo != -10 || y.Hi != 20 || stacks[1][1].Hi != 20 {
		t.Fatalf("signed stack domain = %+v, stack = %+v", y, stacks)
	}
	data := samples()
	before := append([]Point{}, data[0].Points...)
	state := State{Selected: -1}
	tt := ui.NewTester(func(c *ui.Context) {
		if _, err := LiveLine(c, CartesianOptions{Options: Options{Label: "Live", State: &state}, Series: data}, 2); err != nil {
			t.Fatal(err)
		}
	}, 400, 300)
	tt.Frame()
	if !reflect.DeepEqual(before, data[0].Points) {
		t.Fatal("live chart mutated caller data")
	}
	state.Handle.Focus()
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("1 · Sales: 8") {
		t.Fatal("live window did not trim the oldest sample")
	}
	data[0].Points = append(data[0].Points, Point{X: 3, Y: 16})
	tt.Frame()
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("3 · Sales: 16") {
		t.Fatal("live window did not consume the host's updated samples")
	}
}

func TestCustomPlotUsesHostThemeAndAnnotationDomain(t *testing.T) {
	saved := theme.PreferencesSnapshot()
	t.Cleanup(func() { theme.Configure(saved) })
	state := State{Selected: 0}
	var ctx PlotContext
	indicatorPainted := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		_, err := Line(c, CartesianOptions{
			Options: Options{State: &state, Indicator: func(p *ui.Painter, r ui.Rect, i int) {
				indicatorPainted = r.W > 0 && r.H > 0 && i == 0
			}},
			Series: samples(),
			Annotations: []Annotation{{Kind: Band, Points: []Point{{X: 0, Y: -4}, {X: 2, Y: 12}},
				Upper: []Point{{X: 0, Y: 4}, {X: 2, Y: 20}}}},
			Build: func(_ *ui.Context, plot PlotContext) { ctx = plot },
		})
		if err != nil {
			t.Fatal(err)
		}
	}, 500, 400)
	for _, ink := range []ui.Color{ui.Hex("#13579b"), ui.Hex("#9b5713")} {
		theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &theme.PalettePatch{Chart1: &ink}}})
		tt.Frame()
		if ctx.Y.Hi != 20 || ctx.Selected != 0 || !indicatorPainted {
			t.Fatalf("custom plot did not receive annotated scales/selection: %+v, paint %v", ctx, indicatorPainted)
		}
		img := tt.Image()
		want := color.RGBA{R: ink.R, G: ink.G, B: ink.B, A: ink.A}
		found := false
		for y := range img.Bounds().Dy() {
			for x := range img.Bounds().Dx() {
				if img.At(x, y) == want {
					found = true
				}
			}
		}
		if !found {
			t.Fatal("native series paint did not inherit the updated host chart token")
		}
	}
}

func TestSunburstDrillAndBreadcrumb(t *testing.T) {
	root := &HierarchyNode{ID: "all", Label: "All", Children: []*HierarchyNode{
		{ID: "branch", Label: "Branch", Children: []*HierarchyNode{{ID: "one", Label: "One", Value: 20}, {ID: "two", Label: "Two", Value: 30}}},
		{ID: "other", Label: "Other", Value: 10},
	}}
	state := SunburstState{State: State{Selected: -1}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if _, err := Sunburst(c, SunburstOptions{Options: Options{Label: "Hierarchy"}, Root: root, Drill: &state}); err != nil {
			t.Fatal(err)
		}
	}, 500, 400)
	tt.Frame()
	if err := tt.Click("Hierarchy"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("Branch: 50") {
		t.Fatalf("branch did not sum descendants: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if state.FocusID != "branch" || !tt.HasText("Branch") {
		t.Fatal("Enter did not drill")
	}
	if err := tt.Click("All"); err != nil {
		t.Fatal(err)
	}
	if state.FocusID != "all" {
		t.Fatal("breadcrumb did not navigate to root")
	}
}

func TestSelectionDoesNotMoveAxesAndInheritsDisabled(t *testing.T) {
	state := State{Selected: -1}
	count := 0
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).FillWidth().Disabled(disabled).Children(func() {
			_, err := Line(c, CartesianOptions{Options: Options{Label: "Stable axes", State: &state, OnSelect: func(int) { count++ }},
				Series: []Series{{ID: "a", Points: []Point{{X: 1, Y: 0}, {X: 2, Y: 10}}}}})
			if err != nil {
				t.Fatal(err)
			}
		})
	}, 600, 400)
	tt.Frame()
	zeroBefore, ok := tt.Find("0")
	if !ok {
		t.Fatal("missing zero tick")
	}
	if err := tt.Click("Stable axes"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnd)
	zeroAfter, _ := tt.Find("0")
	if zeroAfter != zeroBefore {
		t.Fatalf("inspection displaced axis: before %+v, after %+v", zeroBefore, zeroAfter)
	}
	disabled = true
	state.Selected = -1
	tt.Frame()
	before := count
	r, _ := tt.Find("Stable axes")
	tt.Move(r.X+r.W*.2, r.Y+r.H*.3)
	if state.Selected != -1 || count != before {
		t.Fatal("disabled ancestor allowed pointer selection")
	}
}

func TestInspectionRespectsTrailingDisabledAncestor(t *testing.T) {
	state := State{Selected: -1}
	calls, activations := 0, 0
	disabled := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			_, err := Pie(c, PieOptions{
				Options: Options{Label: "Disabled pie", State: &state,
					OnSelect: func(int) { calls++ }, OnActivate: func(int) { activations++ }},
				Values: []Value{{ID: "a", Value: 1}, {ID: "b", Value: 1}},
			})
			if err != nil {
				t.Fatal(err)
			}
		}).Disabled(disabled)
	}, 400, 300)
	tt.Frame()
	r, ok := tt.Find("Disabled pie")
	if !ok {
		t.Fatal("chart missing")
	}
	tt.Move(r.X+r.W/2+40, r.Y+r.H/2)
	if state.Selected != -1 || calls != 0 {
		t.Fatal("disabled ancestor allowed inspection")
	}
	disabled = false
	tt.Frame()
	tt.Move(r.X+r.W/2+45, r.Y+r.H/2)
	if state.Selected != 0 || calls != 1 {
		t.Fatalf("enabled chart did not inspect: selected %d, calls %d", state.Selected, calls)
	}
	state.Handle.Focus()
	tt.Key(0, ui.KeyEnter)
	if activations != 1 {
		t.Fatalf("enabled chart activated %d times", activations)
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Move(r.X+r.W/2-45, r.Y+r.H/2)
	if state.Selected != 0 || calls != 1 || activations != 1 {
		t.Fatal("disabled chart inspected or activated")
	}
}

func TestHugeFiniteProgressAndHierarchyRenderWithoutOverflow(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).FillWidth().Children(func() {
			if _, err := Ring(c, RingOptions{Options: Options{Height: 120}, Values: []Value{{ID: "huge", Value: 1e308, Max: 1e308}}}); err != nil {
				t.Fatal(err)
			}
			if _, err := Sunburst(c, SunburstOptions{Options: Options{Height: 120}, Root: &HierarchyNode{ID: "root", Children: []*HierarchyNode{
				{ID: "one", Value: 1e308}, {ID: "two", Value: 1e307},
			}}}); err != nil {
				t.Fatal(err)
			}
		})
	}, 400, 400)
	tt.Frame()
	if tt.Image() == nil {
		t.Fatal("no rendered result")
	}
}

func TestAggregateOverflowAndUnboundedCalendarAreErrors(t *testing.T) {
	cases := []func(*ui.Context) (ui.Element, error){
		func(c *ui.Context) (ui.Element, error) {
			return Pie(c, PieOptions{Values: []Value{{ID: "a", Value: 1e308}, {ID: "b", Value: 1e308}}})
		},
		func(c *ui.Context) (ui.Element, error) {
			return Sankey(c, SankeyOptions{Nodes: []FlowNode{{ID: "a"}, {ID: "b"}, {ID: "c"}}, Links: []FlowLink{{Source: "a", Target: "b", Value: 1e308}, {Source: "a", Target: "c", Value: 1e308}}})
		},
		func(c *ui.Context) (ui.Element, error) {
			return Line(c, CartesianOptions{Series: []Series{{ID: "a", Points: []Point{{X: 0, Y: 1}, {X: 1e20, Y: 2}}}}, Annotations: []Annotation{{Kind: Weekends}}})
		},
	}
	for i, build := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var err error
			tt := ui.NewTester(func(c *ui.Context) { _, err = build(c) }, 400, 300)
			tt.Frame()
			if err == nil {
				t.Fatal("unrepresentable geometry accepted")
			}
		})
	}
}

func TestLongFlowHasForwardRibbons(t *testing.T) {
	o := SankeyOptions{}
	for i := range 25 {
		o.Nodes = append(o.Nodes, FlowNode{ID: fmt.Sprint(i)})
		if i > 0 {
			o.Links = append(o.Links, FlowLink{Source: fmt.Sprint(i - 1), Target: fmt.Sprint(i), Value: 1})
		}
	}
	_, ribbons, err := layoutFlow(o)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range ribbons {
		if r.x2 <= r.x1 {
			t.Fatalf("backward ribbon: %+v", r)
		}
	}
}
