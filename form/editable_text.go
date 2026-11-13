package form

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// EditableTextOptions configures text that edits in place. OnChange runs
// after a commit that changed the value.
type EditableTextOptions struct {
	Placeholder string
	AriaLabel   string
	Disabled    bool
	FullWidth   bool
	OnChange    func()
}

// editableState is the transient edit of one EditableText: the draft, and
// the handle that returns the focus to the text after a keyboard commit.
type editableState struct {
	draft   string
	editing bool
	focused bool
	text    ui.Handle
}

type editableStateKey struct{}

// EditableText shows *value as text. A click, Enter, Space or F2 turns it
// into an input of the same size; Enter or moving the focus away commits the
// draft to *value, and Escape drops it.
func EditableText(c *ui.Context, value *string, o EditableTextOptions) ui.Element {
	value = inputTextAnchor(value)
	root := ui.Row(c).MinWidth(0).Shrink(0)
	if o.FullWidth {
		root.FillWidth()
	} else {
		root.Width(theme.Scale(c, inputDefaultWidth))
	}
	st := ui.Local(root, editableStateKey{}, func() editableState { return editableState{} })
	if o.Disabled {
		st.editing = false
	}
	root.Children(func() {
		if st.editing {
			editableEditor(c, value, st, o)
			return
		}
		editableDisplay(c, value, st, o)
	})
	return root
}

// editableDisplay is the resting text: a button with the input's box, so the
// text stays in place when editing starts.
func editableDisplay(c *ui.Context, value *string, st *editableState, o EditableTextOptions) {
	t := theme.Colors(c)
	radius := theme.RadiusScale().Field
	b := ui.ButtonBase(c).Bind(&st.text).FillWidth().MinWidth(0).Justify(ui.Start).
		PaddingX(theme.Space(3)).PaddingY(theme.Scale(c, 5)).MinHeight(theme.Scale(c, 32)).
		Cursor(ui.CursorText).Label(inputAccessibleName("", o.AriaLabel, *value)).Value(*value)
	b = style.Disabled(b, o.Disabled)
	b = style.Background(b, radius, ui.Transparent, t.Accent, t.Accent)
	start := func() {
		st.draft, st.editing, st.focused = *value, true, false
	}
	b.OnClick(start).OnShortcut(0, ui.KeyF2, start)
	b.Children(func() {
		text, ink := *value, t.Foreground
		if text == "" {
			text, ink = o.Placeholder, t.FieldPlaceholder
		}
		fieldText(ui.Text(c, text), c).TextColor(ink).SingleLine().MinWidth(0)
	})
}

// editableEditor is the input holding the draft.
func editableEditor(c *ui.Context, value *string, st *editableState, o EditableTextOptions) {
	finish := func(commit bool) {
		st.editing = false
		if !commit || st.draft == *value {
			return
		}
		*value = st.draft
		if o.OnChange != nil {
			o.OnChange()
		}
	}
	e := inputEditor(c, &st.draft, inputCoreOptions{
		Placeholder: o.Placeholder,
		Name:        inputAccessibleName("", o.AriaLabel, o.Placeholder),
		FullWidth:   true,
		OnSubmit: func() {
			finish(true)
			st.text.Focus(c)
		},
	}).AutoFocus()
	e.OnTextKey(func(k ui.TextKey) bool {
		if k.Key != ui.KeyEscape {
			return false
		}
		finish(false)
		st.text.Focus(c)
		return true
	})
	// The focus arrives a frame after the input appears; only its loss
	// afterwards is a blur.
	switch focused := e.Focused(); {
	case focused:
		st.focused = true
	case st.focused:
		finish(true)
		c.Invalidate()
	}
}
