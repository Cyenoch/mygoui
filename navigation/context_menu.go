package navigation

import "github.com/egoist/mygo/ui"

// ContextMenu attaches a native context menu to anchor. MyGo positions it at
// the secondary click and also opens it with the platform menu key or
// Shift+F10. The operating system owns its appearance and submenu placement.
func ContextMenu(anchor ui.Element, o MenuOptions) ui.Element {
	if o.Selection != MenuSelectNone && o.Selected == nil {
		o.Selected = ui.Local(anchor, "context-menu-selected", func() []string { return nil })
	}
	return anchor.ContextMenu(func(m *ui.Menu) {
		sc := &menuScope{selection: o.Selection, selected: o.Selected, onSelect: o.OnSelect}
		for _, section := range o.Sections {
			contextMenuItems(m, section.Items, sc)
		}
	})
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
		if sc.selection != MenuSelectNone || item.Checked {
			entry.Checked(checked)
		}
		if entry.Chosen() {
			activateMenuItem(sc, true, key, item)
		}
	}
}
