package form

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestSwitchTogglesAndReports(t *testing.T) {
	on := false
	var reports []bool
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			Switch(c, SwitchOptions{
				On:       &on,
				Label:    "Wi-Fi",
				Size:     SwitchMedium,
				OnChange: func(v bool) { reports = append(reports, v) },
			})
		})
	}, 280, 120)

	if err := tt.Click("Wi-Fi"); err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Fatal("a click did not turn the switch on")
	}
	if len(reports) != 1 || !reports[0] {
		t.Fatalf("OnChange reported %v, want [true]", reports)
	}
	if !tt.Focused("Wi-Fi") {
		t.Fatal("the switch does not have the focus")
	}
	tt.Key(0, ui.KeySpace)
	if on || len(reports) != 2 || reports[1] {
		t.Fatalf("Space: on %v, reports %v", on, reports)
	}
}

func TestSwitchLabelPositionAndDisabled(t *testing.T) {
	small, large, locked := false, true, false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).Children(func() {
			Switch(c, SwitchOptions{On: &small, Label: "Small", Size: SwitchSmall})
			Switch(c, SwitchOptions{
				On: &large, Label: "Before", Size: SwitchLarge,
				LabelPosition: SwitchLabelBefore,
			})
			Switch(c, SwitchOptions{On: &locked, Label: "Locked", Disabled: true})
		})
	}, 320, 220)

	if err := tt.Click("Locked"); err != nil {
		t.Fatal(err)
	}
	if locked {
		t.Fatal("a disabled switch toggled")
	}
	// The label-before row still toggles from the label.
	if err := tt.Click("Before"); err != nil {
		t.Fatal(err)
	}
	if large {
		t.Fatal("the label-before switch did not toggle")
	}
}

func TestSwitchGroup(t *testing.T) {
	a, b := false, true
	tt := ui.NewTester(func(c *ui.Context) {
		SwitchGroup(c, SwitchGroupOptions{
			Label:       "Alerts",
			Description: "Choose alerts",
			Orientation: GroupVertical,
		}, func() {
			Switch(c, SwitchOptions{On: &a, Label: "Email"})
			Switch(c, SwitchOptions{On: &b, Label: "SMS"})
		})
	}, 320, 240)
	if !tt.HasText("Alerts") {
		t.Fatalf("group label missing: %q", tt.Texts())
	}
	if err := tt.Click("Email"); err != nil {
		t.Fatal(err)
	}
	if !a {
		t.Fatal("a group switch did not toggle")
	}
}

// TestSwitchThumbStretchesWhilePressed pins coss press feedback: a held
// switch widens its thumb over the track, and reduced motion still shows the
// held state without animating it.
func TestSwitchThumbStretchesWhilePressed(t *testing.T) {
	for _, reduced := range []bool{false, true} {
		on := false
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			ui.Column(c).Padding(20).Children(func() {
				Switch(c, SwitchOptions{On: &on, Label: "Wi-Fi"})
			})
		}, 280, 120)
		tt.SetPreferences(ui.Preferences{TextScale: 1, ReduceMotion: reduced})
		r, ok := tt.Find("Wi-Fi")
		if !ok {
			t.Fatal("switch not found")
		}
		y := r.Y + r.H/2
		// Track pixels visible beside the thumb shrink as the thumb widens.
		trackPixels := func() int {
			img, n := tt.Image(), 0
			for x := int(r.X); x < int(r.X)+40; x++ {
				if v := img.RGBAAt(x, int(y)).R; v > 200 && v < 240 {
					n++
				}
			}
			return n
		}
		rest := trackPixels()
		tt.Press(r.X+8, y)
		tt.Advance(300 * time.Millisecond)
		held := trackPixels()
		tt.Release(r.X+8, y)
		tt.Advance(300 * time.Millisecond)
		if held >= rest {
			t.Fatalf("reduced=%v: held thumb did not widen: track pixels rest=%d held=%d", reduced, rest, held)
		}
	}
}
