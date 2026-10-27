package overlay

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDrawerInsetAndCustomSizes(t *testing.T) {
	for _, tc := range []struct {
		name        string
		width       int
		side        DrawerSide
		inset       bool
		panelWidth  float32
		panelHeight float32
		want        ui.Rect
	}{
		{"left inset", 640, DrawerLeft, true, 0, 0, ui.Rect{X: 16, Y: 16, W: 448, H: 448}},
		{"right inset custom", 640, DrawerRight, true, 240, 0, ui.Rect{X: 384, Y: 16, W: 240, H: 448}},
		{"top inset custom", 640, DrawerTop, true, 0, 160, ui.Rect{X: 16, Y: 16, W: 608, H: 160}},
		{"bottom inset custom", 640, DrawerBottom, true, 0, 160, ui.Rect{X: 16, Y: 304, W: 608, H: 160}},
		{"narrow inset", 360, DrawerRight, true, 0, 0, ui.Rect{X: 48, Y: 0, W: 312, H: 480}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				if p := Drawer(c, &open, DrawerOptions{
					Side: tc.side, Inset: tc.inset, Width: tc.panelWidth, Height: tc.panelHeight,
				}); p.Valid() {
					panel = p.Bounds()
				}
			}, tc.width, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			if panel != tc.want {
				t.Fatalf("drawer geometry=%+v, want %+v", panel, tc.want)
			}
		})
	}
}

func TestDrawerRoundsCornersFacingWindow(t *testing.T) {
	for _, tc := range []struct {
		name    string
		side    DrawerSide
		rounded [4]bool
	}{
		{"left", DrawerLeft, [4]bool{false, true, true, false}},
		{"right", DrawerRight, [4]bool{true, false, false, true}},
		{"top", DrawerTop, [4]bool{false, false, true, true}},
		{"bottom", DrawerBottom, [4]bool{true, true, false, false}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Box(c).Fill().Background(ui.RGB(40, 70, 110))
				if p := Drawer(c, &open, DrawerOptions{
					Side: tc.side, Height: 180, Backdrop: ModalBackdropTransparent,
				}); p.Valid() {
					panel = p.Bounds()
				}
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			img := tt.Image()
			inside := img.RGBAAt(int(panel.X+panel.W/2), int(panel.Y+panel.H/2))
			points := [4][2]int{
				{int(panel.X) + 2, int(panel.Y) + 2},
				{int(panel.X+panel.W) - 3, int(panel.Y) + 2},
				{int(panel.X+panel.W) - 3, int(panel.Y+panel.H) - 3},
				{int(panel.X) + 2, int(panel.Y+panel.H) - 3},
			}
			for i, point := range points {
				rounded := sheetColorDistance(img.RGBAAt(point[0], point[1]), inside) > 40
				if rounded != tc.rounded[i] {
					t.Fatalf("corner %d rounded=%v, want %v (panel=%+v)", i, rounded, tc.rounded[i], panel)
				}
			}
		})
	}
}

func TestDrawerCloseButtonOptions(t *testing.T) {
	for _, tc := range []struct {
		name     string
		show     bool
		hide     bool
		hasClose bool
	}{
		{"default", false, false, false},
		{"shown", true, false, true},
		{"hidden override", true, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			dismissals := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				Drawer(c, &open, DrawerOptions{
					Title: "Panel", ShowClose: tc.show, HideClose: tc.hide,
					OnDismiss: func() { dismissals++ },
				})
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if _, found := tt.Find("Close"); found != tc.hasClose {
				t.Fatalf("close button present=%v, want %v", found, tc.hasClose)
			}
			if tc.hasClose {
				if err := tt.Click("Close"); err != nil {
					t.Fatal(err)
				}
			} else {
				tt.Key(0, ui.KeyEscape)
			}
			if open || dismissals != 1 || !tt.Focused("Open") {
				t.Fatalf("dismissal: open=%v dismissals=%d focused=%v", open, dismissals, tt.Focused("Open"))
			}
		})
	}
}

type drawerSwipeFixture struct {
	open       bool
	dismissals int
	panel      ui.Rect
	text       string
	actions    int
}

func newDrawerSwipeFixture(t *testing.T, side DrawerSide, prevent, reduce bool) (*ui.Tester, *drawerSwipeFixture) {
	t.Helper()
	f := &drawerSwipeFixture{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Box(c).Fill().Background(ui.RGB(240, 240, 240)).Children(func() {
			ui.Box(c).Absolute().Left(0).Bottom(0).Width(48).Height(100).Background(ui.RGB(0, 0, 0))
			ui.Button(c, "Open").OnClick(func() { f.open = true })
			if p := Drawer(c, &f.open, DrawerOptions{
				Side: side, Title: "Panel", Width: 260, Height: 240, PreventDismiss: prevent,
				OnDismiss: func() { f.dismissals++ },
				Body: func() {
					ui.TextInput(c, &f.text).Label("Editor")
					ui.Box(c).Size(24, 24).Background(ui.RGB(255, 0, 0))
					ui.Button(c, "Inside").OnClick(func() { f.actions++ })
				},
			}); p.Valid() {
				f.panel = p.Bounds()
			}
		})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	motionSettle(tt, 2*drawerMotionDuration)
	if err := tt.Click("Editor"); err != nil {
		t.Fatal(err)
	}
	return tt, f
}

func drawerSwipeStart(side DrawerSide, r ui.Rect) (x, y, dx, dy float32) {
	switch side {
	case DrawerLeft:
		return r.X + r.W - 12, r.Y + r.H/2, -1, 0
	case DrawerRight:
		return r.X + 12, r.Y + r.H/2, 1, 0
	case DrawerTop:
		return r.X + r.W/2, r.Y + r.H - 12, 0, -1
	default:
		return r.X + r.W/2, r.Y + 12, 0, 1
	}
}

func drawerSwipe(tt *ui.Tester, x, y, dx, dy float32) {
	tt.Move(x, y)
	tt.Press(x, y)
	tt.Advance(16 * time.Millisecond)
	tt.Move(x+dx, y+dy)
	tt.Frame()
}

func TestDrawerSwipeDirectionThresholdAndFocus(t *testing.T) {
	for _, tc := range []struct {
		name string
		side DrawerSide
	}{
		{"left", DrawerLeft}, {"right", DrawerRight}, {"top", DrawerTop}, {"bottom", DrawerBottom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt, f := newDrawerSwipeFixture(t, tc.side, false, false)
			rest := f.panel
			x, y, dx, dy := drawerSwipeStart(tc.side, rest)
			distance := func(r ui.Rect) float32 { return (r.X-rest.X)*dx + (r.Y-rest.Y)*dy }

			drawerSwipe(tt, x, y, -40*dx, -40*dy)
			if d := distance(f.panel); d != 0 {
				t.Fatalf("inward swipe moved the drawer by %v", d)
			}
			tt.Advance(60 * time.Millisecond)
			tt.Release(x-40*dx, y-40*dy)
			motionSettle(tt, 2*drawerMotionDuration)

			drawerSwipe(tt, x, y, 40*dx, 40*dy)
			if d := distance(f.panel); d < 38 || d > 42 {
				t.Fatalf("short swipe moved the drawer by %v, want 40", d)
			}
			if !tt.Focused("Editor") {
				t.Fatal("the grab moved focus away from the editor")
			}
			tt.Advance(60 * time.Millisecond)
			tt.Release(x+40*dx, y+40*dy)
			motionSettle(tt, 2*drawerMotionDuration)
			if !f.open || f.dismissals != 0 || distance(f.panel) != 0 || !tt.Focused("Editor") {
				t.Fatalf("short swipe: open=%v dismissals=%d panel=%+v editor focused=%v",
					f.open, f.dismissals, f.panel, tt.Focused("Editor"))
			}

			drawerSwipe(tt, x, y, 120*dx, 120*dy)
			if d := distance(f.panel); d < 118 || d > 122 {
				t.Fatalf("long swipe moved the drawer by %v, want 120", d)
			}
			tt.Advance(60 * time.Millisecond)
			tt.Release(x+120*dx, y+120*dy)
			if f.open || f.dismissals != 1 || !tt.Focused("Open") {
				t.Fatalf("long swipe: open=%v dismissals=%d opener focused=%v",
					f.open, f.dismissals, tt.Focused("Open"))
			}
			motionSettle(tt, 2*drawerMotionDuration)
			if f.dismissals != 1 || tt.FramePending() {
				t.Fatalf("after exit: dismissals=%d pending=%v", f.dismissals, tt.FramePending())
			}
		})
	}
}

func TestDrawerSwipePreventionAndReducedMotion(t *testing.T) {
	for _, tc := range []struct {
		name string
		side DrawerSide
	}{
		{"left", DrawerLeft}, {"right", DrawerRight}, {"top", DrawerTop}, {"bottom", DrawerBottom},
	} {
		for _, prevent := range []bool{true, false} {
			name := tc.name + "/reduced"
			if prevent {
				name = tc.name + "/prevented"
			}
			t.Run(name, func(t *testing.T) {
				tt, f := newDrawerSwipeFixture(t, tc.side, prevent, true)
				rest := f.panel
				x, y, dx, dy := drawerSwipeStart(tc.side, rest)
				drawerSwipe(tt, x, y, 120*dx, 120*dy)
				tt.Advance(60 * time.Millisecond)
				tt.Release(x+120*dx, y+120*dy)
				tt.Frame()
				if tt.FramePending() {
					t.Fatal("reduced-motion swipe left an animation pending")
				}
				if prevent {
					if !f.open || f.dismissals != 0 || f.panel != rest || !tt.Focused("Editor") {
						t.Fatalf("prevented swipe: open=%v dismissals=%d panel=%+v focused=%v",
							f.open, f.dismissals, f.panel, tt.Focused("Editor"))
					}
				} else if f.open || f.dismissals != 1 || !tt.Focused("Open") {
					t.Fatalf("reduced-motion dismissal: open=%v dismissals=%d focused=%v",
						f.open, f.dismissals, tt.Focused("Open"))
				}
			})
		}
	}
}

func TestDrawerHandleShowsAndSwipes(t *testing.T) {
	for _, tc := range []struct {
		name string
		side DrawerSide
	}{
		{"left", DrawerLeft}, {"right", DrawerRight}, {"top", DrawerTop}, {"bottom", DrawerBottom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				if p := Drawer(c, &open, DrawerOptions{
					Side: tc.side, Width: 260, Height: 240, Title: "Panel", Handle: true,
				}); p.Valid() {
					panel = p.Bounds()
				}
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			x, y, dx, dy := drawerSwipeStart(tc.side, panel)
			switch tc.side {
			case DrawerLeft:
				x = panel.X + panel.W - 16
			case DrawerRight:
				x = panel.X + 16
			case DrawerTop:
				y = panel.Y + panel.H - 10
			case DrawerBottom:
				y = panel.Y + 26
			}
			pal := theme.ColorsForMode(false)
			if got := tt.Image().RGBAAt(int(x), int(y)); !sameColor(got, pal.Input.Over(pal.Popover), 4) {
				t.Fatalf("handle pixel=%v at (%v,%v), panel=%+v", got, x, y, panel)
			}
			drawerSwipe(tt, x, y, 120*dx, 120*dy)
			tt.Advance(60 * time.Millisecond)
			tt.Release(x+120*dx, y+120*dy)
			if open || !tt.Focused("Open") {
				t.Fatalf("handle swipe: open=%v opener focused=%v", open, tt.Focused("Open"))
			}
		})
	}
}

func TestDrawerSwipeFromInset(t *testing.T) {
	for _, tc := range []struct {
		name string
		side DrawerSide
	}{
		{"left", DrawerLeft}, {"right", DrawerRight}, {"top", DrawerTop}, {"bottom", DrawerBottom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				if p := Drawer(c, &open, DrawerOptions{
					Side: tc.side, Width: 260, Height: 240, Inset: true, Title: "Panel",
				}); p.Valid() {
					panel = p.Bounds()
				}
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			x, y, dx, dy := drawerSwipeStart(tc.side, panel)
			drawerSwipe(tt, x, y, 120*dx, 120*dy)
			tt.Advance(60 * time.Millisecond)
			tt.Release(x+120*dx, y+120*dy)
			if open || !tt.Focused("Open") {
				t.Fatalf("inset swipe: open=%v opener focused=%v", open, tt.Focused("Open"))
			}
		})
	}
}

func TestDrawerSwipeFadesDefaultBackdrop(t *testing.T) {
	tt, f := newDrawerSwipeFixture(t, DrawerRight, false, false)
	rest := f.panel
	beforeImage := tt.Image()
	before := beforeImage.RGBAAt(64, 440)
	blurBefore := beforeImage.RGBAAt(46, 440).R
	x, y, dx, dy := drawerSwipeStart(DrawerRight, rest)
	drawerSwipe(tt, x, y, 120*dx, 120*dy)
	duringImage := tt.Image()
	during := duringImage.RGBAAt(64, 440)
	blurDuring := duringImage.RGBAAt(46, 440).R
	if sumRGB(during) <= sumRGB(before) {
		t.Fatalf("default backdrop did not fade while swiping: before=%v during=%v", before, during)
	}
	if blurBefore == 0 || blurDuring >= blurBefore {
		t.Fatalf("default blur did not fade while swiping: before=%d during=%d", blurBefore, blurDuring)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Release(x+120*dx, y+120*dy)
	motionSettle(tt, 2*drawerMotionDuration)
	afterImage := tt.Image()
	after := afterImage.RGBAAt(64, 440)
	if sumRGB(after) <= sumRGB(during) {
		t.Fatalf("default backdrop did not clear on exit: during=%v after=%v", during, after)
	}
	if blurAfter := afterImage.RGBAAt(46, 440).R; blurAfter != 0 {
		t.Fatalf("default blur remained after exit: pixel=%d", blurAfter)
	}
}

func TestDrawerSwipePreventionAnimatesBackAndCanBeCaught(t *testing.T) {
	tt, f := newDrawerSwipeFixture(t, DrawerRight, true, false)
	rest := f.panel
	x, y, dx, dy := drawerSwipeStart(DrawerRight, rest)
	drawerSwipe(tt, x, y, 120*dx, 120*dy)
	if moved := f.panel.X - rest.X; moved < 118 {
		t.Fatalf("prevented drag did not follow the pointer: %v", moved)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Release(x+120*dx, y+120*dy)
	tt.Advance(30 * time.Millisecond)
	tt.Frame()
	moving := f.panel.X - rest.X
	if moving <= 0 || moving >= 120 {
		t.Fatalf("prevented drag did not animate back: %v", moving)
	}
	x, y, _, _ = drawerSwipeStart(DrawerRight, f.panel)
	drawerSwipe(tt, x, y, -20, 0)
	if got := f.panel.X - rest.X; got < moving-22 || got > moving-18 {
		t.Fatalf("second grab did not catch the moving drawer: before=%v after=%v", moving, got)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Release(x-20, y)
	motionSettle(tt, 2*drawerMotionDuration)
	if !f.open || f.dismissals != 0 || f.panel != rest || !tt.Focused("Editor") {
		t.Fatalf("prevented swipe: open=%v dismissals=%d panel=%+v editor focused=%v",
			f.open, f.dismissals, f.panel, tt.Focused("Editor"))
	}
}

func TestDrawerSwipeExitSnapshotIsInert(t *testing.T) {
	tt, f := newDrawerSwipeFixture(t, DrawerRight, false, false)
	rest := f.panel
	x, y, dx, dy := drawerSwipeStart(DrawerRight, rest)
	drawerSwipe(tt, x, y, 120*dx, 120*dy)
	inside, ok := tt.Find("Inside")
	if !ok {
		t.Fatal("the drawer action is missing")
	}
	redX := func() int {
		img := tt.Image()
		for px := 0; px < img.Bounds().Max.X; px++ {
			for py := 0; py < img.Bounds().Max.Y; py++ {
				pixel := img.RGBAAt(px, py)
				if int(pixel.R)-int(pixel.G) > 20 && int(pixel.R)-int(pixel.B) > 20 {
					return px
				}
			}
		}
		return -1
	}
	dragged := redX()
	if dragged < 0 {
		t.Fatal("the dragged panel marker was not painted")
	}
	tt.Advance(60 * time.Millisecond)
	tt.Release(x+120*dx, y+120*dy)
	if f.open || !tt.Focused("Open") {
		t.Fatal("swipe did not dismiss and restore focus")
	}
	tt.ClickAt(inside.X+inside.W/4, inside.Y+inside.H/2)
	if f.actions != 0 {
		t.Fatal("the exiting snapshot accepted a click")
	}
	tt.Advance(20 * time.Millisecond)
	tt.Frame()
	if departing := redX(); departing <= dragged {
		t.Fatalf("exit did not continue toward the edge: dragged x=%d departing x=%d", dragged, departing)
	}
	motionSettle(tt, 2*drawerMotionDuration)
	if remaining := redX(); remaining >= 0 {
		t.Fatalf("the exiting marker remained at x=%d", remaining)
	}
}

func TestDrawerReopensDuringSwipeExit(t *testing.T) {
	tt, f := newDrawerSwipeFixture(t, DrawerRight, false, false)
	rest := f.panel
	x, y, dx, dy := drawerSwipeStart(DrawerRight, rest)
	drawerSwipe(tt, x, y, 120*dx, 120*dy)
	tt.Advance(60 * time.Millisecond)
	tt.Release(x+120*dx, y+120*dy)
	if f.open || f.dismissals != 1 {
		t.Fatalf("swipe: open=%v dismissals=%d", f.open, f.dismissals)
	}
	tt.Advance(50 * time.Millisecond)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	motionSettle(tt, 2*drawerMotionDuration)
	if !f.open || f.panel != rest {
		t.Fatalf("reopened drawer did not settle: open=%v panel=%+v rest=%+v", f.open, f.panel, rest)
	}
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	if f.actions != 1 {
		t.Fatalf("reopened drawer action count=%d", f.actions)
	}
	tt.Key(0, ui.KeyEscape)
	if f.open || f.dismissals != 2 || !tt.Focused("Open") {
		t.Fatalf("reopened drawer dismissal: open=%v dismissals=%d focused=%v",
			f.open, f.dismissals, tt.Focused("Open"))
	}
}

func TestDrawerBodyScrollAndControlsDoNotSwipe(t *testing.T) {
	var open bool
	var panel ui.Rect
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		if p := Drawer(c, &open, DrawerOptions{
			Side: DrawerRight, Width: 300, Title: "Panel",
			Body: func() {
				for range 25 {
					ui.Text(c, "Scrollable row")
				}
				ui.Button(c, "Last action").OnClick(func() { actions++ })
			},
		}); p.Valid() {
			panel = p.Bounds()
		}
	}, 640, 480)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	motionSettle(tt, 2*drawerMotionDuration)
	rest := panel
	x, y := rest.X+100, rest.Y+200
	drawerSwipe(tt, x, y, 120, 0)
	tt.Advance(60 * time.Millisecond)
	tt.Release(x+120, y)
	motionSettle(tt, 2*drawerMotionDuration)
	if !open || panel != rest {
		t.Fatalf("dragging the body moved or dismissed it: open=%v panel=%+v rest=%+v", open, panel, rest)
	}
	tt.Scroll(x, y, 0, 1000)
	tt.Frame()
	if err := tt.Click("Last action"); err != nil {
		t.Fatal(err)
	}
	if !open || actions != 1 {
		t.Fatalf("scroll/action: open=%v actions=%d", open, actions)
	}
}

func TestDrawerGrabDoesNotCoverCloseButtonInNarrowWindow(t *testing.T) {
	var open bool
	dismissals := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Drawer(c, &open, DrawerOptions{
			Side: DrawerBottom, Height: 160, Title: "Panel", ShowClose: true,
			OnDismiss: func() { dismissals++ },
		})
	}, 80, 300)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if open || dismissals != 1 || !tt.Focused("Open") {
		t.Fatalf("close in narrow window: open=%v dismissals=%d focused=%v",
			open, dismissals, tt.Focused("Open"))
	}
}
