package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// GroupOrientation is the layout direction of a choice group. The zero value
// is HeroUI's vertical default.
type GroupOrientation uint8

const (
	GroupVertical GroupOrientation = iota
	GroupHorizontal
)

// CheckboxOptions configures a native check box. Checked is caller-owned: a
// click or Space flips it and OnChange reports the new value. A nil Checked
// renders a disabled check box, since there is nowhere to store selection.
//
// Indeterminate paints the mixed mark and tells assistive technology the
// check box is mixed. Activating it clears the flag and checks the box.
type CheckboxOptions struct {
	Checked       *bool
	Indeterminate *bool
	Label         string
	Description   string
	Error         string
	Required      bool
	Invalid       bool
	Disabled      bool
	Variant       ControlVariant
	OnChange      func(bool)
	Handle        *ui.Handle
}

// Checkbox is HeroUI's check box: a focusable row whose box, label and
// description follow the caller's value. The whole row activates the native
// check box, including Space while it has the focus.
func Checkbox(c *ui.Context, o CheckboxOptions) ui.Element {
	checked := o.Checked
	if checked == nil {
		var off bool
		checked = &off
		o.Disabled = true
	}
	t := theme.Colors(c)
	disabled := o.Disabled
	errText := choiceErrorText(o.Invalid, o.Error)

	return choiceField(c, o.Description, errText, disabled, theme.Space(7), func() {
		row := ui.CheckboxBase(c, checked).Gap(theme.Space12).Shrink(0).
			FocusRing(false).TextColor(t.Foreground).Label(o.Label)
		if o.Handle != nil {
			row.Bind(o.Handle)
		}
		if o.Description != "" {
			row.Description(o.Description)
		}
		if errText != "" {
			row.Error(errText)
		}
		if o.Indeterminate != nil && *o.Indeterminate {
			row.Mixed()
		}
		if o.Checked != nil && !disabled && (o.OnChange != nil || o.Indeterminate != nil) {
			row.OnChange(func() {
				if o.Indeterminate != nil {
					*o.Indeterminate = false
				}
				if o.OnChange != nil {
					o.OnChange(*checked)
				}
			})
		}

		row.Children(func() {
			checkboxControl(c, checked, o.Indeterminate, o.Invalid, o.Variant, row)
			if o.Label != "" {
				Label(c, o.Label, LabelOptions{
					Required: o.Required, Invalid: o.Invalid,
				})
			}
		})
	})
}

// checkboxControl draws HeroUI's 16 DIP square and its check or mixed mark.
// owner supplies the live hover, pressed and focus state.
func checkboxControl(c *ui.Context, checked *bool, indeterminate *bool, invalid bool, variant ControlVariant, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	size := float32(theme.Space16)
	radius := theme.RadiusScale().MD
	stroke := theme.Scale(c, 2.5)
	box := ui.Box(c).Size(size, size).Shrink(0)
	// Snapshot the row's states in the build: a retained copy of a leaving
	// element paints in later passes, where owner has expired.
	hovered, pressed, focused := owner.Hovered(), owner.Pressed(), owner.FocusVisible()
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		selected := *checked
		mixed := indeterminate != nil && *indeterminate
		fill := t.FieldBackground
		border := t.FieldBorder
		if variant == ControlSecondary && !selected && !mixed {
			fill = t.Default
		}
		switch {
		case selected || mixed:
			border = ui.Transparent
			switch {
			case invalid:
				fill = t.Danger
			case mixed:
				fill = t.Accent
				if pressed {
					fill = t.AccentHover
				}
			default:
				fill = t.Accent
				if hovered || pressed {
					fill = t.AccentHover
				}
			}
		case invalid:
			border = t.Danger
		case hovered:
			border = t.FieldBorderHover
		}
		if variant != ControlSecondary {
			p.Shadow(r, radius, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.08))
		}
		p.Fill(r, fill, radius)
		if border.A != 0 {
			p.Stroke(r, border, radius, theme.Hairline)
		}
		if selected || mixed {
			fg := t.AccentForeground
			if invalid {
				fg = t.DangerForeground
			}
			if mixed {
				h := float32(1.5)
				p.Fill(ui.Rect{X: r.X + r.W*0.25, Y: r.Y + (r.H-h)/2, W: r.W * 0.5, H: h}, fg, h/2)
			} else {
				p.StrokePath(choiceCheckPath(r), stroke, fg)
			}
		}
		if focused {
			choiceFocusRing(p, r, radius, invalid, t)
		}
	})
	return box
}

// choiceCheckPath is the upstream check mark, within a 16 DIP box.
func choiceCheckPath(r ui.Rect) *ui.Path {
	x, y, s := r.X+r.W*0.19, r.Y+r.H*0.19, r.W*0.62
	var p ui.Path
	return p.MoveTo(x+s*0.22, y+s*0.52).
		LineTo(x+s*0.42, y+s*0.71).
		LineTo(x+s*0.78, y+s*0.31)
}

// choiceFocusRing paints HeroUI's two-pixel ring, in Danger for an invalid
// control, outside the control's box.
func choiceFocusRing(p *ui.Painter, r ui.Rect, radius float32, invalid bool, t theme.Palette) {
	if invalid {
		p.Stroke(choiceExpand(r, 2), t.Danger, radius+2, 2)
		return
	}
	p.FocusRing(r, choiceRadii(radius))
}

// choiceRadii returns the uniform radius Painter.FocusRing takes.
func choiceRadii(radius float32) [4]float32 {
	return [4]float32{radius, radius, radius, radius}
}

// choiceExpand grows r by by on every side.
func choiceExpand(r ui.Rect, by float32) ui.Rect {
	return ui.Rect{X: r.X - by, Y: r.Y - by, W: r.W + 2*by, H: r.H + 2*by}
}

// choiceField lays out a choice row built by body with HeroUI's field-level
// description and error, indented under the label.
func choiceField(c *ui.Context, description, errText string, disabled bool, indent float32, body func()) ui.Element {
	outer := ui.Column(c).MinWidth(0).Gap(theme.Space4).Children(func() {
		body()
		if description != "" {
			Description(c, description).PaddingStart(indent)
		}
		if errText != "" {
			FieldError(c, errText).PaddingStart(indent)
		}
	})
	return style.Disabled(outer, disabled)
}

// choiceErrorText resolves a caller's error text, falling back to the
// library's localized invalid message.
func choiceErrorText(invalid bool, text string) string {
	if text != "" {
		return text
	}
	if invalid {
		return locale.M().Invalid
	}
	return ""
}
