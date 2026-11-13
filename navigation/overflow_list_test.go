package navigation_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var overflowLabels = []string{"Overview", "Activity", "Settings", "Billing", "Members"}

func overflowTester(collapse navigation.OverflowCollapse, chosen *string, width int) *ui.Tester {
	items := make([]navigation.OverflowItem, len(overflowLabels))
	for i, label := range overflowLabels {
		items[i] = navigation.OverflowItem{Label: label}
	}
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.OverflowList(c, navigation.OverflowListOptions{
			Label: "Sections", Items: items, Collapse: collapse,
			OnSelect: func(key string) { *chosen = key },
		})
	}, width, 240)
}

// overflowShown reports the labels an OverflowList shows as buttons.
func overflowShown(tt *ui.Tester) []string {
	var shown []string
	for _, label := range overflowLabels {
		if _, ok := tt.AccessByRole(ui.RoleButton, label); ok {
			shown = append(shown, label)
		}
	}
	return shown
}

func TestOverflowListCollapsesWhatDoesNotFitIntoItsMenu(t *testing.T) {
	var chosen string
	tt := overflowTester(navigation.OverflowCollapseEnd, &chosen, 640)
	if got := overflowShown(tt); len(got) != len(overflowLabels) {
		t.Fatalf("a wide row shows %v", got)
	}
	if _, ok := tt.AccessByRole(ui.RoleMenuButton, "More"); ok {
		t.Fatal("a row that fits shows a more button")
	}
	tt.SetSize(250, 240)
	shown := overflowShown(tt)
	if len(shown) == 0 || len(shown) == len(overflowLabels) || shown[0] != "Overview" {
		t.Fatalf("a narrow row shows %v", shown)
	}
	more, ok := tt.AccessByRole(ui.RoleMenuButton, "More")
	if !ok || more.Bounds.X+more.Bounds.W > 250 {
		t.Fatalf("more button: %+v, found %v", more, ok)
	}
	if last, _ := tt.AccessByRole(ui.RoleButton, shown[len(shown)-1]); last.Bounds.X+last.Bounds.W > more.Bounds.X {
		t.Fatalf("the last shown item %v runs under the more button %v", last.Bounds, more.Bounds)
	}
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	if _, ok := tt.AccessByRole(ui.RoleMenuItem, "Members"); !ok {
		t.Fatal("the menu does not list a collapsed item")
	}
	if _, ok := tt.AccessByRole(ui.RoleMenuItem, "Overview"); ok {
		t.Fatal("the menu lists an item the row shows")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if chosen != "Members" {
		t.Fatalf("menu choice reported %q", chosen)
	}
	tt.SetSize(640, 240)
	if got := overflowShown(tt); len(got) != len(overflowLabels) {
		t.Fatalf("after widening the row shows %v", got)
	}
}

func TestOverflowListCollapsesFromTheStart(t *testing.T) {
	var chosen string
	tt := overflowTester(navigation.OverflowCollapseStart, &chosen, 250)
	shown := overflowShown(tt)
	if len(shown) == 0 || len(shown) == len(overflowLabels) || shown[len(shown)-1] != "Members" {
		t.Fatalf("a narrow row collapsing from the start shows %v", shown)
	}
	more, _ := tt.AccessByRole(ui.RoleMenuButton, "More")
	first, _ := tt.AccessByRole(ui.RoleButton, shown[0])
	if more.Bounds.X >= first.Bounds.X {
		t.Fatalf("more button at %v is not before the first shown item at %v", more.Bounds, first.Bounds)
	}
	if err := tt.Click(shown[0]); err != nil {
		t.Fatal(err)
	}
	if chosen != shown[0] {
		t.Fatalf("row choice reported %q, want %q", chosen, shown[0])
	}
}

// A narrow row sends the lowest priority into the menu first, never a pinned
// item, and a toggle shows its state in the row and in the menu.
func TestOverflowListKeepsPinnedAndHighPriorityItemsAndShowsToggles(t *testing.T) {
	bold := true
	items := []navigation.OverflowItem{
		{Key: "bold", Label: "Bold", Checked: &bold, Priority: 2},
		{Key: "link", Label: "Link", Priority: 1},
		{Key: "quote", Label: "Quote", Priority: 0},
		{Key: "help", Label: "Help", Pinned: true},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.OverflowList(c, navigation.OverflowListOptions{Label: "Format", Items: items})
	}, 640, 240)
	if info, ok := tt.AccessByRole(ui.RoleToggleButton, "Bold"); !ok || !info.Checked {
		t.Fatalf("a checked item in the row: %+v, found %v", info, ok)
	}
	tt.SetSize(150, 240)
	if _, ok := tt.AccessByRole(ui.RoleButton, "Help"); !ok {
		t.Fatal("a pinned item left the row")
	}
	if _, ok := tt.AccessByRole(ui.RoleButton, "Quote"); ok {
		t.Fatal("the lowest priority stayed while the row had no room")
	}
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if _, ok := tt.AccessByRole(ui.RoleMenuItem, "Quote"); !ok {
		t.Fatal("the lowest priority is not in the menu")
	}
}

// Only a toggle is a checkbox in the menu: a screen reader must not announce
// a plain action as "not checked".
func TestOverflowListMenuChecksOnlyItsToggles(t *testing.T) {
	bold := true
	items := []navigation.OverflowItem{
		{Key: "bold", Label: "Bold", Checked: &bold},
		{Key: "quote", Label: "Quote"},
		{Key: "help", Label: "Help", Pinned: true},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.OverflowList(c, navigation.OverflowListOptions{Label: "Format", Items: items})
	}, 120, 240)
	if err := tt.Click("More"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if info, ok := tt.AccessByRole(ui.RoleMenuItemCheckBox, "Bold"); !ok || !info.Checked {
		t.Fatalf("the toggle in the menu: %+v, found %v", info, ok)
	}
	if _, ok := tt.AccessByRole(ui.RoleMenuItem, "Quote"); !ok {
		t.Fatal("a plain action is not a plain menu item")
	}
}
