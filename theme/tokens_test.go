package theme_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestHeadingRolesCarryTrackingTight pins HeroUI's type scale: text-4xl to
// text-base headings are semibold with tracking-tight, body copy is not, and
// the typography roles keep their larger line boxes.
func TestHeadingRolesCarryTrackingTight(t *testing.T) {
	for _, tc := range []struct {
		name string
		role theme.SemanticType
		size float32
		line float32
	}{
		{"h1", theme.TypeH1, 36, 40},
		{"h2", theme.TypeH2, 30, 36},
		{"h3", theme.TypeH3, 24, 32},
		{"h4", theme.TypeH4, 20, 28},
		{"h5", theme.TypeH5, 18, 28},
		{"h6", theme.TypeH6, 16, 24},
	} {
		style := tc.role.Style()
		if style.Size != tc.size || style.LineHeight != tc.line {
			t.Errorf("%s = %v/%v, want %v/%v", tc.name, style.Size, style.LineHeight, tc.size, tc.line)
		}
		if style.Weight != theme.WeightSemibold {
			t.Errorf("%s weight = %d, want %d", tc.name, style.Weight, theme.WeightSemibold)
		}
		if style.Tracking != theme.TrackingTight {
			t.Errorf("%s tracking = %v, want %v", tc.name, style.Tracking, theme.TrackingTight)
		}
	}

	// Prose roles: body 16/28, body-sm 14/24 and body-xs 12/20, all
	// untracked.
	for _, tc := range []struct {
		name string
		role theme.SemanticType
		size float32
		line float32
	}{
		{"body", theme.TypeBody, 16, 28},
		{"body-sm", theme.TypeBodySM, 14, 24},
		{"body-xs", theme.TypeBodyXS, 12, 20},
		{"code", theme.TypeCode, 14, 20},
	} {
		style := tc.role.Style()
		if style.Size != tc.size || style.LineHeight != tc.line {
			t.Errorf("%s = %v/%v, want %v/%v", tc.name, style.Size, style.LineHeight, tc.size, tc.line)
		}
		if style.Tracking != 0 {
			t.Errorf("%s tracking = %v, want none", tc.name, style.Tracking)
		}
	}
	if !theme.TypeCode.Style().Mono {
		t.Error("code role is not monospaced")
	}
}

// TestFieldBorderDefaultsTransparent pins upstream's common --field-border
// value and proves a host-injected border derives the hover and focus shades
// from itself instead of the HeroUI default.
func TestFieldBorderDefaultsTransparent(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	if got := theme.ColorsForMode(false).FieldBorder; got != ui.Transparent {
		t.Fatalf("default field border = %v, want transparent", got)
	}

	border := ui.RGB(120, 120, 130)
	ink := ui.RGB(20, 20, 20)
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{
		Light: &theme.PalettePatch{FieldBorder: &border, Foreground: &ink},
	}})
	got := theme.ColorsForMode(false)
	if got.FieldBorder != border {
		t.Fatalf("field border = %v, want the injected %v", got.FieldBorder, border)
	}
	toBorder := channelDistance(got.FieldBorderHover, border)
	toInk := channelDistance(got.FieldBorderHover, ink)
	if toBorder >= toInk {
		t.Errorf("field border hover = %v is not a shade of the injected border %v", got.FieldBorderHover, border)
	}
	if channelDistance(got.FieldBorderFocus, ink) >= channelDistance(got.FieldBorderHover, ink) {
		t.Errorf("field border focus = %v is not stronger than hover %v", got.FieldBorderFocus, got.FieldBorderHover)
	}
	if got.FieldBorderHover.A != 255 {
		t.Errorf("field border hover alpha = %d, want an opaque shade of the opaque border", got.FieldBorderHover.A)
	}
}

// TestShadowTablesAreReadOnly pins the shared elevation tables: a caller that
// appends to the slice it got must not reach the package's own array, and
// dark surfaces have no drop shadow, only the inset overlay hairline.
func TestShadowTablesAreReadOnly(t *testing.T) {
	light := theme.SurfaceShadows(false)
	if len(light) == 0 {
		t.Fatal("light surfaces have no shadow")
	}
	// The table is shared, so its capacity must not invite an append into the
	// package's own array.
	if cap(light) != len(light) {
		t.Fatalf("surface shadows expose spare capacity: len %d, cap %d", len(light), cap(light))
	}
	_ = append(light, theme.Shadow{})
	if got := len(theme.SurfaceShadows(false)); got != 3 {
		t.Fatalf("appending to the returned shadows changed the table: %d layers", got)
	}
	if theme.SurfaceShadows(true) != nil {
		t.Error("dark surfaces have a drop shadow")
	}
	dark := theme.OverlayShadows(true)
	if len(dark) != 1 || !dark[0].Inset {
		t.Fatalf("dark overlay shadows = %+v, want one inset hairline", dark)
	}
	if field := theme.FieldShadows(false); len(field) != len(theme.SurfaceShadows(false)) {
		t.Error("light field elevation is not the surface elevation")
	}
}
