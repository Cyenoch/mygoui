package navigation

import (
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PaginationSize chooses the height of the page controls.
type PaginationSize uint8

const (
	PaginationMedium PaginationSize = iota
	PaginationSmall
	PaginationLarge
)

// PaginationOptions configures page navigation. Page is 1-based and
// caller-owned; nil keeps it internally. Siblings sets the pages shown around
// the current one (default 0) and Boundaries the pages at each end (default
// 1). PreviousLabel and NextLabel override the locale's names. Summary builds
// the line beside the controls, as "1-10 of 120 results".
type PaginationOptions struct {
	Page          *int
	Pages         int
	Size          PaginationSize
	Siblings      int
	Boundaries    int
	Label         string
	PreviousLabel string
	NextLabel     string
	// Summary builds a muted line before the controls. With it, the row
	// splits between the summary and the controls on a wide pagination and
	// stacks them below 640 DIP, as upstream's sm breakpoint does.
	Summary  func()
	OnChange func(page int)
}

// Pagination builds previous and next controls around numbered pages with
// ellipses for gaps. Previous is disabled on the first page and Next on the
// last; clicking the current page does nothing.
func Pagination(c *ui.Context, o PaginationOptions) ui.Element {
	t := theme.Colors(c)
	pages := max(o.Pages, 1)
	root := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(scale(c, theme.Space16)).
		Role(ui.RoleGroup).MinWidth(0)
	if o.Summary == nil {
		root.Justify(ui.Center)
	}
	if o.Label != "" {
		root.Label(o.Label)
	}
	// stacked is the narrow arrangement, measured once per crossing of the
	// upstream breakpoint.
	stacked := ui.Local(root, "pagination-stacked", func() bool { return false })
	if o.Summary != nil && *stacked {
		root.Column().AlignItems(ui.Start)
	}
	page := o.Page
	if page == nil {
		page = ui.Local(root, "pagination-page", func() int { return 1 })
	}
	*page = max(1, min(*page, pages))

	size, font, navPad, radius := float32(32), theme.SM(), float32(10), theme.RadiusScale().MD
	switch o.Size {
	case PaginationSmall:
		size, font, navPad = 28, theme.XS(), 8
	case PaginationLarge:
		size, font, navPad = 36, theme.Base(), 12
	}
	size = scale(c, size)

	previous := o.PreviousLabel
	if previous == "" {
		previous = locale.M().Previous
	}
	nextLabel := o.NextLabel
	if nextLabel == "" {
		nextLabel = locale.M().Next
	}

	windowWidth, _ := c.Size()
	narrow := windowWidth < 640
	nav := func(label string, icon *ui.SVG, to int, enabled, trailing bool) {
		padding := navPad
		if narrow {
			padding = 0
		}
		btn := ui.ButtonBase(c).Height(size).MinWidth(size).PaddingX(scale(c, padding)).
			Center().Gap(scale(c, theme.Space6)).Shrink(0).Radius(radius).FocusRing(true).Label(label)
		btn = style.Background(btn, radius, ui.Transparent, t.Accent, t.Accent)
		btn.TextColor(t.DefaultForeground)
		if enabled {
			target := to
			btn.OnClick(func() {
				if *page == target {
					return
				}
				*page = target
				if o.OnChange != nil {
					o.OnChange(target)
				}
			})
		}
		btn.Children(func() {
			if icon != nil && !trailing {
				ui.Icon(c, icon).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
			}
			if !narrow {
				style.Text(ui.Text(c, label), c, font).FontWeight(theme.WeightMedium)
			}
			if icon != nil && trailing {
				ui.Icon(c, icon).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
			}
		})
		style.Disabled(btn, !enabled)
	}

	root.Children(func() {
		if o.Summary != nil {
			summary := ui.Row(c).AlignItems(ui.Center).Gap(scale(c, theme.Space8)).Shrink(0)
			summary = style.Text(summary, c, font).TextColor(t.Muted)
			summary.Children(o.Summary)
		}
		content := ui.Row(c).AlignItems(ui.Center).Gap(scale(c, theme.Space4)).Shrink(0)
		if o.Summary != nil && !*stacked {
			content.MarginStart(ui.Auto)
		}
		build := func() {
			nav(previous, lucide.ChevronLeft(), *page-1, *page > 1, false)
			for _, n := range pageSlots(*page, pages, o.Siblings, o.Boundaries) {
				if n == 0 {
					ellipsis := ui.Box(c).Size(scale(c, 28), size).Center().Shrink(0).
						Role(ui.RoleText).Label(locale.M().MorePages).TextColor(t.Foreground)
					ellipsis.Children(func() {
						ui.Icon(c, lucide.MoreHorizontal()).Size(scale(c, 16), scale(c, 16))
					})
					continue
				}
				number := n
				active := number == *page
				label := strconv.Itoa(number)
				if active {
					label = label + ", " + locale.M().Selected
				}
				btn := ui.ButtonBase(c).Size(size, size).Center().Shrink(0).Radius(radius).
					FocusRing(true).Label(label)
				bg := ui.Transparent
				hover := t.Accent
				if active {
					bg, hover = t.Popover, t.Accent.Alpha(0.5)
					if c.Theme().Dark {
						bg, hover = t.Input.Alpha(0.32), t.Input.Alpha(0.64)
					}
					btn.Border(hairline, t.Input)
					if !btn.Pressed() {
						btn = style.Shadow(btn, c, theme.FieldShadows(c.Theme().Dark))
					}
				}
				btn = style.Background(btn, radius, bg, hover, hover)
				btn.TextColor(t.DefaultForeground)
				if !active {
					btn.OnClick(func() {
						*page = number
						if o.OnChange != nil {
							o.OnChange(number)
						}
					})
				}
				btn.Children(func() {
					style.Text(ui.Text(c, strconv.Itoa(number)), c, font).FontWeight(theme.WeightMedium)
				})
			}
			nav(nextLabel, lucide.ChevronRight(), *page+1, *page < pages, true)
		}
		content.Children(build)
	})
	if o.Summary != nil {
		root.DrawOver(func(_ *ui.Painter, r ui.Rect) {
			// Upstream's sm breakpoint, in window DIPs: breakpoints do not
			// follow the text scale.
			narrow := r.W < 640
			if narrow != *stacked {
				*stacked = narrow
				c.Invalidate()
			}
		})
	}
	return root
}

// pageSlots returns the pages to show in order, with 0 for an ellipsis. The
// three ranges ascend, so they are merged by skipping what an earlier range
// already held: no map and no sort per build.
func pageSlots(page, pages, siblings, boundaries int) []int {
	if pages <= 0 {
		return nil
	}
	if siblings < 0 {
		siblings = 0
	}
	if boundaries <= 0 {
		boundaries = 1
	}
	if page < 1 {
		page = 1
	}
	if page > pages {
		page = pages
	}
	out := make([]int, 0, 2*boundaries+2*siblings+3)
	add := func(n int) {
		if n < 1 || n > pages {
			return
		}
		if len(out) > 0 {
			last := out[len(out)-1]
			if n <= last {
				return
			}
			if n-last > 1 {
				out = append(out, 0)
			}
		}
		out = append(out, n)
	}
	for n := 1; n <= boundaries; n++ {
		add(n)
	}
	for n := page - siblings; n <= page+siblings; n++ {
		add(n)
	}
	for n := pages - boundaries + 1; n <= pages; n++ {
		add(n)
	}
	return out
}
