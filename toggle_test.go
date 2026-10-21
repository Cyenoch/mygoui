package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

func TestToggleBindsNativeStateAndDisabledInput(t *testing.T) {
	on, disabled, changes := false, false, 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		mygoui.Toggle(c, &on, mygoui.ToggleOptions{
			Label: "Bold", Disabled: disabled,
			OnChange: func(value bool) {
				if value != on {
					t.Fatal("callback observed pre-input value")
				}
				changes++
			},
		})
	}, 320, 160)
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if !on || changes != 1 {
		t.Fatalf("click on=%v changes=%d", on, changes)
	}
	tt.Key(0, myui.KeySpace)
	if on || changes != 2 {
		t.Fatalf("keyboard on=%v changes=%d", on, changes)
	}
	disabled = true
	tt.Frame()
	if err := tt.Click("Bold"); err != nil {
		t.Fatal(err)
	}
	if on || changes != 2 {
		t.Fatalf("disabled on=%v changes=%d", on, changes)
	}
}

func TestToggleGroupSelectionAndRovingFocus(t *testing.T) {
	state := mygoui.ToggleGroupState{}
	mode := mygoui.ToggleSelectionSingle
	items := []mygoui.ToggleGroupItem{
		{Value: "bold", Label: "Bold"},
		{Value: "locked", Label: "Unavailable", Disabled: true},
		{Value: "italic", Label: "Italic"},
	}
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		myui.Column(c).Children(func() {
			mygoui.ToggleGroup(c, &state, mygoui.ToggleGroupOptions{Label: "Formatting", Selection: mode}, items)
			mygoui.Button(c, mygoui.ButtonOptions{Label: "After", OnClick: func() {}})
		})
	}, 420, 240)
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("Bold") {
		t.Fatal("Tab did not enter toggle group")
	}
	tt.Key(0, myui.KeyRight)
	if !tt.Focused("Italic") || len(state.Values) != 0 {
		t.Fatal("arrows must move focus past disabled item without changing selection")
	}
	tt.Key(0, myui.KeySpace)
	if !state.Selected("italic") || len(state.Values) != 1 {
		t.Fatalf("keyboard selection=%v", state.Values)
	}
	tt.Key(0, myui.KeyHome)
	if !tt.Focused("Bold") {
		t.Fatal("Home did not move to first enabled toggle")
	}
	tt.Key(0, myui.KeySpace)
	if !state.Selected("bold") || state.Selected("italic") {
		t.Fatalf("single selection=%v", state.Values)
	}
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("toggle group consumed a second Tab stop")
	}
	mode = mygoui.ToggleSelectionMultiple
	tt.Frame()
	if err := tt.Click("Italic"); err != nil {
		t.Fatal(err)
	}
	if !state.Selected("bold") || !state.Selected("italic") {
		t.Fatalf("multiple selection=%v", state.Values)
	}
	items[0], items[2] = items[2], items[0]
	tt.Frame()
	if err := tt.Click("Italic"); err != nil {
		t.Fatal(err)
	}
	if state.Selected("italic") || !state.Selected("bold") {
		t.Fatalf("reordering moved selection identity: %v", state.Values)
	}
}

func TestVerticalToggleGroupUsesVerticalArrows(t *testing.T) {
	state := mygoui.ToggleGroupState{}
	tt := myui.NewTester(func(c *myui.Context) {
		mygoui.ToggleGroup(c, &state, mygoui.ToggleGroupOptions{Vertical: true}, []mygoui.ToggleGroupItem{
			{Value: "a", Label: "Top"}, {Value: "b", Label: "Bottom"},
		})
	}, 320, 200)
	tt.Key(0, myui.KeyTab)
	tt.Key(0, myui.KeyDown)
	if !tt.Focused("Bottom") {
		t.Fatal("Down did not move within vertical toggle group")
	}
	tt.Key(0, myui.KeyEnd)
	tt.Key(0, myui.KeySpace)
	if !state.Selected("b") {
		t.Fatal("vertical group did not activate focused value")
	}
}
