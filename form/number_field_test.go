package form

import (
	"testing"
	"time"

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
	tester.Key(ui.Cmd, ui.KeyA)
	tester.Type("50")
	tester.Key(0, ui.KeyEnter)
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
	tester.Key(0, ui.KeyEnter)
	if state.Value != 9 || state.Text != "9" {
		t.Fatalf("invalid text gave %v (%q)", state.Value, state.Text)
	}

	tester.Key(ui.Cmd, ui.KeyA)
	tester.Type("-5")
	tester.Key(0, ui.KeyEnter)
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

func TestNumberFieldClampsOnCommitNotPerKey(t *testing.T) {
	lo, hi := 10.0, 100.0
	state := NumberFieldState{Value: 20}
	var after string
	tt := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{Label: "Count", Min: &lo, Max: &hi})
		Input(c, &after, InputOptions{Label: "After"})
	})
	formTextClick(t, tt, "Count")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("50")
	if state.Value != 50 || state.Text != "50" {
		t.Fatalf("typing 50 gave %v (%q)", state.Value, state.Text)
	}
	tt.Type("0")
	if state.Value != 50 || state.Text != "500" {
		t.Fatalf("typing past Max gave %v (%q) before the commit", state.Value, state.Text)
	}
	tt.Key(0, ui.KeyTab)
	if state.Value != 100 || state.Text != "100" {
		t.Fatalf("leaving the field committed %v (%q), want 100", state.Value, state.Text)
	}
}

func TestNumberFieldSteppersKeepFocusInEditor(t *testing.T) {
	state := NumberFieldState{Value: 1}
	var after string
	tt := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{Label: "Count"})
		Input(c, &after, InputOptions{Label: "After"})
	})
	formTextClick(t, tt, "Count")
	formTextClick(t, tt, "Increase")
	if state.Value != 2 || !tt.Focused("Count") {
		t.Fatalf("Increase gave %v, editor focused %v", state.Value, tt.Focused("Count"))
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab stopped on a stepper")
	}
}

func TestNumberFieldPageKeysStepTenSteps(t *testing.T) {
	state := NumberFieldState{Value: 5}
	tt := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{Label: "Count", Step: 0.5})
	})
	formTextClick(t, tt, "Count")
	tt.Key(0, ui.KeyPageUp)
	if state.Value != 10 {
		t.Fatalf("PageUp gave %v", state.Value)
	}
	tt.Key(0, ui.KeyPageDown)
	if state.Value != 5 {
		t.Fatalf("PageDown gave %v", state.Value)
	}
}

func TestNumberFieldStepperRepeatsWhileHeld(t *testing.T) {
	state := NumberFieldState{Value: 0}
	tt := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{Label: "Count"})
	})
	r, ok := tt.Find("Increase")
	if !ok {
		t.Fatal("no Increase stepper")
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	if state.Value != 1 {
		t.Fatalf("pressing stepped to %v, want 1", state.Value)
	}
	tt.Advance(time.Second)
	tt.Release(x, y)
	held := state.Value
	if held < 5 {
		t.Fatalf("holding for a second stepped to %v", held)
	}
	tt.Advance(time.Second)
	if state.Value != held {
		t.Fatalf("the stepper kept stepping after the release: %v", state.Value)
	}
}

func TestNumberFieldAdoptsProgrammaticValueWhileFocused(t *testing.T) {
	state := NumberFieldState{Value: 5}
	changes := 0
	var after string
	tt := formTextTester(t, func(c *ui.Context) {
		NumberField(c, &state, NumberFieldOptions{Label: "Count", OnChange: func() { changes++ }})
		Input(c, &after, InputOptions{Label: "After"})
	})
	formTextClick(t, tt, "Count")
	state.Value = 42
	tt.Frame()
	if state.Text != "42" {
		t.Fatalf("focused editor shows %q after the application set 42", state.Text)
	}
	tt.Key(0, ui.KeyTab)
	if state.Value != 42 || state.Text != "42" || changes != 0 {
		t.Fatalf("leaving the field gave %v (%q), %d changes; want 42, none", state.Value, state.Text, changes)
	}
}
