package theme

import (
	"slices"

	"github.com/egoist/mygo/ui"
)

// Sans is MyGo's native interface family list.
const Sans = "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"

// Mono is MyGo's native monospace family list.
const Mono = "ui-monospace, SF Mono, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace"

// The built-in spacing rhythm is 4 DIPs, named by DIP value.
const (
	SpaceUnit = 4
	Space2    = 2
	Space4    = 4
	Space6    = 6
	Space8    = 8
	Space10   = 10
	Space12   = 12
	Space14   = 14
	Space16   = 16
	Space20   = 20
	Space24   = 24
	Space32   = 32
	Space40   = 40
	Space48   = 48
	Space64   = 64
)

// Space returns n spacing units of the published theme. The constants above
// are the steps of the built-in 4 DIP unit and do not follow a theme.
func Space(n float32) float32 { return n * snapshot().theme.Spacing.Unit }

// Text weights map to Element.FontWeight.
const (
	WeightNormal   = 400
	WeightMedium   = 500
	WeightSemibold = 600
	WeightBold     = 700
)

// TrackingTight is Tailwind's tracking-tight, a letter spacing of -0.025em.
// coss ui sets it on every heading role and on prose headings.
const TrackingTight = -0.025

// TextStyle describes an unscaled type role. Tracking is letter
// spacing as a fraction of the font size: zero keeps the font's own spacing,
// TrackingTight tightens headings.
type TextStyle struct {
	Size       float32
	LineHeight float32
	Weight     int
	Tracking   float32
	Mono       bool
	Heading    bool
}

// The type roles read the published theme.

// XS is text-xs: 12/16 by default.
func XS() TextStyle { return snapshot().theme.Typography.XS }

// SM is text-sm: 14/20 by default.
func SM() TextStyle { return snapshot().theme.Typography.SM }

// Base is text-base: 16/24 by default.
func Base() TextStyle { return snapshot().theme.Typography.Base }

// LG is text-lg: 18/28 by default.
func LG() TextStyle { return snapshot().theme.Typography.LG }

// XL is text-xl: 20/28 by default.
func XL() TextStyle { return snapshot().theme.Typography.XL }

// H1 is the first heading role: 36/40 semibold and tracking-tight by default.
func H1() TextStyle { return snapshot().theme.Typography.H1 }

// H2 is the second heading role: 30/36 by default.
func H2() TextStyle { return snapshot().theme.Typography.H2 }

// H3 is the third heading role: 24/32 by default.
func H3() TextStyle { return snapshot().theme.Typography.H3 }

// H4 is the fourth heading role: 20/28 by default.
func H4() TextStyle { return snapshot().theme.Typography.H4 }

// H5 is the fifth heading role: 18/28 by default.
func H5() TextStyle { return snapshot().theme.Typography.H5 }

// H6 is the sixth heading role: 16/24 by default.
func H6() TextStyle { return snapshot().theme.Typography.H6 }

// Code is the inline code role.
func Code() TextStyle { return snapshot().theme.Typography.Code }

// Label is the form label role.
func Label() TextStyle { return snapshot().theme.Typography.Label }

// HelpText is the form description role.
func HelpText() TextStyle { return snapshot().theme.Typography.HelpText }

// SemanticType is a coss ui Typography type.
type SemanticType uint8

const (
	TypeBody SemanticType = iota
	TypeBodySM
	TypeBodyXS
	TypeH1
	TypeH2
	TypeH3
	TypeH4
	TypeH5
	TypeH6
	TypeCode
)

// Style resolves a semantic Typography type.
func (t SemanticType) Style() TextStyle {
	switch t {
	case TypeBodySM:
		s := SM()
		s.LineHeight = 24
		return s
	case TypeBodyXS:
		s := XS()
		s.LineHeight = 20
		return s
	case TypeH1:
		return H1()
	case TypeH2:
		return H2()
	case TypeH3:
		return H3()
	case TypeH4:
		return H4()
	case TypeH5:
		return H5()
	case TypeH6:
		return H6()
	case TypeCode:
		return Code()
	}
	s := Base()
	s.LineHeight = 28
	return s
}

// RadiusBase returns the radius base in effect.
func RadiusBase() float32 { return snapshot().radius }

// Radii is the radius scale derived from the configured base.
type Radii struct{ XS, SM, MD, LG, XL, XXL, XXXL, XXXXL, Field float32 }

// radii follows coss's offset scale. At the default 10 DIP base the small,
// medium, large and extra-large radii are 6, 8, 10 and 14 DIPs.
func radii(base float32) Radii {
	return Radii{
		XS: max(0, base-6), SM: max(0, base-4), MD: max(0, base-2), LG: base,
		XL: base + 4, XXL: base + 6, XXXL: base + 14, XXXXL: base + 22,
		Field: max(0, base-2),
	}
}

// RadiusScale returns the scale published with the current preferences, so
// components do not re-derive it on every build.
func RadiusScale() Radii { return snapshot().radii }

// Component metrics in desktop DIPs. coss ui's mobile button heights are 4 DIPs
// taller; this library targets desktop windows.
const (
	HeightButtonSM = 28
	HeightButtonMD = 32
	HeightButtonLG = 36

	HeightInput    = 32
	HeightTextarea = 38
	HeightChipSM   = 20
	HeightChipMD   = 24
	HeightChipLG   = 32
	HeightBadgeSM  = 16
	HeightBadgeMD  = 18
	HeightBadgeLG  = 22
	HeightAvatarSM = 24
	HeightAvatarMD = 32
	HeightAvatarLG = 40
	HeightKbd      = 20
	HeightRow      = 40
	Hairline       = 1

	IconXS = 12
	IconSM = 14
	IconMD = 16
	IconLG = 20
)

// ButtonSize is coss ui's button size scale.
type ButtonSize uint8

const (
	ButtonMedium ButtonSize = iota
	ButtonSmall
	ButtonLarge
	ButtonExtraSmall
	ButtonExtraLarge
)

// ButtonHeight resolves a button size.
func ButtonHeight(size ButtonSize) float32 {
	switch size {
	case ButtonSmall:
		return HeightButtonSM
	case ButtonLarge:
		return HeightButtonLG
	case ButtonExtraSmall:
		return 24
	case ButtonExtraLarge:
		return 40
	}
	return HeightButtonMD
}

// Shadow is one native drop-shadow layer.
type Shadow struct {
	X, Y, Blur, Spread float32
	Color              ui.Color
	Inset              bool
}

// SurfaceShadows returns an independent copy of the card elevation.
func SurfaceShadows(dark bool) []Shadow {
	return slices.Clone(snapshot().theme.Elevation.mode(dark).Surface)
}

// OverlayShadows returns an independent copy of the popup and dialog
// elevation.
func OverlayShadows(dark bool) []Shadow {
	return slices.Clone(snapshot().theme.Elevation.mode(dark).Overlay)
}

// TooltipShadows returns an independent copy of the lighter tooltip
// elevation.
func TooltipShadows(dark bool) []Shadow {
	return slices.Clone(snapshot().theme.Elevation.mode(dark).Tooltip)
}

// FieldShadows returns an independent copy of the fine edge and lower shadow
// of an outlined control.
func FieldShadows(dark bool) []Shadow {
	return slices.Clone(snapshot().theme.Elevation.mode(dark).Field)
}
