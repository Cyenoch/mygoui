package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type AlertVariant uint8

const (
	AlertDefault AlertVariant = iota
	AlertDestructive
)

type AlertOptions struct {
	Title, Description string
	Variant            AlertVariant
	Icon               *ui.SVG
	Action             *ButtonOptions
}

// Alert is a wrapping callout. Its optional action is a real shared native
// button; title and description remain readable at the caller's text scale.
func Alert(c *ui.Context, o AlertOptions) ui.Element {
	t := theme.Colors(c)
	color := t.CardForeground
	if o.Variant == AlertDestructive {
		color = t.Destructive
	}
	return ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Wrap().
		Gap(theme.Space12).Padding(theme.Space16).Radius(theme.Radii().XXL).
		Border(theme.HairlineWidth, t.Border).Background(t.Card).TextColor(color).
		Label(o.Title).Description(o.Description).Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconMedium), theme.Scale(c, theme.IconMedium)).Shrink(0)
		}
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space4).Children(func() {
			if o.Title != "" {
				ui.Text(c, o.Title).FontSize(theme.Scale(c, theme.Body().Size)).Bold()
			}
			if o.Description != "" {
				ui.Text(c, o.Description).FontSize(theme.Scale(c, theme.Small().Size)).Selectable()
			}
		})
		if o.Action != nil {
			Button(c, *o.Action)
		}
	})
}
