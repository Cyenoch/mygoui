package gallery

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/overlay"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDefaultPaginationIsCenteredAndIndependent(t *testing.T) {
	page, ok := componentPage("pagination")
	if !ok || len(page.Examples) < 5 || page.Examples[0].Title != "Default" {
		t.Fatal("Pagination does not start with the default example")
	}
	for _, width := range []int{360, 1180} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			s := New()
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				ui.Column(c).Fill().Children(func() { page.Examples[0].Demo(c, s) })
			}, width, 300)
			tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: true})
			st := navigationStateOf(s)
			if st.pageDefault != 2 || st.pageSummary != 2 {
				t.Fatalf("initial pagination pages: default=%d summary=%d", st.pageDefault, st.pageSummary)
			}
			previous, prevOK := tt.Find("Previous")
			next, nextOK := tt.Find("Next")
			if !prevOK || !nextOK ||
				math.Abs(float64((previous.X+next.X+next.W)/2-float32(width)/2)) > 1 {
				t.Fatalf("pagination is not centered: previous=%+v next=%+v", previous, next)
			}
			for _, label := range []string{"1", "2, Selected", "3", "More pages", "5"} {
				if _, ok := tt.Find(label); !ok {
					t.Fatalf("default pagination is missing %q", label)
				}
			}
			if err := tt.Click("Next"); err != nil {
				t.Fatal(err)
			}
			if st.pageDefault != 3 || st.pageSummary != 2 {
				t.Fatalf("Next changed default=%d summary=%d", st.pageDefault, st.pageSummary)
			}
			tt.Key(0, ui.KeyEnter)
			if st.pageDefault != 4 {
				t.Fatalf("Enter on focused Next selected page %d", st.pageDefault)
			}
			if err := tt.Click("1"); err != nil {
				t.Fatal(err)
			}
			if err := tt.Click("Previous"); err != nil {
				t.Fatal(err)
			}
			if st.pageDefault != 1 {
				t.Fatalf("disabled Previous selected page %d", st.pageDefault)
			}
		})
	}
}

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

func TestDrawerDemoCloseRestoresItsTrigger(t *testing.T) {
	state := New()
	state.Selected = "drawer"
	tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 900, 900)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	revealContent(t, tt, state, "Open Right drawer", 900, 900)
	if err := tt.Click("Open Right drawer"); err != nil {
		t.Fatal(err)
	}
	if !overlaysStateOf(state).drawerOpen {
		t.Fatal("drawer did not open")
	}
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Frame()
	if overlaysStateOf(state).drawerOpen || overlaysStateOf(state).drawerNote != "Right drawer dismissed" {
		t.Fatal("close did not dismiss the drawer")
	}
	if !tt.Focused("Open Right drawer") {
		t.Fatal("closing the drawer did not restore focus to its trigger")
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

func TestToastKeepsItsRendererAcrossPageNavigation(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(fmt.Sprint(dark), func(t *testing.T) {
			state := New()
			state.Selected, state.Dark = "toast", dark
			tt := ui.NewTester(func(c *ui.Context) { View(c, state) }, 1180, 900)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			toastPageClick(t, tt, state, "Action toast")
			before, ok := tt.Find("Invite sent")
			if !ok {
				t.Fatal("toast missing")
			}
			face := tt.Image().RGBAAt(int(before.X)+8, int(before.Y+before.H/2))
			page, _ := pageFor("previewcard")
			selectPage(state, page)
			tt.Frame()
			after, ok := tt.Find("Invite sent")
			if !ok || after != before {
				t.Fatalf("toast renderer/layout changed on navigation: before=%+v after=%+v", before, after)
			}
			if got := tt.Image().RGBAAt(int(after.X)+8, int(after.Y+after.H/2)); got != face {
				t.Fatalf("toast face changed on navigation: %v -> %v", face, got)
			}
			if err := tt.Click("Undo"); err != nil {
				t.Fatal(err)
			}
			if overlaysStateOf(state).toastNote != "undo pressed; invitation revoked" || !tt.HasText("Invitation revoked") {
				t.Fatal("navigation lost the toast action or its resulting notification")
			}
		})
	}
}

func TestToastExamplesRepeatedArrivalAndExplicitUpdate(t *testing.T) {
	st := &overlaysState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		st.toasts.View(c)
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
				ui.Row(c.Key("leaving")).Transition(ui.ElementTransition{
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
