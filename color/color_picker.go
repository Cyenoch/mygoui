package color

import (
	"time"

	mygoui "github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ColorPickerOptions configures HeroUI's color picker: a trigger that opens a
// panel with a color area, channel sliders, a hex field and a palette, all
// editing the same ColorState.
type ColorPickerOptions struct {
	// Label names the trigger for assistive technology.
	Label string
	State *ColorState
	// Open is the caller-owned popover state. A nil Open leaves the picker
	// showing only its trigger.
	Open *bool
	// Placement places the panel; zero is below the trigger, aligned start.
	Placement mygoui.Placement
	// Swatches replaces the default palette. An empty, non-nil slice hides
	// the palette.
	Swatches    []ui.Color
	SwatchNames []string
	// HideArea, HideHue, HideAlpha, HideField and HideSwatches omit parts of
	// the panel.
	HideArea, HideHue, HideAlpha, HideField, HideSwatches bool
	Disabled                                              bool
	// OnChange reports that the user picked another color.
	OnChange func()
}

// defaultSwatchList is the shared palette the picker reads without copying;
// DefaultSwatches hands callers their own copy.
var defaultSwatchList = [...]ui.Color{
	ui.Hex("#ef4444"), ui.Hex("#f97316"), ui.Hex("#eab308"),
	ui.Hex("#22c55e"), ui.Hex("#14b8a6"), ui.Hex("#3b82f6"),
	ui.Hex("#6366f1"), ui.Hex("#a855f7"), ui.Hex("#ec4899"),
	ui.Hex("#71717a"), ui.Hex("#000000"), ui.Hex("#ffffff"),
}

// DefaultSwatches returns the palette a picker shows unless configured
// otherwise: a full hue ring plus black, white and a neutral.
func DefaultSwatches() []ui.Color {
	out := make([]ui.Color, len(defaultSwatchList))
	copy(out, defaultSwatchList[:])
	return out
}

// ColorPicker is HeroUI v3's ColorPicker. The trigger shows the current color
// and hex value; the panel combines the area, hue and alpha sliders, a hex
// field and a swatch palette, all synchronized through State.
func ColorPicker(c *ui.Context, o ColorPickerOptions) ui.Element {
	if c == nil || o.State == nil {
		return ui.Element{}
	}
	s := o.State
	label := o.Label
	if label == "" {
		label = "Color picker"
	}
	t := theme.Colors(c)
	// Keep the overlay layer present while closed so the panel's fade in and
	// out has a parent in both frames.
	ui.Overlay(c, func() {})
	root := ui.Row(c).Shrink(0).AlignItems(ui.Center)
	root.Children(func() {
		trigger := ui.ButtonBase(c).
			Radius(theme.RadiusScale().SM).Gap(theme.Space12).
			Role(ui.RoleColorWell).Label(label).
			Expanded(o.Open != nil && *o.Open).
			Value(s.Hex())
		trigger = style.Disabled(trigger, o.Disabled)
		trigger.OnClick(func() {
			if o.Open != nil && !o.Disabled {
				*o.Open = !*o.Open
			}
		})
		trigger.Children(func() {
			current := s.Value
			// HeroUI's default trigger swatch is MD: 32 DIP, 16 DIP round.
			ui.Box(c).Size(32, 32).Shrink(0).Radius(16).
				Role(ui.RoleNone).Clip().Draw(func(p *ui.Painter, r ui.Rect) {
				paintSwatch(p, r, current, 16)
				paintSwatchBorder(p, r, 16)
			})
			style.Text(ui.Text(c, s.Hex()), c, theme.SM()).
				TextColor(t.Foreground).SingleLine()
		})
		if o.Open == nil {
			return
		}
		panelRadius := theme.RadiusBase() * 2.5
		if panelRadius > 32 {
			panelRadius = 32
		}
		place := o.Placement
		if place.Side == mygoui.SideAuto && place.Align == mygoui.AlignCenter {
			place.Align = mygoui.AlignStart
		}
		ui.PopoverBase(c, trigger, o.Open, func(panel ui.Element) {
			panel = anchored.Apply(panel, trigger, place, mygoui.SideBottom)
			panel.Background(t.Overlay).Radius(panelRadius).
				Width(248).Shrink(0).Column().Gap(theme.Space12).
				Padding(theme.Space8, theme.Space8, theme.Space12, theme.Space8).
				Label(label)
			panel = style.Shadow(panel, c, theme.OverlayShadows(c.Theme().Dark))
			if !theme.ReduceMotion(c) {
				fade := ui.Motion{Opacity: 0}
				panel = panel.Transition(ui.ElementTransition{
					Duration: 150 * time.Millisecond, Ease: ui.EaseOut,
					Enter: &fade, Exit: &fade,
				})
			}
			panel.Children(func() {
				if !o.HideArea {
					// HeroUI's picker stretches the area to the panel's
					// content width (248 - 16 padding).
					ColorArea(c, ColorAreaOptions{
						Label: "Color area", State: s, Size: 232,
						Disabled: o.Disabled, OnChange: o.OnChange,
					})
				}
				if !o.HideHue {
					ColorSlider(c, ColorSliderOptions{
						Label: "Hue", State: s, Channel: ChannelHue,
						ShowOutput: true, Disabled: o.Disabled,
						OnChange: o.OnChange,
					})
				}
				if !o.HideAlpha {
					ColorSlider(c, ColorSliderOptions{
						Label: "Opacity", State: s, Channel: ChannelAlpha,
						ShowOutput: true, Disabled: o.Disabled,
						OnChange: o.OnChange,
					})
				}
				if !o.HideField {
					ColorField(c, ColorFieldOptions{
						Label: "Hex", State: s, FullWidth: true,
						Disabled: o.Disabled, OnChange: o.OnChange,
					})
				}
				if !o.HideSwatches {
					swatches := o.Swatches
					if swatches == nil {
						swatches = defaultSwatchList[:]
					}
					if len(swatches) > 0 {
						picked := s.Value
						ColorSwatchPicker(c, ColorSwatchPickerOptions{
							Label: "Swatches", Colors: swatches,
							Names: o.SwatchNames, Value: &picked,
							Size: SwatchSM, Disabled: o.Disabled,
							OnChange: func() {
								s.Set(picked)
								if o.OnChange != nil {
									o.OnChange()
								}
							},
						})
					}
				}
			})
		})
	})
	return root
}
