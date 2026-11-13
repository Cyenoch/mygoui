package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/egoist/mygo/ui"
)

func TestPasswordFieldToggleRevealsTypedText(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		PasswordField(c, &value, PasswordFieldOptions{Label: "Password"})
	})
	formTextClick(t, tester, "Password")
	tester.Type("s3cret")
	if value != "s3cret" {
		t.Fatalf("password committed %q", value)
	}
	if tester.HasText("s3cret") {
		t.Fatal("a concealed password is readable")
	}
	toggle := locale.M().ShowPassword
	if info, ok := tester.AccessByRole(ui.RoleToggleButton, toggle); !ok || info.Checked {
		t.Fatalf("concealed toggle: found %v, pressed %v", ok, info.Checked)
	}
	formTextClick(t, tester, toggle)
	info, _ := tester.AccessByRole(ui.RoleTextField, "Password")
	if info.Value != "s3cret" {
		t.Fatalf("revealed field reads %q, want the password", info.Value)
	}
	if info, _ := tester.AccessByRole(ui.RoleToggleButton, toggle); !info.Checked {
		t.Fatal("the toggle is not pressed while the password shows")
	}
}

func TestPasswordFieldFollowsCallerRevealedFlag(t *testing.T) {
	value, revealed := "hunter2", true
	tester := formTextTester(t, func(c *ui.Context) {
		PasswordField(c, &value, PasswordFieldOptions{Label: "Password", Revealed: &revealed})
	})
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Password"); info.Value != "hunter2" {
		t.Fatalf("a revealed field reads %q", info.Value)
	}
	formTextClick(t, tester, locale.M().ShowPassword)
	if revealed {
		t.Fatal("the toggle left the caller's flag set")
	}
	if info, _ := tester.AccessByRole(ui.RoleTextField, "Password"); info.Value == "hunter2" {
		t.Fatal("a concealed field still reads the password")
	}
}

func TestPasswordFieldDisabledToggleDoesNothing(t *testing.T) {
	value, revealed := "hunter2", false
	tester := formTextTester(t, func(c *ui.Context) {
		PasswordField(c, &value, PasswordFieldOptions{Label: "Password", Revealed: &revealed, Disabled: true})
	})
	_ = tester.Click(locale.M().ShowPassword)
	tester.Frame()
	if revealed {
		t.Fatal("a disabled field revealed its password")
	}
}

func TestPasswordFieldConcealsAgainWhenFocusLeaves(t *testing.T) {
	var value, other string
	tester := formTextTester(t, func(c *ui.Context) {
		Input(c, &other, InputOptions{Label: "Other"})
		PasswordField(c, &value, PasswordFieldOptions{Label: "Password"})
	})
	toggle := locale.M().ShowPassword
	formTextClick(t, tester, "Password")
	tester.Type("s3cret")
	formTextClick(t, tester, toggle)
	if info, _ := tester.AccessByRole(ui.RoleToggleButton, toggle); !info.Checked {
		t.Fatal("the toggle did not reveal the password")
	}
	formTextClick(t, tester, "Other")
	tester.Frame()
	if info, _ := tester.AccessByRole(ui.RoleToggleButton, toggle); info.Checked {
		t.Fatal("the password stays revealed after the focus left the field")
	}
}

func TestPasswordFieldReportsStrengthAsTyped(t *testing.T) {
	var value string
	tester := formTextTester(t, func(c *ui.Context) {
		PasswordField(c, &value, PasswordFieldOptions{
			Label: "Password",
			Strength: func(v string) (int, string) {
				if len(v) < 8 {
					return 1, "Too short"
				}
				return 4, "Strong"
			},
		})
	})
	formTextClick(t, tester, "Password")
	tester.Type("abc")
	if !tester.HasText("Too short") {
		t.Fatal("the hint for a short password is not shown")
	}
	tester.Type("defghij")
	if !tester.HasText("Strong") || tester.HasText("Too short") {
		t.Fatal("the hint did not follow the value")
	}
}
