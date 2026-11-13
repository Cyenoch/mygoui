package form

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TagsInputOptions configures a free-text tag field. Max caps the number of
// tags when positive. A tag already present is rejected unless
// AllowDuplicates is set; tags compare exactly, after trimming spaces.
// A draft still in the editor when the focus leaves the field is committed
// as a tag, unless KeepDraftOnBlur is set.
type TagsInputOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Required, Disabled        bool
	Invalid, FullWidth        bool
	Max                       int
	AllowDuplicates           bool
	KeepDraftOnBlur           bool
	OnChange                  func()
}

// tagsState is the per-field state kept across builds.
type tagsState struct {
	draft string
	// refused is the draft text a commit was refused for; the field shows it
	// until the draft or the tags change that.
	refused string
	attempt bool
	focus   int
	pasted  []string
}

const (
	tagsFocusNone   = -1
	tagsFocusEditor = -2
)

type tagsStateKey struct{}

// TagsInput edits the caller's tags as removable tokens before a native
// editor. Enter or a comma commits what was typed, and so does the focus
// leaving the field; pasted text splits at commas, tabs and line breaks into
// tags. A tag the options refuse stays in the editor and the field shows it
// as invalid, with the existing token tinted for a duplicate, until the
// draft changes. At Max the editor stops taking text.
//
// Backspace in the empty editor removes the last tag, and Left (Right in a
// right-to-left layout) moves to the last token. The tokens are one group the arrows move
// through, Backspace and Delete remove the focused one and pass the focus
// to its neighbour, and the arrow past the last returns to the editor.
func TagsInput(c *ui.Context, tags *[]string, o TagsInputOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	if tags == nil {
		tags = new([]string)
		o.Disabled = true
	}
	fieldOptions := InputOptions{
		Label: o.Label, Description: o.Description, Error: o.Error,
		Required: o.Required, Disabled: o.Disabled, Invalid: o.Invalid, FullWidth: o.FullWidth,
	}
	if !inputNeedsField(fieldOptions) {
		return tagsShell(c, tags, o)
	}
	return inputField(c, fieldOptions, true, func() { tagsShell(c, tags, o) })
}

// tagsShell builds the field face holding the tokens and the editor.
func tagsShell(c *ui.Context, tags *[]string, o TagsInputOptions) ui.Element {
	shell := ui.Row(c).Wrap().AlignItems(ui.Center).Gap(theme.Space(1)).
		Padding(theme.Space(1), theme.Space(1.5)).MinHeight(theme.Scale(c, 32)).Shrink(0).
		Cursor(ui.CursorText)
	if o.FullWidth {
		shell.FillWidth()
	} else {
		shell.Width(theme.Scale(c, 280))
	}
	st := ui.Local(shell, tagsStateKey{}, func() tagsState { return tagsState{focus: tagsFocusNone} })
	changed := func() {
		if o.OnChange != nil {
			o.OnChange()
		}
	}
	// The refused draft is shown while it would still be refused.
	text, _, _ := strings.Cut(st.draft, ",")
	if text = strings.TrimSpace(text); text != st.refused {
		st.attempt = false
	}
	duplicate := -1
	if st.attempt && !o.AllowDuplicates {
		duplicate = slices.Index(*tags, text)
	}
	full := o.Max > 0 && len(*tags) >= o.Max
	refusing := st.attempt && (duplicate >= 0 || full)
	if !refusing {
		st.attempt = false
	}
	rtl := c.Direction() == ui.RightToLeft
	var editor ui.Element
	remove, backspaced := -1, false
	shell.Children(func() {
		for i := 0; i < len(*tags); i++ {
			at, tag := i, (*tags)[i]
			chip := tagsToken(c.Key(tagsKeyAt(*tags, i)), tag, o.Disabled, i == duplicate, func() {
				// Backspace may have removed tags since this row was built.
				if at >= len(*tags) || (*tags)[at] != tag {
					return
				}
				*tags = slices.Delete(*tags, at, at+1)
				changed()
			})
			chip.Role(ui.RoleListItem).Label(tag).PositionInSet(i+1, len(*tags))
			if o.Disabled {
				continue
			}
			// A token takes the focus only once the keyboard moved it there:
			// a field's label focuses its first control, which is the editor.
			if st.focus == i || chip.Focused() {
				chip.Focusable().TabStop(false)
			}
			if st.focus == i {
				chip.Focus()
				st.focus = tagsFocusNone
			}
			if !chip.Focused() {
				continue
			}
			back, forward := ui.KeyLeft, ui.KeyRight
			if rtl {
				back, forward = forward, back
			}
			backspace, removal := chip.Shortcut(0, ui.KeyBackspace), chip.Shortcut(0, ui.KeyDelete)
			switch {
			case chip.Shortcut(0, back):
				st.focus = max(i-1, 0)
			case chip.Shortcut(0, forward):
				st.focus = i + 1
				if i+1 == len(*tags) {
					st.focus = tagsFocusEditor
				}
			case backspace || removal:
				// The tags stay as built; the removal applies after the loop.
				remove, backspaced = i, backspace
			}
		}
		// The editor keeps its identity, whatever tokens come before it.
		ui.Box(c.Key(tagsStateKey{})).Grow(1).Basis(0).MinWidth(theme.Scale(c, 80)).Children(func() {
			editor = tagsEditor(c, tags, st, o, full, changed)
		})
		if o.Max > 0 {
			tagsCount(c, len(*tags), o.Max, refusing && full)
		}
	})
	if remove >= 0 {
		*tags = slices.Delete(*tags, remove, remove+1)
		changed()
		c.Invalidate()
		// Backspace leans back, as in a text field; Delete forward.
		st.focus = min(remove, len(*tags)-1)
		if backspaced && remove > 0 {
			st.focus = remove - 1
		}
		if len(*tags) == 0 {
			st.focus = tagsFocusEditor
		}
	}
	if !o.Disabled {
		shell.OnClick(func() { editor.Focus() })
		if st.focus == tagsFocusEditor {
			editor.Focus()
			st.focus = tagsFocusNone
		}
		if shell.FocusLeft() && !o.KeepDraftOnBlur {
			tagsCommit(tags, st, o)
			c.Invalidate()
		}
	}
	return formstyle.SurfaceFocused(shell, c, editor.Focused(), o.Invalid || refusing, o.Disabled)
}

// tagsEditor is the native editor after the tokens. It commits the draft on
// Enter and at each comma, and gives Backspace to the last tag while empty.
func tagsEditor(c *ui.Context, tags *[]string, st *tagsState, o TagsInputOptions, full bool, changed func()) ui.Element {
	e := ui.TextInputBase(c, &st.draft).FillWidth().MinWidth(0).
		PaddingX(theme.Space(1.5)).PaddingY(theme.Scale(c, 2))
	e = fieldText(e, c)
	if o.Placeholder != "" && len(*tags) == 0 {
		e.Placeholder(o.Placeholder)
	}
	if name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder); name != "" {
		e.Label(name)
	}
	if full {
		e.ReadOnly(true)
	}
	back := ui.KeyLeft
	if c.Direction() == ui.RightToLeft {
		back = ui.KeyRight
	}
	e.OnTextKey(func(k ui.TextKey) bool {
		if k.Mods != 0 {
			return false
		}
		switch {
		case k.Key == ui.KeyEscape && k.Text() != "":
			k.SetText("")
		case (k.Key == ui.KeyBackspace || k.Key == back) && k.Text() == "" && len(*tags) > 0:
			if k.Key == ui.KeyBackspace {
				*tags = (*tags)[:len(*tags)-1]
				changed()
			} else {
				st.focus = len(*tags) - 1
			}
		default:
			return false
		}
		return true
	})
	e.OnPaste(func(text string) bool {
		parts := tagsSplit(text)
		if len(parts) == 0 || !strings.ContainsAny(text, ",\t\r\n") {
			return false
		}
		st.pasted = append(st.pasted, parts...)
		return true
	})
	before := len(*tags)
	if len(st.pasted) > 0 {
		// A pasted tag the options refuse stays in the editor ahead of the draft.
		kept := tagsFlush(st, tags, st.pasted, o)
		if draft := strings.TrimSpace(st.draft); len(kept) > 0 && draft != "" {
			kept = append(kept, draft)
		}
		if len(kept) > 0 {
			st.draft = strings.Join(kept, ", ")
		}
		st.pasted = st.pasted[:0]
	}
	if strings.Contains(st.draft, ",") {
		parts := strings.Split(st.draft, ",")
		last := parts[len(parts)-1]
		kept := tagsFlush(st, tags, parts[:len(parts)-1], o)
		if len(kept) == 0 {
			st.draft = strings.TrimLeft(last, " ")
		} else {
			if last = strings.TrimSpace(last); last != "" {
				kept = append(kept, last)
			}
			st.draft = strings.Join(kept, ", ")
		}
	}
	if e.Submitted() {
		tagsCommit(tags, st, o)
	}
	if len(*tags) != before {
		changed()
		c.Invalidate()
	}
	return e
}

// tagsCommit turns the draft into tags. A piece the options refuse stays.
func tagsCommit(tags *[]string, st *tagsState, o TagsInputOptions) {
	kept := tagsFlush(st, tags, strings.Split(st.draft, ","), o)
	st.draft = strings.Join(kept, ", ")
}

// tagsFlush adds the pieces in order and returns those the options refuse,
// the first of which the field then shows as refused.
func tagsFlush(st *tagsState, tags *[]string, pieces []string, o TagsInputOptions) (refused []string) {
	for _, piece := range pieces {
		if tagsAdd(tags, piece, o) == tagsRefused {
			refused = append(refused, strings.TrimSpace(piece))
		}
	}
	if len(refused) > 0 {
		st.refused, st.attempt = refused[0], true
	}
	return refused
}

type tagsResult int

const (
	tagsAdded tagsResult = iota
	tagsEmpty
	tagsRefused
)

// tagsAdd appends text as a tag when the options accept it.
func tagsAdd(tags *[]string, text string, o TagsInputOptions) tagsResult {
	text = strings.TrimSpace(text)
	switch {
	case text == "":
		return tagsEmpty
	case o.Max > 0 && len(*tags) >= o.Max:
	case !o.AllowDuplicates && slices.Contains(*tags, text):
	default:
		*tags = append(*tags, text)
		return tagsAdded
	}
	return tagsRefused
}

// tagsSplit cuts pasted text into tag candidates at commas, tabs and line
// breaks.
func tagsSplit(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return r == ',' || r == '\t' || r == '\n' || r == '\r'
	})
}

// tagsTokenKey identifies a token by its text, so removing one tag leaves
// the others in place and only the removed token fades. A repeated tag
// keeps its place by its position among its equals.
type tagsTokenKey struct {
	tag  string
	same int
}

// tagsKeyAt returns the key of tags[i].
func tagsKeyAt(tags []string, i int) tagsTokenKey {
	same := 0
	for _, earlier := range tags[:i] {
		if earlier == tags[i] {
			same++
		}
	}
	return tagsTokenKey{tags[i], same}
}

// tagsCount shows how many of the allowed tags are in use.
func tagsCount(c *ui.Context, n, max int, refused bool) {
	t := theme.Colors(c)
	ink := t.MutedForeground
	if refused {
		ink = t.Destructive
	}
	count := style.Text(ui.Text(c, strconv.Itoa(n)+"/"+strconv.Itoa(max)), c, theme.XS()).
		SingleLine().Shrink(0).PaddingEnd(theme.Space(1))
	count.TextColor(motion.Color(count, "ink", ink))
}

// tagsToken is one tag with its remove button, a stop of the token group
// that Tab skips. The button keeps the focus where it was, so Backspace
// keeps removing tags from the editor. A refused duplicate tints its token.
func tagsToken(c *ui.Context, tag string, disabled, refused bool, remove func()) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().SM
	chip := ui.Row(c).AlignItems(ui.Center).Gap(theme.Space(0.5)).Shrink(0).MaxWidthPercent(100).
		PaddingStart(theme.Space(2)).PaddingEnd(theme.Space(0.5)).Height(theme.Scale(c, 22)).
		Radius(radius).FocusRing(false).
		Transition(ui.ElementTransition{
			Duration: theme.Motion().State, Ease: theme.Motion().Ease,
			Enter: &ui.Motion{}, Exit: &ui.Motion{},
		})
	tint := ui.Transparent
	if refused {
		tint = t.Destructive.Alpha(0.2)
	}
	tint = motion.Color(chip, "refused", tint)
	focused := chip.Focused()
	chip.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, t.Secondary, radius)
		p.Fill(r, tint, radius)
		if focused {
			p.Stroke(r, t.Ring, radius, 1.5)
		}
	})
	chip.Children(func() {
		style.Text(ui.Text(c, tag), c, theme.XS()).FontWeight(theme.WeightMedium).
			TextColor(t.SecondaryForeground).SingleLine().MinWidth(0)
		size := theme.Scale(c, 18)
		b := ui.ButtonBase(c).Size(size, size).Center().Shrink(0).Padding(0).
			FocusRing(false).TabStop(false).KeepFocus().
			Label(fmt.Sprintf(locale.M().RemoveNamed, tag)).TextColor(t.MutedForeground)
		b = style.Disabled(b, disabled)
		b = style.Background(b, radius-2, ui.Transparent, t.Accent, t.Accent)
		b.Children(func() {
			ui.Icon(c, lucide.X()).Size(theme.Scale(c, theme.IconXS), theme.Scale(c, theme.IconXS)).Shrink(0)
		})
		if !disabled {
			b.OnClick(remove)
		}
	})
	return chip
}
