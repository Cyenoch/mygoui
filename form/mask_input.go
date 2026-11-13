package form

import (
	"math"
	"slices"
	"unicode"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// MaskInputOptions configures a masked input. In Mask, 9 takes a digit, A a
// letter and * either; every other rune is a literal the input writes itself.
// The mask has no escape, so 9, A and * cannot be literals.
type MaskInputOptions struct {
	Mask                      string
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	OnChange, OnSubmit        func()
}

// maskState is the per-field state kept across builds.
type maskState struct {
	// selection is the selection before this frame's input, which places the
	// edit that input made.
	selection [2]int
	// undo and redo hold raw values. kind is the direction of the last
	// change that joined the step at the top of undo: one character typed
	// after another, or deleted after another, makes one step.
	undo, redo []string
	kind       int
	// pasted is text the user pasted that the editor takes in this frame,
	// cut to what the mask holds. whole marks a paste that carries a full
	// value, which replaces the text instead of joining it.
	pasted string
	whole  bool
	// toEnd puts the caret after the text an undo or a redo set.
	toEnd bool
}

type maskStateKey struct{}

// maskUndoDepth caps the steps of undo.
const maskUndoDepth = 100

// MaskInput is a native editor constrained to Mask. The caller's value holds
// only the characters typed into the mask's slots, so "99/99/9999" showing
// 12/25/1999 stores "12251999"; MaskFormat renders it elsewhere.
//
// Literals appear as the slot after them fills, and those ending the mask
// once it is full. The slots still empty show as underscores while the
// editor has the focus or text, unless a Placeholder owns the empty field.
// Text typed or pasted anywhere is fitted to the slots from that point: a
// character its slot does not take is dropped, as is whatever overflows the
// mask, and removing a character moves the later ones back. A paste that
// holds more characters than the mask has slots, as a phone number with its
// country code does, keeps the last ones. Backspace and Delete over a
// literal remove the character beyond it, and a removal that would move a
// later character into a slot that rejects it is refused. The caret stays
// with the edit. Undo and redo step through the raw values, a run of typing
// or deleting being one step.
func MaskInput(c *ui.Context, value *string, o MaskInputOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	value = inputTextAnchor(value)
	fieldOptions := InputOptions{
		Label: o.Label, Description: o.Description, Error: o.Error,
		Required: o.Required, Disabled: o.Disabled, Invalid: o.Invalid, FullWidth: o.FullWidth,
	}
	if !inputNeedsField(fieldOptions) {
		return maskEditor(c, value, o)
	}
	var editor ui.Element
	inputField(c, fieldOptions, true, func() { editor = maskEditor(c, value, o) })
	return editor
}

// MaskFormat renders raw slot characters through mask, as MaskInput shows
// them. Characters their slot does not take are skipped.
func MaskFormat(mask, raw string) string {
	text, _ := maskFormat([]rune(mask), []rune(raw))
	return string(text)
}

// maskEditor binds the editor to the formatted text of the frame and turns
// what the user did to that text into the next raw value.
func maskEditor(c *ui.Context, value *string, o MaskInputOptions) ui.Element {
	mask := []rune(o.Mask)
	shown, raw := maskFormat(mask, []rune(*value))
	text := string(shown)
	var e ui.Element
	// The ghost slots sit over the editor, in a box that is its size.
	frame := ui.Box(c).MinWidth(0)
	if o.FullWidth {
		frame.FillWidth()
	}
	frame.Children(func() {
		e = inputEditor(c, &text, inputCoreOptions{
			Placeholder: o.Placeholder,
			Name:        inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder),
			Variant:     o.Variant,
			Disabled:    o.Disabled,
			ReadOnly:    o.ReadOnly,
			Invalid:     o.Invalid,
			FullWidth:   o.FullWidth,
			OnSubmit:    o.OnSubmit,
		})
		st := ui.Local(e, maskStateKey{}, func() maskState { return maskState{} })
		if !o.ReadOnly && !o.Disabled {
			maskBindKeys(e, st, value, mask, o.OnChange)
		}
		if st.toEnd {
			e.SetTextSelection(len(shown), len(shown))
			st.toEnd = false
		}
		selection := st.selection
		pasting := st.pasted != ""
		if pasting {
			start, end := e.TextSelection()
			if st.whole {
				start, end = 0, math.MaxInt
			}
			e.ReplaceText(start, end, st.pasted)
			st.pasted, st.whole = "", false
		}
		if edited := []rune(text); !slices.Equal(edited, shown) {
			_, caret := e.TextSelection()
			next, caret := maskEdit(mask, raw, shown, edited, selection, caret)
			formatted, _ := maskFormat(mask, next)
			if !slices.Equal(formatted, edited) {
				e.ReplaceText(0, math.MaxInt, string(formatted))
			}
			e.SetTextSelection(caret, caret)
			if !slices.Equal(next, raw) {
				st.record(string(raw), len(next)-len(raw))
				raw = next
				*value = string(next)
				if o.OnChange != nil {
					o.OnChange()
				}
			}
			pasting = false
		}
		// A paste the editor has not yet shown is read against the selection
		// it replaced, which the next frame still needs.
		if !pasting {
			start, end := e.TextSelection()
			st.selection = [2]int{start, end}
		}
		maskGhost(c, e, mask, raw, o)
	})
	return e
}

// maskBindKeys gives the editor's undo, redo and paste to the raw value: the
// editor's own would bring back text as it stood between an edit and its
// formatting, which reads as a new edit.
func maskBindKeys(e ui.Element, st *maskState, value *string, mask []rune, changed func()) {
	e.OnTextKey(func(k ui.TextKey) bool {
		if k.Mods&ui.Cmd == 0 || k.Mods&^(ui.Cmd|ui.Shift) != 0 {
			return false
		}
		var from, to *[]string
		switch {
		case k.Key == ui.KeyZ && k.Mods == ui.Cmd:
			from, to = &st.undo, &st.redo
		case k.Key == ui.KeyZ, k.Key == ui.KeyY && k.Mods == ui.Cmd:
			from, to = &st.redo, &st.undo
		default:
			return false
		}
		if len(*from) == 0 {
			return true
		}
		last := len(*from) - 1
		step := (*from)[last]
		*from = (*from)[:last]
		*to = append(*to, *value)
		*value, st.kind, st.toEnd = step, 0, true
		k.SetText(MaskFormat(string(mask), step))
		if changed != nil {
			changed()
		}
		return true
	})
	e.OnPaste(func(text string) bool {
		kept := maskPasteTail(mask, text)
		if kept == "" {
			return false
		}
		st.pasted, st.whole = kept, true
		return true
	})
}

// record notes that raw was replaced by a value delta runes longer, as a
// step of undo; a run of typing or of deleting is a single step.
func (s *maskState) record(raw string, delta int) {
	kind := 0
	if delta == 1 || delta == -1 {
		kind = delta
	}
	s.redo = s.redo[:0]
	if kind != 0 && kind == s.kind && len(s.undo) > 0 {
		return
	}
	s.kind = kind
	if len(s.undo) == maskUndoDepth {
		s.undo = s.undo[1:]
	}
	s.undo = append(s.undo, raw)
}

// maskPasteTail returns the end of pasted text that the mask can hold when
// the text carries at least as many characters for its slots as the mask has,
// as a phone number does with its country code, and "" when the paste is
// shorter and joins the text at the caret.
func maskPasteTail(mask []rune, text string) string {
	slots := maskSlotsBefore(mask, len(mask))
	var taken []rune
	for _, r := range text {
		for _, m := range mask {
			if maskSlot(m) && maskAccepts(m, r) {
				taken = append(taken, r)
				break
			}
		}
	}
	if len(taken) < slots {
		return ""
	}
	return string(taken[len(taken)-slots:])
}

// maskGhost shows the slots raw leaves empty as underscores after the text,
// by drawing the text itself invisibly in front of them.
func maskGhost(c *ui.Context, e ui.Element, mask, raw []rune, o MaskInputOptions) {
	if o.Disabled || o.ReadOnly {
		return
	}
	text, _ := maskFormat(mask, raw)
	if len(text) == 0 && (o.Placeholder != "" || !e.Focused()) || len(text) >= len(mask) {
		return
	}
	ghost := make([]rune, 0, len(mask)-len(text))
	for _, m := range mask[len(text):] {
		if maskSlot(m) {
			m = '_'
		}
		ghost = append(ghost, m)
	}
	ink := theme.Colors(c).MutedForeground
	fieldText(ui.RichText(c,
		ui.Span{Text: string(text), Color: ui.Transparent},
		ui.Span{Text: string(ghost), Color: ink},
	), c).SingleLine().Absolute().Start(theme.Space12).Top(theme.Scale(c, 5)).
		IgnoreInput(true).AccessHidden(true)
}

// maskSlot reports whether a mask rune takes a typed character.
func maskSlot(m rune) bool { return m == '9' || m == 'A' || m == '*' }

// maskAccepts reports whether slot m takes r.
func maskAccepts(m, r rune) bool {
	switch m {
	case '9':
		return unicode.IsDigit(r)
	case 'A':
		return unicode.IsLetter(r)
	}
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

// maskNextSlot returns the index of the first slot at or after i, or -1.
func maskNextSlot(mask []rune, i int) int {
	for ; i < len(mask); i++ {
		if maskSlot(mask[i]) {
			return i
		}
	}
	return -1
}

// maskSlotAt returns the index of slot number k, or -1 without one.
func maskSlotAt(mask []rune, k int) int {
	for i, m := range mask {
		if !maskSlot(m) {
			continue
		}
		if k == 0 {
			return i
		}
		k--
	}
	return -1
}

// maskSlotsBefore counts the slots in mask[:i].
func maskSlotsBefore(mask []rune, i int) int {
	n := 0
	for _, m := range mask[:min(i, len(mask))] {
		if maskSlot(m) {
			n++
		}
	}
	return n
}

// maskFormat fits raw into mask and returns the text shown with the raw
// characters kept. Rune i of the text sits at mask[i].
func maskFormat(mask, raw []rune) (text, kept []rune) {
	at := 0
	for _, r := range raw {
		slot := maskNextSlot(mask, at)
		if slot < 0 {
			break
		}
		if !maskAccepts(mask[slot], r) {
			continue
		}
		text = append(append(text, mask[at:slot]...), r)
		kept = append(kept, r)
		at = slot + 1
	}
	if len(kept) > 0 && maskNextSlot(mask, at) < 0 {
		text = append(text, mask[at:]...)
	}
	return text, kept
}

// maskSpan finds what the change from shown to edited kept at either end. It
// reads the edit from the selection before it and the caret after it, since
// a comparison of the texts cannot place an edit among equal characters; it
// compares them when the two do not account for the change, as after a drop.
func maskSpan(shown, edited []rune, selection [2]int, caret int) (prefix, suffix int) {
	prefix = min(selection[0], caret)
	suffix = len(edited) - caret
	if prefix >= 0 && suffix >= 0 && prefix+suffix <= min(len(shown), len(edited)) &&
		slices.Equal(shown[:prefix], edited[:prefix]) &&
		slices.Equal(shown[len(shown)-suffix:], edited[caret:]) {
		return prefix, suffix
	}
	prefix, suffix = 0, 0
	for prefix < len(shown) && prefix < len(edited) && shown[prefix] == edited[prefix] {
		prefix++
	}
	for suffix < min(len(shown), len(edited))-prefix &&
		shown[len(shown)-1-suffix] == edited[len(edited)-1-suffix] {
		suffix++
	}
	return prefix, suffix
}

// maskEdit turns the change from shown, the formatted text of raw, to edited
// into the next raw value and the caret in its formatted text. selection is
// the selection shown had before the change and caret the caret after it.
func maskEdit(mask, raw, shown, edited []rune, selection [2]int, caret int) ([]rune, int) {
	prefix, suffix := maskSpan(shown, edited, selection, caret)
	typed := edited[prefix : len(edited)-suffix]
	// raw[from:to] are the slot characters the change removed.
	from, to := maskSlotsBefore(mask, prefix), maskSlotsBefore(mask, len(shown)-suffix)
	if len(typed) == 0 && from == to && selection[0] == selection[1] {
		// Only a literal went: take the character beyond it instead.
		if selection[0] > prefix {
			from = max(from-1, 0)
		} else {
			to = min(to+1, len(raw))
		}
	}
	next := slices.Clone(raw[:from])
	// fit appends the runes their slots take and reports whether all fit.
	fit := func(runes []rune) bool {
		all := true
		for _, r := range runes {
			slot := maskSlotAt(mask, len(next))
			if slot >= 0 && maskAccepts(mask[slot], r) {
				next = append(next, r)
			} else {
				all = false
			}
		}
		return all
	}
	fit(typed)
	at := len(next)
	if !fit(raw[to:]) && len(typed) == 0 {
		// Closing the gap would push a later character into a slot that
		// rejects it: keep the text rather than lose that character.
		return raw, min(selection[0], len(shown))
	}
	return next, maskCaret(mask, next, at, len(typed) > 0)
}

// maskCaret places the caret after at raw characters of next's formatted
// text: past the literals that follow when the user typed, before them when
// the user deleted.
func maskCaret(mask, next []rune, at int, typed bool) int {
	if typed && at >= len(next) {
		text, _ := maskFormat(mask, next)
		return len(text)
	}
	if typed {
		return maskSlotAt(mask, at)
	}
	return maskSlotAt(mask, at-1) + 1
}
