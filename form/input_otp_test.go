package form

import (
	"slices"
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

func TestInputOTPReportsAcceptedChanges(t *testing.T) {
	var state OTPState
	var changes, completed []string
	tester := formTextTester(t, func(c *ui.Context) {
		InputOTP(c, &state, OTPOptions{
			Length:     2,
			Label:      "Verification code",
			OnChange:   func(value string) { changes = append(changes, value) },
			OnComplete: func(value string) { completed = append(completed, value) },
		})
	})
	formTextClick(t, tester, "Verification code")
	tester.Type("1")
	tester.Type("2")
	tester.Type("3") // The code is already full.
	tester.Key(0, ui.KeyBackspace)
	tester.Type("é") // The filter rejects the character.
	if !slices.Equal(changes, []string{"1", "12", "1"}) {
		t.Fatalf("accepted changes = %q", changes)
	}
	if !slices.Equal(completed, []string{"12"}) || state.Value != "1" || state.Complete {
		t.Fatalf("value %q, complete %v, completions %q", state.Value, state.Complete, completed)
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

func TestInputOTPReadOnlyPreservesCode(t *testing.T) {
	state := OTPState{Value: "12"}
	var changes, completed int
	tester := formTextTester(t, func(c *ui.Context) {
		InputOTP(c, &state, OTPOptions{
			Length: 4, Label: "Verification code", ReadOnly: true,
			OnChange: func(string) { changes++ }, OnComplete: func(string) { completed++ },
		})
	})
	formTextClick(t, tester, "Verification code")
	tester.Type("34")
	tester.Key(0, ui.KeyBackspace)
	if state.Value != "12" || changes != 0 || completed != 0 {
		t.Fatalf("read-only: value %q changes %d completions %d", state.Value, changes, completed)
	}
}
