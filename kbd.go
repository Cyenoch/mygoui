package mygoui

import (
	"strings"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// KbdOptions configures a keyboard chord.
type KbdOptions struct {
	Keys []string
	// Light removes the default background.
	Light bool
}

// Kbd renders one coss ui keycap. Multiple keys form a chord separated by a
// thin gap, with the concatenated chord as its accessible text.
func Kbd(c *ui.Context, o KbdOptions) ui.Element {
	t := theme.Colors(c)
	background := t.MutedBackground
	if o.Light {
		background = ui.Transparent
	}
	label := strings.Join(o.Keys, " ")
	e := ui.Row(c).Shrink(0).Center().AlignSelf(ui.Start).Gap(theme.Space4).
		MinHeight(theme.Scale(c, theme.HeightKbd)).MinWidth(theme.Scale(c, theme.HeightKbd)).
		PaddingX(theme.Space4).Radius(4).Background(background).
		TextColor(t.MutedForeground).Label(label)
	e.Children(func() {
		ts := theme.XS()
		ts.Weight = theme.WeightMedium
		for _, key := range o.Keys {
			style.Text(ui.Text(c, key), c, ts).SingleLine()
		}
	})
	return e
}

// KbdGroup renders each key in a chord as a separate keycap.
func KbdGroup(c *ui.Context, o KbdOptions) ui.Element {
	e := ui.Row(c).Shrink(0).AlignSelf(ui.Start).AlignItems(ui.Center).
		Gap(theme.Space4).Label(strings.Join(o.Keys, " "))
	return e.Children(func() {
		for _, key := range o.Keys {
			Kbd(c, KbdOptions{Keys: []string{key}, Light: o.Light})
		}
	})
}
