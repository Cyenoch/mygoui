package navigation

import (
	"slices"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PaginationSize is HeroUI's pagination scale.
type PaginationSize uint8

const (
	PaginationMedium PaginationSize = iota
	PaginationSmall
	PaginationLarge
)

// PaginationOptions configures page navigation. Page is 1-based and
// caller-owned; nil keeps it internally. Siblings sets the pages shown around
// the current one (default 1) and Boundaries the pages at each end (default
// 1). PreviousLabel and NextLabel override the locale's names.
type PaginationOptions struct {
	Page          *int
	Pages         int
	Size          PaginationSize
	Siblings      int
	Boundaries    int
	Label         string
	PreviousLabel string
	NextLabel     string
	OnChange      func(page int)
}

// Pagination builds previous and next controls around numbered pages with
// ellipses for gaps. Previous is disabled on the first page and Next on the
// last; clicking the current page does nothing.
func Pagination(c *ui.Context, o PaginationOptions) ui.Element {
	t := theme.Colors(c)
	pages := max(o.Pages, 1)
	root := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(scale(c, theme.Space16)).
		Role(ui.RoleGroup).MinWidth(0)
	if o.Label != "" {
		root.Label(o.Label)
	}
	page := o.Page
	if page == nil {
		page = ui.Local(root, "pagination-page", func() int { return 1 })
	}
	*page = max(1, min(*page, pages))

	size, font, navPad, radius := float32(32), theme.SM(), float32(10), theme.RadiusScale().XXXL
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

	nav := func(label string, icon *ui.SVG, to int, enabled bool) {
		btn := ui.ButtonBase(c).Height(size).MinWidth(size).PaddingX(scale(c, navPad)).
			Center().Gap(scale(c, theme.Space6)).Shrink(0).Radius(radius).FocusRing(true).Label(label)
		btn = style.Background(btn, radius, ui.Transparent, t.DefaultHover, t.DefaultHover)
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
			if icon != nil {
				ui.Icon(c, icon).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
			}
			style.Text(ui.Text(c, label), c, font)
		})
		style.Disabled(btn, !enabled)
	}

	root.Children(func() {
		content := ui.Row(c).AlignItems(ui.Center).Gap(scale(c, theme.Space4)).Shrink(0)
		build := func() {
			nav(previous, lucide.ChevronLeft(), *page-1, *page > 1)
			for _, n := range pageSlots(*page, pages, o.Siblings, o.Boundaries) {
				if n == 0 {
					ellipsis := ui.Box(c).Size(size, size).Center().Shrink(0).Role(ui.RoleNone)
					ellipsis.TextColor(t.Muted)
					ellipsis.Children(func() {
						style.Text(ui.Text(c, "\u2026"), c, font).TextColor(t.Muted)
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
				if active {
					bg = t.Default
				}
				btn = style.Background(btn, radius, bg, t.DefaultHover, t.DefaultHover)
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
			nav(nextLabel, lucide.ChevronRight(), *page+1, *page < pages)
		}
		content.Children(build)
	})
	return root
}

// pageSlots returns the pages to show in order, with 0 for an ellipsis.
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
	seen := make(map[int]bool, pages)
	add := func(n int) {
		if n >= 1 && n <= pages {
			seen[n] = true
		}
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
	out := make([]int, 0, len(seen)+2)
	prev := 0
	for _, n := range sortedKeys(seen) {
		if prev != 0 && n-prev > 1 {
			out = append(out, 0)
		}
		out = append(out, n)
		prev = n
	}
	return out
}

func sortedKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
