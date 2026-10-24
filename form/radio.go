package form

import (
	"time"

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
// the live hover, pressed and focus state. The dot fills the control while
// unselected and shrinks to the 6 DIP mark when chosen, as upstream's
// indicator does.
func radioControl(c *ui.Context, selected, invalid bool, variant ControlVariant, owner ui.Element) ui.Element {
	t := theme.Colors(c)
	size := float32(theme.Space16)
	radius := theme.RadiusScale().LG
	// Snapshot the row's states in the build: a retained copy of a leaving
	// element paints in later passes, where owner has expired.
	hovered, pressed, focused := owner.Hovered(), owner.Pressed(), owner.FocusVisible()
	reduceMotion := theme.ReduceMotion(c)

	baseFill := t.FieldBackground
	if variant == ControlSecondary {
		baseFill = t.Default
	}
	fill := baseFill
	switch {
	case selected:
		fill = t.Accent
		if pressed {
			fill = t.AccentHover
		}
	case hovered:
		if variant == ControlSecondary {
			fill = t.DefaultHover
		} else {
			fill = t.FieldHover
		}
	}
	// The upstream dot is the 14 DIP inside the border, which scales to
	// 0.4286 (6 DIP) when chosen and 0.5714 (8 DIP) while pressed.
	dotScale := float32(1)
	switch {
	case selected && pressed:
		dotScale = 0.5714
	case selected:
		dotScale = 0.4286
	}
	dotFill := baseFill
	switch {
	case selected:
		dotFill = t.AccentForeground
	case hovered && variant == ControlSecondary:
		dotFill = t.DefaultHover
	case hovered:
		dotFill = t.FieldHover
	}

	box := ui.Box(c).Size(size, size).Shrink(0).Center()
	if !reduceMotion {
		// A color change is part of the switch between faces upstream, so
		// it follows the same 200 ms curve as the dot's scale.
		dotFill = box.AnimateColor("radio-dot", dotFill, ui.Animation{
			Duration: 200 * time.Millisecond, Ease: ui.EaseOut,
		})
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		border := t.FieldBorder
		if !selected && !invalid && hovered {
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
	box.Children(func() {
		dot := ui.Box(c).Size(size-2*theme.Hairline, size-2*theme.Hairline).Shrink(0)
		dot.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, dotFill, r.W/2)
		})
		dot.AnimatePose(ui.Pose{Scale: dotScale}, ui.Animation{
			Duration: 200 * time.Millisecond, Ease: ui.EaseOut,
		})
	})
	return box
}
