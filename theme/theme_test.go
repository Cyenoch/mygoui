package theme_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestPaletteModesDiffer(t *testing.T) {
	light, dark := theme.ColorsForMode(false), theme.ColorsForMode(true)
	if light.Background.R <= dark.Background.R {
		t.Fatal("light mode is not lighter than dark mode")
	}
	if light.Surface == dark.Surface {
		t.Fatal("light and dark surfaces are identical")
	}
	if light.Accent.R != dark.Accent.R || light.Accent.G != dark.Accent.G {
		t.Fatal("accent is shared between modes, as upstream specifies")
	}
}

func TestAccentCustomizationDerivesSoftAndHover(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	custom := ui.RGB(220, 20, 60)
	theme.Configure(theme.Preferences{Accent: custom})
	got := theme.ColorsForMode(false)
	if got.Accent.R != custom.R || got.Accent.G != custom.G || got.Accent.B != custom.B {
		t.Fatalf("accent = %v, want %v", got.Accent, custom)
	}
	if got.AccentSoft.A == 0 || got.AccentSoft.A == 255 {
		t.Fatalf("soft accent alpha = %d, want translucent", got.AccentSoft.A)
	}
	if got.AccentHover == got.Accent {
		t.Fatal("hover accent was not derived")
	}
}

func TestLightStatusSoftForegroundStaysNearStatus(t *testing.T) {
	light := theme.ColorsForMode(false)
	cases := []struct {
		name             string
		status, soft, fg ui.Color
	}{
		{"warning", light.Warning, light.WarningSoftForeground, light.Foreground},
		{"success", light.Success, light.SuccessSoftForeground, light.Foreground},
		{"danger", light.Danger, light.DangerSoftForeground, light.Foreground},
	}
	for _, tc := range cases {
		toStatus, toForeground := channelDistance(tc.soft, tc.status), channelDistance(tc.soft, tc.fg)
		if toStatus >= toForeground {
			t.Errorf("%s soft foreground = %v: distance to status %d >= distance to foreground %d",
				tc.name, tc.soft, toStatus, toForeground)
		}
	}
}

func channelDistance(a, b ui.Color) int {
	d := 0
	for _, c := range [][2]uint8{{a.R, b.R}, {a.G, b.G}, {a.B, b.B}} {
		if c[0] > c[1] {
			d += int(c[0] - c[1])
		} else {
			d += int(c[1] - c[0])
		}
	}
	return d
}

func TestRadiusAndTextPreferences(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	theme.Configure(theme.Preferences{Radius: 16, FontSize: 20})
	if got := theme.RadiusScale().XXXL; got != 48 {
		t.Fatalf("radius 3xl = %v, want 48", got)
	}
	if got := theme.FontSize(); got != 20 {
		t.Fatalf("font size = %v, want 20", got)
	}
}
