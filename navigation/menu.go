package navigation

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// MenuSelection is how a menu tracks checked items.
type MenuSelection uint8

const (
	MenuSelectNone MenuSelection = iota
	MenuSelectSingle
	MenuSelectMultiple
)

// MenuItemVariant is a menu item's emphasis.
type MenuItemVariant uint8

const (
	MenuItemDefault MenuItemVariant = iota
	MenuItemDanger
)

// MenuItemOptions configures one menu row. Key identifies it in a selection;
// it defaults to Label. Submenu turns the row into a submenu trigger opening
// the given items to the side. Separator renders a rule and ignores the rest.
type MenuItemOptions struct {
	Key         string
	Label       string
	Description string
	Shortcut    string
	Icon        *ui.SVG
	Variant     MenuItemVariant
	Disabled    bool
	Separator   bool
	Checked     bool
	Submenu     []MenuItemOptions
	OnSelect    func()
}

// MenuSectionOptions configures a group of menu items with an optional
// heading. Key keeps the group stable when headings are absent.
type MenuSectionOptions struct {
	Key   string
	Title string
	Items []MenuItemOptions
}

// MenuOptions configures menu contents. Selected is caller-owned and holds
// the keys of checked items; nil keeps it internally. OnSelect receives the
// key of an activated item, in addition to its own callback.
type MenuOptions struct {
	Label     string
	Selection MenuSelection
	Selected  *[]string
	Sections  []MenuSectionOptions
	OnSelect  func(key string)
}

// menuScope is the selection and closing state a menu shares with its items.
type menuScope struct {
	selection  MenuSelection
	selected   *[]string
	onSelect   func(string)
	close      *bool
	focusFirst *bool
}

// menuExtra carries wiring that only a root (Dropdown or SubMenu) has.
type menuExtra struct {
	close      *bool
	focusFirst *bool
}

// Menu builds a menu: a vertical stack of sections and items with Arrow, Home
// and End moving focus, Enter and Space activating, and checked semantics for
// single or multiple selection. It is the panel content of Dropdown, and can
// stand alone when a caller anchors it.
func Menu(c *ui.Context, o MenuOptions) ui.Element {
	return buildMenu(c, o, menuExtra{}, nil)
}

func buildMenu(c *ui.Context, o MenuOptions, extra menuExtra, parent *menuScope) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Gap(scale(c, theme.Space4)).
		Padding(scale(c, theme.Space4)).Role(ui.RoleMenu).FocusGroup(ui.Vertical).Clip()
	if o.Label != "" {
		root.Label(o.Label)
	}
	selection := o.Selection
	selected := o.Selected
	onSelect := o.OnSelect
	closeTo := extra.close
	focusFirst := extra.focusFirst
	if parent != nil {
		if selection == MenuSelectNone {
			selection = parent.selection
		}
		if selected == nil {
			selected = parent.selected
		}
		if onSelect == nil {
			onSelect = parent.onSelect
		}
		if closeTo == nil {
			closeTo = parent.close
		}
		if focusFirst == nil {
			focusFirst = parent.focusFirst
		}
	}
	if selection != MenuSelectNone && selected == nil {
		selected = ui.Local(root, "menu-selected", func() []string { return nil })
	}
	sc := &menuScope{
		selection:  selection,
		selected:   selected,
		onSelect:   onSelect,
		close:      closeTo,
		focusFirst: focusFirst,
	}
	root.Children(func() {
		for i := range o.Sections {
			buildMenuSection(c, o.Sections[i], sc)
		}
	})
	return root
}

// buildMenuSection builds a titled group of menu items.
func buildMenuSection(c *ui.Context, o MenuSectionOptions, sc *menuScope) ui.Element {
	t := theme.Colors(c)
	if o.Key != "" {
		c.Key(o.Key)
	}
	section := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Gap(0).Role(ui.RoleGroup)
	if o.Title != "" {
		section.Label(o.Title)
	}
	section.Children(func() {
		if o.Title != "" {
			heading := ui.Row(c).FillWidth().AlignItems(ui.Center).Shrink(0).
				Padding(scale(c, theme.Space6), scale(c, theme.Space8), scale(c, theme.Space4), scale(c, theme.Space8)).
				Role(ui.RoleText)
			heading.Children(func() {
				style.Text(ui.Text(c, o.Title), c, theme.XS()).TextColor(t.Muted).
					FontWeight(theme.WeightMedium).SingleLine()
			})
		}
		for i := range o.Items {
			buildMenuItem(c, o.Items[i], sc)
		}
	})
	return section
}

// buildMenuItem builds one focusable row. The menu scope is explicit, so no
// ambient stack or free compound constructor is needed.
func buildMenuItem(c *ui.Context, o MenuItemOptions, sc *menuScope) ui.Element {
	t := theme.Colors(c)
	if o.Separator {
		// Upstream insets the menu's separator to 3% on each side; a plain
		// full-width box paints that rule in its own paint.
		sep := ui.Box(c).FillWidth().Shrink(0).Height(hairline).Role(ui.RoleNone)
		sep.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + r.W*0.03, Y: r.Y, W: r.W * 0.94, H: r.H}, t.Separator, 2)
		})
		return sep
	}
	inMenu := true
	key := o.Key
	if key == "" {
		key = o.Label
	}
	role := ui.RoleMenuItem
	switch {
	case inMenu && sc.selection == MenuSelectSingle:
		role = ui.RoleMenuItemRadio
	case inMenu && sc.selection == MenuSelectMultiple:
		role = ui.RoleMenuItemCheckBox
	case o.Checked:
		role = ui.RoleMenuItemCheckBox
	}
	on := o.Checked
	if inMenu && sc.selection != MenuSelectNone && sc.selected != nil {
		on = containsKey(*sc.selected, key)
	}
	hasIndicator := (inMenu && sc.selection != MenuSelectNone) || o.Checked
	hasSubmenu := len(o.Submenu) > 0
	radius := theme.RadiusScale().SM
	if o.Key != "" {
		c.Key(o.Key)
	}
	el := ui.ButtonBase(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Justify(ui.Start).
		Gap(scale(c, theme.Space8)).Padding(scale(c, theme.Space4), scale(c, theme.Space8)).
		MinHeight(scale(c, 28)).Radius(radius).FocusRing(true).
		Role(role).Label(o.Label).Cursor(ui.CursorPointer)
	el.Checked(on)
	el = style.Background(el, radius, ui.Transparent, t.Accent, t.Accent)
	// Upstream's indicator is absolute: the row pads its start by 28 DIP to
	// clear it, and both sides where a submenu keeps it at the end.
	if hasIndicator {
		el.PaddingStart(scale(c, 28))
		if hasSubmenu {
			el.PaddingStart(scale(c, theme.Space8))
			el.PaddingEnd(scale(c, 28))
		}
	}
	fg := t.Foreground
	indicatorColor := t.Muted
	if o.Variant == MenuItemDanger {
		fg, indicatorColor = t.DestructiveForeground, t.DestructiveForeground
	}
	if o.Description != "" {
		el.Description(o.Description)
	}
	el.TextColor(fg)
	var subOpen *bool
	if hasSubmenu {
		subOpen = ui.Local(el, "submenu", func() bool { return false })
		el.Expanded(*subOpen)
		el.OnClick(func() { *subOpen = !*subOpen })
		el.OnShortcut(0, ui.KeyRight, func() { *subOpen = true })
	} else {
		el.OnClick(func() { activateMenuItem(sc, inMenu, key, o) })
	}
	// Checkbox and radio menu items share the same Coss check indicator.
	if hasIndicator {
		target := float32(0)
		if on {
			target = 1
		}
		shown := target
		if !theme.ReduceMotion(c) {
			shown = el.AnimateTo("menu-indicator", target, ui.Animation{Duration: menuCheckMotion, Ease: ui.EaseOut})
		}
		glyph := &ui.Path{}
		iconSize := scale(c, theme.IconMD)
		el.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if shown <= 0.001 {
				return
			}
			start := r.X + scale(c, theme.Space8)
			end := r.X + r.W - scale(c, theme.Space8) - iconSize
			if c.Direction() == ui.RightToLeft {
				start, end = end, start
			}
			x := start
			if hasSubmenu {
				x = end
			}
			box := ui.Rect{X: x, Y: r.Y + (r.H-iconSize)/2, W: iconSize, H: iconSize}
			drawCheck(p, glyph, box, shown, indicatorColor)
		})
	}
	el.Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
		}
		text := ui.Column(c).MinWidth(0).Grow(1).Gap(0)
		text.Children(func() {
			style.Text(ui.Text(c, o.Label), c, theme.SM()).TextColor(fg).SingleLine()
			if o.Description != "" {
				style.Text(ui.Text(c, o.Description), c, theme.XS()).TextColor(t.Muted).SingleLine()
			}
		})
		if o.Shortcut != "" {
			style.Text(ui.Text(c, o.Shortcut), c, theme.SM()).TextColor(t.Muted).
				FontWeight(theme.WeightMedium).Shrink(0)
		}
		if hasSubmenu {
			chev := ui.Box(c).Size(scale(c, theme.IconSM), scale(c, theme.IconSM)).Shrink(0)
			glyph := &ui.Path{}
			chev.Draw(func(p *ui.Painter, r ui.Rect) {
				angle := float32(0)
				if c.Direction() == ui.RightToLeft {
					angle = 180
				}
				drawChevron(p, glyph, r, angle, 1.5, t.Muted)
			})
		}
	})
	if inMenu && sc.focusFirst != nil && *sc.focusFirst && !o.Disabled {
		el.Focus()
		*sc.focusFirst = false
	}
	el = style.Disabled(el, o.Disabled)
	if hasSubmenu {
		SubMenu(c, el, SubMenuOptions{Open: subOpen, Items: o.Submenu, Label: o.Label, parent: sc})
	}
	return el
}

// activateMenuItem applies a menu item's selection, callbacks and close.
func activateMenuItem(sc *menuScope, inMenu bool, key string, o MenuItemOptions) {
	if inMenu && sc.selection != MenuSelectNone && sc.selected != nil {
		switch sc.selection {
		case MenuSelectSingle:
			*sc.selected = append((*sc.selected)[:0], key)
		case MenuSelectMultiple:
			if containsKey(*sc.selected, key) {
				removeKey(sc.selected, key)
			} else {
				*sc.selected = append(*sc.selected, key)
			}
		}
	}
	if o.OnSelect != nil {
		o.OnSelect()
	}
	if inMenu {
		if sc.onSelect != nil {
			sc.onSelect(key)
		}
		if sc.close != nil {
			*sc.close = false
		}
	}
}

// SubMenuOptions configures a submenu panel. Open is the caller-owned (or
// item-owned) expansion; nil renders nothing.
type SubMenuOptions struct {
	Open  *bool
	Items []MenuItemOptions
	Label string

	parent *menuScope
}

// SubMenu builds the panel of a submenu beside anchor while *Open is true.
// Escape or Left closes it; activating an item closes the enclosing dropdown.
func SubMenu(c *ui.Context, anchor ui.Element, o SubMenuOptions) ui.Element {
	if o.Open == nil {
		return ui.Element{}
	}
	keepOverlay(c)
	radius := theme.RadiusScale().LG
	var panel ui.Element
	panel = ui.PopoverBase(c, anchor, o.Open, func(p ui.Element) {
		p.MaxWidth(scale(c, 360)).MinWidth(scale(c, 128)).Padding(0)
		overlayPanel(c, p, radius)
		focus := ui.Local(p, "submenu-focus", func() bool { return true })
		scroll := ui.Scroll(c).MaxHeight(scale(c, 360)).MinWidth(0).Shrink(0)
		scroll.Children(func() {
			buildMenu(c, MenuOptions{Label: o.Label, Sections: []MenuSectionOptions{{Items: o.Items}}},
				menuExtra{focusFirst: focus}, o.parent)
		})
		if p.Shortcut(0, ui.KeyLeft) {
			*o.Open = false
		}
	})
	if panel.Valid() {
		// 4 DIP is the menu's own padding: the submenu's edge then meets the
		// parent panel's edge, as upstream's flush submenu does.
		panel = anchored.Apply(panel, anchor, mygoui.Placement{
			Side: mygoui.SideRight, Align: mygoui.AlignStart, Offset: scale(c, theme.Space4),
		}, mygoui.SideRight)
		panel = entryFade(c, panel, panelEntryMotion)
	}
	return panel
}

// DropdownOptions configures a drop-down menu. Open is caller-owned; nil keeps
// it internally. Label is both the trigger text and its accessible name.
type DropdownOptions struct {
	Open      *bool
	Label     string
	Icon      *ui.SVG
	Disabled  bool
	Placement mygoui.Placement
	Menu      MenuOptions
}

// Dropdown builds a menu button and its anchored menu panel. Click, Enter or
// Space opens it; a press outside or Escape closes it; choosing an item closes
// it and reports the item. Menu items receive the keyboard focus as it opens.
func Dropdown(c *ui.Context, o DropdownOptions) ui.Element {
	trigger := mygoui.Button(c, mygoui.ButtonOptions{
		Label:    o.Label,
		Icon:     o.Icon,
		EndIcon:  lucide.ChevronDown(),
		Variant:  mygoui.ButtonSecondary,
		Disabled: o.Disabled,
	})
	open := o.Open
	if open == nil {
		open = ui.Local(trigger, "dropdown-open", func() bool { return false })
	}
	keepOverlay(c)
	prev := ui.Local(trigger, "dropdown-was-open", func() bool { return *open })
	enter := *open && !*prev
	*prev = *open
	trigger.Role(ui.RoleMenuButton).Expanded(*open)
	if !o.Disabled {
		trigger.OnClick(func() { *open = !*open })
	}
	radius := theme.RadiusScale().LG
	panel := ui.PopoverBase(c, trigger, open, func(p ui.Element) {
		p.MinWidth(scale(c, 128)).MaxWidth(scale(c, 320)).Padding(0)
		overlayPanel(c, p, radius)
		focus := enter
		scroll := ui.Scroll(c).MaxHeight(scale(c, 420)).MinWidth(0).Shrink(0)
		scroll.Children(func() {
			buildMenu(c, o.Menu, menuExtra{close: open, focusFirst: &focus}, nil)
		})
	})
	if panel.Valid() {
		placement := o.Placement
		if placement.Offset == 0 {
			placement.Offset = 4
		}
		panel = anchored.Apply(panel, trigger, placement, mygoui.SideBottom)
		entryFade(c, panel, panelEntryMotion)
	}
	return trigger
}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func removeKey(keys *[]string, key string) {
	v := *keys
	for i, k := range v {
		if k == key {
			*keys = append(v[:i], v[i+1:]...)
			return
		}
	}
}
