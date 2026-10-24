// Package formstyle owns HeroUI's form-field surface and typography. Public
// form and date components share it so native controls paint one field chrome.
package formstyle

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Surface paints HeroUI's field background, border and hover/focus/invalid
// states around a native editor. The editor keeps its own focus ring.
func Surface(e ui.Element, c *ui.Context, secondary, invalid, disabled bool) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().Field
	base, hover := t.FieldBackground, t.FieldHover
	if secondary {
		base, hover = t.Default, t.DefaultHover
	}
	// Snapshot the states while the view builds: a retained copy of a
	// leaving element paints in later passes, where an Element captured
	// here has expired.
	focused, hovered := e.FocusWithin(), e.Hovered()
	e = e.Radius(radius).Clip()
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		fill := base
		switch {
		case focused:
			fill = t.FieldFocus
		case hovered:
			fill = hover
		}
		p.Fill(r, fill, radius)
		border := t.FieldBorder
		switch {
		case invalid:
			border = t.Danger
		case focused:
			border = t.FieldBorderFocus
		case hovered:
			border = t.FieldBorderHover
		}
		if border.A != 0 {
			p.Stroke(r, border, radius, 1)
		}
	})
	if !secondary {
		e = style.Shadow(e, c, theme.FieldShadows(c.Theme().Dark))
	}
	if focused {
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			ring := t.Focus
			if invalid {
				ring = t.Danger
			}
			p.Stroke(r, ring, radius, 2)
		})
	}
	return style.Disabled(e, disabled)
}

// Text applies HeroUI's 14/20 field typography and foreground.
func Text(e ui.Element, c *ui.Context) ui.Element {
	return style.Text(e, c, theme.SM()).TextColor(theme.Colors(c).FieldForeground)
}
