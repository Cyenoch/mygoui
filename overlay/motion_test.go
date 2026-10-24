package overlay

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// motionShot copies the last frame: the tester's pixels are its own.
func motionShot(tt *ui.Tester) *image.RGBA {
	frame := tt.Image()
	shot := image.NewRGBA(frame.Bounds())
	for y := frame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
		start := frame.PixOffset(0, y)
		copy(shot.Pix[shot.PixOffset(0, y):], frame.Pix[start:start+frame.Bounds().Dx()*4])
	}
	return shot
}

func motionArea(r ui.Rect) image.Rectangle {
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
}

func motionDiff(a, b *image.RGBA, r image.Rectangle) int {
	r = r.Intersect(a.Bounds()).Intersect(b.Bounds())
	n := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if a.RGBAAt(x, y) != b.RGBAAt(x, y) {
				n++
			}
		}
	}
	return n
}

func motionAreaSize(r image.Rectangle) int { return max(r.Dx()*r.Dy(), 1) }

// motionLabel returns the painted area of label, which stays where it is
// while a surface fades.
func motionLabel(t *testing.T, tt *ui.Tester, label string) image.Rectangle {
	t.Helper()
	r, ok := tt.Find(label)
	if !ok {
		t.Fatalf("missing %q; labels=%q", label, tt.Texts())
	}
	return motionArea(r)
}

// motionSettle moves past d and builds the frame at rest.
func motionSettle(tt *ui.Tester, d time.Duration) {
	tt.Advance(d)
	tt.Frame()
}

// motionAssertSurface checks that a surface paints on the way in and out of
// the frame it shows: the frames before it settles and after it departs must
// differ inside area, as must the settled frame and the departed one.
func motionAssertSurface(t *testing.T, tt *ui.Tester, area image.Rectangle, duration time.Duration, dismiss func()) {
	t.Helper()
	start := motionShot(tt)
	if !tt.FramePending() {
		t.Error("arrival does not ask for frames")
	}
	tt.Advance(duration / 3)
	mid := motionShot(tt)
	motionSettle(tt, duration)
	settled := motionShot(tt)
	if motionDiff(start, mid, area) == 0 {
		t.Error("arrival paints nothing on the way in")
	}
	if motionDiff(mid, settled, area) == 0 {
		t.Error("arrival never reaches its paint")
	}
	dismiss()
	if !tt.FramePending() {
		t.Error("departure does not ask for frames")
	}
	exitStart := motionShot(tt)
	if motionDiff(settled, exitStart, area) > motionAreaSize(area)/20 {
		t.Error("departure does not paint the frame it last showed")
	}
	tt.Advance(duration / 3)
	exitMid := motionShot(tt)
	if motionDiff(exitStart, exitMid, area) == 0 {
		t.Error("departure paints nothing on the way out")
	}
	motionSettle(tt, duration)
	exitEnd := motionShot(tt)
	if motionDiff(settled, exitEnd, area) == 0 {
		t.Error("the surface stays painted after its departure")
	}
	motionSettle(tt, duration)
	if again := motionShot(tt); motionDiff(exitEnd, again, area) != 0 {
		t.Error("the surface keeps painting after it departed")
	}
}

func TestModalMotionPaintsArrivalAndDeparture(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Modal(c, &open, ModalOptions{
			Title: "Settings",
			Body:  func() { ui.Button(c, "Confirm") },
		})
	}, 640, 480)

	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Settings") {
		t.Fatal("the modal did not open")
	}
	area := motionLabel(t, tt, "Settings")
	motionAssertSurface(t, tt, area, dialogMotionDuration, func() {
		tt.Key(0, ui.KeyEscape)
		if open || tt.HasText("Settings") {
			t.Fatal("Escape did not close the modal")
		}
	})
	if !tt.Focused("Open") {
		t.Fatal("the modal did not give the focus back")
	}
}

// TestEdgeSurfacesPaintArrivalAndDeparture covers every edge a drawer comes
// from and an edge dialog: the native backdrop is absolute in the window,
// and its motion must survive that layout. The bottom sheet slides instead,
// so it has its own arrival and departure test.
func TestEdgeSurfacesPaintArrivalAndDeparture(t *testing.T) {
	for _, tc := range []struct {
		name     string
		duration time.Duration
		build    func(c *ui.Context, open *bool)
	}{
		{"drawer left", drawerMotionDuration, func(c *ui.Context, open *bool) {
			Drawer(c, open, DrawerOptions{Side: DrawerLeft, Title: "Panel"})
		}},
		{"drawer right", drawerMotionDuration, func(c *ui.Context, open *bool) {
			Drawer(c, open, DrawerOptions{Side: DrawerRight, Title: "Panel"})
		}},
		{"drawer top", drawerMotionDuration, func(c *ui.Context, open *bool) {
			Drawer(c, open, DrawerOptions{Side: DrawerTop, Title: "Panel"})
		}},
		{"drawer bottom", drawerMotionDuration, func(c *ui.Context, open *bool) {
			Drawer(c, open, DrawerOptions{Side: DrawerBottom, Title: "Panel"})
		}},
		{"edge dialog", dialogMotionDuration, func(c *ui.Context, open *bool) {
			Dialog(c, open, DialogOptions{
				Placement: mygoui.Placement{Side: mygoui.SideRight},
				Title:     "Panel",
			})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				tc.build(c, &open)
			}, 640, 480)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			if !tt.HasText("Panel") {
				t.Fatal("the surface did not open")
			}
			area := motionLabel(t, tt, "Panel")
			motionAssertSurface(t, tt, area, tc.duration, func() {
				tt.Key(0, ui.KeyEscape)
				if open || tt.HasText("Panel") {
					t.Fatal("Escape did not close the surface")
				}
			})
			if !tt.Focused("Open") {
				t.Fatal("the surface did not give the focus back")
			}
		})
	}
}

func TestPopoverMotionPaintsArrivalAndDeparture(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			anchor := ui.Button(c, "Options").OnClick(func() { open = !open })
			Popover(c, anchor, &open, PopoverOptions{
				Title: "Actions",
				Body:  func() { ui.Button(c, "Rename") },
			})
		})
	}, 640, 480)

	if err := tt.Click("Options"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Actions") {
		t.Fatal("the popover did not open")
	}
	area := motionLabel(t, tt, "Actions")
	motionAssertSurface(t, tt, area, popupMotionDuration, func() {
		tt.Key(0, ui.KeyEscape)
		if open || tt.HasText("Actions") {
			t.Fatal("Escape did not close the popover")
		}
	})
	if !tt.Focused("Options") {
		t.Fatal("the popover did not give the focus back")
	}
}

func TestTooltipMotionPaintsArrivalAndDeparture(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		anchor := ui.Button(c, "Hover me")
		Tooltip(c, anchor, TooltipOptions{Text: "Help"})
	}, 640, 480)

	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if !tt.HasText("Help") || !tt.Focused("Hover me") {
		t.Fatalf("the keyboard focus did not show the tooltip; labels=%q", tt.Texts())
	}
	area := motionLabel(t, tt, "Help")
	motionAssertSurface(t, tt, area, popupMotionDuration, func() {
		tt.Key(0, ui.KeyEscape)
		if tt.HasText("Help") {
			t.Fatal("Escape did not dismiss the tooltip")
		}
	})
	if !tt.Focused("Hover me") {
		t.Fatal("dismissing the tooltip lost the keyboard focus")
	}
}

func TestModalReopenDuringExitSettlesOnce(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Modal(c, &open, ModalOptions{
			Title: "Settings",
			Body:  func() { ui.Button(c, "Confirm") },
		})
	}, 640, 480)

	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	motionSettle(tt, dialogMotionDuration)
	area := motionLabel(t, tt, "Settings")
	settled := motionShot(tt)

	open = false
	tt.Frame()
	if tt.HasText("Settings") {
		t.Fatal("the modal did not close")
	}
	tt.Advance(dialogMotionDuration / 4)
	if motionDiff(motionShot(tt), settled, area) == 0 {
		t.Fatal("the modal departed without painting the frame it showed")
	}

	open = true
	tt.Frame()
	if !tt.HasText("Settings") {
		t.Fatal("the modal did not reopen")
	}
	if tt.Focused("Open") {
		t.Fatal("the reopened modal left the focus on the page")
	}
	motionSettle(tt, dialogMotionDuration)
	reopened := motionShot(tt)
	motionSettle(tt, dialogMotionDuration)
	if again := motionShot(tt); motionDiff(reopened, again, area) != 0 {
		t.Error("the reopened modal does not come to rest once")
	}
}

func TestAlertDialogDepartureKeepsChoiceAndFocus(t *testing.T) {
	var open bool
	cancelled, accepted := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Delete file").OnClick(func() { open = true })
		AlertDialog(c, &open, AlertDialogOptions{
			Title:       "Delete report?",
			Description: "This cannot be undone.",
			CancelLabel: "Cancel",
			ActionLabel: "Delete",
			Destructive: true,
			OnCancel:    func() { cancelled++ },
			OnAction:    func() { accepted++ },
		})
	}, 640, 480)

	if err := tt.Click("Delete file"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Cancel") {
		t.Fatal("Cancel did not take the initial focus")
	}
	motionSettle(tt, dialogMotionDuration)
	title := motionLabel(t, tt, "Delete report?")
	settled := motionShot(tt)
	cancel, ok := tt.Find("Cancel")
	if !ok {
		t.Fatal("the alert has no cancel action")
	}

	// An ignored dismissal neither closes the alert nor starts its departure.
	tt.Key(0, ui.KeyEscape)
	tt.ClickAt(4, 4)
	tt.Frame()
	if !open || cancelled != 0 || accepted != 0 {
		t.Fatalf("an ignored dismissal changed the alert: open=%v cancelled=%d accepted=%d", open, cancelled, accepted)
	}
	tt.Advance(dialogMotionDuration / 2)
	if motionDiff(settled, motionShot(tt), title) != 0 {
		t.Error("an ignored dismissal animated the alert away")
	}

	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if open || cancelled != 1 || accepted != 0 {
		t.Fatalf("Cancel: open=%v cancelled=%d accepted=%d", open, cancelled, accepted)
	}
	if !tt.Focused("Delete file") {
		t.Fatal("the alert did not give the focus back")
	}
	// The copy still painting the alert is inert: a press on it repeats no
	// choice.
	tt.ClickAt(cancel.X+cancel.W/2, cancel.Y+cancel.H/2)
	tt.Advance(dialogMotionDuration)
	if cancelled != 1 || accepted != 0 {
		t.Fatalf("a press on the departing alert repeated the choice: cancelled=%d accepted=%d", cancelled, accepted)
	}
}

func TestReducedMotionSurfacesSettleAtOnce(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Modal(c, &open, ModalOptions{
			Title: "Settings",
			Body:  func() { ui.Button(c, "Confirm") },
		})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Settings") {
		t.Fatal("the modal did not open")
	}
	area := motionLabel(t, tt, "Settings")
	opened := motionShot(tt)
	if tt.FramePending() {
		t.Error("reduced motion asked for frames as the modal opened")
	}
	motionSettle(tt, dialogMotionDuration)
	if settled := motionShot(tt); motionDiff(opened, settled, area) != 0 {
		t.Error("reduced motion animated the modal in")
	}

	tt.Key(0, ui.KeyEscape)
	closed := motionShot(tt)
	if tt.FramePending() {
		t.Error("reduced motion asked for frames as the modal closed")
	}
	tt.Advance(dialogMotionDuration)
	if later := motionShot(tt); motionDiff(closed, later, area) != 0 {
		t.Error("reduced motion kept painting after the modal closed")
	}
	if tt.HasText("Settings") {
		t.Fatal("the modal stayed open")
	}
}
