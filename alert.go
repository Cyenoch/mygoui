package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AlertOptions configures a HeroUI status callout.
type AlertOptions struct {
	Title, Description string
	Status             StatusColor
	Icon               *ui.SVG
	Action             func()
}

// Alert renders a surface callout with the status indicator and foreground.
func Alert(c *ui.Context, o AlertOptions) ui.Element {
	t := theme.Colors(c)
	ink := StatusColors(c, o.Status).Foreground
	e := ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space16).
		Padding(theme.Space12, theme.Space16, theme.Space12, theme.Space16).
		Radius(min(theme.Space32, theme.RadiusScale().XXXL)).Background(t.Surface).
		TextColor(t.SurfaceForeground).Label(o.Title).Description(o.Description)
	e = style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark))
	e.Children(func() {
		ui.Row(c).Padding(theme.Space4).Shrink(0).Children(func() {
			icon := o.Icon
			if icon == nil {
				icon = alertIcon(o.Status)
			}
			ui.Icon(c, icon).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
		}).TextColor(ink)
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
			if o.Title != "" {
				ts := theme.SM()
				ts.LineHeight, ts.Weight = 24, theme.WeightMedium
				style.Text(ui.Text(c, o.Title), c, ts).TextColor(ink)
			}
			if o.Description != "" {
				ts := theme.SM()
				style.Text(ui.Text(c, o.Description), c, ts).TextColor(t.Muted).Selectable()
			}
			if o.Action != nil {
				o.Action()
			}
		})
	})
	return e
}

func alertIcon(color StatusColor) *ui.SVG {
	switch color {
	case ColorAccent:
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
