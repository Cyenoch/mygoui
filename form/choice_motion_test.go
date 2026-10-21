package form

import (
	"bytes"
	"image/color"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSwitchPaintInterpolatesAndReversesFromCurrentPosition(t *testing.T) {
	for _, tc := range []struct {
		name       string
		small, rtl bool
	}{
		{name: "normal"}, {name: "small", small: true},
		{name: "rtl", rtl: true}, {name: "small rtl", small: true, rtl: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value := false
			var bounds ui.Rect
			var ink ui.Color
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ink = theme.Colors(c).Background
				ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
					Switch(c, &value, SwitchOptions{Small: tc.small}).Label("Airplane mode").
						DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
				})
			}, 200, 90)
			if tc.rtl {
				tt.SetDirection(ui.RightToLeft)
			}
			tt.Advance(time.Second)
			start := switchThumbCenter(t, tt, bounds, ink)
			height := float32(20)
			if tc.small {
				height = 16
			}
			trackPixel := func() color.RGBA {
				return tt.Image().RGBAAt(int(bounds.X+bounds.W/2), int(bounds.Y+(bounds.H-height)/2+1))
			}
			offTrack := trackPixel()
			if err := tt.Click("Airplane mode"); err != nil {
				t.Fatal(err)
			}
			if !value {
				t.Fatal("pointer activation did not change the native value immediately")
			}
			tt.Advance(60 * time.Millisecond)
			middle := switchThumbCenter(t, tt, bounds, ink)
			middleTrack := trackPixel()
			if math.Abs(float64(middle-start)) < 1 {
				t.Fatal("switch thumb did not paint an intermediate position")
			}
			tt.Key(0, ui.KeySpace)
			if value {
				t.Fatal("Space did not reverse the native switch")
			}
			if reversed := switchThumbCenter(t, tt, bounds, ink); math.Abs(float64(reversed-middle)) > 1 {
				t.Fatalf("rapid reversal jumped from %.1f to %.1f", middle, reversed)
			}
			tt.Advance(60 * time.Millisecond)
			returning := switchThumbCenter(t, tt, bounds, ink)
			if math.Abs(float64(returning-start)) >= math.Abs(float64(middle-start)) {
				t.Fatal("reversed thumb did not move back toward its current target")
			}
			tt.Advance(time.Second)
			if got := switchThumbCenter(t, tt, bounds, ink); math.Abs(float64(got-start)) > 0.5 {
				t.Fatalf("switch did not settle back at its off position: %.1f != %.1f", got, start)
			}
			value = true
			tt.Frame()
			tt.Advance(time.Second)
			end := switchThumbCenter(t, tt, bounds, ink)
			if middle <= min(start, end) || middle >= max(start, end) {
				t.Fatalf("midframe %.1f must be between endpoints %.1f and %.1f", middle, start, end)
			}
			if middleTrack == offTrack || middleTrack == trackPixel() {
				t.Fatal("track color jumped instead of painting an intermediate color")
			}
			if tt.FramePending() {
				t.Fatal("settled switch continued scheduling animation frames")
			}
		})
	}
}

func TestSwitchStartsSelectedAndReducedMotionSnaps(t *testing.T) {
	value := true
	var bounds ui.Rect
	var ink ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ink = theme.Colors(c).Background
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			Switch(c, &value, SwitchOptions{}).Label("Selected switch").
				DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
		})
	}, 160, 90)
	selected := switchThumbCenter(t, tt, bounds, ink)
	first := bytes.Clone(tt.Image().Pix)
	tt.Advance(time.Second)
	if !bytes.Equal(first, tt.Image().Pix) || tt.FramePending() {
		t.Fatal("initially selected switch played a false entrance")
	}
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	if err := tt.Click("Selected switch"); err != nil {
		t.Fatal(err)
	}
	if value || switchThumbCenter(t, tt, bounds, ink) >= selected-5 {
		t.Fatal("reduced-motion switch did not snap to its new value")
	}
	reduced := bytes.Clone(tt.Image().Pix)
	tt.Advance(time.Second)
	if !bytes.Equal(reduced, tt.Image().Pix) || tt.FramePending() {
		t.Fatal("reduced-motion switch retained animated position or colors")
	}
}

func TestSwitchDarkThumbColorInterpolatesWithTrack(t *testing.T) {
	value := false
	var bounds ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			Switch(c, &value, SwitchOptions{}).Label("Dark switch").
				DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r })
		})
	}, 160, 90)
	tt.SetDark(true)
	tt.Advance(time.Second)
	// This point is inside the thumb at both endpoints and throughout its
	// travel, so a color change cannot be confused with movement over track.
	pixel := func() color.RGBA {
		return tt.Image().RGBAAt(int(bounds.X+16), int(bounds.Y+bounds.H/2))
	}
	off := pixel()
	value = true
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	middle := pixel()
	tt.Advance(time.Second)
	on := pixel()
	if off == on || middle == off || middle == on {
		t.Fatalf("dark thumb color must interpolate between distinct states: %v -> %v -> %v", off, middle, on)
	}
}

func TestSwitchLabelInputAndReadOnlyDoNotDuplicateCallbacks(t *testing.T) {
	value, readonly, disabled := false, false, false
	changes := []bool{}
	var bounds ui.Rect
	var ink ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ink = theme.Colors(c).Background
		Switch(c, &value, SwitchOptions{
			ControlOptions: ControlOptions{Label: "Airplane mode", ReadOnly: readonly, Disabled: disabled},
			OnChange:       func(on bool) { changes = append(changes, on) },
		}).DrawOver(func(_ *ui.Painter, r ui.Rect) { bounds = r; bounds.W = 32 })
	}, 260, 90)
	if err := tt.Click("Airplane mode"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time.Second)
	tt.Key(0, ui.KeySpace)
	tt.Advance(time.Second)
	if value || len(changes) != 2 || !changes[0] || changes[1] {
		t.Fatalf("pointer/Space callbacks must report each bool once: %v, value=%v", changes, value)
	}
	readonly = true
	tt.Frame()
	before := switchThumbCenter(t, tt, bounds, ink)
	tt.Click("Airplane mode")
	tt.Key(0, ui.KeySpace)
	tt.Advance(time.Second)
	if value || len(changes) != 2 || switchThumbCenter(t, tt, bounds, ink) != before {
		t.Fatal("read-only activation changed the value, thumb or callbacks")
	}
	readonly, disabled = false, true
	tt.Frame()
	tt.Click("Airplane mode")
	tt.Key(0, ui.KeySpace)
	if value || len(changes) != 2 {
		t.Fatal("disabled switch changed value or callbacks")
	}
}

func switchThumbCenter(t *testing.T, tt *ui.Tester, bounds ui.Rect, ink ui.Color) float32 {
	t.Helper()
	var sum, count int
	for x := int(bounds.X + 1); x < int(bounds.X+bounds.W-1); x++ {
		got := tt.Image().RGBAAt(x, int(bounds.Y+bounds.H/2))
		if absChannel(got.R, ink.R) <= 2 && absChannel(got.G, ink.G) <= 2 && absChannel(got.B, ink.B) <= 2 {
			sum += x
			count++
		}
	}
	if count == 0 {
		t.Fatal("no painted thumb found in switch track")
	}
	return float32(sum) / float32(count)
}

func absChannel(a, b uint8) int { return absInt(int(a) - int(b)) }

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
