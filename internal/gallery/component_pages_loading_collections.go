package gallery

import (
	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// collectionLoadingState drives the Loading example of each collection: one
// switch sets the flag, the other chooses whether the caller has rows yet.
type collectionLoadingState struct {
	loading, loaded bool
	table           data.TableState
	listBox         list.ListBoxState
	tree            list.TreeState
	command         navigation.CommandState
	ready           bool
}

func collectionLoadingOf(s *State) *collectionLoadingState {
	st := StateOf[collectionLoadingState](s, "collection-loading")
	if !st.ready {
		st.ready = true
		st.loaded = true
		st.table.Selection.Row = -1
		st.tree.SetExpanded("src", true)
	}
	return st
}

func collectionLoadingSwitches(c *ui.Context, st *collectionLoadingState) {
	ui.Row(c).FillWidth().Gap(theme.Space24).Children(func() {
		form.Switch(c, form.SwitchOptions{On: &st.loading, Label: "Loading"})
		form.Switch(c, form.SwitchOptions{On: &st.loaded, Label: "Has rows"})
	})
}

func tableLoadingDemo(c *ui.Context, s *State) {
	st := collectionLoadingOf(s)
	tbl := StateOf[tableState](s, "table")
	tbl.ensure()
	rows := 0
	if st.loaded {
		rows = len(tbl.Rows)
	}
	options := tbl.tableExampleOptions(c, &st.table, rows)
	options.Loading = st.loading
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
		collectionLoadingSwitches(c, st)
		tableExampleFrame(c, options, func() string { return "" })
	})
}

func listBoxLoadingDemo(c *ui.Context, s *State) {
	st := collectionLoadingOf(s)
	var items []list.ListBoxItem
	if st.loaded {
		items = []list.ListBoxItem{
			{Key: "apple", Label: "Apple"}, {Key: "banana", Label: "Banana"},
			{Key: "cherry", Label: "Cherry"}, {Key: "grape", Label: "Grape"},
		}
	}
	formsFieldStack(c, func() {
		collectionLoadingSwitches(c, st)
		list.ListBox(c, list.ListBoxOptions{
			Label: "Fruit", Items: items, State: &st.listBox, Selection: list.SelectionSingle,
			Empty: "No fruit", Loading: st.loading,
		})
	})
}

func treeLoadingDemo(c *ui.Context, s *State) {
	st := collectionLoadingOf(s)
	var nodes []list.TreeNode
	if st.loaded {
		nodes = []list.TreeNode{
			{Key: "src", Label: "src", Children: []list.TreeNode{
				{Key: "main.go", Label: "main.go"}, {Key: "util.go", Label: "util.go"},
			}},
			{Key: "docs", Label: "docs", Children: []list.TreeNode{{Key: "readme", Label: "README.md"}}},
			{Key: "go.mod", Label: "go.mod"},
		}
	}
	formsFieldStack(c, func() {
		collectionLoadingSwitches(c, st)
		list.Tree(c, list.TreeOptions{
			Label: "Files", Nodes: nodes, State: &st.tree, Selection: list.SelectionSingle,
			Loading: st.loading,
		})
	})
}

func commandLoadingDemo(c *ui.Context, s *State) {
	st := collectionLoadingOf(s)
	var items []navigation.CommandItem
	if st.loaded {
		items = []navigation.CommandItem{
			{Key: "new", Label: "New project", Shortcut: "N"},
			{Key: "open", Label: "Open file", Shortcut: "O"},
			{Key: "settings", Label: "Settings"},
		}
	}
	formsFieldStack(c, func() {
		collectionLoadingSwitches(c, st)
		navigation.Command(c, navigation.CommandOptions{
			State: &st.command, Label: "Search actions", Placeholder: "Type a command",
			Items: items, Loading: st.loading,
		})
	})
}
