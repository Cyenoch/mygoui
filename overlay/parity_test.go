package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestPopoverEntranceZoomsToFullSize pins HeroUI's zoom-in-90: a panel arrives
// at 90% of its size and reaches its layout size over the popup duration,
// without leaving a frame pending.
func TestPopoverEntranceZoomsToFullSize(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			anchor := ui.Button(c, "Options").OnClick(func() { open = !open })
			Popover(c, anchor, &open, PopoverOptions{
				Title: "Details",
				Body:  func() { ui.Button(c, "Inside") },
			})
		})
	}, 640, 480)

	if err := tt.Click("Options"); err != nil {
		t.Fatal(err)
	}
	entering := motionLabel(t, tt, "Details")
	if !tt.FramePending() {
		t.Error("the popover arrival asks for no frames")
	}
	tt.Advance(popupMotionDuration)
	tt.Frame()
	settled := motionLabel(t, tt, "Details")
	if settled.Dx() <= entering.Dx() {
		t.Fatalf("the panel did not grow to its layout size: entering=%v settled=%v", entering, settled)
	}
	ratio := float64(entering.Dx()) / float64(settled.Dx())
	if ratio < 0.87 || ratio > 0.93 {
		t.Fatalf("entrance scale=%v, want HeroUI's zoom-in-90", ratio)
	}
	if tt.FramePending() {
		t.Error("the popover keeps asking for frames after it settled")
	}
}

// TestPopoverArrowFollowsFlippedPanel puts the anchor too low for the panel
// below it, so the attach engine flips the panel above the anchor. The arrow
// belongs on the edge facing the anchor, not on the side the caller asked
// for.
func TestPopoverArrowFollowsFlippedPanel(t *testing.T) {
	var open bool
	var box ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		pal := theme.Colors(c)
		ui.Column(c).Fill().Background(pal.Background).Children(func() {
			ui.Box(c).Height(360).Shrink(0)
			anchor := ui.Button(c, "Details").OnClick(func() { open = !open })
			if p := Popover(c, anchor, &open, PopoverOptions{
				Title:     "Information",
				Arrow:     true,
				Placement: mygoui.Placement{Side: mygoui.SideBottom},
				Body:      func() { ui.Box(c).Height(120) },
			}); p.Valid() {
				box = p.Bounds()
			}
		})
	}, 640, 480)

	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(2 * popupMotionDuration)
	tt.Frame()
	anchorRect, ok := tt.Find("Details")
	if !ok {
		t.Fatal("the trigger vanished")
	}
	if box.W <= 0 {
		t.Fatal("the popover has no committed box")
	}
	if box.Y+box.H > anchorRect.Y {
		t.Fatalf("the panel did not flip above its anchor: panel=%+v anchor=%+v", box, anchorRect)
	}
	pal := theme.ColorsForMode(false)
	img := tt.Image()
	cx := int(box.X + box.W/2)
	near := img.RGBAAt(cx, int(box.Y+box.H)+3)
	far := img.RGBAAt(cx, int(box.Y)-3)
	if !sameColor(near, pal.Overlay, 3) {
		t.Fatalf("no arrow on the anchor-facing edge: %v, want the overlay color", near)
	}
	if sameColor(far, pal.Overlay, 3) {
		t.Fatalf("the arrow stayed on the requested edge: %v", far)
	}
}

// TestPopoverArrowPointsAtItsAnchor aligns a wide panel's start with a narrow
// anchor: HeroUI's OverlayArrow sits at the trigger's center, not at the
// middle of the panel.
func TestPopoverArrowPointsAtItsAnchor(t *testing.T) {
	var open bool
	var box ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		pal := theme.Colors(c)
		ui.Column(c).Fill().AlignItems(ui.Start).Padding(40).Background(pal.Background).Children(func() {
			anchor := ui.Button(c, "Few words").OnClick(func() { open = !open })
			if p := Popover(c, anchor, &open, PopoverOptions{
				Title:     "Wide panel",
				Arrow:     true,
				Width:     300,
				Placement: mygoui.Placement{Side: mygoui.SideBottom, Align: mygoui.AlignStart},
			}); p.Valid() {
				box = p.Bounds()
			}
		})
	}, 640, 480)
	if err := tt.Click("Few words"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(2 * popupMotionDuration)
	tt.Frame()
	anchorRect, ok := tt.Find("Few words")
	if !ok || box.W <= 0 {
		t.Fatalf("the popover did not open: anchor=%v box=%v", anchorRect, box)
	}
	pal := theme.ColorsForMode(false)
	img := tt.Image()
	first, last := -1, -1
	for x := int(box.X) - 8; x < int(box.X+box.W)+8; x++ {
		if sameColor(img.RGBAAt(x, int(box.Y)-3), pal.Overlay, 3) {
			if first < 0 {
				first = x
			}
			last = x
		}
	}
	if first < 0 {
		t.Fatal("no arrow above the panel")
	}
	arrow := float32(first+last) / 2
	want := anchorRect.X + anchorRect.W/2
	if diff := arrow - want; diff < -2 || diff > 2 {
		t.Fatalf("arrow at %v, want the anchor's center %v", arrow, want)
	}
}

// TestModalBackdropKeepsItsHitAreaWhileEntering presses the backdrop in the
// frame the modal is still fading in: the scrim is never scaled, so its whole
// window stays pressable.
func TestModalBackdropKeepsItsHitAreaWhileEntering(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Modal(c, &open, ModalOptions{Title: "Settings"})
	}, 640, 480)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.FramePending() {
		t.Fatalf("the modal is not entering: open=%v", open)
	}
	tt.ClickAt(4, 4)
	if open || tt.HasText("Settings") {
		t.Fatal("the entering backdrop did not take the press")
	}
	if !tt.Focused("Open") {
		t.Fatal("the backdrop dismissal did not restore the focus")
	}
}

// TestPopoverReopenDuringExitIsWhole reopens a popover while its exit copy
// still paints: the panel returns at its layout size, takes input, and does
// not replay the entrance.
func TestPopoverReopenDuringExitIsWhole(t *testing.T) {
	var open bool
	inside := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			anchor := ui.Button(c, "Options").OnClick(func() { open = !open })
			Popover(c, anchor, &open, PopoverOptions{
				Title: "Details",
				Body:  func() { ui.Button(c, "Inside").OnClick(func() { inside++ }) },
			})
		})
	}, 640, 480)

	if err := tt.Click("Options"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(popupMotionDuration)
	tt.Frame()
	settled := motionLabel(t, tt, "Details")
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not close the popover")
	}
	tt.Advance(popupMotionDuration / 3)
	if err := tt.Click("Options"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.HasText("Details") {
		t.Fatal("the popover did not reopen during its exit")
	}
	// The reopened panel comes back from the exit copy's 95% instead of
	// replaying the 90% entrance.
	if again := motionLabel(t, tt, "Details"); again.Dx() < settled.Dx()*95/100 {
		t.Fatalf("the reopened panel restarted its entrance: settled=%v again=%v", settled, again)
	}
	tt.Advance(popupMotionDuration)
	tt.Frame()
	if got := motionLabel(t, tt, "Details"); got != settled {
		t.Fatalf("the reopened panel did not settle at its layout size: got=%v want=%v", got, settled)
	}
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	if inside != 1 {
		t.Fatalf("the reopened panel did not take the press: %d", inside)
	}
}

// TestModalCoverAndExtraSmallGeometry pins the two sizes added from HeroUI's
// scale: max-w-xs and the cover that fills the container's padding box.
func TestModalCoverAndExtraSmallGeometry(t *testing.T) {
	const winW, winH, pad = 800, 600, 40
	for _, tc := range []struct {
		name string
		size ModalSize
		want ui.Rect
	}{
		{"extra small", ModalExtraSmall, ui.Rect{X: (winW - 320) / 2, Y: 0, W: 320}},
		{"cover", ModalCover, ui.Rect{X: pad, Y: pad, W: winW - 2*pad, H: winH - 2*pad}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				if p := Modal(c, &open, ModalOptions{Size: tc.size, Title: "Panel"}); p.Valid() {
					panel = p.Bounds()
				}
			}, winW, winH)
			tt.Frame()
			tt.Frame()
			if tc.size == ModalExtraSmall {
				if panel.W != tc.want.W || panel.X != tc.want.X {
					t.Fatalf("extra small panel = %+v, want width %v centered", panel, tc.want.W)
				}
				return
			}
			if panel != tc.want {
				t.Fatalf("cover panel = %+v, want %+v", panel, tc.want)
			}
		})
	}
}

// TestDrawerWidthBreakpoints checks HeroUI's w-80 / w-96 pair within the
// 16 DIP window margin and its 85% cap for the floating edge panels.
func TestDrawerWidthBreakpoints(t *testing.T) {
	for _, tc := range []struct {
		winW, want float32
	}{
		{560, 320},
		{640, 384},
		{1000, 384},
		{360, 306},
	} {
		open := true
		var panel ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			if p := Drawer(c, &open, DrawerOptions{Side: DrawerRight, Title: "Filters"}); p.Valid() {
				panel = p.Bounds()
			}
		}, int(tc.winW), 480)
		tt.Frame()
		tt.Frame()
		if panel.W != tc.want {
			t.Fatalf("drawer width at %v = %v, want %v", tc.winW, panel.W, tc.want)
		}
		if panel.X+panel.W != tc.winW-theme.Space16 {
			t.Fatalf("drawer at %v is not inset from the right edge: %+v", tc.winW, panel)
		}
	}
}

// TestDrawerHandleGutter pins HeroUI's handle spacing: the 4 DIP bar follows
// the panel's 24 DIP padding at the top of a bottom drawer and leaves 8 DIPs
// before the content, while a top drawer keeps its bar 8 DIPs off the bottom
// edge.
func TestDrawerHandleGutter(t *testing.T) {
	const winH = 480
	for _, side := range []DrawerSide{DrawerTop, DrawerBottom} {
		open := true
		var panel ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			if p := Drawer(c, &open, DrawerOptions{Side: side, Title: "Filters", Handle: true}); p.Valid() {
				panel = p.Bounds()
			}
		}, 640, winH)
		tt.Frame()
		tt.Frame()
		pal := theme.ColorsForMode(false)
		img := tt.Image()
		top, bottom := -1, -1
		for y := int(panel.Y); y < int(panel.Y+panel.H); y++ {
			if !sameColor(img.RGBAAt(int(panel.X+panel.W/2), y), pal.Separator, 4) {
				continue
			}
			if top < 0 {
				top = y
			}
			bottom = y
		}
		if top < 0 {
			t.Fatalf("%v: the handle bar is missing", side)
		}
		if thickness := bottom - top + 1; thickness < 3 || thickness > 5 {
			t.Fatalf("%v: handle bar is %d DIPs thick, want 4", side, thickness)
		}
		want := panel.Y + 24
		if side == DrawerTop {
			want = panel.Y + panel.H - 8 - 4
		}
		if diff := float32(top) - want; diff < -1.5 || diff > 1.5 {
			t.Fatalf("%v: handle bar starts at %v, want %v", side, float32(top), want)
		}
	}
}
