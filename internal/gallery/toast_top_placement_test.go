package gallery

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// toastPageClick scrolls the page's content until label is inside the window
// and clicks it. The Toast page stacks three example frames, so its controls
// sit below the fold at the default scroll offset.
func toastPageClick(t *testing.T, tt *ui.Tester, state *State, label string) {
	t.Helper()
	for _, off := range []float32{0, 300, 600, 900, 1200, 1500, 1800, 2100, 2400, 2700} {
		state.Content.Y = off
		tt.Frame()
		r, ok := tt.Find(label)
		if !ok {
			continue
		}
		x, y := r.X+r.W/2, r.Y+r.H/2
		if x < 0 || y < 0 || x >= 1180 || y >= 900 {
			continue
		}
		tt.ClickAt(x, y)
		tt.Frame()
		return
	}
	t.Fatalf("%s: no clickable element", label)
}

// TestToastTopPlacementDismissDuringExit is the regression for the reported
// gallery crash: on the Toast page, anchor the region at the top, fire a
// toast, let its entry transition settle, then dismiss it with the close
// button or the action. The exit transition paints the retained copy in a
// later build pass, where the Elements captured by the close button's and
// the action's Draw callbacks had expired: "ui: element from pass N used in
// pass N+2", which quit the app.
func TestToastTopPlacementDismissDuringExit(t *testing.T) {
	state := New()
	state.Selected = "toast"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()

	// Move the region to the top edge; this fires a placement toast too.
	toastPageClick(t, tt, state, "Anchor toasts at top")
	if !overlaysStateOf(state).toastTop {
		t.Fatal("toast region did not anchor at the top")
	}

	toastPageClick(t, tt, state, "Default toast")
	tt.Advance(300 * time.Millisecond)
	tt.Frame()

	if _, ok := tt.Find("Close"); !ok {
		t.Fatal("fired toast shows no close button")
	}
	toastPageClick(t, tt, state, "Close")
	// The exit transition paints the retained copy in later passes.
	tt.Advance(60 * time.Millisecond)
	tt.Frame()

	// The action path paints the same retained copies.
	toastPageClick(t, tt, state, "Action toast")
	tt.Advance(300 * time.Millisecond)
	tt.Frame()
	toastPageClick(t, tt, state, "Undo")
	tt.Advance(60 * time.Millisecond)
	tt.Frame()

	// Close every toast the action left, each exit painting a retained copy.
	for range 4 {
		if _, ok := tt.Find("Close"); !ok {
			break
		}
		toastPageClick(t, tt, state, "Close")
		tt.Advance(300 * time.Millisecond)
		tt.Frame()
	}
	if got := len(overlaysStateOf(state).toasts.Items()); got != 0 {
		t.Fatalf("%d toasts still retained", got)
	}
}
