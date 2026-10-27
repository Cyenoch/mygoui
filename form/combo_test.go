package form_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func comboItems() []list.ListBoxItem {
	return []list.ListBoxItem{{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Beta"}}
}

func TestComboBoxPointerChoosesOnce(t *testing.T) {
	state := form.ComboBoxState{Open: true}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.ComboBox(c, form.ComboBoxOptions{
			Label: "Letters", Items: comboItems(), State: &state,
			Selection: list.SelectionMultiple,
		})
	}, 360, 420)
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	keys := state.SelectedKeys()
	if len(keys) != 1 || keys[0] != "a" {
		t.Fatalf("multiple ComboBox pointer selection = %v, want one key", keys)
	}
}

func TestAutocompletePointerChoosesOnce(t *testing.T) {
	state := form.AutocompleteState{Open: true}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Autocomplete(c, form.AutocompleteOptions{
			Label: "Letters", Items: comboItems(), State: &state,
			Selection: list.SelectionMultiple,
		})
	}, 360, 480)
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	keys := state.SelectedKeys()
	if len(keys) != 1 || keys[0] != "a" {
		t.Fatalf("multiple Autocomplete pointer selection = %v, want one key", keys)
	}
}

func TestComboBoxEscapeClosesAndStaysClosed(t *testing.T) {
	state := form.ComboBoxState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.ComboBox(c, form.ComboBoxOptions{
				Label: "Letters", Items: comboItems(), State: &state,
			})
		})
	}, 360, 480)

	// The label focuses the editor, which opens the popup while focused.
	if err := tt.Click("Letters"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !state.Open {
		t.Fatal("label click did not focus the editor")
	}
	tt.Type("Al")
	if state.Text != "Al" {
		t.Fatalf("typing after label click committed %q", state.Text)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open {
		t.Fatal("Escape reopened the popup it closed")
	}
	// The editor keeps the focus; typing opens the popup again.
	tt.Type("p")
	if state.Text != "Alp" || !state.Open {
		t.Fatalf("after Escape: text %q, open %v", state.Text, state.Open)
	}
}

func TestComboBoxKeyboardChoosesAndKeepsFocus(t *testing.T) {
	var state form.ComboBoxState
	var changed []string
	var focused bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.ComboBox(c, form.ComboBoxOptions{
				Label: "Letters", Items: comboItems(), State: &state,
				OnChange: func(keys []string) { changed = keys },
			})
			focused = state.Input.Focused(c)
		})
	}, 360, 480)
	if err := tt.Click("Letters"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Be")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.Text != "Beta" || state.Open || !focused || len(changed) != 1 || changed[0] != "b" {
		t.Fatalf("keyboard choice: text %q open %v focused %v callback %v", state.Text, state.Open, focused, changed)
	}
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	if !state.Open || !tt.HasText("Alpha") || !tt.HasText("Beta") {
		t.Fatalf("reopened choices: open %v text %q", state.Open, tt.Texts())
	}
}

func TestComboBoxSelectionCallbackCanEditText(t *testing.T) {
	var state form.ComboBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.ComboBox(c, form.ComboBoxOptions{
			Label: "Letters", Items: comboItems(), State: &state,
			OnChange: func([]string) { state.Text = "custom" },
		})
	}, 360, 480)
	if err := tt.Click("Letters"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Al")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.Text != "custom" {
		t.Fatalf("text after callback = %q", state.Text)
	}
}

func TestEditablePopupKeyboardFiltersGroupedOptions(t *testing.T) {
	for _, kind := range []string{"autocomplete", "combobox"} {
		t.Run(kind, func(t *testing.T) {
			var autocomplete form.AutocompleteState
			var combo form.ComboBoxState
			var changed []string
			items := []list.ListBoxItem{
				{Key: "a", Label: "Alpha"},
				{Key: "disabled", Label: "Bee", Disabled: true},
			}
			sections := []list.ListBoxSection{{
				Label: "Group",
				Items: []list.ListBoxItem{{Key: "b", Label: "Beta"}, {Key: "c", Label: "Bravo"}},
			}}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				ui.Column(c).Fill().Padding(16).Children(func() {
					onChange := func(keys []string) { changed = keys }
					if kind == "autocomplete" {
						form.Autocomplete(c, form.AutocompleteOptions{
							Label: "Fruit", State: &autocomplete, Items: items, Sections: sections, OnChange: onChange,
						})
					} else {
						form.ComboBox(c, form.ComboBoxOptions{
							Label: "Fruit", State: &combo, Items: items, Sections: sections, OnChange: onChange,
						})
					}
				})
			}, 360, 480)
			if err := tt.Click("Fruit"); err != nil {
				t.Fatal(err)
			}
			tt.Type("b")
			tt.Key(0, ui.KeyDown)
			tt.Key(0, ui.KeyEnter)
			if len(changed) != 1 || changed[0] != "c" {
				t.Fatalf("filtered keyboard selection = %v, want Bravo", changed)
			}
		})
	}
}
