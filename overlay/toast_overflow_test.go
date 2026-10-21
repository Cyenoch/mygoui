package overlay

import (
	"fmt"
	"slices"
	"testing"
	"time"

	myui "github.com/egoist/mygo/ui"
)

func TestToastOverflowHiddenRecordExpiresWithoutLosingNewest(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	closed := 0
	f.add(ToastOptions{ID: "overflow-expiring", Title: "Hidden finite notification", Timeout: 3 * time.Second, OnClose: func() { closed++ }})
	for i := range 3 {
		f.add(ToastOptions{ID: fmt.Sprintf("overflow-sticky-%d", i), Title: fmt.Sprintf("Sticky visible %d", i), Timeout: -1})
	}
	f.tt.Move(0, 0)
	f.tt.Advance(4 * time.Second)
	if closed != 1 || slices.Contains(f.state.ids, "overflow-expiring") || len(f.state.ids) != 3 {
		t.Fatalf("hidden expiry lifecycle: close=%d ids=%q", closed, f.state.ids)
	}
	for i := range 3 {
		if !f.tt.HasText(fmt.Sprintf("Sticky visible %d", i)) {
			t.Fatalf("hidden expiry displaced live visible record %d: %q", i, f.tt.Texts())
		}
	}
	f.tt.Advance(4 * time.Second)
	if closed != 1 {
		t.Fatalf("hidden expiry callback repeated: %d", closed)
	}
}

func TestToastOverflowHiddenSameIDUpdateOwnsCurrentPayload(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	oldClosed, currentClosed, actions := 0, 0, 0
	f.add(ToastOptions{ID: "overflow-hidden", Title: "Old hidden loading", Kind: ToastLoading, OnClose: func() { oldClosed++ }})
	for i := range 3 {
		f.add(ToastOptions{ID: fmt.Sprintf("overflow-front-%d", i), Title: fmt.Sprintf("Front card %d", i), Timeout: -1})
	}
	f.tt.Advance(2 * time.Second)
	f.add(ToastOptions{ID: "overflow-hidden", Title: "Updated hidden result", Description: "The current description must be readable when this record is reached.", Kind: ToastSuccess,
		Action: "Inspect hidden result", Timeout: -1, OnAction: func() { actions++ }, OnClose: func() { currentClosed++ }})
	if len(f.state.ids) != 4 || oldClosed != 0 || currentClosed != 0 {
		t.Fatalf("hidden update changed total or closed record: ids=%q old=%d current=%d", f.state.ids, oldClosed, currentClosed)
	}
	// Closing a newer entry naturally exposes the authoritative native record;
	// no copied payload or wrapper history is used to recover its content.
	f.dismiss("overflow-front-2")
	f.tt.Advance(2 * time.Second)
	if f.tt.HasText("Old hidden loading") || !f.tt.HasText("Updated hidden result") || !f.tt.HasText("The current description must be readable when this record is reached.") {
		t.Fatalf("native hidden update lost current payload: %q", f.tt.Texts())
	}
	notificationFocus(t, f.tt, "Inspect hidden result")
	f.tt.Advance(2 * time.Second)
	f.tt.Key(0, myui.KeyEnter)
	f.tt.Advance(2 * time.Second)
	if actions != 1 || currentClosed != 1 || oldClosed != 0 || slices.Contains(f.state.ids, "overflow-hidden") {
		t.Fatalf("hidden same-ID callback ownership: action=%d old=%d current=%d ids=%q", actions, oldClosed, currentClosed, f.state.ids)
	}
}

func TestToastOverflowThirtyTwoRecordsRemainNativeAndBounded(t *testing.T) {
	f := newNotificationFixture(t, 640, 850)
	closed := make([]int, ToastCapacity)
	f.pending = append(f.pending, func(c *myui.Context) {
		for i := range ToastCapacity {
			Toast(c, &f.state, ToastOptions{ID: fmt.Sprintf("overflow-capacity-%d", i), Title: fmt.Sprintf("Retained notification %02d", i), Timeout: -1, OnClose: func() { closed[i]++ }})
		}
	})
	f.tt.Frame()
	f.tt.Advance(2 * time.Second)
	if len(f.state.ids) != ToastCapacity {
		t.Fatalf("retained count=%d want=%d", len(f.state.ids), ToastCapacity)
	}
	for i, count := range closed {
		if count != 0 {
			t.Fatalf("retained message %d was closed at capacity: %d", i, count)
		}
	}
	// Native removal must release one callback per retained entry even when
	// most of those entries were never in the rendered newest-three window.
	f.dismiss("")
	f.tt.Advance(2 * time.Second)
	for i, count := range closed {
		if count != 1 {
			t.Fatalf("clear-all retained message %d callback=%d want=1", i, count)
		}
	}
	if len(f.state.ids) != 0 || f.tt.HasText("Retained notification") || f.tt.FramePending() {
		t.Fatalf("overflow clear retained native data/animation: ids=%q texts=%q", f.state.ids, f.tt.Texts())
	}
}
