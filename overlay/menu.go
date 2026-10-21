package overlay

import (
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
	"strings"
	"time"
)

type MenuItemKind uint8

const (
	MenuAction MenuItemKind = iota
	MenuCheckbox
	MenuRadio
	MenuSeparator
	MenuLabel
)

// MenuActionID is a stable caller identity, separate from visible labels.
type MenuActionID string

// MenuChoice distinguishes an exclusive choice value from an action identity.
type MenuChoice string

type MenuItem struct {
	ID                    MenuActionID
	Label                 string
	Kind                  MenuItemKind
	Disabled, Destructive bool
	Checked               *bool
	Radio                 *MenuChoice
	Value                 MenuChoice
	Icon                  *myui.SVG
	Shortcut              string
	OnSelect              func()
	Items                 []MenuItem
}

type MenuOptions struct {
	Label    string
	Disabled bool
	Handle   *myui.Handle
}

// MenuState is caller-owned presentation state. IDs must be unique within each
// menu; labels are used as identity when ID is absent. Do not copy active state.
type MenuState struct {
	Open    bool
	root    menuLevel
	trigger myui.Handle
}

type menuLevel struct {
	active  MenuActionID
	focus   bool
	prefix  string
	typed   time.Time
	child   *menuLevel
	childID MenuActionID
	closing bool
}

func menuID(it MenuItem) MenuActionID {
	if it.ID != "" {
		return it.ID
	}
	return MenuActionID(it.Label)
}
func selectable(it MenuItem) bool {
	return !it.Disabled && it.Kind != MenuSeparator && it.Kind != MenuLabel
}
func firstItem(items []MenuItem, backwards bool) int {
	for j := range items {
		i := j
		if backwards {
			i = len(items) - 1 - j
		}
		if selectable(items[i]) {
			return i
		}
	}
	return -1
}
func activeItem(items []MenuItem, id MenuActionID) int {
	for i, it := range items {
		if menuID(it) == id && selectable(it) {
			return i
		}
	}
	return -1
}
func moveItem(items []MenuItem, i, delta int) int {
	if len(items) == 0 {
		return -1
	}
	for range len(items) {
		i = (i + delta + len(items)) % len(items)
		if selectable(items[i]) {
			return i
		}
	}
	return -1
}
func (s *MenuState) show(items []MenuItem, backwards bool) {
	s.Open = true
	s.root = menuLevel{focus: true}
	if i := firstItem(items, backwards); i >= 0 {
		s.root.active = menuID(items[i])
	}
}

func DropdownMenu(c *myui.Context, s *MenuState, o MenuOptions, items []MenuItem) myui.Element {
	if s == nil {
		return myui.Element{}
	}
	h := o.Handle
	if h == nil {
		h = &s.trigger
	}
	trigger := mygoui.Button(c, mygoui.ButtonOptions{Label: o.Label, Disabled: o.Disabled, Variant: mygoui.ButtonOutline, Handle: h,
		OnClick: func() {
			if s.Open {
				s.Open = false
			} else {
				s.show(items, false)
			}
		}})
	trigger.Role(myui.RoleMenuButton).Expanded(s.Open).HandleInput(func(ev myui.InputEvent) bool {
		if ev.Kind == myui.InputKeyDown && ev.Mods == 0 && (ev.Key == myui.KeyDown || ev.Key == myui.KeyUp) && !o.Disabled {
			s.show(items, ev.Key == myui.KeyUp)
			return true
		}
		return false
	})
	if s.Open && !o.Disabled {
		menuPanel(c, trigger, s, &s.root, items, false, nil)
	}
	return trigger
}

func menuPanel(c *myui.Context, anchor myui.Element, s *MenuState, level *menuLevel, items []MenuItem, nested bool, switchMenu func(int)) {
	if level.child != nil && level.child.closing {
		level.child = nil
		level.childID = ""
		level.focus = true
	}
	if activeItem(items, level.active) < 0 {
		if i := firstItem(items, false); i >= 0 {
			level.active = menuID(items[i])
			level.focus = true
		}
	}
	myui.Overlay(c, func() {
		p := theme.Colors(c)
		w, h := c.Size()
		panel := myui.Column(c).Role(myui.RoleMenu).Label(locales.ComponentMenu()).MinWidth(144).MaxWidth(max(0, w-16)).MaxHeight(max(0, h-16)).
			Padding(4).Gap(2).Radius(16).Background(p.Popover).TextColor(p.PopoverForeground).Border(1, p.Border).OnClick(func() {})
		placement := Placement{Side: SideBottom, Align: AlignStart}
		if nested {
			placement.Side = SideRight
		}
		place(panel, anchor, placement, SideBottom)
		var childAnchor myui.Element
		var childItems []MenuItem
		panel.Children(func() {
			myui.Scroll(c).MinHeight(0).Children(func() {
				for i, it := range items {
					id := menuID(it)
					if it.Kind == MenuSeparator {
						myui.Box(c.Key(id)).Height(1).Background(p.Border).MarginY(4)
						continue
					}
					if it.Kind == MenuLabel {
						myui.Text(c.Key(id), it.Label).Padding(4, 8).FontSize(theme.Scale(c, 12)).TextColor(p.MutedForeground)
						continue
					}
					b := myui.ButtonBase(c.Key(id)).Label(it.Label).Disabled(it.Disabled).Role(myui.RoleMenuItem).MinHeight(theme.Scale(c, 28)).
						Padding(6, 8).Gap(8).Radius(12).Justify(myui.Start).TextColor(p.PopoverForeground)
					if it.Kind == MenuCheckbox {
						b.Role(myui.RoleMenuItemCheckBox).Checked(it.Checked != nil && *it.Checked)
					}
					if it.Kind == MenuRadio {
						b.Role(myui.RoleMenuItemRadio).Checked(it.Radio != nil && *it.Radio == it.Value)
					}
					if it.Destructive {
						b.TextColor(p.Destructive)
					}
					if it.Disabled {
						b.TextColor(p.DisabledForeground)
					}
					// Pointer hover paints independently of keyboard selection.
					if b.Focused() || b.Hovered() && !it.Disabled {
						b.Background(p.Accent).TextColor(p.AccentForeground)
					}
					if level.focus && level.active == id && !it.Disabled {
						b.Focus()
						level.focus = false
					}
					b.HandleInput(func(ev myui.InputEvent) bool { return menuInput(ev, s, level, items, i, nested, switchMenu) })
					b.OnClick(func() {
						if len(it.Items) > 0 {
							openChild(level, it)
						} else {
							activate(it)
							s.Open = false
						}
					}).Children(func() {
						if it.Icon != nil {
							myui.Icon(c, it.Icon).Size(16, 16)
						}
						myui.Text(c, it.Label).Grow(1)
						if it.Shortcut != "" {
							myui.Text(c, it.Shortcut).FontSize(theme.Scale(c, 12)).TextColor(p.MutedForeground)
						}
						if len(it.Items) > 0 {
							myui.Text(c, "›")
						} else if it.Kind == MenuCheckbox && it.Checked != nil && *it.Checked || it.Kind == MenuRadio && it.Radio != nil && *it.Radio == it.Value {
							myui.Text(c, "✓")
						}
					})
					if len(it.Items) > 0 {
						b.Expanded(level.child != nil && level.childID == id)
					}
					if level.child != nil && level.childID == id && !it.Disabled {
						childAnchor = b
						childItems = it.Items
					}
				}
			})
		})
		if childAnchor.Valid() {
			menuPanel(c, childAnchor, s, level.child, childItems, true, switchMenu)
		}
		if level.child != nil && level.child.closing && childAnchor.Valid() {
			childAnchor.Focus()
		}
		if panel.OverlayShortcut(0, myui.KeyEscape) {
			if nested {
				level.closing = true
			} else {
				s.Open = false
			}
		}
		if panel.PressedOutside() {
			s.Open = false
		}
	})
	if level.child != nil && level.child.closing {
		level.child = nil
		level.childID = ""
		level.focus = true
	}
}

func openChild(level *menuLevel, it MenuItem) {
	level.active = menuID(it)
	child := &menuLevel{focus: true}
	if i := firstItem(it.Items, false); i >= 0 {
		child.active = menuID(it.Items[i])
	}
	level.child, level.childID = child, menuID(it)
}
func activate(it MenuItem) {
	if !selectable(it) {
		return
	}
	if it.Kind == MenuCheckbox && it.Checked != nil {
		*it.Checked = !*it.Checked
	}
	if it.Kind == MenuRadio && it.Radio != nil {
		*it.Radio = it.Value
	}
	if it.OnSelect != nil {
		it.OnSelect()
	}
}

func menuInput(ev myui.InputEvent, s *MenuState, level *menuLevel, items []MenuItem, index int, nested bool, switchMenu func(int)) bool {
	if ev.Kind != myui.InputKeyDown || ev.Mods != 0 {
		return false
	}
	var next = -1
	switch ev.Key {
	case myui.KeyDown:
		next = moveItem(items, index, 1)
	case myui.KeyUp:
		next = moveItem(items, index, -1)
	case myui.KeyHome:
		next = firstItem(items, false)
	case myui.KeyEnd:
		next = firstItem(items, true)
	case myui.KeyRight:
		if len(items[index].Items) > 0 {
			openChild(level, items[index])
			return true
		}
		if switchMenu != nil {
			switchMenu(1)
			return true
		}
		return false
	case myui.KeyLeft:
		if nested {
			level.closing = true
			return true
		}
		if switchMenu != nil {
			switchMenu(-1)
			return true
		}
		return false
	case myui.KeyTab:
		s.Open = false
		return false
	default:
		var r rune
		if ev.Key >= myui.KeyA && ev.Key <= myui.KeyZ {
			r = 'a' + rune(ev.Key-myui.KeyA)
		} else if ev.Key >= myui.Key0 && ev.Key <= myui.Key9 {
			r = '0' + rune(ev.Key-myui.Key0)
		} else {
			return false
		}
		if ev.Time.Sub(level.typed) > time.Second {
			level.prefix = ""
		}
		level.typed = ev.Time
		letter := string(r)
		if level.prefix != letter {
			level.prefix += letter
		}
		for step := 1; step <= len(items); step++ {
			i := (index + step) % len(items)
			if selectable(items[i]) && strings.HasPrefix(strings.ToLower(items[i].Label), level.prefix) {
				next = i
				break
			}
		}
	}
	if next >= 0 {
		level.active = menuID(items[next])
		level.focus = true
		level.child = nil
		level.childID = ""
	}
	return true
}

// ContextMenu deliberately uses the system's menu, including Shift+F10, arrows,
// typeahead and nested navigation. Its system appearance cannot use theme tokens;
// radio choices use the platform check indicator rather than a custom dot.
func ContextMenu(anchor myui.Element, items []MenuItem) myui.Element {
	return anchor.ContextMenu(func(m *myui.Menu) { nativeItems(m, items) })
}
func nativeItems(m *myui.Menu, items []MenuItem) {
	for _, it := range items {
		if it.Kind == MenuSeparator {
			m.Separator()
			continue
		}
		if len(it.Items) > 0 {
			m.Submenu(it.Label, func(sub *myui.Menu) { nativeItems(sub, it.Items) }).Disabled(it.Disabled)
			continue
		}
		item := m.Item(it.Label).Disabled(it.Disabled || it.Kind == MenuLabel)
		if it.Kind == MenuCheckbox {
			item.Checked(it.Checked != nil && *it.Checked)
		}
		if it.Kind == MenuRadio {
			item.Checked(it.Radio != nil && *it.Radio == it.Value)
		}
		if item.Chosen() {
			activate(it)
		}
	}
}

type BarMenu struct {
	ID       MenuActionID
	Label    string
	Disabled bool
	Items    []MenuItem
}
type MenubarState struct {
	Menu   MenuState
	active int
	focus  bool
}

func Menubar(c *myui.Context, s *MenubarState, menus []BarMenu) myui.Element {
	if s == nil {
		return myui.Element{}
	}
	switchTo := func(delta int) {
		for range len(menus) {
			s.active = (s.active + delta + len(menus)) % len(menus)
			if !menus[s.active].Disabled {
				s.focus = true
				if s.Menu.Open {
					s.Menu.show(menus[s.active].Items, false)
				}
				return
			}
		}
	}
	p := theme.Colors(c)
	bar := myui.Row(c).Role(myui.RoleMenuBar).FocusGroup(myui.Horizontal).Gap(2).Padding(3).Radius(16).Border(1, p.Border)
	bar.Children(func() {
		for i, m := range menus {
			id := m.ID
			if id == "" {
				id = MenuActionID(m.Label)
			}
			b := mygoui.Button(c.Key(id), mygoui.ButtonOptions{Label: m.Label, Disabled: m.Disabled, Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM,
				OnClick: func() {
					if s.Menu.Open && s.active == i {
						s.Menu.Open = false
					} else {
						s.active = i
						s.Menu.show(m.Items, false)
					}
				}})
			b.Role(myui.RoleMenuItem).Expanded(s.Menu.Open && s.active == i)
			if s.focus && s.active == i {
				b.Focus()
				s.focus = false
			}
			if s.Menu.Open && b.Hovered() && s.active != i && !m.Disabled {
				s.active = i
				s.Menu.show(m.Items, false)
			}
			b.HandleInput(func(ev myui.InputEvent) bool {
				if ev.Kind != myui.InputKeyDown || ev.Mods != 0 {
					return false
				}
				s.active = i
				switch ev.Key {
				case myui.KeyRight:
					switchTo(1)
				case myui.KeyLeft:
					switchTo(-1)
				case myui.KeyHome:
					s.active = len(menus) - 1
					switchTo(1)
				case myui.KeyEnd:
					s.active = 0
					switchTo(-1)
				case myui.KeyDown, myui.KeyUp:
					s.Menu.show(m.Items, ev.Key == myui.KeyUp)
				default:
					return false
				}
				return true
			})
			if s.Menu.Open && s.active == i && !m.Disabled {
				menuPanel(c, b, &s.Menu, &s.Menu.root, m.Items, false, switchTo)
			}
		}
	})
	return bar
}
