// Package formstyle shares Coss field chrome between native form and date controls.
package formstyle

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// RingMode selects the focus condition of a field shell.
type RingMode uint8

const (
	RingOnFocus RingMode = iota
	RingOnKeyboardFocus
)

// Surface styles a native editor's shell.
func Surface(e ui.Element, c *ui.Context, secondary, invalid, disabled bool) ui.Element {
	return SurfaceMode(e, c, RingOnFocus, secondary, invalid, disabled)
}

// SurfaceMode paints Coss's bordered input and focus or invalid state.
// Secondary is retained for existing callers; Coss uses the same field face.
func SurfaceMode(e ui.Element, c *ui.Context, ring RingMode, secondary, invalid, disabled bool) ui.Element {
	focused := false
	if !disabled {
		if ring == RingOnKeyboardFocus {
			focused = e.FocusVisible()
		} else {
			focused = e.FocusWithin()
		}
	}
	return SurfaceFocused(e, c, focused, invalid, disabled)
}

// Face holds the coss field colors for one state. Ring is transparent
// unless the field is focused.
type Face struct{ Fill, Border, Ring ui.Color }

// FieldFace resolves the coss field face for the current appearance.
func FieldFace(c *ui.Context, focused, invalid bool) Face {
	t := theme.Colors(c)
	dark := c.Theme().Dark
	f := Face{Fill: t.Background, Border: t.Input}
	if dark {
		f.Fill = t.Input.Alpha(0.32)
	}
	if invalid {
		f.Border = t.Destructive.Alpha(0.36)
	}
	if focused {
		f.Border, f.Ring = t.Ring, t.Ring.Alpha(0.24)
		if invalid {
			f.Border, f.Ring = t.Destructive.Alpha(0.64), t.Destructive.Alpha(0.16)
			if dark {
				f.Ring = t.Destructive.Alpha(0.24)
			}
		}
	}
	return f
}

// SurfaceFocused styles a compound shell whose editor owns the focus condition.
func SurfaceFocused(e ui.Element, c *ui.Context, focused, invalid, disabled bool) ui.Element {
	radius := theme.RadiusScale().Field
	focused = focused && !disabled
	face := FieldFace(c, focused, invalid)
	e = e.Radius(radius).Clip()
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, face.Fill, radius)
		p.Stroke(r, face.Border, radius, theme.Hairline)
	})
	if focused {
		e.DrawOver(func(p *ui.Painter, r ui.Rect) { FieldRing(p, r, radius, face.Ring) })
	} else if !disabled && !invalid {
		e.Shadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.05))
	}
	return style.Disabled(e, disabled)
}

// FieldRing paints the three-DIP translucent focus ring outside a field.
func FieldRing(p *ui.Painter, r ui.Rect, radius float32, color ui.Color) {
	const width = 3
	box := ui.Rect{X: r.X - width, Y: r.Y - width, W: r.W + 2*width, H: r.H + 2*width}
	p.Stroke(box, color, radius+width, width)
}

// Text applies Coss desktop field typography.
func Text(e ui.Element, c *ui.Context) ui.Element {
	return style.Text(e, c, theme.SM()).TextColor(theme.Colors(c).Foreground)
}

// ErrorText is a field's visible error: the caller's text, or the localized
// invalid message when the field is invalid without one.
func ErrorText(invalid bool, text string) string {
	if text == "" && invalid {
		return locale.M().Invalid
	}
	return text
}
