package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

// Kbd displays a key hint, without adding an input target or fake shortcut.
// Source: https://ui.shadcn.com/r/styles/base-rhea/kbd.json
func Kbd(c *myui.Context, label string) myui.Element {
	t := theme.Colors(c)
	return myui.Row(c).Center().MinHeight(max(theme.Space20, theme.Scale(c, theme.Small().Size)+theme.Space4)).
		MinWidth(theme.Space20).PaddingX(theme.Space4).Radius(theme.Radii().MD).Background(t.Muted).
		TextColor(t.MutedForeground).Children(func() {
		myui.Text(c, label).SingleLine().Font(theme.FontFamily(c)).
			FontSize(theme.Scale(c, theme.Small().Size)).FontWeight(500)
	})
}

func KbdGroup(c *myui.Context, body func()) myui.Element {
	return myui.Row(c).Gap(theme.Space4).Children(body)
}
