// Package theme owns coss ui's semantic tokens, geometry and typography.
//
// Tokens are derived once from upstream and published as immutable snapshots:
// components read live values with Colors(c) and never cache a palette keyed by
// window. Host customization goes through Configure before Apply.
package theme

import (
	"reflect"
	"sync/atomic"

	"github.com/egoist/mygo/ui"
)

// Preferences is the process-wide coss ui customization surface. Zero fields
// select native defaults: the system font, 16 DIP base text, a 10 DIP radius
// base and the neutral primary action color.
//
// Semantic colors resolve by precedence, highest first: Preferences.Palette,
// Preferences.Accent, coss ui's published palette. Apply installs the theme
// composed from that resolution, so an injected palette survives a view-root
// Apply and reaches native widgets (selections, scrollbars, text cursors)
// too. Library paint reads Colors, never the context theme, so a host that
// styles native widgets through its own scoped theme either installs that
// theme after Apply or skips Apply.
type Preferences struct {
	FontFamily        string
	HeadingFontFamily string
	MonoFontFamily    string
	FontSize          float32
	// Radius is the base radius; the coss ui scale is derived from it.
	Radius float32
	// Accent replaces the primary action color when opaque.
	Accent ui.Color
	// Palette overrides individual semantic tokens per appearance. The zero
	// value keeps the coss ui defaults. Configure replaces the published
	// preferences as a whole, so a later call that changes only typography
	// must publish the same Palette again; UpdatePreferences edits in place.
	Palette PaletteOverrides
}

type state struct {
	prefs  Preferences
	family string
	radii  Radii
	light  Palette
	dark   Palette
}

var current atomic.Pointer[state]

func init() { Configure(Preferences{}) }

// Configure publishes host preferences and recomputes both palettes. Call on
// the UI thread, then invalidate affected windows. Publishing the same
// preferences again is a no-op, so a view root may call it on every build.
func Configure(p Preferences) {
	if p.FontSize == 0 {
		p.FontSize = 16
	}
	if p.Radius == 0 {
		p.Radius = 10
	}
	if p.Accent.A == 0 {
		p.Accent = ui.Color{}
	}
	// Compare before deriving the family string: Configure is on the build
	// path, and the published value may be unchanged.
	if existing := current.Load(); existing != nil && reflect.DeepEqual(existing.prefs, p) {
		return
	}
	family := Sans
	if p.FontFamily != "" {
		family = p.FontFamily + ", " + Sans
	}
	p.Palette.Light = clonePatch(p.Palette.Light)
	p.Palette.Dark = clonePatch(p.Palette.Dark)
	current.Store(&state{
		prefs:  p,
		family: family,
		radii:  radii(p.Radius),
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
func PreferencesSnapshot() Preferences {
	p := snapshot().prefs
	p.Palette.Light = clonePatch(p.Palette.Light)
	p.Palette.Dark = clonePatch(p.Palette.Dark)
	return p
}

// FontFamily returns the configured family followed by native fallbacks.
func FontFamily(*ui.Context) string { return snapshot().family }

// HeadingFontFamily resolves coss's font-heading separately from body text.
func HeadingFontFamily(c *ui.Context) string {
	if family := snapshot().prefs.HeadingFontFamily; family != "" {
		return family + ", " + Sans
	}
	return FontFamily(c)
}

// MonoFontFamily resolves the caller's code face with native fallbacks.
func MonoFontFamily() string {
	if family := snapshot().prefs.MonoFontFamily; family != "" {
		return family + ", " + Mono
	}
	return Mono
}

// FontSize returns the configured base text size in DIPs.
func FontSize() float32 { return snapshot().prefs.FontSize }

// TextScale combines the OS text scale with the configured base relative to
// coss ui's 16 DIP root.
func TextScale(c *ui.Context) float32 {
	return c.Preferences().TextScale * FontSize() / 16
}

// Scale converts an unscaled coss ui text, line-box or icon metric to DIPs.
func Scale(c *ui.Context, dip float32) float32 { return dip * TextScale(c) }

// ReduceMotion reflects the mandatory OS accessibility request.
func ReduceMotion(c *ui.Context) bool { return c.Preferences().ReduceMotion }

// Apply installs coss ui for the current native appearance. Call once at the
// view root before constructing children. The installed theme is composed
// from the published preferences, so an injected Palette is preserved and
// also reaches native widgets.
func Apply(c *ui.Context) *ui.Theme { return ApplyMode(c, c.Theme().Dark) }

// ApplyMode installs coss ui for an explicit mode. Previews and tests use it.
// Reinstalling the theme already in place is a no-op, so a view root may call
// it on every build without allocating a theme or restyling the root.
func ApplyMode(c *ui.Context, dark bool) *ui.Theme {
	want := nativeTheme(ColorsForMode(dark), dark)
	want.Font, want.FontSize = FontFamily(c), FontSize()
	if installed := c.Theme(); installed != nil && *installed == want {
		return installed
	}
	installed := new(ui.Theme)
	*installed = want
	c.SetTheme(installed)
	return installed
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
	surface := p.Popover
	// A host sheet is also the native control face: MyGo paints buttons and
	// inputs in Surface, and a host palette uses one sheet for both.
	if patch := snapshot().prefs.Palette.patch(dark); patch != nil && patch.Surface != nil {
		surface = p.Surface
	}
	return ui.Theme{
		Dark: dark, Background: p.Background, Surface: surface,
		SurfaceHover: p.Accent.Over(surface), SurfacePressed: p.Accent.Over(surface),
		Border: p.Border, Text: p.Foreground, TextMuted: p.Muted,
		Accent: p.Primary, AccentHover: opacity(p.Primary, 0.9), AccentPressed: opacity(p.Primary, 0.9),
		AccentText: p.PrimaryForeground, Danger: p.Destructive,
		Warning: p.Warning, Success: p.Success,
		Selection: p.Selection, Focus: p.Focus,
		Inverse: p.Foreground, InverseText: p.Background,
		Scrollbar: p.Scrollbar, ScrollbarWidth: Space6, Radius: RadiusBase(),
		Spacing: SpaceUnit, FontSize: FontSize(),
	}
}
