package form

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// InputOptions configures a single-line HeroUI input. Label, Description and
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
	OnChange, OnSubmit        func()
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

// inputEditor builds the HeroUI input surface around one native single-line
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
	e = e.PaddingX(theme.Space12).PaddingY(theme.Space8).
		MinHeight(theme.Scale(c, theme.HeightInput)).Shrink(0)
	e = fieldSurface(e, c, o.Variant, o.Invalid, o.Disabled)
	if o.OnChange != nil {
		e.OnChange(o.OnChange)
	}
	if o.OnSubmit != nil {
		e.OnSubmit(o.OnSubmit)
	}
	return e
}

// inputField composes label, description and error around an input. The
// description is hidden while invalid, as HeroUI's field roots do.
func inputField(c *ui.Context, o InputOptions, hideDescriptionOnInvalid bool, body func()) ui.Element {
	description := o.Description
	if hideDescriptionOnInvalid && (o.Invalid || o.Error != "") {
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

// Input renders HeroUI's single-line input on a native editor: typing commits
// to value, and Enter submits without a fabricated action.
func Input(c *ui.Context, value *string, o InputOptions) ui.Element {
	value = inputTextAnchor(value)
	build := func() ui.Element {
		return inputEditor(c, value, inputCoreOptions{
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
		})
	}
	if !inputNeedsField(o) {
		return build()
	}
	return inputField(c, o, true, func() { build() })
}
