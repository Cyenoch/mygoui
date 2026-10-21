package overlay

import (
	"fmt"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
	"slices"
	"testing"
	"time"
)

func TestToastExpiresPausesAndAnnouncesOnce(t *testing.T) {
	var state ToastState
	closed := 0
	added := false
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		if !added {
			added = true
			Toast(c, &state, ToastOptions{Title: "Saved", Timeout: time.Second, OnClose: func() { closed++ }})
		}
		Toaster(c)
	}, 640, 480)
	if announcements := tt.Announcements(); !slices.Equal(announcements, []string{"Saved"}) {
		t.Fatalf("initial toast announcement=%q", announcements)
	}
	tt.Frame()
	tt.Frame()
	if len(tt.Announcements()) != 0 {
		t.Fatal("ordinary rebuild repeated toast announcement")
	}
	r, ok := tt.Find("Saved")
	if !ok {
		t.Fatal("native toast did not show")
	}
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(2 * time.Second)
	if !tt.HasText("Saved") || closed != 0 {
		t.Fatal("hover did not pause native expiration")
	}
	tt.Move(0, 0)
	tt.Advance(time.Second)
	if tt.HasText("Saved") || closed != 1 || len(state.ids) != 0 {
		t.Fatal("expired toast or callback retained")
	}
}

func TestToastFocusPauseActionDismissalAndCapacity(t *testing.T) {
	var state ToastState
	added := false
	actions, closed := 0, 0
	tt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		if !added {
			added = true
			Toast(c, &state, ToastOptions{Title: "Removed", Timeout: time.Second, Action: "Undo", OnAction: func() { actions++ }, OnClose: func() { closed++ }})
		}
		Toaster(c)
	}, 640, 480)
	tt.Move(0, 0)
	tt.Key(0, myui.KeyTab)
	tt.Advance(2 * time.Second)
	if !tt.HasText("Removed") || !tt.Focused("Undo") {
		t.Fatal("focused action did not pause expiration")
	}
	tt.Key(0, myui.KeyEnter)
	if actions != 1 || closed != 1 || len(state.ids) != 0 {
		t.Fatalf("action=%d close=%d retained=%d", actions, closed, len(state.ids))
	}
	added = false
	tt.Frame()
	tt.Key(0, myui.KeyTab)
	tt.Key(0, myui.KeyEscape)
	if tt.HasText("Removed") {
		t.Fatal("Escape did not dismiss focused toast")
	}

	var bounded ToastState
	inserted := false
	evicted := 0
	bt := myui.NewTester(func(c *myui.Context) {
		theme.Apply(c)
		if !inserted {
			inserted = true
			for i := range ToastCapacity + 5 {
				Toast(c, &bounded, ToastOptions{Title: fmt.Sprintf("Message %d", i), Timeout: -1, OnClose: func() { evicted++ }})
			}
		}
		Toaster(c)
	}, 640, 480)
	if len(bounded.ids) != ToastCapacity || evicted != 5 {
		t.Fatalf("collection retained=%d evicted=%d", len(bounded.ids), evicted)
	}
	if bt.HasText("Message 0") || !bt.HasText("Message 36") {
		t.Fatal("viewport must show only newest three bounded toasts")
	}
}
