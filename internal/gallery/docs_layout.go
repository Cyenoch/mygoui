package gallery

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const (
	docsSidebarWidth = 260
	docsMaxWidth     = 1000
	docsCompactBelow = 760
)

type docsChromeState struct {
	MarkdownCopiedAt time.Time
	ImportCopiedAt   time.Time
	UsageCopiedAt    time.Time
}

func docsChromeOf(s *State, id string) *docsChromeState {
	return StateOf[docsChromeState](s, "docs:"+id)
}

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

func docsContentWidth(c *ui.Context, s *State) float32 {
	width, _ := c.Size()
	if !docsCompact(c) && !s.SidebarClosed {
		width -= theme.Scale(c, docsSidebarWidth)
	}
	return max(0, min(theme.Scale(c, docsMaxWidth), width)-2*docsPadding(c))
}

func pageWidth(c *ui.Context, s *State) float32     { return docsContentWidth(c, s) }
func overviewWidth(c *ui.Context, s *State) float32 { return docsContentWidth(c, s) }

func docsPageColumn(c *ui.Context, s *State) {
	page, ok := currentPage(s)
	pad := docsPadding(c)
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, docsMaxWidth)).
		AlignSelf(ui.Center).Padding(pad, pad, theme.Scale(c, theme.Space64), pad).
		Children(func() {
			if !ok {
				mygoui.EmptyState(c, mygoui.EmptyStateOptions{Title: "No components found"})
				return
			}
			ui.Column(c.Key(page.ID)).FillWidth().MinWidth(0).Children(func() {
				page.Build(c, s)
			})
		})
}
