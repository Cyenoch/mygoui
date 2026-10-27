package theme

import (
	"github.com/egoist/mygo/ui"
)

// Palette contains coss semantic colors. Translucent tokens stay translucent:
// borders and subtle fills composite against the actual surface beneath them.
// MutedBackground is CSS --muted; MutedForeground is CSS --muted-foreground.
type Palette struct {
	Background                 ui.Color
	Foreground                 ui.Color
	Card                       ui.Color
	CardForeground             ui.Color
	Popover                    ui.Color
	PopoverForeground          ui.Color
	Primary                    ui.Color
	PrimaryForeground          ui.Color
	Secondary                  ui.Color
	SecondaryForeground        ui.Color
	MutedBackground            ui.Color
	MutedForeground            ui.Color
	Accent                     ui.Color
	AccentForeground           ui.Color
	Destructive                ui.Color
	DestructiveForeground      ui.Color
	Info                       ui.Color
	InfoForeground             ui.Color
	Success                    ui.Color
	SuccessForeground          ui.Color
	Warning                    ui.Color
	WarningForeground          ui.Color
	Border                     ui.Color
	Input                      ui.Color
	Ring                       ui.Color
	Chart1                     ui.Color
	Chart2                     ui.Color
	Chart3                     ui.Color
	Chart4                     ui.Color
	Chart5                     ui.Color
	Sidebar                    ui.Color
	SidebarForeground          ui.Color
	SidebarPrimary             ui.Color
	SidebarPrimaryForeground   ui.Color
	SidebarAccent              ui.Color
	SidebarAccentForeground    ui.Color
	SidebarBorder              ui.Color
	SidebarRing                ui.Color
	Code                       ui.Color
	CodeForeground             ui.Color
	CodeHighlight              ui.Color
	BackgroundSecondary        ui.Color
	BackgroundTertiary         ui.Color
	BackgroundInverse          ui.Color
	Surface                    ui.Color
	SurfaceForeground          ui.Color
	SurfaceSecondary           ui.Color
	SurfaceSecondaryForeground ui.Color
	SurfaceTertiary            ui.Color
	SurfaceTertiaryForeground  ui.Color
	SurfaceHover               ui.Color
	Overlay                    ui.Color
	OverlayForeground          ui.Color
	Muted                      ui.Color
	Default                    ui.Color
	DefaultForeground          ui.Color
	DefaultHover               ui.Color
	DefaultSoft                ui.Color
	DefaultSoftForeground      ui.Color
	DefaultSoftHover           ui.Color
	AccentHover                ui.Color
	AccentSoft                 ui.Color
	AccentSoftForeground       ui.Color
	AccentSoftHover            ui.Color
	SuccessHover               ui.Color
	SuccessSoft                ui.Color
	SuccessSoftForeground      ui.Color
	SuccessSoftHover           ui.Color
	WarningHover               ui.Color
	WarningSoft                ui.Color
	WarningSoftForeground      ui.Color
	WarningSoftHover           ui.Color
	Danger                     ui.Color
	DangerForeground           ui.Color
	DangerHover                ui.Color
	DangerSoft                 ui.Color
	DangerSoftForeground       ui.Color
	DangerSoftHover            ui.Color
	InfoHover                  ui.Color
	InfoSoft                   ui.Color
	InfoSoftForeground         ui.Color
	InfoSoftHover              ui.Color
	Segment                    ui.Color
	SegmentForeground          ui.Color
	BorderSecondary            ui.Color
	BorderTertiary             ui.Color
	Separator                  ui.Color
	SeparatorSecondary         ui.Color
	SeparatorTertiary          ui.Color
	Selection                  ui.Color
	Focus                      ui.Color
	Link                       ui.Color
	FieldBackground            ui.Color
	FieldForeground            ui.Color
	FieldPlaceholder           ui.Color
	FieldBorder                ui.Color
	FieldHover                 ui.Color
	FieldFocus                 ui.Color
	FieldBorderHover           ui.Color
	FieldBorderFocus           ui.Color
	Backdrop                   ui.Color
	Scrollbar                  ui.Color
}

// PalettePatch overrides named tokens for one appearance. Nil inherits the
// default; a pointer to ui.Transparent deliberately removes a fill or border.
type PalettePatch struct {
	Background               *ui.Color
	Foreground               *ui.Color
	Card                     *ui.Color
	CardForeground           *ui.Color
	Popover                  *ui.Color
	PopoverForeground        *ui.Color
	Primary                  *ui.Color
	PrimaryForeground        *ui.Color
	Secondary                *ui.Color
	SecondaryForeground      *ui.Color
	MutedBackground          *ui.Color
	MutedForeground          *ui.Color
	Accent                   *ui.Color
	AccentForeground         *ui.Color
	Destructive              *ui.Color
	DestructiveForeground    *ui.Color
	Info                     *ui.Color
	InfoForeground           *ui.Color
	Success                  *ui.Color
	SuccessForeground        *ui.Color
	Warning                  *ui.Color
	WarningForeground        *ui.Color
	Border                   *ui.Color
	Input                    *ui.Color
	Ring                     *ui.Color
	Chart1                   *ui.Color
	Chart2                   *ui.Color
	Chart3                   *ui.Color
	Chart4                   *ui.Color
	Chart5                   *ui.Color
	Sidebar                  *ui.Color
	SidebarForeground        *ui.Color
	SidebarPrimary           *ui.Color
	SidebarPrimaryForeground *ui.Color
	SidebarAccent            *ui.Color
	SidebarAccentForeground  *ui.Color
	SidebarBorder            *ui.Color
	SidebarRing              *ui.Color
	Code                     *ui.Color
	CodeForeground           *ui.Color
	CodeHighlight            *ui.Color
	Surface                  *ui.Color
	Muted                    *ui.Color
	FieldBorder              *ui.Color
	Danger                   *ui.Color
	DangerForeground         *ui.Color
	Hover                    *ui.Color
	Selected                 *ui.Color
	Overlay                  *ui.Color
	Backdrop                 *ui.Color
}

// PaletteOverrides customizes light and dark modes independently.
type PaletteOverrides struct{ Light, Dark *PalettePatch }

func (o PaletteOverrides) patch(dark bool) *PalettePatch {
	if dark {
		return o.Dark
	}
	return o.Light
}

func (p *PalettePatch) apply(c *Palette) {
	pairs := []struct{ dst, src *ui.Color }{
		{&c.Background, p.Background},
		{&c.Foreground, p.Foreground},
		{&c.Card, p.Card},
		{&c.CardForeground, p.CardForeground},
		{&c.Popover, p.Popover},
		{&c.PopoverForeground, p.PopoverForeground},
		{&c.Primary, p.Primary},
		{&c.PrimaryForeground, p.PrimaryForeground},
		{&c.Secondary, p.Secondary},
		{&c.SecondaryForeground, p.SecondaryForeground},
		{&c.MutedBackground, p.MutedBackground},
		{&c.MutedForeground, p.MutedForeground},
		{&c.Accent, p.Accent},
		{&c.AccentForeground, p.AccentForeground},
		{&c.Destructive, p.Destructive},
		{&c.DestructiveForeground, p.DestructiveForeground},
		{&c.Info, p.Info},
		{&c.InfoForeground, p.InfoForeground},
		{&c.Success, p.Success},
		{&c.SuccessForeground, p.SuccessForeground},
		{&c.Warning, p.Warning},
		{&c.WarningForeground, p.WarningForeground},
		{&c.Border, p.Border},
		{&c.Input, p.Input},
		{&c.Ring, p.Ring},
		{&c.Chart1, p.Chart1},
		{&c.Chart2, p.Chart2},
		{&c.Chart3, p.Chart3},
		{&c.Chart4, p.Chart4},
		{&c.Chart5, p.Chart5},
		{&c.Sidebar, p.Sidebar},
		{&c.SidebarForeground, p.SidebarForeground},
		{&c.SidebarPrimary, p.SidebarPrimary},
		{&c.SidebarPrimaryForeground, p.SidebarPrimaryForeground},
		{&c.SidebarAccent, p.SidebarAccent},
		{&c.SidebarAccentForeground, p.SidebarAccentForeground},
		{&c.SidebarBorder, p.SidebarBorder},
		{&c.SidebarRing, p.SidebarRing},
		{&c.Code, p.Code},
		{&c.CodeForeground, p.CodeForeground},
		{&c.CodeHighlight, p.CodeHighlight},
		{&c.Surface, p.Surface},
		{&c.Muted, p.Muted},
		{&c.FieldBorder, p.FieldBorder},
		{&c.Danger, p.Danger},
		{&c.DangerForeground, p.DangerForeground},
		{&c.Overlay, p.Overlay},
		{&c.Backdrop, p.Backdrop},
	}
	for _, pair := range pairs {
		if pair.src != nil {
			*pair.dst = *pair.src
		}
	}
}

// clonePatch detaches all caller-owned color pointers at the publication boundary.
func clonePatch(p *PalettePatch) *PalettePatch {
	if p == nil {
		return nil
	}
	out := *p
	fields := []**ui.Color{
		&out.Background,
		&out.Foreground,
		&out.Card,
		&out.CardForeground,
		&out.Popover,
		&out.PopoverForeground,
		&out.Primary,
		&out.PrimaryForeground,
		&out.Secondary,
		&out.SecondaryForeground,
		&out.MutedBackground,
		&out.MutedForeground,
		&out.Accent,
		&out.AccentForeground,
		&out.Destructive,
		&out.DestructiveForeground,
		&out.Info,
		&out.InfoForeground,
		&out.Success,
		&out.SuccessForeground,
		&out.Warning,
		&out.WarningForeground,
		&out.Border,
		&out.Input,
		&out.Ring,
		&out.Chart1,
		&out.Chart2,
		&out.Chart3,
		&out.Chart4,
		&out.Chart5,
		&out.Sidebar,
		&out.SidebarForeground,
		&out.SidebarPrimary,
		&out.SidebarPrimaryForeground,
		&out.SidebarAccent,
		&out.SidebarAccentForeground,
		&out.SidebarBorder,
		&out.SidebarRing,
		&out.Code,
		&out.CodeForeground,
		&out.CodeHighlight,
		&out.Surface,
		&out.Muted,
		&out.FieldBorder,
		&out.Danger,
		&out.DangerForeground,
		&out.Hover,
		&out.Selected,
		&out.Overlay,
		&out.Backdrop,
	}
	for _, field := range fields {
		if *field != nil {
			value := **field
			*field = &value
		}
	}
	return &out
}

var white = ui.RGB(255, 255, 255)
var black = ui.RGB(0, 0, 0)

// buildPalette expresses coss's published styling tokens as native colors.
// The MIT apps/ui/content/docs/(root)/styling.mdx at revision
// 3d8ca28e7d90123e9e0f44665944dc240b026ccb is the token reference;
// docs/references.md records the docs site's runtime color differences.
func buildPalette(dark bool, prefs Preferences) Palette {
	ink := ui.Oklch(0.269, 0, 0)
	p := Palette{
		Background: white, Foreground: ink, Card: white, CardForeground: ink,
		Popover: white, PopoverForeground: ink, Primary: ink,
		PrimaryForeground: ui.Oklch(0.985, 0, 0),
		Secondary:         black.Alpha(0.04), SecondaryForeground: ink,
		MutedBackground: black.Alpha(0.04),
		MutedForeground: mix(ui.Oklch(0.556, 0, 0), black, 0.1),
		Accent:          black.Alpha(0.04), AccentForeground: ink,
		Destructive:           ui.Oklch(0.637, 0.237, 25.331),
		DestructiveForeground: ui.Oklch(0.505, 0.213, 27.518),
		Info:                  ui.Oklch(0.623, 0.214, 259.815), InfoForeground: ui.Oklch(0.488, 0.243, 264.376),
		Success: ui.Oklch(0.696, 0.17, 162.48), SuccessForeground: ui.Oklch(0.508, 0.118, 165.612),
		Warning: ui.Oklch(0.769, 0.188, 70.08), WarningForeground: ui.Oklch(0.555, 0.163, 48.998),
		Border: black.Alpha(0.08), Input: black.Alpha(0.10), Ring: ui.Oklch(0.708, 0, 0),
		Chart1: ui.Oklch(0.646, 0.222, 41.116), Chart2: ui.Oklch(0.6, 0.118, 184.704),
		Chart3: ui.Oklch(0.398, 0.07, 227.392), Chart4: ui.Oklch(0.828, 0.189, 84.429),
		Chart5:  ui.Oklch(0.769, 0.188, 70.08),
		Sidebar: ui.Oklch(0.985, 0, 0), SidebarPrimary: ink,
		SidebarPrimaryForeground: ui.Oklch(0.985, 0, 0),
		SidebarAccent:            black.Alpha(0.04), SidebarAccentForeground: ink,
		SidebarBorder: black.Alpha(0.06), SidebarRing: ui.Oklch(0.708, 0, 0),
		Code: white, CodeForeground: ink, CodeHighlight: black.Alpha(0.04),
	}
	if dark {
		p.Background = mix(ui.Oklch(0.145, 0, 0), white, 0.04)
		p.Foreground = ui.Oklch(0.97, 0, 0)
		p.Card = mix(p.Background, white, 0.02)
		p.CardForeground = p.Foreground
		p.Popover = mix(p.Background, white, 0.04)
		p.PopoverForeground = p.Foreground
		p.Primary = p.Foreground
		p.PrimaryForeground = ink
		p.Secondary = white.Alpha(0.04)
		p.SecondaryForeground = p.Foreground
		p.MutedBackground = white.Alpha(0.04)
		p.MutedForeground = mix(ui.Oklch(0.556, 0, 0), white, 0.1)
		p.Accent = white.Alpha(0.04)
		p.AccentForeground = p.Foreground
		p.Destructive = mix(p.Destructive, white, 0.1)
		p.DestructiveForeground = ui.Oklch(0.704, 0.191, 22.216)
		p.InfoForeground = ui.Oklch(0.707, 0.165, 254.624)
		p.SuccessForeground = ui.Oklch(0.765, 0.177, 163.223)
		p.WarningForeground = ui.Oklch(0.828, 0.189, 84.429)
		p.Border = white.Alpha(0.06)
		p.Input = white.Alpha(0.08)
		p.Ring = ui.Oklch(0.556, 0, 0)
		p.Chart1 = ui.Oklch(0.488, 0.243, 264.376)
		p.Chart2 = p.Success
		p.Chart3 = p.Warning
		p.Chart4 = ui.Oklch(0.627, 0.265, 303.9)
		p.Chart5 = ui.Oklch(0.645, 0.246, 16.439)
		p.Sidebar = mix(ui.Oklch(0.145, 0, 0), white, 0.03)
		p.SidebarPrimary = p.Primary
		p.SidebarPrimaryForeground = ink
		p.SidebarAccent = white.Alpha(0.04)
		p.SidebarAccentForeground = p.Foreground
		p.SidebarBorder = white.Alpha(0.05)
		p.Code = mix(p.Background, white, 0.02)
		p.CodeForeground = p.Foreground
		p.CodeHighlight = white.Alpha(0.04)
	}
	p.SidebarForeground = mix(p.Foreground, p.Sidebar, 0.36)
	if prefs.Accent.A != 0 {
		p.Primary = prefs.Accent
		p.PrimaryForeground = contrastInk(p.Primary)
	}
	patch := prefs.Palette.patch(dark)
	if patch != nil {
		patch.apply(&p)
		if patch.Foreground != nil {
			if patch.CardForeground == nil {
				p.CardForeground = p.Foreground
			}
			if patch.PopoverForeground == nil {
				p.PopoverForeground = p.Foreground
			}
			if patch.SecondaryForeground == nil {
				p.SecondaryForeground = p.Foreground
			}
			if patch.AccentForeground == nil {
				p.AccentForeground = p.Foreground
			}
			if patch.CodeForeground == nil {
				p.CodeForeground = p.Foreground
			}
		}
	}
	// These names are component-family conveniences, not a second token system.
	p.BackgroundSecondary = p.MutedBackground.Over(p.Background)
	p.BackgroundTertiary = p.Input.Over(p.Background)
	p.BackgroundInverse = p.Foreground
	p.Surface = p.Card
	p.SurfaceForeground = p.CardForeground
	p.SurfaceSecondary = p.BackgroundSecondary
	p.SurfaceSecondaryForeground = p.Foreground
	p.SurfaceTertiary = p.BackgroundTertiary
	p.SurfaceTertiaryForeground = p.Foreground
	p.SurfaceHover = p.Accent
	p.Overlay = p.Popover
	p.OverlayForeground = p.PopoverForeground
	p.Muted = p.MutedForeground
	p.Default = p.Secondary
	p.DefaultForeground = p.SecondaryForeground
	p.DefaultHover = opacity(p.Secondary, 0.9)
	p.DefaultSoft = p.Secondary
	p.DefaultSoftForeground = p.SecondaryForeground
	p.DefaultSoftHover = p.DefaultHover
	p.AccentHover = opacity(p.Accent, 0.9)
	p.AccentSoft = p.Accent
	p.AccentSoftForeground = p.AccentForeground
	p.AccentSoftHover = p.AccentHover
	alpha := float32(0.08)
	if dark {
		alpha = 0.16
	}
	p.SuccessHover = opacity(p.Success, 0.9)
	p.SuccessSoft = opacity(p.Success, alpha)
	p.SuccessSoftForeground = p.SuccessForeground
	p.SuccessSoftHover = opacity(p.Success, alpha+0.04)
	p.WarningHover = opacity(p.Warning, 0.9)
	p.WarningSoft = opacity(p.Warning, alpha)
	p.WarningSoftForeground = p.WarningForeground
	p.WarningSoftHover = opacity(p.Warning, alpha+0.04)
	p.Danger = p.Destructive
	p.DangerForeground = white
	p.DangerHover = opacity(p.Destructive, 0.9)
	p.DangerSoft = opacity(p.Destructive, alpha)
	p.DangerSoftForeground = p.DestructiveForeground
	p.DangerSoftHover = opacity(p.Destructive, alpha+0.04)
	p.InfoHover = opacity(p.Info, 0.9)
	p.InfoSoft = opacity(p.Info, alpha)
	p.InfoSoftForeground = p.InfoForeground
	p.InfoSoftHover = opacity(p.Info, alpha+0.04)
	p.Segment = p.Background
	if dark {
		p.Segment = p.Input.Over(p.Background)
	}
	p.SegmentForeground = p.Foreground
	p.BorderSecondary = opacity(p.Border, 1.5)
	p.BorderTertiary = opacity(p.Border, 2)
	p.Separator = p.Border
	p.SeparatorSecondary = p.BorderSecondary
	p.SeparatorTertiary = p.BorderTertiary
	p.Selection = opacity(p.Primary, 0.20)
	p.Focus = p.Ring
	p.Link = p.Foreground
	p.FieldBackground = p.Background
	if dark {
		p.FieldBackground = opacity(p.Input, 0.32)
	}
	p.FieldForeground = p.Foreground
	p.FieldPlaceholder = p.MutedForeground
	p.FieldBorder = p.Input
	p.FieldHover = p.FieldBackground
	p.FieldFocus = p.FieldBackground
	p.FieldBorderHover = p.Input
	p.FieldBorderFocus = p.Ring
	p.Backdrop = black.Alpha(0.32)
	p.Scrollbar = opacity(p.Foreground, 0.2)
	if patch != nil {
		applyConveniences(&p, patch)
	}
	return p
}

func applyConveniences(p *Palette, o *PalettePatch) {
	if o.Surface != nil {
		p.Surface = *o.Surface
		p.Card = *o.Surface
		p.FieldBackground = *o.Surface
	}
	if o.Muted != nil {
		p.Muted = *o.Muted
		p.MutedForeground = *o.Muted
		p.FieldPlaceholder = *o.Muted
	}
	if o.FieldBorder != nil {
		p.FieldBorder = *o.FieldBorder
		p.Input = *o.FieldBorder
	}
	if o.Danger != nil {
		p.Danger = *o.Danger
		p.Destructive = *o.Danger
		p.DangerHover = opacity(p.Danger, 0.9)
		p.DangerSoft = opacity(p.Danger, 0.08)
		p.DangerSoftHover = opacity(p.Danger, 0.12)
	}
	if o.DangerForeground != nil {
		p.DangerForeground = *o.DangerForeground
	}
	if o.Hover != nil {
		p.SurfaceHover = *o.Hover
		p.Default = *o.Hover
		p.DefaultHover = opacity(*o.Hover, 0.9)
	}
	if o.Selected != nil {
		p.Selection = *o.Selected
	}
	if o.Overlay != nil {
		p.Overlay = *o.Overlay
		p.Popover = *o.Overlay
	}
	if o.Backdrop != nil {
		p.Backdrop = *o.Backdrop
	}
}

func contrastInk(c ui.Color) ui.Color {
	if float32(c.R)*0.2126+float32(c.G)*0.7152+float32(c.B)*0.0722 > 150 {
		return ui.RGB(38, 38, 38)
	}
	return white
}

func opacity(c ui.Color, factor float32) ui.Color { return c.Alpha(factor) }

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
