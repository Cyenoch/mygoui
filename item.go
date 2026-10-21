package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type ItemVariant uint8

const (
	ItemDefault ItemVariant = iota
	ItemOutline
	ItemMuted
)

type ItemSize uint8

const (
	ItemSizeDefault ItemSize = iota
	ItemSM
	ItemXS
)

type ItemOptions struct {
	Label, Description string
	Icon               *myui.SVG
	Variant            ItemVariant
	Size               ItemSize
	Media, Content     func()
	Actions            func()
	Header, Footer     func()
	OnClick            func()
	Disabled           bool
	Handle             *myui.Handle
}

// Item is a content row with independent actions. With OnClick it becomes a
// native press target; nested action buttons keep their own input and AX role.
// Source: https://ui.shadcn.com/r/styles/base-rhea/item.json
func Item(c *myui.Context, o ItemOptions) myui.Element {
	t := theme.Colors(c)
	var e myui.Element
	if o.OnClick != nil {
		e = myui.ButtonBase(c).OnClick(o.OnClick).Disabled(o.Disabled).Justify(myui.Start)
	} else {
		e = myui.Row(c).Disabled(o.Disabled)
	}
	x, y, gap := float32(theme.Space16), float32(theme.Space14), float32(theme.Space14)
	switch o.Size {
	case ItemSM:
		x, y = theme.Space14, theme.Space12
	case ItemXS:
		x, y, gap = theme.Space10, theme.Space8, theme.Space8
	}
	e.FillWidth().MinWidth(0).Column().AlignItems(myui.Stretch).Padding(y, x).Gap(gap).
		Radius(theme.Radii().XXL).TextColor(t.Foreground).Label(o.Label).Description(o.Description)
	if o.Handle != nil {
		e.Bind(o.Handle)
	}
	if o.Variant == ItemOutline {
		e.Border(theme.HairlineWidth, t.Border)
	}
	if o.Variant == ItemMuted {
		e.Background(theme.Over(t.Muted, t.Background, 0.5))
	}
	if o.OnClick != nil && !o.Disabled {
		if e.Pressed() {
			e.Background(t.AccentPressed)
		} else if e.Hovered() {
			e.Background(t.AccentHover)
		}
	}
	if o.Disabled {
		e.TextColor(t.DisabledForeground)
	}
	e.Children(func() {
		if o.Header != nil {
			myui.Row(c).FillWidth().Gap(theme.Space8).Children(o.Header)
		}
		myui.Row(c).FillWidth().MinWidth(0).Gap(gap).AlignItems(myui.Start).Children(func() {
			if o.Media != nil {
				o.Media()
			} else if o.Icon != nil {
				myui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconMedium), theme.Scale(c, theme.IconMedium))
			}
			myui.Column(c).Grow(1).MinWidth(0).Gap(theme.Space4).Children(func() {
				if o.Content != nil {
					o.Content()
					return
				}
				myui.Text(c, o.Label).Font(theme.FontFamily(c)).FontSize(theme.Scale(c, theme.Body().Size)).
					FontWeight(500).SingleLine()
				if o.Description != "" {
					myui.Text(c, o.Description).FontSize(theme.Scale(c, theme.Body().Size)).
						TextColor(t.MutedForeground).MaxLines(2)
				}
			})
			if o.Actions != nil {
				myui.Row(c).Gap(theme.Space8).Wrap().Children(o.Actions)
			}
		})
		if o.Footer != nil {
			myui.Row(c).FillWidth().Gap(theme.Space8).Children(o.Footer)
		}
	})
	return e
}

func ItemGroup(c *myui.Context, size ItemSize, body func()) myui.Element {
	gap := float32(theme.Space16)
	switch size {
	case ItemSM:
		gap = theme.Space10
	case ItemXS:
		gap = theme.Space8
	}
	return myui.Column(c).FillWidth().MinWidth(0).Gap(gap).Role(myui.RoleList).Children(body)
}

func ItemSeparator(c *myui.Context) myui.Element {
	return Separator(c, SeparatorOptions{}).MarginY(theme.Space8)
}
