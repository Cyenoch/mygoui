package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type BadgeVariant uint8

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeDestructive
	BadgeOutline
	BadgeGhost
	BadgeLink
)

type BadgeOptions struct {
	Label              string
	Icon, TrailingIcon *myui.SVG
	Variant            BadgeVariant
	Disabled           bool
	Invalid            bool
	OnClick            func()
	Handle             *myui.Handle
}

// Badge is text by default; supplying an action creates a real focusable button.
// Source: https://ui.shadcn.com/r/styles/base-rhea/badge.json
func Badge(c *myui.Context, o BadgeOptions) myui.Element {
	var b myui.Element
	if o.OnClick != nil {
		b = myui.ButtonBase(c).OnClick(o.OnClick).Disabled(o.Disabled)
	} else {
		b = myui.Row(c).Center().Shrink(0).Disabled(o.Disabled)
	}
	b.Label(o.Label).MinHeight(max(theme.Space20, theme.Scale(c, theme.Small().Size)+theme.Space4)).
		Padding(2, theme.Space8).Gap(theme.Space4).Radius(theme.Radii().XXL).
		Font(theme.FontFamily(c)).FontSize(theme.Scale(c, theme.Small().Size)).FontWeight(500)
	if o.Handle != nil {
		b.Bind(o.Handle)
	}
	variant := ButtonDefault
	switch o.Variant {
	case BadgeSecondary:
		variant = ButtonSecondary
	case BadgeDestructive:
		variant = ButtonDestructive
	case BadgeOutline:
		variant = ButtonOutline
	case BadgeGhost:
		variant = ButtonGhost
	case BadgeLink:
		variant = ButtonLink
	}
	paintAction(c, b, variant, o.Disabled)
	if o.Invalid {
		b.BorderColor(theme.Colors(c).Destructive).Error(o.Label)
	}
	b.Children(func() {
		if o.Icon != nil {
			myui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconMicro), theme.Scale(c, theme.IconMicro))
		}
		text := myui.Text(c, o.Label).SingleLine()
		if o.Variant == BadgeLink && o.OnClick != nil && b.Hovered() {
			text.Underline()
		}
		if o.TrailingIcon != nil {
			myui.Icon(c, o.TrailingIcon).Size(theme.Scale(c, theme.IconMicro), theme.Scale(c, theme.IconMicro))
		}
	})
	return b
}
