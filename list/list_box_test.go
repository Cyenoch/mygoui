package list_test

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func items() []list.ListBoxItem {
	return []list.ListBoxItem{
		{Key: "apple", Label: "Apple"},
		{Key: "banana", Label: "Banana", Description: "Yellow"},
		{Key: "cherry", Label: "Cherry", Disabled: true},
	}
}

func TestListBoxSingleSelectionAndAction(t *testing.T) {
	var state list.ListBoxState
	var chosen []string
	acted := ""
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Fruit", Items: items(), Selection: list.SelectionSingle,
			State: &state, OnSelectionChange: func(keys []string) { chosen = keys },
			OnAction: func(key string) { acted = key },
		})
	}, 320, 300)
	if err := tt.Click("Banana"); err != nil {
		t.Fatal(err)
	}
	if state.SelectedKey != "banana" || !slices.Equal(chosen, []string{"banana"}) {
		t.Fatalf("single selection = %q, callback %v", state.SelectedKey, chosen)
	}
	if acted != "banana" {
		t.Fatalf("action = %q", acted)
	}
	if err := tt.Click("Cherry"); err != nil {
		t.Fatal(err)
	}
	if state.SelectedKey != "banana" {
		t.Fatalf("disabled row changed selection to %q", state.SelectedKey)
	}
}

func TestListBoxMultipleSections(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Food", Selection: list.SelectionMultiple, State: &state,
			Sections: []list.ListBoxSection{{
				Label: "Fruit", Items: items()[:2],
			}, {
				Label: "Vegetables", Items: []list.ListBoxItem{{Key: "kale", Label: "Kale"}},
			}},
		})
	}, 320, 360)
	if !tt.HasText("Fruit") || !tt.HasText("Vegetables") {
		t.Fatal("section headers missing")
	}
	if err := tt.Click("Apple"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Kale"); err != nil {
		t.Fatal(err)
	}
	keys := state.SelectedKeys(list.SelectionMultiple)
	if len(keys) != 2 || !state.IsSelected("apple") || !state.IsSelected("kale") {
		t.Fatalf("multiple selection = %v", keys)
	}
	if err := tt.Click("Apple"); err != nil {
		t.Fatal(err)
	}
	if state.IsSelected("apple") {
		t.Fatal("second click did not remove the key")
	}
}

func TestListBoxKeyboardNavigationSkipsDisabled(t *testing.T) {
	var state list.ListBoxState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Fruit", Items: items(), Selection: list.SelectionSingle, State: &state,
		})
	}, 320, 300)
	tt.Click("Apple")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "banana" {
		t.Fatalf("keyboard selection = %q, want banana (disabled cherry skipped)", state.SelectedKey)
	}
}
