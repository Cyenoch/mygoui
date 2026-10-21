package theme

import (
	"runtime"

	"github.com/egoist/mygo/ui"
)

const (
	// Sans is MyGo's unquoted native interface family list.
	Sans = "-apple-system, BlinkMacSystemFont, Segoe UI, system-ui, sans-serif"
	// Mono is MyGo's unquoted native monospace family list.
	Mono = "ui-monospace, SF Mono, SFMono-Regular, Menlo, Consolas, Liberation Mono, monospace"
	// SpaceUnit is the shared four-DIP spacing rhythm.
	SpaceUnit = 4
	// Space4 is one spacing unit.
	Space4 = SpaceUnit
	// Space6 is one and a half spacing units.
	Space6 = 6
	// Space8 is two spacing units.
	Space8 = 8
	// Space10 is two and a half spacing units.
	Space10 = 10
	// Space12 is three spacing units.
	Space12 = 12
	// Space14 is three and a half spacing units.
	Space14 = 14
	// Space16 is four spacing units.
	Space16 = 16
	// Space20 is five spacing units.
	Space20 = 20
	// Space24 is six spacing units.
	Space24 = 24
	// Space32 is eight spacing units.
	Space32 = 32

	// RadiusControl is the shared MD radius.
	RadiusControl = RadiusCard * 8 / 10
	// RadiusCard is the shared LG radius and radius-scale base.
	RadiusCard = 10
	// RadiusPanel is the shared XL radius at 1.4 times the base.
	RadiusPanel = RadiusCard * 14 / 10
	// RadiusChip is the shared SM radius for compact chips.
	RadiusChip = RadiusCard * 6 / 10
	// RadiusBubble is the 18-DIP conversation surface radius.
	RadiusBubble = RadiusCard * 18 / 10

	// HeightButton is the compact control height.
	HeightButton = 28
	// HeightInput is the default input height.
	HeightInput = 32
	// HeightChip is the compact chip height.
	HeightChip = 24
	// HeightHeader is the default reusable panel header height.
	HeightHeader = 52
	// HeightSidebarRow is the native dense navigation row target.
	HeightSidebarRow = 34
	// HeightSectionHeader is the compact navigation heading height.
	HeightSectionHeader = 32
	// HeightRowCompact is a compact list row target.
	HeightRowCompact = 36
	// HeightWorkRow is a dense activity row minimum height.
	HeightWorkRow = 24
	// SidebarRowRadius is the navigation row's control radius.
	SidebarRowRadius = RadiusControl
	// SidebarRowPadding is the navigation row's outer vertical padding.
	SidebarRowPadding = 2
	// SidebarRowInsetY is the inner navigation row's vertical inset.
	SidebarRowInsetY = Space8
	// SidebarRowInsetX is the navigation row's horizontal inset.
	SidebarRowInsetX = Space10

	// HairlineWidth is the 1 DIP separator width.
	HairlineWidth = 1
	// IconMicro is the small inline and section-heading glyph size.
	IconMicro = 12
	// IconSmall is the compact-control glyph size.
	IconSmall = 14
	// IconMedium is the standard-control glyph size.
	IconMedium = 16
	// IconLarge is the prominent-control glyph size.
	IconLarge = 20

	// SidebarWidth is the default navigation pane width.
	SidebarWidth = 256
	// SidebarMinWidth is the minimum navigation pane width.
	SidebarMinWidth = 208
	// FontTiny is the minimum native badge text size.
	FontTiny = 10
	// FontSmall is the native metadata text size.
	FontSmall = 11
	// FontBody is the native macOS body size; Body selects 14 elsewhere.
	FontBody = 13
	// LineHeightMeta is the metadata line box height in DIPs.
	LineHeightMeta = 16
	// LineHeightBadge is the badge line box height in DIPs.
	LineHeightBadge = 14
)

// Space returns n units of the four-DIP spacing rhythm, including half units.
func Space(n float32) float32 { return n * SpaceUnit }

// TextStyle holds native text metrics; it is deliberately separate from ui.Theme.
// Apply LineHeight with Element.LineHeight, or LineHeightPixels with
// Element.FixedLineHeight, not both. Scale Size and LineHeightPixels together
// for accessibility; the multiplier does not change.
type TextStyle struct {
	// Size is the font size in unscaled DIPs.
	Size float32
	// LineHeight is a dimensionless multiplier of Size.
	LineHeight float32
	// LineHeightPixels is the corresponding unscaled line height in DIPs.
	LineHeightPixels float32
	// Weight is the numeric font weight for Element.FontWeight.
	// MyGo's FontWeight accepts 100–900; CoreText's nsWeight
	// interpolates those CSS weights rather than collapsing to regular/bold.
	Weight int
	// Font is the comma-joined family list for Element.Font.
	Font string
	// Mono requests the system monospace family rather than the system UI font.
	Mono bool
}

// Meta uses 11-DIP text with a 16-DIP line box.
func Meta() TextStyle {
	return TextStyle{Size: FontSmall, LineHeight: float32(LineHeightMeta) / FontSmall,
		LineHeightPixels: LineHeightMeta, Weight: 400, Font: Sans}
}

// Small uses 12-DIP text with a 16-DIP line box.
func Small() TextStyle {
	return TextStyle{Size: 12, LineHeight: 16.0 / 12, LineHeightPixels: 16, Weight: 400, Font: Sans}
}

// Badge uses 10-DIP text with a 14-DIP line box.
func Badge() TextStyle {
	return TextStyle{Size: FontTiny, LineHeight: float32(LineHeightBadge) / FontTiny,
		LineHeightPixels: LineHeightBadge, Weight: 400, Font: Sans}
}

// Body keeps MyGo's platform body size (13 on macOS, 14 elsewhere).
// Its relaxed line-height multiplier is 1.625.
func Body() TextStyle {
	size := float32(FontBody)
	if runtime.GOOS != "darwin" {
		size = 14
	}
	return TextStyle{Size: size, LineHeight: 1.625, LineHeightPixels: size * 1.625, Weight: 400, Font: Sans}
}

// Prose uses 14-DIP text with a relaxed 1.625 line-height multiplier.
func Prose() TextStyle {
	return TextStyle{Size: 14, LineHeight: 1.625, LineHeightPixels: 22.75, Weight: 400, Font: Sans}
}

// Title uses 14-DIP text with a 1.3 line-height multiplier and semibold weight.
func Title() TextStyle {
	return TextStyle{Size: 14, LineHeight: 1.3, LineHeightPixels: 18.2, Weight: 600, Font: Sans}
}

// Code uses the native body size with a snug 1.375 line-height multiplier.
func Code() TextStyle {
	s := Body()
	s.LineHeight = 1.375
	s.LineHeightPixels = s.Size * s.LineHeight
	s.Mono = true
	s.Font = Mono
	return s
}

// RadiusTokens follows shadcn's multiplicative radius scale at a 10-DIP base.
// App-specific chrome geometry remains owned by its view, independent of this.
type RadiusTokens struct{ SM, MD, LG, XL, XXL, XXXL, XXXXL float32 }

// Radii returns value metrics; no CSS length calculation occurs in a view.
func Radii() RadiusTokens {
	return RadiusTokens{SM: RadiusChip, MD: RadiusControl, LG: RadiusCard,
		XL: RadiusPanel, XXL: RadiusBubble, XXXL: RadiusCard * 22 / 10, XXXXL: RadiusCard * 26 / 10}
}

// MetricTokens measures text-bearing library controls in native DIPs.
// Fixed application tab/titlebar geometry does not use these metrics.
type MetricTokens struct{ ControlHeight, CompactHeight, InputHeight, IconSize, LineHeight float32 }

// Metrics grows control minimums to fit scaled text and ordinary fixed padding;
// spacing and corner radii themselves do not grow with the interface font.
func Metrics(c *ui.Context) MetricTokens {
	line := Scale(c, Small().LineHeightPixels)
	return MetricTokens{ControlHeight: max(HeightInput, line+Space16),
		CompactHeight: max(HeightButton, line+Space12), InputHeight: max(HeightInput, line+Space16),
		IconSize: Scale(c, IconMedium), LineHeight: line}
}
