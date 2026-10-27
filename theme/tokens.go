package theme

import "github.com/egoist/mygo/ui"

// Sans is MyGo's native interface family list.
const Sans = "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"

// Mono is MyGo's native monospace family list.
const Mono = "ui-monospace, SF Mono, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace"

// coss ui's spacing rhythm is 0.25rem at a 16 DIP root.
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

// Space returns n coss ui spacing units.
func Space(n float32) float32 { return n * SpaceUnit }

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

// TextStyle describes an unscaled coss ui type role. Tracking is letter
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

// XS is text-xs: 12/16.
func XS() TextStyle { return TextStyle{Size: 12, LineHeight: 16, Weight: WeightNormal} }

// SM is text-sm: 14/20.
func SM() TextStyle { return TextStyle{Size: 14, LineHeight: 20, Weight: WeightNormal} }

// Base is text-base: 16/24.
func Base() TextStyle { return TextStyle{Size: 16, LineHeight: 24, Weight: WeightNormal} }

// LG is text-lg: 18/28.
func LG() TextStyle { return TextStyle{Size: 18, LineHeight: 28, Weight: WeightNormal} }

// XL is text-xl: 20/28.
func XL() TextStyle { return TextStyle{Size: 20, LineHeight: 28, Weight: WeightNormal} }

// H1 is coss ui's text-4xl heading role, 36/40 semibold and tracking-tight.
func H1() TextStyle {
	return TextStyle{Size: 36, LineHeight: 40, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
}

// H2 is coss ui's text-3xl heading role, 30/36 semibold and tracking-tight.
func H2() TextStyle {
	return TextStyle{Size: 30, LineHeight: 36, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
}

// H3 is coss ui's text-2xl heading role, 24/32 semibold and tracking-tight.
func H3() TextStyle {
	return TextStyle{Size: 24, LineHeight: 32, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
}

// H4 is coss ui's text-xl heading role, 20/28 semibold and tracking-tight.
func H4() TextStyle {
	return TextStyle{Size: 20, LineHeight: 28, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
}

// H5 is coss ui's text-lg heading role, 18/28 semibold and tracking-tight.
func H5() TextStyle {
	return TextStyle{Size: 18, LineHeight: 28, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
}

// H6 is coss ui's text-base heading role, 16/24 semibold and tracking-tight.
func H6() TextStyle {
	return TextStyle{Size: 16, LineHeight: 24, Weight: WeightSemibold, Tracking: TrackingTight, Heading: true}
}

// Code is the inline code role.
func Code() TextStyle { return TextStyle{Size: 14, LineHeight: 20, Weight: WeightNormal, Mono: true} }

// Label is the form label role.
func Label() TextStyle { return TextStyle{Size: 14, LineHeight: 20, Weight: WeightMedium} }

// HelpText is the form description role.
func HelpText() TextStyle { return TextStyle{Size: 12, LineHeight: 16, Weight: WeightNormal} }

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

// RadiusBase returns the configured coss ui radius base.
func RadiusBase() float32 { return snapshot().prefs.Radius }

// Radii is coss ui's radius scale derived from the configured base.
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

// Coss layers a low-alpha bottom shadow with a crisp inset edge.
var (
	lightSurfaceShadows = []Shadow{{Y: 1, Blur: 2, Color: ui.RGBA(0, 0, 0, 0.05)}, {Y: 1, Color: ui.RGBA(0, 0, 0, 0.04), Inset: true}}
	darkSurfaceShadows  = []Shadow{{Y: -1, Color: ui.RGBA(255, 255, 255, 0.06), Inset: true}}
	lightOverlayShadows = []Shadow{{Y: 4, Blur: 6, Spread: -1, Color: ui.RGBA(0, 0, 0, 0.1)}, {Y: 2, Blur: 4, Spread: -2, Color: ui.RGBA(0, 0, 0, 0.1)}}
	lightFieldShadows   = []Shadow{{Y: 1, Blur: 2, Color: ui.RGBA(0, 0, 0, 0.05)}, {Y: 1, Color: ui.RGBA(0, 0, 0, 0.06), Inset: true}}
)

// SurfaceShadows returns an independent copy of the coss card elevation.
func SurfaceShadows(dark bool) []Shadow {
	if dark {
		return append([]Shadow{}, darkSurfaceShadows...)
	}
	return append([]Shadow{}, lightSurfaceShadows...)
}

// OverlayShadows returns the coss popup elevation in either appearance.
func OverlayShadows(dark bool) []Shadow { return append([]Shadow{}, lightOverlayShadows...) }

// FieldShadows returns the fine edge and lower shadow of an outlined input.
func FieldShadows(dark bool) []Shadow {
	if dark {
		return append([]Shadow{}, darkSurfaceShadows...)
	}
	return append([]Shadow{}, lightFieldShadows...)
}
