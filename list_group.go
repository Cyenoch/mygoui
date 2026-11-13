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

// ListGroup is a rounded surface containing related rows: a list named by
// Label, whose ListGroupItems are its items.
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
	radius := min(theme.Space(8), theme.RadiusScale().XXXL)
	e := ui.Column(c).FillWidth().MinWidth(0).Padding(0).Radius(radius).Clip().
		Background(background).TextColor(ink).Role(ui.RoleList).Label(o.Label)
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

// ListGroupItem renders a 16 DIP padded list item with an optional chevron,
// which points along the reading direction. With OnClick the row is a native
// button inside the item. Prefix and suffix callbacks own their layout.
func ListGroupItem(c *ui.Context, o ListGroupItemOptions) ui.Element {
	if o.OnClick == nil || o.Disabled {
		return listGroupRow(c, ui.Row(c).Role(ui.RoleListItem), o)
	}
	item := ui.Column(c).FillWidth().MinWidth(0).Role(ui.RoleListItem).Label(o.Title)
	return item.Children(func() { listGroupRow(c, ui.ButtonBase(c).Justify(ui.Start), o) })
}

func listGroupRow(c *ui.Context, e ui.Element, o ListGroupItemOptions) ui.Element {
	t := theme.Colors(c)
	e.FillWidth().MinWidth(0).AlignItems(ui.Center).Gap(theme.Space(3)).
		Padding(theme.Space(4)).Label(o.Title).Description(o.Description)
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
			chevron := lucide.ChevronRight()
			if c.Direction() == ui.RightToLeft {
				chevron = lucide.ChevronLeft()
			}
			ui.Icon(c, chevron).Size(theme.Scale(c, theme.IconMD),
				theme.Scale(c, theme.IconMD)).TextColor(t.Muted)
		}
	})
	return style.Disabled(e, o.Disabled)
}
