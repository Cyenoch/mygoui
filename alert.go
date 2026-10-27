package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AlertOptions configures a coss ui status callout.
type AlertOptions struct {
	Title, Description string
	Status             StatusColor
	Icon               *ui.SVG
	// ShowIcon uses the status glyph when Icon is nil.
	ShowIcon bool
	Action   func()
}

// Alert renders a status callout with an optional icon and action.
// It exposes ui.RoleStatus; MyGo does not have a separate alert role.
func Alert(c *ui.Context, o AlertOptions) ui.Element {
	t := theme.Colors(c)
	border, fill, iconInk := t.Border, ui.Transparent, t.MutedForeground
	if o.Status != ColorDefault {
		status := StatusColors(c, o.Status)
		border, fill, iconInk = status.Solid.Alpha(0.32), status.Solid.Alpha(0.04), status.Solid
	} else if c.Theme().Dark {
		fill = t.Input.Alpha(0.32)
	}
	e := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space8).
		Padding(theme.Space12, theme.Space14).
		Radius(theme.RadiusScale().XL).Background(fill).Border(theme.Hairline, border).
		TextColor(t.SurfaceForeground).Role(ui.RoleStatus).Label(o.Title).Description(o.Description)
	e.Children(func() {
		icon := o.Icon
		if icon == nil && o.ShowIcon {
			icon = alertIcon(o.Status)
		}
		if icon != nil {
			ui.Row(c).PaddingY(theme.Space2).Shrink(0).Children(func() {
				ui.Icon(c, icon).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
			}).TextColor(iconInk)
		}
		ui.Row(c).Grow(1).Basis(0).MinWidth(0).Wrap().Gap(theme.Space8).
			AlignItems(ui.Center).Children(func() {
			if o.Title != "" || o.Description != "" {
				ui.Column(c).Grow(1).Basis(theme.Scale(c, 120)).MinWidth(0).Gap(theme.Space2).
					AlignItems(ui.Start).Children(func() {
					if o.Title != "" {
						ts := theme.SM()
						ts.Weight = theme.WeightMedium
						style.Text(ui.Text(c, o.Title), c, ts).TextColor(t.Foreground)
					}
					if o.Description != "" {
						style.Text(ui.Text(c, o.Description), c, theme.SM()).
							TextColor(t.MutedForeground).Selectable()
					}
				})
			}
			if o.Action != nil {
				ui.Row(c).Shrink(0).AlignItems(ui.Center).Children(o.Action)
			}
		})
	})
	return e
}

func alertIcon(color StatusColor) *ui.SVG {
	switch color {
	case ColorAccent, ColorInfo:
		return lucide.Info()
	case ColorSuccess:
		return lucide.CircleCheck()
	case ColorWarning:
		return lucide.TriangleAlert()
	case ColorDanger:
		return lucide.CircleAlert()
	}
	return lucide.Info()
}
