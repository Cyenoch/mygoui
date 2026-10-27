package form

import (
	"image"
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestFieldNativeHelpTypographyAndFocus(t *testing.T) {
	for _, tc := range []struct {
		name  string
		scale float32
	}{
		{name: "default", scale: 1},
		{name: "larger text", scale: 1.5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var value, other string
			var inputHandle ui.Handle
			var inputFocused bool
			var originalFont, restoredFont float32
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				originalFont = c.Theme().FontSize
				ui.Column(c).Fill().Padding(20).Gap(10).Children(func() {
					Field(c, FieldOptions{
						Label: "Account", Description: "Helpful description", Error: "Invalid input",
					}, func() {
						ui.TextInputBase(c, &value).Width(220).Height(28).Bind(&inputHandle)
					})
					inputFocused = inputHandle.Focused(c)
					Field(c, FieldOptions{Label: "Second", Description: "Another description"}, func() {
						ui.TextInputBase(c, &other).Width(220).Height(28)
					})
					restoredFont = c.Theme().FontSize
					ui.Text(c, "Neighbor text")
				})
			}, 400, 350)
			tt.SetPreferences(ui.Preferences{TextScale: tc.scale})
			if originalFont != restoredFont {
				t.Fatalf("field changed theme font size from %v to %v", originalFont, restoredFont)
			}

			for _, text := range []struct {
				value  string
				size   float32
				weight int
				color  func(*ui.Theme) ui.Color
			}{
				{"Account", 14 * tc.scale, theme.WeightMedium, func(t *ui.Theme) ui.Color { return t.Text }},
				{"Helpful description", 12 * tc.scale, theme.WeightNormal, func(t *ui.Theme) ui.Color { return t.TextMuted }},
				{"Invalid input", 12 * tc.scale, theme.WeightNormal, func(t *ui.Theme) ui.Color { return t.Danger }},
				{"Another description", 12 * tc.scale, theme.WeightNormal, func(t *ui.Theme) ui.Color { return t.TextMuted }},
				{"Neighbor text", originalFont, theme.WeightNormal, func(t *ui.Theme) ui.Color { return t.Text }},
			} {
				ref := ui.NewTester(func(c *ui.Context) {
					theme.ApplyMode(c, false)
					ui.Column(c).Fill().Padding(20).Children(func() {
						ui.Text(c, text.value).FontSize(text.size).FontWeight(text.weight).TextColor(text.color(c.Theme()))
					})
				}, 400, 100)
				ref.SetPreferences(ui.Preferences{TextScale: tc.scale})
				actual, actualOK := tt.Find(text.value)
				expected, expectedOK := ref.Find(text.value)
				if !actualOK || !expectedOK {
					t.Fatalf("missing text %q", text.value)
				}
				if math.Abs(float64(actual.W-expected.W)) > 0.5 || math.Abs(float64(actual.H-expected.H)) > 0.5 {
					t.Fatalf("%q bounds %+v, want size %v x %v", text.value, actual, expected.W, expected.H)
				}
				gotInk := fieldTextInk(tt.Image(), actual)
				wantInk := fieldTextInk(ref.Image(), expected)
				if wantInk == 0 || math.Abs(gotInk-wantInk) > wantInk*0.08 {
					t.Fatalf("%q ink %.0f differs from expected regular/medium face %.0f", text.value, gotInk, wantInk)
				}
			}

			if err := tt.Click("Account"); err != nil {
				t.Fatal(err)
			}
			if !inputFocused {
				t.Fatal("clicking the field label did not focus its native editor")
			}
			tt.Type("ok")
			if value != "ok" {
				t.Fatalf("focused editor received %q", value)
			}
		})
	}
}

func fieldTextInk(img *image.RGBA, r ui.Rect) float64 {
	var ink float64
	for y := int(r.Y); y < int(math.Ceil(float64(r.Y+r.H))); y++ {
		for x := int(r.X); x < int(math.Ceil(float64(r.X+r.W))); x++ {
			p := img.RGBAAt(x, y)
			ink += float64(255 - int(p.R) + 255 - int(p.G) + 255 - int(p.B))
		}
	}
	return ink
}
