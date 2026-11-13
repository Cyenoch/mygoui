package chart

import (
	"bytes"
	"image"
	"math"
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// lineChart builds the sample line chart with state inside a column that
// counts the Escapes reaching it, as a dialog closing on Escape would.
func lineChart(t *testing.T, o Options, escapes *int) *ui.Tester {
	t.Helper()
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(20).OnShortcut(0, ui.KeyEscape, func() { *escapes++ }).Children(func() {
			if _, err := Line(c, CartesianOptions{Options: o, Series: samples()}); err != nil {
				t.Fatal(err)
			}
		})
	}, 600, 500)
}

// TestChartEscapeReachesTheDialogWithNothingSelected checks that Escape
// clears a selection and, with nothing selected, goes on to the dialog
// around the chart.
func TestChartEscapeReachesTheDialogWithNothingSelected(t *testing.T) {
	state := State{Selected: -1}
	escapes := 0
	tt := lineChart(t, Options{Label: "Revenue", State: &state}, &escapes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEscape)
	if state.Selected != -1 || escapes != 0 {
		t.Fatalf("Escape with a selection: selected %d, escapes reaching the dialog %d", state.Selected, escapes)
	}
	tt.Key(0, ui.KeyEscape)
	if escapes != 1 {
		t.Fatalf("Escape with nothing selected reached the dialog %d times, want 1", escapes)
	}
}

// TestChartSelectionFromThePointerEndsWhenItLeaves checks that hovering
// selects and leaving the plot clears that selection, while a selection
// made with the keys stays.
func TestChartSelectionFromThePointerEndsWhenItLeaves(t *testing.T) {
	state := State{Selected: -1}
	escapes := 0
	tt := lineChart(t, Options{Label: "Revenue", State: &state}, &escapes)
	plot, ok := tt.Find("Revenue")
	if !ok {
		t.Fatal("chart missing")
	}
	tt.Move(plot.X+plot.W-4, plot.Y+plot.H/2)
	if state.Selected != 2 {
		t.Fatalf("hover selected %d, want 2", state.Selected)
	}
	tt.Move(2, 2)
	if state.Selected != -1 {
		t.Fatalf("after leaving the plot selected %d, want -1", state.Selected)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Move(plot.X+plot.W/2, plot.Y+plot.H/2)
	tt.Move(2, 2)
	if state.Selected != -1 {
		t.Fatalf("a hover after the keys did not end with the pointer: selected %d", state.Selected)
	}
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	if state.Selected != 0 {
		t.Fatalf("a selection made with the keys did not stay: selected %d", state.Selected)
	}
}

// TestLoadingChartIsBusyAndKeepsItsTabStop checks F17: a loading chart is
// a busy region that Tab still reaches.
func TestLoadingChartIsBusyAndKeepsItsTabStop(t *testing.T) {
	state := State{Selected: -1}
	escapes := 0
	tt := lineChart(t, Options{Label: "Revenue", State: &state, Loading: true}, &escapes)
	info, ok := tt.AccessByRole(ui.RoleGroup, "Revenue")
	if !ok || !info.Busy || !info.Focusable {
		t.Fatalf("loading chart: found %v, busy %v, focusable %v", ok, info.Busy, info.Focusable)
	}
}

// TestChartFocusRingShowsOnlyForKeyboardFocus checks that a click focuses
// the plot without a ring, and Tab shows it.
func TestChartFocusRingShowsOnlyForKeyboardFocus(t *testing.T) {
	state := State{Selected: -1}
	escapes := 0
	tt := lineChart(t, Options{Label: "Revenue", State: &state}, &escapes)
	rest := tt.Image()
	plot, _ := tt.Find("Revenue")
	tt.ClickAt(plot.X+plot.W/2, plot.Y+plot.H/2)
	tt.Move(2, 2)
	if !tt.Focused("Revenue") {
		t.Fatal("a click did not focus the chart")
	}
	if !sameImage(rest, tt.Image()) {
		t.Fatal("pointer focus painted a focus ring")
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Revenue") || sameImage(rest, tt.Image()) {
		t.Fatal("keyboard focus painted no focus ring")
	}
}

func sameImage(a, b *image.RGBA) bool { return bytes.Equal(a.Pix, b.Pix) }

// TestChartValueSpansTheDataWithoutASelection checks that a chart with
// nothing selected, or no State at all, still exposes its data: the first
// and last measurements.
func TestChartValueSpansTheDataWithoutASelection(t *testing.T) {
	state := State{Selected: -1}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			_, _ = Line(c, CartesianOptions{Options: Options{Label: "Revenue", State: &state}, Series: samples()})
			_, _ = Pie(c, PieOptions{Options: Options{Label: "Share", Height: 80},
				Values: []Value{{ID: "a", Label: "A", Value: 1}, {ID: "b", Label: "B", Value: 3}}})
		})
	}, 600, 600)
	for label, want := range map[string]string{"Revenue": "0 · Sales: -4 – 2 · Sales: 12", "Share": "A: 1"} {
		info, _ := tt.AccessByRole(ui.RoleGroup, label)
		if !strings.HasPrefix(info.Value, want) {
			t.Fatalf("%s value = %q, want it to start %q", label, info.Value, want)
		}
	}
}

// TestChartKeysFollowTheDrawing checks that the keys move the way the
// measurements are drawn: horizontal bars run up the plot, PageDown moves a
// tenth of the data, and right to left the plot, which does not mirror,
// still steps forward with Right. One X reads every series at it.
func TestChartKeysFollowTheDrawing(t *testing.T) {
	points := make([]Point, 30)
	other := make([]Point, 30)
	for i := range points {
		points[i] = Point{X: float64(i), Y: float64(i)}
		other[i] = Point{X: float64(i), Y: float64(2 * i)}
	}
	series := []Series{{ID: "a", Label: "A", Points: points}, {ID: "b", Label: "B", Points: other}}
	bars, line := State{Selected: -1}, State{Selected: -1}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			_, _ = Bar(c, CartesianOptions{Options: Options{Label: "Bars", State: &bars, Height: 120}, Series: series[:1], Horizontal: true})
			_, _ = Line(c, CartesianOptions{Options: Options{Label: "Line", State: &line, Height: 120}, Series: series})
		})
	}, 600, 600)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyUp)
	if bars.Selected != 1 {
		t.Fatalf("Up on horizontal bars selected %d, want 1", bars.Selected)
	}
	tt.Key(0, ui.KeyDown)
	if bars.Selected != 0 {
		t.Fatalf("Down on horizontal bars selected %d, want 0", bars.Selected)
	}
	tt.Key(0, ui.KeyTab)
	tt.SetDirection(ui.RightToLeft)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageDown)
	if line.Selected != 4 {
		t.Fatalf("Right then PageDown right to left selected %d, want 4", line.Selected)
	}
	tt.Key(0, ui.KeyPageUp)
	if line.Selected != 1 {
		t.Fatalf("PageUp selected %d, want 1", line.Selected)
	}
	info, _ := tt.AccessByRole(ui.RoleGroup, "Line")
	if info.Value != "1 · A: 1; 1 · B: 2" {
		t.Fatalf("readout = %q, want every series at the X", info.Value)
	}
}

// TestHeatmapArrowsMoveAcrossRowsAndColumns checks that Up and Down change
// rows and the arrows pass over missing cells.
func TestHeatmapArrowsMoveAcrossRowsAndColumns(t *testing.T) {
	state := State{Selected: -1}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		_, err := Heatmap(c, HeatmapOptions{Options: Options{Label: "Load", State: &state},
			Columns: []string{"Mon", "Tue", "Wed"}, Rows: []string{"AM", "PM"},
			Cells: []Cell{{Column: 0, Row: 0, Value: 1}, {Column: 2, Row: 0, Value: 2}, {Column: 0, Row: 1, Value: 3}}})
		if err != nil {
			t.Fatal(err)
		}
	}, 600, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyUp)
	if state.Selected != 2 {
		t.Fatalf("Up from Mon AM selected %d, want Mon PM (2)", state.Selected)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if state.Selected != 1 {
		t.Fatalf("Right from Mon AM selected %d, want Wed AM (1) past the missing Tue", state.Selected)
	}
	tt.Key(0, ui.KeyRight)
	if state.Selected != 1 {
		t.Fatalf("Right at the last column moved to %d", state.Selected)
	}
}

// TestGaugeIsAMeter checks that a gauge exposes its value as a meter
// without being selected.
func TestGaugeIsAMeter(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		_, _ = Gauge(c, GaugeOptions{Options: Options{Label: "Disk"}, Value: 70})
	}, 400, 300)
	info, ok := tt.AccessByRole(ui.RoleMeter, "Disk")
	if !ok || info.Min != 0 || info.Max != 100 || info.Now != 70 {
		t.Fatalf("gauge meter: found %v, range %v-%v, now %v", ok, info.Min, info.Max, info.Now)
	}
}

// TestSankeyRightFollowsTheFlow checks that Right goes from a node to its
// outgoing link and on to the link's target, Left back, and that links read
// by node label.
func TestSankeyRightFollowsTheFlow(t *testing.T) {
	state := State{Selected: -1}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		_, err := Sankey(c, SankeyOptions{Options: Options{Label: "Flow", State: &state},
			Nodes: []FlowNode{{"in", "Visits"}, {"out", "Orders"}, {"gone", "Exits"}},
			Links: []FlowLink{{"in", "out", 3}, {"in", "gone", 7}}})
		if err != nil {
			t.Fatal(err)
		}
	}, 600, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyRight)
	info, _ := tt.AccessByRole(ui.RoleGroup, "Flow")
	if state.Selected != 3 || info.Value != "Visits → Orders: 3" {
		t.Fatalf("Right from Visits selected %d reading %q", state.Selected, info.Value)
	}
	tt.Key(0, ui.KeyRight)
	if state.Selected != 1 {
		t.Fatalf("Right from the link selected %d, want Orders (1)", state.Selected)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if state.Selected != 0 {
		t.Fatalf("Left twice from Orders selected %d, want Visits (0)", state.Selected)
	}
}

// TestSunburstEscapeGoesUpALevel checks that Escape with nothing selected
// leaves a drilled level, and that the path is a breadcrumb whose current
// node is not a link.
func TestSunburstEscapeGoesUpALevel(t *testing.T) {
	root := &HierarchyNode{ID: "all", Label: "All", Children: []*HierarchyNode{
		{ID: "branch", Label: "Branch", Children: []*HierarchyNode{{ID: "one", Label: "One", Value: 20}}},
		{ID: "other", Label: "Other", Value: 5},
	}}
	state := SunburstState{State: State{Selected: -1}, FocusID: "branch"}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if _, err := Sunburst(c, SunburstOptions{Options: Options{Label: "Hierarchy"}, Root: root, Drill: &state}); err != nil {
			t.Fatal(err)
		}
	}, 500, 400)
	if _, ok := tt.AccessByRole(ui.RoleLink, "All"); !ok {
		t.Fatal("the root is not a link in the path")
	}
	if _, ok := tt.AccessByRole(ui.RoleLink, "Branch"); ok {
		t.Fatal("the current node is a link")
	}
	if a, ok := tt.Access("Branch"); !ok || a.Current != ui.CurrentPage {
		t.Fatalf("the current node = %+v, %v; want marked current", a, ok)
	}
	if err := tt.Click("Hierarchy"); err != nil {
		t.Fatal(err)
	}
	tt.Move(2, 2)
	tt.Key(0, ui.KeyEscape)
	if state.FocusID != "all" {
		t.Fatalf("Escape left the focus at %q, want all", state.FocusID)
	}
}

// TestChoroplethWheelPinchAndDragMoveTheMap checks that the wheel and a
// pinch zoom and a drag pans the caller's viewport.
func TestChoroplethWheelPinchAndDragMoveTheMap(t *testing.T) {
	square := func(lon float64) Polygon {
		return Polygon{Rings: [][]GeoPoint{{{lon, 0}, {lon + 10, 0}, {lon + 10, 10}, {lon, 10}}}}
	}
	view := MapState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		_, err := Choropleth(c, ChoroplethOptions{Options: Options{Label: "Map"}, Viewport: &view,
			Regions: []Region{{ID: "a", Label: "A", Value: 1, Polygons: []Polygon{square(0)}},
				{ID: "b", Label: "B", Value: 2, Polygons: []Polygon{square(20)}}}})
		if err != nil {
			t.Fatal(err)
		}
	}, 600, 500)
	plot, ok := tt.Find("Map")
	if !ok {
		t.Fatal("map missing")
	}
	cx, cy := plot.X+plot.W/2, plot.Y+plot.H/2
	tt.Scroll(cx, cy, 0, -70)
	if math.Abs(view.Zoom-math.Exp(0.7)) > 1e-6 || view.PanX != 0 || view.PanY != 0 {
		t.Fatalf("wheel at the center: %+v", view)
	}
	tt.Pinch(cx, cy, 0.1)
	if view.Zoom != 1 {
		t.Fatalf("pinching out stopped at %v, want the fit (1)", view.Zoom)
	}
	tt.Press(cx, cy)
	tt.Move(cx+60, cy-25)
	tt.Release(cx+60, cy-25)
	if math.Abs(view.PanX*float64(plot.W)-60) > 0.5 || math.Abs(view.PanY*float64(plot.H)+25) > 0.5 {
		t.Fatalf("a drag of 60, -25 panned %v, %v DIPs", view.PanX*float64(plot.W), view.PanY*float64(plot.H))
	}
}

// TestChartTextsComeFromTheLocale checks the empty text, the zoom buttons
// and the brush name.
func TestChartTextsComeFromTheLocale(t *testing.T) {
	locale.Set(locale.Messages{NoData: "Aucune donnée", ZoomIn: "Agrandir", Range: "Plage"})
	t.Cleanup(func() { locale.Set(locale.Messages{}) })
	brush := BrushState{From: 0, To: 1}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			_, _ = Line(c, CartesianOptions{Options: Options{Label: "Empty"}})
			_, _ = Choropleth(c, ChoroplethOptions{Options: Options{Label: "Map"}, Viewport: &MapState{}})
			_, _ = Brush(c, &brush, 0, 1, false)
		})
	}, 600, 700)
	if !tt.HasText("Aucune donnée") {
		t.Fatal("the empty text is not the locale's")
	}
	if _, ok := tt.Find("Agrandir"); !ok {
		t.Fatal("the zoom button is not named by the locale")
	}
	if _, ok := tt.Find("Plage"); !ok {
		t.Fatal("the brush is not named by the locale")
	}
}
