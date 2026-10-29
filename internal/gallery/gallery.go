package gallery

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/fonts"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TitleBarHeight is the gallery's native title-bar height in DIPs.
const TitleBarHeight = 52

// View builds the workbench. Theme application happens here so appearance
// switches affect every specimen immediately.
func View(c *ui.Context, s *State) {
	zoom(c, s)
	applyAppearance(c, s)
	overlaysStateOf(s).toasts.View(c)
	t := theme.Colors(c)
	compact := docsCompact(c)
	c.OnShortcut(ui.Cmd, ui.KeyK, func() {
		if compact {
			s.NavOpen = true
		} else {
			s.SidebarClosed = false
		}
		s.SearchFocus.Focus(c)
	})
	ui.Column(c).Fill().MinWidth(0).Clip().Background(t.Background).TextColor(t.Foreground).Children(func() {
		titleBar(c)
		ui.Row(c).Grow(1).Basis(0).MinHeight(0).MinWidth(0).AlignItems(ui.Stretch).Children(func() {
			if !compact && !s.SidebarClosed {
				sidebar(c, s)
			}
			content(c, s, compact)
		})
	})
}

func titleBar(c *ui.Context) {
	t := theme.Colors(c)
	bar := c.TitleBar()
	width, _ := c.Size()
	height := max(float32(TitleBarHeight)/c.Zoom(), bar.Height)
	padding := theme.Scale(c, 16)
	available := max(0, width-bar.Left-bar.Right-2*padding)
	heading := theme.H6()
	headingSize := theme.Scale(c, heading.Size)
	brandWidth, _ := c.MeasureText(0, ui.Span{
		Text: "coss ui", Font: theme.HeadingFontFamily(c), Size: headingSize,
		Weight: heading.Weight, LetterSpacing: heading.Tracking * headingSize,
	})
	body := theme.TypeBodyXS.Style()
	subtitleWidth, _ := c.MeasureText(0, ui.Span{
		Text: "MyGo", Font: theme.FontFamily(c), Size: theme.Scale(c, body.Size), Weight: body.Weight,
	})
	inline := available >= brandWidth && theme.Scale(c, heading.LineHeight) <= height
	brand := func(available float32) {
		mygoui.Typography(c, "coss ui", mygoui.TypographyOptions{Type: theme.TypeH6, Truncate: true}).
			MaxWidth(max(0, available))
		if available >= brandWidth+theme.Space8+subtitleWidth {
			mygoui.Typography(c, "MyGo", mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
		}
	}
	title := ui.Row(c).FillWidth().Height(height).Shrink(0).DragWindow().
		Padding(0, bar.Right+padding, 0, bar.Left+padding).
		AlignItems(ui.Center).Gap(theme.Space8).
		Background(t.Sidebar).Label("Gallery title bar")
	if inline {
		title.BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.SidebarBorder)
	}
	title.Children(func() {
		if inline {
			brand(available)
		}
	})
	if !inline {
		ui.Row(c).FillWidth().Height(theme.Scale(c, heading.LineHeight+8)).Shrink(0).DragWindow().
			PaddingX(padding).AlignItems(ui.Center).Gap(theme.Space8).
			Background(t.Sidebar).BorderWidth(0, 0, theme.Hairline, 0).
			BorderColor(t.SidebarBorder).Children(func() { brand(max(0, width-2*padding)) })
	}
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
	if s.SystemFont {
		family = ""
	}
	theme.Configure(theme.Preferences{
		FontFamily: family, HeadingFontFamily: family,
		Accent: accent, Radius: s.Radius,
	})
	theme.ApplyMode(c, s.Dark)
}

func sidebar(c *ui.Context, s *State) {
	t := theme.Colors(c)
	ui.Column(c.Key("gallery-sidebar")).Width(theme.Scale(c, 260)).FillHeight().Shrink(0).MinWidth(0).
		Padding(theme.Space16, theme.Space12, theme.Space12, theme.Space12).Gap(theme.Space12).
		Background(t.Sidebar).
		BorderWidth(0, theme.Hairline, 0, 0).BorderColor(t.SidebarBorder).Children(func() {
		search(c, s)
		mygoui.ScrollArea(c, &s.Sidebar, mygoui.ScrollAreaOptions{
			Axis:  mygoui.ScrollVertical,
			Label: "Components",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Edges: mygoui.FadeTop | mygoui.FadeBottom, Surface: t.Sidebar,
			},
		}, func() {
			pageList(c, s, func(page Page) { selectPage(s, page) })
		}).Grow(1).Basis(0).MinHeight(0)
		ui.Row(c).Padding(theme.Space8).Children(func() {
			mygoui.Typography(c, "⌘/Ctrl K", mygoui.TypographyOptions{
				Type: theme.TypeBodyXS, Muted: true,
			})
		})
	})
}

func buttonFor(on bool) mygoui.ButtonVariant {
	if on {
		return mygoui.ButtonSecondary
	}
	return mygoui.ButtonGhost
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

// pageList renders the grouped page links, shared by the sidebar and the
// compact navigation drawer.
func pageList(c *ui.Context, s *State, choose func(Page)) {
	groups := pageGroups(filterPagesInOrder(s.filterPages()))
	searching := strings.TrimSpace(s.Search) != ""
	ui.Column(c).Gap(theme.Space4).Children(func() {
		if len(groups) == 0 {
			mygoui.Typography(c, "No components found", mygoui.TypographyOptions{
				Type: theme.TypeBodySM, Muted: true,
			})
		}
		for _, group := range groups {
			ui.Column(c.Key(group.name)).Gap(theme.Space2).Children(func() {
				if group.name == "Overview" {
					page := group.pages[0]
					navRow(c.Key(page.ID), page.Title, page.ID == s.Selected, func() { choose(page) })
					return
				}
				if searching {
					groupHeading(c, group.name)
				} else {
					categoryRow(c, s, group)
				}
				if !searching && s.OpenGroup != group.name {
					return
				}
				ui.Column(c).PaddingStart(theme.Space12).Gap(theme.Space2).Children(func() {
					for _, page := range group.pages {
						navRow(c.Key(page.ID), page.Title, page.ID == s.Selected, func() { choose(page) })
					}
				})
			})
		}
	})
}

func categoryRow(c *ui.Context, s *State, group pageGroup) {
	t := theme.Colors(c)
	open := s.OpenGroup == group.name
	active, ok := pageFor(s.Selected)
	fill := ui.Transparent
	if ok && active.Group == group.name {
		fill = t.SidebarAccent
	}
	radius := theme.RadiusScale().LG
	e := ui.ButtonBase(c.Key("category")).FillWidth().Height(theme.Scale(c, 36)).
		PaddingX(theme.Scale(c, theme.Space8)).Gap(theme.Space8).
		AlignItems(ui.Center).Justify(ui.Start).FocusRing(false).
		Role(ui.RoleToggleButton).Checked(open).Label(group.name).TextColor(t.SidebarForeground)
	if ok && active.Group == group.name {
		e.Description(locale.M().Selected)
	}
	e = style.Background(e, radius, fill, t.SidebarAccent, t.SidebarAccent)
	if e.FocusVisible() {
		e.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Stroke(r, t.SidebarRing, radius, 2) })
	}
	e.Children(func() {
		icon := lucide.ChevronRight()
		if open {
			icon = lucide.ChevronDown()
		}
		ui.Icon(c, icon).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
		style.Text(ui.Text(c, group.name), c, theme.SM()).Grow(1).MinWidth(0).SingleLine()
		mygoui.Typography(c, fmt.Sprint(len(group.pages)), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
	e.OnClick(func() {
		if open {
			s.OpenGroup = ""
		} else {
			s.OpenGroup = group.name
		}
	})
}

// filterPagesInOrder keeps the navigation order of a filtered page list.
func filterPagesInOrder(list []Page) []Page {
	sort.SliceStable(list, func(i, j int) bool { return list[i].navRank < list[j].navRank })
	return list
}

// groupHeading is the native library's small muted section label above a nav group.
func groupHeading(c *ui.Context, text string) {
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Muted).
		Padding(theme.Space4, theme.Space8, theme.Space4, theme.Space8).FillWidth().SingleLine()
}

// navRow is one nav link: a 32 DIP rounded row that fills with the accent
// soft colors while its page is the selected one.
func navRow(c *ui.Context, title string, selected bool, onClick func()) {
	t := theme.Colors(c)
	fill, hover, ink := ui.Transparent, t.Default, t.Foreground
	if selected {
		fill, hover, ink = t.SidebarAccent, t.SidebarAccent, t.SidebarAccentForeground
	}
	radius := theme.RadiusScale().LG
	ts := theme.SM()
	ts.Weight = theme.WeightMedium
	e := ui.ButtonBase(c).FillWidth().MinHeight(theme.Scale(c, theme.HeightButtonSM)).
		PaddingX(theme.Scale(c, theme.Space10)).PaddingY(0).Radius(radius).
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

func search(c *ui.Context, s *State) {
	compact := docsCompact(c)
	form.Input(c.Key("gallery-search"), &s.Search, form.InputOptions{
		Placeholder: "Search components", AriaLabel: "Search components", FullWidth: true,
		OnChange: func() { s.Sidebar.Y, s.NavScroll.Y = 0, 0 },
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
			s.Sidebar.Y, s.NavScroll.Y = 0, 0
			return true
		}
		if compact && s.NavOpen {
			s.NavOpen = false
			return true
		}
		return false
	}).Bind(&s.SearchFocus)
}

func selectPage(s *State, page Page) {
	if s.flyout != nil {
		s.flyout.Close()
	}
	s.Selected = page.ID
	s.OpenGroup = page.Group
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
			})
		}
	})
}

func content(c *ui.Context, s *State, compact bool) {
	t := theme.Colors(c)
	ui.Column(c.Key("workbench")).Grow(1).Basis(0).MinWidth(0).FillHeight().Children(func() {
		workbenchToolbar(c, s, compact)
		if s.AppearanceOpen {
			appearanceControls(c, s)
		}
		if compact && s.NavOpen {
			compactNavigation(c, s)
		}
		mygoui.ScrollArea(c.Key("workbench-content"), &s.Content, mygoui.ScrollAreaOptions{
			Axis: mygoui.ScrollVertical, Label: "Component content",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Edges: mygoui.FadeTop | mygoui.FadeBottom, Surface: t.Background,
			},
		}, func() { docsPageColumn(c, s) }).Grow(1).Basis(0).MinHeight(0)
	})
}

func workbenchToolbar(c *ui.Context, s *State, compact bool) {
	t := theme.Colors(c)
	page, ok := currentPage(s)
	ui.Column(c).FillWidth().Shrink(0).Padding(theme.Space12, theme.Space16, theme.Space12, theme.Space16).
		Gap(theme.Space12).BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
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
			}).Role(ui.RoleToggleButton).Checked(open)
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
				if compact && s.NavOpen {
					mygoui.Typography(c, "Browse", mygoui.TypographyOptions{Type: theme.TypeH5})
				} else if ok {
					mygoui.Typography(c, page.Title, mygoui.TypographyOptions{Type: theme.TypeH5})
				}
			})
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Appearance", Icon: lucide.SlidersHorizontal(), IconOnly: true,
				Variant: buttonFor(s.AppearanceOpen), Size: theme.ButtonSmall,
				OnClick: func() { s.AppearanceOpen = !s.AppearanceOpen },
			}).Role(ui.RoleToggleButton).Checked(s.AppearanceOpen)
		})
		if _, isComponent := componentPage(page.ID); ok && isComponent {
			workbenchViews(c, s, compact)
		}
	})
}

func workbenchViews(c *ui.Context, s *State, compact bool) {
	t := theme.Colors(c)
	labels := [...]string{"Preview", "Code", "API", "Tokens"}
	width, _ := c.Size()
	if !compact && !s.SidebarClosed {
		width -= theme.Scale(c, docsSidebarWidth)
	}
	available := max(0, width-2*theme.Space16-theme.Scale(c, 4))
	gap := theme.Scale(c, 2)
	var widths [len(labels)]float32
	total := gap * (float32(len(labels)) - 1)
	for i, label := range labels {
		textWidth, _ := c.MeasureText(0, ui.Span{
			Text: label, Font: theme.FontFamily(c),
			Size: theme.Scale(c, theme.SM().Size), Weight: theme.WeightMedium,
		})
		widths[i] = textWidth + theme.Scale(c, 20)
		total += widths[i]
	}
	wrap := total > available
	parts := ui.SegmentedBase(c.Key("workbench-views"), &s.Panel, len(labels))
	track := parts.Track.FillWidth().MinWidth(0).Wrap().
		Padding(theme.Scale(c, 2)).Gap(gap).
		Radius(theme.RadiusScale().LG).Background(t.MutedBackground).Label("Component view")
	track.Children(func() {
		for i, label := range labels {
			minWidth := min(widths[i], available)
			if wrap {
				minWidth = max(minWidth, max(0, (available-gap)/2))
			}
			segment := parts.Segment(i).Shrink(0).Height(theme.Scale(c, 30)).
				MinWidth(minWidth).PaddingX(theme.Scale(c, 10)).Radius(theme.RadiusScale().MD).
				FocusRing(true).Label(label)
			if wrap {
				segment.Grow(1)
			}
			ink := t.MutedForeground.Alpha(0.72)
			if segment.Hovered() {
				ink = t.MutedForeground
			}
			if i == s.Panel {
				ink = t.Foreground
				face := t.Background
				if c.Theme().Dark {
					face = t.Input
				}
				segment = style.Shadow(segment.Background(face), c, theme.SurfaceShadows(c.Theme().Dark))
			}
			segment.Children(func() {
				style.Text(ui.Text(c, label), c, theme.SM()).
					FontWeight(theme.WeightMedium).TextColor(ink).SingleLine()
			})
		}
	})
	if track.Changed() {
		s.Content.Y = 0
	}
}

func appearanceControls(c *ui.Context, s *State) {
	t := theme.Colors(c)
	ui.Column(c).FillWidth().Shrink(0).Padding(theme.Space12, theme.Space16, theme.Space12, theme.Space16).
		Gap(theme.Space8).Background(t.SurfaceSecondary).
		BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).Children(func() {
		ui.Row(c).Wrap().Gap(theme.Space4).Children(func() {
			for _, option := range []struct {
				label string
				on    bool
				set   func()
			}{
				{"Light", !s.Dark, func() { s.setDark(false) }},
				{"Dark", s.Dark, func() { s.setDark(true) }},
				{"Cal Sans", !s.SystemFont, func() { s.SystemFont = false }},
				{"System", s.SystemFont, func() { s.SystemFont = true }},
			} {
				mygoui.Button(c.Key(option.label), mygoui.ButtonOptions{
					Label: option.label, Variant: buttonFor(option.on), Size: theme.ButtonSmall, OnClick: option.set,
				}).Role(ui.RoleToggleButton).Checked(option.on)
			}
		})
		ui.Row(c).Wrap().Gap(theme.Space4).Children(func() {
			for _, radius := range []float32{4, 10, 16} {
				mygoui.Button(c.Key(radius), mygoui.ButtonOptions{
					Label: fmt.Sprintf("Radius %d", int(radius)), Variant: buttonFor(s.Radius == radius),
					Size: theme.ButtonSmall, OnClick: func() { s.Radius = radius },
				}).Role(ui.RoleToggleButton).Checked(s.Radius == radius)
			}
		})
	})
}

func compactNavigation(c *ui.Context, s *State) {
	t := theme.Colors(c)
	_, height := c.Size()
	ui.Column(c).FillWidth().Shrink(0).Padding(theme.Space12).Gap(theme.Space8).
		Background(t.Sidebar).BorderWidth(0, 0, theme.Hairline, 0).
		BorderColor(t.SidebarBorder).Children(func() {
		search(c, s)
		mygoui.ScrollArea(c, &s.NavScroll, mygoui.ScrollAreaOptions{
			Axis: mygoui.ScrollVertical, Height: min(theme.Scale(c, 300), height*0.35), Label: "Components",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Edges: mygoui.FadeTop | mygoui.FadeBottom, Surface: t.Sidebar,
			},
		}, func() { pageList(c, s, func(page Page) { selectPage(s, page) }) }).FillWidth()
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
