package form

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// RadioGroupState is the caller-owned single choice shared by a group's radio
// buttons. Selected must point at stable caller storage; OnChange, when set,
// reports the value chosen.
type RadioGroupState[T comparable] struct {
	Selected *T
	OnChange func(T)
}

// Value returns the selected value and whether there is one.
func (s *RadioGroupState[T]) Value() (T, bool) {
	var zero T
	if s == nil || s.Selected == nil {
		return zero, false
	}
	return *s.Selected, true
}

// Select stores value and reports it through OnChange.
func (s *RadioGroupState[T]) Select(value T) {
	if s == nil || s.Selected == nil {
		return
	}
	*s.Selected = value
	if s.OnChange != nil {
		s.OnChange(value)
	}
}

// RadioOptions configures one native radio button. State carries the group's
// caller-owned selection; a nil State or Selected renders a disabled radio.
type RadioOptions[T comparable] struct {
	State       *RadioGroupState[T]
	Value       T
	Label       string
	Description string
	Error       string
	Required    bool
	Invalid     bool
	Disabled    bool
	Variant     ControlVariant
	OnChange    func(T)
}

// Radio is a focusable row that selects Value into the
// shared state when clicked or with Space. Its control paints from the live
// pointer and focus state of the row.
func Radio[T comparable](c *ui.Context, o RadioOptions[T]) ui.Element {
	selected := o.State
	if selected == nil || selected.Selected == nil {
		o.Disabled = true
		var fallback T
		selected = &RadioGroupState[T]{Selected: &fallback}
	}

	t := theme.Colors(c)
	disabled := o.Disabled
	errText := choiceErrorText(o.Invalid, o.Error)
	on := *selected.Selected == o.Value

	return choiceField(c, o.Description, errText, disabled, theme.Space24, func() {
		row := ui.RadioBase(c, selected.Selected, o.Value).Gap(theme.Space8).Shrink(0).
			FocusRing(false).TextColor(t.Foreground).Label(o.Label)
		if o.Description != "" {
			row.Description(o.Description)
		}
		if errText != "" {
			row.Error(errText)
		}
		if !disabled && selected.Selected != nil {
			row.OnChange(func() {
				if selected.OnChange != nil {
					selected.OnChange(o.Value)
				}
				if o.OnChange != nil {
					o.OnChange(o.Value)
				}
			})
		}

		row.Children(func() {
			radioControl(c, on, o.Invalid, o.Variant, row)
			if o.Label != "" {
				Label(c, o.Label, LabelOptions{
					Required: o.Required, Invalid: o.Invalid,
				})
			}
		})
	})
}

// radioControl paints the 16 DIP Coss radio and six-DIP selected dot.
func radioControl(c *ui.Context, selected, invalid bool, variant ControlVariant, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	size := theme.Scale(c, 16)
	radius := size / 2
	focused := owner.FocusVisible()
	fill, border := t.Background, t.Input
	if c.Theme().Dark {
		fill = t.Input.Alpha(0.32)
	}
	if selected {
		fill, border = t.Primary, t.Primary
	} else if invalid {
		border = t.Destructive.Alpha(0.36)
	}
	box := ui.Box(c).Size(size, size).Shrink(0).Center()
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, fill, radius)
		p.Stroke(r, border, radius, theme.Hairline)
		if selected {
			dot := theme.Scale(c, 6)
			p.Fill(ui.Rect{X: r.X + (r.W-dot)/2, Y: r.Y + (r.H-dot)/2, W: dot, H: dot}, t.PrimaryForeground, dot/2)
		}
		if focused {
			choiceFocusRing(p, r, radius, invalid, t)
		}
	})
	return box
}
