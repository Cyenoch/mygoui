package theme

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestDarkSurfaceHierarchy(t *testing.T) {
	p := ColorsForMode(true, Blue)
	if luminance(p.Card) <= luminance(p.Background) {
		t.Fatal("card must be elevated above the canvas")
	}
	if p.Accent.A == 0 || p.Border.A == 0 || p.SidebarRowSelected.A == 0 {
		t.Fatal("dark interaction and divider treatments must remain paintable")
	}
	if p.Primary == p.Accent || p.Selection == p.Accent {
		t.Fatal("action, selection and neutral hover roles must remain independent")
	}
}

func TestOverPreservesBorderContrast(t *testing.T) {
	t.Parallel()
	p := ColorsForMode(true, Blue)
	for _, tc := range []struct {
		name       string
		color      ui.Color
		wantStrong bool
	}{
		{name: "composited contrast", color: Over(p.Border, p.Background, 0.7), wantStrong: true},
		{name: "old double multiplication", color: p.Border.Alpha(float32(p.Border.A) / 255 * 0.7)},
		{name: "plain alpha", color: p.Border.Alpha(0.7), wantStrong: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			painted := tc.color.Over(p.Background)
			delta := max(
				int(painted.R)-int(p.Background.R),
				int(painted.G)-int(p.Background.G),
				int(painted.B)-int(p.Background.B),
			)
			if (delta >= 6) != tc.wantStrong {
				t.Errorf("painted border = %+v, canvas = %+v, maximum channel delta = %d; strong = %v",
					painted, p.Background, delta, tc.wantStrong)
			}
			if tc.name == "plain alpha" && (tc.color.A == 0 || tc.color.A >= p.Border.A) {
				t.Errorf("Alpha must multiply existing opacity, got A=%d from A=%d", tc.color.A, p.Border.A)
			}
			t.Logf("painted border RGB=(%d,%d,%d), canvas RGB=(%d,%d,%d), delta=%d, source alpha=%d",
				painted.R, painted.G, painted.B, p.Background.R, p.Background.G, p.Background.B, delta, tc.color.A)
		})
	}
}

func TestSurfaceBorderContrast(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		dark bool
	}{
		{name: "light"},
		{name: "dark", dark: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ColorsForMode(tc.dark, Blue)
			cardFill := Over(p.Card, p.Background, 0.7)
			cardBorder := Over(p.Border, cardFill, 0.8)
			if cardBorder == cardFill || cardBorder == p.Background {
				t.Errorf("plan border %+v must differ from fill %+v and canvas %+v", cardBorder, cardFill, p.Background)
			}
			codeFill := p.Accent.Over(p.Background)
			oldCodeBorder := p.Border.Alpha(float32(p.Border.A) / 255 * 0.7)
			newCodeBorder := Over(p.Border, codeFill, 0.7)
			if tc.dark {
				oldCodeBorder, newCodeBorder = ui.Color{}, ui.Color{}
			}
			oldPixel, newPixel := oldCodeBorder.Over(codeFill), newCodeBorder.Over(codeFill)
			t.Logf("plan fill RGB=(%d,%d,%d), border RGB=(%d,%d,%d), canvas RGB=(%d,%d,%d)",
				cardFill.R, cardFill.G, cardFill.B, cardBorder.R, cardBorder.G, cardBorder.B,
				p.Background.R, p.Background.G, p.Background.B)
			t.Logf("code fill RGB=(%d,%d,%d), old border RGB=(%d,%d,%d) delta=(%d,%d,%d), new border RGB=(%d,%d,%d) delta=(%d,%d,%d)",
				codeFill.R, codeFill.G, codeFill.B, oldPixel.R, oldPixel.G, oldPixel.B,
				int(oldPixel.R)-int(codeFill.R), int(oldPixel.G)-int(codeFill.G), int(oldPixel.B)-int(codeFill.B),
				newPixel.R, newPixel.G, newPixel.B,
				int(newPixel.R)-int(codeFill.R), int(newPixel.G)-int(codeFill.G), int(newPixel.B)-int(codeFill.B))
		})
	}
}

func TestMixPremultipliesAlpha(t *testing.T) {
	// Transparent red contributes no red when mixed with opaque blue.
	got := mix(ui.RGBA(255, 0, 0, 0), ui.RGB(0, 0, 255), 0.5)
	want := ui.RGBA(0, 0, 255, 0.5)
	if got != want {
		t.Errorf("mix = %+v, want %+v", got, want)
	}
	if got := mix(ui.Color{}, ui.Color{}, 0.5); got != (ui.Color{}) {
		t.Errorf("transparent mix = %+v", got)
	}
}

func TestApplyFollowsPreferencesWithoutCompounding(t *testing.T) {
	previous := preferenceSnapshot()
	SetPreferences(Preferences{})
	t.Cleanup(func() { SetPreferences(previous) })
	var sizes []float32
	var applied ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		applied = *Apply(c)
		sizes = append(sizes, applied.FontSize)
		ui.Text(c, "Readable")
	}, 120, 80)
	tt.SetDark(true)
	tt.SetPreferences(ui.Preferences{TextScale: 2, HighContrast: true, ReduceMotion: true})
	tt.Frame()
	if applied.Background != ColorsForMode(true, Neutral).Background || applied.Border == ColorsForMode(true, Neutral).Border || applied.Focus.A != 255 {
		t.Errorf("high contrast dark theme not applied: %+v", applied)
	}
	if applied.FontSize != 16*2 {
		t.Errorf("scaled font = %g, want %g", applied.FontSize, float32(16*2))
	}
	if n := len(sizes); n < 2 || sizes[n-1] != sizes[n-2] {
		t.Errorf("font size compounded across frames: %v", sizes)
	}
	accent := ui.RGB(220, 180, 20)
	tt.SetPreferences(ui.Preferences{Accent: accent, TextScale: 1})
	if applied.Accent != accent || applied.AccentText == ui.RGB(255, 255, 255) {
		t.Errorf("OS accent and readable accent text were lost: %+v", applied)
	}
}
