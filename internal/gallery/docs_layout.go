package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const (
	docsSidebarWidth = 260
	docsMaxWidth     = 1000
	landingMaxWidth  = 1120
	docsCompactBelow = 760
	docsIndexWidth   = 208
	// docsIndexAbove is the width beside the sidebar from which a component
	// page keeps its index next to the content.
	docsIndexAbove = 900
)

func docsCompact(c *ui.Context) bool {
	width, _ := c.Size()
	return width < theme.Scale(c, docsCompactBelow)
}

func docsPadding(c *ui.Context) float32 {
	if docsCompact(c) {
		return theme.Scale(c, theme.Space16)
	}
	return theme.Scale(c, theme.Space32)
}

// isLanding reports whether the landing page is showing. It owns the whole
// window: no sidebar, a wider centered column and its own title bar.
func isLanding(s *State) bool {
	page, ok := currentPage(s)
	return ok && page.ID == "overview"
}

// sidebarShown reports whether the sidebar sits beside the content.
func sidebarShown(c *ui.Context, s *State) bool {
	return !docsCompact(c) && !s.SidebarClosed && !isLanding(s)
}

// docsAreaWidth is the width right of the sidebar.
func docsAreaWidth(c *ui.Context, s *State) float32 {
	width, _ := c.Size()
	if sidebarShown(c, s) {
		width -= theme.Scale(c, docsSidebarWidth)
	}
	return width
}

// indexedPage is the component page whose index shows beside the content.
func indexedPage(c *ui.Context, s *State) (ComponentPage, bool) {
	if docsAreaWidth(c, s) < theme.Scale(c, docsIndexAbove) {
		return ComponentPage{}, false
	}
	page, ok := currentPage(s)
	if !ok {
		return ComponentPage{}, false
	}
	return componentPage(page.ID)
}

func docsContentWidth(c *ui.Context, s *State) float32 {
	width := docsAreaWidth(c, s)
	if _, ok := indexedPage(c, s); ok {
		width -= theme.Scale(c, docsIndexWidth)
	}
	return max(0, min(theme.Scale(c, docsMaxWidth), width)-2*docsPadding(c))
}

func pageWidth(c *ui.Context, s *State) float32 { return docsContentWidth(c, s) }

// overviewWidth is the landing page's centered column.
func overviewWidth(c *ui.Context, _ *State) float32 {
	width, _ := c.Size()
	return max(0, min(theme.Scale(c, landingMaxWidth), width)-2*docsPadding(c))
}

// docsPageColumn is the page itself: title, description and content, all on
// the content column's left edge.
func docsPageColumn(c *ui.Context, s *State) {
	page, ok := currentPage(s)
	pad := docsPadding(c)
	maxWidth := float32(docsMaxWidth)
	if ok && page.ID == "overview" {
		maxWidth = landingMaxWidth
	}
	// The column sits in the middle of the room between the sidebar and
	// the page index, as the coss docs do, not against the sidebar.
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, maxWidth)).
		AlignSelf(ui.Center).Padding(pad, pad, theme.Scale(c, theme.Space64), pad).
		Gap(theme.Scale(c, theme.Space24)).Children(func() {
		if !ok {
			mygoui.EmptyState(c, mygoui.EmptyStateOptions{Title: "No components found"})
			return
		}
		if page.ID != "overview" {
			ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space8).Children(func() {
				mygoui.Typography(c, page.Title, mygoui.TypographyOptions{Type: theme.TypeH2})
				if page.Description != "" {
					mygoui.Typography(c, page.Description, mygoui.TypographyOptions{
						Type: theme.TypeBody, Muted: true,
					})
				}
			})
		}
		ui.Column(c.Key(page.ID)).FillWidth().MinWidth(0).Children(func() {
			page.Build(c, s)
		})
	})
}
