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
