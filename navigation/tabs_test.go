package navigation

import (
	"math"
	"slices"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

func TestTabsIntrinsicStripIndependentOfPanelWidth(t *testing.T) {
	for _, variant := range []TabsVariant{TabsDefault, TabsLine} {
		name := "default"
		if variant == TabsLine {
			name = "line"
		}
		t.Run(name, func(t *testing.T) {
			value := "overview"
			var root ui.Handle
			var strip ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				Tabs(c, &value, TabsOptions[string]{Label: "Views", Variant: variant}, []Tab[string]{
					{Value: "overview", Label: "Overview", Icon: lucide.AppWindow(), Content: func() { ui.Text(c, "Overview panel").FillWidth() }},
					{Value: "analytics", Label: "Analytics"},
					{Value: "reports", Label: "Reports"},
				}).Bind(&root)
			}, 900, 400)
			strip, ok := tt.Find("Views")
			if !ok || strip.W <= 100 || strip.W >= 500 {
				t.Fatalf("strip did not fit its labels in a wide root: %+v", strip)
			}
			initialWidth := strip.W
			panel, ok := tt.Find("Overview panel")
			if !ok || panel.W < 850 || panel.Y < strip.Y+strip.H {
				t.Fatalf("panel did not retain the wide root below its strip: %+v / %+v", panel, strip)
			}
			tt.SetSize(1200, 400)
			strip, _ = tt.Find("Views")
			if math.Abs(float64(strip.W-initialWidth)) > 1 || strip.X > 1 {
				t.Fatalf("root resize stretched or centered list: before=%g after=%+v", initialWidth, strip)
			}
			tt.SetDirection(ui.RightToLeft)
			strip, _ = tt.Find("Views")
			if math.Abs(float64(strip.W-initialWidth)) > 1 || math.Abs(float64(strip.X+strip.W-1200)) > 1 {
				t.Fatalf("RTL list lost intrinsic width or leading anchor: %+v", strip)
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyLeft)
			if value != "analytics" || !tt.Focused("Analytics") {
				t.Fatal("RTL arrows lost native selection")
			}
		})
	}
}

func TestTabsConstrainedParentLargeTextAndVertical(t *testing.T) {
	for _, variant := range []TabsVariant{TabsDefault, TabsLine} {
		name := "default"
		if variant == TabsLine {
			name = "line"
		}
		t.Run(name, func(t *testing.T) {
			value := "account"
			vertical := false
			tabs := []Tab[string]{
				{Value: "account", Label: "Account", Icon: lucide.User(), Content: func() {}},
				{Value: "security", Label: "Security", Content: func() {}},
				{Value: "disabled", Label: "Unavailable", Disabled: true, Content: func() { t.Error("disabled panel built") }},
			}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Width(308).Label("Constrained tabs parent").Children(func() {
					orientation := ui.Horizontal
					if vertical {
						orientation = ui.Vertical
					}
					Tabs(c, &value, TabsOptions[string]{Label: "Settings", Variant: variant, Orientation: orientation}, tabs)
				})
			}, 900, 500)
			tt.SetPreferences(ui.Preferences{TextScale: 2, ReduceMotion: true})
			for _, label := range []string{"Account", "Security", "Unavailable"} {
				r, ok := tt.Find(label)
				if !ok || r.X < 0 || r.X+r.W > 309 || r.H < 28 {
					t.Fatalf("200%% tab outside constrained parent: %s %+v", label, r)
				}
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			if value != "security" {
				t.Fatal("constrained list lost selection")
			}
			vertical = true
			tt.Frame()
			strip, _ := tt.Find("Settings")
			if strip.W >= 270 {
				t.Fatalf("vertical list stretched to parent: %+v", strip)
			}
			a, _ := tt.Find("Account")
			b, _ := tt.Find("Security")
			if b.Y < a.Y+a.H {
				t.Fatalf("vertical tabs overlap: %+v %+v", a, b)
			}
			tt.Key(0, ui.KeyHome)
			if value != "account" {
				t.Fatal("vertical Home did not select first enabled")
			}
			tt.Key(0, ui.KeyDown)
			if value != "security" {
				t.Fatal("vertical Down did not retain native navigation")
			}
			tt.SetDirection(ui.RightToLeft)
			a, _ = tt.Find("Account")
			parent, ok := tt.Find("Constrained tabs parent")
			if !ok || a.X < parent.X || a.X+a.W > parent.X+parent.W+1 {
				t.Fatalf("vertical RTL escaped parent: tab=%+v parent=%+v", a, parent)
			}
		})
	}
}

func TestTabsSelectionCallbacksDisabledAndPanelFocus(t *testing.T) {
	value := "overview"
	changes := []string{}
	tabs := []Tab[string]{
		{Value: "overview", Label: "Overview"},
		{Value: "disabled", Label: "Unavailable", Disabled: true},
		{Value: "reports", Label: "Reports"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Tabs(c, &value, TabsOptions[string]{Label: "Views", OnChange: func(v string) { changes = append(changes, v) }}, tabs)
		mygoui.Button(c, mygoui.ButtonOptions{Label: "After tabs", OnClick: func() { changes = append(changes, "after") }})
	}, 800, 350)
	if err := tt.Click("Reports"); err != nil {
		t.Fatal(err)
	}
	if value != "reports" || !slices.Equal(changes, []string{"reports"}) {
		t.Fatalf("selection=%s callbacks=%v", value, changes)
	}
	if err := tt.Click("Unavailable"); err != nil {
		t.Fatal(err)
	}
	if value != "reports" || len(changes) != 1 {
		t.Fatal("disabled tab selected")
	}
	// A pointer press on a disabled control clears native keyboard focus.
	// Re-enter the radio group before checking focus across a data reorder.
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Reports") {
		t.Fatal("Tab did not re-enter at selected tab after disabled pointer press")
	}
	slices.Reverse(tabs)
	tt.Frame()
	if !tt.Focused("Reports") {
		t.Fatal("reorder lost stable focus")
	}
	tt.Key(0, ui.KeyRight)
	if value != "overview" {
		t.Fatal("keyboard did not skip disabled")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After tabs") {
		t.Fatal("tab list created multiple Tab stops")
	}
}
