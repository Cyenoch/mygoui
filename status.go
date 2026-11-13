package mygoui

import (
	"math"
	"sync"

	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusPalette is the resolved color set for one status. Components
// and applications use it to paint custom status surfaces consistently.
type StatusPalette struct {
	Solid, SolidForeground ui.Color
	Hover                  ui.Color
	Soft, SoftForeground   ui.Color
	SoftHover              ui.Color
	Foreground             ui.Color
}

// StatusColors resolves a status against the live theme.
func StatusColors(c *ui.Context, color StatusColor) StatusPalette {
	t := theme.Colors(c)
	return statusColors(&t, color)
}

// statusColors resolves a status against t, which callers holding the
// palette pass instead of copying it again.
func statusColors(t *theme.Palette, color StatusColor) StatusPalette {
	// Every solid hovers the same way: the fill at the theme's hover opacity.
	dim := theme.Effects().HoverOpacity
	hover := t.Primary.Alpha(dim)
	switch color {
	case ColorAccent:
		return StatusPalette{
			Solid: t.Primary, SolidForeground: t.PrimaryForeground, Hover: hover,
			Soft: t.AccentSoft, SoftForeground: t.AccentSoftForeground,
			SoftHover: t.AccentSoftHover, Foreground: t.AccentSoftForeground,
		}
	case ColorSuccess:
		return StatusPalette{
			Solid: t.Success, SolidForeground: solidForeground(t.Success, t.Background), Hover: t.Success.Alpha(dim),
			Soft: t.SuccessSoft, SoftForeground: t.SuccessSoftForeground,
			SoftHover: t.SuccessSoftHover, Foreground: t.SuccessSoftForeground,
		}
	case ColorWarning:
		return StatusPalette{
			Solid: t.Warning, SolidForeground: solidForeground(t.Warning, t.Background), Hover: t.Warning.Alpha(dim),
			Soft: t.WarningSoft, SoftForeground: t.WarningSoftForeground,
			SoftHover: t.WarningSoftHover, Foreground: t.WarningSoftForeground,
		}
	case ColorDanger:
		return StatusPalette{
			Solid: t.Danger, SolidForeground: t.DangerForeground, Hover: t.Danger.Alpha(dim),
			Soft: t.DangerSoft, SoftForeground: t.DangerSoftForeground,
			SoftHover: t.DangerSoftHover, Foreground: t.DangerSoftForeground,
		}
	case ColorInfo:
		return StatusPalette{
			Solid: t.Info, SolidForeground: solidForeground(t.Info, t.Background), Hover: t.Info.Alpha(dim),
			Soft: t.InfoSoft, SoftForeground: t.InfoSoftForeground,
			SoftHover: t.InfoSoftHover, Foreground: t.InfoSoftForeground,
		}
	}
	return StatusPalette{
		Solid: t.Primary, SolidForeground: t.PrimaryForeground, Hover: hover,
		Soft: t.DefaultSoft, SoftForeground: t.DefaultSoftForeground,
		SoftHover: t.DefaultSoftHover, Foreground: t.Foreground,
	}
}

// solidForeground picks the higher-contrast black or white ink for a status
// fill; the status foreground tokens are intended for pale surfaces.
// Its colors come from the published palettes, which change only with the
// preferences, so each pair is resolved once: compositing Oklch colors is
// costly on every build.
func solidForeground(fill, background ui.Color) ui.Color {
	key := [2]ui.Color{fill, background}
	solidInks.Lock()
	defer solidInks.Unlock()
	if ink, ok := solidInks.m[key]; ok {
		return ink
	}
	fill = fill.Over(background)
	channel := func(v uint8) float64 {
		s := float64(v) / 255
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	luminance := 0.2126*channel(fill.R) + 0.7152*channel(fill.G) + 0.0722*channel(fill.B)
	ink := ui.RGB(255, 255, 255)
	if luminance > 0.179 {
		ink = ui.RGB(0, 0, 0)
	}
	// A host that edits its accent live publishes a palette per change.
	if solidInks.m == nil || len(solidInks.m) >= 64 {
		solidInks.m = map[[2]ui.Color]ui.Color{}
	}
	solidInks.m[key] = ink
	return ink
}

var solidInks struct {
	sync.Mutex
	m map[[2]ui.Color]ui.Color
}

// statusMark names the mark whose colour shows status, as an image read by
// the status's word; a status without a meaning of its own leaves the mark
// out.
func statusMark(e ui.Element, color StatusColor) ui.Element {
	m := locale.M()
	name := ""
	switch color {
	case ColorSuccess:
		name = m.Success
	case ColorWarning:
		name = m.LevelWarning
	case ColorDanger:
		name = m.LevelError
	case ColorInfo:
		name = m.LevelInfo
	}
	if name == "" {
		return e.Role(ui.RoleNone)
	}
	return e.Role(ui.RoleImage).Label(name)
}

// StatusIcon is the coss icon for a status: info, check, triangle alert, or
// circle alert for danger.
func StatusIcon(color StatusColor) *ui.SVG {
	switch color {
	case ColorSuccess:
		return lucide.CircleCheck()
	case ColorWarning:
		return lucide.TriangleAlert()
	case ColorDanger:
		return lucide.CircleAlert()
	}
	return lucide.Info()
}
