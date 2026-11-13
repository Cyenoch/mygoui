// Package theme owns the design tokens: semantic colors, geometry, type,
// elevation, motion and effects.
//
// A Theme names a whole design language as one value; Default is the
// built-in one and Use publishes another. Tokens are published as immutable
// snapshots: components read live values with Colors(c), Motion() and the
// other accessors, and never cache them keyed by window. Host customization
// goes through Use and Configure before Apply.
package theme

import (
	"reflect"
	"sync/atomic"

	"github.com/egoist/mygo/ui"
)

// Preferences is the process-wide customization surface. Zero fields
// inherit the published theme: by default the system font, 16 DIP base text,
// a 10 DIP radius base and the neutral primary action color.
//
// Semantic colors resolve by precedence, highest first: Preferences.Palette,
// Preferences.Accent, the published palette. Apply installs the theme
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
	// Radius is the base radius; the radius scale is derived from it.
	Radius float32
	// Accent replaces the primary action color when opaque.
	Accent ui.Color
	// Palette overrides individual semantic tokens per appearance. The zero
	// value keeps the default palette. Configure replaces the published
	// preferences as a whole, so a later call that changes only typography
	// must publish the same Palette again; UpdatePreferences edits in place.
	Palette PaletteOverrides
}

type state struct {
	// prefs are as the host published them: a zero FontSize or Radius
	// inherits the theme.
	prefs Preferences
	theme Theme

	// motion is theme.Motion as components read it.
	motion MotionTokens

	family, heading, mono string
	size, radius          float32
	radii                 Radii
	light, dark           Palette
}

var current atomic.Pointer[state]

func init() { publish(Default(), Preferences{}) }

// Configure publishes host preferences over the published theme and
// recomputes both palettes. Call on the UI thread, then invalidate affected
// windows. Publishing the same preferences again is a no-op, so a view root
// may call it on every build.
//
// A FontSize or Radius equal to the published theme's value means inherit,
// as zero does: Configure(PreferencesSnapshot()) pins neither, and a theme
// published later still sets them.
func Configure(p Preferences) {
	if p.Accent.A == 0 {
		p.Accent = ui.Color{}
	}
	// Compare before resolving: Configure is on the build path, and the
	// published value may be unchanged.
	existing := snapshot()
	if p.FontSize == existing.theme.Typography.Size {
		p.FontSize = 0
	}
	if p.Radius == existing.theme.Radius {
		p.Radius = 0
	}
	if reflect.DeepEqual(existing.prefs, p) {
		return
	}
	p.Palette.Light = clonePatch(p.Palette.Light)
	p.Palette.Dark = clonePatch(p.Palette.Dark)
	publish(existing.theme, p)
}

// publish resolves p over t and stores the result. It owns both values.
func publish(t Theme, p Preferences) {
	s := &state{prefs: p, theme: t, motion: t.Motion.resolved(), size: t.Typography.Size, radius: t.Radius}
	if p.FontSize != 0 {
		s.size = p.FontSize
	}
	if p.Radius != 0 {
		s.radius = p.Radius
	}
	s.radii = radii(s.radius)
	s.family = fontStack(p.FontFamily, t.Typography.Family, Sans)
	s.heading = s.family
	if p.HeadingFontFamily != "" || t.Typography.HeadingFamily != "" {
		s.heading = fontStack(p.HeadingFontFamily, t.Typography.HeadingFamily, Sans)
	}
	s.mono = fontStack(p.MonoFontFamily, t.Typography.MonoFamily, Mono)
	s.light = resolvePalette(t.Light, false, p, t.Effects.HoverOpacity)
	s.dark = resolvePalette(t.Dark, true, p, t.Effects.HoverOpacity)
	current.Store(s)
}

// fontStack is the preferred family, else the theme's, before the native
// fallbacks.
func fontStack(preferred, themed, fallback string) string {
	if preferred == "" {
		preferred = themed
	}
	if preferred == "" {
		return fallback
	}
	return preferred + ", " + fallback
}

func snapshot() *state { return current.Load() }

// UpdatePreferences edits the published preferences and republishes them, so
// a host that changes only typography keeps its injected palette. A font
// size or radius the edit leaves at the theme's value keeps following the
// theme.
func UpdatePreferences(edit func(*Preferences)) {
	p := PreferencesSnapshot()
	edit(&p)
	Configure(p)
}

// PreferencesSnapshot returns the published preferences, with the font size
// and radius in effect.
func PreferencesSnapshot() Preferences {
	s := snapshot()
	p := s.prefs
	p.FontSize, p.Radius = s.size, s.radius
	p.Palette.Light = clonePatch(p.Palette.Light)
	p.Palette.Dark = clonePatch(p.Palette.Dark)
	return p
}

// FontFamily returns the configured family followed by native fallbacks.
func FontFamily(*ui.Context) string { return snapshot().family }

// HeadingFontFamily resolves the heading family separately from body text;
// without one, headings use the body family.
func HeadingFontFamily(*ui.Context) string { return snapshot().heading }

// MonoFontFamily resolves the caller's code face with native fallbacks.
func MonoFontFamily() string { return snapshot().mono }

// FontSize returns the base text size in DIPs.
func FontSize() float32 { return snapshot().size }

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
	hover := opacity(p.Primary, Effects().HoverOpacity)
	return ui.Theme{
		Dark: dark, Background: p.Background, Surface: surface,
		SurfaceHover: p.Accent.Over(surface), SurfacePressed: p.Accent.Over(surface),
		Border: p.Border, Text: p.Foreground, TextMuted: p.Muted,
		Accent: p.Primary, AccentHover: hover, AccentPressed: hover,
		AccentText: p.PrimaryForeground, Danger: p.Destructive,
		Warning: p.Warning, Success: p.Success,
		Selection: p.Selection, Focus: p.Focus,
		Inverse: p.Foreground, InverseText: p.Background,
		Scrollbar: p.Scrollbar, ScrollbarWidth: Space6, Radius: RadiusBase(),
		Spacing: Space(1), FontSize: FontSize(),
	}
}
