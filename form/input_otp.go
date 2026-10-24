package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// OTPState is the caller-owned value of a one-time-code field. Value holds the
// accepted ASCII characters; Complete reports whether Value fills Length.
type OTPState struct {
	Value    string
	Complete bool

	completed string
}

// OTPOptions configures HeroUI's one-time-code field. Length defaults to 6
// and is capped at 12.
// Filter, when set, decides which characters are accepted; the default keeps
// printable ASCII without spaces. OnComplete runs once each time the value
// reaches Length.
type OTPOptions struct {
	Label, Description, Error string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	Length                    int
	Filter                    func(rune) bool
	OnChange                  func(value string)
	OnComplete                func(value string)
}

// InputOTP renders one native editor over HeroUI's slot guides. The editor
// keeps native focus, caret, selection and input methods; the guides, values
// and caret are painted from its committed text.
func InputOTP(c *ui.Context, state *OTPState, o OTPOptions) ui.Element {
	if state == nil {
		state = &OTPState{}
	}
	length := o.Length
	if length <= 0 {
		length = 6
	}
	if length > 12 {
		length = 12
	}
	filter := o.Filter
	if filter == nil {
		filter = otpPrintableASCII
	}

	// Keep the caller's value within the field before the editor reads it.
	state.Value = otpSanitize(state.Value, length, filter)
	state.Complete = len(state.Value) >= length
	if state.Complete {
		state.completed = state.Value
	} else {
		state.completed = ""
	}

	t := theme.Colors(c)
	radius := theme.RadiusScale().Field
	slotHeight := theme.Scale(c, 40)
	gap := float32(theme.Space8)
	reduceMotion := theme.ReduceMotion(c)

	build := func() ui.Element {
		runes := []rune(state.Value)
		filled := t.FieldFocus
		base, hover := t.FieldBackground, t.FieldHover
		if o.Variant == ControlSecondary {
			base, hover, filled = t.Default, t.DefaultHover, t.Default
		}
		container := ui.Row(c).Gap(gap).Height(slotHeight).AlignItems(ui.Center).FillWidth().
			MinWidth(theme.Scale(c, float32(length)*38+float32(length-1)*8))

		var editor ui.Element
		var selStart, selEnd, activeIndex int
		focusedNow := false
		container.Children(func() {
			for index := range length {
				char := ""
				if index < len(runes) {
					char = string(runes[index])
				}
				slot := ui.Box(c).Key(index).Grow(1).Basis(0).MinWidth(0).Height(slotHeight).
					Center().Radius(radius).Clip()
				if o.Variant == ControlPrimary {
					slot = style.Shadow(slot, c, theme.FieldShadows(c.Theme().Dark))
				}
				// Snapshot the slot's hover in the build: a retained copy of
				// a leaving element paints in later passes, where slot has
				// expired.
				slotHovered := !o.Disabled && slot.Hovered()
				slot.Draw(func(p *ui.Painter, r ui.Rect) {
					active := focusedNow && index == activeIndex
					selected := selStart != selEnd && index >= selStart && index < selEnd
					fill := base
					switch {
					case selected:
						fill = t.AccentSoft
					case char != "" || active:
						fill = filled
					case slotHovered:
						fill = hover
					}
					p.Fill(r, fill, radius)
					border := t.FieldBorder
					switch {
					case o.Invalid:
						border = t.Danger
					case active:
						border = t.FieldBorderFocus
					case slotHovered:
						border = t.FieldBorderHover
					}
					if border.A != 0 {
						p.Stroke(r, border, radius, theme.Hairline)
					}
					if active {
						p.FocusRing(r, [4]float32{radius, radius, radius, radius})
						if char == "" {
							visible := true
							if !reduceMotion {
								visible = (p.Now().UnixMilli()/500)%2 == 0
								p.AnimationFrame()
							}
							if visible {
								w := theme.Scale(c, 2)
								h := theme.Scale(c, 16)
								p.Fill(ui.Rect{
									X: r.X + (r.W-w)/2,
									Y: r.Y + (r.H-h)/2,
									W: w,
									H: h,
								}, t.FieldPlaceholder, w/2)
							}
						}
					}
				})
				if char != "" {
					slot.Children(func() {
						style.Text(ui.Text(c, char), c, theme.LG()).FontWeight(theme.WeightSemibold).
							TextColor(t.FieldForeground).LetterSpacing(theme.Scale(c, -0.27))
					})
				}
			}

			editor = ui.TextInputBase(c, &state.Value).Absolute().
				Left(0).Top(0).Right(0).Bottom(0).Padding(0).Opacity(0)
			name := inputAccessibleName(o.Label, o.AriaLabel, "")
			if name == "" {
				name = "Verification code"
			}
			editor.Label(name)
			if o.ReadOnly {
				editor.ReadOnly(true)
			}
			editor.OnChange(func() {
				sanitized := otpSanitize(state.Value, length, filter)
				changed := sanitized != state.Value
				state.Value = sanitized
				if len(sanitized) >= length {
					state.Complete = true
					if state.completed != sanitized {
						state.completed = sanitized
						if o.OnComplete != nil {
							o.OnComplete(sanitized)
						}
					}
				} else {
					state.Complete = false
					state.completed = ""
				}
				if changed && o.OnChange != nil {
					o.OnChange(sanitized)
				}
			})
		})
		focusedNow = editor.Valid() && editor.Focused()
		activeIndex = -1
		if focusedNow {
			selStart, selEnd = editor.TextSelection()
			activeIndex = min(selEnd, length-1)
		}
		if o.Disabled {
			container = style.Disabled(container, true)
		}
		return container
	}

	fieldOptions := InputOptions{
		Label:       o.Label,
		Description: o.Description,
		Error:       o.Error,
		Required:    o.Required,
		Disabled:    o.Disabled,
		Invalid:     o.Invalid,
	}
	if !inputNeedsField(fieldOptions) {
		return build()
	}
	return inputField(c, fieldOptions, true, func() { build() })
}

// otpPrintableASCII keeps the characters a one-time code can hold: visible
// ASCII, without spaces or control characters.
func otpPrintableASCII(r rune) bool { return r >= 0x21 && r <= 0x7e }

// otpSanitize filters s and keeps at most length accepted runes.
func otpSanitize(s string, length int, filter func(rune) bool) string {
	out := make([]rune, 0, min(len(s), length))
	for _, r := range s {
		if filter(r) {
			if len(out) == length {
				break
			}
			out = append(out, r)
		}
	}
	return string(out)
}
