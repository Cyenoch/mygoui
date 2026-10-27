package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ListGroupOptions configures a native list group.
type ListGroupOptions struct {
	Variant SurfaceVariant
	Label   string
}

// ListGroup is a rounded surface containing related rows.
func ListGroup(c *ui.Context, o ListGroupOptions, body func()) ui.Element {
	t := theme.Colors(c)
	background, ink, shadowed := t.Surface, t.SurfaceForeground, true
	switch o.Variant {
	case SurfaceSecondary:
		background, ink = t.SurfaceSecondary, t.SurfaceSecondaryForeground
	case SurfaceTertiary:
		background, ink = t.SurfaceTertiary, t.SurfaceTertiaryForeground
	case SurfaceTransparent:
		background, ink, shadowed = ui.Transparent, t.Foreground, false
	}
	radius := min(theme.Space32, theme.RadiusScale().XXXL)
	e := ui.Column(c).FillWidth().MinWidth(0).Padding(0).Radius(radius).Clip().
		Background(background).TextColor(ink).Label(o.Label)
	if shadowed {
		e = style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark))
	}
	return e.Children(body)
}

// ListGroupItemOptions configures one row.
type ListGroupItemOptions struct {
	Title, Description string
	Prefix, Suffix     func()
	Chevron            bool
	Disabled           bool
	OnClick            func()
}

// ListGroupItem renders a 16 DIP padded row with an optional chevron and real
// native activation. Prefix and suffix callbacks own their layout.
func ListGroupItem(c *ui.Context, o ListGroupItemOptions) ui.Element {
	t := theme.Colors(c)
	var e ui.Element
	if o.OnClick != nil && !o.Disabled {
		e = ui.ButtonBase(c).Justify(ui.Start)
	} else {
		e = ui.Row(c)
	}
	e.FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Space12).
		Padding(theme.Space16).Label(o.Title).Description(o.Description)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick).FocusRing(true).Cursor(ui.CursorPointer)
		e = style.Background(e, 0, ui.Transparent, t.SurfaceHover, t.SurfaceHover)
	}
	e.Children(func() {
		if o.Prefix != nil {
			o.Prefix()
		}
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
			if o.Title != "" {
				// The title inherits the group's surface foreground.
				style.Text(ui.Text(c, o.Title), c, theme.Base()).
					FontWeight(theme.WeightMedium).SingleLine()
			}
			if o.Description != "" {
				style.Text(ui.Text(c, o.Description), c, theme.SM()).
					TextColor(t.Muted).SingleLine()
			}
		})
		if o.Suffix != nil {
			o.Suffix()
		}
		if o.Chevron {
			ui.Icon(c, lucide.ChevronRight()).Size(theme.Scale(c, theme.IconMD),
				theme.Scale(c, theme.IconMD)).TextColor(t.Muted)
		}
	})
	return style.Disabled(e, o.Disabled)
}
