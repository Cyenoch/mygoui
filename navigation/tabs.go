package navigation

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TabsVariant is HeroUI's tab-list treatment.
type TabsVariant uint8

const (
	// TabsPrimary uses the selected pill on a default surface.
	TabsPrimary TabsVariant = iota
	// TabsSecondary uses an accent underline under a flat list.
	TabsSecondary
)

// TabsAlign places a tab's label within its box.
type TabsAlign uint8

const (
	TabsAlignCenter TabsAlign = iota
	TabsAlignStart
	TabsAlignEnd
)

// TabOptions configures one tab: its label, its panel body and whether it can
// be chosen.
type TabOptions struct {
	Label    string
	Content  func()
	Disabled bool
}

// TabsOptions configures a tab set. Selected is caller-owned; nil keeps the
// choice internally. NoListContainer drops the pill/scroll container and
// renders a bare tab list. MaxListHeight bounds a vertical list so it scrolls.
type TabsOptions struct {
	Selected        *int
	Items           []TabOptions
	Variant         TabsVariant
	Orientation     mygoui.Orientation
	Align           TabsAlign
	NoListContainer bool
	MaxListHeight   float32
	Label           string
	OnChange        func(index int)
}

// Tabs builds a tab list and the panel of the selected tab. Click chooses a
// tab, as do Left/Right (or Up/Down) and Home/End while a tab has the focus,
// which follows the choice. A horizontal list scrolls when it overflows; a
// vertical one scrolls up to MaxListHeight.
func Tabs(c *ui.Context, o TabsOptions) ui.Element {
	t := theme.Colors(c)
	vertical := o.Orientation == mygoui.Vertical
	root := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).Gap(scale(c, theme.Space8))
	if vertical {
		root.Row().AlignItems(ui.Start)
	}
	if o.Label != "" {
		root.Label(o.Label)
	}
	if len(o.Items) == 0 {
		return root
	}
	selected := o.Selected
	if selected == nil {
		selected = ui.Local(root, "tabs-selected", func() int { return 0 })
	}
	// viewportWidth is the space the horizontal tab strip has, measured
	// after layout.
	viewportWidth := ui.Local(root, "tabs-viewport", func() float32 { return 0 })
	var parts ui.TabsParts
	root.Children(func() {
		containerRadius := theme.RadiusBase() * 2.5
		var container ui.Element
		if vertical {
			container = ui.Column(c).Shrink(0).MinWidth(0)
		} else {
			container = ui.Column(c).FillWidth().MinWidth(0).Shrink(0)
		}
		if o.NoListContainer {
			container.Background(ui.Transparent).Radius(0)
		} else if o.Variant == TabsSecondary {
			container.Background(ui.Transparent).Radius(0)
			if vertical {
				container.BorderWidth(0, 0, 0, hairline).BorderColor(t.Border)
			} else {
				container.BorderWidth(0, 0, hairline, 0).BorderColor(t.Border)
			}
		} else {
			container.Background(t.Default).Radius(containerRadius).Clip().
				Padding(scale(c, theme.Space4))
		}

		// buildList creates the native tab list under whatever container or
		// scroller is current and fills it with styled tabs. minWidth keeps
		// a horizontal list at least as wide as its viewport, so grown tabs
		// fill it, while longer labels can still overflow and scroll.
		buildList := func(minWidth float32) {
			parts = ui.TabsBase(c, selected, len(o.Items))
			list := parts.List.MinWidth(0)
			if vertical {
				list.Column().Gap(scale(c, theme.Space4)).AlignItems(ui.Stretch)
			} else {
				list.FillWidth().Gap(0)
				if minWidth > 0 {
					list.MinWidth(minWidth)
				}
			}
			list.Children(func() {
				for i := range o.Items {
					buildTab(c, parts, i, o.Items[i], selected, o.Variant, o.Orientation, o.Align)
				}
			})
		}
		if !o.NoListContainer && vertical && o.MaxListHeight > 0 {
			container.Children(func() {
				ui.Scroll(c).MaxHeight(scale(c, o.MaxListHeight)).MinWidth(0).Shrink(0).
					Children(func() { buildList(0) })
			})
		} else if !vertical {
			container.Children(func() {
				viewport := ui.ScrollHorizontal(c).FillWidth().MinWidth(0).Shrink(0)
				// A horizontal list cannot take a percentage of a scroller's
				// unbounded width, so record the viewport after layout and
				// rebuild when it changes: the list then floors at that
				// width and its tabs divide it, while longer labels still
				// overflow and scroll.
				viewport.DrawOver(func(_ *ui.Painter, r ui.Rect) {
					if d := r.W - *viewportWidth; d > 0.5 || d < -0.5 {
						*viewportWidth = r.W
						c.Invalidate()
					}
				})
				viewport.Children(func() { buildList(*viewportWidth) })
			})
		} else {
			container.Children(func() { buildList(0) })
		}
		if parts.List.Changed() && o.OnChange != nil {
			o.OnChange(*selected)
		}
		if *selected >= 0 && *selected < len(o.Items) {
			item := o.Items[*selected]
			panel := ui.Column(c).FillWidth().MinWidth(0).Shrink(0).
				Padding(scale(c, theme.Space8)).Role(ui.RoleGroup)
			if item.Label != "" {
				panel.Label(item.Label)
			}
			if vertical {
				panel.MarginStart(scale(c, theme.Space16))
			} else {
				panel.Margin(scale(c, theme.Space16), 0, 0, 0)
			}
			panel.Children(item.Content)
		}
	})
	return root
}

// buildTab renders one row of a native tab list. Tabs owns the list state, so
// no ambient scope or free compound constructor is needed.
func buildTab(c *ui.Context, parts ui.TabsParts, index int, o TabOptions, chosen *int,
	variant TabsVariant, orientation mygoui.Orientation, align TabsAlign) ui.Element {
	t := theme.Colors(c)
	vertical := orientation == mygoui.Vertical
	selected := index == *chosen
	height := scale(c, theme.Space32)
	radius := theme.RadiusScale().XXXL
	el := parts.Tab(index).Row().AlignItems(ui.Center).Gap(scale(c, theme.Space6)).
		Height(height).PaddingX(scale(c, theme.Space16)).MinWidth(0).Shrink(0).
		Radius(radius).FocusRing(true).Label(o.Label).Cursor(ui.CursorPointer)
	if vertical {
		el.MinWidth(scale(c, theme.Space64)).Justify(ui.Center)
	} else {
		el.Grow(1).MinWidth(0).Justify(ui.Center)
	}
	switch align {
	case TabsAlignStart:
		el.Justify(ui.Start)
	case TabsAlignEnd:
		el.Justify(ui.End)
	}
	fg := t.Muted
	if selected {
		if variant == TabsSecondary {
			fg = t.Foreground
		} else {
			fg = t.SegmentForeground
		}
	}
	el.TextColor(fg)
	prevSelected := index > 0 && index-1 == *chosen
	shadows := theme.SurfaceShadows(c.Theme().Dark)
	// Snapshot the tab's hover in the build: a retained copy of a leaving
	// element paints in later passes, where el has expired.
	hovered := el.Hovered() && !o.Disabled
	el.Draw(func(p *ui.Painter, r ui.Rect) {
		if selected {
			if variant == TabsSecondary {
				thickness := float32(2)
				if vertical {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: thickness, H: r.H}, t.Accent, 0)
				} else {
					p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - thickness, W: r.W, H: thickness}, t.Accent, 0)
				}
			} else {
				p.Fill(r, t.Segment, radius)
				for _, sh := range shadows {
					p.Shadow(r, radius, sh.X, sh.Y, sh.Blur, sh.Spread, sh.Color)
				}
			}
			return
		}
		if variant == TabsSecondary {
			return
		}
		if hovered {
			p.Fill(r, t.Foreground.Alpha(0.04), radius)
		}
		if index > 0 && !prevSelected {
			if vertical {
				p.Fill(ui.Rect{X: r.X + r.W*0.05, Y: r.Y, W: r.W * 0.9, H: hairline}, t.Muted.Alpha(0.25), 0)
			} else {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H*0.25, W: hairline, H: r.H * 0.5}, t.Muted.Alpha(0.25), 0)
			}
		}
	})
	el.Children(func() {
		style.Text(ui.Text(c, o.Label), c, theme.SM()).FontWeight(theme.WeightMedium).SingleLine()
	})
	return style.Disabled(el, o.Disabled)
}
