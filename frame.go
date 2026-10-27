package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// FrameOptions configures the inset between a frame and its independent panels.
type FrameOptions struct {
	Padding, Gap float32
	Label        string
}

// Frame builds coss's muted outer frame. Panels remain independently composable.
func Frame(c *ui.Context, o FrameOptions, body func()) ui.Element {
	padding, gap := o.Padding, o.Gap
	if padding == 0 {
		padding = theme.Space4
	}
	if gap == 0 {
		gap = theme.Space4
	}
	t := theme.Colors(c)
	return ui.Column(c).FillWidth().MinWidth(0).Padding(padding).Gap(gap).
		Radius(theme.RadiusScale().XXL).Background(t.MutedBackground.Alpha(0.72)).
		Label(o.Label).Children(body)
}

// FramePanel builds a bordered panel; callers can override native layout on return.
func FramePanel(c *ui.Context, body func()) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Space20).
		Radius(theme.RadiusScale().XL).Border(theme.Hairline, t.Border).
		Background(t.Background).TextColor(t.Foreground)
	return style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark)).Children(body)
}

// FrameHeader groups a frame's title and description.
func FrameHeader(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Space16, theme.Space20).Children(body)
}

// FrameTitle renders a compact semibold frame title.
func FrameTitle(c *ui.Context, title string) ui.Element {
	ts := theme.SM()
	ts.Weight = theme.WeightSemibold
	return style.Text(ui.Text(c, title), c, ts).TextColor(theme.Colors(c).Foreground)
}

// FrameDescription renders the frame's secondary text.
func FrameDescription(c *ui.Context, text string) ui.Element {
	return style.Text(ui.Text(c, text), c, theme.SM()).TextColor(theme.Colors(c).MutedForeground)
}

// FrameFooter hosts frame-level actions outside its panels.
func FrameFooter(c *ui.Context, body func()) ui.Element {
	return ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Space8).
		Padding(theme.Space16, theme.Space20).Children(body)
}

// GroupOptions configures joined controls. Each item receives its own corners
// before construction, so its paint and keyboard focus use the same geometry.
type GroupOptions struct {
	Orientation Orientation
	FullWidth   bool
	Label       string
}

// GroupItem builds one control with the supplied outer-only corner radii.
// A mixed group can combine buttons, editors and application-defined controls.
type GroupItem struct {
	Key   string
	Build func(corners [4]float32) ui.Element
}

// Group composes native controls without an ambient provider or child inspection.
func Group(c *ui.Context, o GroupOptions, items ...GroupItem) ui.Element {
	vertical := o.Orientation == Vertical
	e := ui.Row(c).Gap(0).MinWidth(0).Role(ui.RoleGroup).Label(o.Label)
	if vertical {
		e.Column()
	}
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.AlignSelf(ui.Start)
	}
	return e.Children(func() {
		for i, item := range items {
			if item.Build == nil {
				panic("mygoui.Group: item has no Build callback")
			}
			corners := joinCorners(vertical, i, len(items), theme.RadiusScale().LG, c.Direction() == ui.RightToLeft)
			key := any(i)
			if item.Key != "" {
				key = item.Key
			}
			wrapper := ui.Box(c.Key(key)).MinWidth(0).Shrink(0)
			if vertical {
				wrapper.AlignSelf(ui.Stretch)
			}
			if o.FullWidth {
				if vertical {
					wrapper.FillWidth()
				} else {
					wrapper.Grow(1).Basis(0)
				}
			}
			wrapper.Children(func() {
				control := item.Build(corners)
				if vertical {
					control.FillWidth()
				}
				if o.FullWidth {
					control.FillWidth()
				}
				if len(items) > 1 {
					if vertical {
						borders := joinBorders(true, i, len(items), false)
						control.BorderWidth(borders[0], borders[1], borders[2], borders[3])
					} else {
						if i > 0 {
							control.BorderStartWidth(0)
						}
						if i < len(items)-1 {
							control.BorderEndWidth(0)
						}
					}
				}
			})
		}
	})
}
