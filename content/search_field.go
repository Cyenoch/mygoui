package content

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// searchField is the query field of a viewer's toolbar: a search icon, the
// text bound to query, and a clear button while it holds text. Escape
// clears it before a shortcut around the viewer sees the key, and clearing
// leaves the focus in the field; onKey, when set, sees the other keys and
// reports whether it took one. It is built here rather than by
// form.SearchField because the viewer's shortcut must focus its editor,
// which that field gives no handle to.
func searchField(c *ui.Context, t theme.Palette, query *string, focus *ui.Handle, onKey func(ui.TextKey) bool) {
	m := locale.M()
	row := ui.Row(c).AlignItems(ui.Center).MinHeight(theme.Scale(c, theme.HeightButtonSM)).Grow(1).MinWidth(0)
	row = formstyle.Surface(row, c, false, false, false)
	row.Children(func() {
		ui.Icon(c, lucide.Search()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
			TextColor(t.FieldPlaceholder).Shrink(0).MarginStart(theme.Space(2))
		editor := formstyle.Text(ui.TextInputBase(c, query), c).Bind(focus).Grow(1).MinWidth(0).
			PaddingStart(theme.Space(2)).PaddingEnd(theme.Space(2)).PaddingY(theme.Scale(c, 3)).
			Placeholder(m.Search).Label(m.Search)
		if onKey != nil {
			editor.OnTextKey(onKey)
		}
		if *query == "" {
			return
		}
		clear := func() {
			*query = ""
			focus.Focus(c)
		}
		editor.OnShortcut(0, ui.KeyEscape, clear)
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: m.Clear, Icon: lucide.X(), IconOnly: true, Variant: mygoui.ButtonGhost,
			Size: theme.ButtonSmall, OnClick: clear,
		}).Shrink(0).MarginEnd(theme.Space(1))
	})
}
