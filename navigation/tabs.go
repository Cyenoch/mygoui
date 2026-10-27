package navigation

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TabsVariant chooses a segmented or underlined tab list.
type TabsVariant uint8

const (
	// TabsPrimary uses the selected pill on a default surface.
	TabsPrimary TabsVariant = iota
	// TabsSecondary uses a foreground underline under a flat list.
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
// choice internally. An invalid or disabled selection falls back to the
// nearest preceding enabled tab, or the first enabled tab, without OnChange.
// If there are no enabled items, Selected becomes -1 until one is available.
// NoListContainer drops the pill/scroll container and renders a bare tab
// list. MaxListHeight bounds a vertical list so it scrolls.
type TabsOptions struct {
	Selected        *int
	Items           []TabOptions
	Variant         TabsVariant
	Orientation     mygoui.Orientation
	Align           TabsAlign
	NoListContainer bool
	// KeepMounted builds inactive panels with Hidden, preserving their native
	// editor/scroll state without taking layout space or keyboard focus.
	KeepMounted   bool
	MaxListHeight float32
	Label         string
	OnChange      func(index int)
}

// tabsIndicator remembers what the indicator was placed from and where: the
// selection the previous build showed (seen), and the selection whose box the
// animated values were placed for (sel) with that box. Before the first
// placement both selections are -1.
type tabsIndicator struct {
	seen int
	sel  int
	box  ui.Rect
}

// tabsChannel tells apart the animated values of the indicator.
type tabsChannel uint8

type tabsPanelKey int

const (
	tabsX tabsChannel = iota
	tabsY
	tabsW
	tabsH
)

// Tabs builds a tab list and the panel of the selected tab. Click chooses a
// tab, as do Left/Right (or Up/Down) and Home/End while a tab has the focus,
// which follows the choice.
//
// The selection indicator is painted under the tabs by the list itself and
// slides between the committed tab boxes over 200 ms; a choice on the way retargets
// the slide from where it is, and reduced motion places it at once. A
// horizontal list scrolls when it overflows; a vertical one scrolls up to
// MaxListHeight.
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
	selected := o.Selected
	if selected == nil {
		selected = ui.Local(root, "tabs-selected", func() int { return 0 })
	}
	if len(o.Items) == 0 {
		*selected = -1
		return root
	}
	// viewportWidth is the space the horizontal tab strip has, measured
	// after layout.
	viewportWidth := ui.Local(root, "tabs-viewport", func() float32 { return 0 })
	// boxes holds the tabs' committed boxes, which place the indicator.
	boxes := ui.Local(root, "tabs-boxes", func() []ui.Rect { return make([]ui.Rect, len(o.Items)) })
	indicator := ui.Local(root, "tabs-indicator", func() tabsIndicator {
		return tabsIndicator{seen: -1, sel: -1}
	})
	nativeSelected := ui.Local(root, "tabs-enabled-selected", func() int { return 0 })
	if len(*boxes) != len(o.Items) {
		*boxes = make([]ui.Rect, len(o.Items))
		*indicator = tabsIndicator{seen: -1, sel: -1}
	}
	radius := theme.RadiusScale().MD
	// One pill per build at most: its shadow layers are resolved here, not
	// per tab and not in the paint callback.
	var shadows []theme.Shadow
	if o.Variant == TabsPrimary {
		shadows = theme.SurfaceShadows(c.Theme().Dark)
	}
	face := segmentBackground(c, t)
	var parts ui.TabsParts
	root.Children(func() {
		enabled := make([]int, 0, len(o.Items))
		chosen := -1
		for i, item := range o.Items {
			if !item.Disabled {
				if i == *selected {
					chosen = len(enabled)
				}
				enabled = append(enabled, i)
			}
		}
		if len(enabled) == 0 {
			*selected = -1
			*nativeSelected = -1
		} else {
			if chosen < 0 {
				chosen = 0
				for i, index := range enabled {
					if index > *selected {
						break
					}
					chosen = i
				}
				*selected = enabled[chosen]
			}
			*nativeSelected = chosen
		}
		// placed is the selection this pass builds for. A click applies in
		// the pass that builds the tabs, so *selected can change halfway
		// through it: the indicator, its boxes and its recorded selection
		// must all come from one choice, or the next pass would jump from a
		// box recorded for another tab.
		placed := *selected
		// The indicator can be painted under the tabs once their boxes are
		// known: from the second frame on. The first frame keeps the pill on
		// the selected tab itself.
		target := ui.Rect{}
		known := placed >= 0 && placed < len(*boxes) && (*boxes)[placed].W > 0
		containerRadius := theme.RadiusScale().LG
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
			container.Background(t.MutedBackground).Radius(containerRadius).Clip().
				Padding(scale(c, 2))
		}

		// buildList creates the native tab list under whatever container or
		// scroller is current and fills it with styled tabs. minWidth keeps
		// a horizontal list at least as wide as its viewport, so grown tabs
		// fill it, while longer labels can still overflow and scroll.
		buildList := func(minWidth float32) {
			parts = ui.TabsBase(c, nativeSelected, len(enabled))
			list := parts.List.MinWidth(0)
			listBox := list.LayoutBounds()
			var paintedList ui.Rect
			if vertical {
				list.Column().Gap(scale(c, theme.Space4)).AlignItems(ui.Stretch)
			} else {
				list.FillWidth().Gap(scale(c, 2))
				if minWidth > 0 {
					list.MinWidth(minWidth)
				}
			}
			list.Children(func() {
				nextEnabled := 0
				for i := range o.Items {
					var tab ui.Element
					if o.Items[i].Disabled {
						tab = ui.ButtonBase(c).Role(ui.RoleTab)
					} else {
						tab = parts.Tab(nextEnabled)
						nextEnabled++
					}
					el := buildTab(c, tab, i, o.Items[i], selected, o.Variant, o.Orientation, o.Align, known)
					// LayoutBounds is in the coordinates a Draw callback
					// paints in, which is where the indicator is painted.
					switch b := el.LayoutBounds(); {
					case b.W > 0 && b.H > 0:
						b.X -= listBox.X
						b.Y -= listBox.Y
						(*boxes)[i] = b
					default:
						(*boxes)[i] = ui.Rect{}
					}
					box := (*boxes)[i]
					el.DrawOver(func(_ *ui.Painter, measured ui.Rect) {
						if paintedList.W <= 0 || paintedList.H <= 0 {
							return
						}
						measured.X -= paintedList.X
						measured.Y -= paintedList.Y
						if layoutChanged(box, measured) {
							c.Invalidate()
						}
					})
				}
			})
			var along, size float32
			if known {
				// The tabs recorded their boxes as they built, so use the
				// freshest committed one: a strip that was still being
				// measured must not paint the pill at its old width.
				if b := (*boxes)[placed]; b.W > 0 {
					target = b
				} else {
					known = false
				}
			}
			if known {
				// The first frame that places the indicator seeds it where it
				// was: the tab the previous build showed as selected. A choice
				// made before the tab boxes were committed (a click on the
				// first frames, or on the pass that applies it) then still
				// slides from the old tab instead of starting at the new one.
				first := indicator.sel < 0
				from := target
				if first && indicator.seen >= 0 && indicator.seen != placed &&
					indicator.seen < len(*boxes) && (*boxes)[indicator.seen].W > 0 {
					from = (*boxes)[indicator.seen]
				}
				if first {
					if vertical {
						list.Jump(tabsY, from.Y)
						list.Jump(tabsH, from.H)
					} else {
						list.Jump(tabsX, from.X)
						list.Jump(tabsW, from.W)
					}
				}
				// A new choice slides from where the indicator is; the same
				// choice keeps its slide going, or rests; a layout change
				// under the same choice (a scroll, a resize) is followed at
				// once, as the indicator belongs to the tab in upstream.
				animated := indicator.sel != placed || indicator.box == target
				follow := func(key any, to float32) float32 {
					if animated {
						return list.AnimateTo(key, to, slideMotion(tabsIndicatorMotion))
					}
					return list.Jump(key, to)
				}
				if vertical {
					along = follow(tabsY, target.Y)
					size = follow(tabsH, target.H)
				} else {
					along = follow(tabsX, target.X)
					size = follow(tabsW, target.W)
				}
				*indicator = tabsIndicator{seen: *selected, sel: placed, box: target}
			} else {
				indicator.seen = *selected
			}
			list.Draw(func(p *ui.Painter, r ui.Rect) {
				paintedList = r
				if !known {
					return
				}
				rect := target
				if vertical {
					rect.Y, rect.H = along, size
				} else {
					rect.X, rect.W = along, size
				}
				rect.X += r.X
				rect.Y += r.Y
				if o.Variant == TabsSecondary {
					const thickness = 2 // h-0.5 / w-0.5
					if vertical {
						p.Fill(ui.Rect{X: rect.X, Y: rect.Y, W: thickness, H: rect.H}, t.Foreground, 0)
					} else {
						p.Fill(ui.Rect{X: rect.X, Y: rect.Y + rect.H - thickness, W: rect.W, H: thickness}, t.Foreground, 0)
					}
					return
				}
				p.Fill(rect, face, radius)
				for _, sh := range shadows {
					p.Shadow(rect, radius, sh.X, sh.Y, sh.Blur, sh.Spread, sh.Color)
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
		if parts.List.Changed() && *nativeSelected >= 0 && *nativeSelected < len(enabled) {
			*selected = enabled[*nativeSelected]
			if o.OnChange != nil {
				o.OnChange(*selected)
			}
		}
		for i, item := range o.Items {
			if item.Content == nil || !o.KeepMounted && i != *selected {
				continue
			}
			panel := ui.Column(c.Key(tabsPanelKey(i))).Hidden(i != *selected).FillWidth().MinWidth(0).Shrink(0).
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
// no ambient scope or free compound constructor is needed. known says the
// list paints the selection indicator from the tabs' boxes; until then the
// selected tab paints its own.
func buildTab(c *ui.Context, tab ui.Element, index int, o TabOptions, chosen *int,
	variant TabsVariant, orientation mygoui.Orientation, align TabsAlign, known bool) ui.Element {
	t := theme.Colors(c)
	vertical := orientation == mygoui.Vertical
	selected := index == *chosen
	height := scale(c, 30)
	radius := theme.RadiusScale().MD
	if variant == TabsSecondary {
		radius = 0
	}
	el := tab.Row().AlignItems(ui.Center).Gap(scale(c, theme.Space6)).
		Height(height).PaddingX(scale(c, 10)).MinWidth(0).Shrink(0).
		Radius(radius).FocusRing(true).Label(o.Label).Cursor(ui.CursorPointer)
	if vertical {
		el.MinWidth(scale(c, theme.Space(20))).Justify(ui.Center)
	} else {
		el.Grow(1).MinWidth(0).Justify(ui.Center)
	}
	switch align {
	case TabsAlignStart:
		el.Justify(ui.Start)
	case TabsAlignEnd:
		el.Justify(ui.End)
	}
	// Snapshot hover: transitions can retain paint after the element expires.
	hovered := el.Hovered() && !o.Disabled
	fg := t.MutedForeground.Alpha(0.72)
	switch {
	case selected && variant == TabsSecondary:
		fg = t.Foreground
	case selected:
		fg = t.Foreground
	case hovered:
		fg = t.MutedForeground
	}
	el.TextColor(fg)
	// The pill the selected tab paints until the list knows the tab boxes;
	// from the second frame the list paints the sliding one. Its shadow
	// layers are resolved once here, not in the paint callback.
	face := segmentBackground(c, t)
	var pill []theme.Shadow
	if selected && !known && variant != TabsSecondary {
		pill = theme.SurfaceShadows(c.Theme().Dark)
	}
	el.Draw(func(p *ui.Painter, r ui.Rect) {
		if variant == TabsSecondary {
			if hovered {
				p.Fill(r, t.Accent, theme.RadiusScale().MD)
			}
			if !known && selected {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 2, W: r.W, H: 2}, t.Foreground, 0)
			}
			return
		}
		if selected {
			if !known {
				p.Fill(r, face, radius)
				for _, sh := range pill {
					p.Shadow(r, radius, sh.X, sh.Y, sh.Blur, sh.Spread, sh.Color)
				}
			}
			return
		}

	})
	el.Children(func() {
		style.Text(ui.Text(c, o.Label), c, theme.SM()).FontWeight(theme.WeightMedium).SingleLine()
	})
	return style.Disabled(el, o.Disabled)
}

func segmentBackground(c *ui.Context, t theme.Palette) ui.Color {
	if c.Theme().Dark {
		return t.Input
	}
	return t.Background
}
