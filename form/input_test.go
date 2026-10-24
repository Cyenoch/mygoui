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

func TestInputDisabledTakesNoInput(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		Input(c, &value, InputOptions{Label: "Locked", Disabled: true})
	})
	if err := tester.Click("Locked"); err != nil {
		t.Fatalf("click disabled label: %v", err)
	}
	tester.Type("ignored")
	if value != "" {
		t.Fatalf("disabled input committed %q", value)
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
