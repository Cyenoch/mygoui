package form

import (
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CheckboxCardOptions configures a card that is one check box. Checked is
// caller-owned; a nil Checked renders the card disabled. Title names the
// card for assistive technology, which reads Description after it. Content
// builds more to show under the texts, such as a price or a feature list; it
// is display only, as the whole card is the control. HideIndicator leaves out
// the check box, the card's border and fill showing the choice.
type CheckboxCardOptions struct {
	Checked       *bool
	Title         string
	Description   string
	Icon          *ui.SVG
	Content       func()
	HideIndicator bool
	Invalid       bool
	Disabled      bool
	FullWidth     bool
	OnChange      func(bool)
}

// CheckboxCard is a native check box drawn as a whole card: a click anywhere
// on it, or Space while it has the focus, flips the caller's value. Group
// several in a CheckboxGroup.
func CheckboxCard(c *ui.Context, o CheckboxCardOptions) ui.Element {
	checked := o.Checked
	if checked == nil {
		checked = new(bool)
		o.Disabled = true
	}
	card := ui.CheckboxBase(c, checked)
	if o.OnChange != nil && !o.Disabled {
		card.OnChange(func() { o.OnChange(*checked) })
	}
	return choiceCard(c, card, choiceCardFace{
		title: o.Title, description: o.Description, icon: o.Icon,
		content: o.Content, hideIndicator: o.HideIndicator,
		on: *checked, invalid: o.Invalid, disabled: o.Disabled, fullWidth: o.FullWidth,
	}, func() {
		checkboxControl(c, checked, nil, o.Invalid, ControlPrimary, ui.Element{})
	})
}

// RadioCardOptions configures a card that is one radio button. State carries
// the group's caller-owned selection; a nil State or Selected renders the
// card disabled. Content and HideIndicator are those of CheckboxCardOptions.
type RadioCardOptions[T comparable] struct {
	State         *RadioGroupState[T]
	Value         T
	Title         string
	Description   string
	Icon          *ui.SVG
	Content       func()
	HideIndicator bool
	Invalid       bool
	Disabled      bool
	FullWidth     bool
	OnChange      func(T)
}

// RadioCard is a native radio button drawn as a whole card that selects
// Value into the shared state. Build the cards of one choice inside a
// RadioGroup, which makes them one Tab stop moved by the arrows.
func RadioCard[T comparable](c *ui.Context, o RadioCardOptions[T]) ui.Element {
	state := o.State
	if state == nil || state.Selected == nil {
		state = &RadioGroupState[T]{Selected: new(T)}
		o.Disabled = true
	}
	card := ui.RadioBase(c, state.Selected, o.Value)
	if !o.Disabled {
		card.OnChange(func() {
			if state.OnChange != nil {
				state.OnChange(o.Value)
			}
			if o.OnChange != nil {
				o.OnChange(o.Value)
			}
		})
	}
	on := *state.Selected == o.Value
	return choiceCard(c, card, choiceCardFace{
		title: o.Title, description: o.Description, icon: o.Icon,
		content: o.Content, hideIndicator: o.HideIndicator,
		on: on, invalid: o.Invalid, disabled: o.Disabled, fullWidth: o.FullWidth,
	}, func() {
		radioControl(c, on, o.Invalid, ControlPrimary, ui.Element{})
	})
}

// choiceCardFace is what a check box card and a radio card show alike.
type choiceCardFace struct {
	title, description string
	icon               *ui.SVG
	content            func()
	hideIndicator      bool
	on                 bool
	invalid, disabled  bool
	fullWidth          bool
}

// choiceCard styles a native check box or radio base as a card with the
// icon, texts and the control built by indicator. The card carries the focus
// ring, so the indicator is built without an owner to ring.
func choiceCard(c *ui.Context, card ui.Element, f choiceCardFace, indicator func()) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().LG
	card = card.Row().AlignItems(ui.Start).Gap(theme.Space12).MinWidth(0).
		Padding(theme.Space12, theme.Space16).FocusRing(false).Cursor(ui.CursorPointer).
		Label(f.title)
	if f.description != "" {
		card.Description(f.description)
	}
	if f.fullWidth {
		card.FillWidth()
	}
	card = style.Disabled(card, f.disabled)
	hovered, focused := card.Hovered(), card.FocusVisible()
	border := t.Input
	switch {
	case f.on:
		border = t.Primary
	case f.invalid:
		border = t.Destructive.Alpha(0.36)
	}
	border = motion.Color(card, "border", border)
	selectedFill := motion.Color(card, "fill", t.Accent.Alpha(0.5*choiceOnOff(f.on)))
	card.Draw(func(p *ui.Painter, r ui.Rect) {
		if hovered && !f.on {
			p.Fill(r, t.Accent.Alpha(0.5), radius)
		}
		p.Fill(r, selectedFill, radius)
		p.Stroke(r, border, radius, theme.Hairline)
		if focused {
			formstyle.FieldRing(p, r, radius, t.Ring.Alpha(0.24))
		}
	})
	card.Children(func() {
		if f.icon != nil {
			size := theme.Scale(c, 20)
			ui.Icon(c, f.icon).Size(size, size).Shrink(0).TextColor(t.MutedForeground)
		}
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space2).Children(func() {
			style.Text(ui.Text(c, f.title), c, theme.Label()).TextColor(t.Foreground)
			if f.description != "" {
				Description(c, f.description)
			}
			if f.content != nil {
				f.content()
			}
		})
		if f.hideIndicator {
			return
		}
		// Two DIPs centre the 16 DIP control on the title's 20 DIP line.
		ui.Box(c).Shrink(0).PaddingY(theme.Scale(c, 2)).Children(indicator)
	})
	return card
}
