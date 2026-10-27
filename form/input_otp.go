package form

import (
	"time"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// OTPState is the caller-owned value of a one-time-code field. Value holds the
// accepted characters; Complete reports whether Value fills Length.
type OTPState struct {
	Value    string
	Complete bool

	completed string
}

// OTPOptions configures a one-time-code field. Length defaults to 6
// and is capped at 12.
// Filter, when set, decides which characters are accepted; the default keeps
// printable ASCII without spaces. OnComplete runs once each time the value
// reaches Length. OnChange reports changes to the accepted value.
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

// InputOTP renders one native editor over slot guides. The editor
// keeps native focus, caret, selection and input methods; the guides, values
// and caret are painted from its committed text. Rejecting characters can
// reposition the caret when the native editor receives the filtered value.
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
	state.Complete = utf8.RuneCountInString(state.Value) >= length
	if state.Complete {
		state.completed = state.Value
	} else {
		state.completed = ""
	}

	t := theme.Colors(c)
	radius := theme.RadiusScale().Field
	slotHeight := theme.Scale(c, 32)
	gap := float32(theme.Space8)
	reduceMotion := theme.ReduceMotion(c)

	build := func() ui.Element {
		// The accepted characters, in a fixed buffer: the value is already
		// filtered and at most 12 runes, so the build copies no slice.
		var runes [12]rune
		n := 0
		for _, r := range state.Value {
			if n == length {
				break
			}
			runes[n] = r
			n++
		}
		filled := t.Background
		base, hover := t.Background, t.Background
		if o.Variant == ControlSecondary {
			base, hover, filled = t.Background, t.Background, t.Background
		}
		container := ui.Row(c).Gap(gap).Height(slotHeight).AlignItems(ui.Center).FillWidth().
			MinWidth(theme.Scale(c, float32(length)*32+float32(length-1)*8))

		var slotShadows []theme.Shadow
		if o.Variant == ControlPrimary {
			// One slice per build, not one per slot.
			slotShadows = theme.FieldShadows(c.Theme().Dark)
		}

		var editor ui.Element
		var selStart, selEnd, activeIndex int
		focusedNow := false
		container.Children(func() {
			for index := range length {
				char := ""
				if index < n {
					char = string(runes[index])
				}
				slot := ui.Box(c).Key(index).Width(slotHeight).Height(slotHeight).
					Center().Radius(radius).Clip()
				if o.Variant == ControlPrimary {
					slot = style.Shadow(slot, c, slotShadows)
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
						fill = t.Selection
					case char != "" || active:
						fill = filled
					case slotHovered:
						fill = hover
					}
					p.Fill(r, fill, radius)
					border := t.Input
					switch {
					case o.Invalid:
						border = t.Destructive.Alpha(0.36)
					case active:
						border = t.Ring
					case slotHovered:
						border = t.Input
					}
					if border.A != 0 {
						p.Stroke(r, border, radius, theme.Hairline)
					}
					if active {
						ring := t.Ring.Alpha(0.24)
						if o.Invalid {
							ring = t.Destructive.Alpha(0.16)
						}
						formstyle.FieldRing(p, r, radius, ring)
						if char == "" {
							// Blink at 500 ms a phase, but repaint only at the
							// boundary instead of at every display frame.
							visible := true
							if !reduceMotion {
								const phase = 500
								now := p.Now().UnixMilli()
								left := now % phase
								visible = left < phase/2
								p.After(time.Duration(phase-left) * time.Millisecond)
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
						style.Text(ui.Text(c, char), c, theme.SM()).FontWeight(theme.WeightMedium).
							TextColor(t.FieldForeground).LetterSpacing(theme.Scale(c, -0.27))
					})
				}
			}

			previous := state.Value
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
				state.Value = sanitized
				if utf8.RuneCountInString(sanitized) >= length {
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
				if sanitized != previous && o.OnChange != nil {
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

// otpSanitize filters s and keeps at most length accepted runes. A value that
// already satisfies the filter and the length is returned untouched, so the
// steady state of a field does not allocate on every build.
func otpSanitize(s string, length int, filter func(rune) bool) string {
	n := 0
	clean := true
	for _, r := range s {
		if !filter(r) {
			clean = false
			continue
		}
		n++
		if n > length {
			clean = false
		}
	}
	if clean {
		return s
	}
	out := make([]rune, 0, min(n, length))
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
