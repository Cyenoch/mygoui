package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestCheckboxTogglesAndReports(t *testing.T) {
	checked := false
	var reports []bool
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			Checkbox(c, CheckboxOptions{
				Checked: &checked,
				Label:   "Sync files",
				OnChange: func(v bool) {
					reports = append(reports, v)
				},
			})
		})
	}, 280, 120)

	if err := tt.Click("Sync files"); err != nil {
		t.Fatal(err)
	}
	if !checked {
		t.Fatal("a click did not check the box")
	}
	if len(reports) != 1 || !reports[0] {
		t.Fatalf("OnChange reported %v, want [true]", reports)
	}
	if !tt.Focused("Sync files") {
		t.Fatal("the clicked box does not have the focus")
	}
	tt.Key(0, ui.KeySpace)
	if checked {
		t.Fatal("Space did not uncheck the focused box")
	}
	if len(reports) != 2 || reports[1] {
		t.Fatalf("Space reported %v, want [true false]", reports)
	}
}

func TestCheckboxIndeterminateAndDisabled(t *testing.T) {
	checked, mixed := false, true
	off := false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(8).Children(func() {
			Checkbox(c, CheckboxOptions{
				Checked:       &checked,
				Indeterminate: &mixed,
				Label:         "All channels",
			})
			Checkbox(c, CheckboxOptions{
				Checked:  &off,
				Label:    "Locked",
				Disabled: true,
			})
		})
	}, 280, 160)

	if err := tt.Click("All channels"); err != nil {
		t.Fatal(err)
	}
	if !checked || mixed {
		t.Fatalf("activating a mixed box: checked %v, mixed %v", checked, mixed)
	}
	if err := tt.Click("Locked"); err != nil {
		t.Fatal(err)
	}
	if off {
		t.Fatal("a disabled box toggled")
	}
}

func TestCheckboxGroupSelectAll(t *testing.T) {
	email, sms := false, false
	state := &CheckboxGroupState{Items: []*bool{&email, &sms}}
	tt := ui.NewTester(func(c *ui.Context) {
		CheckboxGroup(c, CheckboxGroupOptions{
			Label:          "Channels",
			State:          state,
			SelectAllLabel: "All channels",
			Orientation:    GroupVertical,
		}, func() {
			Checkbox(c, CheckboxOptions{Checked: &email, Label: "Email"})
			Checkbox(c, CheckboxOptions{Checked: &sms, Label: "SMS"})
		})
	}, 320, 240)

	if err := tt.Click("Email"); err != nil {
		t.Fatal(err)
	}
	if !email || sms {
		t.Fatalf("first item: email %v, sms %v", email, sms)
	}
	if state.Checked() || !state.Indeterminate() {
		t.Fatalf("one of two: checked %v, mixed %v", state.Checked(), state.Indeterminate())
	}
	if err := tt.Click("SMS"); err != nil {
		t.Fatal(err)
	}
	if !state.Checked() || state.Indeterminate() {
		t.Fatalf("both items: checked %v, mixed %v", state.Checked(), state.Indeterminate())
	}
	if err := tt.Click("All channels"); err != nil {
		t.Fatal(err)
	}
	if email || sms {
		t.Fatalf("select all did not clear: email %v, sms %v", email, sms)
	}
}

func TestCheckboxGroupHorizontal(t *testing.T) {
	a, b := true, false
	tt := ui.NewTester(func(c *ui.Context) {
		CheckboxGroup(c, CheckboxGroupOptions{
			Label:       "Channels",
			Orientation: GroupHorizontal,
			Description: "Pick any",
		}, func() {
			Checkbox(c, CheckboxOptions{Checked: &a, Label: "Email"})
			Checkbox(c, CheckboxOptions{Checked: &b, Label: "SMS"})
		})
	}, 400, 160)
	if !tt.HasText("Channels") || !tt.HasText("Pick any") {
		t.Fatalf("group texts: %q", tt.Texts())
	}
	r1, ok1 := tt.Find("Email")
	r2, ok2 := tt.Find("SMS")
	if !ok1 || !ok2 {
		t.Fatal("group items are not rendered")
	}
	if r1.X == r2.X {
		t.Fatalf("horizontal items stack at %v and %v", r1, r2)
	}
}

func TestCheckboxError(t *testing.T) {
	checked := false
	tt := ui.NewTester(func(c *ui.Context) {
		Checkbox(c, CheckboxOptions{
			Checked: &checked,
			Label:   "Accept",
			Invalid: true,
			Error:   "This field is required",
		})
	}, 300, 120)
	if !tt.HasText("This field is required") {
		t.Fatalf("error missing from %q", tt.Texts())
	}
}
