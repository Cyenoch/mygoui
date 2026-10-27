package color

import (
	"github.com/Cyenoch/mygoui/form"
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ColorFieldOptions configures the hex color field. Its editing buffer
// lives in State, so one field edits one ColorState.
type ColorFieldOptions struct {
	// Label names the input; it defaults to "Hex".
	Label       string
	Description string
	// Error replaces the built-in invalid-hex message.
	Error       string
	Placeholder string
	State       *ColorState
	Disabled    bool
	Required    bool
	FullWidth   bool
	// OnChange reports that the text parsed as another color.
	OnChange func()
}

// ColorField is a hex color field. It accepts #rgb, #rgba, #rrggbb
// and #rrggbbaa; text that does not parse keeps the last valid color and shows
// an error announced with RoleStatus.
func ColorField(c *ui.Context, o ColorFieldOptions) ui.Element {
	if c == nil || o.State == nil {
		return ui.Element{}
	}
	s := o.State
	syncFieldText(s)
	label := o.Label
	if label == "" {
		label = "Hex"
	}
	invalid := o.Error != "" || s.field.invalid
	errorText := o.Error
	if errorText == "" && invalid {
		errorText = locale.M().Invalid
	}
	placeholder := o.Placeholder
	if placeholder == "" {
		placeholder = "#000000"
	}
	description := o.Description
	if invalid {
		description = ""
	}
	return form.Field(c, form.FieldOptions{
		Label: label, Description: description, Error: errorText,
		Required: o.Required, Disabled: o.Disabled, Invalid: invalid,
		FullWidth: o.FullWidth,
	}, func() {
		group := ui.Row(c).Height(theme.Scale(c, theme.HeightInput)).MinWidth(0).
			AlignItems(ui.Center).Gap(theme.Space8).PaddingX(theme.Space12)
		if o.FullWidth {
			group.FillWidth()
		}
		group = formstyle.Surface(group, c, false, invalid, o.Disabled)
		group.Children(func() {
			current := s.Value
			// The ColorField prefix is an XS swatch: 16 DIP, 8 DIP round.
			swatchRadius := theme.RadiusScale().LG
			swatchSize := theme.Scale(c, 16)
			ui.Box(c).Size(swatchSize, swatchSize).Shrink(0).
				Radius(swatchRadius).Clip().
				Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
				paintSwatch(p, r, current, swatchRadius)
				paintSwatchBorder(p, r, swatchRadius)
			})
			input := ui.TextInputBase(c, &s.field.text).Grow(1).MinWidth(theme.Scale(c, 120)).
				Placeholder(placeholder)
			input = formstyle.Text(input, c)
			input.OnChange(func() {
				parsed, err := ParseHex(s.field.text)
				if err != nil {
					s.field.invalid = true
					return
				}
				s.field.invalid = false
				if parsed == s.Value {
					s.field.seen = s.Value
					return
				}
				s.Set(parsed)
				s.field.seen = s.Value
				if o.OnChange != nil {
					o.OnChange()
				}
			})
		})
	})
}

// syncFieldText updates State's editing buffer when the color changed
// elsewhere.
func syncFieldText(s *ColorState) {
	if !s.field.ready || s.Value != s.field.seen {
		s.field.text = s.Hex()
		s.field.seen = s.Value
		s.field.ready = true
		s.field.invalid = false
	}
}
