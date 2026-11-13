package list_test

import (
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func treeNodes() []list.TreeNode {
	return []list.TreeNode{
		{Key: "src", Label: "src", Children: []list.TreeNode{
			{Key: "src/main.go", Label: "main.go"},
			{Key: "src/util.go", Label: "util.go"},
			{Key: "src/gen.go", Label: "gen.go", Disabled: true},
		}},
		{Key: "go.mod", Label: "go.mod"},
		{Key: "vendor", Label: "vendor", Disabled: true, Children: []list.TreeNode{
			{Key: "vendor/x", Label: "x"},
		}},
	}
}

func treeTester(o *list.TreeOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		list.Tree(c, *o)
	}, 320, 320)
}

func TestTreeKeyboardExpandsMovesAndSelects(t *testing.T) {
	var state list.TreeState
	var toggled []string
	o := list.TreeOptions{
		Label: "Files", Nodes: treeNodes(), Selection: list.SelectionSingle, State: &state,
		OnExpandedChange: func(key string, open bool) { toggled = append(toggled, key+"="+strconv.FormatBool(open)) },
	}
	tt := treeTester(&o)
	if tt.HasText("main.go") {
		t.Fatal("a closed branch shows its children")
	}
	if err := tt.Click("src"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if !state.IsExpanded("src") || !tt.HasText("main.go") {
		t.Fatal("Right did not open the branch")
	}
	tt.Key(0, ui.KeyRight) // to the first child
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown) // gen.go is disabled: on to go.mod
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "go.mod" {
		t.Fatalf("Enter chose %q, want go.mod past the disabled row", state.SelectedKey)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyU) // typeahead to util.go
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "src/util.go" {
		t.Fatalf("typeahead chose %q", state.SelectedKey)
	}
	tt.Key(0, ui.KeyLeft) // to the parent
	tt.Key(0, ui.KeyLeft) // closes it
	if state.IsExpanded("src") || tt.HasText("main.go") {
		t.Fatal("Left did not close the branch")
	}
	if want := []string{"src=true", "src=false"}; !slices.Equal(toggled, want) {
		t.Fatalf("expansion reports = %v, want %v", toggled, want)
	}
	if err := tt.Click("vendor"); err == nil && (state.SelectedKey == "vendor" || state.IsExpanded("vendor")) {
		t.Fatal("a disabled branch took a click")
	}
}

func TestTreeChecksCascadeAndMixParents(t *testing.T) {
	var state list.TreeState
	state.SetExpanded("src", true)
	var reported []string
	o := list.TreeOptions{
		Label: "Files", Nodes: treeNodes(), Checkable: true, State: &state,
		OnCheckedChange: func(keys []string) { reported = keys },
	}
	tt := treeTester(&o)
	if err := tt.Click("main.go"); err != nil {
		t.Fatal(err)
	}
	if !state.IsChecked("src/main.go") || !slices.Equal(reported, []string{"src/main.go"}) {
		t.Fatalf("leaf check = %v, reported %v", state.IsChecked("src/main.go"), reported)
	}
	parent, ok := tt.AccessByRole(ui.RoleTreeItem, "src")
	if !ok || !parent.Mixed {
		t.Fatalf("parent of a partly checked branch: %+v, found %v", parent, ok)
	}
	tt.Key(0, ui.KeyLeft) // from main.go to its parent
	tt.Key(0, ui.KeySpace)
	if !state.IsChecked("src/util.go") || state.IsChecked("src/gen.go") {
		t.Fatalf("parent check: util %v, disabled gen %v",
			state.IsChecked("src/util.go"), state.IsChecked("src/gen.go"))
	}
	// The disabled leaf stays unchecked, so the branch is still mixed, and a
	// second press clears the leaves it can reach.
	tt.Key(0, ui.KeySpace)
	if len(state.CheckedKeys()) != 0 {
		t.Fatalf("second parent press left %v checked", state.CheckedKeys())
	}
}

func TestTreeExposesHierarchyToAssistiveTechnology(t *testing.T) {
	var state list.TreeState
	state.SetExpanded("src", true)
	state.SelectedKey = "src/util.go"
	o := list.TreeOptions{Label: "Files", Nodes: treeNodes(), Selection: list.SelectionSingle, State: &state}
	tt := treeTester(&o)
	if _, ok := tt.AccessByRole(ui.RoleTree, "Files"); !ok {
		t.Fatal("no tree named Files")
	}
	branch, _ := tt.AccessByRole(ui.RoleTreeItem, "src")
	if branch.Level != 1 || !branch.Expanded || branch.PosInSet != 1 || branch.SetSize != 3 {
		t.Fatalf("branch: %+v", branch)
	}
	leaf, _ := tt.AccessByRole(ui.RoleTreeItem, "util.go")
	if leaf.Level != 2 || !leaf.Selected || leaf.PosInSet != 2 || leaf.SetSize != 3 {
		t.Fatalf("leaf: %+v", leaf)
	}
}

func TestTreeBuildsOnlyTheViewportAndFollowsRevision(t *testing.T) {
	nodes := []list.TreeNode{{Key: "root", Label: "root"}}
	for i := range 5000 {
		nodes[0].Children = append(nodes[0].Children, list.TreeNode{Label: "node " + strconv.Itoa(i)})
	}
	var state list.TreeState
	state.SetExpanded("root", true)
	o := list.TreeOptions{Label: "Big", Nodes: nodes, State: &state, Height: 280}
	tt := treeTester(&o)
	if !tt.HasText("node 0") || tt.HasText("node 4999") {
		t.Fatal("the tree did not build just the rows in view")
	}
	if built := len(tt.Texts()); built > 100 {
		t.Fatalf("%d texts built for a 280 DIP viewport", built)
	}
	nodes[0].Children[0].Label = "renamed"
	state.Revision++
	tt.Frame()
	if !tt.HasText("renamed") {
		t.Fatal("a revision bump did not show the in-place edit")
	}
}

func TestTreeLoadingRefreshKeepsRowsAndIgnoresInput(t *testing.T) {
	var state list.TreeState
	o := list.TreeOptions{Label: "Files", Nodes: treeNodes(), Selection: list.SelectionSingle, State: &state, Loading: true}
	tt := treeTester(&o)
	tt.Advance(400 * time.Millisecond)
	if !tt.HasText("go.mod") {
		t.Fatalf("rows left while refreshing: %v", tt.Texts())
	}
	_ = tt.Click("go.mod")
	if state.SelectedKey != "" {
		t.Fatalf("refreshing tree took a click: %q", state.SelectedKey)
	}
	o.Loading = false
	tt.Frame()
	tt.Advance(400 * time.Millisecond)
	if err := tt.Click("go.mod"); err != nil || state.SelectedKey != "go.mod" {
		t.Fatalf("tree stayed inert after loading: %v %q", err, state.SelectedKey)
	}
}

func TestTreeLoadingFirstLoadShowsNoNodes(t *testing.T) {
	var state list.TreeState
	o := list.TreeOptions{Label: "Files", State: &state, Loading: true}
	tt := treeTester(&o)
	tt.Advance(400 * time.Millisecond)
	if !slices.Equal(tt.Texts(), []string{"Files"}) {
		t.Fatalf("loading tree shows text: %v", tt.Texts())
	}
}

// The tree keeps the keyboard focus on its list, so assistive technology
// learns the row the keys are on only through the active descendant.
func TestTreeReportsTheKeyboardRowToAssistiveTechnology(t *testing.T) {
	var state list.TreeState
	o := list.TreeOptions{Label: "Files", Nodes: treeNodes(), Selection: list.SelectionSingle, State: &state}
	tt := treeTester(&o)
	if err := tt.Click("src"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyDown)
	if a, ok := tt.AccessByRole(ui.RoleTreeItem, "go.mod"); !ok || !a.Focused {
		t.Fatalf("the row under the keys reads focused=%v (found=%v)", a.Focused, ok)
	}
}

func multiTreeNodes() []list.TreeNode {
	return []list.TreeNode{
		{Key: "a", Label: "alpha"}, {Key: "b", Label: "bravo"}, {Key: "c", Label: "charlie", Disabled: true},
		{Key: "d", Label: "delta"}, {Key: "e", Label: "echo"},
	}
}

func selectedKeys(s *list.TreeState) []string {
	var keys []string
	for k := range s.Selected.All() {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func TestTreeMultipleSelectionFollowsDesktopRules(t *testing.T) {
	var state list.TreeState
	var reports int
	o := list.TreeOptions{
		Label: "Files", Nodes: multiTreeNodes(), Selection: list.SelectionMultiple, State: &state,
		OnSelectionChange: func([]string) { reports++ },
	}
	tt := treeTester(&o)
	want := func(step string, keys ...string) {
		t.Helper()
		if got := selectedKeys(&state); !slices.Equal(got, keys) {
			t.Fatalf("%s: selected %v, want %v", step, got, keys)
		}
	}
	_ = tt.Click("alpha")
	_ = tt.Click("delta")
	want("plain click replaces", "d")
	_ = tt.ClickWith(ui.Cmd, "alpha")
	want("cmd click adds", "a", "d")
	_ = tt.ClickWith(ui.Cmd, "delta")
	want("cmd click removes", "a")
	_ = tt.Click("bravo")
	_ = tt.ClickWith(ui.Shift, "echo")
	want("shift click skips the disabled row", "b", "d", "e")
	tt.Key(ui.Shift, ui.KeyUp)
	want("shift arrow shrinks the range", "b", "d")
	tt.Key(ui.Cmd, ui.KeyA)
	want("select all skips the disabled row", "a", "b", "d", "e")
	tt.Key(0, ui.KeyEscape)
	want("escape keeps the selection", "a", "b", "d", "e")
	before := reports
	tt.Key(ui.Cmd, ui.KeyA)
	if reports != before {
		t.Fatal("select all on a full selection reported a change")
	}
}

func TestTreeDoubleClickOpensABranchWhileSelecting(t *testing.T) {
	var state list.TreeState
	o := list.TreeOptions{Label: "Files", Nodes: treeNodes(), Selection: list.SelectionSingle, State: &state}
	tt := treeTester(&o)
	_ = tt.Click("src")
	if state.IsExpanded("src") || state.SelectedKey != "src" {
		t.Fatalf("a single click opened the branch or missed it: %v %q", state.IsExpanded("src"), state.SelectedKey)
	}
	_ = tt.Click("src")
	if !state.IsExpanded("src") {
		t.Fatal("a double click did not open the branch")
	}
}

func TestTreeTypeaheadMatchesAPrefixWithinTheWindow(t *testing.T) {
	var state list.TreeState
	nodes := []list.TreeNode{{Label: "gamma"}, {Label: "group"}, {Label: "grid"}, {Label: "hotel"}}
	o := list.TreeOptions{Label: "Files", Nodes: nodes, Selection: list.SelectionSingle, State: &state}
	tt := treeTester(&o)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyG)
	tt.Key(0, ui.KeyR)
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "group" {
		t.Fatalf("typing gr chose %q, want group", state.SelectedKey)
	}
	tt.Key(0, ui.KeyI)
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "grid" {
		t.Fatalf("typing on to gri chose %q, want grid", state.SelectedKey)
	}
	tt.Advance(time.Second)
	tt.Key(0, ui.KeyH)
	tt.Key(0, ui.KeyEnter)
	if state.SelectedKey != "hotel" {
		t.Fatalf("a new prefix after the pause chose %q, want hotel", state.SelectedKey)
	}
}

func TestTreeLazyBranchesAskForTheirChildrenAndReportProgress(t *testing.T) {
	var state list.TreeState
	var opened []string
	nodes := []list.TreeNode{{Key: "remote", Label: "remote", Lazy: true}, {Key: "leaf", Label: "leaf"}}
	o := list.TreeOptions{
		Label: "Files", Nodes: nodes, State: &state,
		OnExpandedChange: func(key string, open bool) { opened = append(opened, key) },
	}
	tt := treeTester(&o)
	_ = tt.Click("remote")
	if !slices.Equal(opened, []string{"remote"}) {
		t.Fatalf("expansion reports = %v", opened)
	}
	nodes[0].Loading = true
	tt.Frame()
	if row, ok := tt.AccessByRole(ui.RoleTreeItem, "remote"); !ok || !row.Busy {
		t.Fatal("a loading node is not reported busy")
	}
	nodes[0].Loading = false
	tt.Frame()
	if row, ok := tt.AccessByRole(ui.RoleTreeItem, "remote"); !ok || row.Busy {
		t.Fatal("the node stayed busy after loading ended")
	}
}

func TestTreeContextMenuChoosesTheRowFirst(t *testing.T) {
	var state list.TreeState
	var menuFor string
	o := list.TreeOptions{
		Label: "Files", Nodes: multiTreeNodes(), Selection: list.SelectionMultiple, State: &state,
		OnContextMenu: func(key string, m *ui.Menu) { menuFor = key; m.Item("Rename") },
	}
	tt := treeTester(&o)
	_ = tt.Click("alpha")
	_ = tt.ClickWith(ui.Cmd, "bravo")
	if err := tt.RightClick("bravo"); err != nil {
		t.Fatal(err)
	}
	if menuFor != "b" || !slices.Equal(selectedKeys(&state), []string{"a", "b"}) {
		t.Fatalf("menu on %q kept selection %v", menuFor, selectedKeys(&state))
	}
	if err := tt.RightClick("delta"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(selectedKeys(&state), []string{"d"}) {
		t.Fatalf("menu on an unchosen row left %v", selectedKeys(&state))
	}
	if err := tt.RightClick("charlie"); err == nil && menuFor == "c" {
		t.Fatal("a disabled row opened a menu")
	}
}

func TestTreeSpaceSelectsWithoutRunningOnAction(t *testing.T) {
	var state list.TreeState
	var acted []string
	o := list.TreeOptions{
		Label: "Files", Nodes: multiTreeNodes(), Selection: list.SelectionMultiple, State: &state,
		OnAction: func(key string) { acted = append(acted, key) },
	}
	tt := treeTester(&o)
	_ = tt.Click("alpha")
	acted = nil
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeySpace)
	if len(acted) != 0 {
		t.Fatalf("Space ran OnAction for %v", acted)
	}
	tt.Key(0, ui.KeyEnter)
	if len(acted) != 1 {
		t.Fatalf("Enter acted %v", acted)
	}
}

func TestTreeUnloadedLazyBranchChecksAsAUnit(t *testing.T) {
	var state list.TreeState
	nodes := []list.TreeNode{{Key: "root", Label: "root", Children: []list.TreeNode{
		{Key: "remote", Label: "remote", Lazy: true},
	}}}
	state.SetExpanded("root", true)
	o := list.TreeOptions{Label: "Files", Nodes: nodes, Checkable: true, State: &state}
	tt := treeTester(&o)
	_ = tt.Click("root")
	if !state.IsChecked("remote") {
		t.Fatal("checking the parent left the unloaded branch unchecked")
	}
	if row, _ := tt.AccessByRole(ui.RoleTreeItem, "root"); row.Mixed {
		t.Fatal("parent is mixed though its only child is checked")
	}
}
