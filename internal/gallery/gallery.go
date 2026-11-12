package gallery

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
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
	t := theme.Colors(c)
	compact := docsCompact(c)
	sidebarOpen := !compact && !s.SidebarClosed
	if !compact {
		// The sheet exists only in compact layouts; do not reopen it later.
		s.NavOpen = false
	}
	c.OnShortcut(ui.Cmd, ui.KeyK, func() {
		if compact {
			s.NavOpen = true
		} else {
			s.SidebarClosed = false
		}
		s.SearchFocus.Focus(c)
	})
	ui.Row(c).Fill().MinWidth(0).Clip().AlignItems(ui.Stretch).
		Background(t.Background).TextColor(t.Foreground).Children(func() {
		if sidebarOpen {
			sidebar(c, s)
		}
		content(c, s, compact, sidebarOpen)
	})
	if compact {
		navigationSheet(c, s)
	}
}

// headerControls is the number of buttons in the page title bar.
const headerControls = 3

// headerBand is the native title band the gallery covers.
func headerBand(c *ui.Context) float32 {
	return max(float32(TitleBarHeight)/c.Zoom(), c.TitleBar().Height)
}

// headerStacked reports whether the page title bar's controls cannot fit
// between the native window controls; they then move below the band.
func headerStacked(c *ui.Context, sidebarOpen bool) bool {
	width, _ := c.Size()
	if sidebarOpen {
		width -= theme.Scale(c, docsSidebarWidth)
	}
	bar := c.TitleBar()
	pad := theme.Scale(c, theme.Space12)
	need := headerControls*theme.Scale(c, theme.HeightButtonSM) + (headerControls-1)*theme.Space8
	if !sidebarOpen {
		width -= bar.Left
	}
	return need > width-bar.Right-2*pad
}

// headerHeight is shared by the sidebar brand row and the page title bar, so
// their edges line up. It covers the native title band and grows with text
// scale rather than letting controls overflow it.
func headerHeight(c *ui.Context, stacked bool) float32 {
	control := theme.Scale(c, theme.HeightButtonSM) + 2*theme.Scale(c, theme.Space8)
	if stacked {
		return headerBand(c) + control - theme.Scale(c, theme.Space8)
	}
	return max(headerBand(c), control)
}

// brand is the gallery mark; it truncates rather than reach the window
// controls.
func brand(c *ui.Context) {
	ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
		mygoui.Typography(c, "coss ui", mygoui.TypographyOptions{Type: theme.TypeH6, Truncate: true}).
			Shrink(1).MinWidth(0)
		mygoui.Typography(c, "MyGo", mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true, Truncate: true}).
			Shrink(2).MinWidth(0)
	})
}

func sidebar(c *ui.Context, s *State) {
	t := theme.Colors(c)
	bar := c.TitleBar()
	pad := theme.Scale(c, theme.Space12)
	ui.Column(c.Key("gallery-sidebar")).Width(theme.Scale(c, docsSidebarWidth)).FillHeight().Shrink(0).MinWidth(0).
		Background(t.Sidebar).BorderWidth(0, theme.Hairline, 0, 0).BorderColor(t.SidebarBorder).Children(func() {
		ui.Row(c).FillWidth().Height(headerHeight(c, headerStacked(c, true))).Shrink(0).DragWindow().
			Padding(0, pad, 0, bar.Left+pad+theme.Scale(c, theme.Space4)).
			AlignItems(ui.Center).Gap(theme.Space8).Label("Sidebar title bar").Children(func() { brand(c) })
		ui.Column(c).FillWidth().Padding(0, pad, theme.Space8, pad).Children(func() { search(c, s) })
		mygoui.ScrollArea(c, &s.Sidebar, mygoui.ScrollAreaOptions{
			Axis:  mygoui.ScrollVertical,
			Label: "Gallery navigation",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Edges: mygoui.FadeTop | mygoui.FadeBottom, Surface: t.Sidebar,
			},
		}, func() {
			ui.Column(c).FillWidth().Padding(theme.Space4, pad, theme.Space16, pad).Children(func() {
				pageList(c, s, func(page Page) { selectPage(s, page) })
			})
		}).Grow(1).Basis(0).MinHeight(0)
	})
}

// navigationToggle opens the sidebar, or the navigation sheet when compact.
func navigationToggle(c *ui.Context, s *State, compact bool) {
	icon, open := lucide.PanelLeft(), !s.SidebarClosed
	if compact {
		icon, open = lucide.Menu(), s.NavOpen
	}
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: "Toggle navigation", Icon: icon, IconOnly: true,
		Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
		OnClick: func() {
			if compact {
				s.NavOpen = !s.NavOpen
			} else {
				s.SidebarClosed = !s.SidebarClosed
			}
		},
	}).Role(ui.RoleToggleButton).Checked(open).Shrink(0).AlignSelf(ui.Center)
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

// groupRank orders the functional groups in the navigation. Pages inside a
// group keep the order they registered with.
var groupRank = map[string]int{
	"Overview": 0, "Foundations": 1, "Actions": 2, "Display": 3,
	"Forms": 4, "Collections": 5, "Navigation": 6, "Overlays": 7,
	"Dates": 8, "Color": 9, "Data": 10, "Charts": 11, "Chart Utilities": 12,
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
func pageList(c *ui.Context, s *State, choose func(Page)) {
	groups := pageGroups(filterPagesInOrder(s.filterPages()))
	ui.Column(c).FillWidth().Gap(theme.Space16).Children(func() {
		if len(groups) == 0 {
			ui.Row(c).Padding(theme.Space8).Children(func() {
				mygoui.Typography(c, "No components found", mygoui.TypographyOptions{
					Type: theme.TypeBodySM, Muted: true,
				})
			})
		}
		for _, group := range groups {
			ui.Column(c.Key(group.name)).FillWidth().Gap(theme.Space2).Children(func() {
				if group.name != "Overview" {
					groupHeading(c, group.name)
				}
				for _, page := range group.pages {
					navRow(c.Key(page.ID), page.Title, page.ID == s.Selected, func() { choose(page) })
				}
			})
		}
	})
}

// filterPagesInOrder keeps the navigation order of a filtered page list.
func filterPagesInOrder(list []Page) []Page {
	sort.SliceStable(list, func(i, j int) bool { return list[i].navRank < list[j].navRank })
	return list
}

// groupHeading is the small muted section label above a nav group.
func groupHeading(c *ui.Context, text string) {
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Muted).
		Padding(theme.Space4, theme.Space8, theme.Space4, theme.Space8).FillWidth().SingleLine().
		Role(ui.RoleHeading)
}

// navRow is one nav link: a 32 DIP rounded row that fills with the sidebar
// accent while its page is the selected one.
func navRow(c *ui.Context, title string, selected bool, onClick func()) {
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
	if e.FocusVisible() {
		e.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Stroke(r, t.SidebarRing, radius, 2) })
	}
	e.Children(func() {
		style.Text(ui.Text(c, title), c, ts).SingleLine()
	})
	e.OnClick(onClick)
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
	s.Panel = 0
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

func content(c *ui.Context, s *State, compact, sidebarOpen bool) {
	t := theme.Colors(c)
	ui.Column(c.Key("workbench")).Grow(1).Basis(0).MinWidth(0).FillHeight().Children(func() {
		pageHeader(c, s, compact, sidebarOpen)
		mygoui.ScrollArea(c.Key("workbench-content"), &s.Content, mygoui.ScrollAreaOptions{
			Axis: mygoui.ScrollVertical, Label: "Component content",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Edges: mygoui.FadeBottom, Surface: t.Background,
			},
		}, func() { docsPageColumn(c, s) }).Grow(1).Basis(0).MinHeight(0)
	})
}

// pageHeader is the draggable title bar over the page: navigation toggle,
// location and window-level actions, clear of the native window controls.
// When the controls cannot fit between those controls they sit below the
// native band instead.
func pageHeader(c *ui.Context, s *State, compact, sidebarOpen bool) {
	t := theme.Colors(c)
	bar := c.TitleBar()
	pad := theme.Scale(c, theme.Space12)
	stacked := headerStacked(c, sidebarOpen)
	top, left, right := float32(0), pad, pad+bar.Right
	if !sidebarOpen {
		left += bar.Left
	}
	if stacked {
		top, left, right = headerBand(c), pad, pad
	}
	ui.Row(c).FillWidth().Height(headerHeight(c, stacked)).Shrink(0).DragWindow().
		Padding(top, right, 0, left).AlignItems(ui.Center).Gap(theme.Space8).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).Label("Gallery title bar").
		Children(func() {
			navigationToggle(c, s, compact)
			breadcrumb(c, s, compact)
			services := c.Services()
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "GitHub", Icon: lucide.Github(), IconOnly: true,
				Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
				OnClick: func() { services.OpenURL(RepositoryURL) },
			}).Shrink(0).AlignSelf(ui.Center)
			appearanceButton(c, s)
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

// navigationSheet is the compact navigation: search stays in the header
// while the sheet body scrolls the page list.
func navigationSheet(c *ui.Context, s *State) {
	overlay.Sheet(c, &s.NavOpen, overlay.SheetOptions{
		Side: mygoui.SideLeft, Width: theme.Scale(c, 300), HideClose: true, BodyLabel: "Gallery navigation",
		Header: func() {
			ui.Column(c).FillWidth().Gap(theme.Space12).Children(func() {
				// The sheet covers the title band, so the brand clears
				// leading window controls; the sheet's 24 DIP padding
				// covers part of a taller band.
				bar := c.TitleBar()
				ui.Row(c).FillWidth().MinHeight(max(theme.Scale(c, theme.HeightButtonSM), bar.Height-theme.Space24)).
					PaddingStart(bar.Left).AlignItems(ui.Center).Children(func() { brand(c) })
				search(c, s)
			})
		},
		Body: func() { pageList(c, s, func(page Page) { selectPage(s, page) }) },
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
