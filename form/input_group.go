package form

import (
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// InputGroupOptions configures one input with inline Prefix and Suffix
// regions, or BlockPrefix and BlockSuffix rows above and below it. The
// callbacks build ordinary elements; actions inside them keep their own
// activation.
type InputGroupOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	Prefix, Suffix            func()
	BlockPrefix, BlockSuffix  func()
	OnChange, OnSubmit        func()
}

// InputGroup wraps one native editor in a group shell. Clicking a
// non-interactive prefix or suffix focuses the editor; actions keep their
// own activation. Block regions span the full width.
func InputGroup(c *ui.Context, value *string, o InputGroupOptions) ui.Element {
	value = inputTextAnchor(value)
	t := theme.Colors(c)

	build := func() ui.Element {
		var editor ui.Element
		focusEditor := func() {
			if !o.Disabled {
				editor.Focus()
			}
		}
		row := 1
		if o.BlockPrefix != nil {
			row = 2
		}
		shell := ui.Grid(c).ColumnTracks(ui.FitContent(), ui.Fr(1), ui.FitContent()).
			AlignItems(ui.Center).MinHeight(theme.Scale(c, 32)).Shrink(0)
		if o.FullWidth {
			shell.FillWidth()
		} else {
			shell.MinWidth(theme.Scale(c, inputDefaultWidth))
		}
		shell.Children(func() {
			paddingStart, paddingEnd := float32(theme.Space12), float32(theme.Space12)
			if o.Prefix != nil {
				paddingStart = theme.Space8
			}
			if o.Suffix != nil {
				paddingEnd = theme.Space8
			}
			editor = ui.TextInputBase(c, value).MinWidth(0).ColumnStart(2).RowStart(row).JustifySelf(ui.Stretch).
				PaddingStart(paddingStart).PaddingEnd(paddingEnd).PaddingY(theme.Scale(c, 5))
			editor = fieldText(editor, c)
			if o.Placeholder != "" {
				editor.Placeholder(o.Placeholder)
			}
			if name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder); name != "" {
				editor.Label(name)
			}
			if o.ReadOnly {
				editor.ReadOnly(true)
			}
			if o.OnChange != nil {
				editor.OnChange(o.OnChange)
			}
			if o.OnSubmit != nil {
				editor.OnSubmit(o.OnSubmit)
			}
			if o.Prefix != nil {
				prefix := ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).Shrink(0).ColumnStart(1).RowStart(row).
					PaddingStart(theme.Space12).MinHeight(theme.Scale(c, 32)).TextColor(t.MutedForeground).
					OnClick(focusEditor)
				prefix.Children(o.Prefix)
			}
			if o.Suffix != nil {
				suffix := ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).Shrink(0).ColumnStart(3).RowStart(row).
					PaddingEnd(theme.Space12).MinHeight(theme.Scale(c, 32)).TextColor(t.MutedForeground).
					OnClick(focusEditor)
				suffix.Children(o.Suffix)
			}
			if o.BlockPrefix != nil {
				ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Justify(ui.Start).
					Gap(theme.Space8).ColumnStart(1).ColumnSpan(3).RowStart(1).
					Padding(theme.Space8, theme.Space12, 0, theme.Space12).OnClick(focusEditor).Children(o.BlockPrefix)
			}
			if o.BlockSuffix != nil {
				ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Justify(ui.Start).
					Gap(theme.Space8).ColumnStart(1).ColumnSpan(3).RowStart(row+1).
					Padding(0, theme.Space12, theme.Space8, theme.Space12).OnClick(focusEditor).Children(o.BlockSuffix)
			}
		})
		return formstyle.SurfaceFocused(shell, c, editor.Focused(), o.Invalid || o.Error != "", o.Disabled)
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
