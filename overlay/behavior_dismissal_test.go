package overlay

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// A surface that must not be dismissed keeps the focus where it was: a
// blocked Escape or backdrop press that moved it out would leave the
// keyboard user behind a modal panel.
func TestPreventDismissHoldsTheSurfaceAndItsFocus(t *testing.T) {
	cases := map[string]func(c *ui.Context, open *bool, body func(), dismiss func()) ui.Element{
		"Modal": func(c *ui.Context, open *bool, body func(), dismiss func()) ui.Element {
			return Modal(c, open, ModalOptions{Title: "Hold", Body: body, PreventDismiss: true, OnDismiss: dismiss})
		},
		"Dialog": func(c *ui.Context, open *bool, body func(), dismiss func()) ui.Element {
			return Dialog(c, open, DialogOptions{Title: "Hold", Body: body, PreventDismiss: true, OnDismiss: dismiss})
		},
		"Drawer": func(c *ui.Context, open *bool, body func(), dismiss func()) ui.Element {
			return Drawer(c, open, DrawerOptions{Title: "Hold", Body: body, PreventDismiss: true, OnDismiss: dismiss})
		},
		"BottomSheet": func(c *ui.Context, open *bool, body func(), dismiss func()) ui.Element {
			return BottomSheet(c, open, BottomSheetOptions{Title: "Hold", Body: body, PreventDismiss: true, OnDismiss: dismiss})
		},
	}
	for name, build := range cases {
		t.Run(name, func(t *testing.T) {
			open, dismissed := false, 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Button(c, "Open").OnClick(func() { open = true })
				build(c, &open, func() { ui.Button(c, "Inside") }, func() { dismissed++ })
			}, 640, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if err := tt.Click("Open"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if err := tt.Click("Inside"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyEscape)
			tt.Frame()
			if !open || dismissed != 0 || !tt.Focused("Inside") {
				t.Fatalf("Escape: open=%v dismissed=%d, focus kept=%v", open, dismissed, tt.Focused("Inside"))
			}
			tt.ClickAt(4, 4)
			tt.Frame()
			if !open || dismissed != 0 || !tt.Focused("Inside") {
				t.Fatalf("backdrop: open=%v dismissed=%d, focus kept=%v", open, dismissed, tt.Focused("Inside"))
			}
		})
	}
}

// The first Escape empties the search, so a user who typed a query can clear
// it without losing the dialog; the second closes it and the focus returns.
func TestShortcutsDialogEscapeClearsTheSearchBeforeClosing(t *testing.T) {
	var state ShortcutsDialogState
	tt := shortcutsTester(&state, []ShortcutGroup{
		{Title: "General", Items: []Shortcut{{Label: "Toggle sidebar", Keys: []string{"⌘", "B"}}}},
	}, 800, 600)
	if err := tt.Click("Open"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Type("side")
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if !state.Open || state.Query != "" || !tt.Focused("Search") {
		t.Fatalf("first Escape: open=%v query=%q, search focused=%v", state.Open, state.Query, tt.Focused("Search"))
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open || !tt.Focused("Open") {
		t.Fatalf("second Escape: open=%v, focus back on the opener=%v", state.Open, tt.Focused("Open"))
	}
}

// The viewer starts on its close button and gives the focus back to the
// control that opened it, by Escape or by the close button.
func TestImageViewerFocusesCloseAndReturnsTheFocus(t *testing.T) {
	open, index, dismissed := false, 0, 0
	images := viewerImages()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "View").OnClick(func() { open = true })
		ImageViewer(c, &open, ImageViewerOptions{Images: images, Index: &index, OnDismiss: func() { dismissed++ }})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("View"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.Focused("Close") {
		t.Fatal("the viewer did not focus its close button")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || dismissed != 1 || !tt.Focused("View") {
		t.Fatalf("Escape: open=%v dismissed=%d, focus back=%v", open, dismissed, tt.Focused("View"))
	}
	if err := tt.Click("View"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open || dismissed != 2 || !tt.Focused("View") {
		t.Fatalf("close button: open=%v dismissed=%d, focus back=%v", open, dismissed, tt.Focused("View"))
	}
}

// A popover is non-modal: the press that closes it also reaches the control
// under the pointer, and Escape from inside returns the focus to the anchor.
func TestPopoverOutsidePressReachesItsTarget(t *testing.T) {
	open, pressed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(200).AlignItems(ui.Start).Children(func() {
			anchor := ui.Button(c, "Details").OnClick(func() { open = true })
			Popover(c, anchor, &open, PopoverOptions{Title: "Info", Body: func(c *ui.Context) { ui.Button(c, "Inside") }})
			ui.Button(c, "Elsewhere").OnClick(func() { pressed++ })
		})
	}, 640, 480)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Elsewhere"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open || pressed != 1 || !tt.Focused("Elsewhere") {
		t.Fatalf("outside press: open=%v pressed=%d, focus on target=%v", open, pressed, tt.Focused("Elsewhere"))
	}
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || !tt.Focused("Details") {
		t.Fatalf("Escape inside: open=%v, focus on anchor=%v", open, tt.Focused("Details"))
	}
}

// A non-modal surface takes no focus as it opens, so closing one with the
// focus inside must hand it on, or the keyboard user is left on nothing: a
// teaching tip gives it to its anchor, a floating panel or sheet to
// ReturnFocus.
func TestNonModalSurfacesReturnTheFocusAsTheyClose(t *testing.T) {
	tip, panel, sheet := true, true, true
	var opener ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Gap(120).AlignItems(ui.Start).Children(func() {
			ui.Button(c, "Opener").Bind(&opener)
			anchor := ui.Button(c, "Filters")
			TeachingTip(c, anchor, &tip, TeachingTipOptions{Title: "Saved filters", Body: "Reuse a filter."})
			FloatingPanel(c, &panel, FloatingPanelOptions{Title: "Inspector", Width: 200, ReturnFocus: &opener,
				Body: func() { ui.Button(c, "In panel") }})
			FloatingSheet(c, &sheet, FloatingSheetOptions{Title: "Details", ReturnFocus: &opener,
				Body: func() { ui.Button(c, "In sheet") }})
		})
	}, 640, 600)
	settle := func() {
		for range 30 {
			tt.Frame()
		}
	}
	settle()
	for _, close := range []ui.Key{ui.KeyEscape, ui.KeyEnter} {
		tip = true
		settle()
		if err := tt.Click("Filters"); err != nil {
			t.Fatal(err)
		}
		tt.Key(0, ui.KeyTab)
		if !tt.Focused("Close") {
			t.Fatal("Tab from the anchor did not enter the tip")
		}
		tt.Key(0, close)
		settle()
		if tip || !tt.Focused("Filters") {
			t.Fatalf("tip, key %v: open=%v, anchor focused=%v", close, tip, tt.Focused("Filters"))
		}
	}
	for _, s := range []struct {
		inside string
		open   *bool
	}{{"In panel", &panel}, {"In sheet", &sheet}} {
		if err := tt.Click(s.inside); err != nil {
			t.Fatal(err)
		}
		tt.Key(0, ui.KeyEscape)
		settle()
		if *s.open || !tt.Focused("Opener") {
			t.Fatalf("%s: open=%v, ReturnFocus focused=%v", s.inside, *s.open, tt.Focused("Opener"))
		}
	}
}

// Escape on an alert dialog is its cancel choice, as on any modal dialog:
// it closes the alert, runs OnCancel once and gives the focus back.
func TestAlertDialogEscapeCancels(t *testing.T) {
	open, accepted, cancelled := false, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Button(c, "Delete file").OnClick(func() { open = true })
		AlertDialog(c, &open, AlertDialogOptions{
			Title: "Delete report?", CancelLabel: "Cancel", ActionLabel: "Delete",
			OnAction: func() { accepted++ }, OnCancel: func() { cancelled++ },
		})
	}, 640, 480)
	if err := tt.Click("Delete file"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	for range 30 {
		tt.Frame()
	}
	if open || cancelled != 1 || accepted != 0 || !tt.Focused("Delete file") {
		t.Fatalf("Escape: open=%v cancelled=%d accepted=%d, opener focused %v", open, cancelled, accepted, tt.Focused("Delete file"))
	}
}
