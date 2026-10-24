package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Orientation is a HeroUI layout direction.
type Orientation uint8

const (
	Horizontal Orientation = iota
	Vertical
)

// SeparatorColor is HeroUI's separator emphasis.
type SeparatorColor uint8

const (
	SeparatorDefault SeparatorColor = iota
	SeparatorSecondary
	SeparatorTertiary
)

// SeparatorOptions configures a semantic divider.
type SeparatorOptions struct {
	Orientation Orientation
	Color       SeparatorColor
}

// Separator draws a 1 DIP rule. Vertical rules stretch to the row height.
func Separator(c *ui.Context, o SeparatorOptions) ui.Element {
	t := theme.Colors(c)
	color := t.Separator
	switch o.Color {
	case SeparatorSecondary:
		color = t.SeparatorSecondary
	case SeparatorTertiary:
		color = t.SeparatorTertiary
	}
	e := ui.Box(c).Shrink(0).Background(color).Role(ui.RoleNone)
	if o.Orientation == Vertical {
		e.Width(theme.Hairline).MinHeight(theme.Space8).AlignSelf(ui.Stretch)
	} else {
		e.Height(theme.Hairline).FillWidth()
	}
	return e
}
