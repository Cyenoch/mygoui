package navigation

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ContextMenu attaches a native context menu to anchor. MyGo positions it at
// the secondary click and also opens it with the platform menu key or
// Shift+F10. The operating system owns its appearance and submenu placement.
func ContextMenu(anchor ui.Element, o MenuOptions) ui.Element {
	return anchor.ContextMenu(contextMenuBuilder(anchor, o))
}

// PopoverMenu attaches a context menu drawn inside the window using the current
// MyGo theme and the library's panel and item radii. It uses the same
// secondary-click and keyboard triggers as ContextMenu, with native focus
// restoration, type-ahead and submenu navigation. Attach only one context menu
// to an anchor. Both support actions, selection, disabled items, separators,
// submenus, descriptions, typed accelerators and danger items; PopoverMenu also
// draws custom shortcut strings. Section headings and item icons are not shown.
func PopoverMenu(anchor ui.Element, o MenuOptions) ui.Element {
	return anchor.PopoverMenu(contextMenuBuilder(anchor, o), ui.PopoverMenuOptions{
		Panel: func(panel ui.Element) { panel.Radius(theme.RadiusScale().LG) },
		Item:  func(item ui.Element) { item.Radius(theme.RadiusScale().SM) },
	})
}

func contextMenuBuilder(anchor ui.Element, o MenuOptions) func(*ui.Menu) {
	if o.Selection != MenuSelectNone && o.Selected == nil {
		o.Selected = ui.Local(anchor, "context-menu-selected", func() []string { return nil })
	}
	return func(m *ui.Menu) {
		sc := &menuScope{selection: o.Selection, selected: o.Selected, onSelect: o.OnSelect}
		for _, section := range o.Sections {
			contextMenuItems(m, section.Items, sc)
		}
	}
}

func contextMenuItems(m *ui.Menu, items []MenuItemOptions, sc *menuScope) {
	for _, item := range items {
		if item.Separator {
			m.Separator()
			continue
		}
		if len(item.Submenu) != 0 {
			m.Submenu(item.Label, func(sub *ui.Menu) {
				contextMenuItems(sub, item.Submenu, sc)
			}).Disabled(item.Disabled)
			continue
		}
		key := item.Key
		if key == "" {
			key = item.Label
		}
		checked := item.Checked
		if sc.selection != MenuSelectNone && sc.selected != nil {
			checked = containsKey(*sc.selected, key)
		}
		entry := m.Item(item.Label).Disabled(item.Disabled)
		if item.Key != "" {
			entry.ID(item.Key)
		}
		if item.Description != "" {
			entry.Subtitle(item.Description)
		}
		switch {
		case item.Shortcut != "":
			entry.Trailing(item.Shortcut)
		case item.Accelerator.Key != 0:
			entry.Shortcut(item.Accelerator.Mods, item.Accelerator.Key)
		}
		if item.Variant == MenuItemDanger {
			entry.Danger()
		}
		switch {
		case sc.selection == MenuSelectSingle:
			entry.Radio(checked)
		case sc.selection == MenuSelectMultiple || item.Checked || item.Checkable:
			entry.Checked(checked)
		}
		if entry.Chosen() {
			activateMenuItem(sc, true, key, item)
		}
	}
}
