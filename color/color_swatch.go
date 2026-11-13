package color

import (
	"slices"
	"strings"

	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// swatchHoverScale is a picker swatch's hover face. The swatch shows the
// color itself, which a tint would change, so hover grows it instead, at once
// as every hover shows.
const swatchHoverScale = 1.1

// ColorSwatchSize is the swatch scale.
type ColorSwatchSize uint8

// Swatch sizes.
const (
	SwatchMD ColorSwatchSize = iota
	SwatchXS
	SwatchSM
	SwatchLG
	SwatchXL
)

// side returns the size's side in DIPs.
func (s ColorSwatchSize) side() float32 {
	switch s {
	case SwatchXS:
		return 16
	case SwatchSM:
		return 24
	case SwatchLG:
		return 36
	case SwatchXL:
		return 40
	}
	return 32
}

// ColorSwatchShape is the swatch shape.
type ColorSwatchShape uint8

// Swatch shapes.
const (
	SwatchCircle ColorSwatchShape = iota
	SwatchSquare
)

// swatchRadius returns the radius for a standalone swatch.
func swatchRadius(size ColorSwatchSize, shape ColorSwatchShape) float32 {
	if shape == SwatchSquare {
		return 6
	}
	switch size {
	case SwatchXS:
		return 8
	case SwatchSM:
		return 12
	case SwatchLG, SwatchXL:
		return 24
	}
	return 16
}

// ColorSwatchOptions configures a color chip. Label names the color for
// assistive technology; it defaults to the hex value.
type ColorSwatchOptions struct {
	Color    ui.Color
	Label    string
	Size     ColorSwatchSize
	Shape    ColorSwatchShape
	Selected bool
}

// ColorSwatch is a color chip: the color over a checkerboard for
// transparency, with a hairline border and a native accessible name.
func ColorSwatch(c *ui.Context, o ColorSwatchOptions) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	side := o.Size.side()
	radius := swatchRadius(o.Size, o.Shape)
	label := o.Label
	if label == "" {
		label = FormatHex(o.Color)
	}
	color := o.Color
	e := ui.Box(c).Size(side, side).Shrink(0).Radius(radius).Clip().
		Role(ui.RoleImage).Label(label)
	// An image has no selected state, so the selection is its description.
	if o.Selected {
		e.Description(locale.M().Selected)
	}
	// The tick draws in as the swatch is marked and fades out when unmarked.
	tick := motion.Progress(e, "selected", o.Selected)
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		paintSwatch(p, r, color, radius)
		paintSwatchBorder(p, r, radius)
		if tick > 0 {
			drawn := float32(1)
			if o.Selected {
				drawn = tick
			}
			ink := ui.RGB(255, 255, 255)
			if luminance(color) > 0.6 {
				ink = ui.RGB(0, 0, 0)
			}
			paintCheckmark(p, r, ink.Alpha(tick), drawn)
		}
	})
	return e
}

// ColorSwatchPickerLayout stacks the swatch grid or runs it as a column.
type ColorSwatchPickerLayout uint8

// Swatch picker layouts.
const (
	SwatchGrid ColorSwatchPickerLayout = iota
	SwatchStack
)

// ColorSwatchPickerOptions configures a single-choice palette.
type ColorSwatchPickerOptions struct {
	// Label names the group for assistive technology.
	Label  string
	Colors []ui.Color
	// Names optionally names each color in turn, such as "Ocean". Missing or
	// empty names fall back to the color's hex value.
	Names []string
	// Value is the caller-owned selection. A click or the arrow keys choose a
	// color into it. A nil Value keeps the choice inside the picker.
	Value    *ui.Color
	Size     ColorSwatchSize
	Shape    ColorSwatchShape
	Layout   ColorSwatchPickerLayout
	Disabled bool
	// OnChange reports that the user chose another color.
	OnChange func()
}

// ColorSwatchPicker is a palette of swatches. It is a native radio
// group: Tab reaches the chosen swatch, arrows rove among them and choose,
// and Space or Enter chooses the focused one.
func ColorSwatchPicker(c *ui.Context, o ColorSwatchPickerOptions) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	side := o.Size.side()
	itemRadius := pickerRadii(o.Size, o.Shape)
	border := float32(2)
	switch o.Size {
	case SwatchXS:
		border = 1
	case SwatchLG, SwatchXL:
		border = 3
	}
	label := o.Label
	if label == "" {
		label = "Swatches"
	}
	cursor := ui.CursorPointer
	if o.Disabled {
		cursor = ui.CursorDefault
	}
	// The border sits inside the item, so the swatch fills the content
	// box and grows on hover without escaping its item.
	content := side - 2*border
	group := ui.RadioGroup(c, func() {})
	// Without a caller's Value the picker keeps the chosen index itself.
	selected := ui.Local(group, swatchChoiceKey{}, func() int { return -1 })
	if o.Value != nil {
		*selected = slices.Index(o.Colors, *o.Value)
	}
	group.Children(func() {
		for i, col := range o.Colors {
			i, col := i, col
			name := pickerName(o, i, col)
			chosen := i == *selected
			borderColor := ui.Transparent
			if chosen {
				borderColor = col
			}
			item := ui.RadioBase(c.Key(colorItemKey(i)), selected, i).
				Size(side, side).Center().Shrink(0).Radius(itemRadius).
				FocusRing(true).
				Cursor(cursor).Label(name).Tooltip(name)
			item.Border(border, motion.Color(item, "border", borderColor))
			// Snapshot the item's hover in the build: the pose then animates
			// the swatch in paint-only frames, as the scale transition
			// does.
			hovered := item.Hovered() && !o.Disabled
			itemColor := col
			item.Children(func() {
				// The chosen swatch shrinks over the state motion; hover
				// grows it at once.
				chosenScale := float32(1)
				if chosen {
					chosenScale = 0.77
				}
				innerRadius := itemRadius
				if o.Shape == SwatchSquare && o.Size == SwatchSM && chosen {
					// The control rounds a selected small square swatch to 6 DIP
					// inside its 8 DIP item.
					innerRadius = 6
				}
				swatch := ui.Box(c).Size(content, content).Role(ui.RoleNone)
				if hovered && !chosen {
					swatch.Scale(swatchHoverScale, swatchHoverScale)
				}
				m := theme.Motion()
				swatch.AnimatePose(ui.Pose{Scale: chosenScale},
					ui.Animation{Duration: m.State, Ease: m.Ease})
				swatch.Draw(func(p *ui.Painter, r ui.Rect) {
					paintSwatch(p, r, itemColor, innerRadius)
				})

				// The tick is its own layer so its arrival (a scale and fade over
				// the state duration) stays centered in the item rather than the
				// swatch.
				ink := ui.RGB(255, 255, 255)
				if luminance(itemColor) > 0.5 {
					ink = ui.RGB(0, 0, 0)
				}
				checkScale, checkFade := float32(0.8), float32(1)
				if chosen {
					checkScale, checkFade = 1, 0
				}
				check := ui.Box(c).Absolute().Role(ui.RoleNone).
					Left(0).Top(0).Size(content, content)
				check.AnimatePose(ui.Pose{Scale: checkScale, Fade: checkFade},
					ui.Animation{Duration: m.State, Ease: m.Ease})
				check.Draw(func(p *ui.Painter, r ui.Rect) {
					paintCheckmark(p, r, ink, 1)
				})
			})
			item.OnChange(func() {
				if o.Value != nil {
					*o.Value = col
				}
				if o.OnChange != nil {
					o.OnChange()
				}
			})
		}
	}).Label(label).Gap(theme.Space(2)).AlignItems(ui.Center).
		Disabled(o.Disabled)
	if o.Layout == SwatchStack {
		group = group.Column().Shrink(0)
	} else {
		group = group.Row().Wrap().FillWidth()
	}
	return style.Disabled(group, o.Disabled)
}

// colorItemKey identifies one swatch of a palette across builds.
type colorItemKey int

// swatchChoiceKey keeps a picker's chosen index.
type swatchChoiceKey struct{}

// pickerRadii returns the item radius of a picker; the inner swatch inherits
// it, as the swatch does.
func pickerRadii(size ColorSwatchSize, shape ColorSwatchShape) (item float32) {
	if shape == SwatchSquare {
		switch size {
		case SwatchXS:
			return 6
		case SwatchSM:
			return 8
		}
		return 12
	}
	switch size {
	case SwatchXS:
		return 8
	case SwatchSM:
		return 12
	case SwatchLG, SwatchXL:
		return 24
	}
	return 16
}

// pickerName returns the accessible name of swatch i.
func pickerName(o ColorSwatchPickerOptions, i int, col ui.Color) string {
	if i < len(o.Names) {
		if name := strings.TrimSpace(o.Names[i]); name != "" {
			return name
		}
	}
	return FormatHex(col)
}
