package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// GroupOrientation is the layout direction of a choice group. The zero value
// is vertical.
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

// Checkbox is a focusable row whose box, label and
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

	return choiceField(c, o.Description, errText, disabled, theme.Space24, func() {
		row := ui.CheckboxBase(c, checked).Gap(theme.Space8).Shrink(0).
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

// checkboxControl paints the 16 DIP Coss checkbox around the native row.
func checkboxControl(c *ui.Context, checked *bool, indeterminate *bool, invalid bool, variant ControlVariant, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	size := theme.Scale(c, 16)
	radius := theme.Scale(c, 4)
	focused := owner.FocusVisible()
	mixed := indeterminate != nil && *indeterminate
	selected := *checked && !mixed
	fill, border, ink := t.Background, t.Input, t.Foreground
	if c.Theme().Dark {
		fill = t.Input.Alpha(0.32)
	}
	if selected {
		fill, border, ink = t.Primary, t.Primary, t.PrimaryForeground
	} else if invalid {
		border = t.Destructive.Alpha(0.36)
	}
	box := ui.Box(c).Size(size, size).Shrink(0).Center()
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, fill, radius)
		p.Stroke(r, border, radius, theme.Hairline)
		if focused {
			choiceFocusRing(p, r, radius, invalid, t)
		}
		if mixed {
			w, h := r.W*0.56, theme.Scale(c, 1.5)
			p.Fill(ui.Rect{X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2, W: w, H: h}, ink, h/2)
		} else if selected {
			mark := r.W * 0.65
			p.StrokePath(choiceCheckPath(ui.Rect{
				X: r.X + (r.W-mark)/2, Y: r.Y + (r.H-mark)/2,
				W: mark, H: mark,
			}), theme.Scale(c, 1.8), ink)
		}
	})
	return box
}

// choiceCheckPath fits the check mark into its box.
func choiceCheckPath(r ui.Rect) *ui.Path {
	var p ui.Path
	return p.MoveTo(r.X+r.W*0.084, r.Y+r.H*0.5).
		LineTo(r.X+r.W*0.417, r.Y+r.H*0.778).
		LineTo(r.X+r.W*0.861, r.Y+r.H*0.222)
}

// choiceFocusRing paints the focus ring, in a destructive color for an invalid
// control, outside the control's box.
func choiceFocusRing(p *ui.Painter, r ui.Rect, radius float32, invalid bool, t theme.Palette) {
	if invalid {
		p.Stroke(choiceExpand(r, 2), t.Destructive.Alpha(0.48), radius+2, 2)
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

// choiceField lays out a choice row built by body with field-level
// description and error, indented under the label.
func choiceField(c *ui.Context, description, errText string, disabled bool, indent float32, body func()) ui.Element {
	outer := ui.Column(c).MinWidth(0).Gap(theme.Space4).Children(func() {
		body()
		if description != "" {
			Description(c, description).PaddingStart(indent)
		}
		if errText != "" {
			choiceErrorMessage(c, errText).PaddingStart(indent)
		}
	})
	return style.Disabled(outer, disabled)
}

// choiceErrorMessage renders the field-level error of a choice control.
func choiceErrorMessage(c *ui.Context, text string) ui.Element {
	return style.Text(ui.Text(c, text), c, theme.HelpText()).
		TextColor(theme.Colors(c).DestructiveForeground).Role(ui.RoleStatus)
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
