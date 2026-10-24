package gallery

import (
	"github.com/egoist/mygo/ui"
)

// colorsComponentState keeps every color page editing the one shared
// ColorState, so area, sliders, field, swatches and picker stay in step.
func colorsComponentState(s *State) *colorsState {
	st := StateOf[colorsState](s, "colors")
	st.ensure()
	return st
}

func init() {
	RegisterComponent(ComponentPage{
		ID: "colorarea", Title: "ColorArea", Group: "Color", Order: 10,
		Description: "A two-dimensional saturation and brightness field for one shared color.",
		Examples: []ComponentExample{
			{Title: "Large with dots", Description: "The full saturation and brightness field at 200 DIP.",
				Demo: func(c *ui.Context, s *State) { colorAreaLargeExample(c, colorsComponentState(s)) }},
			{Title: "Compact", Description: "A smaller field without the dot pattern.",
				Demo: func(c *ui.Context, s *State) { colorAreaSmallExample(c, colorsComponentState(s)) }},
			{Title: "Disabled", Description: "A read-only area keeps the current color.",
				Demo: func(c *ui.Context, s *State) { colorAreaDisabledExample(c, colorsComponentState(s)) }},
		},
		Code: `color.ColorArea(c, color.ColorAreaOptions{
    Label: "Fill color", State: &state,
    ShowDots: true, Size: 200,
})`,
		Params: []ComponentParam{
			{"Label", "string", "\"Color area\"", "Names the area for assistive technology."},
			{"State", "*color.ColorState", "required", "The color edited and shown."},
			{"ShowDots", "bool", "false", "Overlays the subtle dot pattern."},
			{"Size", "float32", "224", "Square side in DIPs."},
			{"Disabled", "bool", "false", "Makes the area read-only."},
			{"OnChange", "func()", "nil", "Reports a newly picked color."},
		},
		Tokens: []string{"Foreground", "Focus", "FieldBackground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "colorslider", Title: "ColorSlider", Group: "Color", Order: 20,
		Description: "One color channel per gradient track, editing the shared ColorState.",
		Examples: []ComponentExample{
			{Title: "Hue", Description: "One gradient track edits the hue channel.",
				Demo: func(c *ui.Context, s *State) { colorSliderHueExample(c, colorsComponentState(s)) }},
			{Title: "Saturation and brightness", Description: "Two channels stay synchronized with the same color.",
				Demo: func(c *ui.Context, s *State) { colorSliderChannelsExample(c, colorsComponentState(s)) }},
			{Title: "Opacity", Description: "The alpha track shows the transparency range.",
				Demo: func(c *ui.Context, s *State) { colorSliderOpacityExample(c, colorsComponentState(s)) }},
			{Title: "Vertical", Description: "The hue track runs from bottom to top.",
				Demo: func(c *ui.Context, s *State) { colorSliderVerticalExample(c, colorsComponentState(s)) }},
		},
		Code: `color.ColorSlider(c, color.ColorSliderOptions{
    Label: "Hue", State: &state,
    Channel: color.ChannelHue, ShowOutput: true,
})`,
		Params: []ComponentParam{
			{"Label", "string", "channel name", "Names the slider for assistive technology."},
			{"State", "*color.ColorState", "required", "The color edited by the channel."},
			{"Channel", "color.ColorChannel", "ChannelHue", "Hue, saturation, brightness or alpha."},
			{"ShowOutput", "bool", "false", "Shows the channel value beside the label."},
			{"Vertical", "bool", "false", "Runs the track vertically."},
			{"Height", "float32", "160", "Vertical slider length in DIPs."},
		},
		Tokens: []string{"Foreground", "Focus", "FieldBackground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "colorfield", Title: "ColorField", Group: "Color", Order: 30,
		Description: "A hex editor that parses #rgb, #rgba, #rrggbb and #rrggbbaa text.",
		Examples: []ComponentExample{
			{Title: "Full width", Description: "The editor grows to its parent with help text below.",
				Demo: func(c *ui.Context, s *State) { colorFieldFullExample(c, colorsComponentState(s)) }},
			{Title: "Compact", Description: "The intrinsic-width hex editor.",
				Demo: func(c *ui.Context, s *State) { colorFieldCompactExample(c, colorsComponentState(s)) }},
			{Title: "Disabled", Description: "A read-only field keeps the current color.",
				Demo: func(c *ui.Context, s *State) { colorFieldDisabledExample(c, colorsComponentState(s)) }},
		},
		Code: `color.ColorField(c, color.ColorFieldOptions{
    Label: "Hex value", State: &state, FullWidth: true,
    Description: "Edits the same color as every other control",
})`,
		Params: []ComponentParam{
			{"Label", "string", "\"Hex\"", "Names the input."},
			{"Description", "string", "\"\"", "Help text shown below the field."},
			{"State", "*color.ColorState", "required", "The color edited by the text."},
			{"FullWidth", "bool", "false", "Grows the field to its parent's width."},
			{"Placeholder", "string", "\"#000000\"", "Text shown while the field is empty."},
			{"Disabled", "bool", "false", "Makes the field read-only."},
		},
		Tokens: []string{"FieldBackground", "FieldHover", "FieldFocus", "FieldBorder", "Danger"},
	})

	RegisterComponent(ComponentPage{
		ID: "colorswatch", Title: "ColorSwatch", Group: "Color", Order: 40,
		Description: "A color chip over a transparency checkerboard, in sizes and shapes.",
		Examples: []ComponentExample{
			{Title: "Sizes", Description: "XS through XL side by side.",
				Demo: func(c *ui.Context, s *State) { colorSwatchSizesExample(c, colorsComponentState(s)) }},
			{Title: "Shapes", Description: "Circles and square corners at three sizes.",
				Demo: func(c *ui.Context, s *State) { colorSwatchShapesExample(c, colorsComponentState(s)) }},
			{Title: "Transparency and selection", Description: "A checkerboard shows alpha, and the checkmark marks the choice.",
				Demo: func(c *ui.Context, s *State) { colorSwatchSelectedExample(c, colorsComponentState(s)) }},
		},
		Code: `color.ColorSwatch(c, color.ColorSwatchOptions{
    Color: state.Value, Size: color.SwatchLG,
    Shape: color.SwatchCircle,
})`,
		Params: []ComponentParam{
			{"Color", "ui.Color", "required", "The color shown by the chip."},
			{"Label", "string", "hex value", "Names the color for assistive technology."},
			{"Size", "color.ColorSwatchSize", "SwatchMD", "XS, SM, MD, LG or XL."},
			{"Shape", "color.ColorSwatchShape", "SwatchCircle", "Circle or square corners."},
			{"Selected", "bool", "false", "Draws the selected checkmark."},
		},
		Tokens: []string{"FieldBackground", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "colorswatchpicker", Title: "ColorSwatchPicker", Group: "Color", Order: 50,
		Description: "A single-choice palette that writes its selection into the shared color.",
		Examples: []ComponentExample{
			{Title: "Circle", Description: "The default circular single-choice palette.",
				Demo: func(c *ui.Context, s *State) { colorSwatchPickerCircleExample(c, colorsComponentState(s)) }},
			{Title: "Square", Description: "Square swatches at a smaller size.",
				Demo: func(c *ui.Context, s *State) { colorSwatchPickerSquareExample(c, colorsComponentState(s)) }},
			{Title: "Read-only", Description: "A nil value leaves the palette unselected.",
				Demo: func(c *ui.Context, s *State) { colorSwatchPickerReadOnlyExample(c, colorsComponentState(s)) }},
		},
		Code: `color.ColorSwatchPicker(c, color.ColorSwatchPickerOptions{
    Label: "Brand palette", Colors: color.DefaultSwatches(),
    Value: &state.Value, Size: color.SwatchLG,
})`,
		Params: []ComponentParam{
			{"Label", "string", "\"Swatches\"", "Names the palette for assistive technology."},
			{"Colors", "[]ui.Color", "required", "The colors offered by the palette."},
			{"Value", "*ui.Color", "nil", "Caller-owned selection; nil is read-only."},
			{"Size", "color.ColorSwatchSize", "SwatchMD", "Swatch side of every item."},
			{"Shape", "color.ColorSwatchShape", "SwatchCircle", "Circle or square corners."},
			{"OnChange", "func()", "nil", "Reports a newly chosen color."},
		},
		Tokens: []string{"FieldBackground", "Focus", "Border"},
	})

	RegisterComponent(ComponentPage{
		ID: "colorpicker", Title: "ColorPicker", Group: "Color", Order: 60,
		Description: "A trigger whose popup combines area, sliders, hex field and palette.",
		Examples: []ComponentExample{
			{Title: "Default palette", Description: "The trigger shows the live color while the full panel edits it.",
				Demo: func(c *ui.Context, s *State) { colorPickerDefaultExample(c, colorsComponentState(s)) }},
			{Title: "Compact panel", Description: "The area, alpha track and hex field are omitted.",
				Demo: func(c *ui.Context, s *State) { colorPickerCompactExample(c, colorsComponentState(s)) }},
			{Title: "Without swatches", Description: "An empty swatch list hides the palette.",
				Demo: func(c *ui.Context, s *State) { colorPickerNoSwatchesExample(c, colorsComponentState(s)) }},
		},
		Code: `color.ColorPicker(c, color.ColorPickerOptions{
    Label: "Theme color", State: &state, Open: &open,
    Swatches: color.DefaultSwatches(),
})`,
		Params: []ComponentParam{
			{"Label", "string", "\"Color picker\"", "Names the trigger for assistive technology."},
			{"State", "*color.ColorState", "required", "The color edited by the panel."},
			{"Open", "*bool", "nil", "Caller-owned popover state."},
			{"Swatches", "[]ui.Color", "DefaultSwatches", "Palette shown in the panel; empty hides it."},
			{"Placement", "mygoui.Placement", "below start", "Placement of the panel."},
			{"OnChange", "func()", "nil", "Reports a newly picked color."},
		},
		Tokens: []string{"Overlay", "Border", "Foreground", "FieldBackground", "Focus"},
	})
}
