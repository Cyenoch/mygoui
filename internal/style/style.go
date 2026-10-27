// Package style centralizes typography and state styling shared by components.
// It is internal so the public surface stays component-oriented.
package style

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Text applies a theme type role to an element. Sizes, line boxes and letter
// spacing all follow the host text scale, which is read once per call.
func Text(e ui.Element, c *ui.Context, s theme.TextStyle) ui.Element {
	family := theme.FontFamily(c)
	if s.Mono {
		family = theme.MonoFontFamily()
	} else if s.Heading {
		family = theme.HeadingFontFamily(c)
	}
	scale := theme.TextScale(c)
	e = e.Font(family).
		FontSize(s.Size * scale).
		FixedLineHeight(s.LineHeight * scale).
		FontWeight(s.Weight)
	if s.Tracking != 0 {
		e = e.LetterSpacing(s.Tracking * s.Size * scale)
	}
	return e
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

// Disabled applies MyGo's disabled behavior: 50% opacity and no activation.
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
// painting. Apply Disabled before Background: the face is chosen when this
// runs, and a disabled element reports neither hover nor press.
func Background(e ui.Element, radius float32, normal, hover, pressed ui.Color) ui.Element {
	fill := normal
	switch {
	case e.Pressed() && pressed.A != 0:
		fill = pressed
	case e.Hovered() && hover.A != 0:
		fill = hover
	}
	return e.Radius(radius).Clip().Background(fill)
}

// CornerFill uses the native per-corner background rather than a square
// child fill that depends on ancestor clipping. The retained exit copy keeps
// its face without resolving a handle in a later build.
func CornerFill(c *ui.Context, e ui.Element, gap float32, corners [4]float32,
	normal, hover, pressed ui.Color, content, overlay func()) ui.Element {
	e = Background(e, 0, normal, hover, pressed).
		Radius(corners[0], corners[1], corners[2], corners[3]).FocusRing(false)
	if e.FocusVisible() {
		e.Draw(func(p *ui.Painter, r ui.Rect) { insetRing(p, r, corners) })
	}
	e.Children(func() {
		layer := ui.Row(c).Gap(gap).Center().Shrink(0)
		layer.Children(content)
		if overlay != nil {
			overlay()
		}
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
