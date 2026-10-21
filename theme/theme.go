// Package theme owns reusable native semantic colors, geometry, and typography.
// Apply installs Neutral; ApplyPreset selects an explicit palette.
// Components consume Colors(c), never appearance-only snapshots.
package theme

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/egoist/mygo/ui"
)

// Preferences contains host policy, separate from read-only OS preferences.
type Preferences struct {
	FontFamily             string
	FontSize               float32
	PanelAnimationDuration time.Duration
}

type preferenceState struct {
	Preferences
	family string
}

// Immutable copies are published on edits; frame reads allocate nothing.
var preferences atomic.Pointer[preferenceState]

// SetPreferences publishes host policy. Call on the UI thread, then invalidate
// affected windows. Zero FontSize selects the source default of 16; zero panel
// duration disables panel transitions, not ongoing activity indicators.
func SetPreferences(p Preferences) {
	if p.FontSize == 0 {
		p.FontSize = 16
	}
	if current := preferences.Load(); current != nil && current.Preferences == p {
		return
	}
	family := Sans
	if custom := strings.TrimSpace(p.FontFamily); custom != "" {
		family = custom + ", " + Sans
	}
	preferences.Store(&preferenceState{Preferences: p, family: family})
}

func preferenceSnapshot() Preferences {
	if current := preferences.Load(); current != nil {
		return current.Preferences
	}
	return Preferences{FontSize: 16}
}

// TextScale combines OS scale with the source interface-size baseline of 16.
func TextScale(c *ui.Context) float32 {
	return c.Preferences().TextScale * preferenceSnapshot().FontSize / 16
}

// FontFamily returns the installed interface family followed by native fallbacks.
func FontFamily(c *ui.Context) string {
	if current := preferences.Load(); current != nil {
		return current.family
	}
	return Sans
}

// ReduceMotion reflects the mandatory OS accessibility request.
func ReduceMotion(c *ui.Context) bool { return c.Preferences().ReduceMotion }

// PanelAnimationDuration returns zero when panels should change immediately.
// Progress indicators retain their own cadence; this preference is panel-only.
func PanelAnimationDuration(c *ui.Context) time.Duration {
	if ReduceMotion(c) {
		return 0
	}
	return preferenceSnapshot().PanelAnimationDuration
}

// Scale converts an unscaled text, line-box, or inline-icon metric to DIPs.
// Do not use it for pane widths, padding, or fixed chrome target dimensions.
func Scale(c *ui.Context, dip float32) float32 { return dip * TextScale(c) }

// Apply installs the neutral library preset for the current native appearance.
// Call once at the view root before constructing child mygoui.
func Apply(c *ui.Context) *ui.Theme { return ApplyPreset(c, Neutral) }

// ApplyPreset installs a preset while retaining live OS accessibility inputs.
// The returned native theme is borrowed and immutable; copy it before making
// a scoped override. Cached values retain no Context, Element, or window.
func ApplyPreset(c *ui.Context, preset Preset) *ui.Theme {
	return applyPreset(c, preset, c.Preferences().HighContrast)
}

// ApplyHighContrast previews the shared accessibility policy in galleries.
// It only strengthens contrast; it never changes the read-only OS preferences.
func ApplyHighContrast(c *ui.Context, preset Preset) *ui.Theme {
	return applyPreset(c, preset, true)
}

func applyPreset(c *ui.Context, preset Preset, highContrast bool) *ui.Theme {
	dark := c.Theme().Dark
	prefs := c.Preferences()
	family, size := FontFamily(c), Scale(c, 16)
	key := rootThemeKey{dark: dark, preset: preset, accent: prefs.Accent,
		highContrast: highContrast, family: family, size: size}
	rootThemes.Lock()
	for _, entry := range rootThemes.entries {
		if entry != nil && entry.key == key {
			rootThemes.Unlock()
			c.SetTheme(&entry.value.theme)
			return &entry.value.theme
		}
	}
	p := sourceColors(dark, preset)
	if prefs.Accent.A > 0 {
		p.Primary = prefs.Accent
		p.Primary.A = 255
		p.Info = p.Primary
		p.SidebarPrimary = p.Primary
	}
	finishColors(&p, dark, highContrast)
	t := nativeTheme(p, dark)
	t.Font, t.FontSize = family, size
	value := rememberColors(t, p, highContrast)
	entry := &rootThemeEntry{key: key, value: value}
	rootThemes.entries[rootThemes.next] = entry
	rootThemes.next = (rootThemes.next + 1) % uint8(len(rootThemes.entries))
	rootThemes.Unlock()
	c.SetTheme(&value.theme)
	return &value.theme
}

type rootThemeKey struct {
	dark, highContrast bool
	preset             Preset
	accent             ui.Color
	family             string
	size               float32
}

type rootThemeEntry struct {
	key   rootThemeKey
	value *resolvedColors
}

var rootThemes struct {
	sync.Mutex
	entries [16]*rootThemeEntry
	next    uint8
}
