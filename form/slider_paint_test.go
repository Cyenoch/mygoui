package form

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSliderThumbRemainsDistinctAcrossSemanticThemes(t *testing.T) {
	for _, preset := range []theme.Preset{theme.Blue, theme.Neutral} {
		for _, dark := range []bool{false, true} {
			for _, contrast := range []bool{false, true} {
				for _, variant := range []struct {
					name                  string
					rangeSlider, vertical bool
					rtl, disabled         bool
				}{
					{name: "single"}, {name: "single rtl", rtl: true},
					{name: "vertical", vertical: true},
					{name: "single disabled", disabled: true},
					{name: "range", rangeSlider: true},
					{name: "range rtl", rangeSlider: true, rtl: true},
					{name: "range disabled", rangeSlider: true, disabled: true},
				} {
					name := fmt.Sprintf("preset=%d/dark=%v/contrast=%v/%s", preset, dark, contrast, variant.name)
					t.Run(name, func(t *testing.T) {
						value, low, high := 50.0, 25.0, 75.0
						var bounds ui.Rect
						var palette theme.ColorTokens
						var inset float32
						tt := ui.NewTester(func(c *ui.Context) {
							if contrast {
								theme.ApplyHighContrast(c, preset)
							} else {
								theme.ApplyPreset(c, preset)
							}
							palette = theme.Colors(c)
							inset = c.Theme().Space(5)/2 + 1
							ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
								control := ControlOptions{Label: "Level", Disabled: variant.disabled}
								var e ui.Element
								if variant.rangeSlider {
									e = RangeSlider(c, &low, &high, RangeSliderOptions{ControlOptions: control, Max: 100})
								} else {
									e = Slider(c, &value, SliderOptions{
										ControlOptions: control, Max: 100, Vertical: variant.vertical,
									})
								}
								e.DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
							})
						}, 320, 230)
						tt.SetDark(dark)
						if variant.rtl {
							tt.SetDirection(ui.RightToLeft)
						}
						tt.Frame()
						center := ui.Point{X: bounds.X + bounds.W/2, Y: bounds.Y + bounds.H/2}
						if variant.rangeSlider {
							center.X = bounds.X + inset + (bounds.W-2*inset)*0.25
							if variant.rtl {
								center.X = bounds.X + inset + (bounds.W-2*inset)*0.75
							}
						}
						thumb := tt.Image().RGBAAt(int(center.X), int(center.Y))
						filled, empty := center, center
						switch {
						case variant.rangeSlider:
							filled.X = bounds.X + bounds.W/2
							empty.X = bounds.X + inset + 12
						case variant.vertical:
							filled.Y += 20
							empty.Y -= 20
						case variant.rtl:
							filled.X += 20
							empty.X -= 20
						default:
							filled.X -= 20
							empty.X += 20
						}
						activeTrack := tt.Image().RGBAAt(int(filled.X), int(filled.Y))
						emptyTrack := tt.Image().RGBAAt(int(empty.X), int(empty.Y))
						want := palette.Background
						if dark {
							want = palette.Foreground
						}
						// Disabled paint composites each operation over existing track
						// paint; its center is not a raw token blended over the root.
						if !variant.disabled && (absChannel(thumb.R, want.R) > 3 ||
							absChannel(thumb.G, want.G) > 3 || absChannel(thumb.B, want.B) > 3) {
							t.Fatalf("painted thumb %v does not match semantic surface %v", thumb, want)
						}
						minimum := 3.0
						if variant.disabled {
							minimum = 1.4
						}
						if ratio := sliderPixelContrast(thumb, emptyTrack); ratio < minimum && dark {
							t.Fatalf("dark thumb disappears against unfilled track: %.2f:1", ratio)
						}
						// In light themes the foreground outline distinguishes the light
						// thumb from the light unfilled track, including high contrast.
						// The circular silhouette extends above and below the track.
						outlineContrast := sliderPixelContrast(thumb, emptyTrack)
						// Rasterized round borders straddle subpixels; inspect the
						// visible upper arc rather than assuming its row is y-7.
						for y := int(center.Y - 9); y < int(center.Y-3); y++ {
							for x := int(center.X - 7); x <= int(center.X+7); x++ {
								outline := tt.Image().RGBAAt(x, y)
								outlineContrast = max(outlineContrast, sliderPixelContrast(outline, emptyTrack))
							}
						}
						if outlineContrast < minimum {
							t.Fatalf("thumb surface and outline disappear against unfilled track: %.2f:1", outlineContrast)
						}
						if activeTrack == emptyTrack {
							t.Fatal("active track disappeared into the unfilled track")
						}
						if !variant.disabled {
							before := bytes.Clone(tt.Image().Pix)
							tt.Key(0, ui.KeyTab)
							if bytes.Equal(before, tt.Image().Pix) {
								t.Fatal("keyboard focus produced no visible ring")
							}
						}
					})
				}
			}
		}
	}
}

func TestSliderPointerMatchesPaintedTrackAcrossDirections(t *testing.T) {
	for _, tc := range []struct {
		name          string
		vertical, rtl bool
	}{
		{name: "horizontal"}, {name: "rtl", rtl: true}, {name: "vertical", vertical: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value := 50.0
			var bounds ui.Rect
			changes := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
					Slider(c, &value, SliderOptions{
						ControlOptions: ControlOptions{Label: "Level"}, Max: 100,
						Vertical: tc.vertical, OnChange: func(float64) { changes++ },
					}).DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
				})
			}, 320, 230)
			if tc.rtl {
				tt.SetDirection(ui.RightToLeft)
			}
			for _, target := range []float32{0, 0.3, 1} {
				point := ui.Point{X: bounds.X + bounds.W/2, Y: bounds.Y + bounds.H/2}
				switch {
				case tc.vertical:
					point.Y = bounds.Y + 9 + (bounds.H-18)*(1-target)
				case tc.rtl:
					point.X = bounds.X + 9 + (bounds.W-18)*(1-target)
				default:
					point.X = bounds.X + 9 + (bounds.W-18)*target
				}
				tt.ClickAt(point.X, point.Y)
				if math.Abs(value-float64(target)*100) > 0.1 {
					t.Fatalf("painted %.0f%% point selected %.2f", target*100, value)
				}
			}
			if changes != 3 {
				t.Fatalf("three pointer edits emitted %d callbacks", changes)
			}
		})
	}
}

func TestRangeSliderPaintAndInputShareRTLThumbPositions(t *testing.T) {
	for _, rtl := range []bool{false, true} {
		t.Run(fmt.Sprintf("rtl=%v", rtl), func(t *testing.T) {
			low, high := 20.0, 80.0
			readonly, disabled := false, false
			changes := 0
			var bounds ui.Rect
			var inset float32
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				inset = c.Theme().Space(5)/2 + 1
				ui.Column(c).Padding(20).Children(func() {
					RangeSlider(c, &low, &high, RangeSliderOptions{
						ControlOptions: ControlOptions{Label: "Price", ReadOnly: readonly, Disabled: disabled},
						Max:            100, OnChange: func(float64, float64) { changes++ },
					}).DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
				})
			}, 320, 90)
			if rtl {
				tt.SetDirection(ui.RightToLeft)
			}
			position := func(fraction float32) ui.Point {
				if rtl {
					fraction = 1 - fraction
				}
				return ui.Point{X: bounds.X + inset + (bounds.W-2*inset)*fraction, Y: bounds.Y + bounds.H/2}
			}
			start, target := position(0.2), position(0.35)
			tt.Move(start.X, start.Y)
			tt.Press(start.X, start.Y)
			tt.Move(target.X, target.Y)
			tt.Release(target.X, target.Y)
			if math.Abs(low-35) > 1 || high != 80 || changes == 0 {
				t.Fatalf("dragging painted low thumb selected wrong native value: %.2f..%.2f changes=%d", low, high, changes)
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyEnd)
			if math.Abs(low-35) > 1 || high != 100 {
				t.Fatalf("Tab did not independently route End to high thumb: %.2f..%.2f", low, high)
			}
			tt.Key(ui.Shift, ui.KeyTab)
			tt.Key(0, ui.KeyHome)
			if low != 0 || high != 100 {
				t.Fatalf("Shift-Tab did not independently route Home to low thumb: %.2f..%.2f", low, high)
			}
			readonly = true
			tt.Frame()
			previousLow, previousHigh, previousChanges := low, high, changes
			tt.Key(0, ui.KeyEnd)
			point := position(0.95)
			tt.ClickAt(point.X, point.Y)
			if low != previousLow || high != previousHigh || changes != previousChanges {
				t.Fatal("read-only range slider changed caller state or callback count")
			}
			readonly, disabled = false, true
			tt.Frame()
			tt.Key(0, ui.KeyHome)
			tt.ClickAt(point.X, point.Y)
			if low != previousLow || high != previousHigh || changes != previousChanges {
				t.Fatal("disabled range slider changed caller state or callback count")
			}
		})
	}
}

func sliderPixelContrast(a, b color.RGBA) float64 {
	luminance := func(c color.RGBA) float64 {
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
