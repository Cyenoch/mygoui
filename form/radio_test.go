package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestRadioGroupSelectsOne(t *testing.T) {
	choice := "free"
	var reports []string
	state := &RadioGroupState[string]{
		Selected: &choice,
		OnChange: func(v string) { reports = append(reports, v) },
	}
	tt := ui.NewTester(func(c *ui.Context) {
		RadioGroup(c, RadioGroupOptions{
			Label:       "Plan",
			Description: "Choose one",
		}, func() {
			Radio(c, RadioOptions[string]{State: state, Value: "free", Label: "Free"})
			Radio(c, RadioOptions[string]{State: state, Value: "pro", Label: "Pro"})
			Radio(c, RadioOptions[string]{
				State: state, Value: "enterprise", Label: "Enterprise", Disabled: true,
			})
		})
	}, 360, 240)

	if err := tt.Click("Pro"); err != nil {
		t.Fatal(err)
	}
	if choice != "pro" {
		t.Fatalf("choice is %q, want pro", choice)
	}
	if len(reports) != 1 || reports[0] != "pro" {
		t.Fatalf("OnChange reported %v, want [pro]", reports)
	}
	if !tt.Focused("Pro") {
		t.Fatal("the chosen radio does not have the focus")
	}
	// The arrows move the choice within the radio group.
	tt.Key(0, ui.KeyUp)
	if choice != "free" {
		t.Fatalf("Up chose %q, want free", choice)
	}
	// A disabled option never selects.
	if err := tt.Click("Enterprise"); err != nil {
		t.Fatal(err)
	}
	if choice != "free" {
		t.Fatalf("a disabled radio chose %q", choice)
	}
}

func TestRadioGroupHorizontal(t *testing.T) {
	choice := "a"
	state := &RadioGroupState[string]{Selected: &choice}
	tt := ui.NewTester(func(c *ui.Context) {
		RadioGroup(c, RadioGroupOptions{Orientation: GroupHorizontal}, func() {
			Radio(c, RadioOptions[string]{State: state, Value: "a", Label: "Alpha"})
			Radio(c, RadioOptions[string]{State: state, Value: "b", Label: "Beta"})
		})
	}, 400, 160)
	a, okA := tt.Find("Alpha")
	b, okB := tt.Find("Beta")
	if !okA || !okB {
		t.Fatal("radios are not rendered")
	}
	if a.X == b.X {
		t.Fatalf("horizontal radios stack at %v and %v", a, b)
	}
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if choice != "b" {
		t.Fatalf("choice is %q, want b", choice)
	}
}

func TestRadioWithoutStateIsDisabled(t *testing.T) {
	var reports []string
	tt := ui.NewTester(func(c *ui.Context) {
		Radio(c, RadioOptions[string]{
			Value: "orphan",
			Label: "Orphan",
			OnChange: func(v string) {
				reports = append(reports, v)
			},
		})
	}, 240, 100)
	if err := tt.Click("Orphan"); err != nil {
		t.Fatal(err)
	}
	if len(reports) != 0 {
		t.Fatalf("a radio with no state reported %v", reports)
	}
}

func TestDisabledRadioGroup(t *testing.T) {
	choice := "a"
	state := &RadioGroupState[string]{Selected: &choice}
	tt := ui.NewTester(func(c *ui.Context) {
		RadioGroup(c, RadioGroupOptions{Label: "Locked", Disabled: true}, func() {
			Radio(c, RadioOptions[string]{State: state, Value: "a", Label: "Alpha"})
			Radio(c, RadioOptions[string]{State: state, Value: "b", Label: "Beta"})
		})
	}, 300, 200)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if choice != "a" {
		t.Fatalf("a disabled group chose %q", choice)
	}
}
