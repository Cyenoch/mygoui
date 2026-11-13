package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func rangeSliderTester(low, high *float64, o RangeSliderOptions) *ui.Tester {
	o.Label = "Price"
	return ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Padding(40).Children(func() { RangeSlider(c, low, high, o) })
	}, 400, 160)
}

func TestRangeSliderKeyboardMovesFocusedThumb(t *testing.T) {
	low, high := 20.0, 60.0
	changes := 0
	tt := rangeSliderTester(&low, &high, RangeSliderOptions{
		Max: 100, Step: 10, OnChange: func() { changes++ },
	})
	for _, name := range []string{"Price minimum", "Price maximum"} {
		if _, ok := tt.AccessByRole(ui.RoleSlider, name); !ok {
			t.Fatalf("no slider named %q", name)
		}
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if low != 30 || high != 60 || changes != 1 {
		t.Fatalf("Right on the low thumb gave %v..%v with %d changes", low, high, changes)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if low != 30 || high != 100 {
		t.Fatalf("End on the high thumb gave %v..%v", low, high)
	}
	tt.Key(0, ui.KeyHome)
	if low != 30 || high != 30 {
		t.Fatalf("Home on the high thumb gave %v..%v, want it stopped at the low thumb", low, high)
	}
}

func TestRangeSliderMinGapStopsEitherThumb(t *testing.T) {
	low, high := 20.0, 60.0
	changes := 0
	tt := rangeSliderTester(&low, &high, RangeSliderOptions{
		Max: 100, Step: 10, MinGap: 20, OnChange: func() { changes++ },
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if low != 40 || high != 60 {
		t.Fatalf("End on the low thumb gave %v..%v, want 40..60", low, high)
	}
	changes = 0
	tt.Key(0, ui.KeyRight)
	if low != 40 || changes != 0 {
		t.Fatalf("a thumb held at the gap moved to %v or reported %d changes", low, changes)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if low != 40 || high != 60 {
		t.Fatalf("Home on the high thumb gave %v..%v, want 40..60", low, high)
	}
}

func TestRangeSliderDisabledIgnoresKeys(t *testing.T) {
	low, high := 20.0, 60.0
	tt := rangeSliderTester(&low, &high, RangeSliderOptions{Max: 100, Disabled: true})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if low != 20 || high != 60 {
		t.Fatalf("a disabled range slider moved to %v..%v", low, high)
	}
}
