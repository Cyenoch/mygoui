package navigation

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// menuListRows is the row count, headings included, from which a menu panel
// builds only the rows near its viewport through a native List. Shorter
// menus build every row in a plain scroll, with sections as groups.
const menuListRows = 32

// menuList is the state of a long menu panel: the native list, its rows
// flattened from the sections, and the rows holding and taking the focus.
// The List builds only the rows near its viewport, so the arrows, Home and
// End move the focus here rather than among the rows built.
type menuList struct {
	list     ui.ListState
	rows     []menuRow
	items    int // the item rows, separators aside
	selected []string
	focused  int // the row holding the focus in the last build, -1 for none
	pending  int // the row to focus as it is built, -1 for none
	keyed    bool
	row      int // the row being built
	// decided and long fix the panel's path for one opening: switching it
	// while open would rebuild the rows and drop their focus and marks.
	decided, long bool
}

func newMenuList() *menuList { return &menuList{focused: -1, pending: -1} }

// menuRow is a section heading (item -1) or an item of a section. pos is
// an item's place among the menu's items from 1, 0 for a heading or a
// separator. gap puts the menu's section gap below the last row of a
// section another follows.
type menuRow struct {
	section, item, pos int
	focusable          bool
	gap                bool
}

// menuRowCount returns the rows a menu shows, headings included.
func menuRowCount(o MenuOptions) int {
	n := 0
	for i := range o.Sections {
		n += len(o.Sections[i].Items)
		if o.Sections[i].Title != "" {
			n++
		}
	}
	return n
}

// buildMenuPanel builds a panel's menu in a scroll viewport. A long menu's
// rows are not all built, so its width is fixed at width rather than taken
// from its widest row, which would change as it scrolls.
func buildMenuPanel(c *ui.Context, o MenuOptions, extra menuExtra, parent *menuScope, ml *menuList, width float32) ui.Element {
	if !ml.decided {
		ml.decided, ml.long = true, menuRowCount(o) >= menuListRows
	}
	if !ml.long {
		scroll := ui.Scroll(c)
		scroll.Children(func() { buildMenu(c, o, extra, parent) })
		return scroll
	}
	return buildMenuList(c, o, extra, parent, ml).Width(width)
}

func buildMenuList(c *ui.Context, o MenuOptions, extra menuExtra, parent *menuScope, ml *menuList) ui.Element {
	ml.flatten(o)
	sc := newMenuScope(o, extra, parent)
	if sc.selection != MenuSelectNone && sc.selected == nil {
		sc.selected = &ml.selected
	}
	sc.list, sc.keyed = ml, &ml.keyed
	ml.focused = -1
	list := ui.List(c, &ml.list, len(ml.rows), func(i int) { ml.build(c, o, sc, i) }).
		Padding(scale(c, theme.Space4)).Gap(0).MinWidth(0).
		Role(ui.RoleMenu).FocusGroup(ui.Vertical)
	if o.Label != "" {
		list.Label(o.Label)
	}
	for _, key := range [...]ui.Key{ui.KeyDown, ui.KeyUp, ui.KeyHome, ui.KeyEnd} {
		list.OnShortcut(0, key, func() { ml.move(key) })
	}
	return list
}

// flatten lays the sections out as rows, reusing the rows' storage.
func (ml *menuList) flatten(o MenuOptions) {
	ml.rows, ml.items = ml.rows[:0], 0
	for si := range o.Sections {
		section := &o.Sections[si]
		if section.Title != "" {
			ml.rows = append(ml.rows, menuRow{section: si, item: -1})
		}
		for ii := range section.Items {
			item := &section.Items[ii]
			r := menuRow{section: si, item: ii, focusable: !item.Separator && !item.Disabled}
			if !item.Separator {
				ml.items++
				r.pos = ml.items
			}
			ml.rows = append(ml.rows, r)
		}
		if n := len(ml.rows); n > 0 && ml.rows[n-1].section == si && si < len(o.Sections)-1 {
			ml.rows[n-1].gap = true
		}
	}
}

func (ml *menuList) build(c *ui.Context, o MenuOptions, sc *menuScope, i int) {
	r := ml.rows[i]
	section := &o.Sections[r.section]
	var row ui.Element
	if r.item < 0 {
		row = buildMenuHeading(c, section.Title)
	} else {
		ml.row = i
		row = buildMenuItem(c, section.Items[r.item], sc)
		if r.pos > 0 {
			// Screen readers read the place among all items, built or not.
			row.PositionInSet(r.pos, ml.items)
		}
	}
	if r.gap {
		row.Margin(0, 0, scale(c, theme.Space4), 0)
	}
}

// focus gives the focus to the enabled item row being built when a key moved
// it there, and notes the row holding the focus.
func (ml *menuList) focus(el ui.Element) {
	if ml.pending == ml.row {
		el.Focus()
		ml.pending = -1
	}
	if el.Focused() {
		ml.focused = ml.row
	}
}

// keyedRing reports whether a focused item draws its focus ring itself. The
// runtime shows the ring after focus moves of its own; the moves of a long
// menu keep its focus-visible state as it was, which a pointer-opened menu
// leaves off. A press hands the ring back to the runtime.
func keyedRing(el ui.Element, keyed *bool) bool {
	if el.Pressed() {
		*keyed = false
	}
	return *keyed && el.Focused() && !el.FocusVisible()
}

// move moves the focus as a focus group's arrows, Home and End do, round
// from the last row to the first, over the rows built or not.
func (ml *menuList) move(key ui.Key) {
	from, step := ml.focused, 1
	switch key {
	case ui.KeyUp:
		step = -1
	case ui.KeyHome:
		from = -1
	case ui.KeyEnd:
		from, step = len(ml.rows), -1
	}
	if from < 0 && step < 0 {
		from = len(ml.rows)
	}
	n := len(ml.rows)
	for d := 1; d <= n; d++ {
		i := ((from+step*d)%n + n) % n
		if ml.rows[i].focusable {
			ml.pending, ml.keyed = i, true
			ml.list.ScrollIntoView(i)
			return
		}
	}
}
