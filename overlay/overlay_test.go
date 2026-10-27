package overlay

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestModalTrapsTheFocusAndRestoresIt keeps Tab inside the modal while it is
// open and gives the focus back to the opener when Escape closes it. Backdrop
// dismissal is covered with its hit area, and the reopen during the exit has
// its own test.
func TestModalTrapsTheFocusAndRestoresIt(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			ui.Button(c, "Open").OnClick(func() { open = true })
			Modal(c, &open, ModalOptions{
				Title: "Settings",
				Body: func() {
					ui.Button(c, "Confirm")
				},
			})
		})
	}, 640, 480)

	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if !open || !tt.HasText("Settings") {
		t.Fatal("modal did not open")
	}
	if tt.Focused("Open") {
		t.Fatal("focus stayed behind the modal")
	}
	for range 4 {
		tt.Key(0, ui.KeyTab)
		if tt.Focused("Open") {
			t.Fatal("Tab escaped the modal focus scope")
		}
	}
	tt.Key(0, ui.KeyEscape)
	if open || tt.HasText("Settings") {
		t.Fatal("Escape did not close the modal")
	}
	if !tt.Focused("Open") {
		t.Fatal("closing the modal did not restore focus to its opener")
	}
}

func TestModalPreventDismiss(t *testing.T) {
	var open bool
	dismissed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		Modal(c, &open, ModalOptions{
			Title:          "Locked",
			PreventDismiss: true,
			OnDismiss:      func() { dismissed++ },
		})
	}, 640, 480)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	tt.ClickAt(4, 4)
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if !open || dismissed != 0 {
		t.Fatalf("prevented dismissal closed the modal: open=%v dismissed=%d", open, dismissed)
	}
	tt.Click("Close")
	if open || dismissed != 1 {
		t.Fatalf("close button: open=%v dismissed=%d", open, dismissed)
	}
}

// TestAlertDialogKeepsItsChoiceAndIgnoresDismissal follows the alert through
// its life: Cancel takes the initial focus, neither the backdrop nor Escape
// answers the question or starts the exit, Cancel runs exactly once and gives
// the focus back, and the copy still painting the alert answers no press.
func TestAlertDialogKeepsItsChoiceAndIgnoresDismissal(t *testing.T) {
	var open bool
	accepted, cancelled := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Delete file").OnClick(func() { open = true })
		AlertDialog(c, &open, AlertDialogOptions{
			Title:       "Delete report?",
			Description: "This cannot be undone.",
			CancelLabel: "Cancel",
			ActionLabel: "Delete",
			Destructive: true,
			OnAction:    func() { accepted++ },
			OnCancel:    func() { cancelled++ },
		})
	}, 640, 480)

	if err := tt.Click("Delete file"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Cancel") {
		t.Fatal("Cancel must receive the initial focus")
	}
	motionSettle(tt, dialogMotionDuration)
	title := motionLabel(t, tt, "Delete report?")
	settled := motionShot(tt)
	cancel, ok := tt.Find("Cancel")
	if !ok {
		t.Fatal("the alert has no cancel action")
	}

	tt.ClickAt(4, 4)
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if !open || accepted != 0 || cancelled != 0 {
		t.Fatalf("a dismissal changed the alert: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
	}
	if !tt.Focused("Cancel") {
		t.Fatal("the ignored dismissal lost the alert's focus")
	}
	tt.Advance(dialogMotionDuration / 2)
	if motionDiff(settled, motionShot(tt), title) != 0 {
		t.Error("an ignored dismissal animated the alert away")
	}

	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if open || cancelled != 1 || accepted != 0 {
		t.Fatalf("Cancel: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
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

func TestAlertDialogActionRequiresCallbackAndHonorsDisabled(t *testing.T) {
	var open bool
	disabled, callback := true, false
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Open").OnClick(func() { open = true })
		o := AlertDialogOptions{
			Title: "Confirm?", ActionLabel: "Accept", ActionDisabled: disabled,
		}
		if callback {
			o.OnAction = func() { actions++ }
		}
		AlertDialog(c, &open, o)
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Accept") {
		t.Fatal("action without callback is visible")
	}
	callback = true
	tt.Frame()
	if err := tt.Click("Accept"); err != nil {
		t.Fatal(err)
	}
	if !open || actions != 0 {
		t.Fatalf("disabled action: open=%v actions=%d", open, actions)
	}
	disabled = false
	tt.Frame()
	if err := tt.Click("Accept"); err != nil {
		t.Fatal(err)
	}
	if open || actions != 1 || !tt.Focused("Open") {
		t.Fatalf("accept: open=%v actions=%d focused=%v", open, actions, tt.Focused("Open"))
	}
}

// TestDrawerEdgeGeometry checks all four edge-attached drawer placements.
func TestDrawerEdgeGeometry(t *testing.T) {
	const winW, winH = 640, 480
	for _, tc := range []struct {
		name string
		side DrawerSide
	}{
		{"left", DrawerLeft},
		{"right", DrawerRight},
		{"top", DrawerTop},
		{"bottom", DrawerBottom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			open := true
			var panel ui.Rect
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				if p := Drawer(c, &open, DrawerOptions{Side: tc.side, Title: "Filters"}); p.Valid() {
					panel = p.Bounds()
				}
			}, winW, winH)
			motionSettle(tt, 2*drawerMotionDuration)
			switch tc.side {
			case DrawerLeft:
				if panel.X != 0 || panel.Y != 0 || panel.W != 448 || panel.H != winH {
					t.Fatalf("left drawer geometry = %+v", panel)
				}
			case DrawerRight:
				if panel.X+panel.W != winW || panel.Y != 0 || panel.W != 448 || panel.H != winH {
					t.Fatalf("right drawer geometry = %+v", panel)
				}
			case DrawerTop:
				if panel.X != 0 || panel.Y != 0 || panel.W != winW || panel.H > winH-48 {
					t.Fatalf("top drawer geometry = %+v", panel)
				}
			case DrawerBottom:
				if panel.X != 0 || panel.Y+panel.H != winH || panel.W != winW || panel.H > winH-48 {
					t.Fatalf("bottom drawer geometry = %+v", panel)
				}
			}
		})
	}
}

func TestPopoverAnchorsAndDismisses(t *testing.T) {
	var open bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			anchor := ui.Button(c, "Open popover").OnClick(func() { open = true })
			Popover(c, anchor, &open, PopoverOptions{
				Title: "Details",
				Body: func() {
					ui.Button(c, "Inside")
				},
			})
		})
	}, 640, 480)

	if err := tt.Click("Open popover"); err != nil {
		t.Fatal(err)
	}
	trigger, ok := tt.Find("Open popover")
	if !ok {
		t.Fatal("trigger vanished")
	}
	panel, ok := tt.Find("Details")
	if !ok {
		t.Fatal("popover did not show")
	}
	if panel.Y < trigger.Y+trigger.H {
		t.Fatalf("popover at %v is not below the trigger at %v", panel, trigger)
	}
	tt.Key(0, ui.KeyEscape)
	if open || tt.HasText("Details") {
		t.Fatal("Escape did not dismiss the popover")
	}
	if !tt.Focused("Open popover") {
		t.Fatal("closing the popover did not restore focus")
	}

	if err := tt.Click("Open popover"); err != nil {
		t.Fatal(err)
	}
	tt.ClickAt(600, 300)
	if open || tt.HasText("Details") {
		t.Fatal("outside press did not dismiss the popover")
	}
}

func TestPopoverKeepsTallContentScrollable(t *testing.T) {
	var open bool
	var panel ui.Rect
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			anchor := ui.Button(c, "Open popover").OnClick(func() { open = true })
			if p := Popover(c, anchor, &open, PopoverOptions{
				Title: "Details", Width: 600,
				Body: func() {
					ui.Column(c).Gap(8).Children(func() {
						for i := range 12 {
							ui.Text(c, fmt.Sprintf("Row %d", i))
						}
						ui.Button(c, "Last action").OnClick(func() { actions++ })
					})
				},
			}); p.Valid() {
				panel = p.Bounds()
			}
		})
	}, 300, 240)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Open popover"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if panel.W > 284 || panel.H > 224 {
		t.Fatalf("popup exceeds the viewport: %+v", panel)
	}
	tt.Scroll(panel.X+panel.W/2, panel.Y+panel.H/2, 0, 1000)
	tt.Frame()
	if err := tt.Click("Last action"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 {
		t.Fatalf("scrolled action ran %d times", actions)
	}
}

func TestToastUpdatesCapsAndDismisses(t *testing.T) {
	var state ToastState
	var ctx *ui.Context
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ctx = c
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			ui.Button(c, "Show four").OnClick(func() {
				for i := range 4 {
					Toast(ctx, &state, ToastOptions{
						ID:      fmt.Sprintf("m%d", i),
						Title:   fmt.Sprintf("Message %d", i),
						Timeout: -1,
					})
				}
			})
			ui.Button(c, "Update").OnClick(func() {
				Toast(ctx, &state, ToastOptions{ID: "m3", Title: "Updated", Timeout: -1})
			})
			ui.Button(c, "Dismiss m3").OnClick(func() { state.Dismiss(ctx, "m3") })
			state.View(c)
		})
	}, 640, 480)

	if err := tt.Click("Show four"); err != nil {
		t.Fatal(err)
	}
	if len(state.Items()) != 4 {
		t.Fatalf("retained toasts = %d, want 4", len(state.Items()))
	}
	if tt.HasText("Message 0") {
		t.Fatal("the toast beyond the three-visible cap was painted")
	}
	for _, want := range []string{"Message 1", "Message 2", "Message 3"} {
		if !tt.HasText(want) {
			t.Fatalf("%q is missing from the visible toasts", want)
		}
	}

	if err := tt.Click("Update"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Message 3") || !tt.HasText("Updated") {
		t.Fatal("a same-ID update replaced the message instead of updating in place")
	}
	if items := state.Items(); len(items) != 4 || items[3].Title != "Updated" {
		t.Fatalf("state after update = %+v", items)
	}

	if err := tt.Click("Dismiss m3"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Updated") || len(state.Items()) != 3 {
		t.Fatalf("dismiss left %+v", state.Items())
	}
}

func TestToastActionRequiresCallback(t *testing.T) {
	var state ToastState
	actions, closed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Children(func() {
			ui.Button(c, "Show unavailable").OnClick(func() {
				Toast(c, &state, ToastOptions{ID: "notice", Title: "Notice", Action: "Undo", Timeout: -1})
			})
			ui.Button(c, "Show available").OnClick(func() {
				Toast(c, &state, ToastOptions{
					ID: "notice", Title: "Notice", Action: "Undo", Timeout: -1,
					OnAction: func() { actions++ }, OnClose: func() { closed++ },
				})
			})
			state.View(c)
		})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Show unavailable"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Undo") || state.Items()[0].Action != "" {
		t.Fatal("toast advertised an action without a callback")
	}
	if err := tt.Click("Show available"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 || closed != 1 || len(state.Items()) != 0 {
		t.Fatalf("action=%d close=%d items=%+v", actions, closed, state.Items())
	}
}

func TestToastActionUsesPrimaryOnHover(t *testing.T) {
	previous := theme.PreferencesSnapshot()
	t.Cleanup(func() { theme.Configure(previous) })
	primary := ui.RGBA(40, 80, 120, 0.5)
	popover := ui.RGB(240, 240, 240)
	theme.Configure(theme.Preferences{Palette: theme.PaletteOverrides{
		Light: &theme.PalettePatch{Primary: &primary, Popover: &popover},
	}})
	var state ToastState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Show").OnClick(func() {
			Toast(c, &state, ToastOptions{
				Title: "Notice", Action: "Undo", OnAction: func() {}, Timeout: -1,
			})
		})
		state.View(c)
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Show"); err != nil {
		t.Fatal(err)
	}
	action, ok := tt.Find("Undo")
	if !ok {
		t.Fatal("toast action is missing")
	}
	x, y := int(action.X+4), int(action.Y+action.H/2)
	normal := tt.Image().RGBAAt(x, y)
	if !sameColor(normal, primary.Over(popover), 4) {
		t.Fatalf("action color=%v, want primary over popover", normal)
	}
	tt.Move(float32(x), float32(y))
	tt.Frame()
	hover := tt.Image().RGBAAt(x, y)
	if !sameColor(hover, primary.Alpha(0.9).Over(popover), 4) {
		t.Fatalf("hover color=%v, want primary at 90%% opacity", hover)
	}
}

// toastTestView builds the toast harness used by the close-button tests: a
// page background distinct from the white toast surface, one trigger and the
// toast region.
func toastTestView(state *ToastState, ctx **ui.Context) func(c *ui.Context) {
	return func(c *ui.Context) {
		theme.Apply(c)
		*ctx = c
		t := theme.Colors(c)
		ui.Column(c).Fill().Background(t.Background).Children(func() {
			ui.Button(c, "Show").OnClick(func() {
				Toast(*ctx, state, ToastOptions{
					Title: "File saved", Description: "All changes are on disk.",
					Action: "Undo", OnAction: func() {}, Timeout: -1,
				})
			})
			state.View(c)
		})
	}
}

func TestToastCloseRevealsOnHoverAndDismisses(t *testing.T) {
	var state ToastState
	var ctx *ui.Context
	tt := ui.NewTester(toastTestView(&state, &ctx), 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Show"); err != nil {
		t.Fatal(err)
	}
	closeRect, ok := tt.Find("Close")
	if !ok {
		t.Fatal("toast shows no close button")
	}
	if darkPixels(tt.Image(), closeRect) != 0 {
		t.Fatal("close button paints before hover or focus")
	}
	tt.Move(closeRect.X+closeRect.W/2, closeRect.Y+closeRect.H/2)
	tt.Frame()
	if darkPixels(tt.Image(), closeRect) == 0 {
		t.Fatal("hovering the toast did not reveal the close button")
	}
	tt.Move(4, 4)
	tt.Frame()
	if darkPixels(tt.Image(), closeRect) != 0 {
		t.Fatal("close button stays painted after the pointer leaves")
	}
	tt.Move(closeRect.X+closeRect.W/2, closeRect.Y+closeRect.H/2)
	tt.Frame()
	if err := tt.Click("Close"); err != nil {
		t.Fatalf("revealed close button: %v", err)
	}
	if tt.HasText("File saved") || len(state.Items()) != 0 {
		t.Fatal("revealed close button did not dismiss the toast")
	}
}

// darkPixels counts the muted ink pixels of the close glyph in r, the paint
// the reveal adds.
func darkPixels(img *image.RGBA, r ui.Rect) int {
	count := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if !image.Pt(x, y).In(img.Bounds()) {
				continue
			}
			c := img.RGBAAt(x, y)
			if c.R < 200 && c.G < 200 && c.B < 200 {
				count++
			}
		}
	}
	return count
}

// TestToastDismissDuringExitKeepsPaintSafe dismisses a toast with both its
// close and its action while the exit transition paints the retained copy.
// The copy's paint used to read Elements from the build pass that made it,
// which panic ("element from pass N used in pass N+2") and killed the app.
func TestToastDismissDuringExitKeepsPaintSafe(t *testing.T) {
	var state ToastState
	var ctx *ui.Context
	tt := ui.NewTester(toastTestView(&state, &ctx), 640, 480)
	if err := tt.Click("Show"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(300 * time.Millisecond)
	tt.Frame()
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	// A frame of the exit transition paints the retained copy in a later
	// pass, where the Elements of the build that made it have expired.
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if len(state.Items()) != 0 {
		t.Fatalf("close left %d retained toasts", len(state.Items()))
	}

	if err := tt.Click("Show"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(300 * time.Millisecond)
	tt.Frame()
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if len(state.Items()) != 0 {
		t.Fatalf("action left %d retained toasts", len(state.Items()))
	}
}
