package gallery

import (
	"fmt"
	"math"
	"runtime"
	"slices"
	"sort"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TitleBarHeight is the gallery's native title-bar height in DIPs.
const TitleBarHeight = 52

// RepositoryURL is the gallery's source and release home.
const RepositoryURL = "https://github.com/Cyenoch/mygoui"

// View builds the workbench. Theme application happens here so appearance
// switches affect every specimen immediately.
func View(c *ui.Context, s *State) {
	zoom(c, s)
	applyAppearance(c, s)
	overlaysStateOf(s).toasts.View(c)
	mode := shellMode(s)
	sheet := layout.Compact(c, mode)
	if !s.NavOpen {
		s.navFocusSet = false
	}
	c.OnShortcut(ui.Cmd, ui.KeyK, func() {
		if sheet {
			s.NavOpen, s.navFocusSet = true, true
		} else {
			s.SidebarClosed = false
		}
		s.SearchFocus.Focus(c)
	})
	layout.Scaffold(c, layout.ScaffoldOptions{
		Compact: mode,
		StartBar: layout.TitleBarOptions{
			Label: "Sidebar title bar", MinHeight: TitleBarHeight,
			Center: func() {
				ui.Row(c).MinWidth(0).PaddingStart(theme.Scale(c, theme.Space4)).Children(func() { brand(c) })
			},
		},
		Start: &layout.PaneOptions{
			Label: "Gallery sidebar", Width: docsSidebarWidth, Fixed: true,
			Collapsed: &s.SidebarClosed, Open: &s.NavOpen,
			Body: func() { sidebar(c, s) },
		},
		TitleBar: pageHeader(c, s, mode),
		Content:  func() { content(c, s) },
	})
}

// shellMode is how the shell shows its sidebar. The landing page has none,
// so it navigates through the sheet at every width, as compact windows do.
func shellMode(s *State) layout.CompactMode {
	if isLanding(s) {
		return layout.CompactAlways
	}
	return layout.CompactAuto
}

// brand is the gallery mark; it truncates rather than reach the window
// controls.
func brand(c *ui.Context) {
	ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
		mygoui.Typography(c, "MyGo UI", mygoui.TypographyOptions{Type: theme.TypeH6, Truncate: true}).
			Shrink(1).MinWidth(0)
		mygoui.Typography(c, "MyGo", mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true, Truncate: true}).
			Shrink(2).MinWidth(0)
	})
}

// sidebar is the navigation pane under its title band: search over the
// scrolling page list. The compact sheet shows the same pane.
func sidebar(c *ui.Context, s *State) {
	t := theme.Colors(c)
	pad := theme.Scale(c, theme.Space12)
	ui.Column(c).FillWidth().Padding(0, pad, theme.Space8, pad).Children(func() { search(c, s) })
	mygoui.ScrollArea(c, &s.Sidebar, mygoui.ScrollAreaOptions{
		Axis:  mygoui.ScrollVertical,
		Label: "Gallery navigation",
		Fade: &mygoui.ScrollFadeOptions{
			Mode: mygoui.SurfaceFade, Edges: mygoui.FadeTop | mygoui.FadeBottom, Surface: t.Sidebar,
		},
	}, func() {
		ui.Column(c).FillWidth().Padding(theme.Space4, pad, theme.Space16, pad).Children(func() {
			pageList(c, s, &s.Sidebar, func(page Page) { selectPage(s, page) })
		})
	}).Grow(1).Basis(0).MinHeight(0)
}

func buttonFor(on bool) mygoui.ButtonVariant {
	if on {
		return mygoui.ButtonSecondary
	}
	return mygoui.ButtonGhost
}

var primarySwatches = []ui.Color{
	ui.RGB(24, 24, 27),
	ui.Oklch(0.6204, 0.195, 253.83),
	ui.Oklch(0.62, 0.22, 300),
	ui.Oklch(0.65, 0.19, 25),
	ui.Oklch(0.62, 0.16, 150),
	ui.Oklch(0.72, 0.16, 80),
}

func applyAppearance(c *ui.Context, s *State) {
	accent := ui.Color{}
	if s.Primary > 0 && s.Primary < len(primarySwatches) {
		accent = primarySwatches[s.Primary]
	}
	family := fonts.Family
	if s.SystemFont || s.SystemFontOnly {
		family = ""
	}
	theme.Configure(theme.Preferences{
		FontFamily: family, HeadingFontFamily: family,
		Accent: accent, Radius: s.Radius,
	})
	theme.ApplyMode(c, s.Dark)
}

// pageGroup is one nav block: a Header and the pages it names.
type pageGroup struct {
	name  string
	pages []Page
}

// has reports whether the group lists the page id.
func (g pageGroup) has(id string) bool {
	return id != "" && slices.ContainsFunc(g.pages, func(p Page) bool { return p.ID == id })
}

// groupRank orders the functional groups in the navigation. Pages inside a
// group keep the order they registered with.
var groupRank = map[string]int{
	"Overview": 0, "Foundations": 1, "Actions": 2, "Display": 3,
	"Forms": 4, "Collections": 5, "Navigation": 6, "Overlays": 7,
	"Dates": 8, "Color": 9, "Data": 10,
	// Composed components follow the primitives they are built from:
	// viewers of a document, then ready-made application patterns.
	"Viewers": 11, "Patterns": 12,
	"Charts": 13, "Chart Utilities": 14,
}

// navRankOf returns a group's navigation position; unknown groups sort last.
func navRankOf(group string) int {
	if rank, ok := groupRank[group]; ok {
		return rank
	}
	return 99
}

// pageGroups folds pages into one block per group, never repeating a heading
// when a search splits a group. The list is in navigation order, so a new
// block starts wherever the group changes.
func pageGroups(list []Page) []pageGroup {
	var groups []pageGroup
	for _, page := range list {
		if n := len(groups); n == 0 || groups[n-1].name != page.Group {
			groups = append(groups, pageGroup{name: page.Group})
		}
		groups[len(groups)-1].pages = append(groups[len(groups)-1].pages, page)
	}
	return groups
}

// pageList renders every matching page under its group heading, shared by
// the sidebar and the compact navigation sheet. Groups never collapse, so
// the list does not change height while browsing.
//
// With view set, only the rows near the viewport are built; spacers of the
// heights the rows last painted at stand in for the rest, so the content
// keeps its height and scroll range. The sheet passes nil: its scroll state
// belongs to the sheet, and it is on screen only while open.
func pageList(c *ui.Context, s *State, view *ui.ScrollState, choose func(Page)) {
	groups := s.sidebarGroups.of(s.Search, func() []pageGroup {
		return pageGroups(filterPagesInOrder(s.filterPages()))
	})
	band := navBandOf(c, s, view)
	// A row holding the focus is built wherever it has scrolled to: culling
	// it would drop the focus.
	focused := s.nav.focused
	s.nav.focused = ""
	// A sheet opens on the current page's row, not on the search field: on
	// a phone that would raise the keyboard. Focus is placed once per opening.
	landing := docsCompact(c) && !s.navFocusSet
	s.navFocusSet = true
	ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
		if len(groups) == 0 {
			ui.Row(c).Padding(theme.Space8).Children(func() {
				mygoui.Typography(c, "No components found", mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Muted: true,
				})
			})
		}
		top := float32(0)
		var rowMeasured, headingMeasured bool
		for _, group := range groups {
			height := s.nav.groupHeight(group)
			if !band.overlaps(top, top+height) && !group.has(focused) {
				ui.Box(c.Key(group.name)).FillWidth().Height(height)
				top += height + theme.Space16
				continue
			}
			ui.Column(c.Key(group.name)).FillWidth().Gap(theme.Space2).Children(func() {
				at, culled := top, float32(0)
				// item places one child of the group at the current offset,
				// folding the run of children above the band into a spacer.
				item := func(height float32, keep bool, build func()) {
					visible := keep || band.overlaps(at, at+height)
					if visible && culled > 0 {
						ui.Box(c).FillWidth().Height(culled - theme.Space2)
						culled = 0
					}
					if visible {
						build()
					} else {
						culled += height + theme.Space2
					}
					at += height + theme.Space2
				}
				if group.name != "Overview" {
					item(s.nav.heading, false, func() {
						groupHeading(c, group.name, measureOnce(&headingMeasured, &s.nav.heading))
					})
				}
				for _, page := range group.pages {
					current := landing && page.ID == s.Selected
					item(s.nav.row, page.ID == focused || current, func() {
						row, held := navRow(c.Key(page.ID), page.Title, page.ID == s.Selected,
							measureOnce(&rowMeasured, &s.nav.row), func() { choose(page) })
						if current {
							row.AutoFocus()
						}
						if held {
							s.nav.focused = page.ID
						}
					})
				}
				if culled > 0 {
					ui.Box(c).FillWidth().Height(culled - theme.Space2)
				}
			})
			top += height + theme.Space16
		}
	})
}

// navMetrics are the heights the sidebar's rows last painted at, before
// scaling and the font are known, estimated until the first paint, and the
// page whose row held the focus in the last build.
type navMetrics struct {
	row, heading float32
	focused      string
}

// groupHeight is the height of a group's column: its heading, rows and gaps.
func (m navMetrics) groupHeight(g pageGroup) float32 {
	h := float32(len(g.pages))*m.row + float32(len(g.pages)-1)*theme.Space2
	if g.name != "Overview" {
		h += m.heading + theme.Space2
	}
	return h
}

// navBand is the span of list offsets worth building.
type navBand struct{ from, to float32 }

// navBandOf spans the scrolled viewport, bounded by the window height,
// plus rows of overscan so Tab always finds a neighbor to move to. A nil
// view spans everything.
func navBandOf(c *ui.Context, s *State, view *ui.ScrollState) navBand {
	if s.nav.row == 0 {
		s.nav.row, s.nav.heading = theme.Scale(c, theme.HeightButtonSM), theme.Scale(c, theme.Space24)
	}
	if view == nil {
		return navBand{from: -math.MaxFloat32, to: math.MaxFloat32}
	}
	_, window := c.Size()
	overscan := 3 * s.nav.row
	y := min(view.Y, view.MaxY)
	return navBand{from: y - overscan, to: y + window + overscan}
}

func (b navBand) overlaps(from, to float32) bool { return to >= b.from && from <= b.to }

// measureOnce returns where the first element built in a pass records its
// painted height, nil for the rest.
func measureOnce(done *bool, into *float32) *float32 {
	if *done {
		return nil
	}
	*done = true
	return into
}

// filterPagesInOrder keeps the navigation order of a filtered page list.
func filterPagesInOrder(list []Page) []Page {
	sort.SliceStable(list, func(i, j int) bool { return list[i].navRank < list[j].navRank })
	return list
}

// groupHeading is the small muted section label above a nav group.
func groupHeading(c *ui.Context, text string, measure *float32) {
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	e := style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Muted).
		Padding(theme.Space4, theme.Space8, theme.Space4, theme.Space8).FillWidth().SingleLine().
		Role(ui.RoleHeading)
	if measure != nil {
		e.DrawOver(func(_ *ui.Painter, r ui.Rect) { *measure = r.H })
	}
}

// navRow is one nav link: a 32 DIP rounded row that fills with the sidebar
// accent while its page is the selected one.
// navRow returns the row and whether it holds the focus.
func navRow(c *ui.Context, title string, selected bool, measure *float32, onClick func()) (ui.Element, bool) {
	t := theme.Colors(c)
	fill, hover, ink := ui.Transparent, t.SidebarAccent, t.SidebarForeground
	if selected {
		fill, ink = t.SidebarAccent, t.SidebarAccentForeground
	}
	radius := theme.RadiusScale().LG
	ts := theme.SM()
	if selected {
		ts.Weight = theme.WeightMedium
	}
	e := ui.ButtonBase(c).FillWidth().MinHeight(theme.Scale(c, theme.HeightButtonSM)).
		PaddingX(theme.Scale(c, theme.Space8)).PaddingY(0).Radius(radius).
		Justify(ui.Start).AlignItems(ui.Center).FocusRing(false).Shrink(0).
		Role(ui.RoleButton).Checked(selected).Label(title).TextColor(ink)
	if selected {
		e.Description(locale.M().Selected)
	}
	e = style.Background(e, radius, fill, hover, hover)
	if focus := e.FocusVisible(); focus || measure != nil {
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if measure != nil {
				*measure = r.H
			}
			if focus {
				p.Stroke(r, t.SidebarRing, radius, 2)
			}
		})
	}
	e.Children(func() {
		style.Text(ui.Text(c, title), c, ts).SingleLine()
	})
	e.OnClick(onClick)
	return e, e.Focused()
}

// search filters the navigation. Its shortcut hint gives way to typing, so
// the field never changes size.
func search(c *ui.Context, s *State) {
	compact := docsCompact(c)
	t := theme.Colors(c)
	ui.Box(c).FillWidth().Children(func() {
		editor := form.Input(c.Key("gallery-search"), &s.Search, form.InputOptions{
			Placeholder: "Search components", AriaLabel: "Search components", FullWidth: true,
			OnChange: func() { s.Sidebar.Y = 0 },
			OnSubmit: func() {
				pages := filterPagesInOrder(s.filterPages())
				if len(pages) > 0 {
					selectPage(s, pages[0])
				}
			},
		}).OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEscape {
				return false
			}
			if k.Text() != "" {
				k.SetText("")
				s.Sidebar.Y = 0
				return true
			}
			if compact && s.NavOpen {
				s.NavOpen = false
				return true
			}
			return false
		}).Bind(&s.SearchFocus)
		if s.Search == "" && !editor.Focused() {
			ui.Row(c).Absolute().Top(0).Bottom(0).Right(theme.Scale(c, theme.Space8)).
				AlignItems(ui.Center).PassThrough().TextColor(t.MutedForeground).Children(func() {
				mygoui.Kbd(c, mygoui.KbdOptions{Keys: shortcutKeys()}).AlignSelf(ui.Center)
			})
		}
	})
}

// shortcutKeys spells the search shortcut for the host platform.
func shortcutKeys() []string {
	if runtime.GOOS == "darwin" {
		return []string{"⌘", "K"}
	}
	return []string{"Ctrl", "K"}
}

func selectPage(s *State, page Page) {
	s.Selected = page.ID
	s.Content.Y = 0
	s.NavOpen = false
}

func swatches(c *ui.Context, s *State) {
	ui.Row(c).Gap(theme.Space4).AlignItems(ui.Center).Children(func() {
		for i := range primarySwatches {
			value := primarySwatches[i]
			if i == 0 {
				value = theme.Colors(c).Foreground
			}
			mygoui.Button(c.Key(i), mygoui.ButtonOptions{
				Label: fmt.Sprintf("Primary %d", i+1), IconOnly: true, Size: theme.ButtonSmall,
				Variant: buttonFor(s.Primary == i),
				OnClick: func() { s.Primary = i },
				Content: func() {
					ui.Box(c).Size(theme.Scale(c, theme.Space12), theme.Scale(c, theme.Space12)).
						Radius(theme.Scale(c, theme.Space12)).Background(value)
				},
			}).Role(ui.RoleToggleButton).Checked(s.Primary == i)
		}
	})
}

func content(c *ui.Context, s *State) {
	t := theme.Colors(c)
	area := ui.Row(c).FillWidth().Grow(1).Basis(0).MinHeight(0)
	// Pages size themselves to the area; while the sidebar folds they follow
	// the room it has left, a frame late.
	laid := area.LayoutBounds().W
	s.areaLaid = laid
	area.DrawOver(func(_ *ui.Painter, r ui.Rect) {
		if d := laid - r.W; d < -0.5 || d > 0.5 {
			c.Invalidate()
		}
	})
	area.Children(func() {
		mygoui.ScrollArea(c.Key("workbench-content"), &s.Content, mygoui.ScrollAreaOptions{
			Axis: mygoui.ScrollVertical, Label: "Component content",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Edges: mygoui.FadeBottom, Surface: t.Background,
			},
		}, func() { docsPageColumn(c, s) }).Grow(1).Basis(0).MinWidth(0).FillHeight()
		if page, ok := indexedPage(c, s); ok {
			docsIndex(c, s, page)
		}
	})
}

// pageHeader is the title bar over the page: navigation toggle, location
// and window-level actions. The landing page shows the brand and the ways
// into the catalog in place of the sidebar it does not have.
func pageHeader(c *ui.Context, s *State, mode layout.CompactMode) layout.TitleBarOptions {
	compact := docsCompact(c)
	landing := isLanding(s) && !compact
	o := layout.TitleBarOptions{
		Label: "Gallery title bar", MinHeight: TitleBarHeight,
		Leading: func() {
			layout.PaneToggle(c, layout.PaneToggleOptions{
				Collapsed: &s.SidebarClosed, Open: &s.NavOpen, Compact: mode, Label: "Toggle navigation",
			})
		},
		Center: func() { breadcrumb(c, s, compact) },
		Trailing: func() {
			if landing {
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Search components", Icon: lucide.Search(), IconOnly: true,
					Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
					OnClick: func() {
						s.NavOpen, s.navFocusSet = true, true
						s.SearchFocus.Focus(c)
					},
				})
			}
			services := c.Services()
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "GitHub", Icon: lucide.Github(), IconOnly: true,
				Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() { services.OpenURL(RepositoryURL) },
			})
			appearanceButton(c, s)
		},
	}
	if landing {
		o.Leading, o.Center = func() { landingLinks(c, s) }, nil
	}
	return o
}

// landingLinks is the landing page's leading title bar: the brand and the
// ways into the catalog.
func landingLinks(c *ui.Context, s *State) {
	ui.Row(c).Shrink(0).AlignItems(ui.Center).Gap(theme.Space4).Children(func() {
		ui.Row(c).Shrink(0).Padding(0, theme.Space12, 0, theme.Space4).Children(func() { brand(c) })
		for _, link := range [...]struct{ label, page string }{
			{"Components", "button"}, {"Theming", "theming"},
		} {
			mygoui.Button(c.Key(link.page), mygoui.ButtonOptions{
				Label: link.label, Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() {
					if page, ok := pageFor(link.page); ok {
						selectPage(s, page)
					}
				},
			}).Shrink(0)
		}
	})
}

// breadcrumb names where the page sits; compact widths keep the title only.
func breadcrumb(c *ui.Context, s *State, compact bool) {
	t := theme.Colors(c)
	page, ok := currentPage(s)
	ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space4).AlignItems(ui.Center).Children(func() {
		if !ok || page.ID == "overview" {
			return
		}
		ts := theme.SM()
		if !compact && page.Group != "" {
			style.Text(ui.Text(c, page.Group), c, ts).TextColor(t.MutedForeground).SingleLine().Shrink(0)
			icon := theme.Scale(c, theme.IconSM)
			ui.Icon(c, lucide.ChevronRight()).Size(icon, icon).Shrink(0).TextColor(t.MutedForeground)
		}
		ts.Weight = theme.WeightMedium
		style.Text(ui.Text(c, page.Title), c, ts).SingleLine().MinWidth(0)
	})
}

// currentPage returns the selected page, or the first registered page when
// the selection names none. It reads the catalog under its lock and by ID, so
// a build does not scan or alias the shared slice.
func currentPage(s *State) (Page, bool) {
	if page, ok := pageFor(s.Selected); ok {
		return page, true
	}
	pagesMu.RLock()
	defer pagesMu.RUnlock()
	if len(pages) == 0 {
		return Page{}, false
	}
	return pages[0], true
}
