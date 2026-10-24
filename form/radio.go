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

// Radio is HeroUI's radio button: a focusable row that selects Value into the
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

	return choiceField(c, o.Description, errText, disabled, theme.Space(7), func() {
		row := ui.RadioBase(c, selected.Selected, o.Value).Gap(theme.Space12).Shrink(0).
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

// radioControl draws HeroUI's 16 DIP circle and its inner dot. owner supplies
// the live hover, pressed and focus state.
func radioControl(c *ui.Context, selected, invalid bool, variant ControlVariant, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	size := float32(theme.Space16)
	radius := theme.RadiusScale().LG
	box := ui.Box(c).Size(size, size).Shrink(0)
	// Snapshot the row's states in the build: a retained copy of a leaving
	// element paints in later passes, where owner has expired.
	hovered, pressed, focused := owner.Hovered(), owner.Pressed(), owner.FocusVisible()
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		fill := t.FieldBackground
		border := t.FieldBorder
		if variant == ControlSecondary {
			fill = t.Default
			if !selected && (hovered || pressed) {
				fill = t.DefaultHover
			}
		}
		switch {
		case selected:
			border = ui.Transparent
			fill = t.Accent
			if pressed {
				fill = t.AccentHover
			}
		case hovered:
			border = t.FieldBorderHover
			if variant != ControlSecondary {
				fill = t.FieldHover
			}
		}
		if variant != ControlSecondary {
			p.Shadow(r, radius, 0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.08))
		}
		p.Fill(r, fill, radius)
		if selected {
			d := r.W * 0.4286
			if pressed {
				d = r.W * 0.5714
			}
			p.Fill(ui.Rect{
				X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d,
			}, t.AccentForeground, d/2)
		}
		if invalid {
			p.Stroke(r, t.Danger, radius, theme.Hairline)
		} else if border.A != 0 {
			p.Stroke(r, border, radius, theme.Hairline)
		}
		if focused {
			choiceFocusRing(p, r, radius, invalid, t)
		}
	})
	return box
}
