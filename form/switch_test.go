package form

import (
	"testing"

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
