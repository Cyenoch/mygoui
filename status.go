package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// StatusPalette is the resolved color set for one HeroUI status. Components
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
			Solid: t.Accent, SolidForeground: t.AccentForeground, Hover: t.AccentHover,
			Soft: t.AccentSoft, SoftForeground: t.AccentSoftForeground,
			SoftHover: t.AccentSoftHover, Foreground: t.AccentSoftForeground,
		}
	case ColorSuccess:
		return StatusPalette{
			Solid: t.Success, SolidForeground: t.SuccessForeground, Hover: t.SuccessHover,
			Soft: t.SuccessSoft, SoftForeground: t.SuccessSoftForeground,
			SoftHover: t.SuccessSoftHover, Foreground: t.SuccessSoftForeground,
		}
	case ColorWarning:
		return StatusPalette{
			Solid: t.Warning, SolidForeground: t.WarningForeground, Hover: t.WarningHover,
			Soft: t.WarningSoft, SoftForeground: t.WarningSoftForeground,
			SoftHover: t.WarningSoftHover, Foreground: t.WarningSoftForeground,
		}
	case ColorDanger:
		return StatusPalette{
			Solid: t.Danger, SolidForeground: t.DangerForeground, Hover: t.DangerHover,
			Soft: t.DangerSoft, SoftForeground: t.DangerSoftForeground,
			SoftHover: t.DangerSoftHover, Foreground: t.DangerSoftForeground,
		}
	}
	return StatusPalette{
		Solid: t.Default, SolidForeground: t.DefaultForeground, Hover: t.DefaultHover,
		Soft: t.DefaultSoft, SoftForeground: t.DefaultSoftForeground,
		SoftHover: t.DefaultSoftHover, Foreground: t.Foreground,
	}
}
