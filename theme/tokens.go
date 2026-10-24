package theme

import "github.com/egoist/mygo/ui"

// Sans is MyGo's native interface family list.
const Sans = "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"

// Mono is MyGo's native monospace family list.
const Mono = "ui-monospace, SF Mono, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace"

// HeroUI's spacing rhythm is 0.25rem at a 16 DIP root.
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

// Space returns n HeroUI spacing units.
func Space(n float32) float32 { return n * SpaceUnit }

// Text weights map to Element.FontWeight.
const (
	WeightNormal   = 400
	WeightMedium   = 500
	WeightSemibold = 600
	WeightBold     = 700
)

// TrackingTight is Tailwind's tracking-tight, a letter spacing of -0.025em.
// HeroUI sets it on every heading role and on prose headings.
const TrackingTight = -0.025

// TextStyle describes an unscaled HeroUI type role. Tracking is letter
// spacing as a fraction of the font size: zero keeps the font's own spacing,
// TrackingTight tightens headings.
type TextStyle struct {
	Size       float32
	LineHeight float32
	Weight     int
	Tracking   float32
	Mono       bool
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

// H1 is HeroUI's text-4xl heading role, 36/40 semibold and tracking-tight.
func H1() TextStyle {
	return TextStyle{Size: 36, LineHeight: 40, Weight: WeightSemibold, Tracking: TrackingTight}
}

// H2 is HeroUI's text-3xl heading role, 30/36 semibold and tracking-tight.
func H2() TextStyle {
	return TextStyle{Size: 30, LineHeight: 36, Weight: WeightSemibold, Tracking: TrackingTight}
}

// H3 is HeroUI's text-2xl heading role, 24/32 semibold and tracking-tight.
func H3() TextStyle {
	return TextStyle{Size: 24, LineHeight: 32, Weight: WeightSemibold, Tracking: TrackingTight}
}

// H4 is HeroUI's text-xl heading role, 20/28 semibold and tracking-tight.
func H4() TextStyle {
	return TextStyle{Size: 20, LineHeight: 28, Weight: WeightSemibold, Tracking: TrackingTight}
}

// H5 is HeroUI's text-lg heading role, 18/28 semibold and tracking-tight.
func H5() TextStyle {
	return TextStyle{Size: 18, LineHeight: 28, Weight: WeightSemibold, Tracking: TrackingTight}
}

// H6 is HeroUI's text-base heading role, 16/24 semibold and tracking-tight.
func H6() TextStyle {
	return TextStyle{Size: 16, LineHeight: 24, Weight: WeightSemibold, Tracking: TrackingTight}
}

// Code is the inline code role.
func Code() TextStyle { return TextStyle{Size: 14, LineHeight: 20, Weight: WeightNormal, Mono: true} }

// Label is the form label role.
func Label() TextStyle { return TextStyle{Size: 14, LineHeight: 20, Weight: WeightMedium} }

// HelpText is the form description role.
func HelpText() TextStyle { return TextStyle{Size: 12, LineHeight: 16, Weight: WeightNormal} }

// SemanticType is a HeroUI Typography type.
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

// RadiusBase returns the configured HeroUI radius base.
func RadiusBase() float32 { return PreferencesSnapshot().Radius }

// Radii is HeroUI's radius scale derived from the configured base.
type Radii struct{ XS, SM, MD, LG, XL, XXL, XXXL, XXXXL, Field float32 }

// radii derives the scale from a base. Upstream values at the 8 DIP base:
// 2, 4, 6, 8, 12, 16, 24, 32, field 12.
func radii(base float32) Radii {
	return Radii{
		XS: base * 0.25, SM: base * 0.5, MD: base * 0.75, LG: base,
		XL: base * 1.5, XXL: base * 2, XXXL: base * 3, XXXXL: base * 4,
		Field: base * 1.5,
	}
}

// RadiusScale returns the scale published with the current preferences, so
// components do not re-derive it on every build.
func RadiusScale() Radii { return snapshot().radii }

// Component metrics in desktop DIPs. HeroUI's mobile button heights are 4 DIPs
// taller; this library targets desktop windows.
const (
	HeightButtonSM = 32
	HeightButtonMD = 36
	HeightButtonLG = 40

	HeightInput    = 36
	HeightTextarea = 38
	HeightChipSM   = 20
	HeightChipMD   = 24
	HeightChipLG   = 32
	HeightBadgeSM  = 16
	HeightBadgeMD  = 28
	HeightBadgeLG  = 32
	HeightAvatarSM = 32
	HeightAvatarMD = 40
	HeightAvatarLG = 48
	HeightKbd      = 24
	HeightRow      = 40
	Hairline       = 1

	IconXS = 12
	IconSM = 14
	IconMD = 16
	IconLG = 20
)

// ButtonSize is HeroUI's button size scale.
type ButtonSize uint8

const (
	ButtonMedium ButtonSize = iota
	ButtonSmall
	ButtonLarge
)

// ButtonHeight resolves a button size.
func ButtonHeight(size ButtonSize) float32 {
	switch size {
	case ButtonSmall:
		return HeightButtonSM
	case ButtonLarge:
		return HeightButtonLG
	}
	return HeightButtonMD
}

// Shadow is one native drop-shadow layer.
type Shadow struct {
	X, Y, Blur, Spread float32
	Color              ui.Color
	Inset              bool
}

// HeroUI's three elevation tokens, layer for layer. The slices are shared by
// every caller, so a build pass allocates nothing to paint a surface; callers
// read them, and an append to one allocates its own array. Dark surfaces
// separate by contrast, so upstream gives them no drop shadow.
var (
	lightSurfaceShadows = []Shadow{
		{X: 0, Y: 2, Blur: 4, Color: ui.RGBA(0, 0, 0, 0.04)},
		{X: 0, Y: 1, Blur: 2, Color: ui.RGBA(0, 0, 0, 0.06)},
		{X: 0, Y: 0, Blur: 1, Color: ui.RGBA(0, 0, 0, 0.06)},
	}
	darkOverlayShadows = []Shadow{
		{X: 0, Y: 0, Blur: 1, Color: ui.RGBA(255, 255, 255, 0.3), Inset: true},
	}
	lightOverlayShadows = []Shadow{
		{X: 0, Y: 2, Blur: 8, Color: ui.RGBA(0, 0, 0, 0.06)},
		{X: 0, Y: -6, Blur: 12, Color: ui.RGBA(0, 0, 0, 0.03)},
		{X: 0, Y: 14, Blur: 28, Color: ui.RGBA(0, 0, 0, 0.08)},
	}
)

// shared returns s with its capacity trimmed to its length, so appending to
// the result allocates instead of writing into the package's table.
func shared(s []Shadow) []Shadow { return s[:len(s):len(s)] }

// SurfaceShadows returns the light card shadow. Dark surfaces separate by
// contrast; upstream disables the shadow. The result is shared and read-only.
func SurfaceShadows(dark bool) []Shadow {
	if dark {
		return nil
	}
	return shared(lightSurfaceShadows)
}

// OverlayShadows returns the floating-panel shadow. Dark overlays use
// HeroUI's inset white hairline instead. The result is shared and read-only.
func OverlayShadows(dark bool) []Shadow {
	if dark {
		return shared(darkOverlayShadows)
	}
	return shared(lightOverlayShadows)
}

// FieldShadows returns the form-field elevation, which is the surface shadow
// in light and a transparent inset placeholder upstream. The result is shared
// and read-only.
func FieldShadows(dark bool) []Shadow { return SurfaceShadows(dark) }
