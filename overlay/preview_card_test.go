package overlay

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestPreviewCardHoverFocusAndDismiss(t *testing.T) {
	var state PreviewCardState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(24).Children(func() {
			anchor := ui.Button(c, "Profile")
			PreviewCard(c, anchor, &state, PreviewCardOptions{
				Body: func() { ui.Button(c, "Follow") },
			})
		})
	}, 480, 320)
	anchor, ok := tt.Find("Profile")
	if !ok {
		t.Fatal("missing preview anchor")
	}
	tt.Move(anchor.X+anchor.W/2, anchor.Y+anchor.H/2)
	tt.Frame()
	if !tt.HasText("Follow") {
		t.Fatal("hover did not show preview")
	}
	tt.Advance(300 * time.Millisecond)
	panel, ok := tt.Find("Follow")
	if !ok {
		t.Fatal("missing interactive preview content")
	}
	tt.Move(panel.X+panel.W/2, panel.Y+panel.H/2)
	tt.Advance(300 * time.Millisecond)
	if !tt.HasText("Follow") {
		t.Fatal("preview closed when pointer entered it")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if state.Open {
		t.Fatal("Escape did not dismiss preview")
	}
	tt.Move(2, 2)
	tt.Advance(200 * time.Millisecond)
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if !state.Open {
		t.Fatal("keyboard focus did not show preview")
	}
}

func TestPreviewCardPointerReentersAndActionWorks(t *testing.T) {
	var state PreviewCardState
	var panel ui.Rect
	actions := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().AlignItems(ui.Start).Padding(24).Children(func() {
			anchor := ui.Button(c, "Profile")
			if p := PreviewCard(c, anchor, &state, PreviewCardOptions{
				Label: "Profile preview",
				Body:  func() { ui.Button(c, "Follow").OnClick(func() { actions++ }) },
			}); p.Valid() {
				panel = p.Bounds()
			}
		})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	anchor, _ := tt.Find("Profile")
	tt.Move(anchor.X+anchor.W/2, anchor.Y+anchor.H/2)
	tt.Frame()
	if panel.W != 256 {
		t.Fatalf("preview width=%v, want 256", panel.W)
	}
	tt.Move(2, 2)
	tt.Advance(100 * time.Millisecond)
	if !state.Open {
		t.Fatal("preview closed before the pointer could return")
	}
	tt.Move(panel.X+panel.W/2, panel.Y+panel.H/2)
	tt.Advance(200 * time.Millisecond)
	if !state.Open {
		t.Fatal("preview closed after the pointer entered its panel")
	}
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 || !state.Open {
		t.Fatalf("preview action: actions=%d open=%v", actions, state.Open)
	}
	tt.ClickAt(2, 2)
	if state.Open {
		t.Fatal("outside press did not dismiss the preview")
	}
}

func TestPreviewCardClosesWhenAnchorDisappears(t *testing.T) {
	var state PreviewCardState
	showAnchor := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		if !showAnchor {
			PreviewCard(c, ui.Element{}, &state, PreviewCardOptions{})
			return
		}
		anchor := ui.Button(c, "Profile")
		PreviewCard(c, anchor, &state, PreviewCardOptions{Body: func() { ui.Text(c, "Preview") }})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	anchor, _ := tt.Find("Profile")
	tt.Move(anchor.X+anchor.W/2, anchor.Y+anchor.H/2)
	tt.Frame()
	if !state.Open {
		t.Fatal("preview did not open")
	}
	showAnchor = false
	tt.Frame()
	if state.Open || tt.HasText("Preview") {
		t.Fatal("preview remained open after its anchor disappeared")
	}
	showAnchor = true
	tt.Move(400, 280)
	tt.Frame()
	if state.Open {
		t.Fatal("preview reopened without hover or focus")
	}
}

func TestPreviewCardClosesWhenPointerLeavesWindow(t *testing.T) {
	var state PreviewCardState
	var panel ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Box(c).Fill().Children(func() {
			anchor := ui.Button(c, "Profile").Absolute().Left(350).Top(24)
			if p := PreviewCard(c, anchor, &state, PreviewCardOptions{
				Body: func() { ui.Button(c, "Follow") },
			}); p.Valid() {
				panel = p.Bounds()
			}
		})
	}, 640, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	anchor, ok := tt.Find("Profile")
	if !ok {
		t.Fatal("missing preview anchor")
	}
	tt.Move(anchor.X+anchor.W/2, anchor.Y+anchor.H/2)
	tt.Frame()
	if !state.Open || panel.W == 0 {
		t.Fatal("preview did not open")
	}
	x, y := panel.X+panel.W/2, panel.Y+panel.H/2
	tt.Move(x, y)
	tt.Advance(time.Second)
	if !state.Open {
		t.Fatal("preview did not stay open over its panel")
	}
	// Tester has no PointerLeave event. Shrinking the window around the
	// stationary pointer exercises the same stale-coordinate condition.
	tt.SetSize(int(x)-1, 320)
	tt.Advance(time.Second)
	tt.Frame()
	if state.Open || tt.HasText("Follow") {
		t.Fatal("preview stayed open after the pointer left the window")
	}
}

func TestPreviewCardStaysOpenWhilePanelHasFocus(t *testing.T) {
	var state PreviewCardState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(24).Children(func() {
			anchor := ui.Button(c, "Profile")
			PreviewCard(c, anchor, &state, PreviewCardOptions{
				Body: func() { ui.Button(c, "Follow") },
			})
			ui.Button(c, "Next")
		})
	}, 480, 320)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Follow") {
		t.Fatal("keyboard focus did not enter the preview")
	}
	tt.Move(-1, -1)
	tt.Advance(time.Second)
	if !state.Open || !tt.HasText("Follow") {
		t.Fatal("preview closed while its panel had focus")
	}
	tt.Key(0, ui.KeyTab)
	tt.Advance(time.Second)
	if state.Open {
		t.Fatal("preview stayed open after focus left the panel")
	}
}
