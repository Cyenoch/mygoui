package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// formTextTester runs a form specimen in a filled column and applies the
// HeroUI theme, as an application's view root does.
func formTextTester(t *testing.T, build func(c *ui.Context)) *ui.Tester {
	t.Helper()
	return ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Gap(12).Children(func() { build(c) })
	}, 360, 320)
}

func formTextClick(t *testing.T, tester *ui.Tester, label string) {
	t.Helper()
	if err := tester.Click(label); err != nil {
		t.Fatalf("click %q: %v", label, err)
	}
}

func TestInputCommitsTypedValue(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		Input(c, &value, InputOptions{Label: "Name", Placeholder: "Your name"})
	})
	formTextClick(t, tester, "Name")
	tester.Type("Ada")
	if value != "Ada" {
		t.Fatalf("input committed %q", value)
	}
}

// TestDisabledTextControlsTakeNoInput checks every text-entry control wires
// its Disabled option through to the native editor: a click focuses nothing
// and typing commits nothing.
func TestDisabledTextControlsTakeNoInput(t *testing.T) {
	for _, tc := range []struct {
		name, label string
		build       func(c *ui.Context) func() string
	}{
		{"input", "Locked", func(c *ui.Context) func() string {
			var value string
			Input(c, &value, InputOptions{Label: "Locked", Disabled: true})
			return func() string { return value }
		}},
		{"input group", "Locked", func(c *ui.Context) func() string {
			var value string
			InputGroup(c, &value, InputGroupOptions{Label: "Locked", Disabled: true})
			return func() string { return value }
		}},
		{"search field", "Find", func(c *ui.Context) func() string {
			var value string
			SearchField(c, &value, SearchFieldOptions{Label: "Find", Disabled: true})
			return func() string { return value }
		}},
		{"otp", "Verification code", func(c *ui.Context) func() string {
			var state OTPState
			InputOTP(c, &state, OTPOptions{Length: 4, Label: "Verification code", Disabled: true})
			return func() string { return state.Value }
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var value func() string
			tester := formTextTester(t, func(c *ui.Context) { value = tc.build(c) })
			formTextClick(t, tester, tc.label)
			tester.Type("ignored")
			if got := value(); got != "" {
				t.Fatalf("a disabled %s committed %q", tc.name, got)
			}
		})
	}
}

func TestInputInvalidAnnouncesErrorAndStaysEditable(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		Input(c, &value, InputOptions{
			Label: "Email", Error: "Enter an email address", Invalid: true,
		})
	})
	if !tester.HasText("Enter an email address") {
		t.Fatal("invalid input shows no error")
	}
	formTextClick(t, tester, "Email")
	tester.Type("a@b")
	if value != "a@b" {
		t.Fatalf("invalid input committed %q", value)
	}
}

func TestInputReadOnlyKeepsValue(t *testing.T) {
	value := "fixed"
	tester := formTextTester(t, func(c *ui.Context) {
		Input(c, &value, InputOptions{Label: "Code", ReadOnly: true})
	})
	formTextClick(t, tester, "Code")
	tester.Type("x")
	if value != "fixed" {
		t.Fatalf("read-only input committed %q", value)
	}
}
