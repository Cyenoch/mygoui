package form

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SearchFieldOptions configures HeroUI's search field. Escape and the clear
// button empty the value; Enter reports OnSubmit.
type SearchFieldOptions struct {
	Label, Description, Error string
	Placeholder               string
	AriaLabel                 string
	Variant                   ControlVariant
	Required, Disabled        bool
	ReadOnly, Invalid         bool
	FullWidth                 bool
	OnChange, OnSubmit        func()
	OnClear                   func()
}

// SearchField renders a native search input with the search icon and clear
// action inside HeroUI's group shell.
func SearchField(c *ui.Context, value *string, o SearchFieldOptions) ui.Element {
	value = inputTextAnchor(value)
	t := theme.Colors(c)
	radius := theme.RadiusScale().Field

	build := func() ui.Element {
		row := ui.Row(c).AlignItems(ui.Center).MinHeight(theme.Scale(c, theme.HeightInput)).Shrink(0)
		if o.FullWidth {
			row.FillWidth()
		} else {
			row.MinWidth(theme.Scale(c, inputDefaultWidth))
		}
		row = fieldSurface(row, c, o.Variant, o.Invalid, o.Disabled)

		var editor ui.Element
		row.Children(func() {
			ui.Icon(c, lucide.Search()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
				Shrink(0).MarginStart(theme.Space12).TextColor(t.FieldPlaceholder)

			editor = ui.TextInputBase(c, value).Grow(1).MinWidth(0).
				PaddingStart(theme.Space8).PaddingEnd(theme.Space8).PaddingY(theme.Space8)
			editor = fieldText(editor, c)
			if o.Placeholder != "" {
				editor.Placeholder(o.Placeholder)
			}
			name := inputAccessibleName(o.Label, o.AriaLabel, o.Placeholder)
			if name == "" {
				name = locale.M().Search
			}
			editor.Label(name)
			if o.ReadOnly {
				editor.ReadOnly(true)
			}
			if o.FullWidth {
				editor.FillWidth()
			}
			editor.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEscape || o.Disabled || o.ReadOnly || k.Text() == "" {
					return false
				}
				k.SetText("")
				if o.OnClear != nil {
					o.OnClear()
				}
				if o.OnChange != nil {
					o.OnChange()
				}
				return true
			})
			if o.OnChange != nil {
				editor.OnChange(o.OnChange)
			}
			if o.OnSubmit != nil {
				editor.OnSubmit(o.OnSubmit)
			}

			clear := ui.ButtonBase(c).Size(theme.Scale(c, theme.Space20), theme.Scale(c, theme.Space20)).
				Center().Shrink(0).Padding(0).MarginEnd(theme.Space8).
				Radius(theme.RadiusScale().SM).Label(locale.M().Clear).Tooltip(locale.M().Clear)
			clear.TextColor(t.FieldPlaceholder)
			clear = style.Background(clear, theme.RadiusScale().SM, ui.Transparent, t.DefaultHover, t.DefaultHover)
			clear.Children(func() {
				ui.Icon(c, lucide.X()).Size(theme.Scale(c, theme.IconXS), theme.Scale(c, theme.IconXS)).Shrink(0)
			})
			switch {
			case *value == "":
				clear.Opacity(0)
			case o.Disabled || o.ReadOnly:
				style.Disabled(clear, true)
			default:
				clear.OnClick(func() {
					*value = ""
					if o.OnClear != nil {
						o.OnClear()
					}
					if o.OnChange != nil {
						o.OnChange()
					}
				})
			}
		})

		if editor.Valid() {
			focusVisible := editor.FocusVisible()
			row.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if focusVisible {
					p.FocusRing(r, [4]float32{radius, radius, radius, radius})
				}
			})
		}
		return row
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
