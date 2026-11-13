package navigation_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestSplitButtonEscapeReturnsFocusToTheArrow protects dismissal of the split
// button's menu: Escape closes it and the keyboard continues from the arrow
// button that opened it.
func TestSplitButtonEscapeReturnsFocusToTheArrow(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		navigation.SplitButton(c, navigation.SplitButtonOptions{Label: "Save", MenuLabel: "Save options", Open: &open,
			Menu: navigation.MenuOptions{Sections: []navigation.MenuSectionOptions{{
				Items: []navigation.MenuItemOptions{{Label: "Save as"}, {Label: "Save a copy"}},
			}}}})
	}, 480, 360)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !open || !tt.Focused("Save as") {
		t.Fatal("Enter on the arrow did not open the menu at its first item")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || tt.HasText("Save a copy") {
		t.Fatal("Escape did not close the menu")
	}
	if a, ok := tt.AccessByRole(ui.RoleMenuButton, "Save options"); !ok || !a.Focused || a.Expanded {
		t.Fatalf("the arrow did not take the focus back collapsed: %+v", a)
	}
}

// focusedName returns the accessible name of the element holding the focus,
// "" when none does.
func focusedName(tt *ui.Tester) string {
	for _, a := range tt.AccessAll() {
		if a.Focused {
			return a.Name
		}
	}
	return ""
}

// menuButtons builds a dropdown and a split button between two buttons.
func menuButtons(dropdown, split *bool) func(*ui.Context) {
	items := navigation.MenuOptions{Sections: []navigation.MenuSectionOptions{{
		Items: []navigation.MenuItemOptions{{Label: "First"}, {Label: "Middle"}, {Label: "Last"}},
	}}}
	return func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Gap(8).Children(func() {
			ui.Button(c, "Before")
			navigation.Dropdown(c, navigation.DropdownOptions{Label: "Actions", Open: dropdown, Menu: items})
			navigation.SplitButton(c, navigation.SplitButtonOptions{Label: "Save", MenuLabel: "Save options", Open: split, Menu: items})
		})
	}
}

// TestMenuButtonsTakeTheFocusBackAsTheirMenuCloses protects the menu button
// pattern: choosing an item or Escape closes the menu and the keyboard
// continues from the button that opened it, while a press outside leaves
// the focus where it landed.
func TestMenuButtonsTakeTheFocusBackAsTheirMenuCloses(t *testing.T) {
	var dropdown, split bool
	tt := ui.NewTester(menuButtons(&dropdown, &split), 480, 480)
	for _, button := range []struct {
		name string
		open *bool
		tabs int
	}{{"Actions", &dropdown, 1}, {"Save options", &split, 3}} {
		for _, close := range []ui.Key{ui.KeyEnter, ui.KeyEscape} {
			if err := tt.Click("Before"); err != nil {
				t.Fatal(err)
			}
			for range button.tabs {
				tt.Key(0, ui.KeyTab)
			}
			tt.Key(0, ui.KeyEnter)
			if !*button.open || focusedName(tt) != "First" {
				t.Fatalf("%s: Enter did not open the menu at its first item", button.name)
			}
			tt.Key(0, close)
			tt.Frame()
			if *button.open {
				t.Fatalf("%s: key %v left the menu open", button.name, close)
			}
			if got := focusedName(tt); got != button.name {
				t.Fatalf("%s: after key %v the focus is on %q", button.name, close, got)
			}
		}
		*button.open = true
		tt.Frame()
		if err := tt.Click("Before"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		if got := focusedName(tt); *button.open || got != "Before" {
			t.Fatalf("%s: a press outside left open=%v and the focus on %q", button.name, *button.open, got)
		}
	}
}

// TestMenuButtonsOpenFromTheArrows protects the menu button pattern's arrow
// keys: Down opens the menu at its first item and Up at its last.
func TestMenuButtonsOpenFromTheArrows(t *testing.T) {
	var dropdown, split bool
	tt := ui.NewTester(menuButtons(&dropdown, &split), 480, 480)
	for _, button := range []struct {
		name string
		open *bool
		tabs int
	}{{"Actions", &dropdown, 1}, {"Save options", &split, 3}} {
		for _, arrow := range []struct {
			key  ui.Key
			want string
		}{{ui.KeyDown, "First"}, {ui.KeyUp, "Last"}} {
			if err := tt.Click("Before"); err != nil {
				t.Fatal(err)
			}
			for range button.tabs {
				tt.Key(0, ui.KeyTab)
			}
			tt.Key(0, arrow.key)
			if got := focusedName(tt); !*button.open || got != arrow.want {
				t.Fatalf("%s: key %v left open=%v and the focus on %q", button.name, arrow.key, *button.open, got)
			}
			tt.Key(0, ui.KeyEscape)
			tt.Frame()
		}
	}
}
