package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// EmptyStateOptions configures a quiet empty result.
type EmptyStateOptions struct {
	Title, Description string
	Icon               *ui.SVG
	Action             func()
}

// EmptyState is HeroUI's low-emphasis "no content" block.
func EmptyState(c *ui.Context, o EmptyStateOptions) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).Gap(theme.Space4).Padding(theme.Space8).MinWidth(0)
	e.Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconLG), theme.Scale(c, theme.IconLG)).
				TextColor(t.Muted)
		}
		if o.Title != "" {
			style.Text(ui.Text(c, o.Title), c, theme.SM()).TextColor(t.Foreground)
		}
		if o.Description != "" {
			style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.Muted)
		}
		if o.Action != nil {
			o.Action()
		}
	})
	return e
}

// HeaderOptions configures a section header.
type HeaderOptions struct {
	Text string
}

// Header is HeroUI's small muted section heading used by menus and lists.
func Header(c *ui.Context, o HeaderOptions) ui.Element {
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	return style.Text(ui.Text(c, o.Text), c, ts).
		TextColor(theme.Colors(c).Muted).
		Padding(theme.Space6, theme.Space8, theme.Space4, theme.Space8).
		FillWidth()
}
