package form_test

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSelectSingleChoiceAndClear(t *testing.T) {
	state := &form.SelectState{}
	var changed []string
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Select(c, form.SelectOptions{
			Label: "Fruit", Placeholder: "Choose a fruit", Clearable: true,
			Items: comboItems(), State: state,
			OnChange: func(keys []string) { changed = keys },
		})
	}, 360, 520)
	if err := tt.Click("Choose a fruit"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	if state.List.SelectedKey != "a" || !slices.Equal(changed, []string{"a"}) {
		t.Fatalf("select key=%q changed=%v", state.List.SelectedKey, changed)
	}
	if !tt.HasText("Alpha") {
		t.Fatal("trigger did not show the chosen value")
	}
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if state.List.SelectedKey != "" || !slices.Equal(changed, []string{}) {
		t.Fatalf("clear key=%q changed=%v", state.List.SelectedKey, changed)
	}
}

func TestSelectMultipleKeepsPopupOpen(t *testing.T) {
	state := &form.SelectState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		form.Select(c, form.SelectOptions{
			Label: "Fruit", Placeholder: "Choose fruit", Items: comboItems(),
			Selection: list.SelectionMultiple, State: state,
		})
	}, 360, 520)
	if err := tt.Click("Choose fruit"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if got := state.SelectedKeys(); len(got) != 2 {
		t.Fatalf("multiple keys = %v", got)
	}
	if !state.Open {
		t.Fatal("multiple select closed after one choice")
	}
}
