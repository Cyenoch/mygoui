package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestNumberFieldClampsStepsAndRejectsText(t *testing.T) {
	lo, hi := 0.0, 10.0
	state := NumberFieldState{Value: 5}
	changes := 0
	tester := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{
			Label:    "Count",
			Min:      &lo,
			Max:      &hi,
			Step:     1,
			OnChange: func() { changes++ },
		})
	})

	formTextClick(t, tester, "Count")
	tester.Type("50")
	if state.Value != hi || state.Text != "10" {
		t.Fatalf("clamped to %v (%q)", state.Value, state.Text)
	}

	formTextClick(t, tester, "Decrease")
	if state.Value != 9 {
		t.Fatalf("decrement gave %v", state.Value)
	}

	formTextClick(t, tester, "Count")
	tester.Key(ui.Cmd, ui.KeyA)
	tester.Type("abc")
	if state.Value != 9 || state.Text != "9" {
		t.Fatalf("invalid text gave %v (%q)", state.Value, state.Text)
	}

	tester.Key(ui.Cmd, ui.KeyA)
	tester.Type("-5")
	if state.Value != lo || state.Text != "0" {
		t.Fatalf("lower clamp gave %v (%q)", state.Value, state.Text)
	}

	formTextClick(t, tester, "Decrease")
	if state.Value != lo {
		t.Fatalf("decrement below minimum gave %v", state.Value)
	}
	formTextClick(t, tester, "Increase")
	if state.Value != 1 {
		t.Fatalf("increment gave %v", state.Value)
	}
	if changes != 4 {
		t.Fatalf("reported %d changes", changes)
	}
}

func TestNumberFieldFractionalStep(t *testing.T) {
	state := NumberFieldState{Value: 0}
	tester := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{Label: "Ratio", Step: 0.5})
	})
	formTextClick(t, tester, "Increase")
	if state.Value != 0.5 || state.Text != "0.5" {
		t.Fatalf("fractional step gave %v (%q)", state.Value, state.Text)
	}

	formTextClick(t, tester, "Ratio")
	tester.Key(0, ui.KeyUp)
	if state.Value != 1 {
		t.Fatalf("Up stepped to %v", state.Value)
	}
	tester.Key(0, ui.KeyDown)
	if state.Value != 0.5 {
		t.Fatalf("Down stepped to %v", state.Value)
	}
}
