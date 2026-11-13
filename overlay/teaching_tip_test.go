package overlay

import (
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestTeachingTipLeavesTheWindowUsableUntilDismissed(t *testing.T) {
	open, pressed, dismissed := true, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(200).AlignItems(ui.Start).Children(func() {
			anchor := mygoui.Button(c, mygoui.ButtonOptions{Label: "Filters"})
			TeachingTip(c, anchor, &open, TeachingTipOptions{
				Title: "Saved filters", Body: "Reuse a filter from here.",
				OnDismiss: func() { dismissed++ },
			})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Beneath", OnClick: func() { pressed++ }})
		})
	}, 480, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()

	tip, ok := tt.AccessByRole(ui.RoleDialog, "Saved filters")
	if !ok || tip.Description != "Reuse a filter from here." {
		t.Fatalf("tip is not exposed as a described dialog: %+v (found %v)", tip, ok)
	}
	if err := tt.Click("Beneath"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if pressed != 1 || !open || !tt.HasText("Saved filters") {
		t.Fatalf("outside press: pressed=%d open=%v", pressed, open)
	}
	// The tip takes no focus, so Escape outside it is not its key; it
	// follows its anchor in the Tab order.
	tt.Key(0, ui.KeyEscape)
	if !open {
		t.Fatal("Escape outside the tip closed it")
	}
	if err := tt.Click("Filters"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || dismissed != 1 || tt.HasText("Saved filters") {
		t.Fatalf("Escape inside the tip: open=%v dismissed=%d", open, dismissed)
	}

	open = true
	tt.Frame()
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if open || dismissed != 2 {
		t.Fatalf("close button: open=%v dismissed=%d", open, dismissed)
	}
}

func tourTester(state *TourState, reduce bool, o TourOptions, steps func(a, b *ui.Handle) []TourStep, pressed *int) *ui.Tester {
	var first, second ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Padding(16).Gap(200).AlignItems(ui.Start).Children(func() {
			ui.Row(c).Gap(240).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Inbox"}).Bind(&first)
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Archive"}).Bind(&second)
			})
			mygoui.Button(c, mygoui.ButtonOptions{Label: "Start", OnClick: func() {
				*state = TourState{Open: true}
				if pressed != nil {
					*pressed++
				}
			}})
			Tour(c, state, steps(&first, &second), o)
		})
	}, 640, 400)
	tt.SetPreferences(ui.Preferences{ReduceMotion: reduce, TextScale: 1})
	tt.Frame()
	return tt
}

func twoSteps(a, b *ui.Handle) []TourStep {
	return []TourStep{
		{Anchor: a, Title: "Your inbox", Body: "New mail arrives here."},
		{Anchor: b, Title: "The archive", Body: "Mail you are done with."},
	}
}

func TestTourIsModalAndDrivenByTheKeyboard(t *testing.T) {
	var state TourState
	finished, dismissed, pressed := 0, 0, 0
	tt := tourTester(&state, true, TourOptions{
		OnFinish: func() { finished++ }, OnDismiss: func() { dismissed++ },
	}, twoSteps, &pressed)
	if err := tt.Click("Start"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Frame()
	if _, ok := tt.AccessByRole(ui.RoleDialog, "Your inbox"); !ok || !tt.Focused("Next") {
		t.Fatalf("start: dialog shown and Next focused: %v", tt.Texts())
	}
	if _, ok := tt.Find("Back"); ok {
		t.Fatal("the first step offers Back")
	}
	// The page waits: a press beside the panel reaches nothing.
	behind, _ := tt.Find("Start")
	tt.ClickAt(behind.X+behind.W/2, behind.Y+behind.H/2)
	tt.Frame()
	if pressed != 1 || state.Step != 0 || !state.Open {
		t.Fatalf("outside press: pressed=%d state=%+v", pressed, state)
	}
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if state.Step != 1 || !tt.HasText("Done") || !tt.Focused("Done") {
		t.Fatalf("Right: state=%+v, Done focused=%v", state, tt.Focused("Done"))
	}
	tt.Key(0, ui.KeyLeft)
	tt.Frame()
	if state.Step != 0 || !tt.Focused("Next") {
		t.Fatalf("Left: state=%+v, Next focused=%v", state, tt.Focused("Next"))
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(time.Second)
	if state.Open || dismissed != 1 || finished != 0 || !tt.Focused("Start") {
		t.Fatalf("Escape: state=%+v dismissed=%d, focus back on Start=%v", state, dismissed, tt.Focused("Start"))
	}

	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.Open || finished != 1 || pressed != 2 {
		t.Fatalf("Done: state=%+v finished=%d pressed=%d", state, finished, pressed)
	}
}

// A step whose target is not on screen must not strand the tour open with
// nothing to press.
func TestTourStepWithoutTargetIsCentred(t *testing.T) {
	state := TourState{Open: true}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Tour(c, &state, []TourStep{{Title: "Hidden panel"}, {Title: "Next one"}}, TourOptions{})
	}, 640, 360)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Frame()
	tt.Frame()
	tip, ok := tt.AccessByRole(ui.RoleDialog, "Hidden panel")
	b := tip.Bounds
	if !ok || math.Abs(float64(b.X+b.W/2-320)) > 1 || math.Abs(float64(b.Y+b.H/2-180)) > 1 {
		t.Fatalf("unanchored step: %+v (found %v)", b, ok)
	}
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	if state.Step != 1 {
		t.Fatalf("step=%d", state.Step)
	}
}

// A step change is one motion: the panel shows in every frame and moves
// steadily toward the next target, then rests without asking for frames.
func TestTourTravelsContinuouslyToTheNextStep(t *testing.T) {
	state := TourState{Open: true}
	tt := tourTester(&state, false, TourOptions{}, twoSteps, nil)
	tt.Advance(time.Second)
	start, _ := tt.AccessByRole(ui.RoleDialog, "Your inbox")
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	last := start.Bounds.X
	for range 24 {
		tt.Advance(16 * time.Millisecond)
		panel, ok := tt.AccessByRole(ui.RoleDialog, "The archive")
		if !ok || panel.Bounds.W <= 0 {
			t.Fatalf("panel missing mid-travel at %v", tt.Now())
		}
		if panel.Bounds.X < last {
			t.Fatalf("panel went back: %v after %v", panel.Bounds.X, last)
		}
		last = panel.Bounds.X
	}
	end, _ := tt.AccessByRole(ui.RoleDialog, "The archive")
	if end.Bounds.X <= start.Bounds.X+150 || tt.FramePending() {
		t.Fatalf("panel x %v -> %v, frames pending at rest %v", start.Bounds.X, end.Bounds.X, tt.FramePending())
	}
}

func TestTourStepsInPlaceUnderReducedMotion(t *testing.T) {
	state := TourState{Open: true}
	tt := tourTester(&state, true, TourOptions{}, twoSteps, nil)
	tt.Frame()
	if err := tt.Click("Next"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	now, _ := tt.AccessByRole(ui.RoleDialog, "The archive")
	tt.Advance(time.Second)
	end, _ := tt.AccessByRole(ui.RoleDialog, "The archive")
	if now.Bounds.X != end.Bounds.X {
		t.Fatalf("reduced motion moved the panel: %v then %v", now.Bounds.X, end.Bounds.X)
	}
}
