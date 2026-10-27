package navigation

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ToolbarOptions configures a toolbar. Items builds its controls; use
// ToolbarSeparator between groups. Attached gives the toolbar a bordered card surface.
type ToolbarOptions struct {
	Orientation mygoui.Orientation
	Attached    bool
	Label       string
	Items       func()
}

// Toolbar groups controls into one focus stop whose arrows move between
// them: Left and Right in a horizontal toolbar, Up and Down in a vertical
// one, with Home and End at the ends. Assistive technology sees a toolbar
// named by Label.
func Toolbar(c *ui.Context, o ToolbarOptions) ui.Element {
	t := theme.Colors(c)
	vertical := o.Orientation == mygoui.Vertical
	root := ui.Column(c).Shrink(0).MinWidth(0).Gap(scale(c, theme.Space8)).
		Role(ui.RoleToolbar).FocusGroup(ui.Horizontal)
	if vertical {
		root.Column().AlignItems(ui.Start).Justify(ui.Start)
		root.FocusGroup(ui.Vertical)
	} else {
		root.Row().AlignItems(ui.Center)
	}
	if o.Attached {
		radius := theme.RadiusScale().XL
		root.Background(t.Card).TextColor(t.CardForeground).Border(theme.Hairline, t.Border).Radius(radius).Clip().
			Padding(scale(c, theme.Space4))
		// shadow-overlay is the whole rim: in dark mode it is an inset
		// hairline, so no extra border is added.
		root = style.Shadow(root, c, theme.OverlayShadows(c.Theme().Dark))
	}
	if o.Label != "" {
		root.Label(o.Label)
	}
	root.Children(o.Items)
	return root
}

// ToolbarSeparator draws a rule between groups of toolbar controls:
// a short vertical rule in a horizontal toolbar and a short horizontal rule
// in a vertical one. Orientation is explicit because Go composition has no
// ambient toolbar context.
func ToolbarSeparator(c *ui.Context, orientation mygoui.Orientation) ui.Element {
	t := theme.Colors(c)
	vertical := orientation == mygoui.Vertical
	e := ui.Box(c).Shrink(0).Background(t.Separator).Role(ui.RoleNone)
	if vertical {
		e.Width(scale(c, theme.Space20)).Height(hairline).AlignSelf(ui.Center)
	} else {
		e.Width(hairline).Height(scale(c, theme.Space20)).AlignSelf(ui.Center)
	}
	return e
}
