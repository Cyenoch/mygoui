package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ToggleButtonVariant is HeroUI's toggle surface treatment.
type ToggleButtonVariant uint8

const (
	ToggleDefault ToggleButtonVariant = iota
	ToggleGhost
)

// ToggleButtonOptions configures a persistent pressed-state button. Selected
// is caller-owned; clicking flips it and reports the new value.
type ToggleButtonOptions struct {
	Label     string
	Icon      *ui.SVG
	Variant   ToggleButtonVariant
	Size      theme.ButtonSize
	IconOnly  bool
	FullWidth bool
	Selected  *bool
	Disabled  bool
	OnChange  func(bool)
	Tooltip   string
	// Corners overrides the uniform radius. Nil keeps it.
	Corners *[4]float32
}

// ToggleButton is a focusable native toggle whose selected paint follows the
// caller's bool. It is not a radio: each button owns its state.
func ToggleButton(c *ui.Context, o ToggleButtonOptions) ui.Element {
	t := theme.Colors(c)
	selected := o.Selected != nil && *o.Selected
	height := theme.ButtonHeight(o.Size)
	radius := theme.RadiusScale().XXXL
	paddingX := float32(theme.Space16)
	glyph := float32(theme.IconMD)
	font := theme.SM()
	font.Weight = theme.WeightMedium
	normal, hover, fg := t.Default, t.DefaultHover, t.DefaultForeground
	if o.Variant == ToggleGhost {
		normal, hover, fg = ui.Transparent, t.Default, t.DefaultForeground
	}
	if selected {
		normal, hover, fg = t.AccentSoft, t.AccentSoftHover, t.AccentSoftForeground
	}
	if o.Size == theme.ButtonSmall {
		paddingX, glyph = theme.Space12, theme.IconMD
	}
	if o.Size == theme.ButtonLarge {
		font = theme.Base()
		font.Weight = theme.WeightMedium
	}
	if o.IconOnly {
		paddingX = 0
	}
	e := ui.ButtonBase(c).MinHeight(theme.Scale(c, height)).Gap(theme.Space8).
		PaddingX(paddingX).PaddingY(0).Center().Shrink(0).Radius(radius).
		FocusRing(true).Label(o.Label).TextColor(fg).
		Role(ui.RoleToggleButton).Checked(selected)
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	if o.IconOnly {
		e.Width(theme.Scale(c, height)).MinWidth(theme.Scale(c, height))
	}
	content := func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, glyph), theme.Scale(c, glyph)).Shrink(0)
		}
		if !o.IconOnly {
			style.Text(ui.Text(c, o.Label), c, font).SingleLine()
		}
	}
	if o.Corners != nil {
		e = style.CornerFill(c, e, theme.Space8, *o.Corners, normal, hover, hover, content)
	} else {
		e = style.Background(e, radius, normal, hover, hover)
		e.Children(content)
	}
	if !o.Disabled && o.Selected != nil {
		e.OnClick(func() {
			*o.Selected = !*o.Selected
			if o.OnChange != nil {
				o.OnChange(*o.Selected)
			}
		})
	}
	if o.Tooltip != "" {
		e.Tooltip(o.Tooltip)
	} else if o.IconOnly && o.Label != "" {
		e.Tooltip(o.Label)
	}
	return style.Disabled(e, o.Disabled)
}

// ToggleButtonGroupOptions controls a joined or detached toggle set.
type ToggleButtonGroupOptions struct {
	Orientation Orientation
	FullWidth   bool
	Detached    bool
}

// ToggleButtonGroup lays out caller-owned toggle states. Each option carries
// its own Selected pointer, so single- and multi-select are both explicit Go
// state rather than an ambient group context.
func ToggleButtonGroup(c *ui.Context, o ToggleButtonGroupOptions, buttons ...ToggleButtonOptions) ui.Element {
	var e ui.Element
	if o.Orientation == Vertical {
		e = ui.Column(c).Gap(0).MinWidth(0)
	} else {
		e = ui.Row(c).Gap(0).MinWidth(0)
	}
	if o.Detached {
		e.Gap(theme.Space4)
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
			if !o.Detached {
				switch o.Orientation {
				case Vertical:
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
			}
			if o.FullWidth && o.Orientation == Vertical {
				b.FullWidth = true
			}
			item := ToggleButton(c, b)
			if o.FullWidth && o.Orientation == Horizontal {
				item.Grow(1).Basis(0).MinWidth(0)
			}
		}
	})
	return e
}
