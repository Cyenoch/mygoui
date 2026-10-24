package color

import (
	"strings"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ColorSwatchSize is HeroUI's swatch scale.
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

// ColorSwatchShape is HeroUI's swatch shape.
type ColorSwatchShape uint8

// Swatch shapes.
const (
	SwatchCircle ColorSwatchShape = iota
	SwatchSquare
)

// swatchRadius returns HeroUI's radius for a standalone swatch.
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

// ColorSwatch is HeroUI v3's color chip: the color over a checkerboard for
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
	e.Draw(func(p *ui.Painter, r ui.Rect) {
		paintSwatch(p, r, color, radius)
		paintSwatchBorder(p, r, radius)
		if o.Selected {
			ink := ui.RGB(255, 255, 255)
			if luminance(color) > 0.6 {
				ink = ui.RGB(0, 0, 0)
			}
			paintCheckmark(p, r, ink)
		}
	})
	return e
}

// ColorSwatchPickerLayout stacks HeroUI's swatch grid or runs it as a column.
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
	// color into it. A nil Value makes the picker read-only.
	Value    *ui.Color
	Size     ColorSwatchSize
	Shape    ColorSwatchShape
	Layout   ColorSwatchPickerLayout
	Disabled bool
	// OnChange reports that the user chose another color.
	OnChange func()
}

// ColorSwatchPicker is HeroUI v3's palette of swatches. It is a native radio
// group: Tab reaches the chosen swatch, arrows rove among them and choose,
// and Space or Enter chooses the focused one.
func ColorSwatchPicker(c *ui.Context, o ColorSwatchPickerOptions) ui.Element {
	if c == nil {
		return ui.Element{}
	}
	selected := -1
	if o.Value != nil {
		for i, col := range o.Colors {
			if col == *o.Value {
				selected = i
				break
			}
		}
	}
	side := o.Size.side()
	itemRadius, swatchRadiusValue := pickerRadii(o.Size, o.Shape)
	fieldBg := theme.Colors(c).FieldBackground
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
	group := ui.RadioGroup(c, func() {
		for i, col := range o.Colors {
			i, col := i, col
			name := pickerName(o, i, col)
			item := ui.RadioBase(c, &selected, i).Size(side, side).
				Shrink(0).Radius(itemRadius).FocusRing(false).
				Cursor(cursor).Label(name).Tooltip(name)
			chosen := i == selected
			// Snapshot the swatch's states in the build: a retained copy of
			// a leaving element paints in later passes, where item has
			// expired.
			hovered, focused := item.Hovered() && !o.Disabled, item.FocusVisible()
			item.Draw(func(p *ui.Painter, r ui.Rect) {
				scale := float32(1)
				if chosen {
					scale = 0.77
				} else if hovered {
					scale = 1.1
				}
				sw := side * scale
				sr := swatchRadiusValue * scale
				box := ui.Rect{
					X: r.X + (r.W-sw)/2, Y: r.Y + (r.H-sw)/2,
					W: sw, H: sw,
				}
				paintSwatch(p, box, col, sr)
				paintSwatchBorder(p, box, sr)
				if chosen {
					ink := visibleColor(col, fieldBg)
					p.Stroke(r, ink, itemRadius, border)
					check := ui.RGB(255, 255, 255)
					if luminance(ink) > 0.5 {
						check = ui.RGB(0, 0, 0)
					}
					paintCheckmark(p, r, check)
				}
				if focused {
					p.FocusRing(r, [4]float32{itemRadius, itemRadius, itemRadius, itemRadius})
				}
			})
			item.OnChange(func() {
				if o.Value == nil {
					return
				}
				if *o.Value != col {
					*o.Value = col
				}
				if o.OnChange != nil {
					o.OnChange()
				}
			})
		}
	}).Label(label).Gap(theme.Space8).AlignItems(ui.Center).
		Disabled(o.Disabled)
	if o.Layout == SwatchStack {
		group = group.Column().Shrink(0)
	} else {
		group = group.Row().Wrap().FillWidth()
	}
	return style.Disabled(group, o.Disabled)
}

// pickerRadii returns the item and inner swatch radii of a picker.
func pickerRadii(size ColorSwatchSize, shape ColorSwatchShape) (item, swatch float32) {
	if shape == SwatchSquare {
		switch size {
		case SwatchXS:
			return 6, 6
		case SwatchSM:
			return 8, 8
		}
		return 12, 8
	}
	switch size {
	case SwatchXS:
		return 8, 8
	case SwatchSM:
		return 12, 12
	case SwatchLG, SwatchXL:
		return 24, 24
	}
	return 16, 16
}

// pickerName returns the accessible name of swatch i.
func pickerName(o ColorSwatchPickerOptions, i int, col ui.Color) string {
	if i < len(o.Names) {
		if name := strings.TrimSpace(o.Names[i]); name != "" {
			return name
		}
	}
	return strings.ToUpper(FormatHex(col))
}
