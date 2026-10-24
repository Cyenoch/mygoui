package gallery

import (
	"strings"
	"testing"

	"github.com/Cyenoch/mygoui/color"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestAPILongValuesWrapInsteadOfClipping pins that a type or default wider
// than its column wraps onto more lines inside the pill instead of clipping at
// a fixed height.
func TestAPILongValuesWrapInsteadOfClipping(t *testing.T) {
	const longType = "[]navigation.AccordionItemOptions"
	const longDefault = "mygoui.ScrollShadowOrientation"
	params := []ComponentParam{
		{"LongOption", longType, longDefault, "A value longer than its column."},
	}
	tt := ui.NewTester(func(c *ui.Context) { apiTable(c, params) }, 1200, 900)
	tt.Frame()
	for _, value := range []string{longType, longDefault} {
		r, ok := tt.Find(value)
		if !ok {
			t.Fatalf("%q is missing from the API table", value)
		}
		if r.W <= 0 || r.H <= 0 {
			t.Fatalf("%q is clipped (rect %+v)", value, r)
		}
		if r.X < 0 || r.X+r.W > 1200 {
			t.Fatalf("%q overflows the page (rect %+v)", value, r)
		}
		if r.H <= 18 {
			t.Fatalf("%q did not wrap onto more lines (rect %+v)", value, r)
		}
	}
}

// TestAPILongValuesWrapInStackedLayout is the narrow counterpart: a value
// wider than the param row wraps instead of running past the pane.
func TestAPILongValuesWrapInStackedLayout(t *testing.T) {
	const longType = "[]navigation.AccordionItemOptions"
	params := []ComponentParam{
		{"LongOption", longType, "nil", "A value longer than its column."},
	}
	tt := ui.NewTester(func(c *ui.Context) { apiTable(c, params) }, 500, 900)
	tt.Frame()
	r, ok := tt.Find(longType)
	if !ok || r.W <= 0 || r.H <= 0 || r.X < 0 || r.X+r.W > 500 {
		t.Fatalf("%q is clipped or overflows the stacked page (rect %+v)", longType, r)
	}
}

// TestThemeTokenCardsShowValuesAndRoles pins the token section: one card per
// token, each naming the token in mono, resolving its light and dark hex
// values from the palette, labeling the L/D swatches and explaining its role.
func TestThemeTokenCardsShowValuesAndRoles(t *testing.T) {
	tokens := []string{"Accent", "AccentForeground", "Default", "DefaultHover", "Danger"}
	light := theme.ColorsForMode(false)
	dark := theme.ColorsForMode(true)
	for _, size := range []struct{ width, height int }{{1100, 900}, {640, 900}} {
		state := New()
		state.Selected = "button"
		var page float32
		tt := ui.NewTester(func(c *ui.Context) {
			page = pageWidth(c)
			View(c, state)
		}, size.width, size.height)
		tt.Frame()

		cards := 0
		for _, text := range tt.Texts() {
			if strings.HasPrefix(text, "token card ") {
				cards++
			}
		}
		if cards != len(tokens) {
			t.Fatalf("%dx%d: rendered %d token cards, want %d", size.width, size.height, cards, len(tokens))
		}

		for _, token := range tokens {
			r := revealContent(t, tt, state, "token card "+token, float32(size.width), float32(size.height))
			if r.W <= 0 || r.H <= 0 {
				t.Fatalf("%dx%d: token %q card is empty (rect %+v)", size.width, size.height, token, r)
			}
			if r.X < 0 || r.X+r.W > float32(size.width) {
				t.Fatalf("%dx%d: token %q card overflows the page (rect %+v)", size.width, size.height, token, r)
			}
			if page >= 600 {
				if r.W > page/2 {
					t.Fatalf("%dx%d: token %q card is wider than half the page (%+v of %v)",
						size.width, size.height, token, r, page)
				}
			} else if r.W < page-4 {
				t.Fatalf("%dx%d: token %q card does not fill the narrow page (%+v of %v)",
					size.width, size.height, token, r, page)
			}
			if !tt.HasText(token) {
				t.Fatalf("%dx%d: token %q card lost its name", size.width, size.height, token)
			}
			lc, ok := paletteColor(&light, token)
			if !ok {
				t.Fatalf("%dx%d: token %q does not resolve in the light palette", size.width, size.height, token)
			}
			dc, ok := paletteColor(&dark, token)
			if !ok {
				t.Fatalf("%dx%d: token %q does not resolve in the dark palette", size.width, size.height, token)
			}
			if !tt.HasText(color.FormatHex(lc)) {
				t.Fatalf("%dx%d: token %q card lost its light hex %s", size.width, size.height, token, color.FormatHex(lc))
			}
			if !tt.HasText(color.FormatHex(dc)) {
				t.Fatalf("%dx%d: token %q card lost its dark hex %s", size.width, size.height, token, color.FormatHex(dc))
			}
		}
		if _, ok := tt.Find("L"); !ok {
			t.Fatalf("%dx%d: token cards lost the L swatch label", size.width, size.height)
		}
		if _, ok := tt.Find("D"); !ok {
			t.Fatalf("%dx%d: token cards lost the D swatch label", size.width, size.height)
		}
		first := revealContent(t, tt, state, "token card "+tokens[0], float32(size.width), float32(size.height))
		second := revealContent(t, tt, state, "token card "+tokens[1], float32(size.width), float32(size.height))
		if page >= 600 {
			if second.X <= first.X || second.Y != first.Y {
				t.Fatalf("%dx%d: token cards do not form two columns (first %+v, second %+v)",
					size.width, size.height, first, second)
			}
		} else if second.X != first.X {
			t.Fatalf("%dx%d: token cards do not stack in one column (first %+v, second %+v)",
				size.width, size.height, first, second)
		}
		if !tt.HasText("Primary action and focus emphasis.") {
			t.Fatalf("%dx%d: token cards lost the Accent role description", size.width, size.height)
		}
	}
}

// TestThemeTokenCardsMarkCustomTokens keeps non-palette tokens honest: they
// render a name, a role and "custom token" rather than a made-up color.
func TestThemeTokenCardsMarkCustomTokens(t *testing.T) {
	state := New()
	state.Selected = "card"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1100, 900)
	tt.Frame()
	light := theme.ColorsForMode(false)
	if _, ok := paletteColor(&light, "SurfaceShadows"); ok {
		t.Fatal("SurfaceShadows unexpectedly resolves to a palette color")
	}
	r := revealContent(t, tt, state, "token card SurfaceShadows", 1100, 900)
	if r.W <= 0 || r.H <= 0 {
		t.Fatalf("custom token has no card (rect %+v)", r)
	}
	if !tt.HasText("custom token") {
		t.Fatal("custom token card does not say custom token")
	}
	if !tt.HasText("Elevation shadow layers for cards and overlays.") {
		t.Fatal("custom token card does not explain its role")
	}
}
