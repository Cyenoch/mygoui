package mygoui_test

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestStateMotionFollowsTheThemeDuration pins that a control is timed by
// the published theme: 60 ms into a state change it is still moving under
// the built-in state duration and at rest under a 40 ms one.
func TestStateMotionFollowsTheThemeDuration(t *testing.T) {
	moving := func() bool {
		on := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{Label: "Bold", Selected: &on})
		}, 200, 80)
		tt.Advance(time.Second)
		on = true
		tt.Frame()
		tt.Advance(60 * time.Millisecond)
		return tt.FramePending()
	}
	if !moving() {
		t.Fatal("under the built-in theme the toggle is at rest 60 ms into its change")
	}

	th := theme.Default()
	th.Motion.State = 40 * time.Millisecond
	useTheme(t, th)
	if moving() {
		t.Fatal("with a 40 ms state duration the toggle still moves after 60 ms")
	}
}

// useTheme publishes th and restores the built-in theme when the test ends.
func useTheme(t *testing.T, th theme.Theme) {
	t.Helper()
	t.Cleanup(func() {
		if err := theme.Use(theme.Default()); err != nil {
			t.Error(err)
		}
	})
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}
}

// TestZeroStateDurationSettlesAtOnce pins what a zero motion token means: no
// motion, not the toolkit's default spring. One 60 Hz frame after the change
// the toggle is at rest.
func TestZeroStateDurationSettlesAtOnce(t *testing.T) {
	th := theme.Default()
	th.Motion.State = 0
	useTheme(t, th)

	on := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.ToggleButton(c, mygoui.ToggleButtonOptions{Label: "Bold", Selected: &on})
	}, 200, 80)
	tt.Advance(time.Second)
	on = true
	tt.Frame()
	tt.Advance(20 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("with a zero state duration the toggle still moves a frame after its change")
	}
}

// TestStatusSolidsHoverAtTheThemeOpacity pins that every status hovers its
// solid fill the same way, so one theme opacity governs them all.
func TestStatusSolidsHoverAtTheThemeOpacity(t *testing.T) {
	th := theme.Default()
	th.Effects.HoverOpacity = 0.5
	useTheme(t, th)

	ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		for _, status := range []mygoui.StatusColor{
			mygoui.ColorDefault, mygoui.ColorAccent, mygoui.ColorSuccess,
			mygoui.ColorWarning, mygoui.ColorDanger, mygoui.ColorInfo,
		} {
			p := mygoui.StatusColors(c, status)
			if want := p.Solid.Alpha(0.5); p.Hover != want {
				t.Errorf("status %v hovers at %v, want its solid at the theme opacity %v", status, p.Hover, want)
			}
		}
	}, 200, 80).Frame()
}
