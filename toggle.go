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
// caller's bool. It is not a radio: each button owns its state. A press
// scales it like a HeroUI button; joined groups keep their size.
func ToggleButton(c *ui.Context, o ToggleButtonOptions) ui.Element {
	var adj joinAdjust
	if o.Corners != nil {
		adj.corners, adj.hasCorners = *o.Corners, true
	}
	return toggleButtonView(c, o, adj)
}

// toggleButtonView is the builder ToggleButton and ToggleButtonGroup share.
func toggleButtonView(c *ui.Context, o ToggleButtonOptions, adj joinAdjust) ui.Element {
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
		paddingX = theme.Space12
	}
	if o.Size == theme.ButtonLarge {
		font = theme.Base()
		font.Weight = theme.WeightMedium
	}
	if o.IconOnly {
		paddingX = 0
	}
	side := theme.Scale(c, height)
	iconSize := theme.Scale(c, glyph)

	e := ui.ButtonBase(c).MinHeight(side).Gap(theme.Space8).
		PaddingX(paddingX).PaddingY(0).Center().Shrink(0).
		FocusRing(true).Label(o.Label).TextColor(fg).
		Role(ui.RoleToggleButton).Checked(selected)
	if adj.hasCorners {
		e = e.Radius(adj.corners[0], adj.corners[1], adj.corners[2], adj.corners[3])
	} else {
		e = e.Radius(radius)
	}
	// Disable before the faces read their pointer state: a disabled toggle
	// paints its rest face, matching HeroUI's pointer-events:none.
	e = style.Disabled(e, o.Disabled)
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	if o.IconOnly {
		e.Width(side).MinWidth(side)
	}
	content := func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0).Margin(0, -2)
		}
		if !o.IconOnly {
			style.Text(ui.Text(c, o.Label), c, font).SingleLine()
		}
	}
	if adj.hasCorners {
		e = style.CornerFill(c, e, theme.Space8, adj.corners, normal, hover, hover, content)
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
	if !adj.grouped {
		e = pressedPose(e, e.Pressed(), buttonPressScale(o.Size))
	}
	if adj.seam {
		e = joinSeam(e, adj.vertical, c.Direction() == ui.RightToLeft, fg)
	}
	return e
}

// ToggleButtonGroupOptions controls a joined or detached toggle set.
type ToggleButtonGroupOptions struct {
	Orientation Orientation
	FullWidth   bool
	Detached    bool
	// Separators paints HeroUI's divider at each shared edge. Detached sets
	// hide it, as HeroUI does.
	Separators bool
}

// ToggleButtonGroup lays out caller-owned toggle states. Each option carries
// its own Selected pointer, so single- and multi-select are both explicit Go
// state rather than an ambient group context. HeroUI removes the press scale
// from every toggle inside a group, detached ones included.
func ToggleButtonGroup(c *ui.Context, o ToggleButtonGroupOptions, buttons ...ToggleButtonOptions) ui.Element {
	vertical := o.Orientation == Vertical
	var e ui.Element
	if vertical {
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
	rtl := c.Direction() == ui.RightToLeft
	n := len(buttons)
	e.Children(func() {
		for i := range buttons {
			b := buttons[i]
			adj := joinAdjust{grouped: true, vertical: vertical}
			if !o.Detached {
				adj.corners, adj.hasCorners = joinCorners(vertical, i, n, radius, rtl), true
				adj.seam = o.Separators && i > 0
			}
			if o.FullWidth && vertical {
				b.FullWidth = true
			}
			item := toggleButtonView(c, b, adj)
			if o.FullWidth && !vertical {
				item.Grow(1).Basis(0).MinWidth(0)
			}
		}
	})
	return e
}
