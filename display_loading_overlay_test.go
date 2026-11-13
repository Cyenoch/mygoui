package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestLoadingOverlayBlocksAndAnnouncesWhileVisible pins the cover: while the
// caller's flag is set the content takes neither a click nor the keyboard
// focus, the container is busy and a status names the wait; cleared, the
// content works again.
func TestLoadingOverlayBlocksAndAnnouncesWhileVisible(t *testing.T) {
	visible, saves := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.LoadingOverlay(c, mygoui.LoadingOverlayOptions{Visible: visible, Label: "Saving invoice"}, func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: func() { saves++ }})
		}).Label("Invoice")
	}, 320, 160)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Save") {
		t.Fatal("Tab does not reach the uncovered content")
	}
	if _, ok := tt.AccessByRole(ui.RoleStatus, "Saving invoice"); ok {
		t.Fatal("a hidden overlay announces a wait")
	}

	visible = true
	tt.Frame()
	if tt.Focused("Save") {
		t.Fatal("covered content kept the keyboard focus")
	}
	if err := tt.Click("Save"); err == nil || saves != 0 {
		t.Fatalf("covered content took a click (err=%v, saves=%d)", err, saves)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if tt.Focused("Save") || saves != 0 {
		t.Fatalf("covered content took the keyboard (saves=%d)", saves)
	}
	if a, _ := tt.AccessByRole(ui.RoleGroup, "Invoice"); !a.Busy {
		t.Fatal("the covered container is not busy")
	}
	if _, ok := tt.AccessByRole(ui.RoleStatus, "Saving invoice"); !ok {
		t.Fatal("the cover does not name the wait")
	}

	visible = false
	tt.Frame()
	if err := tt.Click("Save"); err != nil || saves != 1 {
		t.Fatalf("uncovered content did not take a click (err=%v, saves=%d)", err, saves)
	}
}

// TestLoadingOverlayLooksLateButHoldsInputAtOnce pins the shared timing: a
// load that ends within the wait paints no scrim, yet the content takes no
// input from the first frame.
func TestLoadingOverlayLooksLateButHoldsInputAtOnce(t *testing.T) {
	visible, saves := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.LoadingOverlay(c, mygoui.LoadingOverlayOptions{Visible: visible}, func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Save", OnClick: func() { saves++ }})
		}).Width(320).Height(160)
	}, 320, 160)
	save := findRect(t, tt, "Save")
	x, y := save.X+3, save.Y+save.H/2
	before := pixel(tt.Image(), x, y)

	visible = true
	tt.Frame()
	if err := tt.Click("Save"); err == nil || saves != 0 {
		t.Fatalf("content took a click in the first frame (err=%v, saves=%d)", err, saves)
	}
	if got := pixel(tt.Image(), x, y); !same(got, before) {
		t.Fatalf("the scrim showed before the wait: %v vs %v", got, before)
	}
	tt.Advance(400 * time.Millisecond)
	if got := pixel(tt.Image(), x, y); same(got, before) {
		t.Fatal("the scrim never appeared")
	}
}

// TestLoadingOverlayGivesTheFocusBackWhenItHides pins the return of the
// focus the cover took: the keyboard continues where the user left it.
func TestLoadingOverlayGivesTheFocusBackWhenItHides(t *testing.T) {
	visible := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.LoadingOverlay(c, mygoui.LoadingOverlayOptions{Visible: visible}, func() {
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Refresh"})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Save"})
		})
	}, 320, 160)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Save") {
		t.Fatal("Tab does not reach Save")
	}
	visible = true
	tt.Frame()
	if tt.Focused("Save") {
		t.Fatal("covered content kept the focus")
	}
	visible = false
	tt.Advance(400 * time.Millisecond)
	if !tt.Focused("Save") {
		t.Fatal("hiding the overlay did not give the focus back")
	}
}
