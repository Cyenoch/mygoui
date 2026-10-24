package form

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TextAreaOptions configures HeroUI's multiline input. MinLines and MaxLines
// make the editor grow with its text to MaxLines, past which it scrolls.
type TextAreaOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	MinLines, MaxLines        int
	OnChange                  func()
}

// TextArea renders a multiline native editor on HeroUI's textarea surface.
func TextArea(c *ui.Context, value *string, o TextAreaOptions) ui.Element {
	value = inputTextAnchor(value)
	build := func() ui.Element {
		e := ui.TextAreaBase(c, value)
		e = fieldText(e, c)
		if o.Placeholder != "" {
			e.Placeholder(o.Placeholder)
		}
		if name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder); name != "" {
			e.Label(name)
		}
		if o.ReadOnly {
			e.ReadOnly(true)
		}
		if o.FullWidth {
			e.FillWidth()
		} else {
			e.MinWidth(theme.Scale(c, inputDefaultWidth))
		}
		e = e.PaddingX(theme.Space12).PaddingY(theme.Space8).Shrink(0)
		if o.MinLines > 0 || o.MaxLines > 0 {
			minLines := o.MinLines
			if minLines <= 0 {
				minLines = 1
			}
			maxLines := o.MaxLines
			if maxLines < minLines {
				maxLines = minLines
			}
			e.Lines(minLines, maxLines)
		} else {
			e.MinHeight(theme.Scale(c, theme.HeightTextarea))
		}
		e = fieldSurface(e, c, o.Variant, o.Invalid, o.Disabled)
		e.FocusRing(true)
		if o.OnChange != nil {
			e.OnChange(o.OnChange)
		}
		return e
	}

	fieldOptions := InputOptions{
		Label:       o.Label,
		Description: o.Description,
		Error:       o.Error,
		Required:    o.Required,
		Disabled:    o.Disabled,
		Invalid:     o.Invalid,
		FullWidth:   o.FullWidth,
	}
	if !inputNeedsField(fieldOptions) {
		return build()
	}
	return inputField(c, fieldOptions, true, func() { build() })
}
