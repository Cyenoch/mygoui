package navigation

import "github.com/egoist/mygo/ui"

// ContextMenu attaches a native context menu to anchor. MyGo positions it at
// the secondary click and also opens it with the platform menu key or
// Shift+F10. The operating system owns its appearance and submenu placement.
func ContextMenu(anchor ui.Element, o MenuOptions) ui.Element {
	return anchor.ContextMenu(contextMenuBuilder(anchor, o))
}

// PopoverMenu attaches a context menu drawn inside the window using the current
// MyGo theme. It uses the same secondary-click and keyboard triggers as
// ContextMenu, with native focus restoration, type-ahead and submenu navigation.
// Attach only one context menu to an anchor. Like ContextMenu, it supports
// actions, selection, disabled items, separators and submenus, not section
// headings or item icons, descriptions, shortcut strings and danger styling.
func PopoverMenu(anchor ui.Element, o MenuOptions) ui.Element {
	return anchor.PopoverMenu(contextMenuBuilder(anchor, o))
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
		switch {
		case sc.selection == MenuSelectSingle:
			entry.Radio(checked)
		case sc.selection == MenuSelectMultiple || item.Checked:
			entry.Checked(checked)
		}
		if entry.Chosen() {
			activateMenuItem(sc, true, key, item)
		}
	}
}
