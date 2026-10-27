package mygoui_test

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSidebarSelectionAndDisabledItems(t *testing.T) {
	selected, collapsed := "home", false
	var selections []string
	var actions []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{
				Label: "Workspace",
				Items: []mygoui.SidebarItemOptions{
					{ID: "home", Label: "Home", Icon: lucide.House()},
					{ID: "locked", Label: "Locked", Icon: lucide.Lock(), Disabled: true,
						OnClick: func() { actions = append(actions, "locked") }},
					{ID: "settings", Label: "Settings", Icon: lucide.Settings(),
						OnClick: func() { actions = append(actions, selected) }},
				},
			}},
			OnSelect: func(id string) { selections = append(selections, id) },
		})
	}, 400, 320)

	if err := tt.Click("Locked"); err != nil {
		t.Fatal(err)
	}
	if selected != "home" || len(selections) != 0 || len(actions) != 0 {
		t.Fatalf("disabled click changed navigation: selected=%q selections=%v actions=%v", selected, selections, actions)
	}

	if err := tt.Click("Settings"); err != nil {
		t.Fatal(err)
	}
	if selected != "settings" || len(selections) != 1 || selections[0] != "settings" ||
		len(actions) != 1 || actions[0] != "settings" || !tt.Focused("Settings") {
		t.Fatalf("pointer activation: selected=%q selections=%v actions=%v focus=%v",
			selected, selections, actions, tt.Focused("Settings"))
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if len(selections) != 3 || len(actions) != 3 {
		t.Fatalf("keyboard reactivation: selections=%v actions=%v", selections, actions)
	}

	if err := tt.Click("Home"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Settings") {
		t.Fatalf("Tab did not skip disabled item: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if selected != "settings" || len(selections) != 5 || len(actions) != 4 {
		t.Fatalf("keyboard navigation: selected=%q selections=%v actions=%v", selected, selections, actions)
	}
}

func TestSidebarCollapseAndSlots(t *testing.T) {
	selected, collapsed := "home", false
	var changes []bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Label: "Navigation", ToggleLabel: "Show navigation",
			Width: 200, CollapsedWidth: 64,
			Header: func() {
				if !collapsed {
					ui.Text(c, "Account")
				}
			},
			Footer: func() {
				if !collapsed {
					ui.Text(c, "Profile")
				}
			},
			Groups: []mygoui.SidebarGroupOptions{{
				Label: "Workspace",
				Items: []mygoui.SidebarItemOptions{
					{ID: "home", Label: "Home", Icon: lucide.House()},
					{ID: "settings", Label: "Settings", Icon: lucide.Settings()},
				},
			}},
			OnCollapseChange: func(value bool) { changes = append(changes, value) },
		})
	}, 400, 320)

	box, ok := tt.Find("Navigation")
	workspaceLabels := 0
	for _, label := range tt.Texts() {
		if label == "Workspace" {
			workspaceLabels++
		}
	}
	if workspaceLabels != 2 {
		t.Fatalf("expanded group is missing its visible label: %v", tt.Texts())
	}
	if !ok || box.W != 200 || !tt.HasText("Workspace") || !tt.HasText("Account") || !tt.HasText("Profile") {
		t.Fatalf("expanded sidebar: bounds=%+v text=%v", box, tt.Texts())
	}
	if err := tt.Click("Show navigation"); err != nil {
		t.Fatal(err)
	}
	box, ok = tt.Find("Navigation")
	workspaceLabels = 0
	for _, label := range tt.Texts() {
		if label == "Workspace" {
			workspaceLabels++
		}
	}
	if workspaceLabels != 1 {
		t.Fatalf("collapsed group should retain only its accessible name: %v", tt.Texts())
	}
	if !ok || box.W != 64 || !collapsed || len(changes) != 1 || !changes[0] {
		t.Fatalf("collapsed sidebar: bounds=%+v collapsed=%v changes=%v", box, collapsed, changes)
	}
	if tt.HasText("Account") || tt.HasText("Profile") {
		t.Fatalf("collapsed content did not adapt: %v", tt.Texts())
	}
	if _, ok := tt.Find("Settings"); !ok {
		t.Fatal("collapsed icon lost its accessible label")
	}
	if err := tt.Click("Settings"); err != nil {
		t.Fatal(err)
	}
	if selected != "settings" {
		t.Fatalf("collapsed navigation selected %q", selected)
	}
	if err := tt.Click("Show navigation"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeySpace)
	box, ok = tt.Find("Navigation")
	if !ok || box.W != 64 || !collapsed || len(changes) != 3 || changes[1] || !changes[2] {
		t.Fatalf("keyboard toggle: bounds=%+v collapsed=%v changes=%v", box, collapsed, changes)
	}

	collapsed = false
	tt.Frame()
	box, ok = tt.Find("Navigation")
	if !ok || box.W != 200 || !tt.HasText("Workspace") || selected != "settings" {
		t.Fatalf("caller expansion: bounds=%+v selected=%q text=%v", box, selected, tt.Texts())
	}
}

func TestSidebarCollapsedTextOnlyNavigation(t *testing.T) {
	selected, collapsed := "", true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{Items: []mygoui.SidebarItemOptions{
				{ID: "help", Label: "Help"},
			}}},
		})
	}, 300, 200)
	box, ok := tt.Find("Sidebar")
	if !ok || box.W != 48 {
		t.Fatalf("collapsed default width: %+v", box)
	}
	if err := tt.Click("Help"); err != nil {
		t.Fatal(err)
	}
	if selected != "help" {
		t.Fatalf("collapsed text-only item selected %q", selected)
	}
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	box, ok = tt.Find("Sidebar")
	if !ok || box.W != 96 {
		t.Fatalf("collapsed width at enlarged text scale: %+v", box)
	}
	if err := tt.Click("Help"); err != nil {
		t.Fatal(err)
	}
	if selected != "help" {
		t.Fatalf("enlarged collapsed item selected %q", selected)
	}
}

func TestSidebarTokenPaintAndFocus(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(fmt.Sprintf("dark=%v", dark), func(t *testing.T) {
			selected, collapsed := "home", false
			var colors theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, dark)
				colors = theme.Colors(c)
				mygoui.Sidebar(c, mygoui.SidebarOptions{
					Selected: &selected, Collapsed: &collapsed,
					Groups: []mygoui.SidebarGroupOptions{{Items: []mygoui.SidebarItemOptions{
						{ID: "home", Label: "Home"},
						{ID: "settings", Label: "Settings"},
					}}},
				})
			}, 400, 240)

			home, _ := tt.Find("Home")
			settings, _ := tt.Find("Settings")
			img := tt.Image()
			wantPixel(t, img, 100, 8, colors.Sidebar, "sidebar background")
			wantPixel(t, img, 100, int(home.Y+home.H/2),
				colors.SidebarAccent.Over(colors.Sidebar), "selected item")
			wantPixel(t, img, 100, int(settings.Y+settings.H/2), colors.Sidebar, "unselected item")

			tt.Move(100, settings.Y+settings.H/2)
			tt.Frame()
			wantPixel(t, tt.Image(), 100, int(settings.Y+settings.H/2),
				colors.SidebarAccent.Over(colors.Sidebar), "hovered item")

			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Toggle Sidebar") {
				t.Fatal("toggle not keyboard focused")
			}
			toggle, _ := tt.Find("Toggle Sidebar")
			wantPixel(t, tt.Image(), int(toggle.X+toggle.W/2), int(toggle.Y)+1,
				colors.SidebarRing.Over(colors.Sidebar), "sidebar focus ring")
		})
	}
}

func TestSidebarUsesInjectedTokens(t *testing.T) {
	old := theme.PreferencesSnapshot()
	defer theme.Configure(old)
	background := ui.RGB(210, 230, 240)
	accent := ui.RGB(30, 90, 160)
	border := ui.RGB(200, 80, 70)
	ring := ui.RGB(40, 180, 80)
	prefs := old
	prefs.Palette.Light = &theme.PalettePatch{
		Sidebar: &background, SidebarAccent: &accent, SidebarBorder: &border, SidebarRing: &ring,
	}
	theme.Configure(prefs)

	selected, collapsed := "home", false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{Items: []mygoui.SidebarItemOptions{
				{ID: "home", Label: "Home"},
			}}},
		})
	}, 400, 200)
	home, _ := tt.Find("Home")
	img := tt.Image()
	wantPixel(t, img, 100, 8, background, "injected sidebar background")
	wantPixel(t, img, 100, int(home.Y+home.H/2), accent, "injected sidebar accent")
	wantPixel(t, img, 255, 100, border, "injected sidebar border")
	tt.Key(0, ui.KeyTab)
	toggle, _ := tt.Find("Toggle Sidebar")
	wantPixel(t, tt.Image(), int(toggle.X+toggle.W/2), int(toggle.Y)+1, ring, "injected sidebar ring")
}

func TestSidebarFocusFollowsStableItems(t *testing.T) {
	selected, collapsed := "", false
	groups := []mygoui.SidebarGroupOptions{
		{Key: "first", Label: "First", Items: []mygoui.SidebarItemOptions{
			{ID: "alpha", Label: "Alpha"}, {ID: "beta", Label: "Beta"},
		}},
		{Key: "second", Label: "Second", Items: []mygoui.SidebarItemOptions{
			{ID: "gamma", Label: "Gamma"},
		}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed, Groups: groups,
		})
	}, 400, 400)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	groups[0].Items[0], groups[0].Items[1] = groups[0].Items[1], groups[0].Items[0]
	groups[0], groups[1] = groups[1], groups[0]
	tt.Frame()
	if !tt.Focused("Beta") {
		t.Fatal("focus was lost or moved when keyed items reordered")
	}
	tt.Key(0, ui.KeyEnter)
	if selected != "beta" {
		t.Fatalf("reordered focus activated %q", selected)
	}
}

func TestSidebarCollapsePreservesFocusedItem(t *testing.T) {
	selected, collapsed := "", false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{Label: "Main", Items: []mygoui.SidebarItemOptions{
				{ID: "home", Label: "Home", Icon: lucide.House()},
			}}},
		})
	}, 320, 200)
	if err := tt.Click("Home"); err != nil {
		t.Fatal(err)
	}
	collapsed = true
	tt.Frame()
	if !tt.Focused("Home") {
		t.Fatal("collapsing lost navigation focus")
	}
	tt.Key(0, ui.KeyEnter)
	if selected != "home" {
		t.Fatalf("focused icon selected %q", selected)
	}
	collapsed = false
	tt.Frame()
	if !tt.Focused("Home") {
		t.Fatal("expanding lost navigation focus")
	}
}

func TestSidebarScrollsGroupsAndKeepsFooter(t *testing.T) {
	selected, collapsed := "", false
	items := make([]mygoui.SidebarItemOptions, 12)
	for i := range items {
		items[i] = mygoui.SidebarItemOptions{ID: fmt.Sprint(i), Label: fmt.Sprintf("Page %d", i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{Items: items}},
			Footer: func() { ui.Text(c, "Footer") },
		})
	}, 400, 220)
	footer, ok := tt.Find("Footer")
	if !ok {
		t.Fatal("footer missing")
	}
	tt.Scroll(100, 100, 0, 1000)
	if err := tt.Click("Page 11"); err != nil {
		t.Fatal(err)
	}
	after, ok := tt.Find("Footer")
	if selected != "11" || !ok || after != footer {
		t.Fatalf("scrolling navigation: selected=%q footer=%+v, was %+v", selected, after, footer)
	}
}

func TestSidebarTabRevealsOffscreenItem(t *testing.T) {
	selected, collapsed := "", false
	items := make([]mygoui.SidebarItemOptions, 10)
	for i := range items {
		items[i] = mygoui.SidebarItemOptions{ID: fmt.Sprint(i), Label: fmt.Sprintf("Item %d", i)}
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Sidebar(c, mygoui.SidebarOptions{
			Selected: &selected, Collapsed: &collapsed,
			Groups: []mygoui.SidebarGroupOptions{{Items: items}},
		})
	}, 320, 160)
	for range 11 {
		tt.Key(0, ui.KeyTab)
	}
	last, ok := tt.Find("Item 9")
	if !ok || !tt.Focused("Item 9") || last.Y < 0 || last.Y+last.H > 160 {
		t.Fatalf("offscreen item was not revealed: focused=%v bounds=%+v", tt.Focused("Item 9"), last)
	}
	tt.Key(0, ui.KeyEnter)
	if selected != "9" {
		t.Fatalf("keyboard selected %q", selected)
	}
}
