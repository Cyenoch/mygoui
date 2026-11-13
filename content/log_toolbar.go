package content

import (
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// logToolbar is the search field and the toggles above the lines.
func logToolbar(c *ui.Context, t theme.Palette, st *LogViewerState, o *LogViewerOptions, total int) {
	m := locale.M()
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(2))).
		Padding(theme.Scale(c, theme.Space(1.5)), theme.Scale(c, theme.Space(2))).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).
		Children(func() {
			searchField(c, t, &st.Query, &st.search, nil)
			// A button keeps its own alignment, so the row around them
			// centers them against the field.
			ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Scale(c, theme.Space(1))).Children(func() {
				toggle := func(label string, icon *ui.SVG, on bool, set func(bool)) {
					mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{
						Label: label, Icon: icon, IconOnly: true, Variant: mygoui.ToggleOutline,
						Size: theme.ButtonSmall, Selected: &on, OnChange: set,
					})
				}
				if o.Level != nil {
					for _, l := range []struct {
						level LogLevel
						label string
						icon  *ui.SVG
					}{
						{LogDebug, m.LevelDebug, lucide.Bug()},
						{LogInfo, m.LevelInfo, lucide.Info()},
						{LogWarning, m.LevelWarning, lucide.TriangleAlert()},
						{LogError, m.LevelError, lucide.CircleX()},
					} {
						toggle(l.label, l.icon, st.Levels.has(l.level), func(on bool) {
							st.Levels = st.Levels.toggled(l.level, on)
						})
					}
				}
				toggle(m.WrapLines, lucide.WrapText(), !st.NoWrap, func(on bool) { st.NoWrap = !on })
				listed := st.refresh(o, total)
				if st.filtering() {
					ts := theme.Code()
					ts.Size, ts.LineHeight = codeFontSize-1, codeLineHeight
					style.Text(ui.Text(c.Key("count"), groupDigits(listed)+" / "+groupDigits(total)), c, ts).
						TextColor(t.Muted).Shrink(0).MarginStart(theme.Scale(c, theme.Space(2))).Unselectable().
						Role(ui.RoleStatus)
				}
			})
		})
}

// groupDigits writes n with a comma between thousands.
func groupDigits(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
