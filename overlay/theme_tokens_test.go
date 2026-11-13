package overlay

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestOverlayMotionFollowsTheThemeAndStillYieldsToReducedMotion pins both
// sides of an overlay token: a longer dialog duration keeps a modal moving
// past the built-in one, and reduced motion still shows it at once.
func TestOverlayMotionFollowsTheThemeAndStillYieldsToReducedMotion(t *testing.T) {
	th := theme.Default()
	th.Motion.Dialog = 2 * time.Second
	t.Cleanup(func() {
		if err := theme.Use(theme.Default()); err != nil {
			t.Error(err)
		}
	})
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}

	for _, reduced := range []bool{false, true} {
		open := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Button(c, "Open").OnClick(func() { open = true })
			Modal(c, &open, ModalOptions{Title: "Settings", Body: func() { ui.Button(c, "Confirm") }})
		}, 640, 480)
		tt.SetPreferences(ui.Preferences{ReduceMotion: reduced, TextScale: 1})
		if err := tt.Click("Open"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		tt.Advance(2 * dialogMotionDuration)
		if moving := tt.FramePending(); moving == reduced {
			t.Errorf("reduced motion %v: modal moving %v after twice the built-in duration", reduced, moving)
		}
	}
}
