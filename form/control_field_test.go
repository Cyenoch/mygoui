package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestControlFieldSwitchRow(t *testing.T) {
	on := false
	var reports []bool
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			ControlField(c, ControlFieldOptions{
				Checked:     &on,
				Label:       "Notifications",
				Description: "Receive account alerts",
				Kind:        ControlFieldSwitch,
				OnChange:    func(v bool) { reports = append(reports, v) },
			})
		})
	}, 360, 140)

	if !tt.HasText("Notifications") || !tt.HasText("Receive account alerts") {
		t.Fatalf("row texts: %q", tt.Texts())
	}
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	if !on {
		t.Fatal("the pressable row did not toggle the switch")
	}
	if len(reports) != 1 || !reports[0] {
		t.Fatalf("OnChange reported %v", reports)
	}
	tt.Key(0, ui.KeySpace)
	if on {
		t.Fatal("Space did not toggle the row")
	}
}

func TestControlFieldCheckboxRowAndDisabled(t *testing.T) {
	agree, locked := false, false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).Children(func() {
			ControlField(c, ControlFieldOptions{
				Checked: &agree,
				Label:   "I agree",
				Kind:    ControlFieldCheckbox,
			})
			ControlField(c, ControlFieldOptions{
				Checked:  &locked,
				Label:    "Locked",
				Kind:     ControlFieldCheckbox,
				Disabled: true,
			})
		})
	}, 360, 180)

	if err := tt.Click("I agree"); err != nil {
		t.Fatal(err)
	}
	if !agree {
		t.Fatal("the checkbox row did not check")
	}
	if err := tt.Click("Locked"); err != nil {
		t.Fatal(err)
	}
	if locked {
		t.Fatal("a disabled control field toggled")
	}
}
