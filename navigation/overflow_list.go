package navigation

import (
	"slices"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// OverflowCollapse selects the end of an OverflowList that gives way.
type OverflowCollapse uint8

const (
	// OverflowCollapseEnd moves the last items into the menu first.
	OverflowCollapseEnd OverflowCollapse = iota
	// OverflowCollapseStart moves the first items into the menu first, as a
	// breadcrumb trail keeps its tail.
	OverflowCollapseStart
)

// OverflowItem is one action of an OverflowList. Key identifies it in
// OnSelect and in the tree; it defaults to Label and must be unique.
//
// Checked makes the item a toggle: the pressed button of the row, a check in
// the overflow menu. It is caller-owned and an activation flips it before
// OnSelect runs; nil is a plain action. Pinned items never move into the
// menu. Priority orders the items that do: the lowest goes first, and among
// equal ones the end the list collapses from. Separator draws a rule between
// groups, in the row and in the menu, and ignores the rest.
type OverflowItem struct {
	Key       string
	Label     string
	Icon      *ui.SVG
	Disabled  bool
	Checked   *bool
	Pinned    bool
	Priority  int
	Separator bool
	OnSelect  func()
}

func (it *OverflowItem) key(i int) string {
	switch {
	case it.Separator:
		return "separator-" + strconv.Itoa(i)
	case it.Key != "":
		return it.Key
	}
	return it.Label
}

// OverflowListOptions configures an OverflowList. OnSelect receives the key
// of an activated item, in the row or in the menu, after its own callback.
type OverflowListOptions struct {
	Label    string
	Items    []OverflowItem
	Collapse OverflowCollapse
	OnSelect func(key string)
}

// overflowState is the menu's open flag, and what the row remembers of its
// items between builds: the full width each had, and which of them the last
// build showed as an icon alone, whose width then says nothing of the full one.
type overflowState struct {
	open bool
	full map[string]float32
	icon map[string]bool
}

// overflowMode is how an item shows.
type overflowMode uint8

const (
	// The zero mode is the whole item; overflowIcon is its icon alone.
	overflowIcon overflowMode = iota + 1
	overflowInMenu
)

// OverflowList builds a row of actions that shows as many as fit its width.
// When they do not fit, the items with an icon first give up their labels,
// then the rest move into a trailing "more" menu, or a leading one when
// collapsing from the start.
//
// The row measures itself from the layout of the frame before and asks for
// another frame when its width or an item's width changed, so a resize
// settles in the next frame; meanwhile the row clips what does not fit
// rather than painting past its bounds.
func OverflowList(c *ui.Context, o OverflowListOptions) ui.Element {
	gap := scale(c, theme.Space(1))
	root := ui.Row(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).Gap(gap).Clip().
		Role(ui.RoleToolbar).FocusGroup(ui.Horizontal)
	if o.Label != "" {
		root.Label(o.Label)
	}
	st := ui.Local(root, "overflow-list", func() overflowState {
		return overflowState{full: map[string]float32{}, icon: map[string]bool{}}
	})
	avail := root.LayoutBounds().W
	root.DrawOver(func(_ *ui.Painter, r ui.Rect) { invalidateOnResize(c, avail, r.W) })

	fromStart := o.Collapse == OverflowCollapseStart
	root.Children(func() {
		var more ui.Element
		if fromStart {
			more = buildOverflowMore(c, st)
		}
		n := len(o.Items)
		items := make([]ui.Element, n)
		full := make([]float32, n)
		keys := make([]string, n)
		for i := range o.Items {
			keys[i] = o.Items[i].key(i)
			items[i] = buildOverflowItem(c, o, i, keys[i], st.icon[keys[i]])
			w := items[i].LayoutBounds().W
			if !st.icon[keys[i]] && w > 0 {
				st.full[keys[i]] = w
			}
			full[i] = st.full[keys[i]]
			items[i].DrawOver(func(_ *ui.Painter, r ui.Rect) { invalidateOnResize(c, w, r.W) })
		}
		if !fromStart {
			more = buildOverflowMore(c, st)
		}
		modes, measured := overflowPlan(o.Items, full, more.LayoutBounds().W, avail, gap, fromStart)
		menu := false
		for i, mode := range modes {
			st.icon[keys[i]] = mode == overflowIcon
			if mode == overflowInMenu {
				items[i].Absolute().Invisible()
				menu = menu || !o.Items[i].Separator
			}
		}
		if !menu {
			more.Absolute().Invisible()
			if measured {
				st.open = false
			}
			return
		}
		align := mygoui.AlignEnd
		if fromStart {
			align = mygoui.AlignStart
		}
		var m MenuOptions
		if st.open {
			m = overflowMenu(o, modes)
		}
		menuPopover(c, more, &st.open, m, mygoui.Placement{Align: align}, true)
	})
	return root
}

// invalidateOnResize asks for another build when a width laid out in this
// frame differs from the one the build measured with.
func invalidateOnResize(c *ui.Context, built, laidOut float32) {
	if d := built - laidOut; d < -0.5 || d > 0.5 {
		c.Invalidate()
	}
}

// overflowPlan decides how each item shows in avail, from the full width of
// each and that of the more button, which an icon alone has too. Everything
// shows, unmeasured, while a width is still unknown: in the first frame, or
// the one that adds an item. Otherwise the labels of items with an icon go
// first, then items move into the menu: the lowest priority first, from the
// collapsing end among equals, never a pinned one.
func overflowPlan(items []OverflowItem, full []float32, more, avail, gap float32, fromStart bool) (modes []overflowMode, measured bool) {
	n := len(items)
	modes = make([]overflowMode, n)
	if avail <= 0 || more <= 0 {
		return modes, false
	}
	for _, w := range full {
		if w <= 0 {
			return modes, false
		}
	}
	separator := make([]bool, n)
	for i := range items {
		separator[i] = items[i].Separator
	}
	shown := make([]bool, n)
	fits := func() bool {
		total, count := float32(0), 0
		menu := false
		for i, mode := range modes {
			shown[i] = mode != overflowInMenu
			menu = menu || !shown[i]
		}
		settleSeparators(separator, shown)
		for i := range items {
			switch {
			case !shown[i]:
			case modes[i] == overflowIcon:
				total, count = total+more, count+1
			default:
				total, count = total+full[i], count+1
			}
		}
		if menu {
			total, count = total+more, count+1
		}
		return total+gap*float32(max(count-1, 0)) <= avail+0.5
	}
	if fits() {
		return withSeparators(modes, shown, separator), true
	}
	for i := range items {
		if items[i].Icon != nil && items[i].Label != "" && !items[i].Separator {
			modes[i] = overflowIcon
		}
	}
	if fits() {
		return withSeparators(modes, shown, separator), true
	}
	order := make([]int, 0, n)
	for i := range items {
		if !items[i].Pinned && !items[i].Separator {
			order = append(order, i)
		}
	}
	slices.SortStableFunc(order, func(a, b int) int {
		if items[a].Priority != items[b].Priority {
			return items[a].Priority - items[b].Priority
		}
		if fromStart {
			return a - b
		}
		return b - a
	})
	for _, i := range order {
		modes[i] = overflowInMenu
		if fits() {
			break
		}
	}
	return withSeparators(modes, shown, separator), true
}

// withSeparators sets the mode of each rule to whether it shows in the row,
// as the last fits call settled it: a rule that does not stands for nothing,
// and counts as out of the row.
func withSeparators(modes []overflowMode, shown, separator []bool) []overflowMode {
	for i := range modes {
		if separator[i] && !shown[i] {
			modes[i] = overflowInMenu
		}
	}
	return modes
}

// buildOverflowItem builds one item of the row: a button, a toggle where it
// has a pressed state, or a rule. The icon of an item that gave up its label
// keeps the label for assistive technology and a tooltip.
func buildOverflowItem(c *ui.Context, o OverflowListOptions, i int, key string, iconOnly bool) ui.Element {
	it := &o.Items[i]
	c.Key(key)
	if it.Separator {
		return ui.Box(c).Width(hairline).Height(scale(c, theme.Space(4))).Shrink(0).
			Margin(0, scale(c, theme.Space(1))).Background(theme.Colors(c).Border).Role(ui.RoleNone)
	}
	if it.Checked != nil {
		return mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
			Label: it.Label, Icon: it.Icon, Size: theme.ButtonSmall, IconOnly: iconOnly,
			Selected: it.Checked, Disabled: it.Disabled, OnChange: func(bool) { selectOverflowItem(o, it, key) },
		})
	}
	return mygoui.Button(c, mygoui.ButtonOptions{
		Label: it.Label, Icon: it.Icon, IconOnly: iconOnly, Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
		Disabled: it.Disabled, OnClick: func() { selectOverflowItem(o, it, key) },
	})
}

func selectOverflowItem(o OverflowListOptions, it *OverflowItem, key string) {
	if it.OnSelect != nil {
		it.OnSelect()
	}
	if o.OnSelect != nil {
		o.OnSelect(key)
	}
}

// buildOverflowMore builds the button opening the menu of collapsed items.
func buildOverflowMore(c *ui.Context, st *overflowState) ui.Element {
	c.Key("overflow-more")
	more := mygoui.Button(c, mygoui.ButtonOptions{
		Label: locale.M().More, Icon: lucide.Ellipsis(), IconOnly: true,
		Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
		OnClick: func() { st.open = !st.open },
	})
	return more.Role(ui.RoleMenuButton).Expanded(st.open)
}

// overflowMenu lists the items moved out of the row, in their order, with a
// rule where a group ends between them. A toggle is a check item in the menu,
// checked while pressed; plain actions stay plain items.
func overflowMenu(o OverflowListOptions, modes []overflowMode) MenuOptions {
	hidden := make([]MenuItemOptions, 0, len(o.Items))
	shown := make([]bool, len(o.Items))
	separator := make([]bool, len(o.Items))
	for i := range o.Items {
		it := &o.Items[i]
		separator[i] = it.Separator
		shown[i] = it.Separator || modes[i] == overflowInMenu
	}
	settleSeparators(separator, shown)
	for i := range o.Items {
		it := &o.Items[i]
		if !shown[i] {
			continue
		}
		if it.Separator {
			hidden = append(hidden, MenuItemOptions{Separator: true})
			continue
		}
		key := it.key(i)
		hidden = append(hidden, MenuItemOptions{
			Key: key, Label: it.Label, Icon: it.Icon, Disabled: it.Disabled,
			Checked: it.Checked != nil && *it.Checked, Checkable: it.Checked != nil,
			OnSelect: func() {
				if it.Checked != nil {
					*it.Checked = !*it.Checked
				}
				if it.OnSelect != nil {
					it.OnSelect()
				}
			},
		})
	}
	if slices.ContainsFunc(hidden, func(m MenuItemOptions) bool { return m.Checkable }) {
		for i := range hidden {
			hidden[i].inset = !hidden[i].Checkable && !hidden[i].Separator
		}
	}
	return MenuOptions{
		Label:    locale.M().More,
		Sections: []MenuSectionOptions{{Items: hidden}},
		OnSelect: o.OnSelect,
	}
}
