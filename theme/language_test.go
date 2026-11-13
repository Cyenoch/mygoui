package theme_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// useTheme publishes t and restores the built-in theme and empty
// preferences when the test ends.
func useTheme(t *testing.T, th theme.Theme) {
	t.Helper()
	t.Cleanup(func() {
		if err := theme.Use(theme.Default()); err != nil {
			t.Error(err)
		}
		theme.Configure(theme.Preferences{})
	})
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}
}

// published is what a component reads from the theme, in comparable form.
type published struct {
	light, dark   theme.Palette
	radii         theme.Radii
	space, size   float32
	body, heading theme.TextStyle
	family, mono  string
	motion        [8]time.Duration
	effects       theme.EffectTokens
	overlayLayers int
	overlayEdge   theme.Shadow
}

func readPublished() published {
	m := theme.Motion()
	overlay := theme.OverlayShadows(true)
	return published{
		light: theme.ColorsForMode(false), dark: theme.ColorsForMode(true),
		radii: theme.RadiusScale(), space: theme.Space(1), size: theme.FontSize(),
		body: theme.Base(), heading: theme.H1(),
		family: theme.FontFamily(nil), mono: theme.MonoFontFamily(),
		motion:        [8]time.Duration{m.State, m.Reveal, m.Roll, m.Shift, m.Page, m.Popup, m.Dialog, m.Drawer},
		effects:       theme.Effects(),
		overlayLayers: len(overlay), overlayEdge: overlay[len(overlay)-1],
	}
}

// An application that publishes the built-in theme sees what one that
// publishes nothing sees.
func TestDefaultThemeIsWhatAnUnconfiguredApplicationSees(t *testing.T) {
	theme.Configure(theme.Preferences{})
	before := readPublished()
	useTheme(t, theme.Default())
	if after := readPublished(); after != before {
		t.Fatalf("publishing theme.Default() changed the tokens:\n got %+v\nwant %+v", after, before)
	}
}

// The built-in theme is the look the library had before themes were values:
// these are the timings, effect parameters and geometry components used as
// literals, so a changed default is a changed look.
func TestDefaultThemeKeepsThePreviousValues(t *testing.T) {
	d := theme.Default()
	ms := time.Millisecond
	for _, tc := range []struct {
		name      string
		got, want time.Duration
	}{
		{"Motion.State", d.Motion.State, 150 * ms},
		{"Motion.Reveal", d.Motion.Reveal, 200 * ms},
		{"Motion.Roll", d.Motion.Roll, 300 * ms},
		{"Motion.Shift", d.Motion.Shift, 150 * ms},
		{"Motion.Page", d.Motion.Page, 250 * ms},
		{"Motion.Popup", d.Motion.Popup, 150 * ms},
		{"Motion.Dialog", d.Motion.Dialog, 200 * ms},
		{"Motion.Drawer", d.Motion.Drawer, 450 * ms},
		{"Effects.SpinnerPeriod", d.Effects.SpinnerPeriod, 1000 * ms},
		{"Effects.SweepPeriod", d.Effects.SweepPeriod, 1500 * ms},
		{"Effects.LoadingDelay", d.Effects.LoadingDelay, 150 * ms},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
	for _, tc := range []struct {
		name      string
		got, want float32
	}{
		{"Effects.HoverOpacity", d.Effects.HoverOpacity, 0.9},
		{"Effects.PopupScale", d.Effects.PopupScale, 0.98},
		{"Effects.BackdropBlur", d.Effects.BackdropBlur, 4},
		{"Effects.FadeEdge", d.Effects.FadeEdge, 24},
		{"Effects.SweepBand", d.Effects.SweepBand, 0.4},
		{"Effects.RefreshOpacity", d.Effects.RefreshOpacity, 0.55},
		{"Effects.ThumbPressScale", d.Effects.ThumbPressScale, 1.2},
		{"Spacing.Unit", d.Spacing.Unit, 4},
		{"Radius", d.Radius, 10},
		{"Typography.Size", d.Typography.Size, 16},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %v, want %v", tc.name, tc.got, tc.want)
		}
	}
	assertCurve(t, "Motion.Ease", d.Motion.Ease, ui.EaseOut)
	assertCurve(t, "Motion.EaseTravel", d.Motion.EaseTravel, ui.EaseInOut)
	assertCurve(t, "Motion.EaseDrawer", d.Motion.EaseDrawer, ui.CubicBezier(0.32, 0.72, 0, 1))
}

// assertCurve compares two easings where they are sampled: funcs do not
// compare.
func assertCurve(t *testing.T, name string, got, want ui.Easing) {
	t.Helper()
	if got == nil {
		t.Errorf("%s is nil", name)
		return
	}
	for _, x := range []float32{0.1, 0.25, 0.5, 0.75, 0.9} {
		if got(x) != want(x) {
			t.Errorf("%s(%v) = %v, want %v", name, x, got(x), want(x))
		}
	}
}

// A saved and restored preference value pins nothing: the font size and
// radius it carries are the theme's, so a theme published later sets them.
func TestRestoredPreferencesStillFollowALaterTheme(t *testing.T) {
	useTheme(t, theme.Default())
	saved := theme.PreferencesSnapshot()
	saved.Accent = ui.RGB(200, 40, 40)
	theme.Configure(saved)

	th := theme.Default()
	th.Radius, th.Typography.Size = 4, 20
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}
	if got := theme.RadiusBase(); got != 4 {
		t.Errorf("RadiusBase = %v, want the later theme's 4", got)
	}
	if got := theme.FontSize(); got != 20 {
		t.Errorf("FontSize = %v, want the later theme's 20", got)
	}
	if got := theme.ColorsForMode(false).Primary; got != saved.Accent {
		t.Errorf("Primary = %v, want the restored accent", got)
	}
}

// A theme that would make a component divide by zero, or make the toolkit
// panic, is refused by name and leaves the published theme in place.
func TestUseRejectsDurationsComponentsCannotRun(t *testing.T) {
	useTheme(t, theme.Default())
	for _, tc := range []struct {
		field string
		edit  func(*theme.Theme)
	}{
		{"Motion.State", func(th *theme.Theme) { th.Motion.State = -time.Millisecond }},
		{"Effects.LoadingDelay", func(th *theme.Theme) { th.Effects.LoadingDelay = -time.Millisecond }},
		{"Effects.SweepPeriod", func(th *theme.Theme) { th.Effects.SweepPeriod = 0 }},
		{"Effects.SpinnerPeriod", func(th *theme.Theme) { th.Effects.SpinnerPeriod = time.Microsecond }},
	} {
		th := theme.Default()
		th.Radius = 3
		tc.edit(&th)
		err := theme.Use(th)
		if err == nil || !strings.Contains(err.Error(), tc.field) {
			t.Errorf("Use with a bad %s returned %v, want an error naming the field", tc.field, err)
		}
		if got := theme.RadiusBase(); got != 10 {
			t.Errorf("a refused theme was published: RadiusBase = %v", got)
		}
	}
}

func TestThemeSuppliesPaletteAndGeometryTogether(t *testing.T) {
	brand, paper := ui.RGB(30, 90, 180), ui.RGB(20, 22, 28)
	th := theme.Default()
	th.Light.Primary, th.Dark.Background = brand, paper
	th.Radius, th.Spacing.Unit = 4, 5
	th.Typography.Size, th.Typography.Base.Size = 14, 15
	useTheme(t, th)

	if got := theme.ColorsForMode(false).Primary; got != brand {
		t.Errorf("light Primary = %v, want the theme's %v", got, brand)
	}
	if got := theme.ColorsForMode(true).Background; got != paper {
		t.Errorf("dark Background = %v, want the theme's %v", got, paper)
	}
	if got := theme.RadiusScale().LG; got != 4 {
		t.Errorf("radius LG = %v, want the theme's base 4", got)
	}
	if got := theme.Space(2); got != 10 {
		t.Errorf("Space(2) = %v, want two 5 DIP units", got)
	}
	if got := theme.FontSize(); got != 14 {
		t.Errorf("FontSize = %v, want the theme's 14", got)
	}
	if got := theme.Base().Size; got != 15 {
		t.Errorf("Base role size = %v, want the theme's 15", got)
	}
}

func TestPartialThemeOverrideInheritsTheRest(t *testing.T) {
	builtin := theme.Default()
	th := theme.Default()
	th.Motion.State = 80 * time.Millisecond
	useTheme(t, th)

	m := theme.Motion()
	if m.State != 80*time.Millisecond {
		t.Fatalf("State = %v, want the override", m.State)
	}
	if m.Reveal != builtin.Motion.Reveal || m.Drawer != builtin.Motion.Drawer {
		t.Errorf("other durations moved: Reveal %v, Drawer %v", m.Reveal, m.Drawer)
	}
	if theme.ColorsForMode(false) != builtin.Light || theme.ColorsForMode(true) != builtin.Dark {
		t.Error("a motion override changed a palette")
	}
	if theme.Effects() != builtin.Effects || theme.RadiusScale().LG != builtin.Radius {
		t.Error("a motion override changed effects or radius")
	}
}

// A theme whose easings are left nil still animates: the built-in curves
// stand in, so a dialog can ask its travel curve for a position.
func TestThemeWithoutEasingsUsesBuiltInCurves(t *testing.T) {
	th := theme.Default()
	th.Motion.Ease, th.Motion.EaseTravel, th.Motion.EaseDrawer = nil, nil, nil
	useTheme(t, th)
	m := theme.Motion()
	assertCurve(t, "Motion.Ease", m.Ease, ui.EaseOut)
	assertCurve(t, "Motion.EaseTravel", m.EaseTravel, ui.EaseInOut)
	assertCurve(t, "Motion.EaseDrawer", m.EaseDrawer, ui.CubicBezier(0.32, 0.72, 0, 1))
}

func TestThemeCopiesCallerOwnedLayers(t *testing.T) {
	th := theme.Default()
	want := th.Elevation.Light.Overlay[0]
	useTheme(t, th)

	th.Elevation.Light.Overlay[0].Y = 99
	if got := theme.OverlayShadows(false)[0]; got != want {
		t.Errorf("editing the value after Use changed the published layer: %+v", got)
	}
	returned := theme.Current()
	returned.Elevation.Light.Overlay[0].Y = 77
	if got := theme.OverlayShadows(false)[0]; got != want {
		t.Errorf("editing the value Current returned changed the published layer: %+v", got)
	}
}

func TestPreferencesResolveOverACustomTheme(t *testing.T) {
	brand, link, paper := ui.RGB(30, 90, 180), ui.RGB(0, 110, 200), ui.RGB(250, 248, 240)
	th := theme.Default()
	th.Light.Primary, th.Light.Link = brand, link
	th.Radius = 6
	useTheme(t, th)

	theme.Configure(theme.Preferences{
		FontFamily: "Test Sans",
		Palette:    theme.PaletteOverrides{Light: &theme.PalettePatch{Background: &paper}},
	})
	light := theme.ColorsForMode(false)
	if light.Background != paper {
		t.Errorf("Background = %v, want the preference %v", light.Background, paper)
	}
	if light.FieldBackground != paper {
		t.Errorf("FieldBackground = %v, want it to follow the overridden Background", light.FieldBackground)
	}
	if light.Primary != brand || light.Link != link {
		t.Errorf("roles the preference does not touch left the theme: Primary %v, Link %v", light.Primary, light.Link)
	}
	if theme.ColorsForMode(true) != th.Dark {
		t.Error("a light override changed the dark palette")
	}

	accent := ui.RGB(200, 40, 40)
	theme.UpdatePreferences(func(p *theme.Preferences) { p.Accent = accent })
	light = theme.ColorsForMode(false)
	if light.Primary != accent || light.Selection != accent.Alpha(0.2) {
		t.Errorf("Accent did not replace Primary and its selection: %v, %v", light.Primary, light.Selection)
	}
	if light.Background != paper || light.Link != link {
		t.Error("an Accent update dropped the palette override or the theme's Link")
	}

	// The update above left radius alone, so radii keep following the theme.
	th.Radius = 4
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}
	if got := theme.RadiusScale().LG; got != 4 {
		t.Errorf("radius LG = %v after a new theme, want 4", got)
	}
	if got := theme.ColorsForMode(false).Background; got != paper {
		t.Errorf("publishing a theme dropped the palette preference: %v", got)
	}
	if got := theme.PreferencesSnapshot().FontFamily; got != "Test Sans" {
		t.Errorf("publishing a theme dropped the font preference: %q", got)
	}

	theme.UpdatePreferences(func(p *theme.Preferences) { p.Radius = 12 })
	if got := theme.RadiusScale().LG; got != 12 {
		t.Errorf("radius LG = %v, want the preference 12 over the theme", got)
	}
}
