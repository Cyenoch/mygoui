package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSliderDragAndKeyboard(t *testing.T) {
	value := 50.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).FillWidth().Children(func() {
			Slider(c, SliderOptions{
				State: SliderState{Value: &value},
				Label: "Volume",
				Min:   0,
				Max:   100,
				Step:  10,
				Marks: []SliderMark{{Value: 0, Label: "0"}, {Value: 50, Label: "50"}, {Value: 100, Label: "100"}},
			})
		})
	}, 360, 160)

	header, ok := tt.Find("Volume")
	if !ok {
		t.Fatal("the slider label is not rendered")
	}
	if !tt.HasText("50") {
		t.Fatalf("output missing: %q", tt.Texts())
	}
	const pad = float32(14)
	left, width := header.X, float32(360)-2*header.X
	y := header.Y + header.H + theme.Space4 + 10
	at := func(f float32) float32 { return left + pad + (width-2*pad)*f }
	// A press at 30% of the padded content box snaps to 30.
	tt.Press(at(0.3), y)
	tt.Release(at(0.3), y)
	if value != 30 {
		t.Fatalf("a press at 30%% set %v", value)
	}
	// Dragging is live: a move before release changes the value.
	tt.Press(at(0.3), y)
	tt.Move(at(0.7), y)
	if value != 70 {
		t.Fatalf("a drag to 70%% set %v", value)
	}
	tt.Release(at(0.7), y)
	// The arrows move by Step; Home and End go to the bounds.
	tt.Key(0, ui.KeyRight)
	if value != 80 {
		t.Fatalf("Right set %v, want 80", value)
	}
	tt.Key(0, ui.KeyHome)
	if value != 0 {
		t.Fatalf("Home set %v", value)
	}
	tt.Key(0, ui.KeyEnd)
	if value != 100 {
		t.Fatalf("End set %v", value)
	}
}

func TestSliderMarksSitInTrackRow(t *testing.T) {
	value := 50.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).FillWidth().Children(func() {
			Slider(c, SliderOptions{
				State: SliderState{Value: &value}, Label: "Volume",
				HideOutput: true, Min: 0, Max: 100,
				Marks: []SliderMark{
					{Value: 0, Label: "0"},
					{Value: 50, Label: "50"},
					{Value: 100, Label: "100"},
				},
			})
		})
	}, 360, 200)

	header, ok := tt.Find("Volume")
	if !ok {
		t.Fatal("the slider label is not rendered")
	}
	const trackH, markH, pad = float32(20), float32(16), float32(14)
	left, width := header.X, float32(360)-2*header.X
	markTop := header.Y + header.H + theme.Space4 + trackH + theme.Space4
	want := []float32{left + pad, left + width/2, left + width - pad}
	for i, s := range []string{"0", "50", "100"} {
		r, ok := tt.Find(s)
		if !ok {
			t.Fatalf("mark %q is not rendered", s)
		}
		// The label box lies within the marks row right below the track, not
		// at the top of the page.
		if r.Y < markTop-0.5 || r.Y+r.H > markTop+markH+0.5 {
			t.Fatalf("mark %q sits at y %v..%v, outside the marks row %v..%v", s, r.Y, r.Y+r.H, markTop, markTop+markH)
		}
		// Its center lines up with the thumb at the same value, inside the
		// track's padded end.
		if c := r.X + r.W/2; c < want[i]-1.5 || c > want[i]+1.5 {
			t.Fatalf("mark %q center is %.1f, want %.1f", s, c, want[i])
		}
	}
}

func TestSliderDisabled(t *testing.T) {
	value := 40.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).FillWidth().Children(func() {
			Slider(c, SliderOptions{
				State: SliderState{Value: &value}, Label: "Locked",
				Min: 0, Max: 100, Disabled: true,
			})
		})
	}, 320, 140)
	header, ok := tt.Find("Locked")
	if !ok {
		t.Fatal("the slider label is not rendered")
	}
	y := header.Y + header.H + theme.Space4 + 10
	tt.ClickAt(header.X+(float32(320)-2*header.X)/2, y)
	if value != 40 {
		t.Fatalf("a disabled slider set %v", value)
	}
}

func TestSliderRangeKeepsThumbsOrdered(t *testing.T) {
	low, high := 20.0, 80.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).FillWidth().Children(func() {
			Slider(c, SliderOptions{
				State: SliderState{Low: &low, High: &high},
				Label: "Price",
				Min:   0,
				Max:   100,
				Step:  10,
			})
		})
	}, 360, 140)

	header, ok := tt.Find("Price")
	if !ok {
		t.Fatal("the range labels are not rendered")
	}
	left, width := header.X, float32(360)-2*header.X
	y := header.Y + header.H + theme.Space4 + 10
	// Pressing at 60% moves the nearest knob, the high one.
	tt.ClickAt(left+10+(width-20)*0.6, y)
	if low != 20 || high != 60 {
		t.Fatalf("a click at 60%%: low %v high %v", low, high)
	}
	// The high knob has the focus; Right moves it by Step.
	tt.Key(0, ui.KeyRight)
	if high != 70 {
		t.Fatalf("Right on the high knob: %v", high)
	}
	if low > high {
		t.Fatalf("thumbs crossed: low %v high %v", low, high)
	}
	// Dragging the low knob past the high one stops at it.
	lowX := left + 10 + (width-20)*(float32(low)/100)
	tt.Press(lowX, y)
	tt.Move(left+width-10, y)
	tt.Release(left+width-10, y)
	if low > high {
		t.Fatalf("drag crossed the thumbs: low %v high %v", low, high)
	}
	if low != high {
		t.Fatalf("the low knob stopped at %v, high is %v", low, high)
	}
}

func TestSliderVertical(t *testing.T) {
	value := 25.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Center).Children(func() {
			Slider(c, SliderOptions{
				State: SliderState{Value: &value}, Label: "Level",
				Min: 0, Max: 100, Vertical: true, Height: 120,
			})
		})
	}, 240, 260)
	r, ok := tt.Find("Level")
	if !ok {
		t.Fatal("the vertical slider is not rendered")
	}
	if r.H <= r.W {
		t.Fatalf("vertical slider is %v, want taller than wide", r)
	}
	// The bottom of the padded content box is lo, the top is hi.
	tt.Press(r.X+r.W/2, r.Y+r.H-14)
	tt.Release(r.X+r.W/2, r.Y+r.H-14)
	if value != 0 {
		t.Fatalf("the bottom set %v", value)
	}
	tt.Key(0, ui.KeyUp)
	if value != 1 {
		t.Fatalf("Up set %v, want 1", value)
	}
	tt.Key(0, ui.KeyEnd)
	if value != 100 {
		t.Fatalf("End set %v", value)
	}
}

func TestSliderFormatAndDescription(t *testing.T) {
	value := 0.0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).FillWidth().Children(func() {
			Slider(c, SliderOptions{
				State:       SliderState{Value: &value},
				Label:       "Opacity",
				Description: "How solid it looks",
				Min:         0,
				Max:         1,
				Step:        0.1,
				Format:      func(v float64) string { return sliderNumberFormat(0.1)(v) + "%" },
			})
		})
	}, 320, 160)
	if !tt.HasText("How solid it looks") {
		t.Fatalf("description missing: %q", tt.Texts())
	}
	if !tt.HasText("0.0%") {
		t.Fatalf("formatted output missing: %q", tt.Texts())
	}
}
