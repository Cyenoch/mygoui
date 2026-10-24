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

// Kbd renders one HeroUI keycap. Multiple keys form a chord separated by a
// thin gap, with the concatenated chord as its accessible text.
func Kbd(c *ui.Context, o KbdOptions) ui.Element {
	t := theme.Colors(c)
	background := t.Default
	if o.Light {
		background = ui.Transparent
	}
	label := strings.Join(o.Keys, " ")
	e := ui.Row(c).Shrink(0).Center().AlignSelf(ui.Start).Gap(theme.Space2).
		MinHeight(theme.Scale(c, theme.HeightKbd)).
		PaddingX(theme.Space8).Radius(theme.RadiusScale().LG).Background(background).
		TextColor(t.Muted).Label(label)
	e.Children(func() {
		ts := theme.SM()
		ts.Weight = theme.WeightMedium
		for _, key := range o.Keys {
			style.Text(ui.Text(c, key), c, ts).SingleLine()
		}
	})
	return e
}
