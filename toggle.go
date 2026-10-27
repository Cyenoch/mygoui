package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ToggleButtonVariant selects a toggle surface, including a native ghost variant.
type ToggleButtonVariant uint8

const (
	ToggleDefault ToggleButtonVariant = iota
	ToggleGhost
	ToggleOutline
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
	radius := theme.RadiusScale().LG
	paddingX := float32(7)
	glyph := float32(theme.IconMD)
	font := theme.SM()
	font.Weight = theme.WeightMedium
	normal, hover, fg := ui.Transparent, t.Accent, t.Foreground
	if o.Variant == ToggleGhost {
		normal, hover, fg = ui.Transparent, t.Default, t.DefaultForeground
	}
	if selected {
		normal, hover, fg = t.Input.Alpha(0.64), t.Input, t.AccentForeground
	}
	if o.Size == theme.ButtonSmall {
		paddingX = 5
	}
	if o.Size == theme.ButtonLarge {
		paddingX = 9
		font = theme.SM()
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
	if o.Variant == ToggleOutline {
		if adj.hasBorders {
			e.BorderWidth(adj.borders[0], adj.borders[1], adj.borders[2], adj.borders[3]).BorderColor(t.Input)
		} else {
			e.Border(theme.Hairline, t.Input)
		}
		if !selected {
			normal = t.Background
			e = style.Shadow(e, c, theme.FieldShadows(c.Theme().Dark))
		}
		if c.Theme().Dark {
			normal, hover = t.Input.Alpha(0.32), t.Input.Alpha(0.64)
			if selected {
				normal, hover = t.Input, t.Input
			}
		}
	}
	// Disable before the faces read their pointer state: a disabled toggle
	// paints its rest face, matching coss ui's pointer-events:none.
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
		e = style.CornerFill(c, e, theme.Space8, adj.corners, normal, hover, hover, content, nil)
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
	if adj.seam {
		e = joinSeam(e, adj.vertical, c.Direction() == ui.RightToLeft, fg)
	}
	return e
}

// ToggleButtonGroupOptions controls a joined or detached toggle set.
type ToggleButtonGroupOptions struct {
	Orientation Orientation
	FullWidth   bool
	// Detached keeps a 4 DIP gap even for outline toggles.
	Detached bool
	// Separators paints a divider at the shared edges of joined outline toggles.
	Separators bool
}

// ToggleButtonGroup lays out caller-owned toggle states. Each option carries
// its own Selected pointer, so single- and multi-select are both explicit Go
// state rather than an ambient group context. Outline-only groups join their
// items; other groups use independent rounded items with a 2 DIP gap.
func ToggleButtonGroup(c *ui.Context, o ToggleButtonGroupOptions, buttons ...ToggleButtonOptions) ui.Element {
	vertical := o.Orientation == Vertical
	joined := !o.Detached && len(buttons) > 0
	for _, b := range buttons {
		if b.Variant != ToggleOutline {
			joined = false
			break
		}
	}
	gap := float32(theme.Space2)
	if o.Detached {
		gap = theme.Space4
	} else if joined {
		gap = 0
	}
	var e ui.Element
	if vertical {
		e = ui.Column(c).Gap(gap).MinWidth(0)
	} else {
		e = ui.Row(c).Gap(gap).MinWidth(0)
	}
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	radius := theme.RadiusScale().LG
	rtl := c.Direction() == ui.RightToLeft
	n := len(buttons)
	e.Children(func() {
		for i := range buttons {
			b := buttons[i]
			adj := joinAdjust{grouped: true, vertical: vertical}
			if joined {
				adj.corners, adj.hasCorners = joinCorners(vertical, i, n, radius, rtl), true
				adj.seam = o.Separators && i > 0
				adj.borders, adj.hasBorders = joinBorders(vertical, i, n, rtl), true
			} else if b.Corners != nil {
				adj.corners, adj.hasCorners = *b.Corners, true
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
