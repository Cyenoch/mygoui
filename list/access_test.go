package list_test

import (
	"runtime"
	"slices"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func fruitBox(o *list.ListBoxOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.ListBox(c, *o)
	}, 320, 400)
}

// TestListBoxIsAListBoxOfOptions pins what assistive technology hears: a
// list box of options, each chosen or disabled itself, with the row the keys
// are on as the active one.
func TestListBoxIsAListBoxOfOptions(t *testing.T) {
	var state list.ListBoxState
	state.Selected.Add("apple")
	o := list.ListBoxOptions{Label: "Fruit", Items: items(), Selection: list.SelectionMultiple, State: &state}
	tt := fruitBox(&o)
	box, ok := tt.AccessByRole(ui.RoleListBox, "Fruit")
	if !ok || !box.Multiselectable {
		t.Fatalf("the list is not a multiselectable list box (%+v)", box)
	}
	apple, ok := tt.AccessByRole(ui.RoleOption, "Apple")
	if !ok || !apple.Selected {
		t.Fatalf("the chosen row is not a selected option (%+v)", apple)
	}
	if cherry, _ := tt.AccessByRole(ui.RoleOption, "Cherry"); !cherry.Disabled {
		t.Fatal("the disabled row is not a disabled option")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if banana, _ := tt.AccessByRole(ui.RoleOption, "Banana"); !banana.Focused {
		t.Fatal("the highlighted row is not the active descendant")
	}
	single := list.ListBoxOptions{Label: "Fruit", Items: items(), Selection: list.SelectionSingle, State: &list.ListBoxState{}}
	if box, _ := fruitBox(&single).AccessByRole(ui.RoleListBox, "Fruit"); box.Multiselectable {
		t.Fatal("a single selection is multiselectable")
	}
}

func manyFruit() []list.ListBoxItem {
	var out []list.ListBoxItem
	for _, l := range []string{"Apple", "Apricot", "Avocado", "Banana", "Blueberry", "Cherry",
		"Date", "Elderberry", "Fig", "Grape", "Guava", "Kiwi", "Lemon", "Lime", "Mango", "Melon"} {
		out = append(out, list.ListBoxItem{Key: l, Label: l})
	}
	return out
}

// TestListBoxTypeaheadBuildsAPrefix pins type-to-select: letters typed
// within a second narrow one prefix, and a pause starts anew.
func TestListBoxTypeaheadBuildsAPrefix(t *testing.T) {
	var state list.ListBoxState
	o := list.ListBoxOptions{Label: "Fruit", Items: manyFruit(), State: &state, Height: 200}
	tt := fruitBox(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyA)
	tt.Key(0, ui.KeyV)
	if got := state.HighlightKey(); got != "Avocado" {
		t.Fatalf("typing av highlighted %q, want Avocado", got)
	}
	tt.Advance(1100 * time.Millisecond)
	tt.Key(0, ui.KeyL)
	tt.Key(0, ui.KeyI)
	if got := state.HighlightKey(); got != "Lime" {
		t.Fatalf("typing li after a pause highlighted %q, want Lime", got)
	}
}

// TestListBoxPagesFollowThePlatform pins PageUp and PageDown: on macOS they
// scroll, as the toolkit's lists do; elsewhere they move the highlight by
// the rows in view.
func TestListBoxPagesFollowThePlatform(t *testing.T) {
	var state list.ListBoxState
	o := list.ListBoxOptions{Label: "Fruit", Items: manyFruit(), State: &state, Height: 120}
	tt := fruitBox(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageDown)
	tt.Advance(400 * time.Millisecond)
	if runtime.GOOS == "darwin" {
		if got := state.HighlightKey(); got != "Apple" {
			t.Fatalf("PageDown moved the highlight to %q on macOS", got)
		}
		if first, _ := state.List.Visible(); first == 0 {
			t.Fatal("PageDown did not scroll on macOS")
		}
		return
	}
	if got := state.HighlightKey(); got == "Apple" || got == "Apricot" {
		t.Fatalf("PageDown highlighted %q, want a page further", got)
	}
	tt.Key(0, ui.KeyPageUp)
	if got := state.HighlightKey(); got != "Apple" {
		t.Fatalf("PageUp highlighted %q, want Apple", got)
	}
}

// TestListBoxExtendsAMultipleSelectionFromTheKeyboard pins the desktop range
// keys of a multiple selection.
func TestListBoxExtendsAMultipleSelectionFromTheKeyboard(t *testing.T) {
	var state list.ListBoxState
	var reported []string
	o := list.ListBoxOptions{
		Label: "Fruit", Items: items(), Selection: list.SelectionMultiple, State: &state,
		OnSelectionChange: func(keys []string) { reported = keys },
	}
	tt := fruitBox(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyDown)
	got := state.SelectedKeys(list.SelectionMultiple)
	slices.Sort(got)
	if !slices.Equal(got, []string{"apple", "banana"}) || len(reported) != 2 {
		t.Fatalf("Shift+Down chose %v (reported %v), want apple and banana", got, reported)
	}
	state.ClearSelection()
	tt.Key(ui.Cmd, ui.KeyA)
	got = state.SelectedKeys(list.SelectionMultiple)
	slices.Sort(got)
	if !slices.Equal(got, []string{"apple", "banana"}) {
		t.Fatalf("Cmd+A chose %v, want every enabled row", got)
	}
}

// TestTreeMarksAMultipleSelectionAndCanBeDisabled pins the tree's
// multiselectable state and its whole-tree Disabled.
func TestTreeMarksAMultipleSelectionAndCanBeDisabled(t *testing.T) {
	var state list.TreeState
	o := list.TreeOptions{Label: "Files", Nodes: treeNodes(), Selection: list.SelectionMultiple, State: &state}
	tt := treeTester(&o)
	if tree, _ := tt.AccessByRole(ui.RoleTree, "Files"); !tree.Multiselectable {
		t.Fatal("a multiple selection does not mark the tree multiselectable")
	}
	o.Disabled = true
	tt.Frame()
	if tree, _ := tt.AccessByRole(ui.RoleTree, "Files"); !tree.Disabled {
		t.Fatal("a disabled tree is not disabled")
	}
	tt.Key(0, ui.KeyTab)
	if tree, _ := tt.AccessByRole(ui.RoleTree, "Files"); tree.Focused {
		t.Fatal("a disabled tree takes the focus")
	}
	if err := tt.Click("go.mod"); err == nil && state.Selected.Len() > 0 {
		t.Fatal("a disabled tree took a click")
	}
}

// TestTreeMenuKeyOpensTheKeyboardRowsMenu pins Shift+F10 and the menu key:
// they open the menu of the row the keys are on, choosing it first.
func TestTreeMenuKeyOpensTheKeyboardRowsMenu(t *testing.T) {
	var state list.TreeState
	var asked []string
	o := list.TreeOptions{
		Label: "Files", Nodes: treeNodes(), Selection: list.SelectionSingle, State: &state,
		OnContextMenu: func(key string, m *ui.Menu) {
			asked = append(asked, key)
			m.Item("Rename")
		},
	}
	tt := treeTester(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(ui.Shift, ui.KeyF10)
	if len(asked) == 0 || asked[len(asked)-1] != "go.mod" {
		t.Fatalf("Shift+F10 asked for the menu of %v, want go.mod", asked)
	}
	if !slices.Contains(tt.Menu(), "Rename") {
		t.Fatal("Shift+F10 opened no menu")
	}
	if state.SelectedKey != "go.mod" {
		t.Fatalf("the menu's row was not chosen first: %q", state.SelectedKey)
	}
}

// TestTreePagesFollowThePlatform pins PageUp and PageDown as ListBox's.
func TestTreePagesFollowThePlatform(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("PageUp and PageDown scroll on macOS; TestListBoxPagesFollowThePlatform covers it")
	}
	var state list.TreeState
	var nodes []list.TreeNode
	for _, l := range []string{"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "b1", "b2", "b3"} {
		nodes = append(nodes, list.TreeNode{Key: l, Label: l})
	}
	o := list.TreeOptions{Label: "Files", Nodes: nodes, Selection: list.SelectionSingle, State: &state, Height: 100}
	tt := treeTester(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageDown)
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey == "a1" || state.SelectedKey == "" {
		t.Fatalf("PageDown chose %q, want a page further", state.SelectedKey)
	}
}
