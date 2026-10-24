package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestInputOTPFiltersASCIIAndCompletesOnce(t *testing.T) {
	var state OTPState
	var completed, changes []string
	tester := formTextTester(t, func(c *ui.Context) {
		InputOTP(c, &state, OTPOptions{
			Length: 4,
			Label:  "Verification code",
			OnComplete: func(value string) {
				completed = append(completed, value)
			},
			OnChange: func(value string) {
				changes = append(changes, value)
			},
		})
	})
	formTextClick(t, tester, "Verification code")
	tester.Type("12é34")
	if state.Value != "1234" {
		t.Fatalf("OTP kept %q", state.Value)
	}
	if !state.Complete || len(completed) != 1 || completed[0] != "1234" {
		t.Fatalf("completion: %v, complete %v", completed, state.Complete)
	}
	if len(changes) != 1 || changes[0] != "1234" {
		t.Fatalf("changes: %v", changes)
	}

	tester.Type("9")
	if state.Value != "1234" || len(completed) != 1 {
		t.Fatalf("OTP grew to %q, completions %v", state.Value, completed)
	}

	tester.Key(0, ui.KeyBackspace)
	if state.Value != "123" || state.Complete {
		t.Fatalf("after backspace: %q, complete %v", state.Value, state.Complete)
	}

	tester.Key(ui.Cmd, ui.KeyA)
	tester.Type("5678")
	if state.Value != "5678" || len(completed) != 2 || completed[1] != "5678" {
		t.Fatalf("replacement: %q, completions %v", state.Value, completed)
	}
}

func TestInputOTPFilterOption(t *testing.T) {
	var state OTPState
	var completed []string
	tester := formTextTester(t, func(c *ui.Context) {
		InputOTP(c, &state, OTPOptions{
			Length: 4,
			Label:  "Digits",
			Filter: func(r rune) bool { return r >= '0' && r <= '9' },
			OnComplete: func(value string) {
				completed = append(completed, value)
			},
		})
	})
	formTextClick(t, tester, "Digits")
	tester.Type("a1b2")
	if state.Value != "12" {
		t.Fatalf("digits-only OTP kept %q", state.Value)
	}
	tester.Type("34")
	if state.Value != "1234" || len(completed) != 1 || completed[0] != "1234" {
		t.Fatalf("value %q, completed %v", state.Value, completed)
	}
}

func TestInputOTPCountsRunesWithCustomFilter(t *testing.T) {
	var state OTPState
	var completed []string
	tester := formTextTester(t, func(c *ui.Context) {
		InputOTP(c, &state, OTPOptions{
			Length: 2,
			Label:  "Nickname code",
			Filter: func(r rune) bool { return r != ' ' },
			OnComplete: func(value string) {
				completed = append(completed, value)
			},
		})
	})
	formTextClick(t, tester, "Nickname code")
	tester.Type("好a")
	if state.Value != "好a" || !state.Complete {
		t.Fatalf("two runes: value %q, complete %v", state.Value, state.Complete)
	}
	if len(completed) != 1 || completed[0] != "好a" {
		t.Fatalf("completion: %v", completed)
	}
}

func TestInputOTPDisabledTakesNoInput(t *testing.T) {
	var state OTPState
	tester := formTextTester(t, func(c *ui.Context) {
		InputOTP(c, &state, OTPOptions{Length: 4, Label: "Verification code", Disabled: true})
	})
	if err := tester.Click("Verification code"); err != nil {
		t.Fatalf("click disabled label: %v", err)
	}
	tester.Type("1234")
	if state.Value != "" {
		t.Fatalf("disabled OTP kept %q", state.Value)
	}
}
