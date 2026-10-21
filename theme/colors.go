package theme

import (
	"math"
	"sync"

	"github.com/egoist/mygo/ui"
)

// ColorTokens is the native semantic vocabulary for reusable mygoui.
// Surface/foreground pairs belong together. Accent is a neutral interaction
// surface; MyGo's actionable Theme.Accent maps to Primary instead.
// Values may carry straight alpha: paint once, or use Over with the real base.
type ColorTokens struct {
	Background, Foreground                     ui.Color
	Card, CardForeground                       ui.Color
	Popover, PopoverForeground                 ui.Color
	Primary, PrimaryForeground                 ui.Color
	Secondary, SecondaryForeground             ui.Color
	Muted, MutedForeground                     ui.Color
	Accent, AccentForeground                   ui.Color
	Destructive, DestructiveForeground         ui.Color
	Border, Input, Ring                        ui.Color
	Sidebar, SidebarForeground                 ui.Color
	SidebarPrimary, SidebarPrimaryForeground   ui.Color
	SidebarAccent, SidebarAccentForeground     ui.Color
	SidebarBorder, SidebarRing                 ui.Color
	Chart1, Chart2, Chart3, Chart4, Chart5     ui.Color
	PrimaryHover, PrimaryPressed               ui.Color
	SecondaryHover, SecondaryPressed           ui.Color
	AccentHover, AccentPressed                 ui.Color
	DestructiveHover, DestructivePressed       ui.Color
	DisabledForeground, DisabledSurface        ui.Color
	Selection, SelectionForeground             ui.Color
	Success, SuccessForeground                 ui.Color
	Warning, WarningForeground                 ui.Color
	Info, InfoForeground                       ui.Color
	CodeBackground, CodeForeground, CodeBorder ui.Color
	Bubble, DiffAdd, DiffDelete                ui.Color
	WarningSurface, ErrorSurface               ui.Color
	SidebarRowHover, SidebarRowSelected        ui.Color
	Overlay, Shadow                            ui.Color
}

// Preset chooses source defaults without changing the semantic vocabulary.
type Preset uint8

const (
	// Blue provides blue primary actions and subdued neutral surfaces.
	Blue Preset = iota
	// Neutral uses shadcn's neutral theme with the same native accessibility policy.
	Neutral
)

var (
	white        = ui.RGB(255, 255, 255)
	black        = ui.RGB(0, 0, 0)
	blueLight    = makeColors(false, Blue)
	blueDark     = makeColors(true, Blue)
	neutralLight = makeColors(false, Neutral)
	neutralDark  = makeColors(true, Neutral)
)

func sourceColors(dark bool, preset Preset) ColorTokens {
	if preset == Neutral {
		if dark {
			return neutralDark
		}
		return neutralLight
	}
	if dark {
		return blueDark
	}
	return blueLight
}

// ColorsForMode supplies a preset for static appearance thumbnails.
// Live components MUST use Colors so system and scoped theme changes apply.
func ColorsForMode(dark bool, preset Preset) ColorTokens { return sourceColors(dark, preset) }

func makeColors(dark bool, preset Preset) ColorTokens {
	var p ColorTokens
	if dark {
		p.Background = ui.RGB(10, 10, 10)
		p.Foreground = ui.Oklch(0.97, 0, 0)
		p.Card = mix(p.Background, white, 0.03)
		p.Secondary = p.Card
		p.Muted = white.Alpha(0.03)
		p.MutedForeground = mix(ui.Oklch(0.556, 0, 0), white, 0.1)
		p.Accent = white.Alpha(0.04)
		p.Border = white.Alpha(0.06)
		p.Input = white.Alpha(0.08)
		p.Primary = ui.Oklch(0.571, 0.21, 264)
		p.Destructive = ui.Oklch(0.704, 0.191, 22.216)
		p.Warning = ui.Oklch(0.769, 0.188, 70.08)
		p.Success = ui.Oklch(0.696, 0.17, 162.48)
		p.Chart1 = ui.Oklch(0.488, 0.243, 264.376)
		p.Chart2 = p.Success
		p.Chart3 = p.Warning
		p.Chart4 = ui.Oklch(0.627, 0.265, 303.9)
		p.Chart5 = ui.Oklch(0.645, 0.246, 16.439)
		if preset == Neutral {
			p.Background = ui.Oklch(0.145, 0, 0)
			p.Foreground = ui.Oklch(0.985, 0, 0)
			p.Card = ui.Oklch(0.205, 0, 0)
			p.Primary = ui.Oklch(0.922, 0, 0)
			p.Secondary = ui.Oklch(0.269, 0, 0)
			p.Muted, p.Accent = p.Secondary, p.Secondary
			p.MutedForeground = ui.Oklch(0.708, 0, 0)
			p.Border, p.Input = white.Alpha(0.1), white.Alpha(0.15)
		}
	} else {
		p.Background = ui.Oklch(0.992, 0, 0)
		p.Foreground = ui.Oklch(0.274, 0.006, 286.033)
		p.Card = white
		p.Secondary = ui.Oklch(0.985, 0, 0)
		p.Muted = ui.Oklch(0.967, 0.001, 286.375)
		p.MutedForeground = ui.Oklch(0.552, 0.016, 285.938)
		p.Accent = p.Muted
		p.Border = ui.Oklch(0.92, 0.004, 286.32)
		p.Input = p.Border
		p.Primary = ui.Oklch(0.488, 0.217, 264)
		p.Destructive = ui.Oklch(0.577, 0.245, 27.325)
		p.Warning = ui.Oklch(0.53, 0.12, 70.08)
		p.Success = ui.Oklch(0.50, 0.12, 162.48)
		p.Chart1 = ui.Oklch(0.646, 0.222, 41.116)
		p.Chart2 = ui.Oklch(0.6, 0.118, 184.704)
		p.Chart3 = ui.Oklch(0.398, 0.07, 227.392)
		p.Chart4 = ui.Oklch(0.828, 0.189, 84.429)
		p.Chart5 = ui.Oklch(0.769, 0.188, 70.08)
		if preset == Neutral {
			p.Background = white
			p.Foreground = ui.Oklch(0.145, 0, 0)
			p.Primary = ui.Oklch(0.205, 0, 0)
			p.Secondary = ui.Oklch(0.97, 0, 0)
			p.Muted, p.Accent = p.Secondary, p.Secondary
			p.MutedForeground = ui.Oklch(0.556, 0, 0)
			p.Border = ui.Oklch(0.922, 0, 0)
			p.Input = p.Border
		}
	}
	if preset == Neutral {
		// The current registry's neutral chart scale differs from the older
		// colorful scaffold still shown on the theming documentation page.
		p.Chart1 = ui.Oklch(0.87, 0, 0)
		p.Chart2 = ui.Oklch(0.556, 0, 0)
		p.Chart3 = ui.Oklch(0.439, 0, 0)
		p.Chart4 = ui.Oklch(0.371, 0, 0)
		p.Chart5 = ui.Oklch(0.269, 0, 0)
	}
	p.Popover = p.Card
	p.CardForeground, p.PopoverForeground = p.Foreground, p.Foreground
	p.SecondaryForeground, p.AccentForeground = p.Foreground, p.Foreground
	p.Sidebar, p.SidebarForeground = p.Secondary, p.Foreground
	if dark {
		p.Sidebar = p.Card
	}
	if preset == Neutral && !dark {
		p.Sidebar = ui.Oklch(0.985, 0, 0)
	}
	p.SidebarAccent, p.SidebarAccentForeground = p.Accent, p.Foreground
	p.SidebarRowHover = p.Accent
	p.SidebarRowSelected = p.Card
	if dark {
		p.SidebarRowSelected = white.Alpha(0.03)
	}
	if preset == Neutral {
		p.SidebarRowSelected = p.Accent
	}
	p.Bubble = p.Accent
	p.CodeBackground = mix(p.Card, p.Background, 0.1)
	p.CodeForeground = p.Foreground
	p.CodeBorder = p.Border
	p.Info = p.Primary
	p.Overlay, p.Shadow = black.Alpha(0.5), black.Alpha(0.4)
	if preset == Neutral && dark {
		p.SidebarPrimary = ui.Oklch(0.488, 0.243, 264.376)
	}
	finishColors(&p, dark, false)
	return p
}

// finishColors resolves readable state groups once, never in a paint callback.
func finishColors(p *ColorTokens, dark, highContrast bool) {
	p.MutedForeground = readableAcross(p.MutedForeground, p.Foreground, 4.5,
		p.Background, p.Card, p.Muted.Over(p.Background), p.Accent.Over(p.Background), p.Sidebar)
	statusAlpha := float32(0.08)
	if dark {
		statusAlpha = 0.16
	}
	p.WarningSurface = p.Warning.Alpha(statusAlpha)
	p.ErrorSurface = p.Destructive.Alpha(statusAlpha)
	p.Destructive = readableAcross(p.Destructive, p.Foreground, 4.5,
		p.Background, p.Card, p.ErrorSurface.Over(p.Background), p.ErrorSurface.Over(p.Card))
	p.Success = readableAcross(p.Success, p.Foreground, 4.5, p.Background, p.Card)
	p.Warning = readableAcross(p.Warning, p.Foreground, 4.5,
		p.Background, p.Card, p.WarningSurface.Over(p.Background), p.WarningSurface.Over(p.Card))
	p.Info = readableAcross(p.Info, p.Foreground, 4.5, p.Background, p.Card)
	p.PrimaryHover = mix(p.Primary, p.Background, 0.08)
	p.PrimaryPressed = mix(p.Primary, p.Foreground, 0.12)
	p.PrimaryForeground = groupInk(p.Primary)
	// Keep every state on the same side of the common ink's contrast boundary.
	p.PrimaryHover = readableAcross(p.PrimaryHover, p.Primary, 4.5, p.PrimaryForeground)
	p.PrimaryPressed = readableAcross(p.PrimaryPressed, p.Primary, 4.5, p.PrimaryForeground)
	p.SecondaryHover = mix(p.Secondary.Over(p.Background), p.Foreground, 0.05)
	p.SecondaryPressed = mix(p.Secondary.Over(p.Background), p.Foreground, 0.09)
	p.AccentHover = p.Accent
	p.AccentPressed = mix(p.Accent.Over(p.Background), p.Foreground, 0.08)
	p.DestructiveHover = mix(p.Destructive, p.Background, 0.06)
	p.DestructivePressed = mix(p.Destructive, p.Foreground, 0.1)
	p.DestructiveForeground = groupInk(p.Destructive)
	p.DestructiveHover = readableAcross(p.DestructiveHover, p.Destructive, 4.5, p.DestructiveForeground)
	p.DestructivePressed = readableAcross(p.DestructivePressed, p.Destructive, 4.5, p.DestructiveForeground)
	p.SuccessForeground = groupInk(p.Success)
	p.WarningForeground = groupInk(p.Warning)
	p.InfoForeground = groupInk(p.Info)
	p.DisabledSurface = p.Muted.Over(p.Background)
	p.DisabledForeground = readableAcross(mix(p.Foreground, p.DisabledSurface, 0.5), p.Foreground, 3, p.DisabledSurface)
	p.Selection = p.Primary.Alpha(0.25)
	if dark {
		p.Selection = p.Primary.Alpha(0.4)
	}
	p.SelectionForeground = readableAcross(p.Foreground, groupInk(p.Selection.Over(p.Background)), 4.5, p.Selection.Over(p.Background))
	p.Ring = readableAcross(p.Primary, p.Foreground, 3, p.Background, p.Card)
	if p.SidebarPrimary.A == 0 {
		p.SidebarPrimary = p.Primary
	}
	p.SidebarPrimaryForeground = groupInk(p.SidebarPrimary)
	p.SidebarBorder, p.SidebarRing = p.Border, readableAcross(p.Ring, p.Foreground, 3, p.Sidebar)
	p.DiffAdd, p.DiffDelete = p.Success, p.Destructive
	// Subdued surfaces retain the source status hue; their semantic ink above
	// is resolved against the actual composites, not alpha-dimmed afterward.
	if highContrast {
		p.Border = readableAcross(p.Border.Over(p.Background), p.Foreground, 3, p.Background, p.Card)
		p.Input, p.CodeBorder, p.SidebarBorder = p.Border, p.Border, p.Border
		p.MutedForeground = readableAcross(p.MutedForeground, p.Foreground, 7, p.Background, p.Card, p.Sidebar)
		p.DisabledForeground = readableAcross(p.DisabledForeground, p.Foreground, 4.5, p.DisabledSurface)
	}
}

func groupInk(surfaces ...ui.Color) ui.Color {
	light, dark := float64(21), float64(21)
	for _, surface := range surfaces {
		light = min(light, contrast(white, surface))
		dark = min(dark, contrast(black, surface))
	}
	if light >= dark {
		return white
	}
	return black
}

func readableAcross(ink, target ui.Color, minimum float64, surfaces ...ui.Color) ui.Color {
	for step := 0; step <= 255; step++ {
		candidate := mix(ink, target, float32(step)/255)
		valid := true
		for _, surface := range surfaces {
			if contrast(candidate.Over(surface), surface) < minimum {
				valid = false
				break
			}
		}
		if valid {
			return candidate
		}
	}
	return target
}

func contrast(a, b ui.Color) float64 {
	l1, l2 := luminance(a), luminance(b)
	return (max(l1, l2) + 0.05) / (min(l1, l2) + 0.05)
}

func luminance(c ui.Color) float64 {
	channel := func(v uint8) float64 {
		s := float64(v) / 255
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
}

type resolvedColors struct {
	theme        ui.Theme
	colors       ColorTokens
	highContrast bool
}

// A bounded cache retains immutable values, never Contexts or Elements. The
// normal Apply -> Colors path compares pointers; allocations occur only when a
// theme changes or a previously unseen scoped override is installed.
var colorCache struct {
	sync.RWMutex
	entries [16]*resolvedColors
	next    uint8
}

func cachedColors(t *ui.Theme, highContrast bool) *resolvedColors {
	colorCache.RLock()
	defer colorCache.RUnlock()
	for _, p := range colorCache.entries {
		if p != nil && (p.highContrast || !highContrast) && (&p.theme == t || p.theme == *t) {
			return p
		}
	}
	return nil
}

func rememberColors(t ui.Theme, colors ColorTokens, highContrast bool) *resolvedColors {
	p := &resolvedColors{theme: t, colors: colors, highContrast: highContrast}
	colorCache.Lock()
	colorCache.entries[colorCache.next] = p
	colorCache.next = (colorCache.next + 1) % uint8(len(colorCache.entries))
	colorCache.Unlock()
	return p
}

// Colors returns the resolved semantic colors for this context. Apply or
// ApplyPreset at the view root; scoped SetTheme overrides remain authoritative
// for native shared roles, and accessibility also applies to library extras.
func Colors(c *ui.Context) ColorTokens {
	t := c.Theme()
	highContrast := c.Preferences().HighContrast
	if p := cachedColors(t, highContrast); p != nil {
		return p.colors
	}
	preset := Neutral
	if t.Background == sourceColors(t.Dark, Blue).Background {
		preset = Blue
	}
	p := sourceColors(t.Dark, preset)
	p.Background, p.Foreground = t.Background, t.Text
	p.Card, p.CardForeground, p.Popover, p.PopoverForeground = t.Surface, t.Text, t.Surface, t.Text
	p.SecondaryForeground, p.SidebarForeground, p.SidebarAccentForeground = t.Text, t.Text, t.Text
	// Scoped native themes share Card for the supporting action surface;
	// unrepresented component surfaces retain the explicit preset defaults.
	p.Secondary = t.Surface
	p.SidebarPrimary = t.Accent
	p.Bubble = t.SurfaceHover
	p.Accent, p.AccentForeground, p.AccentPressed = t.SurfaceHover, t.Text, t.SurfacePressed
	p.Primary, p.PrimaryForeground = t.Accent, t.AccentText
	p.MutedForeground, p.Border, p.Input, p.CodeBorder = t.TextMuted, t.Border, t.Border, t.Border
	p.Destructive, p.Success, p.Warning, p.Info = t.Danger, t.Success, t.Warning, t.Accent
	finishColors(&p, t.Dark, highContrast)
	p.Primary, p.PrimaryHover, p.PrimaryPressed = t.Accent, t.AccentHover, t.AccentPressed
	p.PrimaryForeground = groupInk(p.Primary)
	if contrast(t.AccentText, p.Primary) >= 4.5 {
		p.PrimaryForeground = t.AccentText
	}
	p.PrimaryHover = readableAcross(p.PrimaryHover, p.Primary, 4.5, p.PrimaryForeground)
	p.PrimaryPressed = readableAcross(p.PrimaryPressed, p.Primary, 4.5, p.PrimaryForeground)
	p.AccentPressed = t.SurfacePressed
	p.Selection, p.SelectionForeground = t.Selection, readableAcross(t.Text, groupInk(t.Selection.Over(t.Background)), 4.5, t.Selection.Over(t.Background))
	p.Ring = readableAcross(t.Focus.Over(t.Background), t.Text, 3, t.Background, t.Surface)
	p.SidebarPrimary, p.SidebarPrimaryForeground = p.Primary, p.PrimaryForeground
	p.SidebarRing = readableAcross(p.Ring, t.Text, 3, p.Sidebar)
	p.CodeBackground, p.CodeForeground = mix(p.Card, p.Background, 0.1), p.Foreground
	return rememberColors(*t, p, highContrast).colors
}

// HighContrast reports the resolved accessibility policy, including explicit
// previews. Hosts can strengthen extension colors using the same policy.
func HighContrast(c *ui.Context) bool {
	if c.Preferences().HighContrast {
		return true
	}
	if p := cachedColors(c.Theme(), false); p != nil {
		return p.highContrast
	}
	return false
}
