package navigation

import (
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type SidebarID string

type SidebarItem struct {
	ID           SidebarID
	Label, Badge string
	Icon         *ui.SVG
	Disabled     bool
	OnSelect     func()
	Action       *mygoui.ButtonOptions
	// Open controls a branch disclosure. Nil retains its state in SidebarState.
	// Branches toggle their children; leaf callbacks own navigation.
	Open     *bool
	Children []SidebarItem
}

type SidebarGroup struct {
	ID     SidebarID
	Label  string
	Open   *bool
	Items  []SidebarItem
	Action *mygoui.ButtonOptions
}

type sidebarEntry struct {
	trigger, panel ui.Handle
	open, wasOpen  bool
}

type SidebarState struct {
	Collapsed, Open       bool
	Selected              SidebarID
	scroll, contentScroll ui.ScrollState
	entries               map[SidebarID]*sidebarEntry
	toggle                ui.Handle
	sheetWasOpen          bool
}

type SidebarOptions struct {
	Label           string
	Width, NarrowAt float32
	// Header and Footer remain outside the scrolling groups. Callers provide
	// their brand/team and account controls, including their compact rail form.
	Header, Footer func(collapsed bool)
	// ContentHeader follows the sidebar toggle in the main pane's fixed header.
	// Content builds the main pane's independently scrolling body.
	ContentHeader, Content func()
	OnSelect               func(SidebarID)
}

// Sidebar builds a bounded navigation and content shell. Its default 480-DIP
// height can be overridden with Height or FillHeight on the returned root.
// Desktop collapse retains a labelled icon rail; narrow windows use a native
// modal sheet with focus trapping, dismissal and opener focus restoration.
func Sidebar(c *ui.Context, s *SidebarState, o SidebarOptions, groups []SidebarGroup) ui.Element {
	p := theme.Colors(c)
	width := o.Width
	if width <= 0 {
		width = theme.SidebarWidth
	}
	threshold := o.NarrowAt
	if threshold <= 0 {
		threshold = 640
	}
	w, h := c.Size()
	narrow := w < threshold
	label := o.Label
	if label == "" {
		label = locales.ComponentSidebar()
	}
	// Modal construction may move focus before Overlay captures its opener.
	// Keep the external toggle as the explicit return target across all dismissals.
	if s.sheetWasOpen && !s.Open {
		s.toggle.Focus(c)
	}
	s.sheetWasOpen = s.Open
	root := ui.Row(c).FillWidth().Height(480).MaxHeight(h).MinWidth(0).MinHeight(0).
		AlignItems(ui.Stretch).Clip().Background(p.Background)
	root.Children(func() {
		if !narrow {
			railWidth := width
			if s.Collapsed {
				railWidth = 56
			}
			panel := ui.Column(c.Key("sidebar")).Width(railWidth).Shrink(0).FillHeight().MinHeight(0).
				Clip().Background(p.Sidebar).TextColor(p.SidebarForeground).
				BorderWidth(0, theme.HairlineWidth, 0, 0).BorderColor(p.SidebarBorder)
			if c.Direction() == ui.RightToLeft {
				panel.BorderWidth(0, 0, 0, theme.HairlineWidth)
			}
			if !theme.ReduceMotion(c) {
				panel.Transition(ui.ElementTransition{Size: true})
			}
			panel.Children(func() { sidebarBody(c, s, o, groups, s.Collapsed, false) })
		}
		ui.Column(c.Key("main")).Grow(1).Basis(0).MinWidth(0).MinHeight(0).Clip().Children(func() {
			ui.Row(c.Key("header")).Shrink(0).MinHeight(theme.HeightHeader).Padding(theme.Space8, theme.Space12).
				Gap(theme.Space12).Children(func() {
				trigger := mygoui.Button(c.Key("toggle"), mygoui.ButtonOptions{
					Label: locales.ComponentToggleNamed(label), Icon: lucide.PanelLeft(),
					Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconSM,
					OnClick: func() {
						if narrow {
							s.Open = !s.Open
						} else {
							s.Collapsed = !s.Collapsed
						}
					},
				})
				trigger.Bind(&s.toggle).Expanded((narrow && s.Open) || (!narrow && !s.Collapsed))
				if o.ContentHeader != nil {
					ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(o.ContentHeader)
				}
			})
			viewport := ui.Scroll(c.Key("content")).TrackScroll(&s.contentScroll).
				Grow(1).Basis(0).MinHeight(0).FillWidth()
			mygoui.ScrollEdges(c, viewport, &s.contentScroll, mygoui.ScrollEdgesOptions{Surface: p.Background})
			if o.Content != nil {
				viewport.Children(o.Content)
			}
		})
		if narrow {
			// DialogBase owns modality and native sheet dismissal. Building the
			// bounded body directly keeps its header/footer outside all scrolling.
			ui.DialogBase(c.Key("sheet"), &s.Open, func(back, panel ui.Element) {
				back.Justify(ui.Start).AlignItems(ui.Start).Background(p.Overlay)
				if !theme.ReduceMotion(c) {
					back.Transition(ui.ElementTransition{Colors: true, Enter: &ui.Motion{}, Exit: &ui.Motion{}})
				}
				panel.Column().Width(min(width, max(0, w-theme.Space32))).FillHeight().MinHeight(0).
					Background(p.Sidebar).TextColor(p.SidebarForeground).Label(label).Clip()
				ui.Row(c.Key("dismiss")).Shrink(0).Padding(theme.Space8).Justify(ui.End).Children(func() {
					// This named close control needs no tooltip overlay: a tooltip
					// would consume Escape before the sheet when keyboard-focused.
					close := ui.ButtonBase(c).Label(locales.ComponentCloseNamed(label)).Size(theme.Metrics(c).CompactHeight, theme.Metrics(c).CompactHeight).
						Radius(theme.Radii().XL).TextColor(p.SidebarForeground).OnClick(func() { s.Open = false })
					if close.Hovered() {
						close.Background(p.SidebarAccent)
					}
					close.Children(func() { ui.Icon(c, lucide.X()).Size(theme.IconMedium, theme.IconMedium) })
				})
				sidebarBody(c, s, o, groups, false, true)
			})
		}
	})
	return root
}

func sidebarBody(c *ui.Context, s *SidebarState, o SidebarOptions, groups []SidebarGroup, collapsed, narrow bool) {
	p := theme.Colors(c)
	if o.Header != nil {
		ui.Column(c.Key("brand")).Shrink(0).MinWidth(0).Padding(theme.Space8).
			Children(func() { o.Header(collapsed) })
	}
	viewport := ui.Scroll(c.Key("navigation")).TrackScroll(&s.scroll).
		Grow(1).Basis(0).MinHeight(0).FillWidth().Padding(theme.Space8)
	mygoui.ScrollEdges(c, viewport, &s.scroll, mygoui.ScrollEdgesOptions{Surface: p.Sidebar})
	viewport.Children(func() {
		ui.Column(c).Role(ui.RoleTree).Label(o.Label).FocusGroup(ui.Vertical).Gap(theme.Space16).Children(func() {
			for _, group := range groups {
				ui.Column(c.Key(group.ID)).Shrink(0).Gap(theme.Space4).Children(func() {
					if !collapsed && group.Label != "" {
						ui.Row(c).Gap(theme.Space4).Children(func() {
							if group.Open == nil {
								ui.Text(c, group.Label).Grow(1).Padding(theme.Space8).
									FontSize(theme.Scale(c, theme.Small().Size)).TextColor(p.MutedForeground)
							} else {
								mygoui.Button(c.Key("disclosure"), mygoui.ButtonOptions{
									Label: group.Label, TrailingIcon: lucide.ChevronDown(),
									Variant: mygoui.ButtonGhost, Size: mygoui.ButtonSM,
									OnClick: func() { *group.Open = !*group.Open },
								}).Expanded(*group.Open).Grow(1).Justify(ui.Start)
							}
							if group.Action != nil {
								mygoui.Button(c.Key("action"), *group.Action)
							}
						})
					}
					if collapsed || group.Open == nil || *group.Open {
						sidebarItems(c, s, o, group.Items, collapsed, narrow)
					}
				})
			}
		})
	})
	if o.Footer != nil {
		ui.Column(c.Key("account")).Shrink(0).MinWidth(0).Padding(theme.Space8).
			Children(func() { o.Footer(collapsed) })
	}
}

func sidebarItems(c *ui.Context, s *SidebarState, o SidebarOptions, items []SidebarItem, collapsed, narrow bool) {
	p := theme.Colors(c)
	for _, item := range items {
		ui.Column(c.Key(item.ID)).Shrink(0).Gap(theme.Space4).Children(func() {
			if s.entries == nil {
				s.entries = make(map[SidebarID]*sidebarEntry)
			}
			entry := s.entries[item.ID]
			if entry == nil {
				entry = &sidebarEntry{}
				s.entries[item.ID] = entry
			}
			open := &entry.open
			if item.Open != nil {
				open = item.Open
			}
			branch := len(item.Children) > 0
			if entry.wasOpen && (!*open || collapsed) && entry.panel.FocusWithin(c) {
				entry.trigger.Focus(c)
			}
			var parts ui.CollapsibleParts
			ui.Row(c).Gap(theme.Space4).Children(func() {
				var b ui.Element
				if branch && !collapsed {
					parts = ui.CollapsibleBase(c.Key("item"), open)
					b = parts.Trigger
				} else {
					b = ui.ButtonBase(c.Key("item"))
				}
				b.Bind(&entry.trigger).Selected(s.Selected == item.ID).Label(item.Label).Disabled(item.Disabled).
					Justify(ui.Start).AlignItems(ui.Center).MinWidth(0).MinHeight(theme.HeightInput).
					Padding(theme.Space6, theme.Space8).Gap(theme.Space8).Radius(theme.Radii().XL).
					Grow(1).Basis(0).FontSize(theme.Scale(c, theme.Body().Size)).TextColor(p.SidebarForeground)
				if !branch {
					b.Role(ui.RoleTreeItem)
				}
				if s.Selected == item.ID || b.Hovered() || b.Pressed() {
					b.Background(p.SidebarAccent).TextColor(p.SidebarAccentForeground)
				}
				if item.Disabled {
					b.TextColor(p.DisabledForeground)
				}
				if collapsed {
					b.Tooltip(item.Label).Center()
				}
				if !theme.ReduceMotion(c) {
					b.Transition(ui.ElementTransition{Colors: true})
				}
				b.Children(func() {
					if item.Icon != nil {
						ui.Icon(c, item.Icon).Size(theme.IconMedium, theme.IconMedium).Shrink(0)
					} else if collapsed {
						for _, r := range item.Label {
							ui.Text(c, string(r))
							break
						}
					}
					if !collapsed {
						ui.Text(c, item.Label).Grow(1).MinWidth(0).SingleLine()
						if item.Badge != "" {
							ui.Text(c, item.Badge).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(p.MutedForeground)
						}
						if branch {
							angle := 90 * parts.Progress()
							if c.Direction() == ui.RightToLeft {
								angle = 180 - angle
							}
							ui.Icon(c, lucide.ChevronRight()).Size(theme.IconMedium, theme.IconMedium).Shrink(0).Rotate(angle)
						}
					}
				})
				if !branch || collapsed {
					b.OnClick(func() {
						if branch {
							s.Collapsed = false
							*open = true
							return
						}
						s.Selected = item.ID
						if narrow {
							s.Open = false
						}
						if o.OnSelect != nil {
							o.OnSelect(item.ID)
						}
						if item.OnSelect != nil {
							item.OnSelect()
						}
					})
				}
				b.HandleInput(func(ev ui.InputEvent) bool {
					if ev.Kind != ui.InputKeyDown || ev.Mods != 0 || ev.Key < ui.KeyA || ev.Key > ui.KeyZ {
						return false
					}
					return focusSidebarMatch(c, s, items, string(rune('a'+ev.Key-ui.KeyA)), collapsed)
				})
				if !collapsed && item.Action != nil {
					mygoui.Button(c.Key("action"), *item.Action)
				}
			})
			if branch && !collapsed {
				parts.Panel(func() {
					guide := ui.Column(c).Bind(&entry.panel).Disabled(!*open || item.Disabled).
						Margin(theme.Space4, theme.Space12).PaddingStart(theme.Space12).Gap(theme.Space4).
						BorderWidth(0, 0, 0, theme.HairlineWidth).BorderColor(p.SidebarBorder)
					if c.Direction() == ui.RightToLeft {
						guide.BorderWidth(0, theme.HairlineWidth, 0, 0)
					}
					guide.Children(func() { sidebarItems(c, s, o, item.Children, false, narrow) })
				})
			}
			entry.wasOpen = *open && !collapsed
		})
	}
}

func focusSidebarMatch(c *ui.Context, s *SidebarState, items []SidebarItem, prefix string, collapsed bool) bool {
	for _, item := range items {
		entry := s.entries[item.ID]
		if item.Disabled || entry == nil {
			continue
		}
		if strings.HasPrefix(strings.ToLower(item.Label), prefix) {
			entry.trigger.Focus(c)
			return true
		}
		open := entry.open
		if item.Open != nil {
			open = *item.Open
		}
		if !collapsed && open && focusSidebarMatch(c, s, item.Children, prefix, false) {
			return true
		}
	}
	return false
}
