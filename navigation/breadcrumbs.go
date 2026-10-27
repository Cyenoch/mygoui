package navigation

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// BreadcrumbOptions configures one breadcrumb. Current marks the page the
// breadcrumbs lead to; Breadcrumbs treats its last item as current without
// needing the flag.
type BreadcrumbOptions struct {
	Label    string
	Current  bool
	Disabled bool
	OnClick  func()
}

// BreadcrumbsOptions configures a path. MaxVisible keeps the first item, the
// last MaxVisible-2 items and an overflow button for the rest; zero shows all.
// OnSelect receives the index of the item chosen, including from overflow.
type BreadcrumbsOptions struct {
	Items      []BreadcrumbOptions
	Label      string
	MaxVisible int
	MoreLabel  string
	OnSelect   func(index int)
}

// Breadcrumb builds one breadcrumb link, or the current page as plain text
// with the current-page color. Separators belong to Breadcrumbs.
func Breadcrumb(c *ui.Context, o BreadcrumbOptions) ui.Element {
	t := theme.Colors(c)
	ts := theme.SM()
	if o.Current {
		e := ui.Text(c, o.Label)
		e = style.Text(e, c, ts).TextColor(t.Foreground).SingleLine()
		e.Role(ui.RoleText)
		return e
	}
	e := ui.ButtonBase(c).Padding(0).Gap(0).Shrink(0).FocusRing(true).Role(ui.RoleLink).Label(o.Label)
	ink := t.MutedForeground
	if !o.Disabled && (e.Hovered() || e.FocusVisible()) {
		ink = t.Foreground
	}
	e.TextColor(ink)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	e.Children(func() {
		style.Text(ui.Text(c, o.Label), c, ts).SingleLine()
	})
	return style.Disabled(e, o.Disabled)
}

// Breadcrumbs builds a path of links separated by chevrons. The last item is
// the current page and is not interactive; earlier items call OnSelect or
// their own OnClick. Assistive technology sees a named group of links.
func Breadcrumbs(c *ui.Context, o BreadcrumbsOptions) ui.Element {
	t := theme.Colors(c)
	root := ui.Row(c).AlignItems(ui.Center).MinWidth(0).Role(ui.RoleGroup)
	if o.Label != "" {
		root.Label(o.Label)
	}
	root.Gap(scale(c, 10))
	if len(o.Items) == 0 {
		return root
	}
	root.Children(func() {
		more := o.MoreLabel
		if more == "" {
			more = locale.M().MorePages
		}
		head, tail := 0, len(o.Items)
		if o.MaxVisible >= 3 && len(o.Items) > o.MaxVisible {
			tail = o.MaxVisible - 2
			head = 1
		}
		var slots []func()
		add := func(i int) {
			it := o.Items[i]
			it.Current = it.Current || i == len(o.Items)-1
			slots = append(slots, func() {
				Breadcrumb(c, BreadcrumbOptions{
					Label:    it.Label,
					Current:  it.Current,
					Disabled: it.Disabled,
					OnClick: func() {
						if it.OnClick != nil {
							it.OnClick()
						}
						if o.OnSelect != nil {
							o.OnSelect(i)
						}
					},
				})
			})
		}
		for i := range head {
			add(i)
		}
		if tail < len(o.Items)-head {
			slots = append(slots, func() {
				btn := ui.ButtonBase(c).Size(scale(c, theme.Space24), scale(c, theme.Space24)).Center().
					Radius(theme.RadiusScale().MD).FocusRing(true).Label(more).
					Background(ui.Transparent).TextColor(t.Muted)
				// Snapshot the button's states in the build, as above.
				hovered, pressed := btn.Hovered(), btn.Pressed()
				btn.Draw(func(p *ui.Painter, r ui.Rect) {
					if hovered || pressed {
						p.Fill(r, t.DefaultHover, theme.RadiusScale().MD)
					}
				})
				btn.Children(func() {
					ui.Icon(c, lucide.Ellipsis()).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).Shrink(0)
				})
				btn.Menu(func(m *ui.Menu) {
					for i := head; i < len(o.Items)-tail; i++ {
						it := o.Items[i]
						mi := m.Item(it.Label).Disabled(it.Disabled)
						if mi.Chosen() {
							if it.OnClick != nil {
								it.OnClick()
							}
							if o.OnSelect != nil {
								o.OnSelect(i)
							}
						}
					}
				})
			})
		}
		for i := len(o.Items) - tail; i < len(o.Items); i++ {
			if i < head {
				continue
			}
			add(i)
		}
		for j := range slots {
			item := ui.Row(c).AlignItems(ui.Center).Gap(scale(c, theme.Space6)).Shrink(0)
			last := j == len(slots)-1
			item.Children(func() {
				slots[j]()
				if last {
					return
				}
				sep := ui.Icon(c, lucide.ChevronRight()).Size(scale(c, theme.IconMD), scale(c, theme.IconMD)).
					Shrink(0).TextColor(t.MutedForeground.Alpha(0.8))
				if c.Direction() == ui.RightToLeft {
					sep.Rotate(180)
				}
			})
		}
	})
	return root
}
