package conversation

import (
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestBubbleInteractiveVariantsKeepSemanticPaintAndContrast(t *testing.T) {
	for _, preset := range []theme.Preset{theme.Blue, theme.Neutral} {
		for _, dark := range []bool{false, true} {
			for _, variant := range []BubbleVariant{BubbleDefault, BubbleSecondary, BubbleMuted, BubbleTinted,
				BubbleOutline, BubbleGhost, BubbleDestructive} {
				t.Run(fmt.Sprintf("preset=%d/dark=%t/variant=%d", preset, dark, variant), func(t *testing.T) {
					var colors theme.ColorTokens
					clicks := 0
					tt := ui.NewTester(func(c *ui.Context) {
						theme.ApplyPreset(c, preset)
						colors = theme.Colors(c)
						ui.Column(c).Fill().Padding(40).Background(colors.Background).Children(func() {
							Bubble(c, BubbleOptions{Variant: variant, Label: "Interactive bubble", OnClick: func() { clicks++ }}, func() {
								ui.Text(c, "Readable message text").Label("Bubble ink")
								ui.Box(c).Height(12)
							})
						})
					}, 390, 240)
					tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
					tt.SetDark(dark)
					r := messageRect(t, tt, "Interactive bubble")
					ink := messageRect(t, tt, "Bubble ink")
					background, hover, foreground := colors.Primary, colors.PrimaryHover, colors.PrimaryForeground
					switch variant {
					case BubbleSecondary:
						background, hover, foreground = colors.Secondary, colors.SecondaryHover, colors.SecondaryForeground
					case BubbleMuted:
						background, foreground = colors.Muted, colors.Foreground
						hover = theme.Over(colors.Foreground, background.Over(colors.Background), 0.05)
					case BubbleTinted:
						background, hover, foreground = colors.Selection, colors.Selection, colors.SelectionForeground
					case BubbleOutline:
						background, hover, foreground = colors.Background, colors.Muted, colors.Foreground
					case BubbleGhost:
						background, hover, foreground = ui.Color{}, colors.Muted, colors.Foreground
					case BubbleDestructive:
						background, hover, foreground = colors.ErrorSurface, colors.ErrorSurface, colors.Destructive
					}
					assertBubblePaint(t, tt, r, ink, background.Over(colors.Background), foreground)
					tt.Move(r.X+r.W/2, r.Y+r.H/2)
					assertBubblePaint(t, tt, r, ink, hover.Over(colors.Background), foreground)
					if err := tt.Click("Interactive bubble"); err != nil {
						t.Fatal(err)
					}
					if clicks != 1 {
						t.Fatal("semantic surface lost native activation")
					}
					if messageRect(t, tt, "Interactive bubble") != r {
						t.Fatal("hover or activation changed the bubble geometry")
					}
				})
			}
		}
	}
}

// Inspect actual surface and fully covered glyph pixels. This catches replacing
// primary hover with a neutral background while preserving primary text ink.
func assertBubblePaint(t *testing.T, tt *ui.Tester, surface, text ui.Rect, expected, ink ui.Color) {
	t.Helper()
	frame := tt.Image()
	background := frame.RGBAAt(int(surface.X+surface.W/2), int(surface.Y+surface.H-5))
	if !bubbleColorNear(background, expected) {
		t.Fatalf("surface paint=%v want semantic color=%+v", background, expected)
	}
	foundInk := false
	for y := int(math.Ceil(float64(text.Y))); y < int(text.Y+text.H); y++ {
		for x := int(math.Ceil(float64(text.X))); x < int(text.X+text.W); x++ {
			if bubbleColorNear(frame.RGBAAt(x, y), ink) {
				foundInk = true
			}
		}
	}
	if !foundInk {
		t.Fatalf("message text did not retain semantic foreground %+v", ink)
	}
	if ratio := bubblePaintContrast(ink, ui.RGB(background.R, background.G, background.B)); ratio < 4.5 {
		t.Fatalf("interactive bubble contrast=%g:1, want at least 4.5:1", ratio)
	}
}

func bubbleColorNear(got color.RGBA, want ui.Color) bool {
	return absBubbleChannel(got.R, want.R) <= 2 && absBubbleChannel(got.G, want.G) <= 2 &&
		absBubbleChannel(got.B, want.B) <= 2 && got.A == want.A
}

func absBubbleChannel(a, b uint8) int { return int(max(a, b)) - int(min(a, b)) }

func bubblePaintContrast(a, b ui.Color) float64 {
	luminance := func(c ui.Color) float64 {
		channel := func(v uint8) float64 {
			x := float64(v) / 255
			if x <= 0.04045 {
				return x / 12.92
			}
			return math.Pow((x+0.055)/1.055, 2.4)
		}
		return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
	}
	x, y := luminance(a), luminance(b)
	return (max(x, y) + 0.05) / (min(x, y) + 0.05)
}
