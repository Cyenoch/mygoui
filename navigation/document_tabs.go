package navigation

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DocumentTab is one open document. Key is its stable identity, unique in
// the strip: the tab keeps its focus as the order changes. Dirty marks
// unsaved changes. Tooltip is shown on pointer rest, as the document's full
// path. Pinned shows the tab as its icon alone, with no close button; the
// caller keeps pinned tabs first, and dragging never takes a tab across the
// boundary.
type DocumentTab struct {
	Key     string
	Title   string
	Icon    *ui.SVG
	Dirty   bool
	Pinned  bool
	Tooltip string
}

// DocumentTabSizing is how tabs take the strip's width.
type DocumentTabSizing uint8

const (
	// DocumentTabsFixed gives each tab the width of its title, up to a limit.
	DocumentTabsFixed DocumentTabSizing = iota
	// DocumentTabsShrink narrows the tabs to share the strip, as a browser
	// does, and scrolls only once they reach their minimum.
	DocumentTabsShrink
)

// DocumentTabsOptions configures a document tab strip. Tabs and Active are
// caller-owned and never modified: the strip reports what the user asked for
// and the caller applies it. OnMove reports that the tab at from goes to
// index to, counted after its removal. A nil OnClose, OnMove or OnAdd leaves
// that ability out. OnCloseMany receives the keys of the tabs the context
// menu's Close others and Close to the right ask for, which never include a
// pinned tab, in one call; without it each key goes to OnClose. Ctrl+Tab,
// Ctrl+Shift+Tab, Ctrl+PageDown and Ctrl+PageUp choose the next and previous
// tab and Cmd/Ctrl+W closes the active one while the focus is in the strip;
// WindowShortcuts has them work wherever the focus is in the window, for a
// window with one strip.
type DocumentTabsOptions struct {
	Tabs            []DocumentTab
	Active          string
	Label           string
	Sizing          DocumentTabSizing
	WindowShortcuts bool
	OnSelect        func(key string)
	OnClose         func(key string)
	OnCloseMany     func(keys []string)
	OnMove          func(from, to int)
	OnAdd           func()
}

type (
	docTabsStateKey struct{}
	docTabShiftKey  struct{}
)

// docTabsState is the drag in progress, as SortableList keeps it, and what
// the strip remembers between builds.
type docTabsState struct {
	held    int
	heldKey string
	dx      float32
	bounds  []ui.Rect
	// settle places the tabs at once in the build after a drop, where the
	// caller's new order already holds them.
	settle bool
	// shown is the active tab the last build showed.
	shown string
	// middle is the tab the middle button is pressing while middleDown.
	middle     string
	middleDown bool
	// refocus is the place of a focused tab that asked to close, and
	// refocusKey its key: once the caller removed it, the tab that took its
	// place takes the focus.
	refocus    int
	refocusKey string

	// scroll is the viewport's scrolling; maxX is the extent the last build
	// saw, so a change of it builds again.
	scroll ui.ScrollState
	maxX   float32
	// crowded is whether the last build offered the all-tabs menu.
	crowded bool
	// listOpen is the all-tabs menu; reveal is the tab to bring into view
	// until it shows, within revealTries builds; focusKey is the tab a key
	// chose, which takes the focus.
	listOpen    bool
	reveal      string
	revealTries int
	focusKey    string
	// addW is the width the add button laid out with.
	addW float32
	// closeMany is the keys a context menu asked to close, applied after the
	// build.
	closeMany []string
	// menuTab is the tab whose context menu asked for the closes: the menu held
	// the focus, which its tab takes back.
	menuTab string
}

// docTabInput is what one tab's build read from the pointer and the keys.
type docTabInput struct {
	dx       float32
	held     bool
	step     int
	selected bool
	closed   bool
	middle   bool
}

var docTabShift = ui.Animation{Duration: 150 * time.Millisecond, Ease: ui.EaseOut}

// A tab opened or closed grows from nothing along the strip and fades. The
// copy of a closed tab takes no room, so the tabs beside it slide into its
// place on a transition of their own: Position follows the layout, except in
// the build after a drop, where the caller's new order already holds them.
var (
	docTabGone   = ui.Motion{Collapse: true}
	docTabFade   = ui.Motion{}
	docTabEase   = ui.EaseOut
	docTabLayout = ui.ElementTransition{
		// Colors and shadows are left out: they fade on the tab's own
		// channel and would be eased twice.
		Position: true, Duration: motion.Reveal, Ease: docTabEase,
		Enter: &docTabGone, Exit: &docTabGone,
	}
	// docTabStrip makes the viewport and the add button follow the tabs, so
	// the viewport does not clip the ones sliding in or out.
	docTabStrip    = ui.ElementTransition{Position: true, Width: true, Duration: motion.Reveal, Ease: docTabEase}
	docTabFadeOnly = ui.ElementTransition{
		Duration: motion.State, Ease: docTabEase, Enter: &docTabFade, Exit: &docTabFade,
	}
)

// DocumentTabs builds a strip of document tabs. A click, Enter or Space
// selects a tab; the arrows, Home and End move the focus along the strip,
// which is one focus stop. A tab closes from its close button, a middle
// click, Cmd/Ctrl+W or Delete, and moves by dragging it or with Alt+Left and
// Alt+Right; letting go of a dragged tab selects it, as a click does. A
// secondary click opens a menu to close the tab, the others or those after
// it. When the caller removes a focused tab, the tab that takes its place
// takes the focus.
//
// The strip scrolls when the tabs overflow: its edges fade where more tabs
// lie beyond, the wheel scrolls it sideways, a menu lists every tab, and a
// newly active tab is brought into view.
func DocumentTabs(c *ui.Context, o DocumentTabsOptions) ui.Element {
	t := theme.Colors(c)
	gap := scale(c, theme.Space2)
	root := ui.Row(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).Gap(gap).
		Padding(gap).Radius(theme.RadiusScale().LG).Background(t.MutedBackground)
	st := ui.Local(root, docTabsStateKey{}, func() docTabsState { return docTabsState{held: -1, refocus: -1} })
	n := len(o.Tabs)
	// Pinned tabs lead: a tab is dragged within its own group.
	split := 0
	for split < n && o.Tabs[split].Pinned {
		split++
	}
	group := func(i int) (lo, hi int) {
		if i < split {
			return 0, split - 1
		}
		return split, n - 1
	}
	if st.held >= n || len(st.bounds) != n || o.OnMove == nil ||
		(st.held >= 0 && o.Tabs[st.held].Key != st.heldKey) {
		st.held = -1
	}
	over := st.held
	if st.held >= 0 {
		lo, hi := group(st.held)
		st.dx = docTabClamp(st.bounds, st.held, st.dx, lo, hi)
		over = docTabOver(st.bounds, st.held, st.dx, lo, hi)
	}
	settle := st.settle || theme.ReduceMotion(c)
	st.settle = false
	if len(st.bounds) != n {
		st.bounds = make([]ui.Rect, n)
	}
	if st.shown != o.Active {
		st.reveal, st.revealTries = o.Active, 0
	}
	st.shown = o.Active
	focusKey := st.focusKey
	st.focusKey = ""
	// focusAt is the tab that takes the focus of one the caller just closed.
	focusAt := -1
	if st.refocus >= 0 && n > 0 && (st.refocus >= n || o.Tabs[st.refocus].Key != st.refocusKey) {
		focusAt = min(st.refocus, n-1)
	}
	st.refocus = -1
	// The theme's muted ink is short of 4.5:1 on the strip in the dark
	// appearance.
	rest := legibleInk(t.MutedForeground, t.Foreground, t.MutedBackground.Over(t.Background))
	// forward is the screen direction of the next tab.
	forward := 1
	if c.Direction() == ui.RightToLeft {
		forward = -1
	}
	// Under DocumentTabsShrink the tabs share what the layout of the last
	// frame left them.
	width, crowded := float32(0), false
	if o.Sizing == DocumentTabsShrink {
		laidOut := root.LayoutBounds().W
		width, crowded = docTabSharedWidth(c, st, laidOut, split, n-split, gap, o.OnAdd != nil)
		root.DrawOver(func(_ *ui.Painter, r ui.Rect) { invalidateOnResize(c, laidOut, r.W) })
	}
	overflowing := st.scroll.MaxX > 0.5
	if o.Sizing == DocumentTabsShrink {
		overflowing = crowded
	}
	if d := st.scroll.MaxX - st.maxX; d > 0.5 || d < -0.5 {
		st.maxX = st.scroll.MaxX
		c.Invalidate()
	}
	// The strip narrows as its menu button comes, which can push the active
	// tab out of view.
	if overflowing != st.crowded {
		st.crowded = overflowing
		st.reveal, st.revealTries = o.Active, 0
	}
	var pick func(i int)
	if o.OnSelect != nil && n > 1 {
		// A window shortcut chooses a tab without taking the focus from the
		// editor: the focus follows only where the strip already had it.
		focusIn := root.FocusWithin()
		pick = func(i int) {
			key := o.Tabs[i].Key
			if focusIn {
				st.focusKey = key
			}
			o.OnSelect(key)
		}
	}
	active := -1
	for i := range o.Tabs {
		if o.Tabs[i].Key == o.Active {
			active = i
		}
	}
	docTabShortcuts(c, root, o, active, pick, st)

	selected, closed, from, to := -1, -1, -1, -1
	// focused is the tab holding the focus.
	focused, focusedAt := "", -1
	pressing := false
	root.Children(func() {
		viewport := ui.ScrollHorizontal(c.Key("strip")).MinWidth(0).Shrink(1).Scrollbars(false).TrackScroll(&st.scroll)
		view := viewport.LayoutBounds()
		viewport.HandleInput(func(ev ui.InputEvent) bool { return docTabWheel(&st.scroll, ev) })
		docTabEdges(c, viewport, &st.scroll)
		if !settle {
			// The viewport follows its tabs, so it does not clip the ones
			// sliding in or out before they get there.
			viewport.Transition(docTabStrip)
		}
		viewport.Children(func() {
			list := ui.Row(c).Shrink(0).Gap(gap).Role(ui.RoleTabList).FocusGroup(ui.Horizontal)
			if o.Label != "" {
				list.Label(o.Label)
			}
			list.Children(func() {
				for i := range o.Tabs {
					tab := &o.Tabs[i]
					isActive := tab.Key == o.Active
					el, in := docTab(c.Key(tab.Key), tab, docTabView{
						rest: rest, active: isActive, lifted: i == st.held, settle: settle, forward: forward,
						width: width, index: i, count: n, many: &st.closeMany, origin: &st.menuTab,
					}, o)
					if tab.Key == st.reveal {
						if docTabShows(view, st.bounds[i], scale(c, theme.Space24), &st.scroll, forward < 0) || st.revealTries >= 6 {
							st.reveal = ""
						} else {
							st.revealTries++
							el.ScrollIntoView(docTabScrollTo(settle, view.W > 0))
							c.Invalidate()
						}
					}
					if i == focusAt || tab.Key == focusKey {
						el.Focus()
					}
					if i == st.held {
						st.dx += in.dx
					}
					switch shift := docTabShiftOf(st.bounds, st.held, over, i, gap) * float32(forward); {
					case i == st.held:
						el.Translate(st.dx, 0).ZIndex(1)
					case settle:
						el.Translate(el.Jump(docTabShiftKey{}, shift), 0)
					default:
						el.Translate(el.AnimateTo(docTabShiftKey{}, shift, docTabShift), 0)
					}
					switch {
					case in.held && st.held < 0:
						st.held, st.heldKey, st.dx = i, tab.Key, 0
					case !in.held && st.held == i:
						from, to = i, over
						st.held, st.dx, st.settle = -1, 0, true
					case in.step != 0 && st.held < 0:
						lo, hi := group(i)
						if j := i + in.step; j >= lo && j <= hi {
							from, to = i, j
						}
					}
					switch {
					case in.middle:
						st.middle, pressing = tab.Key, true
					case st.middleDown && st.middle == tab.Key && el.Hovered():
						// The middle button let go over the tab it pressed.
						in.closed = !tab.Pinned
					}
					if el.FocusWithin() {
						focused, focusedAt = tab.Key, i
					}
					if in.closed {
						closed = i
						if el.FocusWithin() {
							st.refocus, st.refocusKey = i, tab.Key
						}
					} else if in.selected {
						selected = i
					}
					st.bounds[i] = el.LayoutBounds()
				}
			})
		})
		if overflowing {
			docTabAllTabs(c, st, o, pick)
		}
		if o.OnAdd != nil {
			add := mygoui.Button(c.Key("add"), mygoui.ButtonOptions{
				Label: locale.M().Add, Icon: lucide.Plus(), IconOnly: true,
				Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, OnClick: o.OnAdd,
			})
			st.addW = add.LayoutBounds().W
			if !settle {
				add.Transition(docTabStrip)
			}
		}
	})
	st.middleDown = pressing
	// The caller changes its tabs here, after they were built.
	switch many := st.closeMany; {
	case closed >= 0:
		st.closeMany, st.menuTab = nil, ""
		o.OnClose(o.Tabs[closed].Key)
	case len(many) > 0:
		st.closeMany = nil
		if st.menuTab != "" {
			focused, focusedAt = st.menuTab, slices.IndexFunc(o.Tabs, func(tab DocumentTab) bool { return tab.Key == st.menuTab })
			st.menuTab = ""
		}
		if slices.Contains(many, focused) {
			st.refocus, st.refocusKey = focusedAt, focused
		}
		docTabCloseKeys(o, many)
	case selected >= 0 && o.OnSelect != nil:
		o.OnSelect(o.Tabs[selected].Key)
	}
	if closed < 0 && from != to {
		o.OnMove(from, to)
	}
	if closed >= 0 || selected >= 0 || from != to {
		c.Invalidate()
	}
	return root
}

// docTabView is how one tab is built: its place in the strip and the strip's
// state around it.
type docTabView struct {
	index, count   int
	rest           ui.Color
	active, lifted bool
	settle         bool
	forward        int
	// width is the tab's width under DocumentTabsShrink, 0 for its own.
	width float32
	// many collects the keys a context menu asks to close.
	many   *[]string
	origin *string
}

// docTabSharedWidth returns the width each unpinned tab takes to share the
// strip, and whether the tabs are crowded: at their minimum, with more than
// the strip holds.
func docTabSharedWidth(c *ui.Context, st *docTabsState, rootW float32, pinned, tabs int, gap float32, add bool) (width float32, crowded bool) {
	if rootW <= 0 || tabs == 0 {
		return 0, false
	}
	room := rootW - 2*gap - float32(pinned)*(scale(c, theme.HeightButtonSM)+gap)
	if add {
		room -= st.addW + gap
	}
	fit := (room - float32(tabs-1)*gap) / float32(tabs)
	least := scale(c, 96)
	return min(max(fit, least), scale(c, 200)), fit < least-0.5
}

// docTabShortcuts binds the keys that choose and close tabs.
func docTabShortcuts(c *ui.Context, root ui.Element, o DocumentTabsOptions, active int, pick func(int), st *docTabsState) {
	bind := func(mods ui.Modifiers, key ui.Key, fn func()) {
		if o.WindowShortcuts {
			c.OnWindowShortcut(mods, key, fn)
		} else {
			root.OnShortcut(mods, key, fn)
		}
	}
	if pick != nil {
		n := len(o.Tabs)
		step := func(by int) func() {
			return func() {
				switch {
				case active >= 0:
					pick(((active+by)%n + n) % n)
				case by > 0:
					pick(0)
				default:
					pick(n - 1)
				}
			}
		}
		bind(ui.Ctrl, ui.KeyTab, step(1))
		bind(ui.Ctrl|ui.Shift, ui.KeyTab, step(-1))
		bind(ui.Ctrl, ui.KeyPageDown, step(1))
		bind(ui.Ctrl, ui.KeyPageUp, step(-1))
	}
	if o.OnClose != nil && active >= 0 && !o.Tabs[active].Pinned {
		key := o.Tabs[active].Key
		bind(ui.Cmd, ui.KeyW, func() { st.closeMany = append(st.closeMany, key) })
	}
}

// docTabCloseKeys hands the keys a menu or a key asked to close to the
// caller: one at a time, or together where it takes them together.
func docTabCloseKeys(o DocumentTabsOptions, keys []string) {
	if o.OnCloseMany != nil && len(keys) > 1 {
		o.OnCloseMany(keys)
		return
	}
	for _, key := range keys {
		o.OnClose(key)
	}
}

// docTabWheel scrolls the strip sideways by the wheel's vertical turn, as a
// browser's tab strip does. Sideways scrolling is the scroll container's own.
func docTabWheel(s *ui.ScrollState, ev ui.InputEvent) bool {
	if ev.Kind != ui.InputScroll || ev.Mods != 0 || s.MaxX <= 0 ||
		math.Abs(float64(ev.DY)) <= math.Abs(float64(ev.DX)) {
		return false
	}
	to := min(max(s.X+ev.DY, 0), s.MaxX)
	if ev.Precise {
		s.X = to
	} else {
		s.ScrollTo(to, 0, ui.Animation{Duration: motion.State, Ease: ui.EaseOut})
	}
	return true
}

// docTabEdges fades the strip's edges where tabs lie beyond them, by how far
// they lie.
func docTabEdges(c *ui.Context, viewport ui.Element, s *ui.ScrollState) {
	w := viewport.LayoutBounds().W
	reach := scale(c, theme.Space24)
	if s.MaxX <= 0.5 || w <= 2*reach {
		return
	}
	before := min(s.X/reach, 1)
	after := min((s.MaxX-s.X)/reach, 1)
	if c.Direction() == ui.RightToLeft {
		before, after = after, before
	}
	if before <= 0 && after <= 0 {
		return
	}
	ink, edge := ui.RGB(0, 0, 0), reach/w
	viewport.Masks(ui.Mask{Image: ui.LinearGradient{Angle: 90, Stops: []ui.GradientStop{
		{Color: ink.Alpha(1 - before), At: 0}, {Color: ink, At: edge},
		{Color: ink, At: 1 - edge}, {Color: ink.Alpha(1 - after), At: 1},
	}}})
}

// docTabShows reports whether a tab lies in the strip clear of its faded
// edges, as the last layout had it: an edge fades only where the strip can
// scroll on.
func docTabShows(view, tab ui.Rect, edge float32, s *ui.ScrollState, rtl bool) bool {
	before, after := edge, edge
	if s.X <= 0.5 {
		before = 0
	}
	if s.X >= s.MaxX-0.5 {
		after = 0
	}
	if rtl {
		before, after = after, before
	}
	return tab.W > 0 && tab.X >= view.X+before && tab.X+tab.W <= view.X+view.W-after
}

// docTabScrollTo centers a tab in the strip, gliding unless it is the first
// layout or motion is reduced.
func docTabScrollTo(still, laidOut bool) ui.ScrollIntoViewOptions {
	o := ui.ScrollIntoViewOptions{Inline: ui.ScrollCenter}
	if !still && laidOut {
		o.Animation = &ui.Animation{Duration: motion.Reveal, Ease: ui.EaseOut}
	}
	return o
}

// docTabAllTabs builds the button of the menu that lists every tab and
// chooses one, and scrolls it into view.
func docTabAllTabs(c *ui.Context, st *docTabsState, o DocumentTabsOptions, pick func(int)) {
	button := mygoui.Button(c.Key("all-tabs"), mygoui.ButtonOptions{
		Label: locale.M().AllTabs, Icon: lucide.ChevronDown(), IconOnly: true,
		Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
		OnClick: func() { st.listOpen = !st.listOpen },
	}).Role(ui.RoleMenuButton).Expanded(st.listOpen)
	var menu MenuOptions
	if st.listOpen {
		items := make([]MenuItemOptions, len(o.Tabs))
		for i, tab := range o.Tabs {
			items[i] = MenuItemOptions{Key: tab.Key, Label: tab.Title, Icon: tab.Icon}
		}
		menu = MenuOptions{
			Label: locale.M().AllTabs, Selection: MenuSelectSingle, Selected: &[]string{o.Active},
			Sections: []MenuSectionOptions{{Items: items}},
			OnSelect: func(key string) {
				st.reveal, st.revealTries = key, 0
				if o.OnSelect != nil {
					o.OnSelect(key)
				}
			},
		}
	}
	menuPopover(c, button, &st.listOpen, menu, mygoui.Placement{Align: mygoui.AlignEnd})
}

// docTab builds one tab and reads its input.
func docTab(c *ui.Context, tab *DocumentTab, v docTabView, o DocumentTabsOptions) (ui.Element, docTabInput) {
	t := theme.Colors(c)
	closable := o.OnClose != nil && !tab.Pinned
	height := scale(c, theme.HeightButtonSM)
	el := ui.ButtonBase(c).Role(ui.RoleTab).AccessSelected(v.active).Label(tab.Title).
		Row().AlignItems(ui.Center).Gap(scale(c, theme.Space6)).Shrink(0).
		Height(height).MaxWidth(scale(c, 200)).
		PaddingStart(scale(c, theme.Space10)).PaddingEnd(scale(c, theme.Space10)).
		Radius(theme.RadiusScale().MD).FocusRing(true)
	switch {
	case tab.Pinned:
		el.Width(height).Padding(0).Center()
	case v.width > 0:
		el.Width(v.width)
	}
	// The Exit stays declared in every build, or the tab would count as gone.
	switch {
	case theme.ReduceMotion(c):
		el.Transition(docTabFadeOnly)
	case v.settle:
		still := docTabLayout
		still.Position = false
		el.Transition(still)
	default:
		layout := docTabLayout
		layout.Width = v.width > 0
		el.Transition(layout)
	}
	if !tab.Pinned && (closable || tab.Dirty) {
		el.PaddingEnd(scale(c, theme.Space4))
	}
	if tab.Dirty {
		el.Description(locale.M().Modified)
	}
	switch {
	case tab.Tooltip != "":
		el.Tooltip(tab.Tooltip)
	case tab.Pinned:
		el.Tooltip(tab.Title)
	}
	// Snapshot hover: transitions can retain paint after the element expires.
	hovered := el.Hovered()
	fg := selectionInk(el, "tab-ink", v.active, hovered || v.lifted,
		v.rest, t.Foreground, t.Foreground)
	el.TextColor(fg)
	// The active tab's face fades in and out with its shadow; a tab picked
	// up takes it at once.
	raised := motion.Progress(el, "face", v.active)
	if v.lifted {
		raised = 1
	}
	if raised > 0 {
		el.Background(segmentBackground(c, t).Alpha(raised))
		el = fadeShadows(el, theme.SurfaceShadows(c.Theme().Dark), raised)
	}
	var in docTabInput
	el.Children(func() {
		switch {
		case tab.Icon != nil:
			glyph := scale(c, theme.IconSM)
			ui.Icon(c, tab.Icon).Size(glyph, glyph).Shrink(0)
		case tab.Pinned:
			style.Text(ui.Text(c, docTabInitial(tab.Title)), c, theme.SM()).FontWeight(theme.WeightMedium)
		}
		if tab.Pinned {
			if tab.Dirty {
				dot := scale(c, theme.Space6)
				ui.Box(c).Absolute().Top(scale(c, theme.Space4)).Right(scale(c, theme.Space4)).
					Size(dot, dot).Radius(dot / 2).Background(fg)
			}
			return
		}
		style.Text(ui.Text(c, tab.Title), c, theme.SM()).FontWeight(theme.WeightMedium).
			SingleLine().MinWidth(0).Shrink(1)
		if closable || tab.Dirty {
			in.closed = docTabClose(c, tab, closable, v.active || hovered, fg)
		}
	})
	if o.OnClose != nil {
		in.closed = in.closed || (!tab.Pinned && el.Shortcut(0, ui.KeyDelete))
	}
	if closable {
		button, pressed := el.PressedButton()
		in.middle = pressed && button == 2
	}
	in.selected = !v.active && !in.closed && el.Clicked()
	if o.OnMove != nil {
		in.dx, _, in.held = el.Dragged()
		switch {
		case el.Shortcut(ui.Alt, ui.KeyLeft):
			in.step = -v.forward
		case el.Shortcut(ui.Alt, ui.KeyRight):
			in.step = v.forward
		}
	}
	if o.OnClose != nil {
		PopoverMenu(el, docTabMenu(o, tab, v))
	}
	return el, in
}

// docTabInitial is what a pinned tab without an icon shows.
func docTabInitial(title string) string {
	for _, r := range title {
		return strings.ToUpper(string(r))
	}
	return ""
}

// docTabMenu is the secondary-click menu of one tab. Its choices are kept
// for the end of the build, as the other closes are.
func docTabMenu(o DocumentTabsOptions, tab *DocumentTab, v docTabView) MenuOptions {
	ask := func(from, to int, skipThis bool) func() {
		return func() {
			*v.origin = tab.Key
			var keys []string
			for i := from; i < to; i++ {
				if other := &o.Tabs[i]; !other.Pinned && (!skipThis || i != v.index) {
					keys = append(keys, other.Key)
				}
			}
			*v.many = append(*v.many, keys...)
		}
	}
	m := locale.M()
	return MenuOptions{Sections: []MenuSectionOptions{{Items: []MenuItemOptions{
		{Key: "close", Label: m.Close, Accelerator: form.Hotkey{Mods: ui.Cmd, Key: ui.KeyW},
			OnSelect: func() { *v.origin, *v.many = tab.Key, append(*v.many, tab.Key) }},
		{Key: "others", Label: m.CloseOthers, Disabled: !docTabHasOthers(o.Tabs, v.index),
			OnSelect: ask(0, v.count, true)},
		{Key: "right", Label: m.CloseToRight, Disabled: !docTabHasOthers(o.Tabs[v.index+1:], -1),
			OnSelect: ask(v.index+1, v.count, false)},
	}}}}
}

// docTabHasOthers reports whether an unpinned tab other than the one at skip
// is among tabs.
func docTabHasOthers(tabs []DocumentTab, skip int) bool {
	for i := range tabs {
		if i != skip && !tabs[i].Pinned {
			return true
		}
	}
	return false
}

// docTabClose builds the slot at a tab's end: the close button while the
// tab is active or hovered, and the dot of unsaved changes until the pointer
// is on it. It reports a press of the close button.
func docTabClose(c *ui.Context, tab *DocumentTab, closable, offered bool, ink ui.Color) (closed bool) {
	t := theme.Colors(c)
	side := scale(c, theme.Space20)
	slot := ui.Box(c).Size(side, side).Center().Shrink(0)
	if closable {
		// The slot stays out of the focus order, so the arrows go from tab
		// to tab; the keys close the focused tab.
		slot.Role(ui.RoleButton).Label(fmt.Sprintf(locale.M().CloseNamed, tab.Title)).IgnoreInput(!offered)
		slot = style.Background(slot, theme.RadiusScale().SM, ui.Transparent, t.Accent, t.Accent)
		closed = slot.Clicked()
	} else {
		slot.Role(ui.RoleNone)
	}
	cross := closable && offered && (!tab.Dirty || slot.Hovered())
	// The dot of unsaved changes fades in and out with the state; the
	// swap with the cross follows the pointer and happens at once.
	dirty := motion.Progress(slot, "dirty", tab.Dirty)
	slot.Children(func() {
		switch {
		case cross:
			glyph := scale(c, theme.IconSM)
			ui.Icon(c, lucide.X()).Size(glyph, glyph).Shrink(0)
		case dirty > 0:
			dot := scale(c, theme.Space8)
			ui.Box(c).Size(dot, dot).Radius(dot / 2).Background(ink.Alpha(dirty))
		}
	})
	return closed
}

// docTabClamp keeps the held tab between the first and the last tab of its
// group, lo and hi.
func docTabClamp(bounds []ui.Rect, held int, dx float32, lo, hi int) float32 {
	a, first, last := bounds[held], bounds[lo], bounds[hi]
	left := min(first.X, last.X) - a.X
	right := max(first.X+first.W, last.X+last.W) - a.X - a.W
	return min(max(dx, left), right)
}

// docTabOver returns the place whose tab center is nearest the held tab's.
func docTabOver(bounds []ui.Rect, held int, dx float32, lo, hi int) int {
	center := bounds[held].X + bounds[held].W/2 + dx
	best, far := held, float32(-1)
	for i := lo; i <= hi; i++ {
		b := bounds[i]
		d := b.X + b.W/2 - center
		if d < 0 {
			d = -d
		}
		if far < 0 || d < far {
			best, far = i, d
		}
	}
	return best
}

// docTabShiftOf is how far tab i moves toward the start for the held tab
// going to over: negative makes room after it.
func docTabShiftOf(bounds []ui.Rect, held, over, i int, gap float32) float32 {
	if held < 0 {
		return 0
	}
	room := bounds[held].W + gap
	switch {
	case held < i && i <= over:
		return -room
	case over <= i && i < held:
		return room
	}
	return 0
}
