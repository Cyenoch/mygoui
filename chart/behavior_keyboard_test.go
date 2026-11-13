package chart

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Tab reaches an interactive chart and skips one without State. The arrows
// step through the measurements without wrapping, Left from no selection
// takes the last, and the chart's accessible value follows the selection.
// Enter with nothing selected activates nothing.
func TestChartArrowsStepWithoutWrappingAndValueFollows(t *testing.T) {
	state := State{Selected: -1}
	activations := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ui.Button(c, "Before")
			_, _ = Pie(c, PieOptions{Options: Options{Label: "Static share", Height: 80},
				Values: []Value{{ID: "a", Label: "A", Value: 1}}})
			_, _ = Line(c, CartesianOptions{Options: Options{Label: "Revenue", State: &state,
				OnActivate: func(int) { activations++ }}, Series: samples()})
		})
	}, 600, 500)
	if err := tt.Click("Before"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Revenue") {
		t.Fatal("Tab did not skip the static chart to reach the interactive one")
	}
	tt.Key(0, ui.KeyEnter)
	if activations != 0 {
		t.Fatal("Enter activated with no measurement selected")
	}
	for _, step := range []struct {
		key  ui.Key
		want int
	}{{ui.KeyLeft, 2}, {ui.KeyRight, 2}, {ui.KeyUp, 1}, {ui.KeyDown, 2}} {
		tt.Key(0, step.key)
		if state.Selected != step.want {
			t.Fatalf("%v selected %d, want %d", step.key, state.Selected, step.want)
		}
	}
	info, _ := tt.AccessByRole(ui.RoleGroup, "Revenue")
	if info.Value != "2 · Sales: 12" {
		t.Fatalf("accessible value = %q", info.Value)
	}
}

// Enter on a branch drills into it, clears the selection and keeps the
// focus on the chart, so the keyboard can keep drilling.
func TestSunburstEnterDrillKeepsFocusOnChart(t *testing.T) {
	root := &HierarchyNode{ID: "all", Label: "All", Children: []*HierarchyNode{
		{ID: "branch", Label: "Branch", Children: []*HierarchyNode{{ID: "one", Label: "One", Value: 20}}},
	}}
	state := SunburstState{State: State{Selected: -1}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if _, err := Sunburst(c, SunburstOptions{Options: Options{Label: "Hierarchy"}, Root: root, Drill: &state}); err != nil {
			t.Fatal(err)
		}
	}, 500, 400)
	if err := tt.Click("Hierarchy"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyEnter)
	if state.FocusID != "branch" || state.Selected != -1 || !tt.Focused("Hierarchy") {
		t.Fatalf("drill: focus id %q, selected %d, chart focused %v", state.FocusID, state.Selected, tt.Focused("Hierarchy"))
	}
}
