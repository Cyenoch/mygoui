package data

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// buildEmpty builds a table with no rows: the same native header a populated
// table shows, so sorting, column resizing and the caller's column order keep
// working, above the empty slot.
//
// A native table without rows has no scroll area for the slot, so the header
// is sized to its own height and the slot follows it. Resetting the native
// selection pointers keeps a zero-row pass from consuming the caller's
// chooser; the next populated pass wires them again.
func buildEmpty(c *ui.Context, o TableOptions, st *TableState, tl *tableLayout, chrome func(p *ui.Painter, r ui.Rect)) {
	t := tl.colors
	board := ui.Column(c).Grow(1).MinWidth(0)
	if chrome != nil {
		board.Draw(chrome)
	}
	if line := tableHeaderLine(tl.colors, len(tl.columns) > 0 || tl.actions >= 0); line != nil {
		board.DrawOver(line)
	}
	board.Children(func() {
		if len(tl.columns) > 0 || tl.actions >= 0 {
			st.List.Sort = &st.nativeSort
			st.List.Selected = nil
			st.List.Selection = nil
			ui.Table(c, &st.List, tl.nativeColumns(), 0, func(row, col int) {}).
				Height(theme.Space(8) + theme.Hairline).
				Shrink(0).MinWidth(0).
				Label(o.Label)
		}
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
