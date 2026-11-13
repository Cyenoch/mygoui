package form

import (
	"github.com/Cyenoch/mygoui/internal/busy"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// InputOptions configures a single-line input. Label, Description and
// Error compose the field around it; without them the control is bare so it can
// sit inside a caller's own Field.
type InputOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth, Password       bool
	// Loading shows a spinner at the field's end. The field stays editable;
	// the caller owns the flag and the work behind it.
	Loading            bool
	OnChange, OnSubmit func()
}

// inputCoreOptions is the shared bare-editor configuration used by Input,
// TextField, InputGroup and SearchField.
type inputCoreOptions struct {
	Placeholder string
	Name        string
	Variant     ControlVariant
	Disabled    bool
	ReadOnly    bool
	Invalid     bool
	Password    bool
	FullWidth   bool
	OnChange    func()
	OnSubmit    func()
}

// inputDefaultWidth is the intrinsic width of a field that is not full width,
// close to the width a browser gives an input.
const inputDefaultWidth = 200

// inputTextAnchor returns the value an editor edits, tolerating a nil pointer
// by binding a throwaway empty string.
func inputTextAnchor(value *string) *string {
	if value != nil {
		return value
	}
	var empty string
	return &empty
}

// inputAccessibleName is the name assistive technology reads for a control
// without a visible field label. The placeholder is a last resort.
func inputAccessibleName(label, ariaLabel, placeholder string) string {
	switch {
	case ariaLabel != "":
		return ariaLabel
	case label != "":
		return label
	}
	return placeholder
}

// inputEditor builds the input surface around one native single-line
// editor. Hover, focus and invalid paint stay live through fieldSurface.
func inputEditor(c *ui.Context, value *string, o inputCoreOptions) ui.Element {
	e := ui.TextInputBase(c, value)
	e = fieldText(e, c)
	if o.Placeholder != "" {
		e.Placeholder(o.Placeholder)
	}
	if o.Name != "" {
		e.Label(o.Name)
	}
	if o.Password {
		e.Password()
	}
	if o.ReadOnly {
		e.ReadOnly(true)
	}
	if o.FullWidth {
		e.FillWidth()
	} else {
		e.MinWidth(theme.Scale(c, inputDefaultWidth))
	}
	e = e.PaddingX(theme.Space12).PaddingY(theme.Scale(c, 5)).
		MinHeight(theme.Scale(c, 32)).Shrink(0)
	e = fieldSurface(e, c, o.Variant, o.Invalid, o.Disabled)
	if o.OnChange != nil {
		e.OnChange(o.OnChange)
	}
	if o.OnSubmit != nil {
		e.OnSubmit(o.OnSubmit)
	}
	return e
}

// loadingInput is inputEditor with a trailing spinner while loading. The
// spinner is painted on the editor, which stays the one element callers size
// and focus; the end padding lasts from the flag to the end of the fade, so
// typed text never runs under the spinner.
func loadingInput(c *ui.Context, value *string, o inputCoreOptions, loading bool) ui.Element {
	editor := inputEditor(c, value, o).Busy(loading)
	shown := busy.Swap(editor, loading)
	if !loading && shown == 0 {
		return editor
	}
	size, inset := theme.Scale(c, busySize), theme.Scale(c, theme.Space12)
	editor.PaddingEnd(inset + size + theme.Scale(c, theme.Space8))
	ink := theme.Colors(c).MutedForeground.Alpha(shown)
	still := theme.ReduceMotion(c)
	rtl := c.Direction() == ui.RightToLeft
	// DrawOver also carries the focus ring, which this replaces.
	focused := editor.FocusWithin() && !o.Disabled
	ring := formstyle.FieldFace(c, focused, o.Invalid).Ring
	return editor.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focused {
			formstyle.FieldRing(p, r, theme.RadiusScale().Field, ring)
		}
		x := r.X + r.W - inset - size/2
		if rtl {
			x = r.X + inset + size/2
		}
		busy.PaintSpinner(p, x, r.Y+r.H/2, size, ink, still)
	})
}

// busySize is the spinner a field shows in place of a 14 DIP affordance.
const busySize = 14

// inputField composes label, description and error around an input. The
// description is hidden while invalid.
func inputField(c *ui.Context, o InputOptions, hideDescriptionOnInvalid bool, body func()) ui.Element {
	description := o.Description
	if hideDescriptionOnInvalid && o.Invalid {
		description = ""
	}
	return Field(c, FieldOptions{
		Label:       o.Label,
		Description: description,
		Error:       o.Error,
		Required:    o.Required,
		Disabled:    o.Disabled,
		Invalid:     o.Invalid,
		FullWidth:   o.FullWidth,
	}, body)
}

// inputNeedsField reports whether options ask for the field composition.
func inputNeedsField(o InputOptions) bool {
	return o.Label != "" || o.Description != "" || o.Error != "" || o.Invalid
}

// Input renders a single-line native editor: typing commits
// to value, and Enter submits without a fabricated action.
func Input(c *ui.Context, value *string, o InputOptions) ui.Element {
	// An error message implies the invalid state.
	o.Invalid = o.Invalid || o.Error != ""
	value = inputTextAnchor(value)
	build := func() ui.Element {
		return loadingInput(c, value, inputCoreOptions{
			Placeholder: o.Placeholder,
			Name:        inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder),
			Variant:     o.Variant,
			Disabled:    o.Disabled,
			ReadOnly:    o.ReadOnly,
			Invalid:     o.Invalid,
			Password:    o.Password,
			FullWidth:   o.FullWidth,
			OnChange:    o.OnChange,
			OnSubmit:    o.OnSubmit,
		}, o.Loading)
	}
	if !inputNeedsField(o) {
		return build()
	}
	return inputField(c, o, true, func() { build() })
}
