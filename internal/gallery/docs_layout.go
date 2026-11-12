package gallery

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
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

// docsPageColumn is the page itself: title, description and, for
// components, the view switch, all on the content column's left edge.
func docsPageColumn(c *ui.Context, s *State) {
	page, ok := currentPage(s)
	pad := docsPadding(c)
	ui.Column(c).FillWidth().MinWidth(0).MaxWidth(theme.Scale(c, docsMaxWidth)).
		AlignSelf(ui.Start).Padding(pad, pad, theme.Scale(c, theme.Space64), pad).
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
			if _, isComponent := componentPage(page.ID); isComponent {
				workbenchViews(c, s)
			}
		}
		ui.Column(c.Key(page.ID)).FillWidth().MinWidth(0).Children(func() {
			page.Build(c, s)
		})
	})
}

// workbenchViews switches a component page between its views. Segments keep
// their text width and wrap two to a row when the column is too narrow.
func workbenchViews(c *ui.Context, s *State) {
	t := theme.Colors(c)
	labels := [...]string{"Preview", "Code", "API", "Tokens"}
	available := docsContentWidth(c, s) - theme.Scale(c, 4)
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
	track := parts.Track.MinWidth(0).Wrap().AlignSelf(ui.Start).
		Padding(theme.Scale(c, 2)).Gap(gap).
		Radius(theme.RadiusScale().LG).Background(t.MutedBackground).Label("Component view")
	if wrap {
		track.FillWidth()
	}
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
