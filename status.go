package mygoui

import (
	"math"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusPalette is the resolved color set for one coss ui status. Components
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
	switch color {
	case ColorAccent:
		return StatusPalette{
			Solid: t.Primary, SolidForeground: t.PrimaryForeground, Hover: t.Primary.Alpha(0.9),
			Soft: t.AccentSoft, SoftForeground: t.AccentSoftForeground,
			SoftHover: t.AccentSoftHover, Foreground: t.AccentSoftForeground,
		}
	case ColorSuccess:
		return StatusPalette{
			Solid: t.Success, SolidForeground: solidForeground(t.Success, t.Background), Hover: t.SuccessHover,
			Soft: t.SuccessSoft, SoftForeground: t.SuccessSoftForeground,
			SoftHover: t.SuccessSoftHover, Foreground: t.SuccessSoftForeground,
		}
	case ColorWarning:
		return StatusPalette{
			Solid: t.Warning, SolidForeground: solidForeground(t.Warning, t.Background), Hover: t.WarningHover,
			Soft: t.WarningSoft, SoftForeground: t.WarningSoftForeground,
			SoftHover: t.WarningSoftHover, Foreground: t.WarningSoftForeground,
		}
	case ColorDanger:
		return StatusPalette{
			Solid: t.Danger, SolidForeground: t.DangerForeground, Hover: t.DangerHover,
			Soft: t.DangerSoft, SoftForeground: t.DangerSoftForeground,
			SoftHover: t.DangerSoftHover, Foreground: t.DangerSoftForeground,
		}
	case ColorInfo:
		return StatusPalette{
			Solid: t.Info, SolidForeground: solidForeground(t.Info, t.Background), Hover: t.InfoHover,
			Soft: t.InfoSoft, SoftForeground: t.InfoSoftForeground,
			SoftHover: t.InfoSoftHover, Foreground: t.InfoSoftForeground,
		}
	}
	return StatusPalette{
		Solid: t.Primary, SolidForeground: t.PrimaryForeground, Hover: t.Primary.Alpha(0.9),
		Soft: t.DefaultSoft, SoftForeground: t.DefaultSoftForeground,
		SoftHover: t.DefaultSoftHover, Foreground: t.Foreground,
	}
}

// solidForeground picks the higher-contrast black or white ink for a status
// fill; the status foreground tokens are intended for pale surfaces.
func solidForeground(fill, background ui.Color) ui.Color {
	fill = fill.Over(background)
	channel := func(v uint8) float64 {
		s := float64(v) / 255
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	luminance := 0.2126*channel(fill.R) + 0.7152*channel(fill.G) + 0.0722*channel(fill.B)
	if luminance > 0.179 {
		return ui.RGB(0, 0, 0)
	}
	return ui.RGB(255, 255, 255)
}
