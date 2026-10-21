package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type EmptyOptions struct {
	Icon               *myui.SVG
	Title, Description string
	Media, Actions     func()
	Outline, Muted     bool
}

// Empty composes native media, readable explanatory text and actual caller
// actions. It never invents an enabled placeholder action.
// Source: https://ui.shadcn.com/r/styles/base-rhea/empty.json
func Empty(c *myui.Context, o EmptyOptions) myui.Element {
	t := theme.Colors(c)
	e := myui.Column(c).FillWidth().MinWidth(0).Padding(theme.Space24).
		AlignItems(myui.Center).Justify(myui.Center).Gap(theme.Space16).Radius(theme.Radii().XXXL)
	if o.Outline {
		e.Border(theme.HairlineWidth, t.Border).BorderStyle(myui.BorderDashed)
	}
	if o.Muted {
		e.Background(t.Muted)
	}
	e.Children(func() {
		myui.Column(c).FillWidth().MaxWidth(384).AlignItems(myui.Center).Gap(theme.Space8).Children(func() {
			if o.Media != nil {
				o.Media()
			} else if o.Icon != nil {
				myui.Row(c).Center().Size(theme.Space20*2, theme.Space20*2).
					Background(t.Muted).TextColor(t.Foreground).Radius(theme.Radii().XL).Children(func() {
					myui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconLarge), theme.Scale(c, theme.IconLarge))
				})
			}
			if o.Title != "" {
				Typography(c, o.Title, Heading4).TextAlign(myui.Center)
			}
			if o.Description != "" {
				Typography(c, o.Description, Muted).FillWidth().TextAlign(myui.Center)
			}
		})
		if o.Actions != nil {
			myui.Row(c).MaxWidth(384).MinWidth(0).Wrap().Justify(myui.Center).Gap(theme.Space8).Children(o.Actions)
		}
	})
	return e
}
