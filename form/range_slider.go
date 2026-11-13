package form

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// RangeSliderOptions configures a two-thumb slider. Min, Max, Step, Marks and
// Format behave as in SliderOptions. MinGap keeps the thumbs at least that
// far apart; it should be a multiple of Step so a stopped thumb stays on the
// grid.
type RangeSliderOptions struct {
	Label       string
	Description string
	Error       string
	Min         float64
	Max         float64
	Step        float64
	MinGap      float64
	Marks       []SliderMark
	Disabled    bool
	Invalid     bool
	HideOutput  bool
	Format      func(float64) string
	OnChange    func()
}

// RangeSlider is the horizontal two-thumb Slider over caller-owned low and
// high values. Each thumb is a Tab stop moved by drags, the arrows, Home and
// End; the thumbs never cross, and a thumb moved inside MinGap stops at the
// gap. A nil low or high renders the slider disabled.
func RangeSlider(c *ui.Context, low, high *float64, o RangeSliderOptions) ui.Element {
	if low == nil || high == nil {
		var lo, hi float64
		low, high = &lo, &hi
		o.Disabled = true
	}
	wasLow, wasHigh := *low, *high
	return Slider(c, SliderOptions{
		State:       SliderState{Low: low, High: high},
		Label:       o.Label,
		Description: o.Description,
		Error:       o.Error,
		Min:         o.Min,
		Max:         o.Max,
		Step:        o.Step,
		Marks:       o.Marks,
		Disabled:    o.Disabled,
		Invalid:     o.Invalid,
		HideOutput:  o.HideOutput,
		Format:      o.Format,
		OnChange: func() {
			// A thumb that entered the gap stops at it; one already
			// inside may still move away.
			if o.MinGap > 0 {
				*low = math.Min(*low, math.Max(wasLow, *high-o.MinGap))
				*high = math.Max(*high, math.Min(wasHigh, *low+o.MinGap))
			}
			if (*low != wasLow || *high != wasHigh) && o.OnChange != nil {
				o.OnChange()
			}
		},
	})
}
