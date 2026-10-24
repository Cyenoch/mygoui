// Package style centralizes HeroUI's paint-time typography and state styling.
// It is internal so the public surface stays component-oriented.
package style

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Text applies a HeroUI type role to an element.
func Text(e ui.Element, c *ui.Context, s theme.TextStyle) ui.Element {
	family := theme.FontFamily(c)
	if s.Mono {
		family = theme.Mono
	}
	return e.Font(family).
		FontSize(theme.Scale(c, s.Size)).
		FixedLineHeight(theme.Scale(c, s.LineHeight)).
		FontWeight(s.Weight)
}

// Shadow applies zero or more native drop-shadow layers.
func Shadow(e ui.Element, c *ui.Context, shadows []theme.Shadow) ui.Element {
	for _, s := range shadows {
		if s.Inset {
			e = e.InsetShadow(s.X, s.Y, s.Blur, s.Spread, s.Color)
		} else {
			e = e.Shadow(s.X, s.Y, s.Blur, s.Spread, s.Color)
		}
	}
	return e
}

// Disabled applies HeroUI's disabled contract: 50% opacity and no activation.
func Disabled(e ui.Element, disabled bool) ui.Element {
	if !disabled {
		return e
	}
	return e.Disabled(true)
}

// Background paints a rounded fill before children. It reads hover and
// pressed while the view builds and keeps the chosen face as a value: MyGo
// paints a retained copy of an element leaving through an
// ElementTransition in later build passes, where an Element bound in the
// build that made the copy has expired. A frame is rebuilt for each hover
// or press change, so the face still follows the pointer. Zero alpha skips
// painting.
func Background(e ui.Element, radius float32, normal, hover, pressed ui.Color) ui.Element {
	fill := normal
	switch {
	case e.Pressed() && pressed.A != 0:
		fill = pressed
	case e.Hovered() && hover.A != 0:
		fill = hover
	}
	e = e.Radius(radius).Clip()
	if fill.A == 0 {
		return e
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, fill, radius)
	})
}

// CornerFill paints the same faces as Background on an element whose
// corners differ, as the outer ends of a joined group do. MyGo's Painter
// fills only with one radius, and an element's Draw runs before its own
// clip, so a square Fill in Draw spills past the rounded corners. The
// state face therefore paints in a single child layer, which the element's
// per-corner clip shapes exactly; the layer is the element's only child and
// hugs the content, so the element's padding and Justify place it as the
// uniform-radius path places its children.
//
// The focus ring is inset, as HeroUI draws it for joined controls: MyGo's
// own ring paints outside the box, where the next segment of a joined group
// paints over it.
func CornerFill(c *ui.Context, e ui.Element, gap float32, corners [4]float32,
	normal, hover, pressed ui.Color, content func()) ui.Element {
	e = e.Radius(corners[0], corners[1], corners[2], corners[3]).Clip().FocusRing(false)
	if normal.A != 0 {
		e = e.Background(normal)
	}
	var owner ui.Handle
	e = e.Bind(&owner)
	e.Children(func() {
		layer := ui.Row(c).Gap(gap).Center().Shrink(0)
		layer.Draw(func(p *ui.Painter, _ ui.Rect) {
			// Resolve the face's owner in the pass painting it: a retained
			// copy of a leaving element paints in a later pass, where the
			// Element captured at build time has expired.
			el := c.Resolve(owner)
			if !el.Valid() {
				return
			}
			fill := ui.Color{}
			switch {
			case el.Pressed() && pressed.A != 0:
				fill = pressed
			case el.Hovered() && hover.A != 0:
				fill = hover
			}
			box := el.LayoutBounds()
			if fill.A != 0 {
				p.Fill(box, fill, 0)
			}
			if el.FocusVisible() {
				insetRing(p, box, corners)
			}
		})
		layer.Children(content)
	})
	return e
}

// insetRing draws a 2 DIP focus ring 1 DIP inside box, rounded like it.
// Square corners take the 3 DIP rounding Painter.FocusRing adds to the radii
// it is given, the least it can draw.
func insetRing(p *ui.Painter, box ui.Rect, corners [4]float32) {
	const inset = 4 // FocusRing's ring starts 3 DIP outside the rect it gets
	if box.W < 2*inset || box.H < 2*inset {
		return
	}
	var rad [4]float32
	for i, c := range corners {
		rad[i] = max(c-inset, 0)
	}
	p.FocusRing(ui.Rect{X: box.X + inset, Y: box.Y + inset,
		W: box.W - 2*inset, H: box.H - 2*inset}, rad)
}
