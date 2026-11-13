package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AlertOptions configures a status callout.
type AlertOptions struct {
	Title, Description string
	Status             StatusColor
	Icon               *ui.SVG
	// ShowIcon uses the status glyph when Icon is nil.
	ShowIcon bool
	Action   func()
}

// messageRole is how urgently a message of status is announced.
func messageRole(status StatusColor) ui.Role {
	if status == ColorDanger {
		return ui.RoleAlert
	}
	return ui.RoleStatus
}

// Alert renders a status callout with an optional icon and action. A danger
// alert is an alert, which screen readers announce at once; the others are a
// polite status.
func Alert(c *ui.Context, o AlertOptions) ui.Element {
	t := theme.Colors(c)
	border, fill, iconInk := t.Border, ui.Transparent, t.MutedForeground
	if o.Status != ColorDefault {
		status := statusColors(&t, o.Status)
		border, fill, iconInk = status.Solid.Alpha(0.32), status.Solid.Alpha(0.04), status.Solid
	} else if c.Theme().Dark {
		fill = t.Input.Alpha(0.32)
	}
	e := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space(2)).
		Padding(theme.Space(3), theme.Space(3.5)).
		Radius(theme.RadiusScale().XL).Background(fill).Border(theme.Hairline, border).
		TextColor(t.SurfaceForeground).Role(messageRole(o.Status)).Label(o.Title).Description(o.Description).
		Transition(appearTransition())
	e.Children(func() {
		icon := o.Icon
		if icon == nil && o.ShowIcon {
			icon = StatusIcon(o.Status)
		}
		if icon != nil {
			ui.Row(c).PaddingY(theme.Space(0.5)).Shrink(0).Children(func() {
				ui.Icon(c, icon).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
			}).TextColor(iconInk)
		}
		ui.Row(c).Grow(1).Basis(0).MinWidth(0).Wrap().Gap(theme.Space(2)).
			AlignItems(ui.Center).Children(func() {
			if o.Title != "" || o.Description != "" {
				ui.Column(c).Grow(1).Basis(theme.Scale(c, 120)).MinWidth(0).Gap(theme.Space(0.5)).
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
