package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ResultOptions configures the outcome of an operation. Status picks the
// tint and the default icon; Icon replaces the glyph and keeps the tint.
type ResultOptions struct {
	Status             StatusColor
	Title, Description string
	Icon               *ui.SVG
	// Actions builds the buttons under the text.
	Actions func()
	// Content builds supporting detail under the actions.
	Content func()
}

// Result centers a status icon, title, description, actions and optional
// detail. Text and actions are an EmptyState; the icon carries the status.
func Result(c *ui.Context, o ResultOptions) ui.Element {
	t := theme.Colors(c)
	status := statusColors(&t, o.Status)
	icon := o.Icon
	if icon == nil {
		icon = StatusIcon(o.Status)
	}
	e := ui.Column(c).Gap(theme.Space24).Padding(theme.Space48, theme.Space24).
		MinWidth(0).FillWidth().AlignItems(ui.Center).Justify(ui.Center)
	return e.Children(func() {
		// EmptyState's icon tile, tinted with the status.
		size, glyph := theme.Scale(c, 36), theme.Scale(c, 18)
		ui.Box(c).Size(size, size).Shrink(0).Center().
			Radius(theme.RadiusScale().MD).Border(theme.Hairline, t.Border).
			Background(status.Soft.Over(t.Card)).Children(func() {
			ui.Icon(c, icon).Size(glyph, glyph).TextColor(status.Foreground)
		})
		if o.Title != "" || o.Description != "" || o.Actions != nil {
			EmptyState(c, EmptyStateOptions{
				Title: o.Title, Description: o.Description, Action: o.Actions,
			}).Padding(0)
		}
		if o.Content != nil {
			ui.Column(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Children(o.Content)
		}
	})
}
