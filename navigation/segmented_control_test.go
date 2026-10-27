package navigation

import (
	"image/color"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSegmentedControlChangesSelectionWithoutBuildingPanels(t *testing.T) {
	selected := 0
	var changes []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		SegmentedControl(c, SegmentedControlOptions{
			Label: "View", Selected: &selected,
			Items: []TabOptions{
				{Label: "List", Content: func() { ui.Text(c, "Unexpected panel") }},
				{Label: "Unavailable", Disabled: true},
				{Label: "Grid"},
			},
			OnChange: func(index int) { changes = append(changes, index) },
		})
	}, 400, 160)
	if tt.HasText("Unexpected panel") {
		t.Fatal("a segment built tab panel content")
	}
	if err := tt.Click("List"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if selected != 2 || !tt.Focused("Grid") {
		t.Fatalf("Right arrow: selected %d, grid focused %v", selected, tt.Focused("Grid"))
	}
	if len(changes) != 1 || changes[0] != 2 {
		t.Fatalf("change callbacks: %v", changes)
	}
	if err := tt.Click("Unavailable"); err != nil {
		t.Fatal(err)
	}
	if selected != 2 || len(changes) != 1 {
		t.Fatalf("disabled segment changed selection: %d, %v", selected, changes)
	}
}

func TestSegmentedControlUpdatesSelectionBeforeReturningOnRelease(t *testing.T) {
	selected := 0
	var observed []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		SegmentedControl(c, SegmentedControlOptions{
			Selected: &selected,
			Items:    []TabOptions{{Label: "List"}, {Label: "Grid"}},
		})
		observed = append(observed, selected)
	}, 320, 120)
	grid, ok := tt.Find("Grid")
	if !ok {
		t.Fatal("missing Grid segment")
	}
	x, y := grid.X+grid.W/2, grid.Y+grid.H/2
	tt.Press(x, y)
	observed = nil
	tt.Release(x, y)
	if len(observed) == 0 || observed[0] != 1 {
		t.Fatalf("selection observed after activation: %v", observed)
	}
}

func TestSegmentedControlIndicatorMovesAndRespectsReducedMotion(t *testing.T) {
	for _, reduced := range []bool{false, true} {
		t.Run(map[bool]string{false: "animated", true: "reduced"}[reduced], func(t *testing.T) {
			selected := 0
			var colors theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
					SegmentedControl(c, SegmentedControlOptions{
						Label: "View", Selected: &selected,
						Items: []TabOptions{{Label: "List"}, {Label: "Grid"}},
					})
				})
			}, 360, 140)
			if reduced {
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			}
			tt.Frame()
			tt.Frame()
			if err := tt.Click("Grid"); err != nil {
				t.Fatal(err)
			}
			tt.Move(2, 2)
			tt.Frame()
			list, _ := tt.Find("List")
			grid, _ := tt.Find("Grid")
			oldX, newX := int(list.X+list.W/2), int(grid.X+grid.W/2)
			y := int(grid.Y + 3)
			image := tt.Image()
			if reduced {
				wantSegmentPixel(t, image.RGBAAt(newX, y), colors.Background, "reduced-motion target")
			} else {
				wantSegmentPixel(t, image.RGBAAt(oldX, y), colors.Background, "initial indicator position")
				tt.Advance(100 * time.Millisecond)
				mid := tt.Image()
				middleX := int((list.X + grid.X + grid.W) / 2)
				wantSegmentPixel(t, mid.RGBAAt(middleX, y), colors.Background, "mid-slide indicator")
				wantSegmentPixel(t, mid.RGBAAt(oldX, y), colors.MutedBackground.Over(colors.Background), "mid-slide old position")
				wantSegmentPixel(t, mid.RGBAAt(newX, y), colors.MutedBackground.Over(colors.Background), "mid-slide new position")
			}
			tt.Advance(300 * time.Millisecond)
			wantSegmentPixel(t, tt.Image().RGBAAt(newX, y), colors.Background, "settled indicator")
			wantSegmentPixel(t, tt.Image().RGBAAt(oldX, y), colors.MutedBackground.Over(colors.Background), "settled old position")
		})
	}
}

func TestTabsKeyboardSkipsDisabledTab(t *testing.T) {
	selected := 1
	var changes []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Tabs(c, TabsOptions{Selected: &selected, Items: []TabOptions{
			{Label: "First disabled", Disabled: true}, {Label: "One"},
			{Label: "Locked", Disabled: true}, {Label: "Three"},
			{Label: "Last disabled", Disabled: true},
		}, OnChange: func(index int) { changes = append(changes, index) }})
	}, 400, 160)
	if err := tt.Click("One"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if selected != 3 || !tt.Focused("Three") {
		t.Fatalf("selected %d, third tab focused %v", selected, tt.Focused("Three"))
	}
	tt.Key(0, ui.KeyHome)
	if selected != 1 || !tt.Focused("One") {
		t.Fatalf("Home selected %d, first enabled tab focused %v", selected, tt.Focused("One"))
	}
	tt.Key(0, ui.KeyEnd)
	if selected != 3 || !tt.Focused("Three") {
		t.Fatalf("End selected %d, last enabled tab focused %v", selected, tt.Focused("Three"))
	}
	if len(changes) != 3 || changes[0] != 3 || changes[1] != 1 || changes[2] != 3 {
		t.Fatalf("change callbacks: %v", changes)
	}
}

func TestTabsKeepsSelectionUsableWhenItemDisappears(t *testing.T) {
	count := 3
	var changes []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		items := []TabOptions{
			{Label: "One", Content: func() { ui.Text(c, "One panel") }},
			{Label: "Two", Content: func() { ui.Text(c, "Two panel") }},
			{Label: "Three", Content: func() { ui.Text(c, "Three panel") }},
		}
		Tabs(c, TabsOptions{Items: items[:count], OnChange: func(i int) { changes = append(changes, i) }})
	}, 480, 240)
	if err := tt.Click("Three"); err != nil {
		t.Fatal(err)
	}
	count = 2
	tt.Frame()
	if !tt.HasText("Two panel") {
		t.Fatalf("no fallback panel after removing the selected tab: %v", tt.Texts())
	}
	if err := tt.Click("One"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("One panel") {
		t.Fatalf("remaining tab could not be selected: %v", tt.Texts())
	}
	if len(changes) != 2 || changes[0] != 2 || changes[1] != 0 {
		t.Fatalf("callbacks for user choices: %v", changes)
	}
}

func TestTabsReconcilesDisabledSelection(t *testing.T) {
	selected := 1
	var changes []int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Tabs(c, TabsOptions{Selected: &selected, Items: []TabOptions{
			{Label: "One", Content: func() { ui.Text(c, "One panel") }},
			{Label: "Locked", Disabled: true, Content: func() { ui.Text(c, "Locked panel") }},
			{Label: "Three", Content: func() { ui.Text(c, "Three panel") }},
		}, OnChange: func(i int) { changes = append(changes, i) }})
	}, 480, 240)
	if selected != 0 || !tt.HasText("One panel") || tt.HasText("Locked panel") {
		t.Fatalf("disabled selection was not reconciled: selected=%d texts=%v", selected, tt.Texts())
	}
	if err := tt.Click("Three"); err != nil {
		t.Fatal(err)
	}
	if selected != 2 || !tt.HasText("Three panel") {
		t.Fatalf("enabled selection failed: selected=%d texts=%v", selected, tt.Texts())
	}
	if len(changes) != 1 || changes[0] != 2 {
		t.Fatalf("callbacks for user choice: %v", changes)
	}
}

func TestTabsRestoresSelectionWhenAnItemBecomesEnabled(t *testing.T) {
	selected := 0
	disabled := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Tabs(c, TabsOptions{Selected: &selected, Items: []TabOptions{
			{Label: "One", Disabled: disabled, Content: func() { ui.Text(c, "One panel") }},
			{Label: "Two", Disabled: true, Content: func() { ui.Text(c, "Two panel") }},
		}})
	}, 480, 240)
	if selected != -1 || tt.HasText("One panel") || tt.HasText("Two panel") {
		t.Fatalf("all-disabled tabs: selected=%d texts=%v", selected, tt.Texts())
	}
	disabled = false
	tt.Frame()
	if selected != 0 || !tt.HasText("One panel") {
		t.Fatalf("enabled tab did not become available: selected=%d texts=%v", selected, tt.Texts())
	}
}

func TestTabsRestoresSelectionWhenItemsReturn(t *testing.T) {
	selected := 0
	count := 1
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		items := []TabOptions{{Label: "One", Content: func() { ui.Text(c, "One panel") }}}
		Tabs(c, TabsOptions{Selected: &selected, Items: items[:count]})
	}, 480, 240)
	count = 0
	tt.Frame()
	if selected != -1 || tt.HasText("One panel") {
		t.Fatalf("empty tabs: selected=%d texts=%v", selected, tt.Texts())
	}
	count = 1
	tt.Frame()
	if selected != 0 || !tt.HasText("One panel") {
		t.Fatalf("restored tabs: selected=%d texts=%v", selected, tt.Texts())
	}
}

func TestSelectionIndicatorsFollowRelayout(t *testing.T) {
	for _, tabs := range []bool{false, true} {
		name := "segments"
		if tabs {
			name = "tabs"
		}
		t.Run(name, func(t *testing.T) {
			selected := 0
			spacer := float32(20)
			var colors theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
					ui.Box(c).Height(spacer).Shrink(0)
					items := []TabOptions{{Label: "List"}, {Label: "Grid"}}
					if tabs {
						Tabs(c, TabsOptions{Selected: &selected, Items: items})
					} else {
						SegmentedControl(c, SegmentedControlOptions{Selected: &selected, Items: items})
					}
				})
			}, 360, 240)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			spacer = 70
			tt.Frame()
			choice, ok := tt.Find("List")
			if !ok {
				t.Fatal("selected item missing after move")
			}
			wantSegmentPixel(t, tt.Image().RGBAAt(int(choice.X+choice.W/2), int(choice.Y+3)),
				colors.Background, "indicator after parent movement")
			tt.Frame()
			if tt.FramePending() {
				t.Fatal("indicator keeps requesting frames after layout settles")
			}
		})
	}
}

func TestSegmentedIndicatorRemeasuresChangedLabels(t *testing.T) {
	selected := 1
	first := "One"
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
			SegmentedControl(c, SegmentedControlOptions{Selected: &selected, Items: []TabOptions{
				{Label: first}, {Label: "Selected"},
			}})
		})
	}, 440, 140)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	first = "A much wider first segment"
	tt.Frame()
	tt.Frame()
	item, ok := tt.Find("Selected")
	if !ok {
		t.Fatal("selected segment missing after label change")
	}
	wantSegmentPixel(t, tt.Image().RGBAAt(int(item.X+item.W/2), int(item.Y+3)),
		colors.Background, "indicator after preceding label changes")
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("indicator keeps requesting frames after remeasurement")
	}
}

func wantSegmentPixel(t *testing.T, got color.RGBA, want ui.Color, where string) {
	t.Helper()
	r, g, b, _ := got.RGBA()
	near := func(a uint32, b uint8) bool {
		d := int(a>>8) - int(b)
		return d >= -3 && d <= 3
	}
	if !near(r, want.R) || !near(g, want.G) || !near(b, want.B) {
		t.Fatalf("%s: got RGB %d,%d,%d want %d,%d,%d", where, r>>8, g>>8, b>>8, want.R, want.G, want.B)
	}
}
