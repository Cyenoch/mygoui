package theme

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

// Independent WCAG calculation: assert the painted pair rather than token RGB.
func paintedContrast(ink, surface, base ui.Color) float64 {
	background := surface.Over(base)
	foreground := ink.Over(background)
	luminance := func(c ui.Color) float64 {
		channel := func(v uint8) float64 {
			x := float64(v) / 255
			if x <= 0.04045 {
				return x / 12.92
			}
			return math.Pow((x+0.055)/1.055, 2.4)
		}
		return .2126*channel(c.R) + .7152*channel(c.G) + .0722*channel(c.B)
	}
	a, b := luminance(foreground), luminance(background)
	return (max(a, b) + .05) / (min(a, b) + .05)
}

func assertReadablePairs(t *testing.T, p ColorTokens) {
	t.Helper()
	for _, pair := range []struct {
		name         string
		ink, surface ui.Color
	}{
		{"body", p.Foreground, p.Background}, {"card", p.CardForeground, p.Card},
		{"popover", p.PopoverForeground, p.Popover}, {"primary", p.PrimaryForeground, p.Primary},
		{"primary hover", p.PrimaryForeground, p.PrimaryHover}, {"primary pressed", p.PrimaryForeground, p.PrimaryPressed},
		{"secondary", p.SecondaryForeground, p.Secondary}, {"secondary hover", p.SecondaryForeground, p.SecondaryHover},
		{"secondary pressed", p.SecondaryForeground, p.SecondaryPressed}, {"muted", p.MutedForeground, p.Muted},
		{"accent", p.AccentForeground, p.Accent}, {"accent hover", p.AccentForeground, p.AccentHover},
		{"accent pressed", p.AccentForeground, p.AccentPressed}, {"destructive", p.DestructiveForeground, p.Destructive},
		{"destructive hover", p.DestructiveForeground, p.DestructiveHover}, {"destructive pressed", p.DestructiveForeground, p.DestructivePressed},
		{"success", p.SuccessForeground, p.Success}, {"warning", p.WarningForeground, p.Warning}, {"info", p.InfoForeground, p.Info},
		{"sidebar", p.SidebarForeground, p.Sidebar}, {"sidebar primary", p.SidebarPrimaryForeground, p.SidebarPrimary},
		{"sidebar accent", p.SidebarAccentForeground, p.SidebarAccent}, {"code", p.CodeForeground, p.CodeBackground},
		{"selection", p.SelectionForeground, p.Selection},
		{"warning notice", p.Warning, p.WarningSurface}, {"error notice", p.Destructive, p.ErrorSurface},
	} {
		if ratio := paintedContrast(pair.ink, pair.surface, p.Background); ratio < 4.5 {
			t.Errorf("%s painted contrast %.3f:1, want >=4.5:1", pair.name, ratio)
		}
	}
	for _, role := range []struct {
		name string
		ink  ui.Color
	}{
		{"destructive ink", p.Destructive}, {"success ink", p.Success}, {"warning ink", p.Warning}, {"info ink", p.Info},
	} {
		for _, surface := range []ui.Color{p.Background, p.Card} {
			if ratio := paintedContrast(role.ink, surface, p.Background); ratio < 4.5 {
				t.Errorf("%s contrast %.3f:1, want >=4.5:1", role.name, ratio)
			}
		}
	}
	for _, notice := range []struct {
		name         string
		ink, surface ui.Color
	}{
		{"warning card notice", p.Warning, p.WarningSurface}, {"error card notice", p.Destructive, p.ErrorSurface},
	} {
		if ratio := paintedContrast(notice.ink, notice.surface, p.Card); ratio < 4.5 {
			t.Errorf("%s contrast %.3f:1, want >=4.5:1", notice.name, ratio)
		}
	}
	if paintedContrast(p.DisabledForeground, p.DisabledSurface, p.Background) < 3 {
		t.Error("ordinary disabled text must remain readable at 3:1")
	}
}

func TestColorsResolveEveryPresetAndAppearance(t *testing.T) {
	previous := preferenceSnapshot()
	SetPreferences(Preferences{})
	t.Cleanup(func() { SetPreferences(previous) })
	for _, tc := range []struct {
		name   string
		preset Preset
		dark   bool
	}{
		{"blue light", Blue, false}, {"blue dark", Blue, true},
		{"neutral light", Neutral, false}, {"neutral dark", Neutral, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var colors ColorTokens
			var native ui.Theme
			tt := ui.NewTester(func(c *ui.Context) {
				native = *ApplyPreset(c, tc.preset)
				colors = Colors(c)
				ui.Text(c, "Readable")
			}, 200, 100)
			tt.SetDark(tc.dark)
			tt.Frame()
			assertReadablePairs(t, colors)
			if colors.Primary != native.Accent || colors.Accent == native.Accent {
				t.Fatal("native actionable Accent must map to Primary, separate from semantic Accent")
			}
			if colors.Primary == colors.PrimaryPressed || colors.Secondary == colors.SecondaryPressed || colors.Accent == colors.AccentPressed {
				t.Fatal("pressed feedback must remain distinguishable from resting surfaces")
			}
			if colors.Chart1.A == 0 || colors.Chart2 == colors.Chart1 || colors.Chart3 == colors.Chart4 || colors.Chart5.A == 0 {
				t.Fatal("chart series must be present and distinct")
			}
		})
	}
}

func TestColorsSystemAccentAndHighContrast(t *testing.T) {
	previous := preferenceSnapshot()
	SetPreferences(Preferences{})
	t.Cleanup(func() { SetPreferences(previous) })
	for _, dark := range []bool{false, true} {
		t.Run(fmt.Sprint(dark), func(t *testing.T) {
			var p ColorTokens
			tt := ui.NewTester(func(c *ui.Context) { Apply(c); p = Colors(c); ui.Text(c, "Readable") }, 200, 100)
			tt.SetDark(dark)
			neutralAccent := p.Accent
			for _, accent := range []ui.Color{ui.RGB(220, 180, 20), ui.RGB(20, 60, 180), ui.RGB(100, 180, 100)} {
				tt.SetPreferences(ui.Preferences{Accent: accent, TextScale: 1, HighContrast: true})
				assertReadablePairs(t, p)
				if p.Primary != accent || p.SidebarPrimary != accent || p.Accent != neutralAccent {
					t.Fatal("OS accent must change primary companions while retaining neutral interaction surfaces")
				}
				for _, surface := range []ui.Color{p.Background, p.Card} {
					if paintedContrast(p.Border, surface, p.Background) < 3 || paintedContrast(p.Ring, surface, p.Background) < 3 {
						t.Fatal("high contrast borders and focus must meet nontext 3:1")
					}
					if paintedContrast(p.MutedForeground, surface, p.Background) < 7 {
						t.Fatal("high contrast secondary text must reach 7:1")
					}
				}
				if p.Border.A != 255 || p.Ring.A != 255 {
					t.Fatal("high contrast must provide opaque boundaries")
				}
				if paintedContrast(p.DisabledForeground, p.DisabledSurface, p.Background) < 4.5 {
					t.Fatal("disabled labels must remain readable under high contrast")
				}
			}
		})
	}
}

func TestContrastPreviewPreservesSystemPreferences(t *testing.T) {
	var p ColorTokens
	var system ui.Preferences
	var resolvedHighContrast bool
	preview := true
	tt := ui.NewTester(func(c *ui.Context) {
		if preview {
			ApplyHighContrast(c, Neutral)
		} else {
			ApplyPreset(c, Neutral)
		}
		p, system = Colors(c), c.Preferences()
		resolvedHighContrast = HighContrast(c)
		ui.Text(c, "Preview")
	}, 200, 100)
	if system.HighContrast || !resolvedHighContrast || paintedContrast(p.Border, p.Background, p.Background) < 3 {
		t.Fatal("preview must strengthen resolved policy without spoofing system input")
	}
	preview = false
	tt.Frame()
	if resolvedHighContrast {
		t.Fatal("leaving preview must restore ordinary contrast policy")
	}
}

func TestColorsFollowScopedOverridesAndCacheCopies(t *testing.T) {
	var p ColorTokens
	var override ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		base := Apply(c)
		override = *base
		override.Surface = mix(base.Surface, base.Text, .03)
		override.Accent = ui.RGB(180, 160, 30)
		override.AccentText = black
		override.AccentHover, override.AccentPressed = override.Accent, override.Accent
		c.SetTheme(&override)
		p = Colors(c)
		ui.Text(c, "Scoped")
	}, 200, 100)
	tt.Frame()
	if p.Card != override.Surface || p.Primary != override.Accent || p.PrimaryForeground != override.AccentText {
		t.Fatal("scoped native theme overrides lost")
	}
	original := p
	p.Primary = ui.Color{}
	tt.Frame()
	if p != original {
		t.Fatal("mutating a returned ColorTokens value changed cached authority")
	}
}

func TestColorsWarmFrameDoesNotAllocate(t *testing.T) {
	var measured bool
	tt := ui.NewTester(func(c *ui.Context) {
		Apply(c)
		Colors(c)
		if !measured {
			measured = true
			if allocs := testing.AllocsPerRun(100, func() { Apply(c); Colors(c) }); allocs != 0 {
				t.Errorf("warm theme frame allocated %g objects", allocs)
			}
		}
		ui.Text(c, "Stable")
	}, 200, 100)
	tt.Frame()
}

func TestMetricsPreserveRhythmAndGrowForText(t *testing.T) {
	previous := preferenceSnapshot()
	SetPreferences(Preferences{})
	t.Cleanup(func() { SetPreferences(previous) })
	var metrics MetricTokens
	tt := ui.NewTester(func(c *ui.Context) { metrics = Metrics(c); ui.Text(c, "Scaled") }, 200, 100)
	baseline := metrics
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	if metrics.ControlHeight <= baseline.ControlHeight || metrics.IconSize != 2*baseline.IconSize || metrics.LineHeight != 2*baseline.LineHeight {
		t.Fatal("text-bearing library controls must grow with scaled text")
	}
	if Space(.5) != 2 || Radii().LG != RadiusCard || Radii().XXXL != 22 {
		t.Fatal("geometry must retain shared spacing and radius authority")
	}
}
