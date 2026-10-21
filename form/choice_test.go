package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestCheckboxMixedAndReadOnlyActivation(t *testing.T) {
	value, mixed := false, true
	readonly, disabled := false, false
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Checkbox(c, &value, CheckboxOptions{
			ControlOptions: ControlOptions{Label: "All files", ReadOnly: readonly, Disabled: disabled},
			Mixed:          &mixed, OnChange: func(bool) { changes++ },
		})
	}, 280, 90)
	tt.Click("All files")
	if !value || mixed || changes != 1 {
		t.Fatalf("mixed activation must check and clear mixed: %v/%v/%d", value, mixed, changes)
	}
	readonly = true
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if !value || changes != 1 {
		t.Fatal("read-only checkbox changed caller state")
	}
	readonly, disabled = false, true
	tt.Frame()
	tt.Click("All files")
	if !value || changes != 1 {
		t.Fatal("disabled checkbox activated")
	}
}

func TestRadioGroupUsesNativeRovingFocusAndDisabledSkip(t *testing.T) {
	value := "one"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		RadioGroup(c, &value, RadioGroupOptions{ControlOptions: ControlOptions{Label: "Plan"}}, []RadioItem{
			{Value: "one", Label: "One"},
			{Value: "two", Label: "Two", Disabled: true},
			{Value: "three", Label: "Three"},
		})
		ui.Button(c, "Next")
	}, 280, 200)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("One") {
		t.Fatal("radio Tab entry did not focus selected option")
	}
	tt.Key(0, ui.KeyDown)
	if value != "three" || !tt.Focused("Three") {
		t.Fatalf("native arrow navigation failed disabled skip: %q", value)
	}
	tt.Key(0, ui.KeyHome)
	if value != "one" {
		t.Fatal("native Home did not choose first enabled radio")
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Next") {
		t.Fatal("radio group consumed more than one Tab stop")
	}
}

func TestSwitchReadOnlyAndDisabled(t *testing.T) {
	value, readonly, disabled := false, false, false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Switch(c, &value, SwitchOptions{ControlOptions: ControlOptions{
			Label: "Alerts", ReadOnly: readonly, Disabled: disabled,
		}})
	}, 280, 90)
	tt.Click("Alerts")
	if !value {
		t.Fatal("native switch did not activate")
	}
	readonly = true
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if !value {
		t.Fatal("read-only switch changed")
	}
	readonly, disabled = false, true
	tt.Frame()
	tt.Click("Alerts")
	if !value {
		t.Fatal("disabled switch changed")
	}
}
