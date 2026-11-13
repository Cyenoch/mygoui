package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SidebarItemOptions describes one navigation action. ID must be nonempty
// and unique within the sidebar; Label is its accessible name. OnClick runs
// after the selection is updated.
type SidebarItemOptions struct {
	ID       string
	Label    string
	Icon     *ui.SVG
	Disabled bool
	OnClick  func()
}

// SidebarGroupOptions groups navigation actions under an optional label.
// Key preserves the group's identity when groups are reordered and must be
// unique among groups when set. Without a key, its position is its identity.
type SidebarGroupOptions struct {
	Key   string
	Label string
	Items []SidebarItemOptions
}

// SidebarOptions configures a sidebar. Selected and Collapsed are required
// caller-owned per-window state, and an empty Selected means no selection.
// Width and CollapsedWidth scale with the text size; nonpositive values use
// coss ui's 256 and 48 DIP defaults. Header and Footer build their own
// content in both states and can read Collapsed to adapt their layout. A
// collapsed item without an icon retains its text so it remains usable.
//
// OnSelect runs for every enabled activation, including the current item,
// before the item's OnClick. ToggleLabel can replace the English accessible
// name of the toggle.
type SidebarOptions struct {
	Selected         *string
	Collapsed        *bool
	Groups           []SidebarGroupOptions
	Header           func()
	Footer           func()
	Width            float32
	CollapsedWidth   float32
	Label            string
	ToggleLabel      string
	OnSelect         func(id string)
	OnCollapseChange func(collapsed bool)
}

// sidebarScrollKey retains the content scroll position that drives its edge fade.
type (
	sidebarScrollKey   struct{}
	sidebarCollapseKey struct{}
)

type sidebarGroupKey struct {
	name  string
	index int
}

// Sidebar builds a scrollable navigation sidebar with a native collapse
// toggle and individually focusable navigation actions. The caller owns the
// selected ID, collapse state and the work performed by an activation.
func Sidebar(c *ui.Context, o SidebarOptions) ui.Element {
	if o.Selected == nil {
		panic("mygoui.Sidebar: nil selected state")
	}
	if o.Collapsed == nil {
		panic("mygoui.Sidebar: nil collapsed state")
	}
	width, collapsedWidth := o.Width, o.CollapsedWidth
	if width <= 0 {
		width = 256
	}
	if collapsedWidth <= 0 {
		collapsedWidth = 48
	}
	collapsed := *o.Collapsed
	label := o.Label
	if label == "" {
		label = "Sidebar"
	}
	t := theme.Colors(c)
	root := ui.Column(c).FillHeight().Shrink(0).
		MinWidth(0).Background(t.Sidebar).TextColor(t.SidebarForeground).
		BorderEndWidth(theme.Hairline).BorderColor(t.SidebarBorder).
		Role(ui.RoleGroup).Label(label)
	// Collapsing moves between the widths over 200 ms, as coss does; a text
	// scale change resizes at once. Labels shown on expand are clipped.
	var progress float32
	if collapsed {
		progress = 1
	}
	// AnimateTo settles at once under reduced motion and keeps its value in
	// step, so turning motion back on does not replay an old change.
	progress = root.AnimateTo(sidebarCollapseKey{}, progress, ui.Animation{
		Duration: 200 * time.Millisecond, Ease: ui.Linear,
	})
	full, narrow := theme.Scale(c, width), theme.Scale(c, collapsedWidth)
	root.Width(full + (narrow-full)*progress)
	if progress > 0 && progress < 1 {
		root.Clip()
	}
	root.Children(func() {
		ui.Column(c.Key("sidebar-header")).Padding(theme.Space8).Gap(theme.Space8).Shrink(0).
			MinWidth(0).Children(func() {
			sidebarToggle(c, o, t, collapsed)
			if o.Header != nil {
				o.Header()
			}
		})
		scroll := ui.Local(root, sidebarScrollKey{}, func() ui.ScrollState { return ui.ScrollState{} })
		content := ui.Scroll(c.Key("sidebar-content")).Column().FillWidth().Grow(1).Basis(0).
			MinWidth(0).MinHeight(0).Gap(theme.Space8).TrackScroll(scroll)
		content.Children(func() {
			for i, group := range o.Groups {
				key := sidebarGroupKey{index: i}
				if group.Key != "" {
					key = sidebarGroupKey{name: group.Key}
				}
				section := ui.Column(c.Key(key)).FillWidth().MinWidth(0).Shrink(0).
					Padding(theme.Space8).Gap(theme.Space4)
				if group.Label != "" {
					section.Role(ui.RoleGroup).Label(group.Label)
				}
				section.Children(func() {
					if group.Label != "" && !collapsed {
						ui.Row(c).Height(theme.Scale(c, 32)).PaddingX(theme.Space8).
							AlignItems(ui.Center).Children(func() {
							style.Text(ui.Text(c, group.Label), c, theme.XS()).
								FontWeight(theme.WeightMedium).SingleLine()
						})
					}
					for _, item := range group.Items {
						sidebarItem(c, o, item, t, collapsed)
					}
				})
			}
		})
		ScrollFade(c, content, scroll, ScrollFadeOptions{
			Mode: SurfaceFade, Edges: FadeTop | FadeBottom, Surface: t.Sidebar,
		})
		if o.Footer != nil {
			ui.Column(c.Key("sidebar-footer")).Padding(theme.Space8).Gap(theme.Space8).
				MinWidth(0).Shrink(0).Children(o.Footer)
		}
	})
	return root
}

func sidebarToggle(c *ui.Context, o SidebarOptions, t theme.Palette, collapsed bool) {
	label := o.ToggleLabel
	if label == "" {
		label = "Toggle Sidebar"
	}
	size := theme.Scale(c, 32)
	iconSize := theme.Scale(c, theme.IconMD)
	e := ui.ButtonBase(c.Key("sidebar-toggle")).Size(size, size).AlignSelf(ui.Start).
		Role(ui.RoleToggleButton).Checked(!collapsed).Label(label).Tooltip(label).
		TextColor(t.SidebarForeground).Cursor(ui.CursorPointer)
	e = sidebarFace(e, t, false)
	e.Children(func() {
		ui.Icon(c, lucide.PanelLeft()).Size(iconSize, iconSize)
	})
	e.OnClick(func() {
		*o.Collapsed = !*o.Collapsed
		if o.OnCollapseChange != nil {
			o.OnCollapseChange(*o.Collapsed)
		}
	})
}

func sidebarItem(
	c *ui.Context,
	o SidebarOptions,
	item SidebarItemOptions,
	t theme.Palette,
	collapsed bool,
) {
	selected := *o.Selected == item.ID
	e := ui.ButtonBase(c.Key(item.ID)).FillWidth().MinWidth(0).
		MinHeight(theme.Scale(c, 32)).Padding(theme.Space8).Gap(theme.Space8).
		AlignItems(ui.Center).Justify(ui.Start).Label(item.Label).
		Role(ui.RoleButton).Checked(selected).Disabled(item.Disabled)
	if !item.Disabled {
		e.Cursor(ui.CursorPointer)
	}
	if selected {
		e.Description(locale.M().Selected)
	}
	// The ink fades with the selection; the pointer changes it at once.
	ink := motion.Color(e, sidebarItemInkKey{}, selectedInk(t, selected))
	if e.Hovered() || e.Pressed() {
		ink = t.SidebarAccentForeground
	}
	e.TextColor(ink)
	e = sidebarFace(e, t, selected)
	if collapsed {
		e.Tooltip(item.Label)
	}
	e.Children(func() {
		if item.Icon != nil {
			size := theme.Scale(c, theme.IconMD)
			ui.Icon(c, item.Icon).Size(size, size).Shrink(0)
		}
		if !collapsed || item.Icon == nil {
			text := style.Text(ui.Text(c, item.Label), c, theme.SM()).
				MinWidth(0).Grow(1).SingleLine()
			if selected {
				text.FontWeight(theme.WeightMedium)
			}
		}
	})
	if !item.Disabled {
		e.OnClick(func() {
			*o.Selected = item.ID
			if o.OnSelect != nil {
				o.OnSelect(item.ID)
			}
			if item.OnClick != nil {
				item.OnClick()
			}
		})
	}
}

type (
	sidebarItemInkKey  struct{}
	sidebarItemFillKey struct{}
)

func selectedInk(t theme.Palette, selected bool) ui.Color {
	if selected {
		return t.SidebarAccentForeground
	}
	return t.SidebarForeground
}

func sidebarFace(e ui.Element, t theme.Palette, selected bool) ui.Element {
	radius := theme.RadiusScale().MD
	fill := ui.Transparent
	if selected {
		fill = t.SidebarAccent
	}
	// The fill fades with the selection; the pointer changes it at once.
	fill = motion.Color(e, sidebarItemFillKey{}, fill)
	e = style.Background(e, radius, fill, t.SidebarAccent, t.SidebarAccent).FocusRing(false)
	if e.FocusVisible() {
		ring := t.SidebarRing
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Stroke(r, ring, radius, 2)
		})
	}
	return e
}
