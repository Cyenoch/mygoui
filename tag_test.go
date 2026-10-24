package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestTagGroupSelectsAndRemoves(t *testing.T) {
	state := &mygoui.TagGroupState{}
	var changed []string
	removed := 0
	tags := []mygoui.TagOptions{
		{Key: "go", Label: "Go", Removable: true, OnRemove: func() { removed++ }},
		{Key: "native", Label: "Native"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.TagGroup(c, mygoui.TagGroupOptions{
			Label: "Topics", Selection: mygoui.TagSelectMultiple, State: state,
			OnChange: func(keys []string) { changed = keys },
		}, tags)
	}, 360, 160)
	if err := tt.Click("Go"); err != nil {
		t.Fatal(err)
	}
	if !state.IsSelected("go") || len(changed) != 1 || changed[0] != "go" {
		t.Fatalf("selection = %v, changed %v", state.Selected, changed)
	}
	if err := tt.Click("Remove"); err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Fatalf("remove calls = %d", removed)
	}
}

func TestTagGroupSingleSelectionReplaces(t *testing.T) {
	state := &mygoui.TagGroupState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.TagGroup(c, mygoui.TagGroupOptions{
			Selection: mygoui.TagSelectSingle, State: state,
		}, []mygoui.TagOptions{{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Beta"}})
	}, 360, 160)
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if state.IsSelected("a") || !state.IsSelected("b") || len(state.Selected) != 1 {
		t.Fatalf("single selection = %v", state.Selected)
	}
}

func TestTagRemovesWithDeleteKey(t *testing.T) {
	state := &mygoui.TagGroupState{}
	removed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.TagGroup(c, mygoui.TagGroupOptions{State: state}, []mygoui.TagOptions{
			{Key: "draft", Label: "Draft", Removable: true, OnRemove: func() { removed++ }},
		})
	}, 320, 160)
	if err := tt.Click("Draft"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDelete)
	if removed != 1 {
		t.Fatalf("Delete removed %d tags, want 1", removed)
	}
}
