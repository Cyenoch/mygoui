package gallery

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestComplexPageDemosRenderLiveControls keeps the complex families honest:
// each page must render the real control its example drives, not an empty
// frame. Titles, groups and order are catalog data covered by
// TestCatalogIsUniqueAndComplete.
func TestComplexPageDemosRenderLiveControls(t *testing.T) {
	for _, tc := range []struct{ id, control string }{
		{"form", "Submit form"},
		{"input", "Username"},
		{"checkbox", "Parent folder"},
		{"select", "Current fruit:"},
		{"listbox", "Fruit: banana"},
		{"table", "Deployments"},
		{"modal", "Open Medium modal"},
		{"drawer", "Open Bottom drawer"},
		{"dialog", "Open dialog"},
		{"toast", "Default toast"},
	} {
		state := New()
		state.Selected = tc.id
		tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
		tt.Frame()
		if !tt.HasText(tc.control) {
			t.Fatalf("%s page does not render its %q control", tc.id, tc.control)
		}
	}
}

// TestToastPageDismissesWithoutExpiredPaint fires the page's toasts and
// dismisses them while the exit transition paints the retained copies. Those
// copies used to run the close button's and the action's Draw callbacks with
// Elements from the build pass that made them, which panicked ("element from
// pass N used in pass N+2") and quit the app.
func TestToastPageDismissesWithoutExpiredPaint(t *testing.T) {
	state := New()
	state.Selected = "toast"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	// Anchor the region at the top: the toast stack then stays clear of the
	// page controls, which sit below the fold in the Examples frames.
	toastPageClick(t, tt, state, "Anchor toasts at top")
	if !overlaysStateOf(state).toastTop {
		t.Fatal("toast region did not anchor at the top")
	}
	toastPageClick(t, tt, state, "Action toast")
	// Settle the entry transitions so the leaving copies exist.
	tt.Advance(300 * time.Millisecond)
	tt.Frame()
	closeRect, ok := tt.Find("Close")
	if !ok {
		t.Fatal("fired toast shows no close button")
	}
	if closeRect.X <= 0 || closeRect.Y <= 0 ||
		closeRect.X+closeRect.W > 1180 || closeRect.Y+closeRect.H > 900 {
		t.Fatalf("close button %+v is outside the page", closeRect)
	}
	toastPageClick(t, tt, state, "Dismiss all toasts")
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if got := len(overlaysStateOf(state).toasts.Items()); got != 0 {
		t.Fatalf("dismiss all left %d retained toasts", got)
	}
}

func TestToastExamplesRepeatedArrivalAndExplicitUpdate(t *testing.T) {
	st := &overlaysState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Children(func() {
			toastVariantsExample(c, st)
			toastActionsExample(c, st)
		})
	}, 900, 600)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	for _, label := range []string{"Default toast", "Default toast", "Action toast", "Action toast"} {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	items := st.toasts.Items()
	if len(items) != 4 || items[0].ID == items[1].ID || items[2].ID == items[3].ID {
		t.Fatalf("repeated triggers replaced earlier notifications: %+v", items)
	}
	if err := tt.Click("Update toast"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Update toast"); err != nil {
		t.Fatal(err)
	}
	items = st.toasts.Items()
	if len(items) != 5 || items[4].Title != "Uploading… 50%" {
		t.Fatalf("explicit update did not replace its own notification: %+v", items)
	}
	if err := tt.Click("Dismiss update"); err != nil {
		t.Fatal(err)
	}
	if len(st.toasts.Items()) != 4 || tt.HasText("Uploading… 50%") {
		t.Fatal("dismissing the update affected independent arrivals")
	}
	if err := tt.Click("Dismiss all toasts"); err != nil {
		t.Fatal(err)
	}
	if len(st.toasts.Items()) != 0 || tt.HasText("Draft saved") || tt.HasText("Invite sent") {
		t.Fatal("dismiss all left a repeated arrival alive")
	}
}

// The modal's retained exit paints the last built body after its native
// dialog closes, without keeping its controls live or reading expired elements.
func TestModalWithButtonClosesWithoutExpiredPaint(t *testing.T) {
	var open bool
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			mygoui.Button(c, mygoui.ButtonOptions{
				Label: "Open modal", OnClick: func() { open = true },
			})
			overlay.Modal(c, &open, overlay.ModalOptions{
				Title: "Settings",
				Body: func() {
					mygoui.Button(c, mygoui.ButtonOptions{
						Label: "Inside", OnClick: func() { clicks++ },
					})
				},
			})
		})
	}, 640, 480)
	if err := tt.Click("Open modal"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Inside"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 {
		t.Fatalf("inside button clicks = %d, want 1", clicks)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(300 * time.Millisecond)
	tt.Frame()
	if open || tt.HasText("Inside") {
		t.Fatalf("modal survived Escape: open=%v", open)
	}
}

// TestButtonPaintsThroughExitTransition is the regression for a component
// inside an ElementTransition exit: MyGo paints the retained copy in later
// build passes, where the Element captured by the component's state face has
// expired. It used to panic in style.Background and style.CornerFill.
func TestButtonPaintsThroughExitTransition(t *testing.T) {
	show := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Children(func() {
			ui.Button(c, "Remove").OnClick(func() { show = false })
			if show {
				ui.Row(c).Key("leaving").Transition(ui.ElementTransition{
					Duration: 250 * time.Millisecond,
					Exit:     &ui.Motion{Opacity: 0},
				}).Children(func() {
					mygoui.Button(c, mygoui.ButtonOptions{Label: "Leaving"})
					mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{Label: "Leaving toggle"})
					mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{},
						mygoui.ButtonOptions{Label: "Group A"},
						mygoui.ButtonOptions{Label: "Group B"},
					)
				})
			}
		})
	}, 640, 480)
	tt.Frame()
	if err := tt.Click("Remove"); err != nil {
		t.Fatal(err)
	}
	// The exit frame paints the retained copy of every state face.
	tt.Advance(60 * time.Millisecond)
	tt.Frame()
	if tt.HasText("Leaving") {
		t.Fatal("leaving row is still shown")
	}
}
