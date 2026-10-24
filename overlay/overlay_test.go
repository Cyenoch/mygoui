package overlay

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestModalEscapeRestoresFocusAndBackdropDismisses(t *testing.T) {
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
				Footer: func() {
					ui.Button(c, "Done").OnClick(func() { open = false })
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

	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	tt.ClickAt(4, 4)
	if open {
		t.Fatal("backdrop press did not dismiss the modal")
	}
	if !tt.Focused("Open") {
		t.Fatal("backdrop dismissal did not restore focus")
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

func TestAlertDialogIgnoresBackdropAndEscape(t *testing.T) {
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
	tt.ClickAt(4, 4)
	tt.Frame()
	if !open || accepted != 0 || cancelled != 0 {
		t.Fatalf("backdrop press changed the alert: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if !open || accepted != 0 || cancelled != 0 {
		t.Fatalf("Escape changed the alert: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
	}
	if !tt.Focused("Cancel") {
		t.Fatal("the ignored dismissal lost the alert's focus")
	}
	if err := tt.Click("Cancel"); err != nil {
		t.Fatal(err)
	}
	if open || cancelled != 1 || accepted != 0 {
		t.Fatalf("Cancel: open=%v accepted=%d cancelled=%d", open, accepted, cancelled)
	}
}

// TestDrawerEdgesAndBottomSheet pins the floating edge geometry: 16 DIP
// gutters on every side, full rounding, a left or right panel keeping w-96
// at the sm breakpoint, and a top or bottom one spanning the inset width.
func TestDrawerEdgesAndBottomSheet(t *testing.T) {
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
			tt.Frame()
			tt.Frame()
			switch tc.side {
			case DrawerLeft:
				if panel.X != 16 || panel.Y != 16 || panel.W != 384 || panel.H != winH-32 {
					t.Fatalf("left drawer geometry = %+v", panel)
				}
			case DrawerRight:
				if panel.X+panel.W != winW-16 || panel.Y != 16 || panel.W != 384 || panel.H != winH-32 {
					t.Fatalf("right drawer geometry = %+v", panel)
				}
			case DrawerTop:
				if panel.X != 16 || panel.Y != 16 || panel.W != winW-32 || panel.H > (winH-32)*0.85+1 {
					t.Fatalf("top drawer geometry = %+v", panel)
				}
			case DrawerBottom:
				if panel.X != 16 || panel.Y+panel.H != winH-16 || panel.W != winW-32 || panel.H > (winH-32)*0.85+1 {
					t.Fatalf("bottom drawer geometry = %+v", panel)
				}
			}
		})
	}

	open := true
	var sheet ui.Rect
	st := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if p := BottomSheet(c, &open, BottomSheetOptions{Title: "Actions"}); p.Valid() {
			sheet = p.Bounds()
		}
	}, winW, winH)
	st.Frame()
	// The sheet slides in from below the window; the tester's clock stands
	// still, so settle the arrival before reading its resting box.
	st.Advance(drawerMotionDuration)
	st.Frame()
	if sheet.X != 16 || sheet.W != winW-32 || sheet.Y+sheet.H != winH-16 ||
		sheet.H > (winH-32)*0.85+1 {
		t.Fatalf("bottom sheet geometry = %+v, want the floating panel inside 16 DIP gutters", sheet)
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

func TestTooltipKeyboardRouteAndEscape(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		anchor := ui.Button(c, "Hover me")
		Tooltip(c, anchor, TooltipOptions{Text: "Help"})
	}, 640, 480)

	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if !tt.Focused("Hover me") {
		t.Fatal("the anchor did not take keyboard focus")
	}
	if !tt.HasText("Help") {
		t.Fatal("the keyboard focus route did not show the tooltip")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Help") {
		t.Fatal("Escape did not dismiss the tooltip")
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

func TestToastReducedMotionTimeoutAndClose(t *testing.T) {
	var state ToastState
	var ctx *ui.Context
	closed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ctx = c
		ui.Column(c).Padding(16).Children(func() {
			ui.Button(c, "Show saved").OnClick(func() {
				Toast(ctx, &state, ToastOptions{
					Title:   "File saved",
					Timeout: time.Second,
					OnClose: func() { closed++ },
				})
			})
			state.View(c)
		})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

	if err := tt.Click("Show saved"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("File saved") {
		t.Fatal("reduced-motion toast did not settle into view")
	}
	if tt.FramePending() {
		t.Fatal("reduced motion still animates the toast entry")
	}
	tt.Advance(1100 * time.Millisecond)
	tt.Frame()
	if tt.HasText("File saved") || closed != 1 || len(state.Items()) != 0 {
		t.Fatalf("timeout: shown=%v closed=%d retained=%d", tt.HasText("File saved"), closed, len(state.Items()))
	}

	if err := tt.Click("Show saved"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("File saved") || closed != 2 {
		t.Fatalf("close button: shown=%v closed=%d", tt.HasText("File saved"), closed)
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
					Action: "Undo", Timeout: -1,
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
