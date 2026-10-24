package data

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// buildEmpty builds the header of a table with no rows and the empty slot.
// The native table is not built without rows because its scroll area cannot
// hold an empty message; this header keeps the columns visible and sortable.
func buildEmpty(c *ui.Context, o TableOptions, st *TableState, tl *tableLayout) {
	t := theme.Colors(c)
	board := ui.Column(c).Grow(1).MinWidth(0).Role(ui.RoleTable)
	if o.Variant == TablePrimary {
		board.Background(t.Surface).TextColor(t.SurfaceForeground)
	}
	if o.Label != "" {
		board.Label(o.Label)
	}
	board.Children(func() {
		head := ui.Row(c).Height(theme.Space(8)).Shrink(0).AlignItems(ui.Stretch).Role(ui.RoleRow)
		head.Children(func() {
			for _, cb := range tl.columns {
				emptyHeader(c, o, st, cb.column)
			}
			if tl.actions >= 0 {
				w := o.ActionsWidth
				if w <= 0 {
					w = 64
				}
				box := ui.Row(c).Width(w).MinWidth(w).Shrink(0).
					Padding(theme.Space(1.5), theme.Space(2.5)).
					AlignItems(ui.Center).Justify(ui.End).
					Role(ui.RoleColumnHeader)
				if o.ActionsHeader != "" {
					box.Label(o.ActionsHeader)
					box.Children(func() {
						emptyHeaderText(c, o.ActionsHeader, st.Sort.Column == actionsColumnID, st.Sort.Descending)
					})
				}
			}
		})
		ui.Divider(c)
		body := ui.Column(c).Grow(1).MinWidth(0).
			AlignItems(ui.Center).Justify(ui.Center).
			Padding(theme.Space24)
		body.Children(func() {
			if o.Empty != nil {
				o.Empty(c)
				return
			}
			style.Text(ui.Text(c, locale.M().Empty), c, theme.SM()).
				TextColor(t.Muted).Role(ui.RoleStatus)
		})
	})
}

func emptyHeader(c *ui.Context, o TableOptions, st *TableState, col TableColumn) {
	box := ui.Row(c).Height(theme.Space(8)).Shrink(0).
		Padding(theme.Space(1.5), theme.Space(2.5)).
		AlignItems(ui.Center).Role(ui.RoleColumnHeader)
	if col.Width > 0 {
		box.Width(col.Width).MinWidth(theme.Space(4))
	} else {
		box.Grow(1).Shrink(1).MinWidth(max(col.MinWidth, theme.Space(15)))
	}
	switch col.Align {
	case ui.Center:
		box.Justify(ui.Center)
	case ui.End:
		box.Justify(ui.End)
	}
	label := col.Header
	if label == "" {
		label = col.ID
	}
	if label != "" {
		box.Label(label)
	}
	if col.Sortable {
		box.Focusable().FocusRing(true).Cursor(ui.CursorPointer)
		id := col.ID
		box.OnClick(func() { toggleSort(st, o, id) })
	}
	box.Children(func() {
		emptyHeaderText(c, col.Header, col.Sorted, col.Descending)
	})
}

func emptyHeaderText(c *ui.Context, text string, sorted, descending bool) {
	t := theme.Colors(c)
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	style.Text(ui.Text(c, text), c, ts).SingleLine().TextColor(t.Muted)
	if !sorted {
		return
	}
	arrow := "▲"
	if descending {
		arrow = "▼"
	}
	style.Text(ui.Text(c, arrow), c, theme.XS()).
		TextColor(t.Muted).Role(ui.RoleNone).
		MarginStart(theme.Space4).Shrink(0)
}

// toggleSort sorts by id, reversing it when it was already first, as a header
// click does in the native table. It is used by the empty table's header.
func toggleSort(st *TableState, o TableOptions, id string) {
	if st.Sort.Column == id {
		st.Sort.Descending = !st.Sort.Descending
	} else {
		st.Sort = TableSort{Column: id}
	}
	st.nativeSort = ui.SortOrder{Column: st.Sort.Column, Descending: st.Sort.Descending}
	st.lastSort = st.Sort
	if o.OnSort != nil {
		o.OnSort(st.Sort)
	}
}
