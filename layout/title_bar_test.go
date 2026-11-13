package layout_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/layout"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func newTester(view func(c *ui.Context), width, height int) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		view(c)
	}, width, height)
	tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
	tt.Frame()
	return tt
}

func rect(t *testing.T, tt *ui.Tester, label string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("%q is not shown; texts: %v", label, tt.Texts())
	}
	return r
}

// near allows for float32 layout rounding.
func near(got, want float32) bool { return got > want-0.01 && got < want+0.01 }

func button(c *ui.Context, label string, onClick func()) {
	mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Variant: mygoui.ButtonGhost, Size: theme.ButtonSmall, OnClick: onClick,
	})
}

func barWith(span layout.TitleBarSpan, clicks *int) func(c *ui.Context) {
	return func(c *ui.Context) {
		layout.TitleBar(c, layout.TitleBarOptions{
			Title: "Inbox", Label: "Window bar", Span: span, MinHeight: 52,
			Leading:  func() { button(c, "Back", func() { *clicks++ }) },
			Trailing: func() { button(c, "Share", nil) },
		})
	}
}

func TestTitleBarReservesTheInsetsItsSpanReaches(t *testing.T) {
	const width = 640
	bar := ui.TitleBar{Height: 40, Left: 80, Right: 138}
	for _, tc := range []struct {
		name        string
		span        layout.TitleBarSpan
		rtl         bool
		left, right float32
	}{
		{"window", layout.SpanWindow, false, 80, 138},
		{"start", layout.SpanStart, false, 80, 0},
		{"end", layout.SpanEnd, false, 0, 138},
		{"start, right to left", layout.SpanStart, true, 0, 138},
		{"end, right to left", layout.SpanEnd, true, 80, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var clicks int
			tt := newTester(barWith(tc.span, &clicks), width, 300)
			if tc.rtl {
				tt.SetDirection(ui.RightToLeft)
			}
			tt.SetTitleBar(bar)
			first, last := rect(t, tt, "Back"), rect(t, tt, "Share")
			if tc.rtl {
				first, last = last, first
			}
			pad := float32(theme.Space12)
			if !near(first.X, tc.left+pad) || !near(last.X+last.W, width-tc.right-pad) {
				t.Fatalf("controls span %v..%v, want %v..%v", first.X, last.X+last.W, tc.left+pad, width-tc.right-pad)
			}
			if r := rect(t, tt, "Window bar"); r.H < 52 || r.W != width {
				t.Fatalf("bar = %+v, want the window's width and at least the 52 DIP band", r)
			}
		})
	}
}

func TestTitleBarFollowsTheWindowBandAndFullScreen(t *testing.T) {
	var clicks int
	tt := newTester(barWith(layout.SpanWindow, &clicks), 640, 300)
	tt.SetTitleBar(ui.TitleBar{Height: 70, Left: 80})
	if r := rect(t, tt, "Window bar"); r.H != 70 {
		t.Fatalf("bar height = %v, want the reported 70 DIP band", r.H)
	}
	windowed := rect(t, tt, "Back")
	tt.SetTitleBar(ui.TitleBar{})
	if r := rect(t, tt, "Back"); r.X >= windowed.X {
		t.Fatalf("full screen kept the inset: %+v, was %+v", r, windowed)
	}
}

// MinHeight is the native band in window DIPs: content zoom does not grow it.
func TestTitleBarCoversTheNativeBandUnderZoom(t *testing.T) {
	tt := newTester(func(c *ui.Context) {
		layout.TitleBar(c, layout.TitleBarOptions{Label: "Window bar", MinHeight: 120})
	}, 640, 300)
	tt.SetZoom(2)
	if r := rect(t, tt, "Window bar"); !near(r.H, 120) {
		t.Fatalf("zoomed bar height = %v, want the 120 DIP band", r.H)
	}
}

func TestTitleBarMovesControlsBelowTheBandWhenTheyDoNotFit(t *testing.T) {
	const width, band = 320, 52
	view := func(c *ui.Context) {
		// The bar is the caller's to decorate: its own drawing does not
		// stand in the way of the measurement.
		layout.TitleBar(c, layout.TitleBarOptions{
			Label: "Window bar", MinHeight: band,
			Leading:  func() { button(c, "Navigation", nil) },
			Trailing: func() { button(c, "Appearance", nil) },
		}).DrawOver(func(*ui.Painter, ui.Rect) {})
	}
	tt := newTester(view, width, 300)
	tt.SetTitleBar(ui.TitleBar{Height: band, Left: 80, Right: 138})
	bar := rect(t, tt, "Window bar")
	for _, label := range []string{"Navigation", "Appearance"} {
		r := rect(t, tt, label)
		if r.Y < band || r.Y+r.H > bar.Y+bar.H || r.X < 0 || r.X+r.W > width {
			t.Fatalf("%s = %+v, want it below the %d DIP band inside the bar %+v", label, r, band, bar)
		}
	}
	// Wide enough again, they return to the band between the insets.
	tt.SetSize(640, 300)
	if r := rect(t, tt, "Navigation"); r.Y >= band || r.X < 80 {
		t.Fatalf("Navigation = %+v, want it back in the band after the inset", r)
	}
	// Without window controls, moving below the band would gain no room.
	tt.SetSize(160, 300)
	tt.SetTitleBar(ui.TitleBar{})
	if r := rect(t, tt, "Navigation"); r.Y >= band {
		t.Fatalf("Navigation = %+v, want it in the band of a window without controls", r)
	}
}

func TestTitleBarControlsTakePressesAndTheBarTakesTheRest(t *testing.T) {
	var clicks int
	tt := newTester(barWith(layout.SpanWindow, &clicks), 640, 300)
	tt.SetTitleBar(ui.TitleBar{Height: 52, Left: 80, Right: 138})
	if err := tt.Click("Back"); err != nil || clicks != 1 {
		t.Fatalf("pressing a control in the bar: clicks=%d err=%v", clicks, err)
	}
	// Click succeeds only when the press at the bar's center reaches the
	// bar itself: the surface that moves the window.
	if err := tt.Click("Window bar"); err != nil {
		t.Fatalf("the bar is not the input target of its own surface: %v", err)
	}
	if clicks != 1 {
		t.Fatalf("a press on the bar reached a control: clicks=%d", clicks)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Back") {
		t.Fatal("Tab did not reach the leading control first")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Share") {
		t.Fatal("Tab did not go on to the trailing control")
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if clicks != 2 {
		t.Fatalf("Enter on the focused control: clicks=%d", clicks)
	}
}
