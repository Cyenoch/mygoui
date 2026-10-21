package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestFieldLabelActivatesCheckboxAndDisabledFieldsetBlocksIt(t *testing.T) {
	value, disabled := false, false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		FieldSet(c, FieldSetOptions{Legend: "Privacy", Disabled: disabled}, func() {
			Field(c, FieldOptions{
				Label: "Keep history", Description: "Retain conversations on this machine", Error: "Choose a preference",
			}, func() { Checkbox(c, &value, CheckboxOptions{}) })
		})
	}, 400, 220)
	if !tt.HasText("Retain conversations on this machine") || !tt.HasText("Choose a preference") {
		t.Fatal("native Field omitted connected description or validation")
	}
	tt.Click("Keep history")
	if !value {
		t.Fatal("native field label failed to activate associated checkbox")
	}
	disabled = true
	tt.Frame()
	tt.Click("Keep history")
	if !value {
		t.Fatal("disabled native fieldset left label activation enabled")
	}
}

func TestExplicitLabelNamesAndFocusesBuiltHandle(t *testing.T) {
	value := ""
	var target ui.Handle
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Input(c, &value, InputOptions{ControlOptions: ControlOptions{Handle: &target}})
		Label(c, "Display name", &target)
	}, 300, 140)
	tt.Click("Display name")
	tt.Type("Ada")
	if value != "Ada" {
		t.Fatal("explicit native Handle label did not focus editor")
	}
}
