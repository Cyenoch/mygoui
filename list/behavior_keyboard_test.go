package list_test

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// listBoxBetweenButtons places a ListBox between two buttons, so Tab can
// enter and leave it.
func listBoxBetweenButtons(o *list.ListBoxOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Children(func() {
			ui.Button(c, "Before")
			list.ListBox(c, *o)
			ui.Button(c, "After")
		})
	}, 320, 320)
}

// An action-only list answers Space and Enter on the highlighted row with
// OnAction and chooses nothing; the list is one Tab stop.
func TestListBoxKeysActivateTheHighlightedRow(t *testing.T) {
	var state list.ListBoxState
	var acted []string
	o := list.ListBoxOptions{
		Label: "Fruit", Items: items(), State: &state,
		OnAction: func(key string) { acted = append(acted, key) },
	}
	tt := listBoxBetweenButtons(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Fruit") || state.HighlightKey() != "apple" {
		t.Fatalf("Tab into the list: focused %v, highlight %q", tt.Focused("Fruit"), state.HighlightKey())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if want := []string{"banana", "apple"}; !slices.Equal(acted, want) {
		t.Fatalf("actions = %v, want %v", acted, want)
	}
	if state.SelectionCount(list.SelectionSingle)+state.SelectionCount(list.SelectionMultiple) != 0 {
		t.Fatalf("an action-only list chose %q %v", state.SelectedKey, state.SelectedKeys(list.SelectionMultiple))
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab did not leave the list for the next control")
	}
}

// Space toggles the highlighted row of a multiple selection.
func TestListBoxSpaceTogglesInMultipleSelection(t *testing.T) {
	var state list.ListBoxState
	var reported []string
	o := list.ListBoxOptions{
		Label: "Fruit", Items: items(), Selection: list.SelectionMultiple, State: &state,
		OnSelectionChange: func(keys []string) { reported = keys },
	}
	tt := listBoxBetweenButtons(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if !state.IsSelected("apple") || !state.IsSelected("banana") || len(reported) != 2 {
		t.Fatalf("after two Spaces: selected %v, reported %v", state.SelectedKeys(list.SelectionMultiple), reported)
	}
	tt.Key(0, ui.KeySpace)
	if state.IsSelected("banana") || !slices.Equal(reported, []string{"apple"}) {
		t.Fatalf("a second Space kept banana: reported %v", reported)
	}
}

// Typing the same letter again moves on to the next enabled row starting
// with it, wrapping at the end.
func TestListBoxTypeaheadCyclesThroughMatchingRows(t *testing.T) {
	var state list.ListBoxState
	o := list.ListBoxOptions{
		Label: "Fruit", State: &state,
		Items: []list.ListBoxItem{
			{Key: "apple", Label: "Apple"}, {Key: "banana", Label: "Banana"},
			{Key: "apricot", Label: "Apricot", Disabled: true}, {Key: "avocado", Label: "Avocado"},
		},
	}
	tt := listBoxBetweenButtons(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	var got []string
	for range 3 {
		tt.Key(0, ui.KeyA)
		got = append(got, state.HighlightKey())
	}
	if want := []string{"avocado", "apple", "avocado"}; !slices.Equal(got, want) {
		t.Fatalf("typing a three times highlighted %v, want %v", got, want)
	}
}
