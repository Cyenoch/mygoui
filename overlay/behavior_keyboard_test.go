package overlay

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// The action bar is one Tab stop; the arrows move between its controls.
func TestActionBarIsOneTabStopWithArrowsInside(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Row"})
			ActionBar(c, true, ActionBarOptions{
				Label: "selected", Count: 2,
				Actions: []navigation.OverflowItem{{Label: "Archive"}, {Label: "Delete"}},
				OnClose: func() {},
			})
		})
	}, 640, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	tt.Frame()
	if err := tt.Click("Row"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Archive") {
		t.Fatal("Tab did not enter the bar at its first action")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Close") {
		t.Fatal("Right did not move through the bar to Close")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Row") {
		t.Fatal("Tab did not leave the bar in one step")
	}
}

// The resize handle answers Up as well as Down, and Escape with the focus
// inside the sheet closes it.
func TestFloatingSheetGrowsByKeyAndClosesOnEscapeInside(t *testing.T) {
	open, height := true, float32(150)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ui.Button(c, "Beneath")
			FloatingSheet(c, &open, FloatingSheetOptions{
				Title: "Details", Height: &height, MinHeight: 120, MaxHeight: 240,
				Body: func() { ui.Text(c, "Body") },
			})
		})
	}, 480, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if err := tt.Click("Beneath"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEscape)
	if !open {
		t.Fatal("Escape outside the sheet closed it")
	}
	for !tt.Focused("Resize") {
		tt.Key(0, ui.KeyTab)
		if tt.Focused("Beneath") {
			t.Fatal("Tab never reached the resize handle")
		}
	}
	for range 8 {
		tt.Key(0, ui.KeyUp)
	}
	if height != 240 {
		t.Fatalf("height after Up %v, want the maximum 240", height)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || tt.HasText("Details") {
		t.Fatalf("Escape inside: open=%v, shown=%v", open, tt.HasText("Details"))
	}
}

// Alt with an arrow moves a floating panel while the focus is inside it, so
// placing it needs no pointer; the window still bounds it.
func TestFloatingPanelMovesByKeyWithinTheWindow(t *testing.T) {
	open := true
	pos := ui.Point{X: 100, Y: 40}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		FloatingPanel(c, &open, FloatingPanelOptions{
			Title: "Inspector", Position: &pos, Width: 200,
			Body: func() { ui.Button(c, "Inside") },
		})
	}, 480, 320)
	tt.Frame()
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Alt, ui.KeyRight)
	tt.Key(ui.Alt, ui.KeyDown)
	if pos != (ui.Point{X: 124, Y: 64}) {
		t.Fatalf("after Alt+Right and Alt+Down the panel is at %+v, want {124 64}", pos)
	}
	for range 10 {
		tt.Key(ui.Alt, ui.KeyUp)
	}
	if pos.Y != 0 || !open {
		t.Fatalf("Alt+Up past the top: %+v, open %v", pos, open)
	}
}

// The resize handle tells assistive technology the sheet's height within
// its limits, and Home and End reach those limits.
func TestFloatingSheetSplitterExposesItsRangeAndReachesItsEnds(t *testing.T) {
	open, height := true, float32(150)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			FloatingSheet(c, &open, FloatingSheetOptions{
				Title: "Details", Height: &height, MinHeight: 120, MaxHeight: 240,
			})
		})
	}, 480, 400)
	tt.Frame()
	tt.Frame()
	if grip, ok := tt.AccessByRole(ui.RoleSplitter, "Resize"); !ok || grip.Min != 120 || grip.Max != 240 || grip.Now != 150 {
		t.Fatalf("splitter range: %+v (found %v)", grip, ok)
	}
	if err := tt.Click("Resize"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnd)
	if height != 240 {
		t.Fatalf("End: height %v, want 240", height)
	}
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	if grip, _ := tt.AccessByRole(ui.RoleSplitter, "Resize"); height != 120 || grip.Now != 120 {
		t.Fatalf("Home: height %v, splitter value %v, want 120", height, grip.Now)
	}
}

// Under reduced motion the action bar appears and leaves at once, as every
// overlay does (DESIGN.md, Still Overlay Rule).
func TestActionBarAppearsAndLeavesAtOnceUnderReducedMotion(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ActionBar(c, open, ActionBarOptions{Label: "selected", Count: 2, OnClose: func() { open = false }})
		})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	open = true
	tt.Frame()
	if !tt.HasText("selected") || tt.FramePending() {
		t.Fatalf("opening: shown %v, still animating %v", tt.HasText("selected"), tt.FramePending())
	}
	open = false
	tt.Frame()
	if tt.HasText("selected") || tt.FramePending() {
		t.Fatalf("closing: shown %v, still animating %v", tt.HasText("selected"), tt.FramePending())
	}
}

// A popover is a non-modal dialog: the focus moves into it as it opens,
// Tab reaches its controls, and Tab past its last control closes it rather
// than leaving it open behind the focus.
func TestPopoverTakesTheFocusAndClosesAsTheFocusMovesOn(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
			anchor := ui.Button(c, "Open popover").OnClick(func() { open = true })
			Popover(c, anchor, &open, PopoverOptions{Title: "Details",
				Body: func(c *ui.Context) { ui.Button(c, "Inside") }})
			ui.Button(c, "After")
		})
	}, 640, 480)
	if err := tt.Click("Open popover"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !open || !tt.Focused("Details") {
		t.Fatal("opening did not move the focus into the panel")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Inside") {
		t.Fatal("Tab from the panel did not reach its control")
	}
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if open || !tt.Focused("After") {
		t.Fatalf("Tab past the popover: open=%v, After focused=%v", open, tt.Focused("After"))
	}
}

// F6 brings the keyboard to the notifications, and closing the toast that
// holds the focus hands it to the newest toast left, so a keyboard user can
// work through the stack without the focus falling to nothing.
func TestToastRegionTakesTheFocusAndKeepsItAsToastsClose(t *testing.T) {
	f := newToastMotionFixture(t, 640, 600)
	for _, id := range []string{"older", "newer"} {
		f.add(ToastOptions{ID: id, Title: id, Timeout: -1})
	}
	f.tt.Advance(time.Second)
	if err := f.tt.Click("Page action"); err != nil {
		t.Fatal(err)
	}
	f.tt.Key(0, ui.KeyF6)
	f.tt.Frame()
	f.tt.Key(0, ui.KeyEnter)
	f.tt.Advance(time.Second)
	if items := f.state.Items(); len(items) != 1 || items[0].ID != "older" {
		t.Fatalf("F6 then Enter did not close the newest toast: %+v", items)
	}
	f.tt.Key(0, ui.KeyEnter)
	f.tt.Advance(time.Second)
	if items := f.state.Items(); len(items) != 0 {
		t.Fatalf("the focus did not move on to the toast left: %+v", items)
	}
}

// A popover that opens on what the user types into its anchor, as a
// suggestion list does, leaves the focus in the field so typing goes on, and
// the focus in the anchor does not close it.
func TestPopoverOpenedWhileTypingInItsAnchorLeavesTheFocus(t *testing.T) {
	open := false
	text := ""
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(8).AlignItems(ui.Start).Children(func() {
			anchor := ui.TextInput(c, &text).Label("Query").Width(200)
			if text != "" {
				open = true
			}
			Popover(c, anchor, &open, PopoverOptions{Title: "Suggestions"})
		})
	}, 640, 480)
	if err := tt.Click("Query"); err != nil {
		t.Fatal(err)
	}
	tt.Type("a")
	tt.Frame()
	tt.Type("b")
	tt.Frame()
	if text != "ab" || !tt.Focused("Query") {
		t.Fatalf("typing through the opening popover gave %q, field focused %v", text, tt.Focused("Query"))
	}
	if !open {
		t.Fatal("the focus in the anchor closed the popover")
	}
}

// Switching to another window and back leaves a popover open, as it was in
// v0.2.0: the window losing the keyboard is not the user moving on.
func TestPopoverAndPopconfirmSurviveWindowDeactivation(t *testing.T) {
	deactivate := func(tt *ui.Tester) {
		tt.SetWindowActive(false)
		tt.SetFocused(false)
		tt.Frame()
		tt.SetFocused(true)
		tt.SetWindowActive(true)
		tt.Frame()
	}
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
			anchor := ui.Button(c, "Open popover").OnClick(func() { open = true })
			Popover(c, anchor, &open, PopoverOptions{Title: "Details"})
		})
	}, 640, 480)
	if err := tt.Click("Open popover"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	deactivate(tt)
	if !open {
		t.Fatal("deactivating the window closed the popover")
	}

	var confirmOpen bool
	var confirmed, cancelled int
	pc := popconfirmTester(&confirmOpen, &confirmed, &cancelled)
	if err := pc.Click("Delete"); err != nil {
		t.Fatal(err)
	}
	pc.Frame()
	deactivate(pc)
	if !confirmOpen || cancelled != 0 {
		t.Fatalf("deactivating the window: open=%v cancelled=%d", confirmOpen, cancelled)
	}
}

// Under reduced motion the bar stays one live element as the preference
// turns on while it shows, and a theme change recolours it at once.
func TestActionBarStaysWholeAndStillAsReducedMotionTurnsOn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Children(func() {
			ActionBar(c, true, ActionBarOptions{Label: "selected", Count: 2,
				Actions: []navigation.OverflowItem{{Label: "Archive"}}})
		})
	}, 480, 320)
	shown := func() int {
		n := 0
		for _, s := range tt.Texts() {
			if s == "Archive" {
				n++
			}
		}
		return n
	}
	tt.Advance(time.Second)
	before := shown()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	if after := shown(); after != before || tt.FramePending() {
		t.Fatalf("reduced motion on: %d texts of the bar, was %d; still animating %v", after, before, tt.FramePending())
	}
	tt.SetDark(true)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("a theme change under reduced motion fades the bar's colours")
	}
}
