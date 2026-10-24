package theme_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// palettePatch returns one fully populated patch, so every injected role can
// be traced from Preferences to the resolved palette and the native theme.
func palettePatch() (theme.PalettePatch, map[string]ui.Color) {
	want := map[string]ui.Color{
		"canvas":     ui.RGB(10, 20, 30),
		"sheet":      ui.RGB(40, 50, 60),
		"ink":        ui.RGB(230, 240, 250),
		"muted":      ui.RGB(120, 130, 140),
		"accent":     ui.RGB(70, 95, 208),
		"accentink":  ui.RGB(255, 255, 255),
		"danger":     ui.RGB(184, 50, 50),
		"dangerink":  ui.RGB(255, 250, 248),
		"success":    ui.RGB(17, 107, 75),
		"successink": ui.RGB(240, 255, 248),
		"warning":    ui.RGB(138, 90, 18),
		"warningink": ui.RGB(255, 250, 240),
		"info":       ui.RGB(60, 120, 220),
		"infoink":    ui.RGB(245, 250, 255),
		"hover":      ui.RGB(30, 40, 55),
		"selected":   ui.RGBA(200, 30, 60, 0.5),
		"overlay":    ui.RGB(50, 60, 70),
		"backdrop":   ui.RGBA(5, 6, 7, 0.6),
	}
	patch := theme.PalettePatch{
		Background: new(want["canvas"]), Surface: new(want["sheet"]),
		Foreground: new(want["ink"]), Muted: new(want["muted"]),
		Border: new(ui.Transparent),
		Accent: new(want["accent"]), AccentForeground: new(want["accentink"]),
		Danger: new(want["danger"]), DangerForeground: new(want["dangerink"]),
		Success: new(want["success"]), SuccessForeground: new(want["successink"]),
		Warning: new(want["warning"]), WarningForeground: new(want["warningink"]),
		Info: new(want["info"]), InfoForeground: new(want["infoink"]),
		Hover: new(want["hover"]), Selected: new(want["selected"]),
		Overlay: new(want["overlay"]), Backdrop: new(want["backdrop"]),
	}
	return patch, want
}

func TestPalettePatchOverridesAndDerives(t *testing.T) {
	defaultLight, defaultDark := theme.ColorsForMode(false), theme.ColorsForMode(true)
	defer theme.Configure(theme.Preferences{})
	patch, want := palettePatch()
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &patch}})
	got := theme.ColorsForMode(false)

	for _, tc := range []struct {
		role string
		got  ui.Color
		want ui.Color
	}{
		{"canvas", got.Background, want["canvas"]},
		{"surface", got.Surface, want["sheet"]},
		{"ink", got.Foreground, want["ink"]},
		{"muted", got.Muted, want["muted"]},
		{"border", got.Border, ui.Transparent},
		{"accent", got.Accent, want["accent"]},
		{"accent foreground", got.AccentForeground, want["accentink"]},
		{"danger", got.Danger, want["danger"]},
		{"danger foreground", got.DangerForeground, want["dangerink"]},
		{"success", got.Success, want["success"]},
		{"warning", got.Warning, want["warning"]},
		{"info", got.Info, want["info"]},
		{"hover", got.SurfaceHover, want["hover"]},
		{"neutral fill", got.Default, want["hover"]},
		{"selected", got.Selection, want["selected"]},
		{"overlay", got.Overlay, want["overlay"]},
		{"backdrop", got.Backdrop, want["backdrop"]},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.role, tc.got, tc.want)
		}
	}

	// Roles that pair with the injected ink follow it.
	for _, tc := range []struct {
		role string
		got  ui.Color
	}{
		{"surface foreground", got.SurfaceForeground},
		{"surface secondary foreground", got.SurfaceSecondaryForeground},
		{"surface tertiary foreground", got.SurfaceTertiaryForeground},
		{"overlay foreground", got.OverlayForeground},
		{"segment foreground", got.SegmentForeground},
		{"default foreground", got.DefaultForeground},
		{"field foreground", got.FieldForeground},
	} {
		if tc.got != want["ink"] {
			t.Errorf("%s = %v, want the injected ink %v", tc.role, tc.got, want["ink"])
		}
	}
	if got.FieldBackground != want["sheet"] {
		t.Errorf("field background = %v, want the injected surface", got.FieldBackground)
	}
	if got.FieldPlaceholder != want["muted"] {
		t.Errorf("field placeholder = %v, want the injected muted ink", got.FieldPlaceholder)
	}
	if got.Focus != want["accent"] {
		t.Errorf("focus = %v, want the injected accent", got.Focus)
	}

	// Derived levels are rebuilt from the injected base, not HeroUI's.
	if got.AccentHover == got.Accent {
		t.Error("accent hover was not derived from the injected accent")
	}
	if got.DefaultHover == got.Default {
		t.Error("neutral hover was not derived from the injected hover fill")
	}
	if got.AccentSoft.A == 0 || got.AccentSoft.A == 255 {
		t.Errorf("accent soft alpha = %d, want translucent", got.AccentSoft.A)
	}
	if got.InfoSoft.A == 0 || got.InfoSoft.A == 255 {
		t.Errorf("info soft alpha = %d, want translucent", got.InfoSoft.A)
	}

	// Tokens the patch does not name keep the HeroUI defaults.
	if got.SurfaceSecondary != defaultLight.SurfaceSecondary {
		t.Errorf("unset surface secondary = %v, want the HeroUI default %v", got.SurfaceSecondary, defaultLight.SurfaceSecondary)
	}
	if got.FieldBorder != defaultLight.FieldBorder {
		t.Errorf("unset field border = %v, want the HeroUI default %v", got.FieldBorder, defaultLight.FieldBorder)
	}
	if dark := theme.ColorsForMode(true); dark.Background != defaultDark.Background {
		t.Errorf("light patch changed the dark canvas: %v, want %v", dark.Background, defaultDark.Background)
	}
}

func TestInfoDefaultsToAccentAndPatchSeparatesIt(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	accent := theme.ColorsForMode(false).Accent
	if got := theme.ColorsForMode(false); got.Info != accent || got.InfoForeground != got.AccentForeground {
		t.Fatalf("default info = %v/%v, want the accent family", got.Info, got.InfoForeground)
	}
	info := ui.RGB(60, 120, 220)
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{
		Light: &theme.PalettePatch{Info: new(info)},
	}})
	got := theme.ColorsForMode(false)
	if got.Info != info {
		t.Errorf("info = %v, want %v", got.Info, info)
	}
	if got.Accent != accent {
		t.Errorf("info patch changed the accent: %v, want %v", got.Accent, accent)
	}
	if got.InfoSoft.R == 0 || got.InfoSoft.A == 255 {
		t.Errorf("info soft = %v, want a translucent tint of the injected info", got.InfoSoft)
	}
}

func TestApplyInstallsInjectedPalette(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	patch, want := palettePatch()
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &patch}})

	var first, second ui.Theme
	var frames int
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		if frames == 0 {
			first = *c.Theme()
		} else {
			second = *c.Theme()
		}
		frames++
	}, 120, 80)
	tt.Frame()
	tt.Frame()

	if first.Background != want["canvas"] {
		t.Errorf("native background = %v, want the injected canvas", first.Background)
	}
	// An injected surface is also the native control face.
	if first.Surface != want["sheet"] {
		t.Errorf("native surface = %v, want the injected sheet", first.Surface)
	}
	if first.Border != ui.Transparent {
		t.Errorf("native border = %v, want transparent", first.Border)
	}
	if first.Text != want["ink"] || first.TextMuted != want["muted"] {
		t.Errorf("native ink = %v/%v, want the injected ink/muted", first.Text, first.TextMuted)
	}
	if first.Accent != want["accent"] || first.AccentText != want["accentink"] {
		t.Errorf("native accent = %v/%v, want the injected pair", first.Accent, first.AccentText)
	}
	if first.Selection != want["selected"] {
		t.Errorf("native selection = %v, want the injected fill", first.Selection)
	}
	if first.Danger != want["danger"] || first.Success != want["success"] || first.Warning != want["warning"] {
		t.Error("native status colors do not follow the injected palette")
	}
	if first != second {
		t.Error("a second Apply changed the installed theme")
	}
}

func TestUpdatePreferencesKeepsInjectedPalette(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	patch, want := palettePatch()
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &patch}})
	theme.UpdatePreferences(func(p *theme.Preferences) { p.FontFamily = "Test Sans" })
	if got := theme.PreferencesSnapshot().FontFamily; got != "Test Sans" {
		t.Fatalf("font family = %q, want the edited value", got)
	}
	if got := theme.ColorsForMode(false).Background; got != want["canvas"] {
		t.Errorf("palette was lost by the typography edit: canvas = %v", got)
	}
	// Configure replaces the whole value, including the palette.
	theme.Configure(theme.Preferences{FontFamily: "Test Sans"})
	if got := theme.ColorsForMode(false).Background; got == want["canvas"] {
		t.Error("Configure kept a palette it was not given")
	}
}
