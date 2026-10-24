package form

import (
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

// InputGroup wraps one native editor in HeroUI's group shell. Prefix and
// suffix content is laid out in the shell; block regions span the full width.
func InputGroup(c *ui.Context, value *string, o InputGroupOptions) ui.Element {
	value = inputTextAnchor(value)
	t := theme.Colors(c)

	group := func() ui.Element {
		row := ui.Grid(c).ColumnTracks(ui.FitContent(), ui.Fr(1), ui.FitContent()).
			AlignItems(ui.Center).MinHeight(theme.Scale(c, theme.HeightInput)).Shrink(0)
		if o.FullWidth {
			row.FillWidth()
		} else {
			row.MinWidth(theme.Scale(c, inputDefaultWidth))
		}
		row = fieldSurface(row, c, o.Variant, o.Invalid, o.Disabled)

		var editor ui.Element
		row.Children(func() {
			// The editor comes first so a field label names and focuses it
			// even when an addon holds a button.
			paddingStart, paddingEnd := float32(theme.Space12), float32(theme.Space12)
			if o.Prefix != nil {
				paddingStart = 0
			}
			if o.Suffix != nil {
				paddingEnd = 0
			}
			editor = ui.TextInputBase(c, value).MinWidth(0).ColumnStart(2).FillWidth().
				PaddingStart(paddingStart).PaddingEnd(paddingEnd).PaddingY(theme.Space8)
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
				prefix := ui.Row(c).AlignItems(ui.Center).Shrink(0).ColumnStart(1).
					PaddingX(theme.Space12).MinHeight(theme.Scale(c, theme.HeightInput)).
					TextColor(t.FieldPlaceholder)
				prefix.Children(o.Prefix)
			}
			if o.Suffix != nil {
				suffix := ui.Row(c).AlignItems(ui.Center).Shrink(0).ColumnStart(3).
					PaddingX(theme.Space12).MinHeight(theme.Scale(c, theme.HeightInput)).
					TextColor(t.FieldPlaceholder)
				suffix.Children(o.Suffix)
			}
		})
		if editor.Valid() {
			focusVisible := editor.FocusVisible()
			radius := theme.RadiusScale().Field
			row.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if focusVisible {
					p.FocusRing(r, [4]float32{radius, radius, radius, radius})
				}
			})
		}
		return row
	}

	build := func() ui.Element {
		if o.BlockPrefix == nil && o.BlockSuffix == nil {
			return group()
		}
		// The group is built first and placed below a block prefix in a
		// grid: ui.Field names its first focusable descendant, so an addon
		// button in a block must not come before the editor in tree order.
		// Keyboard order follows the tree, ending at the addon actions.
		return ui.Grid(c).ColumnTracks(ui.Fr(1)).Gap(theme.Space4).
			FillWidth().MinWidth(0).Children(func() {
			groupRow := 1
			if o.BlockPrefix != nil {
				groupRow = 2
			}
			group().RowStart(groupRow)
			if o.BlockPrefix != nil {
				ui.Column(c).FillWidth().MinWidth(0).RowStart(1).Children(o.BlockPrefix)
			}
			if o.BlockSuffix != nil {
				ui.Column(c).FillWidth().MinWidth(0).RowStart(groupRow + 1).Children(o.BlockSuffix)
			}
		})
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
