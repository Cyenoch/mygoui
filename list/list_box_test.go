package list_test

import (
	"slices"
	"testing"
	"time"

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

func TestListBoxLoadingFirstLoadShowsNoEmptyMessage(t *testing.T) {
	var state list.ListBoxState
	loading := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{Label: "Fruit", State: &state, Empty: "No fruit", Loading: loading})
	}, 320, 300)
	if tt.HasText("No fruit") {
		t.Fatalf("empty message shown while loading: %v", tt.Texts())
	}
	loading = false
	tt.Frame()
	if !tt.HasText("No fruit") {
		t.Fatalf("empty message missing after loading: %v", tt.Texts())
	}
}

func TestListBoxLoadingRefreshKeepsRowsAndIgnoresInput(t *testing.T) {
	var state list.ListBoxState
	loading := true
	acted := ""
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Fruit", Items: items(), State: &state, Loading: loading,
			OnAction: func(key string) { acted = key },
		})
	}, 320, 300)
	tt.Advance(400 * time.Millisecond)
	if !tt.HasText("Apple") {
		t.Fatalf("rows left while refreshing: %v", tt.Texts())
	}
	_ = tt.Click("Apple")
	if acted != "" {
		t.Fatalf("refreshing list took a click: %q", acted)
	}
	loading = false
	tt.Frame()
	tt.Advance(400 * time.Millisecond)
	if err := tt.Click("Apple"); err != nil || acted != "apple" {
		t.Fatalf("list stayed inert after loading: %v %q", err, acted)
	}
}
