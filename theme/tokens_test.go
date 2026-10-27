package theme_test

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestHeadingRolesCarryTrackingTight pins the native Coss type scale: text-4xl to
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
		if !style.Heading {
			t.Errorf("%s does not select the heading font", tc.name)
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
		if style.Heading {
			t.Errorf("%s selects the heading font", tc.name)
		}
	}
	if !theme.TypeCode.Style().Mono {
		t.Error("code role is not monospaced")
	}
}

func TestRadiusScaleUsesOffsets(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	for _, tc := range []struct {
		name string
		base float32
		want theme.Radii
	}{
		{
			name: "default",
			want: theme.Radii{XS: 4, SM: 6, MD: 8, LG: 10, XL: 14, XXL: 16, XXXL: 24, XXXXL: 32, Field: 8},
		},
		{
			name: "custom", base: 16,
			want: theme.Radii{XS: 10, SM: 12, MD: 14, LG: 16, XL: 20, XXL: 22, XXXL: 30, XXXXL: 38, Field: 14},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			theme.Configure(theme.Preferences{Radius: tc.base})
			if got := theme.RadiusScale(); got != tc.want {
				t.Errorf("radii = %+v, want %+v", got, tc.want)
			}
			if got := theme.RadiusBase(); got != tc.want.LG {
				t.Errorf("radius base = %v, want %v", got, tc.want.LG)
			}
		})
	}
}

func TestFieldBorderFollowsInputAndRing(t *testing.T) {
	defer theme.Configure(theme.Preferences{})
	input, ring := ui.RGBA(120, 120, 130, 0.4), ui.RGB(20, 30, 140)
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{
		Light: &theme.PalettePatch{Input: &input, Ring: &ring},
	}})
	got := theme.ColorsForMode(false)
	if got.FieldBorder != input || got.FieldBorderHover != input {
		t.Errorf("field border/hover = %v/%v, want the unflattened input token %v", got.FieldBorder, got.FieldBorderHover, input)
	}
	if got.FieldBorderFocus != ring || got.Focus != ring {
		t.Errorf("field focus/focus = %v/%v, want the independent ring token %v", got.FieldBorderFocus, got.Focus, ring)
	}
}

func TestShadowSnapshotsDoNotLeakCallerMutation(t *testing.T) {
	for _, tc := range []struct {
		name string
		get  func(bool) []theme.Shadow
	}{
		{name: "surface", get: theme.SurfaceShadows},
		{name: "overlay", get: theme.OverlayShadows},
		{name: "field", get: theme.FieldShadows},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, mode := range []struct {
				name string
				dark bool
			}{
				{name: "light"},
				{name: "dark", dark: true},
			} {
				t.Run(mode.name, func(t *testing.T) {
					shadows := tc.get(mode.dark)
					if len(shadows) == 0 {
						t.Fatal("elevation has no shadow or edge layers")
					}
					want := slices.Clone(shadows)
					shadows[0] = theme.Shadow{Blur: 99, Color: ui.RGB(255, 0, 255)}
					if got := tc.get(mode.dark); !slices.Equal(got, want) {
						t.Errorf("caller mutation changed published elevation: %+v, want %+v", got, want)
					}
				})
			}
		})
	}
}

func TestCossElevationKeepsCrispEdges(t *testing.T) {
	for _, tc := range []struct {
		name string
		get  func(bool) []theme.Shadow
	}{
		{name: "surface", get: theme.SurfaceShadows},
		{name: "field", get: theme.FieldShadows},
	} {
		t.Run(tc.name, func(t *testing.T) {
			layers := tc.get(true)
			if len(layers) != 1 || !layers[0].Inset || layers[0].Color.A == 0 {
				t.Errorf("dark elevation = %+v, want a visible inset edge", layers)
			}
		})
	}
	if layers := theme.OverlayShadows(true); len(layers) == 0 || layers[0].Inset {
		t.Errorf("dark popup elevation = %+v, want a drop shadow", layers)
	}
}
