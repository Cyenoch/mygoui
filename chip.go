package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusColor selects a semantic color for native status components.
type StatusColor uint8

const (
	ColorDefault StatusColor = iota
	ColorAccent
	ColorSuccess
	ColorWarning
	ColorDanger
	// ColorInfo uses the theme's informational color.
	ColorInfo
)

// ChipVariant selects a native chip surface treatment.
type ChipVariant uint8

const (
	ChipSecondary ChipVariant = iota
	ChipPrimary
	ChipSoft
	ChipTertiary
)

// ChipSize selects the chip's density.
type ChipSize uint8

const (
	ChipMedium ChipSize = iota
	ChipSmall
	ChipLarge
)

// ChipOptions configures a compact native chip.
type ChipOptions struct {
	Label   string
	Color   StatusColor
	Variant ChipVariant
	Size    ChipSize
	Icon    *ui.SVG
}

// Chip renders a compact status label. Its icon follows the label size.
func Chip(c *ui.Context, o ChipOptions) ui.Element {
	bg, fg := chipColors(c, o.Color, o.Variant)
	paddingX, paddingY := float32(theme.Space8), float32(theme.Space2)
	ts := theme.XS()
	ts.LineHeight, ts.Weight = 20, theme.WeightMedium
	switch o.Size {
	case ChipSmall:
		paddingX, paddingY = theme.Space4, 0
	case ChipLarge:
		paddingX, paddingY = theme.Space12, theme.Space4
		ts = theme.SM()
		ts.Weight = theme.WeightMedium
	}
	iconSize := theme.Scale(c, ts.Size)
	e := ui.Row(c).Shrink(0).Center().Gap(theme.Space2).AlignSelf(ui.Start).
		Padding(paddingY, paddingX, paddingY, paddingX).Radius(theme.RadiusScale().SM).
		Background(bg).TextColor(fg).Label(o.Label)
	e.Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0)
		}
		label := style.Text(ui.Text(c, o.Label), c, ts)
		label.PaddingX(theme.Space2).SingleLine()
	})
	return e
}

// chipColors resolves the native chip color matrix. A default chip keeps the
// default foreground for its secondary and tertiary variants.
func chipColors(c *ui.Context, color StatusColor, variant ChipVariant) (bg, fg ui.Color) {
	t := theme.Colors(c)
	p := StatusColors(c, color)
	switch variant {
	case ChipPrimary:
		return p.Solid, p.SolidForeground
	case ChipSoft:
		return p.Soft, p.SoftForeground
	case ChipTertiary:
		if color == ColorDefault {
			return ui.Transparent, t.DefaultForeground
		}
		return ui.Transparent, p.Foreground
	}
	if color == ColorDefault {
		return t.Default, t.DefaultForeground
	}
	return t.Default, p.Foreground
}

// BadgeVariant selects a coss ui badge surface treatment.
type BadgeVariant uint8

const (
	BadgePrimary BadgeVariant = iota
	BadgeSecondary
	BadgeOutline
	BadgeDestructive
	BadgeError
	BadgeInfo
	BadgeSuccess
	BadgeWarning
)

// BadgeSize is coss ui's badge scale.
type BadgeSize uint8

const (
	BadgeMedium BadgeSize = iota
	BadgeSmall
	BadgeLarge
)

// BadgePlacement is the corner used by BadgeAnchor.
type BadgePlacement uint8

const (
	BadgeTopRight BadgePlacement = iota
	BadgeTopLeft
	BadgeBottomRight
	BadgeBottomLeft
)

// BadgeOptions configures a coss ui badge.
type BadgeOptions struct {
	Label   string
	Variant BadgeVariant
	Size    BadgeSize
	Icon    *ui.SVG
}

// Badge renders a compact COSS status label with a size-dependent minimum.
func Badge(c *ui.Context, o BadgeOptions) ui.Element {
	minSize, radius := float32(theme.HeightBadgeMD), theme.RadiusScale().SM
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	switch o.Size {
	case BadgeSmall:
		minSize, radius = theme.HeightBadgeSM, 4
		ts = theme.TextStyle{Size: 10, LineHeight: 14, Weight: theme.WeightMedium}
	case BadgeLarge:
		minSize, radius = theme.HeightBadgeLG, theme.RadiusScale().SM
		ts = theme.SM()
		ts.Weight = theme.WeightMedium
	}
	t := theme.Colors(c)
	bg, fg := t.Primary, t.PrimaryForeground
	border := ui.Transparent
	switch o.Variant {
	case BadgeSecondary:
		bg, fg = t.Secondary, t.SecondaryForeground
	case BadgeOutline:
		bg, fg, border = t.Background, t.Foreground, t.Input
	case BadgeError:
		bg, fg = t.DangerSoft, t.DestructiveForeground
	case BadgeInfo:
		bg, fg = t.InfoSoft, t.InfoForeground
	case BadgeSuccess:
		bg, fg = t.SuccessSoft, t.SuccessForeground
	case BadgeWarning:
		bg, fg = t.WarningSoft, t.WarningForeground
	case BadgeDestructive:
		bg, fg = t.Destructive, ui.RGB(255, 255, 255)
	}
	padding := float32(3)
	if o.Size == BadgeLarge {
		padding = 5
	}
	e := ui.Row(c).Shrink(0).Center().AlignSelf(ui.Start).
		MinHeight(theme.Scale(c, minSize)).MinWidth(theme.Scale(c, minSize)).
		Radius(radius).Background(bg).TextColor(fg).
		Border(theme.Hairline, border).Gap(theme.Space4).PaddingX(padding).Label(o.Label)
	e.Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, 12), theme.Scale(c, 12)).
				Shrink(0).Opacity(0.8).PassThrough()
		}
		style.Text(ui.Text(c, o.Label), c, ts).SingleLine().PassThrough()
	})
	return e
}

// BadgeAnchorOptions places a badge over caller content.
type BadgeAnchorOptions struct {
	Placement BadgePlacement
	Badge     BadgeOptions
}

// BadgeAnchor places a badge at one corner without changing the content's
// layout. The badge can extend beyond the box and follows ancestor clipping.
func BadgeAnchor(c *ui.Context, o BadgeAnchorOptions, body func()) ui.Element {
	content := ui.Box(c).Shrink(0).AlignSelf(ui.Start)
	content.Children(func() {
		if body != nil {
			body()
		}
		size := float32(theme.HeightBadgeMD)
		switch o.Badge.Size {
		case BadgeSmall:
			size = theme.HeightBadgeSM
		case BadgeLarge:
			size = theme.HeightBadgeLG
		}
		over := theme.Scale(c, size) / 4
		badge := Badge(c, o.Badge).Absolute().PassThrough()
		switch o.Placement {
		case BadgeTopLeft:
			badge.Top(-over).Left(-over)
		case BadgeBottomRight:
			badge.Bottom(-over).Right(-over)
		case BadgeBottomLeft:
			badge.Bottom(-over).Left(-over)
		default:
			badge.Top(-over).Right(-over)
		}
	})
	return content
}
