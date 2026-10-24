package gallery

import (
	"fmt"
	"sort"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// View builds the workbench. Theme application happens here so appearance
// switches affect every specimen immediately.
func View(c *ui.Context, s *State) {
	applyAppearance(c, s)
	t := theme.Colors(c)
	width, _ := c.Size()
	compact := width < theme.Scale(c, 760)
	ui.Row(c).Fill().MinWidth(0).AlignItems(ui.Stretch).Clip().Background(t.Background).TextColor(t.Foreground).Children(func() {
		if !compact {
			sidebar(c, s)
		}
		content(c, s, compact)
	})
}

var accentSwatches = []ui.Color{
	ui.Oklch(0.6204, 0.195, 253.83),
	ui.Oklch(0.62, 0.22, 300),
	ui.Oklch(0.65, 0.19, 25),
	ui.Oklch(0.62, 0.16, 150),
	ui.Oklch(0.72, 0.16, 80),
}

func applyAppearance(c *ui.Context, s *State) {
	accent := ui.Color{}
	if s.Accent > 0 && s.Accent < len(accentSwatches) {
		accent = accentSwatches[s.Accent]
	}
	theme.Configure(theme.Preferences{Accent: accent, Radius: s.Radius})
	theme.ApplyMode(c, s.Dark)
}

func sidebar(c *ui.Context, s *State) {
	t := theme.Colors(c)
	ui.Column(c).Width(theme.Scale(c, 260)).FillHeight().Shrink(0).MinWidth(0).
		Padding(theme.Space12, theme.Space16, theme.Space12, theme.Space16).Gap(theme.Space12).
		Background(t.Background).
		BorderWidth(0, theme.Hairline, 0, 0).BorderColor(t.Border).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space2).Children(func() {
				mygoui.Typography(c, "MyGo UI", mygoui.TypographyOptions{Type: theme.TypeH6})
				mygoui.Typography(c, "HeroUI components, native Go", mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Muted: true,
				})
			})
			label := "Dark"
			if s.Dark {
				label = "Light"
			}
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: label, Icon: lucide.SunMoon(), IconOnly: true, Variant: mygoui.ButtonGhost,
				Size: theme.ButtonSmall, OnClick: func() { s.Dark = !s.Dark },
			})
		})
		search(c, s)
		mygoui.ScrollShadow(c, &s.Sidebar, mygoui.ScrollShadowOptions{
			Orientation: mygoui.ScrollShadowVertical,
			Surface:     t.Background,
			Label:       "Components",
		}, func() {
			pageList(c, s, func(page Page) {
				s.Selected = page.ID
				s.Content.Y = 0
			})
		}).Grow(1).MinHeight(0)
		ui.Column(c).Gap(theme.Space6).Children(func() {
			mygoui.Separator(c, mygoui.SeparatorOptions{})
			swatches(c, s)
			ui.Row(c).Gap(theme.Space4).Children(func() {
				for _, radius := range []float32{4, 8, 12} {
					mygoui.Button(c.Key(radius), mygoui.ButtonOptions{
						Label: fmt.Sprintf("R%d", int(radius)), Variant: buttonFor(s.Radius == radius),
						Size: theme.ButtonSmall, OnClick: func() { s.Radius = radius },
					})
				}
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
	"Dates": 8, "Color": 9, "Data": 10,
}

// navPages returns the pages in navigation order: by group rank, then the
// order each page registered with.
func navPages() []Page {
	list := Pages()
	sort.SliceStable(list, func(i, j int) bool {
		ri, iok := groupRank[list[i].Group]
		rj, jok := groupRank[list[j].Group]
		if !iok {
			ri = 99
		}
		if !jok {
			rj = 99
		}
		return ri < rj
	})
	return list
}

// pageGroups folds pages into one block per group, never repeating a heading
// when a search splits a group.
func pageGroups(list []Page) []pageGroup {
	var groups []pageGroup
	index := map[string]int{}
	for _, page := range list {
		i, ok := index[page.Group]
		if !ok {
			i = len(groups)
			index[page.Group] = i
			groups = append(groups, pageGroup{name: page.Group})
		}
		groups[i].pages = append(groups[i].pages, page)
	}
	return groups
}

// pageList renders the grouped page links, shared by the sidebar and the
// compact navigation drawer.
func pageList(c *ui.Context, s *State, choose func(Page)) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		for _, group := range pageGroups(filterPagesInOrder(s.filterPages())) {
			ui.Column(c).Gap(theme.Space2).Children(func() {
				groupHeading(c, group.name)
				ui.Column(c).Gap(theme.Space2).Children(func() {
					for _, page := range group.pages {
						navRow(c, page.Title, page.ID == s.Selected, func() { choose(page) })
					}
				})
			})
		}
	})
}

// filterPagesInOrder keeps the navigation order of a filtered page list.
func filterPagesInOrder(list []Page) []Page {
	rank := func(p Page) int {
		if r, ok := groupRank[p.Group]; ok {
			return r
		}
		return 99
	}
	sort.SliceStable(list, func(i, j int) bool { return rank(list[i]) < rank(list[j]) })
	return list
}

// groupHeading is HeroUI's small muted section label above a nav group.
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
		fill, hover, ink = t.AccentSoft, t.AccentSoftHover, t.AccentSoftForeground
	}
	radius := theme.Scale(c, theme.Space12)
	ts := theme.SM()
	ts.Weight = theme.WeightMedium
	e := ui.ButtonBase(c).FillWidth().MinHeight(theme.Scale(c, theme.HeightButtonSM)).
		PaddingX(theme.Scale(c, theme.Space10)).PaddingY(0).Radius(radius).
		Justify(ui.Start).AlignItems(ui.Center).FocusRing(true).Shrink(0).
		Label(title).TextColor(ink)
	e = style.Background(e, radius, fill, hover, hover)
	e.Children(func() {
		style.Text(ui.Text(c, title), c, ts).SingleLine()
	})
	e.OnClick(onClick)
}

func search(c *ui.Context, s *State) {
	form.SearchField(c, &s.Search, form.SearchFieldOptions{
		Placeholder: "Search components",
		AriaLabel:   "Search components",
		FullWidth:   true,
	})
}

func swatches(c *ui.Context, s *State) {
	ui.Row(c).Gap(theme.Space4).AlignItems(ui.Center).Children(func() {
		for i := range accentSwatches {
			mygoui.Button(c.Key(i), mygoui.ButtonOptions{
				Label: fmt.Sprintf("Accent %d", i+1), IconOnly: true, Size: theme.ButtonSmall,
				Variant: buttonFor(s.Accent == i),
				OnClick: func() { s.Accent = i },
				Content: func() {
					ui.Box(c).Size(theme.Scale(c, theme.Space12), theme.Scale(c, theme.Space12)).
						Radius(theme.Scale(c, theme.Space12)).Background(accentSwatches[i])
				},
			})
		}
	})
}

func content(c *ui.Context, s *State, compact bool) {
	t := theme.Colors(c)
	ui.Column(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Children(func() {
		if compact {
			ui.Column(c).BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border).Children(func() {
				ui.Row(c).Padding(theme.Space12).Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Toggle navigation", Icon: lucide.Menu(), IconOnly: true,
						Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall,
						OnClick: func() { s.NavOpen = !s.NavOpen },
					})
					mygoui.Typography(c, "MyGo UI", mygoui.TypographyOptions{Type: theme.TypeH6})
				})
				if s.NavOpen {
					mygoui.ScrollShadow(c, &s.NavScroll, mygoui.ScrollShadowOptions{
						Orientation: mygoui.ScrollShadowVertical,
						Surface:     t.Background,
						Height:      theme.Scale(c, 420),
						Label:       "Components",
					}, func() {
						ui.Column(c).Padding(0, theme.Space12, theme.Space12, theme.Space12).
							Children(func() {
								pageList(c, s, func(page Page) {
									s.Selected = page.ID
									s.Content.Y = 0
									s.NavOpen = false
								})
							})
					})
				}
			})
		}
		ui.Row(c).Grow(1).Basis(0).MinHeight(0).MinWidth(0).
			AlignItems(ui.Stretch).Children(func() {
			ui.Scroll(c).TrackScroll(&s.Content).Grow(1).Basis(0).MinWidth(0).
				Children(func() { docsPageColumn(c, s) })
			if page := currentPage(s); page != nil && docsRailVisible(c) {
				if cp, ok := componentPage(page.ID); ok {
					docsRail(c, s, cp)
				}
			}
		})
	})
}

func currentPage(s *State) *Page {
	for i := range pages {
		if pages[i].ID == s.Selected {
			return &pages[i]
		}
	}
	if len(pages) > 0 {
		return &pages[0]
	}
	return nil
}
