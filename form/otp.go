package form

import (
	"strings"
	"unicode/utf8"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type OTPPattern uint8

const (
	OTPDigits OTPPattern = iota
	OTPAlphanumeric
)

// OTPState retains only caller data and selection correction; MyGo retains
// the single editor's undo, selection, composition and clipboard behavior.
type OTPState struct {
	Value             string
	Handle            ui.Handle
	anchor, caret     int
	selectionPending  bool
	complete          string
	fontSize, advance float32
}

type OTPOptions struct {
	ControlOptions
	Length         int
	Pattern        OTPPattern
	SeparatorAfter int
	OnChange       func(string)
	OnComplete     func(string)
}

func InputOTP(c *ui.Context, state *OTPState, o OTPOptions) ui.Element {
	length := o.Length
	if length <= 0 {
		length = 6
	}
	colors := theme.Colors(c)
	fontSize := theme.Scale(c, 14)
	slot := max(float32(32), theme.Scale(c, 24)+8)
	if state.fontSize != fontSize || state.advance == 0 {
		state.fontSize = fontSize
		state.advance, _ = c.MeasureText(0, ui.Span{Text: "0", Font: theme.Mono, Size: fontSize})
	}
	e := ui.TextInputBase(c, &state.Value).Font(theme.Mono).FontSize(fontSize).
		TextColor(colors.Foreground).Direction(ui.LeftToRight).
		LetterSpacing(max(float32(0), slot-state.advance)).
		Padding(4, (slot-state.advance)/2).Width(float32(length) * slot).Height(slot).
		Radius(16).ReadOnly(o.ReadOnly).FocusRing(true)
	control := o.ControlOptions
	if control.Handle == nil {
		control.Handle = &state.Handle
	}
	configureControl(c, e, control)
	if state.selectionPending && !e.Composing() {
		e.SetTextSelection(state.anchor, state.caret)
		state.selectionPending = false
	}
	// Do not transform preedit. Committed input is reformatted through the
	// maintained editor's application-value path, which merges it with the
	// triggering undo step rather than creating a replacement editor.
	if !e.Composing() {
		normalized := normalizeOTP(state.Value, length, o.Pattern)
		if normalized != state.Value {
			anchor, caret := e.TextAnchorCaret()
			state.anchor = normalizedOTPOffset(state.Value, anchor, length, o.Pattern)
			state.caret = normalizedOTPOffset(state.Value, caret, length, o.Pattern)
			state.Value = normalized
			state.selectionPending = true
			c.Invalidate()
		}
	}
	composing := e.Composing()
	e.OnChange(func() {
		if composing || o.ReadOnly {
			return
		}
		if o.OnChange != nil {
			o.OnChange(state.Value)
		}
		if len(state.Value) == length && state.Value != state.complete {
			state.complete = state.Value
			if o.OnComplete != nil {
				o.OnComplete(state.Value)
			}
		} else if len(state.Value) != length {
			state.complete = ""
		}
	})
	focus := e.FocusVisible()
	_, caret := e.TextAnchorCaret()
	active := min(caret, length-1)
	invalid := o.Invalid || o.Error != ""
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, colors.Input.Alpha(0.5), 16)
		p.Stroke(r, colors.Input, 16, 1)
		for i := 1; i < length; i++ {
			x := r.X + float32(i)*slot
			width := float32(1)
			if o.SeparatorAfter > 0 && i%o.SeparatorAfter == 0 {
				width = 3
			}
			p.Fill(ui.Rect{X: x - width/2, Y: r.Y, W: width, H: r.H}, colors.Input, 0)
		}
		if focus || invalid {
			color := colors.Ring
			if invalid {
				color = colors.Destructive
			}
			p.Stroke(ui.Rect{X: r.X + float32(active)*slot, Y: r.Y, W: slot, H: r.H}, color, 8, 2)
		}
	})
	return e
}

func otpAccept(r rune, pattern OTPPattern) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return pattern == OTPAlphanumeric && (r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z')
}

func normalizeOTP(value string, length int, pattern OTPPattern) string {
	valid := true
	count := 0
	for _, r := range value {
		if !otpAccept(r, pattern) || count >= length {
			valid = false
			break
		}
		count++
	}
	if valid {
		return value
	}
	var out strings.Builder
	out.Grow(min(len(value), length))
	for _, r := range value {
		if otpAccept(r, pattern) {
			out.WriteRune(r)
			if out.Len() == length {
				break
			}
		}
	}
	return out.String()
}

func normalizedOTPOffset(value string, offset, length int, pattern OTPPattern) int {
	at, accepted := 0, 0
	for len(value) > 0 && at < offset && accepted < length {
		r, size := utf8.DecodeRuneInString(value)
		value = value[size:]
		at++
		if otpAccept(r, pattern) {
			accepted++
		}
	}
	return accepted
}
