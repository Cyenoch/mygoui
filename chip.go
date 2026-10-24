package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusColor is HeroUI's status color vocabulary.
type StatusColor uint8

const (
	ColorDefault StatusColor = iota
	ColorAccent
	ColorSuccess
	ColorWarning
	ColorDanger
)

// ChipVariant selects a HeroUI chip surface treatment.
type ChipVariant uint8

const (
	ChipSecondary ChipVariant = iota
	ChipPrimary
	ChipSoft
	ChipTertiary
)

// ChipSize selects HeroUI's chip density.
type ChipSize uint8

const (
	ChipMedium ChipSize = iota
	ChipSmall
	ChipLarge
)

// ChipOptions configures a compact HeroUI chip.
type ChipOptions struct {
	Label   string
	Color   StatusColor
	Variant ChipVariant
	Size    ChipSize
	Icon    *ui.SVG
}

// Chip renders a 12 DIP medium capsule with HeroUI's color/variant matrix.
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
		ts.Weight = theme.WeightMedium
	}
	e := ui.Row(c).Shrink(0).Center().Gap(theme.Space2).AlignSelf(ui.Start).
		Padding(paddingY, paddingX, paddingY, paddingX).Radius(theme.RadiusScale().XXL).
		Background(bg).TextColor(fg).Label(o.Label)
	e.Children(func() {
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconXS), theme.Scale(c, theme.IconXS)).Shrink(0)
		}
		label := style.Text(ui.Text(c, o.Label), c, ts)
		label.PaddingX(theme.Space2).SingleLine()
	})
	return e
}

func chipColors(c *ui.Context, color StatusColor, variant ChipVariant) (bg, fg ui.Color) {
	p := StatusColors(c, color)
	switch variant {
	case ChipPrimary:
		return p.Solid, p.SolidForeground
	case ChipSoft:
		return p.Soft, p.SoftForeground
	case ChipTertiary:
		return ui.Transparent, p.Foreground
	}
	return theme.Colors(c).Default, p.Foreground
}

// BadgeVariant selects a HeroUI badge surface treatment.
type BadgeVariant = ChipVariant

// BadgeSize is HeroUI's badge scale.
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

// BadgeOptions configures a HeroUI badge.
type BadgeOptions struct {
	Label   string
	Color   StatusColor
	Variant BadgeVariant
	Size    BadgeSize
}

// Badge renders a minimum 28 DIP capsule with a 1 DIP background-colored ring.
func Badge(c *ui.Context, o BadgeOptions) ui.Element {
	minSize, radius := float32(theme.HeightBadgeMD), theme.RadiusScale().XXXL
	ts := theme.XS()
	ts.Weight = theme.WeightMedium
	switch o.Size {
	case BadgeSmall:
		minSize, radius = theme.HeightBadgeSM, theme.RadiusScale().XL
		ts = theme.TextStyle{Size: 10, LineHeight: 14, Weight: theme.WeightMedium}
	case BadgeLarge:
		minSize, radius = theme.HeightBadgeLG, theme.RadiusScale().XXL
		ts = theme.SM()
		ts.Weight = theme.WeightMedium
	}
	bg, fg := chipColors(c, o.Color, o.Variant)
	e := ui.Row(c).Shrink(0).Center().AlignSelf(ui.Start).
		MinHeight(theme.Scale(c, minSize)).MinWidth(theme.Scale(c, minSize)).
		Radius(radius).Background(bg).TextColor(fg).
		Border(theme.Hairline, theme.Colors(c).Background).Label(o.Label)
	e.Children(func() {
		style.Text(ui.Text(c, o.Label), c, ts).PaddingX(theme.Space2).SingleLine()
	})
	return e
}

// BadgeAnchorOptions places a badge over caller content.
type BadgeAnchorOptions struct {
	Placement BadgePlacement
	Badge     BadgeOptions
}

// BadgeAnchor overlays a badge at one corner without changing the content's
// layout. The badge lives in the native overlay layer, so it can extend beyond
// the content box.
func BadgeAnchor(c *ui.Context, o BadgeAnchorOptions, body func()) ui.Element {
	content := ui.Box(c).Shrink(0).AlignSelf(ui.Start)
	content.Children(body)
	ui.Overlay(c, func() {
		at, self := ui.AnchorTopEnd, ui.AnchorTopEnd
		switch o.Placement {
		case BadgeTopLeft:
			at, self = ui.AnchorTopStart, ui.AnchorTopStart
		case BadgeBottomRight:
			at, self = ui.AnchorBottomEnd, ui.AnchorBottomEnd
		case BadgeBottomLeft:
			at, self = ui.AnchorBottomStart, ui.AnchorBottomStart
		}
		size := float32(theme.HeightBadgeMD)
		switch o.Badge.Size {
		case BadgeSmall:
			size = theme.HeightBadgeSM
		case BadgeLarge:
			size = theme.HeightBadgeLG
		}
		over := size / 4
		badge := Badge(c, o.Badge)
		switch o.Placement {
		case BadgeTopLeft:
			badge.Margin(-over, 0, 0, -over)
		case BadgeBottomRight:
			badge.Margin(0, -over, -over, 0)
		case BadgeBottomLeft:
			badge.Margin(0, 0, -over, -over)
		default:
			badge.Margin(-over, -over, 0, 0)
		}
		badge.AttachTo(content, at, self)
	})
	return content
}
