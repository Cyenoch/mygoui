package form

import (
	"github.com/egoist/mygo/ui"
)

// TextFieldOptions describes HeroUI's complete text field: a label,
// description, error and one single-line input around the caller's value.
type TextFieldOptions struct {
	Value                     *string
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth, Password       bool
	OnChange, OnSubmit        func()
}

// TextField composes a labelled single-line field. Typing writes through to
// Value; Enter reports OnSubmit and never invents a submission.
func TextField(c *ui.Context, o TextFieldOptions) ui.Element {
	value := inputTextAnchor(o.Value)
	input := InputOptions{
		Label:       o.Label,
		Description: o.Description,
		Error:       o.Error,
		Required:    o.Required,
		Disabled:    o.Disabled,
		Invalid:     o.Invalid,
		FullWidth:   o.FullWidth,
	}
	return inputField(c, input, true, func() {
		inputEditor(c, value, inputCoreOptions{
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
	})
}
