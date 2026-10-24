package gallery

import (
	"time"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// The docs shell follows the HeroUI documentation page: a fixed navigation
// column, a readable main column and, when the window is wide enough, a
// sticky "On this page" rail beside the scrolling prose.
const (
	docsSidebarWidth = 260
	docsRailWidth    = 220
	docsMaxWidth     = 880
	docsRailFrom     = 1000
	docsCompactBelow = 760
)

// docsTOCSection is one entry of the right rail. The ID names the section
// state and the label is what the reader sees.
type docsTOCSection struct{ id, label string }

var docsTOCSections = []docsTOCSection{
	{"usage", "Usage"},
	{"examples", "Examples"},
	{"api", "API Reference"},
	{"tokens", "Theme tokens"},
}

// docsChromeState is the per-page state of the docs shell: handles that
// remember where each section landed and the copy feedback of the header
// actions. Handles survive rebuilds, unlike Elements.
type docsChromeState struct {
	Page             ui.Handle
	Usage            ui.Handle
	Examples         ui.Handle
	API              ui.Handle
	Tokens           ui.Handle
	MarkdownCopiedAt time.Time
	ImportCopiedAt   time.Time
	UsageCopiedAt    time.Time
}

// docsChromeOf returns the docs chrome state for a page.
func docsChromeOf(s *State, id string) *docsChromeState {
	return StateOf[docsChromeState](s, "docs:"+id)
}

// sectionHandle returns the handle of one of the four tracked sections.
func (d *docsChromeState) sectionHandle(id string) *ui.Handle {
	switch id {
	case "usage":
		return &d.Usage
	case "examples":
		return &d.Examples
	case "api":
		return &d.API
	case "tokens":
		return &d.Tokens
	}
	return nil
}

// docsCompact reports whether the window is too narrow for a persistent
// navigation column.
func docsCompact(c *ui.Context) bool {
	width, _ := c.Size()
	return width < theme.Scale(c, docsCompactBelow)
}

// docsRailVisible reports whether the "On this page" rail fits beside the
// prose. It hides before the main column would become unreadable.
func docsRailVisible(c *ui.Context) bool {
	width, _ := c.Size()
	return width >= theme.Scale(c, docsRailFrom)
}

// docsPadding is the page inset: 40 DIP at desktop width, 24 DIP when the
// window is compact.
func docsPadding(c *ui.Context) float32 {
	if docsCompact(c) {
		return theme.Scale(c, theme.Space24)
	}
	return theme.Scale(c, theme.Space40)
}

// docsContentWidth returns the width available to the page column.
func docsContentWidth(c *ui.Context, rail bool) float32 {
	width, _ := c.Size()
	available := width
	if !docsCompact(c) {
		available -= theme.Scale(c, docsSidebarWidth)
	}
	if rail {
		available -= theme.Scale(c, docsRailWidth)
	}
	// The max width caps the column's border box, padding included.
	column := min(theme.Scale(c, docsMaxWidth), max(0, available))
	return max(0, column-2*docsPadding(c))
}

// pageWidth is the width of a component page's main column, rail included.
func pageWidth(c *ui.Context) float32 {
	return docsContentWidth(c, docsRailVisible(c))
}

// overviewWidth is the width of the landing page, which has no rail.
func overviewWidth(c *ui.Context) float32 {
	return docsContentWidth(c, false)
}

// docsPageColumn builds the scrolling main column: header, sections and the
// project footer room. Binding the page handle lets the rail translate a
// section's committed bounds into a scroll offset.
func docsPageColumn(c *ui.Context, s *State) {
	page := currentPage(s)
	pad := docsPadding(c)
	e := ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, docsMaxWidth)).
		AlignSelf(ui.Center).Padding(pad, pad, theme.Scale(c, theme.Space64), pad).
		Gap(theme.Scale(c, theme.Space32))
	e.Children(func() {
		if page == nil {
			mygoui.EmptyState(c, mygoui.EmptyStateOptions{
				Title: "No components match", Description: s.Search,
			})
			return
		}
		if page.ID != "overview" {
			docsPageHeader(c, s, *page)
		}
		page.Build(c, s)
	})
	if page != nil {
		if _, ok := componentPage(page.ID); ok {
			e.Bind(&docsChromeOf(s, page.ID).Page)
		}
	}
}

// docsPageHeader is the title block of every documentation page: the group
// label, a 30/36 title, a 16/28 description and functional copy pills.
func docsPageHeader(c *ui.Context, s *State, page Page) {
	cp, component := componentPage(page.ID)
	ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space16)).Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Scale(c, theme.Space24)).
			AlignItems(ui.End).Children(func() {
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Scale(c, theme.Space8)).Children(func() {
				mygoui.Typography(c, page.Group, mygoui.TypographyOptions{
					Type: theme.TypeBodyXS, Weight: theme.WeightMedium, Muted: true,
				})
				mygoui.Typography(c, page.Title, mygoui.TypographyOptions{Type: theme.TypeH2})
				mygoui.Typography(c, page.Description, mygoui.TypographyOptions{
					Type: theme.TypeBody, Muted: true,
				})
			})
			if component {
				docsHeaderActions(c, cp, docsChromeOf(s, cp.ID))
			}
		})
		mygoui.Separator(c, mygoui.SeparatorOptions{})
	})
}

// docsHeaderActions renders the two copy pills. Both do real work: one
// copies a generated Markdown document, the other the page's import block.
func docsHeaderActions(c *ui.Context, p ComponentPage, chrome *docsChromeState) {
	ui.Row(c).Shrink(0).Wrap().AlignItems(ui.Center).
		Gap(theme.Scale(c, theme.Space8)).Children(func() {
		docsCopyPill(c, &chrome.MarkdownCopiedAt, "Copy Markdown", lucide.FileText(),
			func() string { return markdownDoc(p) })
		docsCopyPill(c, &chrome.ImportCopiedAt, "Copy import", lucide.Code(),
			func() string { return p.importStatement() })
	})
}

// docsCopyPill is a header action that reports "Copied" for the same window
// as the code block's copy control.
func docsCopyPill(c *ui.Context, at *time.Time, label string, icon *ui.SVG, text func() string) {
	shown, glyph := label, icon
	if !at.IsZero() {
		elapsed := c.Now().Sub(*at)
		if elapsed < codeCopyFeedback {
			shown, glyph = "Copied", lucide.Check()
			c.After(codeCopyFeedback - elapsed)
		} else {
			*at = time.Time{}
		}
	}
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: shown, Icon: glyph, Variant: mygoui.ButtonSecondary,
		Size: theme.ButtonSmall,
		OnClick: func() {
			c.WriteClipboard(text())
			*at = c.Now()
			c.Invalidate()
		},
	})
}

// docsRail renders the sticky "On this page" column. It lives beside the
// scroll container, so it never scrolls with the prose.
func docsRail(c *ui.Context, s *State, p ComponentPage) {
	chrome := docsChromeOf(s, p.ID)
	active := docsActiveSection(c, s, chrome)
	ui.Column(c).Width(theme.Scale(c, docsRailWidth)).Shrink(0).FillHeight().
		Padding(theme.Scale(c, theme.Space40), theme.Scale(c, theme.Space24),
			theme.Scale(c, theme.Space24), theme.Scale(c, theme.Space8)).
		Gap(theme.Scale(c, theme.Space12)).Children(func() {
		mygoui.Typography(c, "On this page", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Weight: theme.WeightMedium, Muted: true,
		})
		ui.Column(c).Gap(theme.Scale(c, theme.Space2)).Children(func() {
			for _, section := range docsTOCSections {
				section := section
				docsTOCLink(c, section, section.id == active, func() {
					docsScrollToSection(c, s, chrome, section.id)
				})
			}
		})
	})
}

// docsTOCLink is one rail entry. Its accessible name includes "Go to" so it
// stays distinct from the section heading it scrolls to.
func docsTOCLink(c *ui.Context, section docsTOCSection, selected bool, onClick func()) {
	t := theme.Colors(c)
	fill, hover, ink := ui.Transparent, t.DefaultSoft, t.Muted
	if selected {
		fill, hover, ink = t.AccentSoft, t.AccentSoftHover, t.AccentSoftForeground
	}
	radius := theme.Scale(c, theme.Space8)
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	e := ui.ButtonBase(c).FillWidth().MinHeight(theme.Scale(c, 28)).
		PaddingX(theme.Scale(c, theme.Space8)).PaddingY(0).Radius(radius).
		Justify(ui.Start).AlignItems(ui.Center).FocusRing(true).Shrink(0).
		Label("Go to " + section.label)
	e = style.Background(e, radius, fill, hover, hover)
	e.Children(func() {
		style.Text(ui.Text(c, section.label), c, ts).TextColor(ink).SingleLine()
	})
	e.OnClick(onClick)
}

// docsSection wraps one tracked section, binding the handle that the rail
// uses to find its committed position.
func docsSection(c *ui.Context, chrome *docsChromeState, id, title string, body func()) {
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Scale(c, theme.Space12))
	e.Children(func() {
		mygoui.Typography(c, title, mygoui.TypographyOptions{Type: theme.TypeH4})
		body()
	})
	if h := chrome.sectionHandle(id); h != nil {
		e.Bind(h)
	}
}

// docsSectionOffset returns the section's offset from the top of the page
// column, in committed bounds: stable while the content scrolls.
func docsSectionOffset(c *ui.Context, chrome *docsChromeState, id string) (float32, bool) {
	h := chrome.sectionHandle(id)
	if h == nil {
		return 0, false
	}
	section := h.Bounds(c)
	page := chrome.Page.Bounds(c)
	if section.H <= 0 || page.H <= 0 {
		return 0, false
	}
	return section.Y - page.Y, true
}

// docsScrollToSection scrolls the content so the section starts just below
// the top edge.
func docsScrollToSection(c *ui.Context, s *State, chrome *docsChromeState, id string) {
	offset, ok := docsSectionOffset(c, chrome, id)
	if !ok {
		return
	}
	s.Content.Y = max(0, offset-theme.Scale(c, theme.Space12))
	c.Invalidate()
}

// docsActiveSection reports the section at the rail's reading position. It
// is derived from committed bounds, not from the click that got there.
func docsActiveSection(c *ui.Context, s *State, chrome *docsChromeState) string {
	active := ""
	limit := s.Content.Y + theme.Scale(c, 220)
	for _, section := range docsTOCSections {
		if offset, ok := docsSectionOffset(c, chrome, section.id); ok && offset <= limit {
			active = section.id
		}
	}
	return active
}
