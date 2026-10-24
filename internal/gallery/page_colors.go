package gallery

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/color"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// colorsState is the single color every control on the page edits, plus the
// popup state of the full picker.
type colorsState struct {
	ready             bool
	Color             color.ColorState
	PickerOpen        bool
	PickerOpenCompact bool
	PickerOpenBare    bool
}

func (st *colorsState) ensure() {
	if st.ready {
		return
	}
	st.ready = true
	st.Color = *color.NewColorState(ui.Hex("#6366f1"))
}

// colorAreaLargeExample is the full saturation/brightness field.
func colorAreaLargeExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		color.ColorArea(c, color.ColorAreaOptions{
			Label: "Fill color", State: &st.Color, ShowDots: true,
			Size: theme.Scale(c, 200),
		})
		ui.Column(c).Grow(1).MinWidth(theme.Scale(c, 160)).Gap(theme.Space8).Children(func() {
			colorsValue(c, st)
		})
	})
}

// colorAreaSmallExample trades the dot pattern for a compact square.
func colorAreaSmallExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		color.ColorArea(c, color.ColorAreaOptions{
			Label: "Small fill color", State: &st.Color, Size: theme.Scale(c, 128),
		})
		colorsValue(c, st)
	})
}

// colorAreaDisabledExample is read-only.
func colorAreaDisabledExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		color.ColorArea(c, color.ColorAreaOptions{
			Label: "Locked fill color", State: &st.Color, Disabled: true,
			Size: theme.Scale(c, 160),
		})
		mygoui.Typography(c, "A disabled area keeps the current color.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorSliderHueExample edits only the hue channel.
func colorSliderHueExample(c *ui.Context, st *colorsState) {
	ui.Column(c).MaxWidth(theme.Scale(c, 360)).Gap(theme.Space16).Children(func() {
		color.ColorSlider(c, color.ColorSliderOptions{
			Label: "Hue", State: &st.Color, Channel: color.ChannelHue, ShowOutput: true,
		})
		colorsValue(c, st)
	})
}

// colorSliderChannelsExample pairs saturation with brightness.
func colorSliderChannelsExample(c *ui.Context, st *colorsState) {
	ui.Column(c).MaxWidth(theme.Scale(c, 360)).Gap(theme.Space16).Children(func() {
		color.ColorSlider(c, color.ColorSliderOptions{
			Label: "Saturation", State: &st.Color, Channel: color.ChannelSaturation, ShowOutput: true,
		})
		color.ColorSlider(c, color.ColorSliderOptions{
			Label: "Brightness", State: &st.Color, Channel: color.ChannelBrightness, ShowOutput: true,
		})
		colorsValue(c, st)
	})
}

// colorSliderOpacityExample edits the alpha channel.
func colorSliderOpacityExample(c *ui.Context, st *colorsState) {
	ui.Column(c).MaxWidth(theme.Scale(c, 360)).Gap(theme.Space16).Children(func() {
		color.ColorSlider(c, color.ColorSliderOptions{
			Label: "Opacity", State: &st.Color, Channel: color.ChannelAlpha, ShowOutput: true,
		})
		colorsValue(c, st)
	})
}

// colorSliderVerticalExample runs the hue track vertically.
func colorSliderVerticalExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space24).AlignItems(ui.Start).Children(func() {
		color.ColorSlider(c, color.ColorSliderOptions{
			Label: "Vertical hue", State: &st.Color, Channel: color.ChannelHue,
			ShowOutput: true, Vertical: true, Height: theme.Scale(c, 180),
		})
		colorsValue(c, st)
	})
}

// colorFieldFullExample grows the hex editor to its parent.
func colorFieldFullExample(c *ui.Context, st *colorsState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		color.ColorField(c, color.ColorFieldOptions{
			Label: "Hex value", State: &st.Color, FullWidth: true,
			Description: "Edits the same color as the area, sliders and palette",
		})
		colorsValue(c, st)
	})
}

// colorFieldCompactExample is the intrinsic-width editor.
func colorFieldCompactExample(c *ui.Context, st *colorsState) {
	color.ColorField(c, color.ColorFieldOptions{
		Label: "Hex value", State: &st.Color, Description: "Accepts #rgb, #rgba, #rrggbb and #rrggbbaa",
	})
}

// colorFieldDisabledExample is read-only.
func colorFieldDisabledExample(c *ui.Context, st *colorsState) {
	color.ColorField(c, color.ColorFieldOptions{
		Label: "Locked hex", State: &st.Color, Disabled: true,
		Description: "A disabled field keeps the current color",
	})
}

// colorSwatchSizesExample shows every swatch size.
func colorSwatchSizesExample(c *ui.Context, st *colorsState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		ui.Row(c).Wrap().Gap(theme.Space8).AlignItems(ui.Center).Children(func() {
			for _, size := range []color.ColorSwatchSize{
				color.SwatchXS, color.SwatchSM, color.SwatchMD,
				color.SwatchLG, color.SwatchXL,
			} {
				color.ColorSwatch(c, color.ColorSwatchOptions{Color: st.Color.Value, Size: size})
			}
		})
		mygoui.Typography(c, "XS through XL at "+st.Color.Hex(), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorSwatchShapesExample compares circles and square corners.
func colorSwatchShapesExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
		for _, shape := range []color.ColorSwatchShape{color.SwatchCircle, color.SwatchSquare} {
			for _, size := range []color.ColorSwatchSize{color.SwatchSM, color.SwatchMD, color.SwatchLG} {
				color.ColorSwatch(c, color.ColorSwatchOptions{
					Color: st.Color.Value, Shape: shape, Size: size,
				})
			}
		}
	})
}

// colorSwatchSelectedExample shows transparency and the selected checkmark.
func colorSwatchSelectedExample(c *ui.Context, _ *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space12).AlignItems(ui.Center).Children(func() {
		color.ColorSwatch(c, color.ColorSwatchOptions{
			Color: ui.Color{R: 0x63, G: 0x66, B: 0xf1, A: 0x66},
			Label: "Translucent swatch", Size: color.SwatchLG, Selected: true,
		})
		color.ColorSwatch(c, color.ColorSwatchOptions{
			Color: ui.Hex("#0ea5e9"), Label: "Opaque swatch",
			Size: color.SwatchLG, Shape: color.SwatchSquare,
		})
		mygoui.Typography(c, "The checkmark marks the current color.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorSwatchPickerCircleExample is the default circular palette.
func colorSwatchPickerCircleExample(c *ui.Context, st *colorsState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		color.ColorSwatchPicker(c, color.ColorSwatchPickerOptions{
			Label:  "Brand palette",
			Colors: color.DefaultSwatches(),
			Value:  &st.Color.Value,
			Size:   color.SwatchLG,
			Shape:  color.SwatchCircle,
		})
		mygoui.Typography(c, "Selected "+st.Color.Hex(), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorSwatchPickerSquareExample uses square corners and a smaller size.
func colorSwatchPickerSquareExample(c *ui.Context, st *colorsState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		color.ColorSwatchPicker(c, color.ColorSwatchPickerOptions{
			Label:  "Square palette",
			Colors: color.DefaultSwatches()[:6],
			Value:  &st.Color.Value,
			Size:   color.SwatchMD,
			Shape:  color.SwatchSquare,
		})
		mygoui.Typography(c, "Selected "+st.Color.Hex(), mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorSwatchPickerReadOnlyExample has no selected value.
func colorSwatchPickerReadOnlyExample(c *ui.Context, _ *colorsState) {
	ui.Column(c).Gap(theme.Space12).Children(func() {
		color.ColorSwatchPicker(c, color.ColorSwatchPickerOptions{
			Label:  "Read-only palette",
			Colors: color.DefaultSwatches()[:6],
			Size:   color.SwatchMD,
		})
		mygoui.Typography(c, "A nil Value makes the palette read-only.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorPickerDefaultExample opens the full panel.
func colorPickerDefaultExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		color.ColorPicker(c, color.ColorPickerOptions{
			Label: "Theme color", State: &st.Color, Open: &st.PickerOpen,
		})
		colorsValue(c, st)
	})
}

// colorPickerCompactExample hides the area and alpha channel.
func colorPickerCompactExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		color.ColorPicker(c, color.ColorPickerOptions{
			Label: "Compact color", State: &st.Color, Open: &st.PickerOpenCompact,
			HideArea: true, HideAlpha: true, HideField: true,
			Swatches: color.DefaultSwatches()[:8],
		})
		mygoui.Typography(c, "A shorter panel keeps the hue track and a small palette.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

// colorPickerNoSwatchesExample hides the palette.
func colorPickerNoSwatchesExample(c *ui.Context, st *colorsState) {
	ui.Row(c).Wrap().Gap(theme.Space16).AlignItems(ui.Center).Children(func() {
		color.ColorPicker(c, color.ColorPickerOptions{
			Label: "Color without palette", State: &st.Color, Open: &st.PickerOpenBare,
			Swatches: []ui.Color{},
		})
		mygoui.Typography(c, "An empty swatch list hides the palette entirely.", mygoui.TypographyOptions{
			Type: theme.TypeBodyXS, Muted: true,
		})
	})
}

func colorsValue(c *ui.Context, st *colorsState) {
	ui.Row(c).AlignItems(ui.Center).Gap(theme.Space8).Children(func() {
		mygoui.Typography(c, "Hex", mygoui.TypographyOptions{Type: theme.TypeBodyXS, Muted: true})
		mygoui.Code(c, st.Color.Hex())
	})
}
