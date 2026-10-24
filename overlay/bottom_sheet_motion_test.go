package overlay

import (
	"image/color"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

const (
	sheetWinW = 640
	sheetWinH = 480
)

// sheetTest drives one BottomSheet: the view records the panel's
// transformed box in every frame it shows, and the dismissal the sheet
// asked for.
type sheetTest struct {
	open      bool
	dismissed int
	shown     bool
	closed    bool
	panel     ui.Rect
	boxes     []ui.Rect
	closedY   float32
}

// newSheetTest builds opts in the view, with "Open" as the opener and the
// sheet's own dismissal counted in the fixture.
func newSheetTest(t *testing.T, build func(c *ui.Context) BottomSheetOptions) (*ui.Tester, *sheetTest) {
	t.Helper()
	f := &sheetTest{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		o := build(c)
		o.OnDismiss = func() { f.dismissed++ }
		ui.Button(c, "Open").OnClick(func() { f.open = true })
		if p := BottomSheet(c, &f.open, o); p.Valid() {
			f.shown = true
			f.panel = p.Bounds()
			f.boxes = append(f.boxes, f.panel)
			return
		}
		if f.shown && !f.closed {
			// The box the sheet last showed before it closed: a swipe must
			// have carried it below the window.
			f.closed = true
			f.closedY = f.panel.Y
		}
	}, sheetWinW, sheetWinH)
	return tt, f
}

// openSheet clicks the opener.
func openSheet(t *testing.T, tt *ui.Tester) {
	t.Helper()
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Actions") {
		t.Fatalf("the sheet did not open; labels=%q", tt.Texts())
	}
}

// settleEntrance lets the arrival finish and returns the resting box.
func (f *sheetTest) settleEntrance(t *testing.T, tt *ui.Tester) ui.Rect {
	t.Helper()
	tt.Advance(drawerMotionDuration * 2)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("the sheet never comes to rest")
	}
	return f.panel
}

// dragGrab presses the sheet's grab strip and moves the pointer by dy. The
// frame after the move is built, so the fixture's box is the moved one.
func dragGrab(tt *ui.Tester, rest ui.Rect, dy float32) (x, y float32) {
	x, y = rest.X+rest.W/2, rest.Y+12
	tt.Move(x, y)
	tt.Press(x, y)
	tt.Advance(16 * time.Millisecond)
	tt.Move(x, y+dy)
	tt.Frame()
	return x, y + dy
}

func sheetColorDistance(a, b color.RGBA) int {
	d := 0
	for _, pair := range [][2]uint8{{a.R, b.R}, {a.G, b.G}, {a.B, b.B}} {
		v := int(pair[0]) - int(pair[1])
		if v < 0 {
			v = -v
		}
		d += v
	}
	return d
}

// TestBottomSheetFloatsInsideTheWindowGutters pins the floating geometry:
// 16 DIPs from every edge the panel touches, the inset width, and all four
// corners rounded, unlike the flush sheet it replaces.
func TestBottomSheetFloatsInsideTheWindowGutters(t *testing.T) {
	tt, f := newSheetTest(t, func(*ui.Context) BottomSheetOptions {
		return BottomSheetOptions{Title: "Actions"}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)
	if rest.X != 16 || rest.W != sheetWinW-32 {
		t.Fatalf("the sheet sits at %+v, want 16 DIPs from either side", rest)
	}
	if rest.Y+rest.H != sheetWinH-16 {
		t.Fatalf("the sheet ends at %v, want 16 DIPs above the window bottom", rest.Y+rest.H)
	}
	if rest.H > (sheetWinH-32)*0.85+1 {
		t.Fatalf("the sheet is %v tall, past 85%% of the window", rest.H)
	}

	img := tt.Image()
	inside := img.RGBAAt(int(rest.X+rest.W/2), int(rest.Y+rest.H/2))
	corners := [][2]int{
		{int(rest.X) + 2, int(rest.Y) + 2},
		{int(rest.X+rest.W) - 3, int(rest.Y) + 2},
		{int(rest.X) + 2, int(rest.Y+rest.H) - 3},
		{int(rest.X+rest.W) - 3, int(rest.Y+rest.H) - 3},
	}
	for _, p := range corners {
		if got := img.RGBAAt(p[0], p[1]); sheetColorDistance(got, inside) < 24 {
			t.Fatalf("corner (%d,%d) = %v matches the panel %v: the corner is not rounded", p[0], p[1], got, inside)
		}
	}
}

// TestBottomSheetSlidesUpFromBelowTheWindow checks the arrival: the first
// frame is below the window, the frames after it move, and the last one
// rests without asking for more.
func TestBottomSheetSlidesUpFromBelowTheWindow(t *testing.T) {
	tt, f := newSheetTest(t, func(c *ui.Context) BottomSheetOptions {
		return BottomSheetOptions{
			Title: "Actions",
			Body: func() {
				ui.Column(c).Gap(theme.Space8).Children(func() {
					for range 3 {
						ui.Text(c, "A row of the sheet")
					}
				})
			},
		}
	})
	openSheet(t, tt)
	if !tt.FramePending() {
		t.Fatal("the arrival asks for no frames")
	}
	first := f.panel
	if first.H == 0 || first.Y < sheetWinH {
		t.Fatalf("the first frame shows the sheet at %+v, want it below the window", first)
	}

	tt.Advance(drawerMotionDuration / 2)
	tt.Frame()
	mid := f.panel
	if mid.Y >= first.Y {
		t.Fatalf("the sheet is at %v halfway through, no higher than its first frame %v", mid, first)
	}

	rest := f.settleEntrance(t, tt)
	if rest.Y >= mid.Y {
		t.Fatalf("the sheet rests at %v, no higher than its halfway frame %v", rest, mid)
	}
	if rest.H != first.H {
		t.Fatalf("the sheet's height changed from %v to %v while it slid in", first.H, rest.H)
	}
	if rest.Y+rest.H != sheetWinH-16 {
		t.Fatalf("the sheet came to rest at %+v, want it above the bottom gutter", rest)
	}
}

// TestBottomSheetReducedMotionStartsAtRest checks that the sheet is in its
// place in its first frame when the desktop asks for less motion.
func TestBottomSheetReducedMotionStartsAtRest(t *testing.T) {
	tt, f := newSheetTest(t, func(*ui.Context) BottomSheetOptions {
		return BottomSheetOptions{Title: "Actions"}
	})
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	openSheet(t, tt)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("reduced motion asks for frames as the sheet opens")
	}
	rest := f.panel
	if rest.Y+rest.H != sheetWinH-16 {
		t.Fatalf("reduced motion rests at %+v, want the sheet in its place at once", rest)
	}
	for _, b := range f.boxes {
		if b.Y >= sheetWinH {
			t.Fatalf("reduced motion showed the sheet below the window at %+v", b)
		}
	}
}

// TestBottomSheetShortDragSpringsBack checks that the grab strip drags the
// sheet with the pointer and that a drag short of the threshold returns it.
func TestBottomSheetShortDragSpringsBack(t *testing.T) {
	tt, f := newSheetTest(t, func(c *ui.Context) BottomSheetOptions {
		return BottomSheetOptions{
			Title:  "Actions",
			Height: 300,
			Body: func() {
				ui.Column(c).Gap(theme.Space8).Children(func() {
					for range 4 {
						ui.Text(c, "A row of the sheet")
					}
				})
			},
		}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)

	x, y := dragGrab(tt, rest, 40)
	if got := f.panel.Y; got < rest.Y+38 || got > rest.Y+42 {
		t.Fatalf("dragged 40 DIPs, the sheet moved to %v from %v", got, rest.Y)
	}
	// Let the pointer rest before letting go: a slow drag, short of the
	// dismissal threshold.
	tt.Advance(60 * time.Millisecond)
	tt.Release(x, y)
	if f.dismissed != 0 || !f.open {
		t.Fatalf("a short drag dismissed the sheet: open=%v dismissed=%d", f.open, f.dismissed)
	}
	tt.Advance(2 * time.Second)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("the spring back never comes to rest")
	}
	if f.panel.Y != rest.Y {
		t.Fatalf("a short drag left the sheet at %v, want its place %v", f.panel.Y, rest.Y)
	}

	// Dragged up, it stays at its place: the sheet only goes down.
	x, y = dragGrab(tt, rest, -40)
	if f.panel.Y != rest.Y {
		t.Fatalf("an upward drag moved the sheet to %v, want %v", f.panel.Y, rest.Y)
	}
	tt.Release(x, y)
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if f.panel.Y != rest.Y || f.dismissed != 0 {
		t.Fatalf("an upward drag left the sheet at %v with %d dismissals", f.panel.Y, f.dismissed)
	}
}

// TestBottomSheetPaintsItsSlideAndDeparture checks the painted frames of
// the arrival and the departure: the sheet climbs into the window rather
// than fading in place, and the backdrop still fades as it goes.
func TestBottomSheetPaintsItsSlideAndDeparture(t *testing.T) {
	tt, f := newSheetTest(t, func(*ui.Context) BottomSheetOptions {
		return BottomSheetOptions{Title: "Actions"}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)
	area := motionArea(rest)

	// Close once to learn the resting area, and open again for its frames.
	tt.Key(0, ui.KeyEscape)
	tt.Advance(2 * drawerMotionDuration)
	tt.Frame()
	if f.open || tt.HasText("Actions") {
		t.Fatal("Escape did not close the sheet")
	}
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Actions") {
		t.Fatal("the sheet did not reopen")
	}
	motionAssertSurface(t, tt, area, drawerMotionDuration, func() {
		tt.Key(0, ui.KeyEscape)
		if f.open || tt.HasText("Actions") {
			t.Fatal("Escape did not close the reopened sheet")
		}
	})
	if !tt.Focused("Open") {
		t.Fatal("the sheet did not give the focus back")
	}
}

// TestBottomSheetGrabCatchesTheSpringingSheet checks that the grab strip
// stays on the handle while the sheet moves: a second grab on the way back
// catches it and follows the pointer.
func TestBottomSheetGrabCatchesTheSpringingSheet(t *testing.T) {
	tt, f := newSheetTest(t, func(*ui.Context) BottomSheetOptions {
		return BottomSheetOptions{Title: "Actions", Height: 300}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)

	// Let go short of the threshold and press again while it springs back.
	x, y := dragGrab(tt, rest, 80)
	tt.Advance(60 * time.Millisecond)
	tt.Release(x, y)
	tt.Advance(30 * time.Millisecond)
	tt.Frame()
	moving := f.panel.Y
	if moving <= rest.Y || moving >= rest.Y+80 {
		t.Fatalf("the sheet is not on its way back: %v from %v", moving, rest.Y)
	}
	hx, hy := rest.X+rest.W/2, moving+12
	tt.Move(hx, hy)
	tt.Press(hx, hy)
	tt.Advance(16 * time.Millisecond)
	tt.Move(hx, hy-20)
	tt.Frame()
	if got := f.panel.Y; got > moving-15 || got < moving-25 {
		t.Fatalf("a grab while springing did not catch the sheet: %v from %v", got, moving)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Release(hx, hy-20)
	tt.Advance(2 * time.Second)
	tt.Frame()
	if !f.open || f.dismissed != 0 {
		t.Fatalf("letting go of the caught sheet: open=%v dismissed=%d", f.open, f.dismissed)
	}
}

// TestBottomSheetCloseButtonSlidesOut checks that the close button runs the
// same exit as the handle swipe: the panel slides below the window and the
// scrim clears with it, instead of vanishing before the mask.
func TestBottomSheetCloseButtonSlidesOut(t *testing.T) {
	tt, f := newSheetTest(t, func(*ui.Context) BottomSheetOptions {
		return BottomSheetOptions{Title: "Actions", Height: 300}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)
	before := tt.Image().RGBAAt(8, 8)
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if !f.open {
		t.Fatal("the close button closed the sheet before it slid out")
	}
	tt.Advance(120 * time.Millisecond)
	tt.Frame()
	if f.panel.Y <= rest.Y+20 {
		t.Fatalf("the sheet did not slide out from the close button: %v from %v", f.panel.Y, rest.Y)
	}
	mid := tt.Image().RGBAAt(8, 8)
	if sumRGB(mid) <= sumRGB(before) {
		t.Fatalf("the scrim did not clear with the sheet: before=%v mid=%v", before, mid)
	}
	tt.Advance(2 * time.Second)
	tt.Frame()
	if f.open || f.dismissed != 1 {
		t.Fatalf("the close button: open=%v dismissed=%d", f.open, f.dismissed)
	}
	if !tt.Focused("Open") {
		t.Fatal("closing the sheet did not restore the focus")
	}
}

func sumRGB(c color.RGBA) int { return int(c.R) + int(c.G) + int(c.B) }

// TestBottomSheetLongDragSlidesOutAndDismisses checks the swipe dismissal:
// the sheet follows the pointer, slides below the window, and only then
// closes and gives the focus back to its opener.
func TestBottomSheetLongDragSlidesOutAndDismisses(t *testing.T) {
	tt, f := newSheetTest(t, func(c *ui.Context) BottomSheetOptions {
		return BottomSheetOptions{
			Title:  "Actions",
			Height: 300,
			Body: func() {
				ui.Column(c).Gap(theme.Space8).Children(func() {
					for range 4 {
						ui.Text(c, "A row of the sheet")
					}
				})
			},
		}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)
	if !tt.Focused("Close") {
		t.Fatalf("the sheet did not take the focus; focused=%q", tt.Texts())
	}

	x, y := dragGrab(tt, rest, 150)
	if f.panel.Y <= rest.Y+100 {
		t.Fatalf("dragged 150 DIPs, the sheet is at %v", f.panel.Y)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Release(x, y)
	if !f.open {
		t.Fatal("letting go closed the sheet before it slid out")
	}
	tt.Advance(120 * time.Millisecond)
	tt.Frame()
	if !f.open || f.panel.Y <= rest.Y+150 {
		t.Fatalf("the sheet did not slide out on its way down: open=%v y=%v", f.open, f.panel.Y)
	}

	tt.Advance(2 * time.Second)
	tt.Frame()
	if f.open || f.dismissed != 1 {
		t.Fatalf("the swipe: open=%v dismissed=%d, want one dismissal", f.open, f.dismissed)
	}
	if !f.closed || f.closedY < sheetWinH-16 {
		t.Fatalf("the sheet closed at y=%v, before it was below the window", f.closedY)
	}
	if !tt.Focused("Open") {
		t.Fatal("the swipe did not give the focus back to the opener")
	}
}

// TestBottomSheetLongDragHonoursPreventDismiss checks that a sheet which
// may not be dismissed comes back instead of sliding out.
func TestBottomSheetLongDragHonoursPreventDismiss(t *testing.T) {
	tt, f := newSheetTest(t, func(c *ui.Context) BottomSheetOptions {
		return BottomSheetOptions{
			Title:          "Actions",
			Height:         300,
			PreventDismiss: true,
		}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)
	x, y := dragGrab(tt, rest, 150)
	tt.Advance(60 * time.Millisecond)
	tt.Release(x, y)
	tt.Advance(2 * time.Second)
	tt.Frame()
	if !f.open || f.dismissed != 0 {
		t.Fatalf("a prevented dismissal closed the sheet: open=%v dismissed=%d", f.open, f.dismissed)
	}
	if f.panel.Y != rest.Y {
		t.Fatalf("the prevented swipe left the sheet at %v, want %v", f.panel.Y, rest.Y)
	}
}

// TestBottomSheetBodyDragDoesNotMoveTheSheet checks that only the grab
// strip owns the swipe: a drag over the body leaves the panel alone.
func TestBottomSheetBodyDragDoesNotMoveTheSheet(t *testing.T) {
	tt, f := newSheetTest(t, func(c *ui.Context) BottomSheetOptions {
		return BottomSheetOptions{
			Title:  "Actions",
			Height: 300,
			Body: func() {
				ui.Column(c).Gap(theme.Space8).Children(func() {
					for range 4 {
						ui.Text(c, "A row of the sheet")
					}
				})
			},
		}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)

	x, y := rest.X+rest.W/2, rest.Y+140
	tt.Move(x, y)
	tt.Press(x, y)
	tt.Advance(16 * time.Millisecond)
	tt.Move(x, y+150)
	tt.Frame()
	tt.Advance(60 * time.Millisecond)
	tt.Release(x, y+150)
	tt.Advance(2 * time.Second)
	tt.Frame()
	if !f.open || f.dismissed != 0 {
		t.Fatalf("a body drag dismissed the sheet: open=%v dismissed=%d", f.open, f.dismissed)
	}
	if f.panel.Y != rest.Y {
		t.Fatalf("a body drag moved the sheet to %v, want %v", f.panel.Y, rest.Y)
	}
}

// TestBottomSheetHideHandleRemovesTheSwipe checks that a sheet without its
// handle has neither the bar nor the gesture: the same swipe over its top
// leaves it open.
func TestBottomSheetHideHandleRemovesTheSwipe(t *testing.T) {
	tt, f := newSheetTest(t, func(c *ui.Context) BottomSheetOptions {
		return BottomSheetOptions{
			Title:      "Actions",
			Height:     300,
			HideHandle: true,
		}
	})
	openSheet(t, tt)
	rest := f.settleEntrance(t, tt)
	x, y := dragGrab(tt, rest, 150)
	tt.Advance(60 * time.Millisecond)
	tt.Release(x, y)
	tt.Advance(2 * time.Second)
	tt.Frame()
	if !f.open || f.dismissed != 0 {
		t.Fatalf("a sheet without a handle dismissed: open=%v dismissed=%d", f.open, f.dismissed)
	}
	if f.panel.Y != rest.Y {
		t.Fatalf("a sheet without a handle moved to %v, want %v", f.panel.Y, rest.Y)
	}
}
