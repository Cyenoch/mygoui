package list

import (
	"runtime"
	"strings"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/busy"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TreeNode is one item of a Tree. Key is its stable identity in the
// expansion, the selection and the checks; it falls back to Label and must
// be unique in the tree. A node with Children is a branch, and so is a Lazy
// one. Disabled stops a node and everything below it from being chosen,
// checked or toggled.
type TreeNode struct {
	Key      string
	Label    string
	Icon     *ui.SVG
	Disabled bool
	Children []TreeNode
	// Lazy makes the node a branch before it has Children: the caller loads
	// them when OnExpandedChange reports it opened. Once loaded, an empty
	// branch shows nothing below it. Bump State.Revision when Children are
	// set in place. While it has no Children it takes a check of its own, so
	// a checked parent counts it; the loaded children do not inherit that
	// check and start unchecked.
	Lazy bool
	// Loading shows a spinner in place of the node's arrow while its
	// children are on their way, and makes the node a branch. It needs no
	// Revision bump.
	Loading bool
}

func (n *TreeNode) key() string {
	if n.Key != "" {
		return n.Key
	}
	return n.Label
}

// TreeOptions configures a Tree. State is required and caller-owned.
type TreeOptions struct {
	// Label names the tree for assistive technology.
	Label string
	// Nodes are the top-level items. Bump State.Revision after editing them
	// in place.
	Nodes []TreeNode
	// Selection selects how rows are chosen. The zero value chooses none: a
	// press then checks a Checkable row, or opens and closes a branch.
	Selection SelectionMode
	// Checkable shows a check box on every row. A branch is checked when all
	// leaves below it are, and mixed when some are.
	Checkable bool
	// Guides draws a hairline at each ancestor level of a row.
	Guides bool
	// State is required and owned by the caller.
	State *TreeState
	// Height fixes the viewport in DIPs; without it the tree sizes to its
	// rows, and MaxHeight caps it.
	Height    float32
	MaxHeight float32
	// Loading marks the nodes as on their way. With no rows the tree shows
	// LoadingRows placeholder rows at a few depths; with rows it dims them,
	// ignores input and draws a progress line along the top edge.
	Loading bool
	// LoadingRows is the number of placeholder rows; 0 shows 5.
	LoadingRows int
	// Disabled dims the tree, takes it out of the Tab order and stops every
	// key and click.
	Disabled bool
	// OnAction reports a row activated by a click or Enter, not by Space.
	OnAction func(key string)
	// OnSelectionChange reports the chosen keys after a user change.
	OnSelectionChange func(keys []string)
	// OnCheckedChange reports the checked leaf keys after a user change.
	OnCheckedChange func(keys []string)
	// OnExpandedChange reports a branch the user opened or closed.
	OnExpandedChange func(key string, expanded bool)
	// OnContextMenu fills the menu opened on the row of key. A row that is
	// not chosen is chosen first, as Finder does; with Selection none the
	// choice is left alone. Disabled rows open none.
	OnContextMenu func(key string, m *ui.Menu)
}

// TreeState is the caller-owned expansion, selection, checks and scroll
// state of one Tree. Keep one per tree and do not copy it.
//
// The tree indexes its nodes once and keeps the rows it shows between
// frames. A new Nodes slice or length is noticed; Revision must change when
// a node or its children are edited in place.
type TreeState struct {
	// List keeps the virtual list's place, row heights and focus.
	List ui.ListState
	// SelectedKey is the chosen key in SelectionSingle mode; empty for none.
	SelectedKey string
	// Selected holds the chosen keys in SelectionMultiple mode.
	Selected ui.Selection[string]
	// Revision tells the tree that nodes changed in place. Increment it
	// when a key, label, disabled flag or child list changes without the
	// top-level slice or its length changing.
	Revision uint64

	expanded ui.Selection[string]
	dim      rowDim
	checked  ui.Selection[string]

	entries []treeEntry // every node, parents before their children
	rows    []int32     // the entries shown, in order
	cursor  int32       // the entry under the keyboard, -1 for none
	// anchor is the key a Shift range starts from.
	anchor string

	// doubled is the entry a double click was seen on, plus one; 0 for none.
	doubled int32

	// typed is the prefix typed within typeWindow of typedAt; letters waits
	// for the next build, which knows the time.
	typed   string
	typedAt time.Time
	letters []rune
	// typing holds the type-ahead key handlers, made once for the state
	// they belong to rather than on every build.
	typing      [len(typeKeys)]func()
	typingOwner *TreeState

	// reveal is the branch the user last opened, -1 for none, and revealAt
	// when its rows first showed; revealNew asks the next build to stamp it.
	reveal    int32
	revealAt  time.Time
	revealNew bool

	source    *TreeNode
	sourceLen int
	revision  uint64
	indexed   bool
	rowsStale bool
	// countsStale asks for the branches' checked counts again.
	countsStale bool
}

// treeEntry is an indexed node. end is the entry after its subtree, pos and
// size place it among its siblings, and leaves and checked count the leaves
// below it (itself for a leaf).
type treeEntry struct {
	node      *TreeNode
	parent    int32
	depth     int32
	end       int32
	pos, size int32
	leaves    int32
	checked   int32
	disabled  bool
}

// unit reports whether the entry holds its own check: a leaf, or a branch
// whose children are not there yet.
func (e *treeEntry) unit() bool { return len(e.node.Children) == 0 }

func (e *treeEntry) branch() bool {
	return len(e.node.Children) > 0 || e.node.Lazy || e.node.Loading
}

// typeWindow is how long typed letters keep building one prefix: the
// toolkit's pause (MyGo ui/input.go, typePause), as a menu's.
const typeWindow = time.Second

// Focus gives the tree the keyboard focus and reports whether it was built
// in the current pass. Call it after building the tree.
func (s *TreeState) Focus(c *ui.Context) bool { return s.List.Focus(c) }

// IsExpanded reports whether the branch of key shows its children.
func (s *TreeState) IsExpanded(key string) bool { return s.expanded.Has(key) }

// SetExpanded opens or closes the branch of key.
func (s *TreeState) SetExpanded(key string, open bool) {
	if s.expanded.Has(key) == open {
		return
	}
	if open {
		s.expanded.Add(key)
	} else {
		s.expanded.Remove(key)
	}
	s.rowsStale = true
}

// IsChecked reports whether the leaf of key is checked.
func (s *TreeState) IsChecked(key string) bool { return s.checked.Has(key) }

// SetChecked checks or unchecks the leaf of key. Branches follow their
// leaves and hold no check of their own.
func (s *TreeState) SetChecked(key string, on bool) {
	if s.checked.Has(key) == on {
		return
	}
	if on {
		s.checked.Add(key)
	} else {
		s.checked.Remove(key)
	}
	s.countsStale = true
}

// CheckedKeys returns the checked leaf keys, in no order.
func (s *TreeState) CheckedKeys() []string {
	keys := make([]string, 0, s.checked.Len())
	for k := range s.checked.All() {
		keys = append(keys, k)
	}
	return keys
}

// IsSelected reports whether key is chosen in either mode.
func (s *TreeState) IsSelected(key string) bool {
	return s.SelectedKey == key || s.Selected.Has(key)
}

// selected reports whether key is chosen in mode, apart from what the other
// mode's store still holds.
func (s *TreeState) selected(mode SelectionMode, key string) bool {
	switch mode {
	case SelectionSingle:
		return key != "" && s.SelectedKey == key
	case SelectionMultiple:
		return s.Selected.Has(key)
	}
	return false
}

// selectedKeys returns the chosen keys for mode.
func (s *TreeState) selectedKeys(mode SelectionMode) []string {
	switch mode {
	case SelectionSingle:
		if s.SelectedKey != "" {
			return []string{s.SelectedKey}
		}
	case SelectionMultiple:
		keys := make([]string, 0, s.Selected.Len())
		for k := range s.Selected.All() {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

// sync brings the index, the checked counts and the shown rows up to date.
// Nothing runs while the nodes, the revision, the expansion and the checks
// are as the last build left them.
func (s *TreeState) sync(o TreeOptions) {
	var source *TreeNode
	if len(o.Nodes) > 0 {
		source = &o.Nodes[0]
	}
	if !s.indexed || source != s.source || len(o.Nodes) != s.sourceLen || s.Revision != s.revision {
		s.reindex(o.Nodes)
		s.source, s.sourceLen, s.revision, s.indexed = source, len(o.Nodes), s.Revision, true
	}
	if s.countsStale && o.Checkable {
		s.recount()
		s.countsStale = false
	}
	if s.rowsStale {
		s.flatten()
		s.rowsStale = false
	}
}

// reindex lists every node, keeping the keyboard's place by key.
func (s *TreeState) reindex(nodes []TreeNode) {
	at := ""
	if s.indexed && s.cursor >= 0 && int(s.cursor) < len(s.entries) {
		at = s.entries[s.cursor].node.key()
	}
	s.entries = s.entries[:0]
	s.index(nodes, -1, 0, false)
	s.cursor = -1
	s.reveal = -1
	if at != "" {
		for i := range s.entries {
			if s.entries[i].node.key() == at {
				if !s.entries[i].disabled {
					s.cursor = int32(i)
				}
				break
			}
		}
	}
	s.rowsStale, s.countsStale = true, true
}

// index appends nodes and their subtrees, and returns their leaf count.
func (s *TreeState) index(nodes []TreeNode, parent, depth int32, disabled bool) int32 {
	var leaves int32
	for i := range nodes {
		n := &nodes[i]
		at := len(s.entries)
		s.entries = append(s.entries, treeEntry{
			node: n, parent: parent, depth: depth,
			pos: int32(i + 1), size: int32(len(nodes)),
			leaves: 1, disabled: disabled || n.Disabled,
		})
		if len(n.Children) > 0 {
			below := s.index(n.Children, int32(at), depth+1, disabled || n.Disabled)
			s.entries[at].leaves = below
		}
		s.entries[at].end = int32(len(s.entries))
		leaves += s.entries[at].leaves
	}
	return leaves
}

// recount counts the checked leaves below every branch. Children follow
// their parent, so a pass from the end sums each subtree before its root.
func (s *TreeState) recount() {
	for i := range s.entries {
		s.entries[i].checked = 0
	}
	for i := len(s.entries) - 1; i >= 0; i-- {
		e := &s.entries[i]
		if e.unit() && s.checked.Has(e.node.key()) {
			e.checked = 1
		}
		if e.parent >= 0 {
			s.entries[e.parent].checked += e.checked
		}
	}
}

// flatten lists the entries shown: the top level, and the children of open
// branches. A keyboard place hidden by a closed branch moves to that branch;
// without one, the keyboard starts at the first enabled row.
func (s *TreeState) flatten() {
	s.rows = s.rows[:0]
	for i := 0; i < len(s.entries); {
		e := &s.entries[i]
		s.rows = append(s.rows, int32(i))
		if e.branch() && !s.expanded.Has(e.node.key()) {
			i = int(e.end)
		} else {
			i++
		}
	}
	if s.cursor >= 0 {
		for p := s.entries[s.cursor].parent; p >= 0; p = s.entries[p].parent {
			if !s.expanded.Has(s.entries[p].node.key()) {
				s.cursor = p
			}
		}
		return
	}
	for _, i := range s.rows {
		if !s.entries[i].disabled {
			s.cursor = i
			return
		}
	}
}

// rowOf returns the row showing entry i, or -1.
func (s *TreeState) rowOf(i int32) int {
	for r, at := range s.rows {
		if at == i {
			return r
		}
	}
	return -1
}

// setCursor moves the keyboard to row r and scrolls it into view.
func (s *TreeState) setCursor(r int) {
	s.cursor = s.rows[r]
	s.List.ScrollIntoView(r)
}

// move steps the keyboard by rows from the current one, over disabled rows.
// From no row it starts at the edge the step comes from.
func (s *TreeState) move(step int) {
	r := s.rowOf(s.cursor)
	if r < 0 && step < 0 {
		r = len(s.rows)
	}
	for r += step; r >= 0 && r < len(s.rows); r += step {
		if !s.entries[s.rows[r]].disabled {
			s.setCursor(r)
			return
		}
	}
}

// page moves the keyboard by the rows in view, stopping at the last enabled
// row on the way.
func (s *TreeState) page(dir int) {
	first, last := s.List.Visible()
	r := s.rowOf(s.cursor)
	to := min(max(r+dir*max(last-first, 1), 0), len(s.rows)-1)
	for t := to; t != r && t >= 0; t -= dir {
		if !s.entries[s.rows[t]].disabled {
			s.setCursor(t)
			return
		}
	}
}

// edge moves the keyboard to the first or the last enabled row.
func (s *TreeState) edge(last bool) {
	s.cursor = -1
	if last {
		s.move(-1)
	} else {
		s.move(1)
	}
}

// typeahead moves the keyboard to the row whose label starts with what was
// typed, as typePrefix finds it.
func (s *TreeState) typeahead(r rune, now time.Time) {
	prefix, next := typePrefix(&s.typed, &s.typedAt, r, now)
	from := max(s.rowOf(s.cursor), 0)
	if next {
		from = s.rowOf(s.cursor) + 1
	}
	n := len(s.rows)
	for d := range n {
		row := (from + d) % n
		e := &s.entries[s.rows[row]]
		if !e.disabled && strings.HasPrefix(strings.ToLower(e.node.Label), prefix) {
			s.setCursor(row)
			return
		}
	}
}

// current returns the entry under the keyboard, or nil.
func (s *TreeState) current() *treeEntry {
	if s.cursor < 0 {
		return nil
	}
	return &s.entries[s.cursor]
}

// toggleExpanded opens or closes the branch e and reports it.
func (s *TreeState) toggleExpanded(o TreeOptions, i int32) {
	key := s.entries[i].node.key()
	open := !s.expanded.Has(key)
	s.SetExpanded(key, open)
	if open {
		s.reveal, s.revealNew = i, true
	}
	// The keys that follow in this pass move among the rows as they now are.
	s.flatten()
	s.rowsStale = false
	if o.OnExpandedChange != nil {
		o.OnExpandedChange(key, open)
	}
}

// descend opens the branch under the keyboard, or moves to its first child.
func (s *TreeState) descend(o TreeOptions) {
	e := s.current()
	if e == nil || !e.branch() {
		return
	}
	if !s.expanded.Has(e.node.key()) {
		s.toggleExpanded(o, s.cursor)
		return
	}
	if r := s.rowOf(s.cursor); r+1 < len(s.rows) && !s.entries[s.rows[r+1]].disabled {
		s.setCursor(r + 1)
	}
}

// ascend closes the branch under the keyboard, or moves to its parent.
func (s *TreeState) ascend(o TreeOptions) {
	e := s.current()
	if e == nil {
		return
	}
	if e.branch() && s.expanded.Has(e.node.key()) {
		s.toggleExpanded(o, s.cursor)
		return
	}
	if e.parent >= 0 && !s.entries[e.parent].disabled {
		s.setCursor(s.rowOf(e.parent))
	}
}

// toggleChecked checks every enabled leaf of entry i's subtree, or unchecks
// them when all already are.
func (s *TreeState) toggleChecked(o TreeOptions, i int32) {
	end := s.entries[i].end
	on := false
	for j := i; j < end; j++ {
		if e := &s.entries[j]; e.unit() && !e.disabled && !s.checked.Has(e.node.key()) {
			on = true
			break
		}
	}
	changed := false
	for j := i; j < end; j++ {
		if e := &s.entries[j]; e.unit() && !e.disabled {
			s.SetChecked(e.node.key(), on)
			changed = true
		}
	}
	if changed && o.OnCheckedChange != nil {
		o.OnCheckedChange(s.CheckedKeys())
	}
}

// choose sets the single selection to key and reports a change.
func (s *TreeState) choose(o TreeOptions, key string) {
	if s.SelectedKey == key {
		return
	}
	s.SelectedKey = key
	if o.OnSelectionChange != nil {
		o.OnSelectionChange(s.selectedKeys(SelectionSingle))
	}
}

// setSelected makes keys the multiple selection, or adds them to it, and
// reports a change.
func (s *TreeState) setSelected(o TreeOptions, keys []string, add bool) {
	changed := false
	if !add {
		keep := make(map[string]bool, len(keys))
		for _, k := range keys {
			keep[k] = true
		}
		for _, k := range s.selectedKeys(SelectionMultiple) {
			if !keep[k] {
				s.Selected.Remove(k)
				changed = true
			}
		}
	}
	for _, k := range keys {
		if !s.Selected.Has(k) {
			s.Selected.Add(k)
			changed = true
		}
	}
	if changed && o.OnSelectionChange != nil {
		o.OnSelectionChange(s.selectedKeys(SelectionMultiple))
	}
}

// enabledKeys lists the keys of the enabled rows from row a to row b, either
// way round.
func (s *TreeState) enabledKeys(a, b int) []string {
	if a > b {
		a, b = b, a
	}
	var keys []string
	for r := max(a, 0); r <= b && r < len(s.rows); r++ {
		if e := &s.entries[s.rows[r]]; !e.disabled {
			keys = append(keys, e.node.key())
		}
	}
	return keys
}

// rowOfKey returns the row showing key, or -1.
func (s *TreeState) rowOfKey(key string) int {
	for r, i := range s.rows {
		if s.entries[i].node.key() == key {
			return r
		}
	}
	return -1
}

// chooseMultiple applies a press on entry i as desktop lists do: it replaces
// the selection, Cmd or Ctrl toggles the row, Shift selects from the anchor
// (adding to the selection with Cmd or Ctrl too).
func (s *TreeState) chooseMultiple(o TreeOptions, i int32, mods ui.Modifiers) {
	key := s.entries[i].node.key()
	toggle := mods&ui.Cmd != 0
	if mods&ui.Shift != 0 {
		from := s.rowOfKey(s.anchor)
		if from < 0 {
			from = s.rowOf(i)
		}
		s.setSelected(o, s.enabledKeys(from, s.rowOf(i)), toggle)
		return
	}
	s.anchor = key
	if !toggle {
		s.setSelected(o, []string{key}, false)
		return
	}
	if s.Selected.Has(key) {
		s.Selected.Remove(key)
		if o.OnSelectionChange != nil {
			o.OnSelectionChange(s.selectedKeys(SelectionMultiple))
		}
		return
	}
	s.setSelected(o, []string{key}, true)
}

// extend moves the keyboard by step and selects from the anchor to it.
func (s *TreeState) extend(o TreeOptions, step int) {
	e := s.current()
	if e == nil {
		return
	}
	if s.rowOfKey(s.anchor) < 0 {
		s.anchor = e.node.key()
	}
	before := s.cursor
	s.move(step)
	if s.cursor != before {
		s.chooseMultiple(o, s.cursor, ui.Shift)
	}
}

// selectAll chooses every enabled row in view.
func (s *TreeState) selectAll(o TreeOptions) {
	s.setSelected(o, s.enabledKeys(0, len(s.rows)-1), false)
}

// selectForMenu chooses the row of key when the menu opens on it and it is
// not chosen yet. It reports whether the choice or the keyboard's place
// changed: MyGo builds an open menu on every pass, so only a change asks for
// a redraw.
func (s *TreeState) selectForMenu(o TreeOptions, key string) bool {
	changed := false
	switch o.Selection {
	case SelectionSingle:
		changed = s.SelectedKey != key
		s.choose(o, key)
	case SelectionMultiple:
		if !s.Selected.Has(key) {
			changed = true
			s.anchor = key
			s.setSelected(o, []string{key}, false)
		}
	}
	if r := s.rowOfKey(key); r >= 0 && s.cursor != s.rows[r] {
		changed = true
		s.cursor = s.rows[r]
	}
	return changed
}

// activate is a click or Enter on entry i: it chooses the row, or without a
// selection checks it, or opens and closes a branch, and reports the action.
// mods are those of a click, none for Enter.
func (s *TreeState) activate(o TreeOptions, i int32, mods ui.Modifiers) {
	s.toggle(o, i, mods)
	if o.OnAction != nil {
		o.OnAction(s.entries[i].node.key())
	}
}

// toggle does what activate does apart from reporting the action; Space
// takes it for a multiple selection.
func (s *TreeState) toggle(o TreeOptions, i int32, mods ui.Modifiers) {
	e := &s.entries[i]
	s.cursor = i
	switch {
	case o.Selection == SelectionSingle:
		s.choose(o, e.node.key())
	case o.Selection == SelectionMultiple:
		s.chooseMultiple(o, i, mods)
	case o.Checkable:
		s.toggleChecked(o, i)
	case e.branch():
		s.toggleExpanded(o, i)
	}
}

func (s *TreeState) keyAt(row int) any {
	if row < 0 || row >= len(s.rows) {
		return row
	}
	return s.entries[s.rows[row]].node.key()
}

// Tree renders an expandable hierarchy. Rows are virtualized through the
// native List and the rows shown are kept between frames, so a tree of
// thousands of nodes costs what its viewport holds.
//
// The tree is one stop of Tab. Up, Down, Home and End move among the rows,
// the arrow toward the content opens a branch or moves to its first child,
// the other closes it or moves to its parent, Enter activates, Space checks
// a Checkable row, and typing letters moves to the next row starting with
// them. A double click opens and closes a branch.
//
// With SelectionMultiple a click chooses one row, Cmd or Ctrl toggles a row,
// Shift selects the range from the last one clicked, Shift with Up or Down
// extends it, and Cmd or Ctrl with A chooses every enabled row in view.
// Opening a branch fades its rows in over a moment and slides the rows below
// down; closing one shows its arrow turn and drops the rows at once.
func Tree(c *ui.Context, o TreeOptions) ui.Element {
	s := o.State
	if s == nil {
		panic("list: Tree requires a caller-owned State")
	}
	s.sync(o)
	if s.doubled > 0 {
		// Seen while the previous pass built its rows, which the list was
		// counted for: the branch changes here, before this pass counts.
		i := s.doubled - 1
		s.doubled = 0
		if int(i) < len(s.entries) {
			s.toggleExpanded(o, i)
		}
	}
	if s.List.Key == nil {
		s.List.Key = s.keyAt
	}
	for _, r := range s.letters {
		s.typeahead(r, c.Now())
	}
	s.letters = s.letters[:0]
	st := treeStyle{
		dim:     s.dim.value(),
		palette: theme.Colors(c),
		radius:  theme.RadiusScale().SM,
		padX:    theme.Scale(c, theme.Space(2)),
		gap:     theme.Scale(c, theme.Space(1.5)),
		indent:  theme.Scale(c, theme.Space(4)),
		height:  theme.Scale(c, 28),
		mark:    theme.Scale(c, theme.Space(4)),
		rtl:     c.Direction() == ui.RightToLeft,
		guides:  o.Guides,
	}
	if s.revealNew {
		s.revealAt, s.revealNew = c.Now(), false
	}
	if s.reveal >= 0 && c.Now().Sub(s.revealAt) < theme.Motion().Reveal+revealSlack {
		st.reveal, st.reduced = true, theme.ReduceMotion(c)
	}
	var focused, focusKnown bool
	var cursorRow ui.Element
	list := ui.List(c, &s.List, len(s.rows), func(row int) {
		if !focusKnown {
			focused = c.Resolve(s.List.Handle).FocusVisible()
			focusKnown = true
		}
		if e := buildTreeRow(c, s, o, row, st, focused); s.cursor == s.rows[row] {
			cursorRow = e
		}
	}).Padding(theme.Scale(c, theme.Space(1))).Gap(0).MinWidth(0).FillWidth().
		Focusable().FocusRing(false).Role(ui.RoleTree).Label(o.Label).
		Multiselectable(o.Selection == SelectionMultiple)
	if o.Height > 0 {
		list.Height(o.Height)
	}
	if o.MaxHeight > 0 {
		list.MaxHeight(o.MaxHeight)
	}
	if cursorRow.Valid() {
		// The list holds the focus; this tells assistive technology which
		// row the keys are on.
		list.ActiveDescendant(cursorRow)
	}
	list.Busy(o.Loading)
	switch {
	case len(s.rows) == 0 && o.Loading:
		list.Children(func() {
			buildSkeletonRows(c, o.LoadingRows, st.height, st.padX, treeSkeletonIndent)
		})
	case len(s.rows) > 0:
		s.dim.set(c, refreshing(c, list, o.Loading))
	}
	if o.Disabled {
		return list.Disabled(true)
	}
	bindTreeKeys(s, list, o, st.rtl)
	if o.OnContextMenu != nil {
		// Shift+F10 and the menu key open the menu of the row the keys are
		// on; a row's own menu takes the pointer's.
		services := c.Services()
		list.ContextMenu(func(m *ui.Menu) {
			e := s.current()
			if e == nil || e.disabled {
				return
			}
			key := e.node.key()
			if s.selectForMenu(o, key) {
				services.Invalidate()
			}
			o.OnContextMenu(key, m)
		})
	}
	return list
}

// treeStyle holds metrics shared by the rows built in one pass.
type treeStyle struct {
	palette              theme.Palette
	radius, padX, gap    float32
	indent, height, mark float32
	rtl                  bool
	// dim is the rows' opacity under a refresh.
	dim    float32
	guides bool
	// reveal is set while the rows of a branch just opened come in.
	reveal, reduced bool
}

// revealSlack keeps the reveal's transitions on past its duration, so the
// last frame of the motion still finds them.
const revealSlack = 100 * time.Millisecond

func buildTreeRow(c *ui.Context, s *TreeState, o TreeOptions, r int, st treeStyle, focused bool) ui.Element {
	i := s.rows[r]
	e := &s.entries[i]
	n := e.node
	key := n.key()
	open := e.branch() && s.expanded.Has(key)
	selected := s.selected(o.Selection, key)
	start := st.padX + float32(e.depth)*st.indent
	row := ui.Row(c).FillWidth().Shrink(0).Gap(st.gap).AlignItems(ui.Center).
		PaddingStart(start).PaddingEnd(st.padX).
		MinHeight(st.height).Radius(st.radius).
		Role(ui.RoleTreeItem).Label(n.Label).Level(int(e.depth)+1).
		PositionInSet(int(e.pos), int(e.size)).Opacity(st.dim)
	// A tree item's content is not read apart from it, and MyGo keeps one
	// state for chosen and checked: a Checkable row reports its check, which
	// alone can be mixed, and any other row its selection.
	switch {
	case o.Checkable && e.checked > 0 && e.checked < e.leaves:
		row.Mixed()
	case o.Checkable:
		row.Checked(e.checked == e.leaves)
	case o.Selection != SelectionNone:
		row.Selected(selected)
	}
	// Selected paints MyGo's accent row; the Coss face is drawn below.
	row.TextColor(st.palette.Foreground).Background(ui.Transparent)
	if e.branch() {
		row.Expanded(open).Busy(n.Loading)
	}
	enabled := !e.disabled && !o.Disabled
	current := s.cursor == i
	depth := e.depth
	row.Draw(func(p *ui.Painter, box ui.Rect) {
		// Chosen outweighs hovered in both appearances: the input tint is
		// well above the accent one, which reads as a pointer resting.
		switch {
		case selected:
			p.Fill(box, st.palette.Input, st.radius)
		case enabled && row.Hovered():
			p.Fill(box, st.palette.Accent, st.radius)
		}
		if st.guides {
			drawTreeGuides(p, box, depth, st)
		}
		if focused && current {
			p.FocusRing(box, [4]float32{st.radius, st.radius, st.radius, st.radius})
		}
	})
	if st.reveal {
		enter := &ui.Motion{}
		if !st.reduced {
			enter = &ui.Motion{Collapse: true}
		}
		in := s.reveal < i && i < s.entries[s.reveal].end
		m := theme.Motion()
		t := ui.ElementTransition{Duration: m.Reveal, Ease: m.Ease, Position: !st.reduced, Colors: st.reduced}
		if in {
			t.Enter = enter
		}
		row.Transition(t)
	}
	// A press on the arrow or the check box moves the keyboard there too.
	take := func() {
		s.List.Focus(c)
		s.cursor = i
	}
	row.Children(func() {
		buildTreeArrow(c, s, o, i, open, st, take)
		if o.Checkable {
			buildTreeCheck(c, s, o, i, st, take)
		}
		if n.Icon != nil {
			ui.Icon(c, n.Icon).Size(st.mark, st.mark).Shrink(0).TextColor(st.palette.Muted)
		}
		// The fill marks the chosen row: a heavier label would be wider
		// and shift what follows it.
		label := style.Text(ui.Text(c, n.Label), c, theme.SM()).SingleLine().Grow(1).MinWidth(0)
		if w := label.Bounds().W; w > 0 {
			span := ui.Span{Text: n.Label, Font: theme.FontFamily(c), Size: theme.Scale(c, theme.SM().Size)}
			if full, _ := c.MeasureText(0, span); full > w+1 {
				label.Tooltip(n.Label)
			}
		}
	})
	if enabled {
		services := c.Services()
		row.Cursor(ui.CursorPointer)
		row.OnClick(func() {
			s.List.Focus(c)
			s.activate(o, i, row.ClickModifiers())
		})
		// A double click is seen while building, a frame after its clicks
		// acted; the tree opens the branch at the start of the next pass. Where a
		// click already opens and closes, two clicks leave it be.
		if e.branch() && (o.Selection != SelectionNone || o.Checkable) && row.DoubleClicked() {
			s.doubled = i + 1
			services.Invalidate()
		}
		if o.OnContextMenu != nil {
			row.ContextMenu(func(m *ui.Menu) {
				if s.selectForMenu(o, key) {
					services.Invalidate()
				}
				o.OnContextMenu(key, m)
			})
		}
	}
	return style.Disabled(row, !enabled)
}

// drawTreeGuides strokes a hairline through the arrow column of each ancestor
// of a row at depth.
func drawTreeGuides(p *ui.Painter, box ui.Rect, depth int32, st treeStyle) {
	for d := range depth {
		x := box.X + st.padX + float32(d)*st.indent + st.mark/2
		if st.rtl {
			x = box.X + box.W - (st.padX + float32(d)*st.indent + st.mark/2)
		}
		p.Fill(ui.Rect{X: x - theme.Hairline/2, Y: box.Y, W: theme.Hairline, H: box.H}, st.palette.Border, 0)
	}
}

// buildTreeArrow builds a branch's disclosure arrow, or the room a leaf
// keeps so labels line up.
func buildTreeArrow(c *ui.Context, s *TreeState, o TreeOptions, i int32, open bool, st treeStyle, take func()) {
	e := &s.entries[i]
	arrow := ui.Box(c).Size(st.mark, st.mark).Shrink(0).Center().Role(ui.RoleNone)
	if !e.branch() {
		return
	}
	angle := float32(0)
	switch {
	case open:
		angle = 90
	case st.rtl:
		angle = 180
	}
	// The branch's rows come and go at once, so its arrow is what shows the
	// change: it turns over the reveal time, or settles under reduced motion.
	m := theme.Motion()
	angle = arrow.AnimateTo("turn", angle, ui.Animation{Duration: m.Reveal, Ease: m.Ease})
	arrow.Children(func() {
		if e.node.Loading {
			busy.Spinner(c, 12, st.palette.Muted)
			return
		}
		icon := ui.Icon(c, treeArrowIcon).Size(st.mark, st.mark).TextColor(st.palette.Muted)
		if angle != 0 {
			icon.Rotate(angle)
		}
	})
	if !e.disabled && !o.Disabled {
		arrow.OnClick(func() {
			take()
			s.toggleExpanded(o, i)
		})
	}
}

// buildTreeCheck builds entry i's check box: checked when every leaf below
// it is, mixed when some are.
func buildTreeCheck(c *ui.Context, s *TreeState, o TreeOptions, i int32, st treeStyle, take func()) {
	e := &s.entries[i]
	t := st.palette
	on := e.checked == e.leaves
	mixed := e.checked > 0 && !on
	fill, border, ink := t.Background, t.Input, t.PrimaryForeground
	if c.Theme().Dark {
		fill = t.Input.Alpha(0.32)
	}
	if on || mixed {
		fill, border = t.Primary, t.Primary
	}
	radius := theme.Scale(c, 4)
	box := ui.Box(c).Size(st.mark, st.mark).Shrink(0).Center().Role(ui.RoleNone)
	// The face fades with the state; the tick and the dash draw in and fade
	// out on channels of their own, so a branch going from mixed to checked
	// swaps one for the other.
	fill, border = motion.Color(box, "fill", fill), motion.Color(box, "border", border)
	tick := motion.Progress(box, "tick", on)
	dash := motion.Progress(box, "dash", mixed)
	glyph := theme.Scale(c, theme.Space(2.5))
	var path *ui.Path
	if tick > 0 || dash > 0 {
		path = &ui.Path{}
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, fill, radius)
		p.Stroke(r, border, radius, theme.Hairline)
		if path == nil {
			return
		}
		mark := ui.Rect{X: r.X + (r.W-glyph)/2, Y: r.Y + (r.H-glyph)/2, W: glyph, H: glyph}
		drawTick(p, path, mark, tick, on, ink)
		drawDash(p, path, mark, dash, mixed, ink)
	})
	if !e.disabled && !o.Disabled {
		box.OnClick(func() {
			take()
			s.toggleChecked(o, i)
		})
	}
}

// bindTreeKeys registers the tree's keyboard behavior on its focusable list.
func bindTreeKeys(s *TreeState, list ui.Element, o TreeOptions, rtl bool) {
	forth, back := ui.KeyRight, ui.KeyLeft
	if rtl {
		forth, back = back, forth
	}
	list.OnShortcut(0, ui.KeyDown, func() { s.move(1) })
	list.OnShortcut(0, ui.KeyUp, func() { s.move(-1) })
	list.OnShortcut(0, ui.KeyHome, func() { s.edge(false) })
	list.OnShortcut(0, ui.KeyEnd, func() { s.edge(true) })
	list.OnShortcut(0, forth, func() { s.descend(o) })
	list.OnShortcut(0, back, func() { s.ascend(o) })
	list.OnShortcut(0, ui.KeyEnter, func() {
		if s.cursor >= 0 {
			s.activate(o, s.cursor, 0)
		}
	})
	list.OnShortcut(0, ui.KeySpace, func() {
		switch {
		case s.cursor < 0:
		case o.Checkable:
			s.toggleChecked(o, s.cursor)
		default:
			// Space toggles a row of a multiple selection, as Cmd-click does.
			s.toggle(o, s.cursor, ui.Cmd)
		}
	})
	if o.Selection == SelectionMultiple {
		list.OnShortcut(ui.Shift, ui.KeyDown, func() { s.extend(o, 1) })
		list.OnShortcut(ui.Shift, ui.KeyUp, func() { s.extend(o, -1) })
		list.OnShortcut(ui.Cmd, ui.KeyA, func() { s.selectAll(o) })
	}
	// As in ListBox, the page keys scroll on macOS and move a page elsewhere.
	if runtime.GOOS != "darwin" {
		list.OnShortcut(0, ui.KeyPageDown, func() { s.page(1) })
		list.OnShortcut(0, ui.KeyPageUp, func() { s.page(-1) })
	}
	if s.typingOwner != s {
		for i := range typeKeys {
			r := typeRunes[i]
			s.typing[i] = func() { s.letters = append(s.letters, r) }
		}
		s.typingOwner = s
	}
	for i, key := range typeKeys {
		list.OnShortcut(0, key, s.typing[i])
	}
}

var treeArrowIcon = lucide.ChevronRight()

// drawTick strokes a check in r. While on it draws in up to progress of its
// length; while off it fades out, never drawing backwards.
func drawTick(p *ui.Painter, path *ui.Path, r ui.Rect, progress float32, on bool, ink ui.Color) {
	if progress <= 0 {
		return
	}
	path.Reset().
		MoveTo(r.X+r.W*0.12, r.Y+r.H*0.54).
		LineTo(r.X+r.W*0.40, r.Y+r.H*0.80).
		LineTo(r.X+r.W*0.88, r.Y+r.H*0.22)
	strokeMark(p, path, progress, on, ink)
}

// drawDash strokes the mixed mark in r with the rules of drawTick.
func drawDash(p *ui.Painter, path *ui.Path, r ui.Rect, progress float32, on bool, ink ui.Color) {
	if progress <= 0 {
		return
	}
	path.Reset().MoveTo(r.X+r.W*0.15, r.Y+r.H/2).LineTo(r.X+r.W*0.85, r.Y+r.H/2)
	strokeMark(p, path, progress, on, ink)
}

func strokeMark(p *ui.Painter, path *ui.Path, progress float32, on bool, ink ui.Color) {
	switch {
	case !on:
		ink = ink.Alpha(progress)
	case progress < 1:
		path.Trim(0, progress)
	}
	p.StrokePath(path, 1.25, ink)
}
