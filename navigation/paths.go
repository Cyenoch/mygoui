package navigation

import (
	"slices"
	"strconv"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	ui "github.com/egoist/mygo/ui"
)

type BreadcrumbID string
type BreadcrumbItem struct {
	ID       BreadcrumbID
	Label    string
	OnSelect func()
}
type BreadcrumbState struct{ menu overlay.MenuState }
type BreadcrumbOptions struct {
	Label      string
	MaxVisible int
	Separator  *ui.SVG
}

// Breadcrumb makes the final item a passive current-page label. Overflow keeps
// the first ancestor and last MaxVisible-2 items, with real actions in between.
func Breadcrumb(c *ui.Context, s *BreadcrumbState, o BreadcrumbOptions, items []BreadcrumbItem) ui.Element {
	p := theme.Colors(c)
	label := o.Label
	if label == "" {
		label = locales.ComponentBreadcrumb()
	}
	separator := o.Separator
	if separator == nil {
		separator = lucide.ChevronRight()
	}
	root := ui.Row(c).Role(ui.RoleGroup).Label(label).Wrap().Gap(theme.Space6)
	root.Children(func() {
		visible := o.MaxVisible
		if visible < 3 {
			visible = 4
		}
		hiddenEnd := 0
		if len(items) > visible {
			hiddenEnd = len(items) - (visible - 2)
		}
		sep := func() { ui.Icon(c, separator).Size(theme.IconSmall, theme.IconSmall).TextColor(p.MutedForeground) }
		for i, item := range items {
			if i > 0 && i < hiddenEnd {
				if i == 1 {
					sep()
					var menu []overlay.MenuItem
					for _, ancestor := range items[1:hiddenEnd] {
						menu = append(menu, overlay.MenuItem{ID: overlay.MenuActionID(ancestor.ID), Label: ancestor.Label, Disabled: ancestor.OnSelect == nil, OnSelect: ancestor.OnSelect})
					}
					overlay.DropdownMenu(c.Key("overflow"), &s.menu, overlay.MenuOptions{Label: locales.ComponentMoreAncestors()}, menu)
				}
				continue
			}
			if i > 0 {
				sep()
			}
			if i == len(items)-1 {
				ui.Text(c.Key(item.ID), item.Label).TextColor(p.Foreground).Description(locales.ComponentCurrentPage())
			} else {
				mygoui.Button(c.Key(item.ID), mygoui.ButtonOptions{Label: item.Label, Variant: mygoui.ButtonLink, Size: mygoui.ButtonSM, OnClick: item.OnSelect}).TextColor(p.MutedForeground)
			}
		}
	})
	return root
}

type PaginationOptions struct {
	Page, Pages                 int
	Disabled, Simple, IconsOnly bool
	OnChange                    func(int)
}

// Pagination uses one-based page numbers. It never rewrites caller state while
// rendering. Page-count zero presents disabled boundaries and no phantom page.
func Pagination(c *ui.Context, o PaginationOptions) ui.Element {
	page := max(1, min(o.Page, max(1, o.Pages)))
	disabled := o.Disabled || o.Pages <= 0 || o.OnChange == nil
	root := ui.Row(c).Role(ui.RoleGroup).Label(locales.ComponentPagination()).Wrap().Gap(theme.Space4)
	root.Children(func() {
		size := mygoui.ButtonSizeDefault
		previous, next := locales.ComponentPrevious(), locales.ComponentNext()
		if o.IconsOnly {
			size = mygoui.ButtonIconSM
		}
		left, right := lucide.ChevronLeft(), lucide.ChevronRight()
		if c.Direction() == ui.RightToLeft {
			left, right = right, left
		}
		mygoui.Button(c.Key("previous"), mygoui.ButtonOptions{Label: previous, Icon: left, Size: size, Variant: mygoui.ButtonGhost, Disabled: disabled || page <= 1, OnClick: func() { o.OnChange(page - 1) }})
		if o.Simple {
			label := locales.ComponentNoPages()
			if o.Pages > 0 {
				label = locales.ComponentPageOf(page, o.Pages)
			}
			ui.Text(c, label).Description(locales.ComponentCurrentPage()).PaddingX(theme.Space8)
		} else if o.Pages > 0 {
			last := 0
			pages := [5]int{1, page - 1, page, page + 1, o.Pages}
			slices.Sort(pages[:])
			for _, n := range pages {
				if n < 1 || n > o.Pages || n == last {
					continue
				}
				if last > 0 && n-last > 1 {
					ui.Icon(c.Key("gap-"+strconv.Itoa(n)), lucide.Ellipsis()).Size(theme.IconMedium, theme.IconMedium).Description(locales.ComponentMorePages())
				}
				current := n == page
				variant := mygoui.ButtonGhost
				if current {
					variant = mygoui.ButtonOutline
				}
				b := mygoui.Button(c.Key(n), mygoui.ButtonOptions{Label: strconv.Itoa(n), Size: mygoui.ButtonSM, Variant: variant, Disabled: disabled, OnClick: func() {
					if !current {
						o.OnChange(n)
					}
				}}).Selected(current)
				if current {
					b.Description(locales.ComponentCurrentPage())
				} else {
					b.Description(locales.ComponentGoToPage(n))
				}
				last = n
			}
		}
		mygoui.Button(c.Key("next"), mygoui.ButtonOptions{Label: next, TrailingIcon: right, Size: size, Variant: mygoui.ButtonGhost, Disabled: disabled || page >= o.Pages, OnClick: func() { o.OnChange(page + 1) }})
	})
	return root
}
