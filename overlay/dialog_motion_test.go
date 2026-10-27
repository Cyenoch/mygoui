package overlay

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func dialogMotionMarker(tt *ui.Tester) ui.Rect {
	img := tt.Image()
	minX, minY := img.Bounds().Max.X, img.Bounds().Max.Y
	maxX, maxY := -1, -1
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			pixel := img.RGBAAt(x, y)
			if int(pixel.R)-int(pixel.G) < 20 || int(pixel.R)-int(pixel.B) < 20 {
				continue
			}
			minX, minY = min(minX, x), min(minY, y)
			maxX, maxY = max(maxX, x), max(maxY, y)
		}
	}
	if maxX < 0 {
		return ui.Rect{}
	}
	return ui.Rect{X: float32(minX), Y: float32(minY), W: float32(maxX - minX + 1), H: float32(maxY - minY + 1)}
}

func TestResponsiveDialogMotion(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, tc := range []struct {
			name   string
			width  int
			mobile bool
		}{{"desktop", 800, false}, {"desktop breakpoint", 640, false}, {"mobile breakpoint", 639, true}, {"mobile", 320, true}} {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				var open bool
				actions, dismissals := 0, 0
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Column(c).Fill().Background(ui.RGB(240, 240, 240)).Children(func() {
						ui.Button(c, "Open").OnClick(func() { open = true })
						body := func() {
							ui.Box(c).Size(200, 64).Label("Marker").Background(ui.RGB(240, 20, 20))
							ui.Button(c, "Inside").OnClick(func() { actions++ })
						}
						if kind == "alert" {
							AlertDialog(c, &open, AlertDialogOptions{
								Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
								OnCancel: func() { dismissals++ },
							})
						} else {
							Dialog(c, &open, DialogOptions{
								Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque,
								OnDismiss: func() { dismissals++ },
							})
						}
					})
				}, tc.width, 480)
				unshaded := tt.Image().RGBAAt(4, 100)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				first, ok := tt.Find("Marker")
				if !ok || !tt.FramePending() {
					t.Fatal("opening did not build an animating marker")
				}
				tt.Advance(dialogMotionDuration / 2)
				middle, _ := tt.Find("Marker")
				halfShade := tt.Image().RGBAAt(4, 100)
				motionSettle(tt, dialogMotionDuration)
				settled, _ := tt.Find("Marker")
				fullShade := tt.Image().RGBAAt(4, 100)
				if unshaded.R <= halfShade.R || halfShade.R <= fullShade.R {
					t.Fatalf("scrim did not fade in: clear=%v half=%v full=%v", unshaded, halfShade, fullShade)
				}
				if tc.mobile {
					if first.Y-settled.Y < 14 || middle.Y <= settled.Y || middle.Y >= first.Y {
						t.Fatalf("mobile entry did not translate from below: first=%+v middle=%+v settled=%+v", first, middle, settled)
					}
				} else if first.W >= middle.W || middle.W >= settled.W || abs32(first.W/settled.W-0.98) > 0.003 {
					t.Fatalf("desktop entry did not scale from 98%%: first=%+v middle=%+v settled=%+v", first, middle, settled)
				}
				painted := dialogMotionMarker(tt)
				if painted.W == 0 {
					t.Fatal("settled marker was not painted")
				}
				inside, _ := tt.Find("Inside")
				if kind == "alert" {
					if !tt.Focused("Cancel") {
						t.Fatal("alert did not focus Cancel")
					}
					tt.Key(0, ui.KeyEnter)
				} else {
					tt.Key(0, ui.KeyEscape)
				}
				if open || dismissals != 1 || !tt.Focused("Open") {
					t.Fatalf("dismissal: open=%v dismissals=%d opener focused=%v", open, dismissals, tt.Focused("Open"))
				}
				tt.ClickAt(inside.X+inside.W/2, inside.Y+inside.H/2)
				if actions != 0 {
					t.Fatal("departing snapshot accepted an action")
				}
				tt.Advance(dialogMotionDuration / 2)
				leaving := dialogMotionMarker(tt)
				exitShade := tt.Image().RGBAAt(4, 100)
				if exitShade.R <= fullShade.R || exitShade.R >= unshaded.R {
					t.Fatalf("scrim did not fade out: full=%v half=%v clear=%v", fullShade, exitShade, unshaded)
				}
				if tc.mobile {
					if leaving.W == 0 || leaving.Y <= painted.Y+3 {
						t.Fatalf("mobile exit did not move down: settled=%+v leaving=%+v", painted, leaving)
					}
				} else if leaving.W == 0 || leaving.W >= painted.W || leaving.X <= painted.X {
					t.Fatalf("desktop exit did not scale down: settled=%+v leaving=%+v", painted, leaving)
				}
				motionSettle(tt, dialogMotionDuration)
				if remaining := dialogMotionMarker(tt); remaining.W != 0 || tt.Image().RGBAAt(4, 100) != unshaded {
					t.Fatalf("departure left paint: marker=%+v background=%v", remaining, tt.Image().RGBAAt(4, 100))
				}
			})
		}
	}
}

func TestResponsiveDialogDismissDuringEntry(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, width := range []int{320, 800} {
			name := "mobile"
			if width >= 640 {
				name = "desktop"
			}
			t.Run(kind+"/"+name, func(t *testing.T) {
				var open bool
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Box(c).Fill().Background(ui.RGB(240, 240, 240)).Children(func() {
						ui.Button(c, "Open").OnClick(func() { open = true })
						body := func() { ui.Box(c).Size(200, 64).Background(ui.RGB(240, 20, 20)) }
						if kind == "alert" {
							AlertDialog(c, &open, AlertDialogOptions{
								Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
							})
						} else {
							Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque})
						}
					})
				}, width, 480)
				clear := tt.Image().RGBAAt(4, 100)
				dismiss := func() {
					if kind == "alert" {
						tt.Key(0, ui.KeyEnter)
					} else {
						tt.Key(0, ui.KeyEscape)
					}
				}
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				motionSettle(tt, dialogMotionDuration)
				rest := dialogMotionMarker(tt)
				dismiss()
				motionSettle(tt, dialogMotionDuration)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				tt.Advance(dialogMotionDuration / 2)
				entering := dialogMotionMarker(tt)
				dismiss()
				atClose := dialogMotionMarker(tt)
				if entering != atClose {
					t.Fatalf("dismissal jumped from the entry frame: entry=%+v close=%+v", entering, atClose)
				}
				tt.Advance(dialogMotionDuration / 2)
				leaving := dialogMotionMarker(tt)
				if width < 640 {
					if entering.Y <= rest.Y || leaving.Y <= entering.Y || leaving.Y > rest.Y+15 {
						t.Fatalf("interrupted mobile motion overshot: rest=%+v entering=%+v leaving=%+v", rest, entering, leaving)
					}
				} else if leaving.W == 0 || leaving.W > entering.W || leaving.W < rest.W-4 {
					t.Fatalf("interrupted desktop motion overshot: rest=%+v entering=%+v leaving=%+v", rest, entering, leaving)
				}
				motionSettle(tt, dialogMotionDuration)
				if open || !tt.Focused("Open") || dialogMotionMarker(tt).W != 0 || tt.Image().RGBAAt(4, 100) != clear {
					t.Fatal("dismissal during entry left the dialog or scrim behind")
				}
			})
		}
	}
}

func TestResponsiveDialogDesktopExitScalesAroundPanel(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		t.Run(kind, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				body := func() { ui.Box(c).Size(200, 64).Background(ui.RGB(240, 20, 20)) }
				if kind == "alert" {
					AlertDialog(c, &open, AlertDialogOptions{
						Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
					})
				} else {
					Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque})
				}
			}, 800, 1200)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			motionSettle(tt, dialogMotionDuration)
			panel, _ := tt.Find("Panel")
			rest := dialogMotionMarker(tt)
			if kind == "alert" {
				tt.Key(0, ui.KeyEnter)
			} else {
				tt.Key(0, ui.KeyEscape)
			}
			tt.Advance(dialogMotionDuration / 2)
			leaving := dialogMotionMarker(tt)
			if leaving.W == 0 || leaving.W >= rest.W {
				t.Fatalf("exit did not scale: rest=%+v leaving=%+v", rest, leaving)
			}
			centerY := panel.Y + panel.H/2
			wantY := centerY + (rest.Y-centerY)*0.99
			if abs32(leaving.Y-wantY) > 1.2 {
				t.Fatalf("exit scaled around the wrong origin: panel=%+v rest=%+v leaving=%+v want y=%v", panel, rest, leaving, wantY)
			}
		})
	}
}

func TestResponsiveDialogReopensAfterInterruptedEntry(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, width := range []int{320, 800} {
			name := "mobile"
			if width >= 640 {
				name = "desktop"
			}
			t.Run(kind+"/"+name, func(t *testing.T) {
				var open bool
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Button(c, "Open").OnClick(func() { open = true })
					body := func() { ui.Box(c).Size(200, 64).Background(ui.RGB(240, 20, 20)) }
					if kind == "alert" {
						AlertDialog(c, &open, AlertDialogOptions{
							Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
						})
					} else {
						Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque})
					}
				}, width, 480)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				tt.Advance(dialogMotionDuration / 2)
				if kind == "alert" {
					tt.Key(0, ui.KeyEnter)
				} else {
					tt.Key(0, ui.KeyEscape)
				}
				tt.Advance(dialogMotionDuration / 4)
				leaving := dialogMotionMarker(tt)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				returning := dialogMotionMarker(tt)
				if leaving.W == 0 || abs32(returning.X-leaving.X) > 1 ||
					abs32(returning.Y-leaving.Y) > 1 || abs32(returning.W-leaving.W) > 1 {
					t.Fatalf("reopening jumped from interrupted exit: leaving=%+v returning=%+v", leaving, returning)
				}
				motionSettle(tt, 2*dialogMotionDuration)
				if !open || !tt.HasText("Panel") || tt.Focused("Open") {
					t.Fatal("reopened dialog did not settle with focus inside")
				}
			})
		}
	}
}

func TestResponsiveDialogRecloseWhileReturning(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		t.Run(kind, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				body := func() { ui.Box(c).Size(200, 64).Label("Marker").Background(ui.RGB(240, 20, 20)) }
				if kind == "alert" {
					AlertDialog(c, &open, AlertDialogOptions{
						Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
					})
				} else {
					Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque})
				}
			}, 320, 480)
			dismiss := func() {
				if kind == "alert" {
					tt.Key(0, ui.KeyEnter)
				} else {
					tt.Key(0, ui.KeyEscape)
				}
			}
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			motionSettle(tt, dialogMotionDuration)
			rest := dialogMotionMarker(tt)
			dismiss()
			tt.Advance(dialogMotionDuration / 2)
			leaving := dialogMotionMarker(tt)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			reopened := dialogMotionMarker(tt)
			if leaving.W == 0 || abs32(reopened.Y-leaving.Y) > 1 {
				t.Fatalf("reopen jumped from the exit: leaving=%+v reopened=%+v", leaving, reopened)
			}
			tt.Advance(dialogMotionDuration / 4)
			returning := dialogMotionMarker(tt)
			if returning.Y <= rest.Y || returning.Y >= reopened.Y {
				t.Fatalf("reopened panel did not return: rest=%+v reopened=%+v returning=%+v", rest, reopened, returning)
			}
			dismiss()
			if closed := dialogMotionMarker(tt); abs32(closed.Y-returning.Y) > 1 || closed.W == 0 {
				t.Fatalf("reclose jumped from the returning frame: returning=%+v closed=%+v", returning, closed)
			}
			tt.Advance(dialogMotionDuration / 2)
			finalExit := dialogMotionMarker(tt)
			if finalExit.W == 0 || finalExit.Y <= returning.Y || finalExit.Y > rest.Y+15 {
				t.Fatalf("reclose did not exit toward its target: rest=%+v returning=%+v exiting=%+v", rest, returning, finalExit)
			}
			motionSettle(tt, dialogMotionDuration)
			if open || tt.HasText("Panel") || !tt.Focused("Open") {
				t.Fatal("reclosed dialog did not finish or restore focus")
			}
		})
	}
}

func TestResponsiveDialogExternalCloseAfterReopen(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, delay := range []time.Duration{0, 20 * time.Millisecond, dialogMotionDuration / 4} {
			t.Run(kind+"/"+delay.String(), func(t *testing.T) {
				var open bool
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Button(c, "Open").OnClick(func() { open = true })
					body := func() { ui.Box(c).Size(200, 64).Background(ui.RGB(240, 20, 20)) }
					if kind == "alert" {
						AlertDialog(c, &open, AlertDialogOptions{
							Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
						})
					} else {
						Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque})
					}
				}, 320, 480)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				motionSettle(tt, dialogMotionDuration)
				rest := dialogMotionMarker(tt)
				open = false
				tt.Frame()
				tt.Advance(dialogMotionDuration / 2)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				if delay > 0 {
					tt.Advance(delay)
				}
				returning := dialogMotionMarker(tt)
				open = false
				tt.Frame()
				closed := dialogMotionMarker(tt)
				if delay == 0 {
					// MyGo has not yet settled the returning backdrop, so a
					// snapshot here would jump to its retained position.
					if closed.W != 0 || !tt.Focused("Open") {
						t.Fatalf("same-frame close left a displaced panel or focus: panel=%+v", closed)
					}
					return
				}
				if closed.W == 0 || abs32(closed.Y-returning.Y) > 1 {
					t.Fatalf("external close jumped: returning=%+v closed=%+v", returning, closed)
				}
				tt.Advance(dialogMotionDuration / 2)
				leaving := dialogMotionMarker(tt)
				if leaving.W == 0 || leaving.Y <= returning.Y || leaving.Y > rest.Y+15 {
					t.Fatalf("external close did not animate out: rest=%+v returning=%+v leaving=%+v", rest, returning, leaving)
				}
				motionSettle(tt, dialogMotionDuration)
				if open || dialogMotionMarker(tt).W != 0 || !tt.Focused("Open") {
					t.Fatal("external close left the dialog or focus behind")
				}
			})
		}
	}
}

func TestResponsiveDialogDesktopRecloseWhileReturning(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		t.Run(kind, func(t *testing.T) {
			var open bool
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				body := func() { ui.Box(c).Size(200, 64).Background(ui.RGB(240, 20, 20)) }
				if kind == "alert" {
					AlertDialog(c, &open, AlertDialogOptions{
						Title: "Panel", Body: body, CancelLabel: "Cancel", Backdrop: ModalBackdropOpaque,
					})
				} else {
					Dialog(c, &open, DialogOptions{Title: "Panel", Body: body, Backdrop: ModalBackdropOpaque})
				}
			}, 800, 480)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			motionSettle(tt, dialogMotionDuration)
			rest := dialogMotionMarker(tt)
			open = false
			tt.Frame()
			tt.Advance(dialogMotionDuration / 2)
			leaving := dialogMotionMarker(tt)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			reopened := dialogMotionMarker(tt)
			if abs32(reopened.X-leaving.X) > 1 || abs32(reopened.Y-leaving.Y) > 1 ||
				abs32(reopened.W-leaving.W) > 1 {
				t.Fatalf("reopen jumped from the exit: leaving=%+v reopened=%+v", leaving, reopened)
			}
			tt.Advance(dialogMotionDuration / 4)
			returning := dialogMotionMarker(tt)
			open = false
			tt.Frame()
			closed := dialogMotionMarker(tt)
			if closed.W == 0 || abs32(closed.X-returning.X) > 1 ||
				abs32(closed.Y-returning.Y) > 1 || abs32(closed.W-returning.W) > 1 {
				t.Fatalf("reclose jumped from the returning frame: returning=%+v closed=%+v", returning, closed)
			}
			tt.Advance(dialogMotionDuration / 2)
			finalExit := dialogMotionMarker(tt)
			if finalExit.W == 0 || finalExit.W > returning.W || finalExit.W < rest.W-5 {
				t.Fatalf("reclose did not scale toward its target: rest=%+v returning=%+v leaving=%+v", rest, returning, finalExit)
			}
			motionSettle(tt, dialogMotionDuration)
			if open || dialogMotionMarker(tt).W != 0 || !tt.Focused("Open") {
				t.Fatal("reclosed dialog did not finish or restore focus")
			}
		})
	}
}

func TestResponsiveDialogReducedMotion(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, width := range []int{320, 800} {
			name := "mobile"
			if width >= 640 {
				name = "desktop"
			}
			t.Run(kind+"/"+name, func(t *testing.T) {
				var open bool
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Button(c, "Open").OnClick(func() { open = true })
					if kind == "alert" {
						AlertDialog(c, &open, AlertDialogOptions{Title: "Panel", CancelLabel: "Cancel"})
					} else {
						Dialog(c, &open, DialogOptions{Title: "Panel"})
					}
				}, width, 480)
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				first := motionShot(tt)
				if tt.FramePending() {
					t.Fatal("reduced-motion arrival requested animation frames")
				}
				tt.Advance(dialogMotionDuration)
				if motionDiff(first, motionShot(tt), first.Bounds()) != 0 {
					t.Fatal("reduced-motion arrival moved")
				}
				if kind == "alert" {
					tt.Key(0, ui.KeyEnter)
				} else {
					tt.Key(0, ui.KeyEscape)
				}
				closed := motionShot(tt)
				if open || tt.FramePending() || !tt.Focused("Open") {
					t.Fatal("reduced-motion dismissal did not finish immediately")
				}
				tt.Advance(dialogMotionDuration)
				if motionDiff(closed, motionShot(tt), closed.Bounds()) != 0 {
					t.Fatal("reduced-motion departure kept painting")
				}
			})
		}
	}
}

func TestResponsiveDialogBackdropDuringEntry(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		for _, width := range []int{320, 800} {
			name := "mobile"
			if width >= 640 {
				name = "desktop"
			}
			t.Run(kind+"/"+name, func(t *testing.T) {
				var open bool
				openings, dismissals := 0, 0
				tt := ui.NewTester(func(c *ui.Context) {
					theme.Apply(c)
					ui.Button(c, "Open").OnClick(func() { open = true; openings++ })
					if kind == "alert" {
						AlertDialog(c, &open, AlertDialogOptions{
							Title: "Panel", CancelLabel: "Cancel", OnCancel: func() { dismissals++ },
						})
					} else {
						Dialog(c, &open, DialogOptions{Title: "Panel", OnDismiss: func() { dismissals++ }})
					}
				}, width, 480)
				if err := tt.Click("Open"); err != nil {
					t.Fatal(err)
				}
				tt.ClickAt(1, 1)
				if openings != 1 {
					t.Fatalf("backdrop leaked a press to the opener: openings=%d", openings)
				}
				if kind == "alert" {
					if !open || dismissals != 0 || !tt.Focused("Cancel") {
						t.Fatalf("alert accepted backdrop dismissal: open=%v dismissals=%d", open, dismissals)
					}
					tt.Key(0, ui.KeyEnter)
				}
				if open || dismissals != 1 || !tt.Focused("Open") {
					t.Fatalf("dismissal: open=%v dismissals=%d opener focused=%v", open, dismissals, tt.Focused("Open"))
				}
			})
		}
	}
}

func TestResponsiveDialogResizesAndReopensDuringExit(t *testing.T) {
	for _, kind := range []string{"dialog", "alert"} {
		t.Run(kind, func(t *testing.T) {
			var open bool
			clicked := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				body := func() {
					ui.Box(c).Size(200, 64).Background(ui.RGB(240, 20, 20))
					for range 30 {
						ui.Text(c, "Scrollable row")
					}
					ui.Button(c, "Inside").OnClick(func() { clicked++ })
				}
				if kind == "alert" {
					AlertDialog(c, &open, AlertDialogOptions{Title: "Panel", Body: body, CancelLabel: "Cancel"})
				} else {
					Dialog(c, &open, DialogOptions{Title: "Panel", Body: body})
				}
			}, 800, 480)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Advance(50 * time.Millisecond)
			tt.SetSize(320, 480)
			motionSettle(tt, dialogMotionDuration)
			mobile, _ := tt.Find("Panel")
			if mobile.X != 0 || mobile.W != 320 || mobile.Y+mobile.H != 480 {
				t.Fatalf("resize to mobile lost placement: %+v", mobile)
			}
			tt.SetSize(800, 480)
			motionSettle(tt, dialogMotionDuration)
			desktop, _ := tt.Find("Panel")
			if desktop.W != 512 || desktop.X != 144 {
				t.Fatalf("resize to desktop lost placement: %+v", desktop)
			}
			if kind == "alert" {
				tt.Key(0, ui.KeyEnter)
			} else {
				tt.Key(0, ui.KeyEscape)
			}
			tt.Advance(dialogMotionDuration / 4)
			leaving := dialogMotionMarker(tt)
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			returning := dialogMotionMarker(tt)
			if leaving.W == 0 || abs32(returning.X-leaving.X) > 1 ||
				abs32(returning.Y-leaving.Y) > 1 || abs32(returning.W-leaving.W) > 1 {
				t.Fatalf("reopening jumped from the departing frame: leaving=%+v returning=%+v", leaving, returning)
			}
			motionSettle(tt, 2*dialogMotionDuration)
			reopened, _ := tt.Find("Panel")
			if reopened != desktop {
				t.Fatalf("reopened panel did not settle: got=%+v want=%+v", reopened, desktop)
			}
			tt.Scroll(reopened.X+100, reopened.Y+100, 0, 2000)
			tt.Frame()
			if err := tt.Click("Inside"); err != nil {
				t.Fatal(err)
			}
			if !open || clicked != 1 {
				t.Fatalf("reopened interaction: open=%v clicks=%d", open, clicked)
			}
			if kind == "alert" {
				if err := tt.Click("Cancel"); err != nil {
					t.Fatal(err)
				}
			} else {
				tt.Key(0, ui.KeyEscape)
			}
			tt.Advance(dialogMotionDuration / 4)
			tt.SetSize(320, 480)
			motionSettle(tt, 2*dialogMotionDuration)
			if open || tt.HasText("Panel") || tt.FramePending() || !tt.Focused("Open") {
				t.Fatal("resize during exit left the dialog, its animation, or focus behind")
			}
		})
	}
}
