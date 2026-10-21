package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSliderBoundsKeyboardVerticalRTLAndReadOnly(t *testing.T) {
	for _, tc := range []struct {
		name     string
		vertical bool
		rtl      bool
		key      ui.Key
	}{
		{name: "horizontal", key: ui.KeyRight},
		{name: "rtl", rtl: true, key: ui.KeyLeft},
		{name: "vertical", vertical: true, key: ui.KeyUp},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value := 5.0
			readonly, disabled := false, false
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				Slider(c, &value, SliderOptions{
					ControlOptions: ControlOptions{Label: "Level", ReadOnly: readonly, Disabled: disabled},
					Min:            0, Max: 10, Step: 2, Vertical: tc.vertical, Unit: "dB",
				})
			}, 300, 220)
			if tc.rtl {
				tt.SetDirection(ui.RightToLeft)
			}
			tt.Key(0, ui.KeyTab)
			tt.Key(0, tc.key)
			if value != 8 {
				t.Fatalf("native keyboard stepping: %v", value)
			}
			tt.Key(0, ui.KeyEnd)
			if value != 10 {
				t.Fatalf("native End bounds: %v", value)
			}
			tt.Key(0, tc.key)
			if value != 10 {
				t.Fatal("slider exceeded maximum")
			}
			readonly = true
			tt.Frame()
			tt.Key(0, ui.KeyHome)
			if value != 10 {
				t.Fatal("readonly slider changed caller value")
			}
			readonly, disabled = false, true
			tt.Frame()
			tt.Key(0, ui.KeyHome)
			if value != 10 {
				t.Fatal("disabled slider changed caller value")
			}
		})
	}
}

func TestRangeSliderThumbBoundsAndIndependentFocus(t *testing.T) {
	low, high := 20.0, 60.0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		RangeSlider(c, &low, &high, RangeSliderOptions{
			ControlOptions: ControlOptions{Label: "Price"}, Min: 0, Max: 100, Step: 10, Unit: "USD",
		})
	}, 300, 90)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if low != 60 || high != 60 {
		t.Fatalf("minimum thumb crossed maximum: %v..%v", low, high)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if low != 60 || high != 100 {
		t.Fatalf("Tab did not independently focus maximum thumb: %v..%v", low, high)
	}
	tt.Key(0, ui.KeyHome)
	if high != low {
		t.Fatalf("maximum thumb crossed minimum: %v..%v", low, high)
	}
}
