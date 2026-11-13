package navigation_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSplitButtonActsAndOffersAlternatives(t *testing.T) {
	var saves int
	var chosen string
	actionDisabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.SplitButton(c, navigation.SplitButtonOptions{
			Label: "Save", MenuLabel: "Save options", ActionDisabled: actionDisabled,
			OnClick: func() { saves++ },
			Menu: navigation.MenuOptions{
				Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
					{Key: "save-as", Label: "Save As…"},
					{Key: "save-copy", Label: "Save a Copy…"},
				}}},
				OnSelect: func(key string) { chosen = key },
			},
		})
	}, 360, 240)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if saves != 1 || tt.HasText("Save As…") {
		t.Fatalf("action ran %d times, menu shown %v", saves, tt.HasText("Save As…"))
	}
	if err := tt.Click("Save options"); err != nil {
		t.Fatal(err)
	}
	arrow, ok := tt.AccessByRole(ui.RoleMenuButton, "Save options")
	if !ok || !arrow.Expanded || !tt.Focused("Save As…") {
		t.Fatalf("arrow: %+v, found %v, first item focused %v", arrow, ok, tt.Focused("Save As…"))
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen != "save-copy" || tt.HasText("Save As…") || saves != 1 {
		t.Fatalf("choice %q, menu shown %v, action ran %d times", chosen, tt.HasText("Save As…"), saves)
	}

	actionDisabled = true
	tt.Frame()
	if err := tt.Click("Save"); err == nil && saves != 1 {
		t.Fatal("a disabled action ran")
	}
	if err := tt.Click("Save options"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Save As…") {
		t.Fatal("the menu does not open while the action alone is disabled")
	}
}

// TestSplitButtonLoadingStopsActionAndKeepsMenu pins the pending contract:
// the action takes no press while the arrow still opens the menu.
func TestSplitButtonLoadingStopsActionAndKeepsMenu(t *testing.T) {
	var saves int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.SplitButton(c, navigation.SplitButtonOptions{
			Label: "Save", MenuLabel: "Save options", Loading: true,
			OnClick: func() { saves++ },
			Menu: navigation.MenuOptions{
				Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
					{Key: "save-as", Label: "Save As…"},
				}}},
			},
		})
	}, 360, 240)
	tt.Advance(200 * time.Millisecond)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if saves != 0 {
		t.Fatalf("a pending action ran %d times", saves)
	}
	if err := tt.Click("Save options"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Save As…") {
		t.Fatal("the menu does not open while the action is pending")
	}
}

func TestSplitButtonRemembersTheAlternativeChosenLast(t *testing.T) {
	selected := "save-copy"
	var chosen string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		navigation.SplitButton(c, navigation.SplitButtonOptions{
			Label: "Save", MenuLabel: "Save options", Selected: &selected,
			Menu: navigation.MenuOptions{
				Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
					{Key: "save-as", Label: "Save As…"},
					{Key: "save-copy", Label: "Save a Copy…"},
				}}},
				OnSelect: func(key string) { chosen = key },
			},
		})
	}, 360, 240)
	if err := tt.Click("Save options"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	if info, ok := tt.AccessByRole(ui.RoleMenuItemRadio, "Save a Copy…"); !ok || !info.Checked {
		t.Fatalf("the remembered item is not marked: %+v, found %v", info, ok)
	}
	if err := tt.Click("Save As…"); err != nil {
		t.Fatal(err)
	}
	if selected != "save-as" || chosen != "save-as" {
		t.Fatalf("selected %q, menu reported %q", selected, chosen)
	}
}
