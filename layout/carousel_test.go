package layout_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCarouselMovesRealPagesWithKeyboardBoundsAndRTL(t *testing.T) {
	state := layout.CarouselState{}
	disabled, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		layout.Carousel(c, &state, layout.CarouselOptions{Label: "Examples", Count: 3, Width: 240, Height: 120, Disabled: disabled, OnChange: func(int) { changes++ }}, func(i int) { ui.Text(c, fmt.Sprintf("Page %d", i)) })
	}, 300, 220)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	first, _ := tt.Find("Page 0")
	if err := tt.Click("Previous slide"); err != nil {
		t.Fatal(err)
	}
	if state.Index != 0 || changes != 0 {
		t.Fatal("previous crossed the start")
	}
	if err := tt.Click("Next slide"); err != nil {
		t.Fatal(err)
	}
	if state.Index != 1 || changes != 1 {
		t.Fatal("next failed to select actual page")
	}
	second, _ := tt.Find("Page 1")
	if second.X != first.X {
		t.Fatalf("page strip did not move: first=%v second=%v", first, second)
	}
	if err := tt.Click("Examples"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	if state.Index != 2 {
		t.Fatal("keyboard crossed final boundary")
	}
	tt.SetDirection(ui.RightToLeft)
	tt.Key(0, ui.KeyRight)
	if state.Index != 1 {
		t.Fatal("RTL physical right did not go to previous page")
	}
	disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyHome)
	if err := tt.Click("Next slide"); err != nil {
		t.Fatal(err)
	}
	if state.Index != 1 {
		t.Fatal("disabled carousel changed page")
	}
}

func TestCarouselNativeDragSnapsAndCanBeCaughtMidMotion(t *testing.T) {
	state := layout.CarouselState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		layout.Carousel(c, &state, layout.CarouselOptions{Label: "Pages", Count: 3, Width: 240, Height: 120}, func(i int) { ui.Text(c, fmt.Sprintf("Slide body %d", i)) })
	}, 300, 220)
	tt.Press(200, 60)
	tt.Move(30, 60)
	tt.Release(30, 60)
	if state.Index == 0 {
		t.Fatal("native drag release did not select projected page")
	}
	// Catch its current motion and release at the same point. Native drag
	// must settle at a valid snap rather than keeping an unbounded offset.
	tt.Press(120, 60)
	tt.Release(120, 60)
	tt.Advance(2 * time.Second)
	if state.Index < 0 || state.Index > 2 || tt.FramePending() {
		t.Fatal("caught drag failed to settle within page boundaries")
	}
	state.Index = 0
	tt.Frame()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Press(20, 60)
	tt.Move(200, 60)
	tt.Release(200, 60)
	if state.Index != 0 {
		t.Fatal("drag beyond start escaped first page")
	}
}
