package theme_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCossPaletteDefaults(t *testing.T) {
	theme.Configure(theme.Preferences{})
	defer theme.Configure(theme.Preferences{})
	for _, mode := range []struct {
		name string
		dark bool
	}{
		{name: "light"},
		{name: "dark", dark: true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			got := theme.ColorsForMode(mode.dark)
			primary, primaryInk := ui.Oklch(0.269, 0, 0), ui.Oklch(0.985, 0, 0)
			accent := ui.RGBA(0, 0, 0, 0.04)
			border, input := ui.RGBA(0, 0, 0, 0.08), ui.RGBA(0, 0, 0, 0.10)
			destructiveInk := ui.Oklch(0.505, 0.213, 27.518)
			infoInk := ui.Oklch(0.488, 0.243, 264.376)
			successInk := ui.Oklch(0.508, 0.118, 165.612)
			warningInk := ui.Oklch(0.555, 0.163, 48.998)
			if mode.dark {
				primary, primaryInk = ui.Oklch(0.97, 0, 0), primary
				accent = ui.RGBA(255, 255, 255, 0.04)
				border, input = ui.RGBA(255, 255, 255, 0.06), ui.RGBA(255, 255, 255, 0.08)
				destructiveInk = ui.Oklch(0.704, 0.191, 22.216)
				infoInk = ui.Oklch(0.707, 0.165, 254.624)
				successInk = ui.Oklch(0.765, 0.177, 163.223)
				warningInk = ui.Oklch(0.828, 0.189, 84.429)
			}
			for _, token := range []struct {
				name      string
				got, want ui.Color
			}{
				{name: "primary", got: got.Primary, want: primary},
				{name: "primary foreground", got: got.PrimaryForeground, want: primaryInk},
				{name: "accent", got: got.Accent, want: accent},
				{name: "accent foreground", got: got.AccentForeground, want: got.Foreground},
				{name: "secondary", got: got.Secondary, want: accent},
				{name: "muted background", got: got.MutedBackground, want: accent},
				{name: "border", got: got.Border, want: border},
				{name: "input", got: got.Input, want: input},
				{name: "info", got: got.Info, want: ui.Oklch(0.623, 0.214, 259.815)},
				{name: "info foreground", got: got.InfoForeground, want: infoInk},
				{name: "success foreground", got: got.SuccessForeground, want: successInk},
				{name: "warning foreground", got: got.WarningForeground, want: warningInk},
				{name: "destructive foreground", got: got.DestructiveForeground, want: destructiveInk},
			} {
				t.Run(token.name, func(t *testing.T) {
					if token.got != token.want {
						t.Errorf("color = %v, want %v", token.got, token.want)
					}
				})
			}
		})
	}
}

func TestAccentPreferenceCustomizesPrimaryOnly(t *testing.T) {
	theme.Configure(theme.Preferences{})
	defer theme.Configure(theme.Preferences{})
	light, dark := theme.ColorsForMode(false), theme.ColorsForMode(true)
	custom := ui.RGB(220, 20, 60)
	theme.Configure(theme.Preferences{Accent: custom})
	for _, mode := range []struct {
		name string
		dark bool
		base theme.Palette
	}{
		{name: "light", base: light},
		{name: "dark", dark: true, base: dark},
	} {
		t.Run(mode.name, func(t *testing.T) {
			got := theme.ColorsForMode(mode.dark)
			if got.Primary != custom || got.PrimaryForeground != ui.RGB(255, 255, 255) {
				t.Errorf("primary = %v/%v, want custom color and contrasting white", got.Primary, got.PrimaryForeground)
			}
			if got.Accent != mode.base.Accent || got.Info != mode.base.Info {
				t.Error("primary preference changed the independent accent or info token")
			}
		})
	}
}

func TestFontAndTextPreferences(t *testing.T) {
	theme.Configure(theme.Preferences{})
	defer theme.Configure(theme.Preferences{})
	if theme.FontSize() != 16 || theme.FontFamily(nil) != theme.Sans {
		t.Fatal("default text is not 16 DIP with native sans fallbacks")
	}
	theme.Configure(theme.Preferences{FontFamily: "Body Face", FontSize: 20})
	if got := theme.HeadingFontFamily(nil); got != theme.FontFamily(nil) {
		t.Fatalf("default heading family = %q, want the configured body face", got)
	}
	theme.UpdatePreferences(func(p *theme.Preferences) {
		p.HeadingFontFamily = "Heading Face"
		p.MonoFontFamily = "Code Face"
	})
	for _, font := range []struct{ name, got, want string }{
		{name: "body", got: theme.FontFamily(nil), want: "Body Face, " + theme.Sans},
		{name: "heading", got: theme.HeadingFontFamily(nil), want: "Heading Face, " + theme.Sans},
		{name: "mono", got: theme.MonoFontFamily(), want: "Code Face, " + theme.Mono},
	} {
		t.Run(font.name, func(t *testing.T) {
			if font.got != font.want {
				t.Errorf("family = %q, want %q", font.got, font.want)
			}
		})
	}

	var scale, metric float32
	var installed ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		installed = *theme.Apply(c)
		scale, metric = theme.TextScale(c), theme.Scale(c, 14)
	}, 120, 80)
	tt.SetPreferences(ui.Preferences{TextScale: 1.5})
	tt.Frame()
	if scale != 1.875 || metric != 26.25 {
		t.Errorf("text scale/14 DIP metric = %v/%v, want 1.875/26.25", scale, metric)
	}
	if installed.Font != "Body Face, "+theme.Sans || installed.FontSize != 20 {
		t.Errorf("native font = %q/%v, want the configured body font/20 DIP", installed.Font, installed.FontSize)
	}
}
