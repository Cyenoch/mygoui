package theme_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Distinct values expose crossed wires between similarly named Coss tokens.
func palettePatch() theme.PalettePatch {
	return theme.PalettePatch{
		Background:               new(ui.Hex("#102030")),
		Foreground:               new(ui.Hex("#e1e2e3")),
		Card:                     new(ui.Hex("#203040")),
		CardForeground:           new(ui.Hex("#d1d2d3")),
		Popover:                  new(ui.Hex("#304050")),
		PopoverForeground:        new(ui.Hex("#c1c2c3")),
		Primary:                  new(ui.Hex("#40506080")),
		PrimaryForeground:        new(ui.Hex("#b1b2b3")),
		Secondary:                new(ui.Hex("#50607080")),
		SecondaryForeground:      new(ui.Hex("#a1a2a3")),
		MutedBackground:          new(ui.Hex("#60708040")),
		MutedForeground:          new(ui.Hex("#919293")),
		Accent:                   new(ui.Hex("#70809060")),
		AccentForeground:         new(ui.Hex("#818283")),
		Destructive:              new(ui.Hex("#c0304080")),
		DestructiveForeground:    new(ui.Hex("#d04050")),
		Info:                     new(ui.Hex("#3060c080")),
		InfoForeground:           new(ui.Hex("#4070d0")),
		Success:                  new(ui.Hex("#30c06080")),
		SuccessForeground:        new(ui.Hex("#40d070")),
		Warning:                  new(ui.Hex("#c0903080")),
		WarningForeground:        new(ui.Hex("#d0a040")),
		Border:                   new(ui.Hex("#8090a020")),
		Input:                    new(ui.Hex("#90a0b030")),
		Ring:                     new(ui.Hex("#a0b0c0")),
		Chart1:                   new(ui.Hex("#112233")),
		Chart2:                   new(ui.Hex("#223344")),
		Chart3:                   new(ui.Hex("#334455")),
		Chart4:                   new(ui.Hex("#445566")),
		Chart5:                   new(ui.Hex("#556677")),
		Sidebar:                  new(ui.Hex("#182838")),
		SidebarForeground:        new(ui.Hex("#e8d8c8")),
		SidebarPrimary:           new(ui.Hex("#283848")),
		SidebarPrimaryForeground: new(ui.Hex("#d8c8b8")),
		SidebarAccent:            new(ui.Hex("#38485850")),
		SidebarAccentForeground:  new(ui.Hex("#c8b8a8")),
		SidebarBorder:            new(ui.Hex("#48586820")),
		SidebarRing:              new(ui.Hex("#586878")),
		Code:                     new(ui.Hex("#213141")),
		CodeForeground:           new(ui.Hex("#d1c1b1")),
		CodeHighlight:            new(ui.Hex("#31415170")),
	}
}

type paletteToken struct {
	name  string
	patch **ui.Color
	color ui.Color
}

func cossTokens(p *theme.PalettePatch, c theme.Palette) []paletteToken {
	return []paletteToken{
		{name: "background", patch: &p.Background, color: c.Background},
		{name: "foreground", patch: &p.Foreground, color: c.Foreground},
		{name: "card", patch: &p.Card, color: c.Card},
		{name: "card foreground", patch: &p.CardForeground, color: c.CardForeground},
		{name: "popover", patch: &p.Popover, color: c.Popover},
		{name: "popover foreground", patch: &p.PopoverForeground, color: c.PopoverForeground},
		{name: "primary", patch: &p.Primary, color: c.Primary},
		{name: "primary foreground", patch: &p.PrimaryForeground, color: c.PrimaryForeground},
		{name: "secondary", patch: &p.Secondary, color: c.Secondary},
		{name: "secondary foreground", patch: &p.SecondaryForeground, color: c.SecondaryForeground},
		{name: "muted background", patch: &p.MutedBackground, color: c.MutedBackground},
		{name: "muted foreground", patch: &p.MutedForeground, color: c.MutedForeground},
		{name: "accent", patch: &p.Accent, color: c.Accent},
		{name: "accent foreground", patch: &p.AccentForeground, color: c.AccentForeground},
		{name: "destructive", patch: &p.Destructive, color: c.Destructive},
		{name: "destructive foreground", patch: &p.DestructiveForeground, color: c.DestructiveForeground},
		{name: "info", patch: &p.Info, color: c.Info},
		{name: "info foreground", patch: &p.InfoForeground, color: c.InfoForeground},
		{name: "success", patch: &p.Success, color: c.Success},
		{name: "success foreground", patch: &p.SuccessForeground, color: c.SuccessForeground},
		{name: "warning", patch: &p.Warning, color: c.Warning},
		{name: "warning foreground", patch: &p.WarningForeground, color: c.WarningForeground},
		{name: "border", patch: &p.Border, color: c.Border},
		{name: "input", patch: &p.Input, color: c.Input},
		{name: "ring", patch: &p.Ring, color: c.Ring},
		{name: "chart 1", patch: &p.Chart1, color: c.Chart1},
		{name: "chart 2", patch: &p.Chart2, color: c.Chart2},
		{name: "chart 3", patch: &p.Chart3, color: c.Chart3},
		{name: "chart 4", patch: &p.Chart4, color: c.Chart4},
		{name: "chart 5", patch: &p.Chart5, color: c.Chart5},
		{name: "sidebar", patch: &p.Sidebar, color: c.Sidebar},
		{name: "sidebar foreground", patch: &p.SidebarForeground, color: c.SidebarForeground},
		{name: "sidebar primary", patch: &p.SidebarPrimary, color: c.SidebarPrimary},
		{name: "sidebar primary foreground", patch: &p.SidebarPrimaryForeground, color: c.SidebarPrimaryForeground},
		{name: "sidebar accent", patch: &p.SidebarAccent, color: c.SidebarAccent},
		{name: "sidebar accent foreground", patch: &p.SidebarAccentForeground, color: c.SidebarAccentForeground},
		{name: "sidebar border", patch: &p.SidebarBorder, color: c.SidebarBorder},
		{name: "sidebar ring", patch: &p.SidebarRing, color: c.SidebarRing},
		{name: "code", patch: &p.Code, color: c.Code},
		{name: "code foreground", patch: &p.CodeForeground, color: c.CodeForeground},
		{name: "code highlight", patch: &p.CodeHighlight, color: c.CodeHighlight},
	}
}

func TestPalettePatchOverridesEveryCossToken(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	for _, mode := range []struct {
		name string
		dark bool
	}{
		{name: "light"},
		{name: "dark", dark: true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			theme.Configure(theme.Preferences{})
			untouched := theme.ColorsForMode(!mode.dark)
			var template theme.PalettePatch
			for i, token := range cossTokens(&template, theme.Palette{}) {
				t.Run(token.name, func(t *testing.T) {
					for _, value := range []struct {
						name  string
						color ui.Color
					}{
						{name: "translucent", color: ui.Hex("#17436580")},
						{name: "transparent", color: ui.Transparent},
					} {
						t.Run(value.name, func(t *testing.T) {
							var patch theme.PalettePatch
							*cossTokens(&patch, theme.Palette{})[i].patch = new(value.color)
							overrides := theme.PaletteOverrides{Light: &patch}
							if mode.dark {
								overrides = theme.PaletteOverrides{Dark: &patch}
							}
							theme.Configure(theme.Preferences{Palette: overrides})
							got := cossTokens(&patch, theme.ColorsForMode(mode.dark))[i].color
							if got != value.color {
								t.Errorf("color = %v, want verbatim override %v", got, value.color)
							}
							if theme.ColorsForMode(!mode.dark) != untouched {
								t.Error("patch changed the other appearance")
							}
						})
					}
				})
			}
		})
	}
}

func TestApplyInstallsFullCossPalette(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	light, dark := palettePatch(), palettePatch()
	dark.Primary = new(ui.Hex("#c8d8e8"))
	dark.PrimaryForeground = new(ui.Hex("#182838"))
	theme.Configure(theme.Preferences{
		Accent:  ui.RGB(255, 0, 0),
		Palette: theme.PaletteOverrides{Light: &light, Dark: &dark},
	})
	for _, mode := range []struct {
		name  string
		dark  bool
		patch *theme.PalettePatch
	}{
		{name: "light", patch: &light},
		{name: "dark", dark: true, patch: &dark},
	} {
		t.Run(mode.name, func(t *testing.T) {
			var got theme.Palette
			var installed *ui.Theme
			var previous ui.Theme
			var frames int
			tt := ui.NewTester(func(c *ui.Context) {
				installed = theme.ApplyMode(c, mode.dark)
				got = theme.Colors(c)
				if again := theme.Apply(c); again != installed {
					t.Error("reapplying the current theme replaced the installed theme")
				}
				if frames > 0 && previous != *installed {
					t.Error("an unchanged frame changed the installed theme values")
				}
				previous = *installed
				frames++
			}, 120, 80)
			tt.Frame()
			tt.Frame()

			for _, token := range cossTokens(mode.patch, got) {
				t.Run(token.name, func(t *testing.T) {
					if token.color != **token.patch {
						t.Errorf("Colors token = %v, want %v", token.color, **token.patch)
					}
				})
			}
			for _, token := range []struct {
				name      string
				got, want ui.Color
			}{
				{name: "background", got: installed.Background, want: *mode.patch.Background},
				{name: "surface", got: installed.Surface, want: *mode.patch.Popover},
				{name: "surface hover", got: installed.SurfaceHover, want: mode.patch.Accent.Over(*mode.patch.Popover)},
				{name: "surface pressed", got: installed.SurfacePressed, want: mode.patch.Accent.Over(*mode.patch.Popover)},
				{name: "border", got: installed.Border, want: *mode.patch.Border},
				{name: "text", got: installed.Text, want: *mode.patch.Foreground},
				{name: "muted", got: installed.TextMuted, want: *mode.patch.MutedForeground},
				{name: "primary", got: installed.Accent, want: *mode.patch.Primary},
				{name: "primary hover", got: installed.AccentHover, want: mode.patch.Primary.Alpha(0.9)},
				{name: "primary pressed", got: installed.AccentPressed, want: mode.patch.Primary.Alpha(0.9)},
				{name: "primary foreground", got: installed.AccentText, want: *mode.patch.PrimaryForeground},
				{name: "selection", got: installed.Selection, want: mode.patch.Primary.Alpha(0.2)},
				{name: "destructive", got: installed.Danger, want: *mode.patch.Destructive},
				{name: "success", got: installed.Success, want: *mode.patch.Success},
				{name: "warning", got: installed.Warning, want: *mode.patch.Warning},
				{name: "focus", got: installed.Focus, want: *mode.patch.Ring},
			} {
				t.Run("native "+token.name, func(t *testing.T) {
					if token.got != token.want {
						t.Errorf("native color = %v, want %v", token.got, token.want)
					}
				})
			}
			if installed.Dark != mode.dark {
				t.Errorf("native dark = %v, want %v", installed.Dark, mode.dark)
			}
		})
	}
}

func TestDerivedColorsPreserveAlphaAndStatusForegrounds(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	patch := palettePatch()
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &patch, Dark: &patch}})
	for _, mode := range []struct {
		name  string
		dark  bool
		alpha float32
	}{
		{name: "light", alpha: 0.08},
		{name: "dark", dark: true, alpha: 0.16},
	} {
		t.Run(mode.name, func(t *testing.T) {
			p := theme.ColorsForMode(mode.dark)
			for _, token := range []struct {
				name      string
				got, want ui.Color
			}{
				{name: "secondary hover", got: p.DefaultHover, want: patch.Secondary.Alpha(0.9)},
				{name: "accent hover", got: p.AccentHover, want: patch.Accent.Alpha(0.9)},
				{name: "accent soft", got: p.AccentSoft, want: *patch.Accent},
				{name: "info hover", got: p.InfoHover, want: patch.Info.Alpha(0.9)},
				{name: "info soft", got: p.InfoSoft, want: patch.Info.Alpha(mode.alpha)},
				{name: "success soft", got: p.SuccessSoft, want: patch.Success.Alpha(mode.alpha)},
				{name: "warning soft", got: p.WarningSoft, want: patch.Warning.Alpha(mode.alpha)},
				{name: "destructive soft", got: p.DangerSoft, want: patch.Destructive.Alpha(mode.alpha)},
				{name: "info soft foreground", got: p.InfoSoftForeground, want: *patch.InfoForeground},
				{name: "success soft foreground", got: p.SuccessSoftForeground, want: *patch.SuccessForeground},
				{name: "warning soft foreground", got: p.WarningSoftForeground, want: *patch.WarningForeground},
				{name: "destructive soft foreground", got: p.DangerSoftForeground, want: *patch.DestructiveForeground},
				{name: "secondary canvas", got: p.BackgroundSecondary, want: patch.MutedBackground.Over(*patch.Background)},
				{name: "tertiary canvas", got: p.BackgroundTertiary, want: patch.Input.Over(*patch.Background)},
			} {
				t.Run(token.name, func(t *testing.T) {
					if token.got != token.want {
						t.Errorf("derived color = %v, want %v", token.got, token.want)
					}
				})
			}
		})
	}
}

func convenienceTokens(p *theme.PalettePatch) []paletteToken {
	return []paletteToken{
		{name: "surface", patch: &p.Surface},
		{name: "muted", patch: &p.Muted},
		{name: "field border", patch: &p.FieldBorder},
		{name: "danger", patch: &p.Danger},
		{name: "danger foreground", patch: &p.DangerForeground},
		{name: "hover", patch: &p.Hover},
		{name: "selected", patch: &p.Selected},
		{name: "overlay", patch: &p.Overlay},
		{name: "backdrop", patch: &p.Backdrop},
	}
}

func TestPaletteSnapshotsDetachCallerPointers(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	light, dark := palettePatch(), palettePatch()
	for _, patch := range []*theme.PalettePatch{&light, &dark} {
		for _, token := range convenienceTokens(patch) {
			*token.patch = new(ui.Hex("#38685880"))
		}
	}
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &light, Dark: &dark}})
	prefs := theme.PreferencesSnapshot()
	var want [2][]ui.Color
	for i, patch := range []*theme.PalettePatch{prefs.Palette.Light, prefs.Palette.Dark} {
		for _, token := range append(cossTokens(patch, theme.Palette{}), convenienceTokens(patch)...) {
			want[i] = append(want[i], **token.patch)
		}
	}
	wantLight, wantDark := theme.ColorsForMode(false), theme.ColorsForMode(true)

	assertUnchanged := func(t *testing.T) {
		t.Helper()
		got := theme.PreferencesSnapshot()
		for i, patch := range []*theme.PalettePatch{got.Palette.Light, got.Palette.Dark} {
			for j, token := range append(cossTokens(patch, theme.Palette{}), convenienceTokens(patch)...) {
				if *token.patch == nil || **token.patch != want[i][j] {
					t.Errorf("caller mutation leaked into mode %d preference %s", i, token.name)
				}
			}
		}
		if theme.ColorsForMode(false) != wantLight || theme.ColorsForMode(true) != wantDark {
			t.Error("caller mutation changed a published palette")
		}
	}
	mutate := func(patch *theme.PalettePatch) {
		for _, token := range append(cossTokens(patch, theme.Palette{}), convenienceTokens(patch)...) {
			**token.patch = ui.Hex("#f0e0d040")
		}
		patch.Primary = new(ui.Transparent)
	}
	t.Run("configure input", func(t *testing.T) {
		mutate(&light)
		mutate(&dark)
		assertUnchanged(t)
	})
	t.Run("preferences snapshot", func(t *testing.T) {
		returned := theme.PreferencesSnapshot()
		mutate(returned.Palette.Light)
		mutate(returned.Palette.Dark)
		assertUnchanged(t)
	})
}

func TestForegroundPatchInheritsWithoutReplacingExplicitPairs(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	ink, explicit := ui.Hex("#d8e8f8"), ui.Transparent
	for _, mode := range []struct {
		name string
		dark bool
	}{
		{name: "light"},
		{name: "dark", dark: true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			patch := theme.PalettePatch{Foreground: &ink}
			for _, override := range []bool{false, true} {
				want := ink
				if override {
					patch.CardForeground, patch.PopoverForeground = &explicit, &explicit
					patch.SecondaryForeground, patch.AccentForeground, patch.CodeForeground = &explicit, &explicit, &explicit
					want = explicit
				}
				o := theme.PaletteOverrides{Light: &patch}
				if mode.dark {
					o = theme.PaletteOverrides{Dark: &patch}
				}
				theme.Configure(theme.Preferences{Palette: o})
				p := theme.ColorsForMode(mode.dark)
				for _, pair := range []struct {
					name  string
					color ui.Color
				}{
					{name: "card", color: p.CardForeground},
					{name: "popover", color: p.PopoverForeground},
					{name: "secondary", color: p.SecondaryForeground},
					{name: "accent", color: p.AccentForeground},
					{name: "code", color: p.CodeForeground},
				} {
					if pair.color != want {
						t.Errorf("%s foreground (explicit=%v) = %v, want %v", pair.name, override, pair.color, want)
					}
				}
			}
		})
	}
}

func TestConveniencePatchStillReachesConsumers(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	surface, ink := ui.Hex("#405060"), ui.Hex("#e0e8f0")
	muted, selected := ui.Hex("#8090a0"), ui.Hex("#c0306080")
	patch := theme.PalettePatch{
		Surface: &surface, Foreground: &ink, Muted: &muted,
		FieldBorder: new(ui.Transparent), Border: new(ui.Transparent),
		Selected: &selected,
	}
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &patch}})
	var got theme.Palette
	var native ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		native = *theme.ApplyMode(c, false)
		got = theme.Colors(c)
	}, 120, 80)
	tt.Frame()
	if got.Surface != surface || got.FieldBackground != surface || native.Surface != surface {
		t.Error("surface customization did not reach library and native control faces")
	}
	if got.FieldForeground != ink || native.Text != ink || got.FieldPlaceholder != muted || native.TextMuted != muted {
		t.Error("consumer ink customization was lost")
	}
	if got.FieldBorder != ui.Transparent || native.Border != ui.Transparent {
		t.Error("explicit transparent border override was lost")
	}
	if got.Selection != selected || native.Selection != selected {
		t.Error("selection customization did not reach native selection")
	}
	if native.SurfaceHover != got.Accent.Over(surface) || native.SurfacePressed != got.Accent.Over(surface) {
		t.Error("native state face was not composited over the explicit surface")
	}
}

func TestUpdatePreferencesKeepsInjectedPalette(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	patch := palettePatch()
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{Light: &patch}})
	want := theme.ColorsForMode(false)
	theme.UpdatePreferences(func(p *theme.Preferences) { p.FontFamily = "Test Sans" })
	if got := theme.PreferencesSnapshot().FontFamily; got != "Test Sans" {
		t.Fatalf("font family = %q, want the edited value", got)
	}
	if got := theme.ColorsForMode(false); got != want {
		t.Error("typography edit changed the injected palette")
	}
	theme.Configure(theme.Preferences{FontFamily: "Test Sans"})
	if got := theme.ColorsForMode(false).Background; got == *patch.Background {
		t.Error("Configure kept a palette it was not given")
	}
}
