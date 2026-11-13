package form

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSelectEscapeReturnsFocusToTrigger(t *testing.T) {
	state := &SelectState{}
	focused := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		Select(c, SelectOptions{Label: "Fruit", Placeholder: "Pick", Items: []list.ListBoxItem{{Key: "a", Label: "Alpha"}}, State: state})
		focused = state.Trigger.Focused(c)
	}, 360, 480)
	formTextClick(t, tt, "Pick")
	tt.Frame()
	if !state.Open {
		t.Fatal("the trigger did not open the popup")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open || !focused {
		t.Fatalf("after Escape: open %v, trigger focused %v", state.Open, focused)
	}
}

func TestSelectClosedTypeaheadCyclesAndMatchesPrefix(t *testing.T) {
	state := &SelectState{}
	items := []list.ListBoxItem{
		{Key: "apple", Label: "Apple"}, {Key: "avocado", Label: "Avocado"},
		{Key: "banana", Label: "Banana"}, {Key: "blueberry", Label: "Blueberry"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		var before string
		Input(c, &before, InputOptions{Label: "Before"})
		Select(c, SelectOptions{Label: "Fruit", Placeholder: "Pick", Items: items, State: state})
	}, 360, 480)
	formTextClick(t, tt, "Before")
	tt.Key(0, ui.KeyTab)
	for _, want := range []string{"apple", "avocado", "apple"} {
		tt.Key(0, ui.KeyA)
		if state.List.SelectedKey != want {
			t.Fatalf("A chose %q, want %q", state.List.SelectedKey, want)
		}
	}
	tt.Advance(2 * time.Second)
	tt.Key(0, ui.KeyB)
	tt.Key(0, ui.KeyL)
	if state.List.SelectedKey != "blueberry" || state.Open {
		t.Fatalf("B L chose %q, open %v", state.List.SelectedKey, state.Open)
	}
}

func TestComboBoxAltArrowsOpenAndClose(t *testing.T) {
	state := &ComboBoxState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ComboBox(c, ComboBoxOptions{Label: "Letters", Items: []list.ListBoxItem{{Key: "a", Label: "Alpha"}}, State: state})
	}, 360, 480)
	formTextClick(t, tt, "Letters")
	tt.Type("Al")
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Frame()
	if state.Open || state.Text != "Al" {
		t.Fatalf("Alt+Up: open %v, text %q", state.Open, state.Text)
	}
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Frame()
	if !state.Open || state.Text != "Al" {
		t.Fatalf("Alt+Down: open %v, text %q", state.Open, state.Text)
	}
}

func TestComboBoxEscapeOnClosedPopupClears(t *testing.T) {
	state := &ComboBoxState{}
	escapes, changes := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Children(func() {
			ComboBox(c, ComboBoxOptions{
				Label: "Letters", Clearable: true, State: state,
				Items:    []list.ListBoxItem{{Key: "a", Label: "Alpha"}},
				OnChange: func([]string) { changes++ },
			})
		}).OnShortcut(0, ui.KeyEscape, func() { escapes++ })
	}, 360, 480)
	formTextClick(t, tt, "Letters")
	tt.Type("Al")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.List.SelectedKey != "a" || state.Open {
		t.Fatalf("Enter chose %q, open %v", state.List.SelectedKey, state.Open)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Text != "" || state.List.SelectedKey != "" || changes != 2 || escapes != 0 {
		t.Fatalf("Escape: text %q, key %q, changes %d, ancestor %d", state.Text, state.List.SelectedKey, changes, escapes)
	}
	tt.Key(0, ui.KeyEscape)
	if escapes != 1 {
		t.Fatalf("Escape on the empty combo box reached the ancestor %d times, want 1", escapes)
	}
}

func TestAutocompleteAltArrowsAndEscapeOnClosedPopup(t *testing.T) {
	state := &AutocompleteState{}
	escapes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Children(func() {
			Autocomplete(c, AutocompleteOptions{
				Label: "City", Clearable: true, State: state,
				Items: []list.ListBoxItem{{Key: "a", Label: "Amsterdam"}},
			})
		}).OnShortcut(0, ui.KeyEscape, func() { escapes++ })
	}, 360, 480)
	formTextClick(t, tt, "City")
	tt.Type("Am")
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Frame()
	if state.Open {
		t.Fatal("Alt+Up left the popup open")
	}
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Frame()
	if !state.Open {
		t.Fatal("Alt+Down did not open the popup")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open || state.Query != "Am" || escapes != 0 {
		t.Fatalf("Escape on the open popup: open %v, query %q, ancestor %d", state.Open, state.Query, escapes)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Query != "" || escapes != 0 {
		t.Fatalf("Escape on the closed popup: query %q, ancestor %d", state.Query, escapes)
	}
	tt.Key(0, ui.KeyEscape)
	if escapes != 1 {
		t.Fatalf("Escape on the empty field reached the ancestor %d times, want 1", escapes)
	}
}

// TestTextPopupsPointTheirInputAtTheHighlightedOption pins the active
// descendant: the input keeps the focus, so it names the option the arrows
// are on, and assistive technology follows them.
func TestTextPopupsPointTheirInputAtTheHighlightedOption(t *testing.T) {
	items := []list.ListBoxItem{{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Beta"}}
	combo := &ComboBoxState{}
	auto := &AutocompleteState{}
	for name, build := range map[string]func(c *ui.Context){
		"ComboBox": func(c *ui.Context) {
			ComboBox(c, ComboBoxOptions{Label: "Letters", Items: items, State: combo})
		},
		"Autocomplete": func(c *ui.Context) {
			Autocomplete(c, AutocompleteOptions{Label: "Letters", Items: items, State: auto})
		},
	} {
		t.Run(name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				build(c)
			}, 360, 480)
			formTextClick(t, tt, "Letters")
			tt.Key(0, ui.KeyDown)
			tt.Frame()
			if beta, _ := tt.AccessByRole(ui.RoleOption, "Beta"); !beta.Focused {
				t.Fatalf("Down from the input: Beta is not the active descendant (%+v)", beta)
			}
			if !tt.Focused("Letters") {
				t.Fatal("the input lost the focus to the list")
			}
		})
	}
}
