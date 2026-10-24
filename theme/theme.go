// Package theme owns HeroUI v3's semantic tokens, geometry and typography.
//
// Tokens are derived once from upstream and published as immutable snapshots:
// components read live values with Colors(c) and never cache a palette keyed by
// window. Host customization goes through Configure before Apply.
package theme

import (
	"sync/atomic"

	"github.com/egoist/mygo/ui"
)

// Preferences is the process-wide HeroUI customization surface. Zero fields
// select upstream defaults: the system font, 16 DIP base text, an 8 DIP radius
// base and the default accent.
//
// Semantic colors resolve by precedence, highest first: Preferences.Palette,
// Preferences.Accent, HeroUI's published palette. Apply installs the theme
// composed from that resolution, so an injected palette survives a view-root
// Apply and reaches native widgets (selections, scrollbars, text cursors)
// too. Library paint reads Colors, never the context theme, so a host that
// styles native widgets through its own scoped theme either installs that
// theme after Apply or skips Apply.
type Preferences struct {
	FontFamily string
	FontSize   float32
	// Radius is the base radius; the HeroUI scale is derived from it.
	Radius float32
	// Accent replaces the default blue accent when opaque.
	Accent ui.Color
	// Vibrant selects the upstream saturated soft-foreground palette.
	Vibrant bool
	// Palette overrides individual semantic tokens per appearance. The zero
	// value keeps the HeroUI defaults. Configure replaces the published
	// preferences as a whole, so a later call that changes only typography
	// must publish the same Palette again; UpdatePreferences edits in place.
	Palette PaletteOverrides
}

type state struct {
	prefs  Preferences
	family string
	light  Palette
	dark   Palette
}

var current atomic.Pointer[state]

func init() { Configure(Preferences{}) }

// Configure publishes host preferences and recomputes both palettes. Call on
// the UI thread, then invalidate affected windows.
func Configure(p Preferences) {
	if p.FontSize == 0 {
		p.FontSize = 16
	}
	if p.Radius == 0 {
		p.Radius = 8
	}
	if p.Accent.A == 0 {
		p.Accent = ui.Color{}
	}
	family := Sans
	if p.FontFamily != "" {
		family = p.FontFamily + ", " + Sans
	}
	if existing := current.Load(); existing != nil &&
		existing.prefs == p && existing.family == family {
		return
	}
	current.Store(&state{
		prefs:  p,
		family: family,
		light:  buildPalette(false, p),
		dark:   buildPalette(true, p),
	})
}

func snapshot() *state {
	s := current.Load()
	if s == nil {
		Configure(Preferences{})
		s = current.Load()
	}
	return s
}

// UpdatePreferences edits the published preferences and republishes them, so
// a host that changes only typography keeps its injected palette.
func UpdatePreferences(edit func(*Preferences)) {
	p := PreferencesSnapshot()
	edit(&p)
	Configure(p)
}

// PreferencesSnapshot returns the published preferences.
func PreferencesSnapshot() Preferences { return snapshot().prefs }

// FontFamily returns the configured family followed by native fallbacks.
func FontFamily(*ui.Context) string { return snapshot().family }

// FontSize returns the configured base text size in DIPs.
func FontSize() float32 { return snapshot().prefs.FontSize }

// TextScale combines the OS text scale with the configured base relative to
// HeroUI's 16 DIP root.
func TextScale(c *ui.Context) float32 {
	return c.Preferences().TextScale * FontSize() / 16
}

// Scale converts an unscaled HeroUI text, line-box or icon metric to DIPs.
func Scale(c *ui.Context, dip float32) float32 { return dip * TextScale(c) }

// ReduceMotion reflects the mandatory OS accessibility request.
func ReduceMotion(c *ui.Context) bool { return c.Preferences().ReduceMotion }

// Apply installs HeroUI for the current native appearance. Call once at the
// view root before constructing children. The installed theme is composed
// from the published preferences, so an injected Palette is preserved and
// also reaches native widgets.
func Apply(c *ui.Context) *ui.Theme { return ApplyMode(c, c.Theme().Dark) }

// ApplyMode installs HeroUI for an explicit mode. Previews and tests use it.
func ApplyMode(c *ui.Context, dark bool) *ui.Theme {
	t := nativeTheme(ColorsForMode(dark), dark)
	t.Font, t.FontSize = FontFamily(c), FontSize()
	c.SetTheme(&t)
	return &t
}

// Colors returns the live tokens for the context's current mode.
func Colors(c *ui.Context) Palette { return ColorsForMode(c.Theme().Dark) }

// ColorsForMode returns tokens without a context, for static thumbnails.
func ColorsForMode(dark bool) Palette {
	if dark {
		return snapshot().dark
	}
	return snapshot().light
}

func nativeTheme(p Palette, dark bool) ui.Theme {
	surface := p.Default
	// A host sheet is also the native control face: MyGo paints buttons and
	// inputs in Surface, and a host palette uses one sheet for both.
	if patch := snapshot().prefs.Palette.patch(dark); patch != nil && patch.Surface != nil {
		surface = p.Surface
	}
	return ui.Theme{
		Dark: dark, Background: p.Background, Surface: surface,
		SurfaceHover: p.DefaultHover, SurfacePressed: p.DefaultHover,
		Border: p.Border, Text: p.Foreground, TextMuted: p.Muted,
		Accent: p.Accent, AccentHover: p.AccentHover, AccentPressed: p.AccentHover,
		AccentText: p.AccentForeground, Danger: p.Danger,
		Warning: p.Warning, Success: p.Success,
		Selection: p.Selection, Focus: p.Focus,
		Inverse: p.Foreground, InverseText: p.Background,
		Scrollbar: p.Scrollbar, ScrollbarWidth: Space6, Radius: RadiusBase(),
		Spacing: SpaceUnit, FontSize: FontSize(),
	}
}
