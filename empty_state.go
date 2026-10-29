package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// EmptyStateOptions configures a quiet empty result.
type EmptyStateOptions struct {
	Title, Description string
	Icon               *ui.SVG
	Action             func()
}

// EmptyState centers an icon, title, description and optional recovery action.
// With no content it shows the localized empty placeholder.
func EmptyState(c *ui.Context, o EmptyStateOptions) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).Gap(theme.Space24).Padding(theme.Space48, theme.Space24).
		MinWidth(0).FillWidth().AlignItems(ui.Center).Justify(ui.Center)
	e.Children(func() {
		if o.Icon == nil && o.Title == "" && o.Description == "" && o.Action == nil {
			style.Text(ui.Text(c, locale.M().Empty), c, theme.SM()).TextColor(t.Muted)
			return
		}
		if o.Icon != nil || o.Title != "" || o.Description != "" {
			ui.Column(c).MaxWidth(384).AlignItems(ui.Center).Children(func() {
				if o.Icon != nil {
					icon := ui.Box(c).Size(36, 36).
						Radius(theme.RadiusScale().MD).Border(theme.Hairline, t.Border).
						Background(t.Card).Center()
					if o.Title != "" || o.Description != "" {
						icon.Margin(0, 0, theme.Space24, 0)
					}
					style.Shadow(icon, c, theme.SurfaceShadows(c.Theme().Dark)).Children(func() {
						ui.Icon(c, o.Icon).Size(18, 18).TextColor(t.Foreground)
					})
				}
				if o.Title != "" {
					ts := theme.XL()
					ts.Heading, ts.Weight = true, theme.WeightBold
					style.Text(ui.Text(c, o.Title), c, ts).TextColor(t.Foreground).TextAlign(ui.Center)
				}
				if o.Description != "" {
					description := style.Text(ui.Text(c, o.Description), c, theme.SM()).
						TextColor(t.MutedForeground).TextAlign(ui.Center)
					if o.Title != "" {
						description.Margin(theme.Space4, 0, 0, 0)
					}
				}
			})
		}
		if o.Action != nil {
			ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Space8).
				Justify(ui.Center).Children(o.Action)
		}
	})
	return e
}

// HeaderOptions configures a section header.
type HeaderOptions struct {
	Text string
}

// Header is a small muted section heading used by menus and lists.
func Header(c *ui.Context, o HeaderOptions) ui.Element {
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	return style.Text(ui.Text(c, o.Text), c, ts).
		TextColor(theme.Colors(c).Muted).
		Padding(theme.Space6, theme.Space8, theme.Space4, theme.Space8).
		FillWidth()
}
