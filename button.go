// Give controls stable identities with Context.Key before construction in collections.
package mygoui

import (
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

type ButtonVariant uint8

const (
	ButtonDefault ButtonVariant = iota
	ButtonOutline
	ButtonSecondary
	ButtonGhost
	ButtonDestructive
	ButtonLink
)

type ButtonSize uint8

const (
	ButtonSizeDefault ButtonSize = iota
	ButtonXS
	ButtonSM
	ButtonLG
	ButtonIcon
	ButtonIconXS
	ButtonIconSM
	ButtonIconLG
)

type ButtonOptions struct {
	Label              string
	Icon, TrailingIcon *myui.SVG
	Variant            ButtonVariant
	Size               ButtonSize
	Disabled, Loading  bool
	Invalid, Rounded   bool
	OnClick            func()
	Handle             *myui.Handle
	Tooltip            string
	Wrap               bool
	// Content replaces the default icons and label during this build. Loading
	// prepends its spinner and still disables activation. Label remains the
	// accessible name; the callback must not retain Context or Element values.
	Content func()
}

// Button handles pointer, Enter, Space, focus and accessibility through ButtonBase.
// Missing actions, loading and disabled controls cannot activate. Icon sizes retain
// Label as their accessible title without drawing it. Wrap allows multiline
// labels within caller width constraints. Content replaces the default body.
// The returned Element is valid only for this build; use Handle for persistent focus.
// Source: https://ui.shadcn.com/r/styles/base-rhea/button.json
func Button(c *myui.Context, o ButtonOptions) myui.Element {
	disabled := o.Disabled || o.Loading || o.OnClick == nil
	b := myui.ButtonBase(c).Label(o.Label).Disabled(disabled)
	if o.Handle != nil {
		b.Bind(o.Handle)
	}
	height, padding, gap, glyph, iconOnly := buttonMetrics(o.Size)
	font := theme.Body().Size
	if o.Size == ButtonXS || o.Size == ButtonIconXS {
		font = theme.Small().Size
	}
	height = max(height, theme.Scale(c, font*1.25)+theme.Space8)
	b.MinHeight(height).Gap(gap).PaddingX(padding).
		Radius(theme.Radii().XXL).Font(theme.FontFamily(c)).FontSize(theme.Scale(c, font)).FontWeight(500)
	if o.Wrap {
		b.MinWidth(0)
	}
	if !iconOnly {
		start, end := padding, padding
		if o.Icon != nil || o.Loading {
			start -= 2
		}
		if o.TrailingIcon != nil {
			end -= 2
		}
		b.Padding(0, end, 0, start)
	}
	if iconOnly {
		b.Size(height, height).Shrink(0).Padding(0).Center()
	}
	if o.Rounded {
		b.Radius(height)
	}
	tooltip := o.Tooltip
	if tooltip == "" && iconOnly {
		tooltip = o.Label
	}
	if tooltip != "" {
		b.Tooltip(tooltip)
	}
	paintAction(c, b, o.Variant, disabled)
	if o.Invalid {
		b.Border(theme.HairlineWidth, theme.Colors(c).Destructive).Error(o.Label)
	}
	if o.Loading {
		b.Description(locales.ComponentLoading())
	}
	b.OnClick(o.OnClick).Children(func() {
		if o.Loading {
			Spinner(c, SpinnerOptions{Size: glyph, Label: locales.ComponentLoading()})
		}
		if o.Content != nil {
			o.Content()
			return
		}
		if !o.Loading && o.Icon != nil {
			myui.Icon(c, o.Icon).Size(theme.Scale(c, glyph), theme.Scale(c, glyph)).Shrink(0)
		}
		if !iconOnly {
			text := myui.Text(c, o.Label).MinWidth(0)
			if !o.Wrap {
				text.SingleLine()
			}
			if o.Variant == ButtonLink && b.Hovered() && !disabled {
				text.Underline()
			}
		}
		if o.TrailingIcon != nil {
			myui.Icon(c, o.TrailingIcon).Size(theme.Scale(c, glyph), theme.Scale(c, glyph)).Shrink(0)
		}
	})
	return b
}

func buttonMetrics(size ButtonSize) (height, padding, gap, glyph float32, icon bool) {
	height, padding, gap, glyph = theme.HeightInput, theme.Space12, theme.Space6, theme.IconMedium
	switch size {
	case ButtonXS, ButtonIconXS:
		height, padding, gap, glyph = theme.HeightChip, theme.Space10, theme.Space4, theme.IconMicro
	case ButtonSM, ButtonIconSM:
		height, padding, gap = theme.HeightButton, theme.Space12, theme.Space4
	case ButtonLG, ButtonIconLG:
		height, padding = theme.HeightRowCompact, theme.Space16
	}
	icon = size >= ButtonIcon
	return
}

func paintAction(c *myui.Context, b myui.Element, variant ButtonVariant, disabled bool) {
	t := theme.Colors(c)
	fill, ink := t.Primary, t.PrimaryForeground
	hover, pressed := t.PrimaryHover, t.PrimaryPressed
	var border myui.Color
	switch variant {
	case ButtonOutline:
		fill, ink, border = t.Background, t.Foreground, t.Border
		if c.Theme().Dark {
			fill = myui.Color{}
		}
		hover, pressed = t.AccentHover, t.AccentPressed
	case ButtonSecondary:
		fill, ink = t.Secondary, t.SecondaryForeground
		hover, pressed = t.SecondaryHover, t.SecondaryPressed
	case ButtonGhost:
		fill, ink = myui.Color{}, t.Foreground
		hover, pressed = t.AccentHover, t.AccentPressed
	case ButtonDestructive:
		strength := float32(0.1)
		if c.Theme().Dark {
			strength = 0.2
		}
		fill, ink = theme.Over(t.Destructive, t.Background, strength), t.Destructive
		hover = theme.Over(t.DestructiveHover, t.Background, 0.2)
		pressed = theme.Over(t.DestructivePressed, t.Background, 0.3)
	case ButtonLink:
		fill, ink = myui.Color{}, t.Primary
		hover, pressed = fill, fill
	}
	if disabled {
		fill, ink = t.DisabledSurface, t.DisabledForeground
	} else if b.Pressed() {
		fill = pressed
	} else if b.Hovered() {
		fill = hover
	}
	b.Background(fill).TextColor(ink).Border(theme.HairlineWidth, border)
}
