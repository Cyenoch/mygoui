package layout_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// app is the caller-owned state of a scaffold under test.
type app struct {
	nav, inspector                   float32
	navCollapsed, navOpen, inspectOn bool
	inspectorCollapsed               bool
	fixed, fixedInspector            bool
}

func (a *app) view(c *ui.Context) {
	layout.Scaffold(c, layout.ScaffoldOptions{
		TitleBar: layout.TitleBarOptions{
			Title: "Inbox", Label: "Window bar", MinHeight: 52,
			Leading: func() {
				layout.PaneToggle(c, layout.PaneToggleOptions{
					Collapsed: &a.navCollapsed, Open: &a.navOpen, Label: "Toggle mailboxes",
				})
			},
			Trailing: func() { button(c, "Compose message", nil) },
		},
		StartBar: layout.TitleBarOptions{Title: "Mail", Label: "Mailboxes bar"},
		Start: &layout.PaneOptions{
			Label: "Mailboxes", Size: &a.nav, Width: 220, Fixed: a.fixed,
			Collapsed: &a.navCollapsed, Open: &a.navOpen,
			Body: func() { ui.Text(c, "Archive") },
		},
		End: &layout.PaneOptions{
			Label: "Details", Size: &a.inspector, Width: 240, Fixed: a.fixedInspector,
			Collapsed: &a.inspectorCollapsed, Open: &a.inspectOn,
			Body: func() { ui.Text(c, "Sender") },
		},
		Content:   func() { ui.Text(c, "Message body") },
		StatusBar: func() { navigation.StatusBar(c, navigation.StatusBarOptions{Message: "Synced", Label: "Status"}) },
	})
}

func TestScaffoldSharesTheTitleBandBetweenPaneAndBar(t *testing.T) {
	const width, height = 1000, 600
	a := &app{}
	tt := newTester(a.view, width, height)
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80, Right: 138})
	pad := float32(theme.Space12)

	band, bar := rect(t, tt, "Mailboxes bar"), rect(t, tt, "Window bar")
	if band.H != bar.H || band.Y != 0 || bar.Y != 0 || !near(band.W+bar.W, width) {
		t.Fatalf("band %+v and bar %+v, want one height across the window's top", band, bar)
	}
	if r := rect(t, tt, "Mail"); !near(r.X, 80+pad) {
		t.Fatalf("the pane's band starts at %v, want after the start window controls", r.X)
	}
	// Beside a docked pane the bar reserves only the end controls.
	if r := rect(t, tt, "Toggle mailboxes"); !near(r.X, 220+pad) {
		t.Fatalf("leading control at %v, want right after the 220 DIP pane", r.X)
	}
	if r := rect(t, tt, "Compose message"); !near(r.X+r.W, width-138-pad) {
		t.Fatalf("trailing control ends at %v, want before the end window controls", r.X+r.W)
	}
	// Panes run from under their band to the strip; the strip spans the window.
	status := rect(t, tt, "Status")
	if status.W != width || status.Y+status.H != height {
		t.Fatalf("status strip = %+v, want the window's bottom edge", status)
	}
	if r := rect(t, tt, "Archive"); r.Y < band.H || r.X > 220 {
		t.Fatalf("start pane body = %+v, want it under its band", r)
	}
	if r := rect(t, tt, "Sender"); r.Y < bar.H || r.X < width-240 {
		t.Fatalf("end pane body = %+v, want it under the title bar at the end edge", r)
	}

	// With the pane gone the bar reaches the start edge and its controls.
	a.navCollapsed = true
	tt.Frame()
	if r := rect(t, tt, "Toggle mailboxes"); !near(r.X, 80+pad) {
		t.Fatalf("leading control at %v with the pane collapsed, want after the start window controls", r.X)
	}
}

func TestScaffoldBandGrowsWithAStackedBar(t *testing.T) {
	a := &app{}
	tt := newTester(a.view, 600, 600)
	tt.SetPreferences(ui.Preferences{TextScale: 2.5, ReduceMotion: true})
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80, Right: 138})
	a.navOpen = true
	tt.Frame()
	bar := rect(t, tt, "Window bar")
	if r := rect(t, tt, "Compose message"); r.Y < 52 || r.Y+r.H > bar.H {
		t.Fatalf("trailing control = %+v, want it below the band inside the %v DIP bar", r, bar.H)
	}
	// The compact start pane's sheet carries the same band, clear of the
	// window controls it covers.
	if r := rect(t, tt, "Mailboxes bar"); r.H != bar.H {
		t.Fatalf("pane band height = %v, want the bar's %v", r.H, bar.H)
	}
	// The sheet is nearly as wide as the window, so its band reaches the end
	// window controls as well.
	if r := rect(t, tt, "Mail"); r.X < 80 || r.X+r.W > 600-138 {
		t.Fatalf("pane band title = %+v, want it clear of the window controls on both sides", r)
	}
}

func TestScaffoldStartDividerLeavesTheContentItsMinimumBesideTheEndPane(t *testing.T) {
	const width = 1200
	for _, tc := range []struct {
		name           string
		fixed, without bool
		start          float32
	}{
		{"resizable end pane", false, false, width - 240 - 240},
		{"fixed end pane", true, false, width - 240 - 240},
		{"collapsed end pane", false, true, width - 240},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := &app{fixedInspector: tc.fixed, inspectorCollapsed: tc.without}
			tt := newTester(a.view, width, 600)
			drag(t, tt, "Mailboxes", 2000)
			if a.nav != tc.start {
				t.Fatalf("start pane = %v, want %v: the window less the content minimum and the end pane", a.nav, tc.start)
			}
			body := rect(t, tt, "Message body")
			if body.X < tc.start || (!tc.without && rect(t, tt, "Sender").X < tc.start+240) {
				t.Fatalf("content at %v, want 240 DIPs of it between the panes", body.X)
			}
		})
	}
}

// While the start pane folds, the title bar keeps its controls clear of the
// start window controls in every frame, and they move with the pane's edge.
func TestScaffoldTitleBarClearsTheWindowControlsWhileThePaneFolds(t *testing.T) {
	const inset = 80
	for _, fixed := range []bool{false, true} {
		name := "resizable pane"
		if fixed {
			name = "fixed pane"
		}
		t.Run(name, func(t *testing.T) {
			a := &app{fixed: fixed}
			tt := newTester(a.view, 1000, 600)
			tt.SetPreferences(ui.Preferences{TextScale: 1})
			tt.SetTitleBar(ui.TitleBar{Height: 52, Left: inset})
			clear := inset + float32(theme.Space12)
			docked := rect(t, tt, "Toggle mailboxes").X
			for _, collapse := range []bool{true, false} {
				last := rect(t, tt, "Toggle mailboxes").X
				a.navCollapsed = collapse
				tt.Frame()
				for range 20 {
					tt.Advance(16 * time.Millisecond)
					x := rect(t, tt, "Toggle mailboxes").X
					if x < clear-0.5 {
						t.Fatalf("collapse=%v: the toggle is at %v, under the window controls (%v)", collapse, x, clear)
					}
					// Collapsing moves it towards the window's edge, expanding away.
					if d := x - last; d > 60 || d < -60 || (collapse && d > 0.5) || (!collapse && d < -0.5) {
						t.Fatalf("collapse=%v: the toggle moved from %v to %v, want it to follow the pane", collapse, last, x)
					}
					last = x
				}
				want := docked
				if collapse {
					want = clear
				}
				if !near(last, want) {
					t.Fatalf("collapse=%v: the toggle rests at %v, want %v", collapse, last, want)
				}
			}
		})
	}
}
