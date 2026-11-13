package color

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestColorAreaArrowsStepSaturationAndBrightness checks the area's keys:
// arrows step one percent, Shift ten, and Home and End reach the ends of
// saturation, each reported once.
func TestColorAreaArrowsStepSaturationAndBrightness(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	s.SetSaturation(0.5)
	s.SetBrightness(0.5)
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Before")
		ColorArea(c, ColorAreaOptions{Label: "Fill", State: s, OnChange: func() { changes++ }})
	}, 320, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Fill") {
		t.Fatal("Tab did not reach the area")
	}
	steps := []struct {
		mods       ui.Modifiers
		key        ui.Key
		sat, value float64
	}{
		{0, ui.KeyRight, 0.51, 0.5},
		{0, ui.KeyDown, 0.51, 0.49},
		{ui.Shift, ui.KeyUp, 0.51, 0.59},
		{ui.Shift, ui.KeyLeft, 0.41, 0.59},
		{0, ui.KeyHome, 0, 0.59},
		{0, ui.KeyEnd, 1, 0.59},
	}
	for i, st := range steps {
		tt.Key(st.mods, st.key)
		if !near(s.Saturation(), st.sat, 0.001) || !near(s.Brightness(), st.value, 0.001) {
			t.Fatalf("after %v: saturation %.3f brightness %.3f, want %.2f %.2f",
				st.key, s.Saturation(), s.Brightness(), st.sat, st.value)
		}
		if changes != i+1 {
			t.Fatalf("after %v OnChange ran %d times, want %d", st.key, changes, i+1)
		}
	}
}

// TestColorSliderKeysStepTheChannel checks a horizontal slider's keys: one
// unit per arrow, ten with Shift or a page key, and the ends with Home and End.
func TestColorSliderKeysStepTheChannel(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	s.SetHue(180)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Before")
		ColorSlider(c, ColorSliderOptions{Label: "Hue", State: s, Channel: ChannelHue})
	}, 320, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	steps := []struct {
		mods ui.Modifiers
		key  ui.Key
		hue  float64
	}{
		{0, ui.KeyRight, 181},
		{0, ui.KeyLeft, 180},
		{ui.Shift, ui.KeyRight, 190},
		{0, ui.KeyPageDown, 180},
		{0, ui.KeyPageUp, 190},
		{0, ui.KeyEnd, 360},
		{0, ui.KeyHome, 0},
	}
	for _, st := range steps {
		tt.Key(st.mods, st.key)
		if !near(s.Hue(), st.hue, 0.001) {
			t.Fatalf("after %v the hue is %v, want %v", st.key, s.Hue(), st.hue)
		}
	}
}
