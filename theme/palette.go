package theme

import "github.com/egoist/mygo/ui"

// Colors is HeroUI v3's semantic palette. Background/foreground pairs belong
// together; soft and hover values are resolved once, not at paint time.
type Palette struct {
	Background, Foreground                               ui.Color
	BackgroundSecondary, BackgroundTertiary              ui.Color
	BackgroundInverse                                    ui.Color
	Surface, SurfaceForeground                           ui.Color
	SurfaceSecondary, SurfaceSecondaryForeground         ui.Color
	SurfaceTertiary, SurfaceTertiaryForeground           ui.Color
	SurfaceHover                                         ui.Color
	Overlay, OverlayForeground                           ui.Color
	Muted                                                ui.Color
	Default, DefaultForeground, DefaultHover             ui.Color
	DefaultSoft, DefaultSoftForeground, DefaultSoftHover ui.Color
	Accent, AccentForeground, AccentHover                ui.Color
	AccentSoft, AccentSoftForeground, AccentSoftHover    ui.Color
	Success, SuccessForeground, SuccessHover             ui.Color
	SuccessSoft, SuccessSoftForeground, SuccessSoftHover ui.Color
	Warning, WarningForeground, WarningHover             ui.Color
	WarningSoft, WarningSoftForeground, WarningSoftHover ui.Color
	Danger, DangerForeground, DangerHover                ui.Color
	DangerSoft, DangerSoftForeground, DangerSoftHover    ui.Color
	// Info has no upstream HeroUI role; it starts as the accent family and a
	// host may separate it through PalettePatch.Info.
	Info, InfoForeground, InfoHover                  ui.Color
	InfoSoft, InfoSoftForeground, InfoSoftHover      ui.Color
	Segment, SegmentForeground                       ui.Color
	Border, BorderSecondary, BorderTertiary          ui.Color
	Separator, SeparatorSecondary, SeparatorTertiary ui.Color
	// Selection fills selected rows and selected text.
	Selection, Focus, Link                             ui.Color
	FieldBackground, FieldForeground, FieldPlaceholder ui.Color
	FieldBorder, FieldHover, FieldFocus                ui.Color
	FieldBorderHover, FieldBorderFocus                 ui.Color
	Backdrop                                           ui.Color
	Scrollbar                                          ui.Color
}

// PalettePatch replaces individual semantic tokens for one appearance. Nil
// fields keep the HeroUI default for that appearance; a non-nil field is
// applied verbatim, so new(ui.Transparent) expresses a borderless border.
// Derived tokens (hover shades, soft fills and the foreground of an injected
// ink) are rebuilt from the patched roles.
type PalettePatch struct {
	// Background is the canvas behind panels and overlays.
	Background *ui.Color
	// Surface is the face of panels, cards, sheets and toasts. Apply also
	// uses it as the native control face, as a host sheet serves both.
	Surface *ui.Color
	// Foreground is the primary ink and, unless named separately, the ink on
	// every surface role.
	Foreground *ui.Color
	// Muted is secondary ink, placeholders and help text.
	Muted *ui.Color
	// Border outlines controls. new(ui.Transparent) is valid and erases them.
	Border *ui.Color

	Accent, AccentForeground   *ui.Color
	Danger, DangerForeground   *ui.Color
	Success, SuccessForeground *ui.Color
	Warning, WarningForeground *ui.Color
	Info, InfoForeground       *ui.Color

	// Hover is the fill of hovered rows, menus and controls. It seeds the
	// neutral Default role those surfaces paint, and its hover level.
	Hover *ui.Color
	// Selected is the fill of selected rows and selected text.
	Selected *ui.Color
	// Overlay is the face of dialogs, modals and drawers. Unset follows
	// Surface, as HeroUI defines it.
	Overlay *ui.Color
	// Backdrop is the scrim behind overlays.
	Backdrop *ui.Color
}

// PaletteOverrides carries one patch per appearance. A nil patch keeps the
// HeroUI defaults for that appearance.
type PaletteOverrides struct {
	Light, Dark *PalettePatch
}

func (o PaletteOverrides) patch(dark bool) *PalettePatch {
	if dark {
		return o.Dark
	}
	return o.Light
}

// apply writes the named tokens. It runs before the derived tokens so that
// hover, soft and foreground levels follow the injected base colors.
func (o *PalettePatch) apply(p *Palette) {
	set := func(dst *ui.Color, src *ui.Color) {
		if src != nil {
			*dst = *src
		}
	}
	set(&p.Background, o.Background)
	set(&p.Surface, o.Surface)
	set(&p.Foreground, o.Foreground)
	set(&p.Muted, o.Muted)
	set(&p.Border, o.Border)
	set(&p.Accent, o.Accent)
	set(&p.AccentForeground, o.AccentForeground)
	set(&p.Danger, o.Danger)
	set(&p.DangerForeground, o.DangerForeground)
	set(&p.Success, o.Success)
	set(&p.SuccessForeground, o.SuccessForeground)
	set(&p.Warning, o.Warning)
	set(&p.WarningForeground, o.WarningForeground)
	set(&p.Info, o.Info)
	set(&p.InfoForeground, o.InfoForeground)
	set(&p.SurfaceHover, o.Hover)
	set(&p.Default, o.Hover)
	set(&p.Selection, o.Selected)
	set(&p.Overlay, o.Overlay)
	set(&p.Backdrop, o.Backdrop)
}

var (
	white   = ui.Oklch(1, 0, 0)
	snow    = ui.Oklch(0.9911, 0, 0)
	eclipse = ui.Oklch(0.2103, 0.0059, 285.89)
)

// buildPalette translates packages/styles/themes/default/variables.css and its
// calculated colors. Values are upstream oklch literals; the CSS color-mix
// results are reproduced with premultiplied sRGB mixing.
func buildPalette(dark bool, prefs Preferences) Palette {
	var p Palette
	p.Accent = ui.Oklch(0.6204, 0.195, 253.83)
	p.AccentForeground = snow
	p.Success = ui.Oklch(0.7329, 0.1935, 150.81)
	p.SuccessForeground = eclipse
	if prefs.Accent.A != 0 {
		p.Accent = prefs.Accent
		p.Accent.A = 255
	}

	if dark {
		p.Background = ui.Oklch(0.12, 0.005, 285.823)
		p.Foreground = snow
		p.Surface = ui.Oklch(0.2103, 0.0059, 285.89)
		p.SurfaceForeground = p.Foreground
		p.SurfaceSecondary = ui.Oklch(0.257, 0.0037, 286.14)
		p.SurfaceSecondaryForeground = p.Foreground
		p.SurfaceTertiary = ui.Oklch(0.2721, 0.0024, 247.91)
		p.SurfaceTertiaryForeground = p.Foreground
		p.Overlay = p.Surface
		p.OverlayForeground = p.Foreground
		p.Muted = ui.Oklch(0.705, 0.015, 286.067)
		p.Default = ui.Oklch(0.274, 0.006, 286.033)
		p.DefaultForeground = snow
		p.Warning = ui.Oklch(0.8203, 0.1388, 76.34)
		p.WarningForeground = eclipse
		p.Danger = ui.Oklch(0.594, 0.1967, 24.63)
		p.DangerForeground = snow
		p.Segment = ui.Oklch(0.3964, 0.01, 285.93)
		p.SegmentForeground = p.Foreground
		p.Border = ui.Oklch(0.28, 0.006, 286.033)
		p.Separator = ui.Oklch(0.25, 0.006, 286.033)
		p.Backdrop = ui.RGBA(0, 0, 0, 0.6)
	} else {
		p.Background = ui.Oklch(0.9702, 0, 0)
		p.Foreground = eclipse
		p.Surface = white
		p.SurfaceForeground = p.Foreground
		p.SurfaceSecondary = ui.Oklch(0.9524, 0.0013, 286.37)
		p.SurfaceSecondaryForeground = p.Foreground
		p.SurfaceTertiary = ui.Oklch(0.9373, 0.0013, 286.37)
		p.SurfaceTertiaryForeground = p.Foreground
		p.Overlay = white
		p.OverlayForeground = p.Foreground
		p.Muted = ui.Oklch(0.5517, 0.0138, 285.94)
		p.Default = ui.Oklch(0.94, 0.001, 286.375)
		p.DefaultForeground = eclipse
		p.Warning = ui.Oklch(0.7819, 0.1585, 72.33)
		p.WarningForeground = eclipse
		p.Danger = ui.Oklch(0.6532, 0.2328, 25.74)
		p.DangerForeground = snow
		p.Segment = white
		p.SegmentForeground = eclipse
		p.Border = ui.Oklch(0.90, 0.004, 286.32)
		p.Separator = ui.Oklch(0.92, 0.004, 286.32)
		p.Backdrop = ui.RGBA(0, 0, 0, 0.5)
	}

	// HeroUI has no separate info role: accent is its info color until a host
	// separates them through PalettePatch.Info.
	p.Info, p.InfoForeground = p.Accent, p.AccentForeground

	patch := prefs.Palette.patch(dark)
	if patch != nil {
		patch.apply(&p)
	}
	// An injected ink reaches every role that pairs with it; HeroUI's own
	// seeds already agree, so this is a no-op without a patch.
	p.SurfaceForeground = p.Foreground
	p.SurfaceSecondaryForeground = p.Foreground
	p.SurfaceTertiaryForeground = p.Foreground
	p.SegmentForeground = p.Foreground
	p.DefaultForeground = p.Foreground
	if patch == nil || patch.Overlay == nil {
		p.Overlay = p.Surface
	}
	p.OverlayForeground = p.Foreground

	p.FieldBackground = p.Surface
	p.FieldForeground = p.Foreground
	p.FieldPlaceholder = p.Muted
	p.FieldBorder = ui.Transparent
	p.Focus = p.Accent
	p.Link = p.Foreground
	p.Scrollbar = p.Foreground.Alpha(0.15)

	hover := patch != nil && patch.Hover != nil
	if !hover {
		p.SurfaceHover = mix(p.Surface, p.SurfaceForeground, 0.08)
	}
	p.BackgroundSecondary = mix(p.Background, p.Foreground, 0.04)
	p.BackgroundTertiary = mix(p.Background, p.Foreground, 0.08)
	p.BackgroundInverse = p.Foreground
	p.DefaultHover = mix(p.Default, p.DefaultForeground, 0.04)
	p.AccentHover = mix(p.Accent, p.AccentForeground, 0.10)
	p.SuccessHover = mix(p.Success, p.SuccessForeground, 0.10)
	p.WarningHover = mix(p.Warning, p.WarningForeground, 0.10)
	p.DangerHover = mix(p.Danger, p.DangerForeground, 0.10)
	p.InfoHover = mix(p.Info, p.InfoForeground, 0.10)
	if patch == nil || patch.Selected == nil {
		p.Selection = p.Accent.Alpha(0.30)
	}

	p.FieldHover = mix(p.FieldBackground, p.FieldForeground, 0.02)
	p.FieldFocus = p.FieldBackground
	p.FieldBorderHover = mix(p.FieldBorder, p.FieldForeground, 0.10)
	p.FieldBorderFocus = mix(p.FieldBorder, p.FieldForeground, 0.22)

	p.DefaultSoft = p.Default.Alpha(0.5)
	p.DefaultSoftForeground = p.DefaultForeground
	p.DefaultSoftHover = p.Default.Alpha(0.6)

	softAlpha, hoverAlpha := float32(0.15), float32(0.20)
	if dark {
		softAlpha, hoverAlpha = 0.12, 0.16
	}
	p.AccentSoft = p.Accent.Alpha(softAlpha)
	p.AccentSoftHover = p.Accent.Alpha(hoverAlpha)
	p.SuccessSoft = p.Success.Alpha(softAlpha)
	p.SuccessSoftHover = p.Success.Alpha(hoverAlpha)
	p.WarningSoft = p.Warning.Alpha(softAlpha)
	p.WarningSoftHover = p.Warning.Alpha(hoverAlpha)
	p.InfoSoft = p.Info.Alpha(softAlpha)
	p.InfoSoftHover = p.Info.Alpha(hoverAlpha)
	if dark {
		p.DangerSoft = p.Danger.Alpha(0.15)
		p.DangerSoftHover = p.Danger.Alpha(0.20)
	} else {
		p.DangerSoft = p.Danger.Alpha(softAlpha)
		p.DangerSoftHover = p.Danger.Alpha(hoverAlpha)
	}

	softWeight := float32(0.30)
	if dark {
		softWeight = 0.273
	}
	p.AccentSoftForeground = softForeground(p.Accent, p.Foreground, prefs.Vibrant, softWeight)
	p.SuccessSoftForeground = softForeground(p.Success, p.Foreground, prefs.Vibrant, softWeight)
	p.WarningSoftForeground = softForeground(p.Warning, p.Foreground, prefs.Vibrant, softWeight)
	p.DangerSoftForeground = softForeground(p.Danger, p.Foreground, prefs.Vibrant, softWeight)
	p.InfoSoftForeground = softForeground(p.Info, p.Foreground, prefs.Vibrant, softWeight)
	if !prefs.Vibrant && !dark {
		p.WarningSoftForeground = mix(p.Warning, p.Foreground, 0.467)
		p.SuccessSoftForeground = mix(p.Success, p.Foreground, 0.429)
		p.DangerSoftForeground = mix(p.Danger, p.Foreground, 0.364)
	}

	p.SeparatorSecondary = mix(p.Surface, p.SurfaceForeground, 0.15)
	p.SeparatorTertiary = mix(p.Surface, p.SurfaceForeground, 0.19)
	p.BorderSecondary = mix(p.Surface, p.SurfaceForeground, 0.22)
	p.BorderTertiary = mix(p.Surface, p.SurfaceForeground, 0.34)
	return p
}

func softForeground(color, fg ui.Color, vibrant bool, weight float32) ui.Color {
	if vibrant {
		weight = 0.08
	}
	return mix(color, fg, weight)
}

// mix blends in premultiplied sRGB. It handles straight alpha exactly once.
func mix(a, b ui.Color, t float32) ui.Color {
	t = max(0, min(t, 1))
	aw, bw := float32(a.A)*(1-t), float32(b.A)*t
	alpha := aw + bw
	if alpha == 0 {
		return ui.Color{}
	}
	channel := func(x, y uint8) uint8 {
		return uint8((float32(x)*aw+float32(y)*bw)/alpha + 0.5)
	}
	return ui.Color{
		R: channel(a.R, b.R), G: channel(a.G, b.G), B: channel(a.B, b.B),
		A: uint8(alpha + 0.5),
	}
}
