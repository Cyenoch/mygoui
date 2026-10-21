package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type ButtonGroupOptions struct {
	Vertical bool
	Label    string
}

// ButtonGroup joins real buttons with one shared border per seam, retaining
// individual Tab stops. A concrete button slice lets native layout set the end
// corners without inspecting children or installing an ambient slot context.
// Source: https://ui.shadcn.com/r/styles/base-rhea/button-group.json
func ButtonGroup(c *myui.Context, o ButtonGroupOptions, buttons ...ButtonOptions) myui.Element {
	group := myui.Row(c).AlignItems(myui.Stretch).Role(myui.RoleGroup).Label(o.Label)
	if o.Vertical {
		group.Column()
	}
	group.Children(func() {
		for i, options := range buttons {
			c.Key(i)
			b := Button(c, options)
			joinControl(b, i, len(buttons), o.Vertical)
		}
	})
	return group
}

func joinControl(b myui.Element, index, count int, vertical bool) {
	r := theme.Radii().XXL
	var tl, tr, br, bl float32
	if vertical {
		if index == 0 {
			tl, tr = r, r
		}
		if index == count-1 {
			bl, br = r, r
		}
		if index > 0 {
			b.BorderWidth(0, 1, 1, 1)
		}
	} else {
		if index == 0 {
			tl, bl = r, r
		}
		if index == count-1 {
			tr, br = r, r
		}
		if index > 0 {
			b.BorderWidth(1, 1, 1, 0)
		}
	}
	b.Radius(tl, tr, br, bl)
}

func ButtonGroupText(c *myui.Context, label string) myui.Element {
	t := theme.Colors(c)
	return myui.Row(c).Center().MinHeight(theme.HeightInput).PaddingX(theme.Space10).
		Background(t.Muted).Border(theme.HairlineWidth, t.Border).Radius(theme.Radii().XXL).
		Children(func() { Typography(c, label, Paragraph) })
}

func ButtonGroupSeparator(c *myui.Context, vertical bool) myui.Element {
	return Separator(c, SeparatorOptions{Vertical: vertical}).Background(theme.Colors(c).Input)
}
