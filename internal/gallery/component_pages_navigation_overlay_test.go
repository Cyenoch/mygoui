package gallery

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// renderComponentPage opens id in a fresh desktop-sized workbench.
func renderComponentPage(t *testing.T, id string) *ui.Tester {
	t.Helper()
	state := New()
	state.Selected = id
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
	tt.Frame()
	return tt
}

func pageByID(id string) (Page, bool) {
	for _, page := range Pages() {
		if page.ID == id {
			return page, true
		}
	}
	return Page{}, false
}

func TestNavigationComponentPagesExistAndRender(t *testing.T) {
	for _, want := range []struct {
		id, group, text string
	}{
		{"accordion", "Navigation", "Members"},
		{"disclosure", "Navigation", "One caller-owned panel plus a disabled trigger."},
		{"disclosuregroup", "Navigation", "Single-open and multi-open coordination from caller-owned flags."},
		{"tabs", "Navigation", "Account overview and recent activity."},
		{"breadcrumbs", "Navigation", "click a link to report its index"},
		{"pagination", "Navigation", "page 4 of 12"},
		{"toolbar", "Navigation", "Text formatting"},
		{"menu", "Navigation", "Reports"},
		{"dropdown", "Navigation", "Click the trigger, choose a row, or press Escape to close the panel."},
		{"submenu", "Navigation", "Escape or Left closes the panel without choosing."},
	} {
		page, ok := pageByID(want.id)
		if !ok {
			t.Fatalf("missing page %q", want.id)
		}
		if page.Group != want.group {
			t.Fatalf("%s group = %q, want %q", want.id, page.Group, want.group)
		}
		if page.Order < 10 {
			t.Fatalf("%s order = %d, want 10 or later", want.id, page.Order)
		}
		tt := renderComponentPage(t, want.id)
		if !tt.HasText(page.Title) {
			t.Fatalf("%s: page title missing", want.id)
		}
		if !tt.HasText(want.text) {
			t.Fatalf("%s: demo does not show %q", want.id, want.text)
		}
	}
}

func TestOverlayComponentPagesExistAndRender(t *testing.T) {
	for _, want := range []struct {
		id, group, text string
	}{
		{"modal", "Overlays", "Open Medium modal"},
		{"alertdialog", "Overlays", "The cancel action receives focus first"},
		{"drawer", "Overlays", "Open Bottom drawer"},
		{"bottomsheet", "Overlays", "Open bottom sheet"},
		{"dialog", "Overlays", "Open dialog"},
		{"popover", "Overlays", "The panel is non-modal"},
		{"tooltip", "Overlays", "Tab to the button to show the tip"},
		{"toast", "Overlays", "Default toast"},
	} {
		page, ok := pageByID(want.id)
		if !ok {
			t.Fatalf("missing page %q", want.id)
		}
		if page.Group != want.group {
			t.Fatalf("%s group = %q, want %q", want.id, page.Group, want.group)
		}
		if page.Order < 10 {
			t.Fatalf("%s order = %d, want 10 or later", want.id, page.Order)
		}
		tt := renderComponentPage(t, want.id)
		if !tt.HasText(page.Title) {
			t.Fatalf("%s: page title missing", want.id)
		}
		if !tt.HasText(want.text) {
			t.Fatalf("%s: demo does not show %q", want.id, want.text)
		}
	}
}

func TestAccordionPageShowsUsageAPIAndTokens(t *testing.T) {
	tt := renderComponentPage(t, "accordion")
	for _, want := range []string{
		"Example", "Usage", "API", "Theme tokens",
		"AccordionOptions{", "HideSeparator", "SurfaceForeground",
	} {
		if !tt.HasText(want) {
			t.Fatalf("accordion page does not show %q", want)
		}
	}
}

// toastPageScrollClick scrolls the page column until label is inside the
// window, then clicks it. The Toast page stacks its controls in three
// Examples frames, so they sit below the fold at the default offset.
func toastPageScrollClick(t *testing.T, tt *ui.Tester, state *State, label string) {
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
	toastPageScrollClick(t, tt, state, "Anchor toasts at top")
	if !overlaysStateOf(state).toastTop {
		t.Fatal("toast region did not anchor at the top")
	}
	toastPageScrollClick(t, tt, state, "Action toast")
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
	toastPageScrollClick(t, tt, state, "Dismiss all toasts")
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
