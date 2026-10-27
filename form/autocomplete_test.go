package form_test

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAutocompleteEditsFiltersAndChoosesWithKeyboard(t *testing.T) {
	var state form.AutocompleteState
	var queries []string
	var changes [][]string
	var focused bool
	items := comboItems()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Placeholder: "Search fruit", Items: items, State: &state,
				OnQueryChange: func(query string) { queries = append(queries, query) },
				OnChange:      func(keys []string) { changes = append(changes, slices.Clone(keys)) },
			})
			focused = state.Trigger.Focused(c) && state.Search.Focused(c)
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Type("b")
	if state.Query != "b" || !state.Open || !focused || !tt.HasText("Beta") || tt.HasText("Alpha") {
		t.Fatalf("query %q open %v focused %v, text %q", state.Query, state.Open, focused, tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.List.SelectedKey != "b" || state.Query != "Beta" || state.Open || !focused {
		t.Fatalf("choice: key %q query %q open %v focused %v", state.List.SelectedKey, state.Query, state.Open, focused)
	}
	if !slices.Equal(queries, []string{"b", "Beta"}) || len(changes) != 1 || !slices.Equal(changes[0], []string{"b"}) {
		t.Fatalf("query callbacks %q, selection callbacks %v", queries, changes)
	}
	tt.Type("x")
	if !state.Open || state.Query != "Betax" {
		t.Fatalf("editing a selection: query %q open %v", state.Query, state.Open)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open || !focused {
		t.Fatalf("escape: open %v focused %v", state.Open, focused)
	}
	tt.Type("y")
	if !state.Open || state.Query != "Betaxy" {
		t.Fatalf("typing after escape: query %q open %v", state.Query, state.Open)
	}
	state.Open = false
	tt.Frame()
	if state.Open || !focused {
		t.Fatalf("caller close: open %v focused %v", state.Open, focused)
	}
}

func TestAutocompleteMultipleSelectionKeepsEditorAndOptions(t *testing.T) {
	var state form.AutocompleteState
	var changes [][]string
	var focused bool
	items := comboItems()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Placeholder: "Search fruit", Items: items, State: &state,
				Selection: list.SelectionMultiple,
				OnChange:  func(keys []string) { changes = append(changes, slices.Clone(keys)) },
			})
			focused = state.Trigger.Focused(c)
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Al")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !state.Open || !focused || state.Query != "" || !slices.Equal(state.SelectedKeys(), []string{"a"}) {
		t.Fatalf("first selection: keys %v query %q open %v focused %v", state.SelectedKeys(), state.Query, state.Open, focused)
	}
	tt.Type("Be")
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !state.Open || !focused || state.Query != "" || len(state.SelectedKeys()) != 2 {
		t.Fatalf("second selection: keys %v query %q open %v focused %v", state.SelectedKeys(), state.Query, state.Open, focused)
	}
	if len(changes) != 2 || !slices.Equal(changes[0], []string{"a"}) || !slices.Equal(changes[1], []string{"a", "b"}) {
		t.Fatalf("selection callbacks %v", changes)
	}
	if items[0].Options.Indicator != nil || items[1].Options.Indicator != nil {
		t.Fatal("autocomplete changed caller-owned item options")
	}
	if err := tt.Click("Remove Alpha"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !focused || !slices.Equal(state.SelectedKeys(), []string{"b"}) || len(changes) != 3 || !slices.Equal(changes[2], []string{"b"}) {
		t.Fatalf("remove: keys %v focused %v callbacks %v", state.SelectedKeys(), focused, changes)
	}
}

func TestAutocompleteDisabledAndClear(t *testing.T) {
	state := form.AutocompleteState{}
	state.List.SelectedKey = "a"
	var queries []string
	var changes [][]string
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Placeholder: "Search fruit", Items: comboItems(), State: &state,
				Clearable: true, Disabled: disabled,
				OnQueryChange: func(query string) { queries = append(queries, query) },
				OnChange:      func(keys []string) { changes = append(changes, slices.Clone(keys)) },
			})
		})
	}, 360, 480)
	if state.Query != "Alpha" {
		t.Fatalf("initial selected label = %q", state.Query)
	}
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if state.Query != "" || state.List.SelectedKey != "" || len(changes) != 1 || len(changes[0]) != 0 || !slices.Equal(queries, []string{""}) {
		t.Fatalf("clear: query %q key %q callbacks %v %q", state.Query, state.List.SelectedKey, changes, queries)
	}
	disabled = true
	tt.Frame()
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Beta")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if state.Query != "" || state.Open || state.List.SelectedKey != "" || len(changes) != 1 || len(queries) != 1 {
		t.Fatalf("disabled: query %q open %v key %q callbacks %v %q", state.Query, state.Open, state.List.SelectedKey, changes, queries)
	}
}

func TestAutocompleteArrowsSkipDisabledAndPreserveFocus(t *testing.T) {
	var state form.AutocompleteState
	var focused bool
	items := []list.ListBoxItem{
		{Key: "a", Label: "Alpha"},
		{Key: "disabled", Label: "Unavailable", Disabled: true},
		{Key: "b", Label: "Beta"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Items: items, State: &state,
			})
			ui.Button(c, "Other")
			focused = state.Trigger.Focused(c)
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.List.SelectedKey != "b" || state.Query != "Beta" || state.Open || !focused {
		t.Fatalf("arrow choice: key %q query %q open %v focused %v", state.List.SelectedKey, state.Query, state.Open, focused)
	}
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	if !state.Open {
		t.Fatal("arrow did not reopen the popup")
	}
	if err := tt.Click("Other"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state.Open {
		t.Fatal("popup remained open after focus moved away")
	}
}

func TestAutocompleteDisabledWhileOpenStopsSelection(t *testing.T) {
	state := form.AutocompleteState{Open: true}
	disabled := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Autocomplete(c, form.AutocompleteOptions{
			Label: "Fruit", Items: comboItems(), State: &state, Disabled: disabled,
		})
	}, 360, 480)
	tt.Frame()
	if state.Open || tt.HasText("Alpha") || tt.HasText("Beta") {
		t.Fatalf("disabled popup: open %v text %q", state.Open, tt.Texts())
	}
}

func TestAutocompleteRefreshesChangedItemAndFilter(t *testing.T) {
	var state form.AutocompleteState
	items := []list.ListBoxItem{{Key: "a", Label: "Alpha"}}
	accept := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Items: items, State: &state,
				Filter: func(string, *list.ListBoxItem) bool { return accept },
			})
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Type("a")
	if !tt.HasText("Alpha") {
		t.Fatalf("initial option missing: %q", tt.Texts())
	}
	items[0].Label = "Almond"
	tt.Frame()
	if !tt.HasText("Almond") || tt.HasText("Alpha") {
		t.Fatalf("updated option not shown: %q", tt.Texts())
	}
	accept = false
	tt.Frame()
	if tt.HasText("Almond") {
		t.Fatalf("updated filter still shows option: %q", tt.Texts())
	}
}

func TestAutocompleteSelectionCallbackCanEditQuery(t *testing.T) {
	var state form.AutocompleteState
	var queries []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Autocomplete(c, form.AutocompleteOptions{
			Label: "Fruit", Items: comboItems(), State: &state,
			OnQueryChange: func(query string) { queries = append(queries, query) },
			OnChange:      func([]string) { state.Query = "custom" },
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Al")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.Query != "custom" || !slices.Equal(queries, []string{"Al", "Alpha"}) {
		t.Fatalf("query %q, callbacks %q", state.Query, queries)
	}
	tt.Type("!")
	if state.Query != "custom!" {
		t.Fatalf("typing after callback = %q", state.Query)
	}
}

func TestAutocompleteQueryCallbackCanEditValue(t *testing.T) {
	var state form.AutocompleteState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Autocomplete(c, form.AutocompleteOptions{
			Label: "Fruit", Items: comboItems(), State: &state,
			OnQueryChange: func(query string) {
				if query == "Alpha" {
					state.Query = "custom"
				}
			},
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Al")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.Query != "custom" {
		t.Fatalf("query after callback = %q", state.Query)
	}
}

func TestAutocompleteReflectsCallerSelection(t *testing.T) {
	var state form.AutocompleteState
	state.List.SelectedKey = "a"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Autocomplete(c, form.AutocompleteOptions{Label: "Fruit", Items: comboItems(), State: &state})
	}, 360, 480)
	if state.Query != "Alpha" {
		t.Fatalf("initial query = %q", state.Query)
	}
	state.List.SelectedKey = "b"
	tt.Frame()
	if state.Query != "Beta" {
		t.Fatalf("updated selection query = %q", state.Query)
	}
	state.List.SelectedKey = ""
	tt.Frame()
	if state.Query != "" {
		t.Fatalf("cleared selection query = %q", state.Query)
	}
	state.Query = "custom"
	state.List.SelectedKey = "a"
	tt.Frame()
	if state.Query != "custom" {
		t.Fatalf("caller query overwritten: %q", state.Query)
	}
}

func TestAutocompleteMultipleChipCanBeRemovedWithKeyboard(t *testing.T) {
	var state form.AutocompleteState
	state.List.SetSelected(list.SelectionMultiple, "a", true)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Items: comboItems(), State: &state, Selection: list.SelectionMultiple,
			})
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Shift, ui.KeyTab)
	if !tt.Focused("Remove Alpha") {
		t.Fatal("chip removal is not keyboard focusable")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(state.SelectedKeys()) != 0 || !tt.Focused("Fruit") {
		t.Fatalf("keyboard removal: keys %v focused fruit %v", state.SelectedKeys(), tt.Focused("Fruit"))
	}
}

func TestAutocompleteClearIsKeyboardAccessible(t *testing.T) {
	var state form.AutocompleteState
	state.List.SelectedKey = "a"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Items: comboItems(), State: &state, Clearable: true,
			})
		})
	}, 360, 480)
	if err := tt.Click("Fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Clear") {
		t.Fatal("clear is not keyboard focusable")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.List.SelectedKey != "" || state.Query != "" || !tt.Focused("Fruit") {
		t.Fatalf("keyboard clear: key %q query %q focused %v", state.List.SelectedKey, state.Query, tt.Focused("Fruit"))
	}
}

func TestAutocompleteMultipleChipsWrapInNarrowField(t *testing.T) {
	var state form.AutocompleteState
	items := []list.ListBoxItem{
		{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Beta"},
		{Key: "c", Label: "Charlie"}, {Key: "d", Label: "Delta"},
	}
	for _, item := range items {
		state.List.SetSelected(list.SelectionMultiple, item.Key, true)
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(16).Children(func() {
			form.Autocomplete(c, form.AutocompleteOptions{
				Label: "Fruit", Items: items, State: &state, Selection: list.SelectionMultiple, FullWidth: true,
			})
		})
	}, 220, 480)
	first, ok := tt.Find("Remove Alpha")
	if !ok {
		t.Fatal("first chip is missing")
	}
	last, ok := tt.Find("Remove Delta")
	if !ok || last.X+last.W > 204 || last.Y <= first.Y {
		t.Fatalf("chips did not wrap inside the field: first %v, last %v", first, last)
	}
	if err := tt.Click("Remove Delta"); err != nil {
		t.Fatal(err)
	}
	if len(state.SelectedKeys()) != 3 {
		t.Fatalf("wrapped chip did not remove: %v", state.SelectedKeys())
	}
}
