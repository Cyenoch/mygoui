// Package formstyle owns HeroUI's form-field surface and typography. Public
// form and date components share it so native controls paint one field chrome.
package formstyle

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// RingMode selects when a field shell paints its focus ring. HeroUI uses
// :focus for inputs and :focus-visible for button-like triggers.
type RingMode uint8

const (
	// RingOnFocus rings whenever the shell or the editor inside it has the
	// focus, as input CSS's :focus does. MyGo reports the editor's focus as
	// the shell's FocusWithin.
	RingOnFocus RingMode = iota
	// RingOnKeyboardFocus rings only for keyboard focus, as a button
	// trigger's :focus-visible does. It applies when the shell itself is the
	// focusable control.
	RingOnKeyboardFocus
)

// Surface paints HeroUI's field chrome with a :focus ring, the mode an editor
// and a color field use. Button-like triggers use SurfaceMode with
// RingOnKeyboardFocus.
func Surface(e ui.Element, c *ui.Context, secondary, invalid, disabled bool) ui.Element {
	return SurfaceMode(e, c, RingOnFocus, secondary, invalid, disabled)
}

// SurfaceMode paints HeroUI's field background, border and hover/focus/invalid
// states around a native editor, and owns the field's focus ring.
//
// HeroUI's field ring is its own token (focus-field-ring): two DIPs against
// the box, not the offset ring `FocusRing(true)` draws. MyGo never rings an
// editor itself, so the shell paints the ring for the editor inside it, on
// the focus condition the caller names: :focus for editors and :focus-visible
// for button-like triggers.
func SurfaceMode(e ui.Element, c *ui.Context, ring RingMode, secondary, invalid, disabled bool) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().Field
	base, hover, focus := t.FieldBackground, t.FieldHover, t.FieldFocus
	if secondary {
		// input--secondary: a default fill at rest and under the pointer, and
		// back to default while focused, with no elevation.
		base, hover, focus = t.Default, t.DefaultHover, t.Default
	}
	// Snapshot the states while the view builds: a retained copy of a
	// leaving element paints in later passes, where an Element captured
	// here has expired. A disabled field paints its resting face and no
	// ring: upstream cannot focus one, while MyGo keeps FocusWithin until
	// the focus moves elsewhere.
	var focused, hovered bool
	if !disabled {
		if ring == RingOnKeyboardFocus {
			focused = e.FocusVisible()
		} else {
			focused = e.FocusWithin()
		}
		hovered = e.Hovered()
	}
	e = e.Radius(radius).Clip()
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		fill := base
		switch {
		case focused:
			fill = focus
		case hovered:
			fill = hover
		}
		p.Fill(r, fill, radius)
		// Upstream fields have no border width by default, so only an
		// injected FieldBorder (or the invalid outline) paints one.
		border := t.FieldBorder
		switch {
		case invalid:
			border = t.Danger
		case focused && border.A != 0:
			border = t.FieldBorderFocus
		case hovered && border.A != 0:
			border = t.FieldBorderHover
		}
		if border.A != 0 {
			p.Stroke(r, border, radius, theme.Hairline)
		}
	})
	if !secondary {
		e = style.Shadow(e, c, theme.FieldShadows(c.Theme().Dark))
	}
	if focused {
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			stroke := t.Focus
			if invalid {
				stroke = t.Danger
			}
			FieldRing(p, r, radius, stroke)
		})
	}
	return style.Disabled(e, disabled)
}

// FieldRing paints HeroUI's two-pixel field focus ring around r: it sits
// against the box, as focus-field-ring's ring-offset-0 does, and is painted
// inside the ring's own rect so no stroke spills past its radius. The Painter
// has no ring-color argument, so a danger ring uses the same geometry.
func FieldRing(p *ui.Painter, r ui.Rect, radius float32, color ui.Color) {
	const width = 2
	box := ui.Rect{X: r.X - width, Y: r.Y - width, W: r.W + 2*width, H: r.H + 2*width}
	p.Stroke(box, color, radius+width, width)
}

// Text applies HeroUI's 14/20 field typography and foreground.
func Text(e ui.Element, c *ui.Context) ui.Element {
	return style.Text(e, c, theme.SM()).TextColor(theme.Colors(c).FieldForeground)
}
