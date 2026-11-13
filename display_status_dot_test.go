package mygoui_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestStatusDotPulsesOnlyWhenAskedAndAllowed pins the motion contract: a
// still dot asks for no frames, a pulsing one repaints without building the
// view again, and reduced motion keeps it still.
func TestStatusDotPulsesOnlyWhenAskedAndAllowed(t *testing.T) {
	for _, tc := range []struct {
		name           string
		pulse, reduced bool
		moves          bool
	}{
		{name: "still"},
		{name: "pulse", pulse: true, moves: true},
		{name: "reduced motion", pulse: true, reduced: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			builds := 0
			tt := ui.NewTesterWith(func(c *ui.Context) {
				builds++
				theme.ApplyMode(c, false)
				mygoui.StatusDot(c, mygoui.StatusDotOptions{
					Status: mygoui.ColorSuccess, Label: "Online", Pulse: tc.pulse,
				})
			}, 200, 60, ui.TesterOptions{Preferences: ui.Preferences{ReduceMotion: tc.reduced, TextScale: 1}})
			if !tt.HasText("Online") {
				t.Fatalf("the label is not shown: %v", tt.Texts())
			}
			if tt.FramePending() != tc.moves {
				t.Fatalf("frame pending = %v, want %v", tt.FramePending(), tc.moves)
			}
			before := append([]byte(nil), tt.Image().Pix...)
			n := builds
			tt.Advance(500 * time.Millisecond)
			if moved := !bytes.Equal(before, tt.Image().Pix); moved != tc.moves {
				t.Fatalf("the dot moved = %v, want %v", moved, tc.moves)
			}
			if builds != n {
				t.Fatalf("the pulse rebuilt the view %d times", builds-n)
			}
		})
	}
}

// TestStatusDotPaintsItsStatusColor pins the color: the dot carries the
// status ink, and the default is the muted one.
func TestStatusDotPaintsItsStatusColor(t *testing.T) {
	status := mygoui.ColorDefault
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		mygoui.StatusDot(c, mygoui.StatusDotOptions{Status: status})
	}, 60, 40)
	if inkCount(tt.Image(), pal.MutedForeground) == 0 {
		t.Fatal("the default dot is not muted")
	}
	status = mygoui.ColorDanger
	tt.Frame()
	tt.Advance(300 * time.Millisecond)
	if inkCount(tt.Image(), pal.Danger) == 0 || inkCount(tt.Image(), pal.MutedForeground) != 0 {
		t.Fatal("a danger dot is not painted in the danger color")
	}
}
