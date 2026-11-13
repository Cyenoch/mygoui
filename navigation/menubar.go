package navigation

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// MenubarMenuOptions is one top-level menu of a Menubar. Key identifies it
// in MenubarOptions.Open; it defaults to Label. Mnemonic is the letter of
// Alt+letter: it defaults to the first letter of Label that no menu before
// uses, and is underlined while Alt is held.
type MenubarMenuOptions struct {
	Key      string
	Label    string
	Mnemonic rune
	Disabled bool
	Menu     MenuOptions
}

func (m *MenubarMenuOptions) key() string {
	if m.Key != "" {
		return m.Key
	}
	return m.Label
}

// MenubarOptions configures a menu bar. Open is caller-owned and holds the
// key of the menu showing, empty for none; nil keeps it internally.
type MenubarOptions struct {
	Label string
	Open  *string
	Menus []MenubarMenuOptions
}

// mnemonicsOn is whether menus open by Alt+letter and Alt or F10 moves the
// focus to the bar: the convention of Windows and Linux, which macOS lacks.
var mnemonicsOn = ui.Cmd != ui.Super

// altTapLimit is how long Alt may be held and still count as a tap.
const altTapLimit = 600 * time.Millisecond

// menubarState is what the bar remembers between builds: a request to take
// the focus, whether the keyboard is driving the bar, and the Alt key's
// press, which asks for the focus when it ends without another key.
type menubarState struct {
	focus    bool
	keyboard bool
	altDown  bool
	altUsed  bool
	altAt    time.Time
}

// Menubar builds a horizontal application menu bar. The bar is one stop of
// Tab whose Left and Right move among the menu titles; a click, Enter, Space
// or Down opens a title's menu. While a menu shows, Left and Right switch to
// the neighbouring menu, moving the pointer over another title switches to
// it, and Escape closes the menu and gives the focus back to its title.
//
// The accelerators of the items work in the window while the bar is built,
// whether or not their menu shows. On Windows and Linux, F10 or a tap of Alt
// moves the focus to the bar and Alt with a menu's mnemonic opens it, its
// letter underlined while Alt is held; macOS has neither.
func Menubar(c *ui.Context, o MenubarOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Row(c).Shrink(0).MinWidth(0).AlignSelf(ui.Start).AlignItems(ui.Center).
		Gap(scale(c, theme.Space(1))).Padding(scale(c, theme.Space(1))).
		Radius(theme.RadiusScale().LG).Border(theme.Hairline, t.Border).Background(t.Background).
		Role(ui.RoleMenuBar).FocusGroup(ui.Horizontal)
	if o.Label != "" {
		root.Label(o.Label)
	}
	open := o.Open
	if open == nil {
		open = ui.Local(root, "menubar-open", func() string { return "" })
	}
	if *open != "" && !menubarOpens(o.Menus, *open) {
		*open = "" // its menu went or was disabled while showing
	}
	menubarAccelerators(c, o.Menus)
	st := ui.Local(root, "menubar-state", func() menubarState { return menubarState{} })
	var letters []rune
	if mnemonicsOn {
		letters = menubarMnemonics(o.Menus)
		menubarKeys(c, root, o, st, open, letters)
	}
	focus := st.focus
	st.focus = false
	root.Children(func() {
		for i := range o.Menus {
			buildMenubarMenu(c, o, i, open, menubarTitle{
				mnemonic: pickLetter(letters, i), underline: st.keyboard && mnemonicsOn, focus: focus,
			})
			if o.Menus[i].Disabled || !focus {
				continue
			}
			focus = false
		}
	})
	return root
}

// menubarTitle is how one title is built: its mnemonic letter, whether to
// underline it, and whether it takes the focus the keyboard asked the bar for.
type menubarTitle struct {
	mnemonic  rune
	underline bool
	focus     bool
}

func pickLetter(letters []rune, i int) rune {
	if i < len(letters) {
		return letters[i]
	}
	return 0
}

// menubarKeys binds F10, Alt taps and Alt+letter, and follows how long Alt
// was held.
func menubarKeys(c *ui.Context, root ui.Element, o MenubarOptions, st *menubarState, open *string, letters []rune) {
	ask := func() { st.focus, st.keyboard = true, true }
	c.OnWindowShortcut(0, ui.KeyF10, ask)
	for i, letter := range letters {
		if letter == 0 || o.Menus[i].Disabled {
			continue
		}
		key := o.Menus[i].key()
		c.OnWindowShortcut(ui.Alt, letterKey(letter), func() {
			st.altUsed, st.keyboard = true, true
			*open = key
		})
	}
	mods := c.Modifiers()
	alt := mods&ui.Alt != 0
	switch {
	case alt && !st.altDown:
		st.altDown, st.altUsed, st.altAt = true, mods != ui.Alt, c.Now()
	case alt && (len(c.KeysDown()) > 0 || mods != ui.Alt):
		// A key or another modifier with Alt, as Alt+Shift switches the
		// layout, makes it no tap.
		st.altUsed = true
	case !alt && st.altDown:
		st.altDown = false
		if !st.altUsed && *open == "" && c.Now().Sub(st.altAt) < altTapLimit {
			ask()
			c.Invalidate()
		}
	}
	if alt {
		st.keyboard = true
	} else if st.keyboard && *open == "" && !root.FocusWithin() && !st.focus {
		st.keyboard = false
	}
}

// menubarMnemonics returns the Alt+letter of each menu: its own Mnemonic, or
// the first letter of its label no earlier menu has, 0 for none.
func menubarMnemonics(menus []MenubarMenuOptions) []rune {
	letters := make([]rune, len(menus))
	taken := map[rune]bool{}
	for i := range menus {
		if r := unicode.ToLower(menus[i].Mnemonic); r != 0 {
			letters[i], taken[r] = r, true
		}
	}
	for i := range menus {
		if letters[i] != 0 {
			continue
		}
		for _, r := range menus[i].Label {
			if r = unicode.ToLower(r); r >= 'a' && r <= 'z' && !taken[r] {
				letters[i], taken[r] = r, true
				break
			}
		}
	}
	return letters
}

// letterKey is the key typing a lowercase ASCII letter, 0 for others.
func letterKey(r rune) ui.Key {
	if r < 'a' || r > 'z' {
		return ui.KeyUnknown
	}
	return ui.KeyA + ui.Key(r-'a')
}

// menubarAccelerators registers the accelerator of every enabled item of the
// enabled menus for the window.
func menubarAccelerators(c *ui.Context, menus []MenubarMenuOptions) {
	for i := range menus {
		if menus[i].Disabled {
			continue
		}
		m := menus[i].Menu
		sc := newMenuScope(m, menuExtra{}, nil)
		for _, section := range m.Sections {
			registerAccelerators(c, section.Items, sc)
		}
	}
}

func registerAccelerators(c *ui.Context, items []MenuItemOptions, sc *menuScope) {
	for _, item := range items {
		switch {
		case item.Separator || item.Disabled:
		case len(item.Submenu) > 0:
			registerAccelerators(c, item.Submenu, sc)
		case item.Accelerator.Key != ui.KeyUnknown:
			key := item.Key
			if key == "" {
				key = item.Label
			}
			c.OnWindowShortcut(item.Accelerator.Mods, item.Accelerator.Key, func() {
				activateMenuItem(sc, true, key, item)
			})
		}
	}
}

// menubarTrigger is one title's state between frames: whether its panel
// shows, whether it did in the last build, and where the pointer was over it.
type menubarTrigger struct {
	shown, was bool
	over       bool
	x, y       float32
}

func buildMenubarMenu(c *ui.Context, o MenubarOptions, i int, open *string, title menubarTitle) {
	m := &o.Menus[i]
	key := m.key()
	t := theme.Colors(c)
	radius := theme.RadiusScale().SM
	c.Key(key)
	el := ui.ButtonBase(c).Shrink(0).Center().
		Padding(scale(c, theme.Space(1)), scale(c, theme.Space(2))).MinHeight(scale(c, 28)).
		Radius(radius).FocusRing(true).Role(ui.RoleMenuItem).Label(m.Label).
		TextColor(t.Foreground).Cursor(ui.CursorPointer)
	st := ui.Local(el, "menubar-trigger", func() menubarTrigger { return menubarTrigger{} })
	if st.was && !st.shown && *open == key {
		*open = "" // an item was chosen
		el.Focus()
	}
	if !m.Disabled && menubarPointerMoved(el, st) && *open != "" && *open != key {
		*open = key
		c.Invalidate() // the menu left closes in the next build
	}
	shown := *open == key
	st.shown, st.was = shown, shown
	el.Expanded(shown)
	el = style.Disabled(el, m.Disabled)
	rest := ui.Transparent
	if shown {
		rest = t.Accent
	}
	el = style.Background(el, radius, rest, t.Accent, t.Accent)
	if title.focus && !m.Disabled {
		el.Focus()
	}
	el.Children(func() { menubarTitleText(c, m.Label, title) })
	toggle := func() {
		if *open == key {
			*open = ""
		} else {
			*open = key
		}
	}
	el.OnClick(toggle)
	if !shown {
		el.OnShortcut(0, ui.KeyDown, toggle)
	}
	panel := menuPopover(c, el, &st.shown, m.Menu, mygoui.Placement{Align: mygoui.AlignStart}, false)
	if !panel.Valid() {
		return
	}
	if panel.Dismissed() {
		*open = ""
		// MyGo gives the focus back to what held it as the first panel
		// opened; a menu reached from another one returns it to its own
		// title. A press outside takes the focus where it lands.
		if !panel.PressedOutside() {
			el.Focus()
		}
		return
	}
	if panel.FocusLeft() {
		*open = "" // Tab moved on from the menu
		return
	}
	forth, back := ui.KeyRight, ui.KeyLeft
	if c.Direction() == ui.RightToLeft {
		forth, back = back, forth
	}
	panel.OnShortcut(0, forth, func() { *open = menubarNeighbour(o.Menus, i, 1) })
	panel.OnShortcut(0, back, func() { *open = menubarNeighbour(o.Menus, i, -1) })
}

// menubarTitleText draws a title, with its mnemonic letter underlined while
// the keyboard is driving the bar.
func menubarTitleText(c *ui.Context, label string, title menubarTitle) {
	text := func(s string) ui.Element {
		return style.Text(ui.Text(c, s), c, theme.SM()).FontWeight(theme.WeightMedium).SingleLine()
	}
	at := strings.IndexFunc(label, func(r rune) bool { return unicode.ToLower(r) == title.mnemonic })
	if title.mnemonic == 0 || at < 0 {
		text(label)
		return
	}
	_, size := utf8.DecodeRuneInString(label[at:])
	ui.Row(c).Shrink(0).Children(func() {
		if at > 0 {
			text(label[:at])
		}
		mark := text(label[at : at+size])
		if title.underline {
			mark.Underline()
		}
		if at+size < len(label) {
			text(label[at+size:])
		}
	})
}

// menubarPointerMoved reports whether the pointer came over el or moved on
// it since the last build. A pointer resting on a title does not take a menu
// opened from the keyboard.
func menubarPointerMoved(el ui.Element, st *menubarTrigger) bool {
	x, y, over := el.PointerPosition()
	moved := over && (!st.over || x != st.x || y != st.y)
	st.over, st.x, st.y = over, x, y
	return moved
}

// menubarOpens reports whether key names an enabled menu.
func menubarOpens(menus []MenubarMenuOptions, key string) bool {
	for i := range menus {
		if menus[i].key() == key {
			return !menus[i].Disabled
		}
	}
	return false
}

// menubarNeighbour returns the key of the next enabled menu from i in
// direction step, going round the ends.
func menubarNeighbour(menus []MenubarMenuOptions, i, step int) string {
	n := len(menus)
	for d := 1; d <= n; d++ {
		j := ((i+step*d)%n + n) % n
		if !menus[j].Disabled {
			return menus[j].key()
		}
	}
	return menus[i].key()
}

// menuPanelKey keys the panel of the menu that open shows.
type menuPanelKey struct{ open *bool }

// menuPopover builds menu in a panel anchored to trigger while *open is
// true, focusing its first item as it opens, and returns the panel, invalid
// while closed. Choosing an item, Escape or a press outside sets *open to
// false. A menu button trigger also opens on Down at the first item and on
// Up at the last, and choosing or Escape gives it the focus back: MyGo
// would return it to the first item, which goes with the panel. A press
// outside leaves the focus where it lands. Dropdown's in-window path is
// the same panel.
func menuPopover(c *ui.Context, trigger ui.Element, open *bool, menu MenuOptions, placement mygoui.Placement, menuButton bool) ui.Element {
	prev := ui.Local(trigger, "menu-was-open", func() bool { return false })
	fromEnd := ui.Local(trigger, "menu-from-end", func() bool { return false })
	enter := *open && !*prev
	closed := *prev && !*open
	ml := ui.Local(trigger, "menu-list", newMenuList)
	if enter {
		*ml = newMenuList() // each opening starts at the top
	}
	first, last := enter && !*fromEnd, enter && *fromEnd
	if enter {
		*fromEnd = false
	}
	if menuButton && !*open {
		trigger.OnShortcut(0, ui.KeyDown, func() { *open = true })
		trigger.OnShortcut(0, ui.KeyUp, func() { *open, *fromEnd = true, true })
	}
	radius := theme.RadiusScale().LG
	panel := ui.PopoverBase(c, trigger, open, func(p ui.Element) {
		// An unkeyed panel takes its identity from the overlays built
		// before it: a trigger's tooltip closing as the menu opens would
		// rebuild every item and drop the focus given to the first.
		p = p.Key(menuPanelKey{open})
		p.MinWidth(scale(c, 128)).MaxWidth(scale(c, 320)).Padding(0)
		anchored.Surface(c, p.Clip(), radius)
		buildMenuPanel(c, menu, menuExtra{close: open, focusFirst: &first, focusLast: &last}, nil, *ml, scale(c, 320)).
			MaxHeight(scale(c, 420)).MinWidth(0).Shrink(0)
	})
	if panel.Valid() {
		anchored.Apply(panel, trigger, placement, mygoui.SideBottom)
		if panel.Dismissed() {
			closed = !panel.PressedOutside()
		}
	}
	*prev = *open
	if closed && menuButton {
		trigger.Focus()
	}
	return panel
}
