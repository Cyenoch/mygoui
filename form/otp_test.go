package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestOTPPasteFilteringSelectionBackspaceUndoAndComplete(t *testing.T) {
	state := OTPState{}
	completions := 0
	readonly, disabled := false, false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		InputOTP(c, &state, OTPOptions{
			ControlOptions: ControlOptions{Label: "Code", ReadOnly: readonly, Disabled: disabled},
			Length:         4, OnComplete: func(string) { completions++ },
		})
	}, 260, 100)
	tt.Click("Code")
	tt.SetClipboard("12-3a456")
	tt.Command("paste")
	if state.Value != "1234" || completions != 1 {
		t.Fatalf("filtered paste/complete: %q, %d", state.Value, completions)
	}
	tt.Command("undo")
	if state.Value != "" {
		t.Fatalf("filtered paste did not remain one native undo step: %q", state.Value)
	}
	tt.Type("123")
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyBackspace)
	if state.Value != "13" {
		t.Fatalf("native arrow/backspace lost caret: %q", state.Value)
	}
	tt.Type("x9")
	if state.Value != "193" {
		t.Fatalf("filtered middle insert changed native selection: %q", state.Value)
	}
	tt.Key(0, ui.KeyBackspace)
	if state.Value != "13" {
		t.Fatalf("filtered insert caret correction failed: %q", state.Value)
	}
	readonly = true
	tt.Frame()
	tt.Command("selectAll")
	tt.SetClipboard("9999")
	tt.Command("paste")
	if state.Value != "13" {
		t.Fatal("readonly OTP accepted paste")
	}
	readonly, disabled = false, true
	tt.Frame()
	tt.Click("Code")
	tt.Type("9999")
	if state.Value != "13" {
		t.Fatal("disabled OTP accepted text")
	}
}

func TestOTPAlphanumericWaitsForIMECommit(t *testing.T) {
	state := OTPState{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		InputOTP(c, &state, OTPOptions{
			ControlOptions: ControlOptions{Label: "Code"}, Length: 6, Pattern: OTPAlphanumeric,
		})
	}, 280, 100)
	tt.Click("Code")
	tt.Compose("ab", 2)
	if state.Value != "" {
		t.Fatal("OTP filtered or committed IME preedit")
	}
	tt.Type("aB2")
	if state.Value != "aB2" {
		t.Fatalf("alphanumeric IME commit: %q", state.Value)
	}
	tt.SetClipboard("X-9_Z")
	tt.Command("paste")
	if state.Value != "aB2X9Z" {
		t.Fatalf("alphanumeric filtering: %q", state.Value)
	}
}
