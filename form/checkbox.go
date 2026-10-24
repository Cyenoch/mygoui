package form

import (
	"time"

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
	// Snapshot the row's states and the value in the build: a retained copy
	// of a leaving element paints in later passes, where owner has expired
	// and the caller may have cleared the check box.
	hovered, pressed, focused := owner.Hovered(), owner.Pressed(), owner.FocusVisible()
	selected, mixed := *checked, indeterminate != nil && *indeterminate
	on := selected || mixed
	reduceMotion := theme.ReduceMotion(c)

	box := ui.Box(c).Size(size, size).Shrink(0).Center()
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		fill := t.FieldBackground
		border := t.FieldBorder
		if variant == ControlSecondary {
			fill = t.Default
		}
		if !on && !invalid && hovered {
			border = t.FieldBorderHover
		}
		if variant != ControlSecondary {
			p.Shadow(r, radius, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.08))
		}
		p.Fill(r, fill, radius)
		if border.A != 0 {
			p.Stroke(r, border, radius, theme.Hairline)
		}
		switch {
		case focused:
			choiceFocusRing(p, r, radius, invalid, t)
		case invalid:
			choiceInvalidRing(p, r, radius, t)
		}
	})
	if on {
		// HeroUI paints the chosen face as a layer over the field and scales
		// it from 70% as the check appears. The native transition keeps the
		// last face alive on the way out, without restructuring the row.
		fill, fg := t.Accent, t.AccentForeground
		switch {
		case invalid:
			fill, fg = t.Danger, t.DangerForeground
		case hovered || pressed:
			fill = t.AccentHover
		}
		box.Children(func() {
			layer := ui.Box(c).Size(size, size).Shrink(0).Center()
			layer.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, fill, radius)
				if mixed {
					// Upstream draws a 9 DIP line in the middle of its 12 DIP
					// indicator; round caps make the mark 10.5 DIP wide.
					w, h := r.W*0.65625, r.H*0.09375
					p.Fill(ui.Rect{
						X: r.X + (r.W-w)/2, Y: r.Y + (r.H-h)/2,
						W: w, H: h,
					}, fg, h/2)
					return
				}
				mark := r.W * 0.625 // the checkmark's 10 DIP box
				p.StrokePath(choiceCheckPath(ui.Rect{
					X: r.X + (r.W-mark)/2, Y: r.Y + (r.H-mark)/2,
					W: mark, H: mark,
				}), stroke, fg)
			})
			if !reduceMotion {
				layer.Transition(ui.ElementTransition{
					Duration: 150 * time.Millisecond,
					Ease:     ui.EaseOut,
					Colors:   true,
					Enter:    &ui.Motion{Scale: 0.7, Opacity: 0},
					Exit:     &ui.Motion{Scale: 0.7, Opacity: 0},
				})
			}
		})
	}
	return box
}

// choiceCheckPath is HeroUI's check mark fitted into a 10 DIP box: the
// upstream 17×18 viewBox drawn without stretching.
func choiceCheckPath(r ui.Rect) *ui.Path {
	var p ui.Path
	return p.MoveTo(r.X+r.W*0.084, r.Y+r.H*0.5).
		LineTo(r.X+r.W*0.417, r.Y+r.H*0.778).
		LineTo(r.X+r.W*0.861, r.Y+r.H*0.222)
}

// choiceInvalidRing paints HeroUI's unfocused invalid outline: a one DIP
// danger ring just outside the control.
func choiceInvalidRing(p *ui.Painter, r ui.Rect, radius float32, t theme.Palette) {
	p.Stroke(choiceExpand(r, 1), t.Danger, radius+1, 1)
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
			choiceErrorMessage(c, errText).PaddingStart(indent)
		}
	})
	return style.Disabled(outer, disabled)
}

// choiceErrorMessage is the muted field-level error a check box, radio or
// switch shows; HeroUI overrides the danger tone for these rows, which mark
// the control and label instead.
func choiceErrorMessage(c *ui.Context, text string) ui.Element {
	return style.Text(ui.Text(c, text), c, theme.HelpText()).
		TextColor(theme.Colors(c).Muted).Role(ui.RoleStatus)
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
