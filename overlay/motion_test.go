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

// TestSurfacesPaintArrivalAndDeparture drives every fading surface through
// the same paint contract: the arrival asks for frames and paints its way in,
// it reaches its resting paint, the departure starts from the pose last shown,
// paints its way out and leaves nothing behind. The sliding drawers, sheet
// and bottom sheet have their own motion tests.
func TestSurfacesPaintArrivalAndDeparture(t *testing.T) {
	clickOpener := func(label string) func(*testing.T, *ui.Tester) {
		return func(t *testing.T, tt *ui.Tester) {
			if err := tt.Click(label); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, tc := range []struct {
		name     string
		duration time.Duration
		label    string
		build    func(c *ui.Context, open *bool)
		open     func(t *testing.T, tt *ui.Tester)
		focus    string
	}{
		{"modal", dialogMotionDuration, "Settings", func(c *ui.Context, open *bool) {
			ui.Button(c, "Open").OnClick(func() { *open = true })
			Modal(c, open, ModalOptions{
				Title: "Settings",
				Body:  func() { ui.Button(c, "Confirm") },
			})
		}, clickOpener("Open"), "Open"},
		{"edge dialog", dialogMotionDuration, "Panel", func(c *ui.Context, open *bool) {
			ui.Button(c, "Open").OnClick(func() { *open = true })
			Dialog(c, open, DialogOptions{
				Placement: mygoui.Placement{Side: mygoui.SideRight},
				Title:     "Panel",
			})
		}, clickOpener("Open"), "Open"},
		{"popover", popupMotionDuration, "Actions", func(c *ui.Context, open *bool) {
			ui.Column(c).Padding(16).Children(func() {
				anchor := ui.Button(c, "Options").OnClick(func() { *open = !*open })
				Popover(c, anchor, open, PopoverOptions{
					Title: "Actions",
					Body:  func() { ui.Button(c, "Rename") },
				})
			})
		}, clickOpener("Options"), "Options"},
		{"tooltip", popupMotionDuration, "Help", func(c *ui.Context, open *bool) {
			anchor := ui.Button(c, "Hover me")
			Tooltip(c, anchor, TooltipOptions{Text: "Help"})
		}, func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyTab)
			tt.Frame()
		}, "Hover me"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				tc.build(c, &open)
			}, 640, 480)
			tc.open(t, tt)
			if !tt.HasText(tc.label) {
				t.Fatalf("the surface did not open; labels=%q", tt.Texts())
			}
			area := motionLabel(t, tt, tc.label)
			motionAssertSurface(t, tt, area, tc.duration, func() {
				tt.Key(0, ui.KeyEscape)
				if open || tt.HasText(tc.label) {
					t.Fatal("Escape did not close the surface")
				}
			})
			if !tt.Focused(tc.focus) {
				t.Fatal("the surface did not give the focus back")
			}
		})
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

// TestReducedMotionSurfacesArriveAtRest checks that the desktop's reduce-motion
// preference puts a surface straight at its resting paint: no frames are
// requested, nothing moves afterwards, and the close releases it at once
// instead of animating it away.
func TestReducedMotionSurfacesArriveAtRest(t *testing.T) {
	clickOpener := func(label string) func(*testing.T, *ui.Tester) {
		return func(t *testing.T, tt *ui.Tester) {
			if err := tt.Click(label); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, tc := range []struct {
		name     string
		duration time.Duration
		label    string
		build    func(c *ui.Context, open *bool)
		open     func(t *testing.T, tt *ui.Tester)
	}{
		{"modal", dialogMotionDuration, "Settings", func(c *ui.Context, open *bool) {
			ui.Button(c, "Open").OnClick(func() { *open = true })
			Modal(c, open, ModalOptions{
				Title: "Settings",
				Body:  func() { ui.Button(c, "Confirm") },
			})
		}, clickOpener("Open")},
		{"bottom sheet", drawerMotionDuration, "Actions", func(c *ui.Context, open *bool) {
			ui.Button(c, "Open").OnClick(func() { *open = true })
			BottomSheet(c, open, BottomSheetOptions{Title: "Actions"})
		}, clickOpener("Open")},
		{"tooltip", popupMotionDuration, "Help", func(c *ui.Context, open *bool) {
			anchor := ui.Button(c, "Hover me")
			Tooltip(c, anchor, TooltipOptions{Text: "Help"})
		}, func(t *testing.T, tt *ui.Tester) {
			tt.Key(0, ui.KeyTab)
			tt.Frame()
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				tc.build(c, &open)
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

			tc.open(t, tt)
			if !tt.HasText(tc.label) {
				t.Fatalf("the surface did not open; labels=%q", tt.Texts())
			}
			rest := motionLabel(t, tt, tc.label)
			area := rest
			arrived := motionShot(tt)
			if tt.FramePending() {
				t.Error("reduced motion asks for frames as the surface opens")
			}
			motionSettle(tt, tc.duration)
			if settled := motionShot(tt); motionDiff(arrived, settled, area) != 0 {
				t.Error("reduced motion animated the surface in")
			}
			if later := motionLabel(t, tt, tc.label); later != rest {
				t.Errorf("reduced motion moved the surface: arrived=%v settled=%v", rest, later)
			}

			tt.Key(0, ui.KeyEscape)
			closed := motionShot(tt)
			if tt.HasText(tc.label) {
				t.Fatal("Escape did not close the surface")
			}
			if tt.FramePending() {
				t.Error("reduced motion asks for frames as the surface closes")
			}
			motionSettle(tt, tc.duration)
			if later := motionShot(tt); motionDiff(closed, later, area) != 0 {
				t.Error("reduced motion kept painting after the surface closed")
			}
		})
	}
}
