package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// ButtonVariant is HeroUI's button color variant.
type ButtonVariant uint8

const (
	ButtonPrimary ButtonVariant = iota
	ButtonSecondary
	ButtonTertiary
	ButtonOutline
	ButtonGhost
	ButtonDanger
	ButtonDangerSoft
)

// ButtonAlign places button content.
type ButtonAlign uint8

const (
	ButtonAlignCenter ButtonAlign = iota
	ButtonAlignStart
	ButtonAlignEnd
)

// ButtonOptions configures a native button. Label is always the accessible
// name, including icon-only layouts and custom content.
type ButtonOptions struct {
	Label     string
	Icon      *ui.SVG
	EndIcon   *ui.SVG
	Variant   ButtonVariant
	Size      theme.ButtonSize
	IconOnly  bool
	FullWidth bool
	Loading   bool
	Disabled  bool
	OnClick   func()
	Handle    *ui.Handle
	Tooltip   string
	Content   func()

	// Align overrides the centered label placement, for navigation rows and
	// other left-aligned actions.
	Align ButtonAlign

	// Corners overrides the uniform HeroUI radius. Nil keeps the default;
	// a zero array makes square corners.
	Corners *[4]float32
}

// Button constructs a HeroUI button on the native focusable button base.
// Hover and pressed paints follow live pointer state; pending buttons do not
// activate and keep their full opacity. HeroUI's pressed scale is not ported
// because MyGo layouts do not transform live element boxes.
func Button(c *ui.Context, o ButtonOptions) ui.Element {
	t := theme.Colors(c)
	disabled := o.Disabled || o.Loading
	height := theme.ButtonHeight(o.Size)
	radius := theme.RadiusScale().XXXL
	paddingX := float32(theme.Space16)
	glyph := float32(theme.IconMD)
	font := theme.SM()
	font.Weight = theme.WeightMedium
	if o.Size == theme.ButtonSmall {
		paddingX = theme.Space12
		glyph = theme.IconMD
	}
	if o.Size == theme.ButtonLarge {
		font = theme.Base()
		font.Weight = theme.WeightMedium
	}
	if o.IconOnly {
		paddingX = 0
	}

	e := ui.ButtonBase(c).MinHeight(theme.Scale(c, height)).Gap(theme.Space8).
		PaddingX(paddingX).PaddingY(0).Center().Shrink(0).
		Radius(radius).FocusRing(true).Label(o.Label)
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	switch o.Align {
	case ButtonAlignStart:
		e.Justify(ui.Start)
	case ButtonAlignEnd:
		e.Justify(ui.End)
	}
	if o.IconOnly {
		e.Width(theme.Scale(c, height)).MinWidth(theme.Scale(c, height))
	}
	e = style.Disabled(e, o.Disabled)
	if o.Handle != nil {
		e.Bind(o.Handle)
	}

	bg, hover, pressed, fg := buttonColors(t, o.Variant)
	if disabled {
		fg = fg.Alpha(0.7)
	}
	e.TextColor(fg)
	if o.Variant == ButtonOutline {
		e.Border(theme.Hairline, t.Border)
	}
	if o.OnClick != nil && !disabled {
		e.OnClick(o.OnClick)
	}
	if o.Tooltip != "" {
		e.Tooltip(o.Tooltip)
	} else if o.IconOnly && o.Label != "" {
		e.Tooltip(o.Label)
	}

	content := func() {
		if o.Loading {
			loading.Arc(c, loading.Options{
				Size: theme.Scale(c, glyph), Color: fg,
			}).Label(locale.M().Loading)
		} else if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, glyph), theme.Scale(c, glyph)).Shrink(0)
		}
		if o.Content != nil {
			o.Content()
		} else if !o.IconOnly {
			style.Text(ui.Text(c, o.Label), c, font).SingleLine()
		}
		if o.EndIcon != nil && !o.IconOnly {
			ui.Icon(c, o.EndIcon).Size(theme.Scale(c, glyph), theme.Scale(c, glyph)).Shrink(0)
		}
	}
	if o.Corners != nil {
		e = style.CornerFill(c, e, theme.Space8, *o.Corners, bg, hover, pressed, content)
	} else {
		e = style.Background(e, radius, bg, hover, pressed)
		e.Children(content)
	}

	return e
}

func buttonColors(t theme.Palette, v ButtonVariant) (bg, hover, pressed, fg ui.Color) {
	switch v {
	case ButtonSecondary:
		return t.Default, t.DefaultHover, t.DefaultHover, t.AccentSoftForeground
	case ButtonTertiary:
		return t.Default, t.DefaultHover, t.DefaultHover, t.DefaultForeground
	case ButtonOutline:
		return ui.Transparent, t.Default.Alpha(0.6), t.Default.Alpha(0.6), t.DefaultForeground
	case ButtonGhost:
		return ui.Transparent, t.Default, t.Default, t.DefaultForeground
	case ButtonDanger:
		return t.Danger, t.DangerHover, t.DangerHover, t.DangerForeground
	case ButtonDangerSoft:
		return t.DangerSoft, t.DangerSoftHover, t.DangerSoftHover, t.DangerSoftForeground
	}
	return t.Accent, t.AccentHover, t.AccentHover, t.AccentForeground
}

// ButtonGroupOrientation selects the join axis.
type ButtonGroupOrientation uint8

const (
	GroupHorizontal ButtonGroupOrientation = iota
	GroupVertical
)

// ButtonGroupOptions controls a joined set of HeroUI buttons.
type ButtonGroupOptions struct {
	Orientation ButtonGroupOrientation
	FullWidth   bool
	Separators  bool
}

// ButtonGroup builds buttons with only the outer corners rounded, as HeroUI
// does. The caller supplies complete button options, avoiding an ambient
// context and keeping each action independently focusable.
func ButtonGroup(c *ui.Context, o ButtonGroupOptions, buttons ...ButtonOptions) ui.Element {
	var e ui.Element
	if o.Orientation == GroupVertical {
		e = ui.Column(c).Gap(0).MinWidth(0)
	} else {
		e = ui.Row(c).Gap(0).MinWidth(0)
	}
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	radius := theme.RadiusScale().XXXL
	n := len(buttons)
	e.Children(func() {
		for i := range buttons {
			b := buttons[i]
			switch o.Orientation {
			case GroupVertical:
				switch {
				case n == 1:
					b.Corners = &[4]float32{radius, radius, radius, radius}
				case i == 0:
					b.Corners = &[4]float32{radius, radius, 0, 0}
				case i == n-1:
					b.Corners = &[4]float32{0, 0, radius, radius}
				default:
					b.Corners = &[4]float32{}
				}
			default:
				switch {
				case n == 1:
					b.Corners = &[4]float32{radius, radius, radius, radius}
				case i == 0:
					b.Corners = &[4]float32{radius, 0, 0, radius}
				case i == n-1:
					b.Corners = &[4]float32{0, radius, radius, 0}
				default:
					b.Corners = &[4]float32{}
				}
			}
			if o.FullWidth && o.Orientation == GroupVertical {
				b.FullWidth = true
			}
			if o.Separators && i > 0 {
				separatorColor := theme.Colors(c).Border
				if o.Orientation == GroupVertical {
					ui.Box(c).Height(theme.Hairline).FillWidth().Background(separatorColor)
				} else {
					ui.Box(c).Width(theme.Hairline).MinHeight(theme.Space20).Background(separatorColor)
				}
			}
			item := Button(c, b)
			if o.FullWidth && o.Orientation == GroupHorizontal {
				item.Grow(1).Basis(0).MinWidth(0)
			}
		}
	})
	return e
}
