package form

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestCheckboxCardTogglesFromAnywhereOnTheCard(t *testing.T) {
	checked := false
	var reported []bool
	tester := formTextTester(t, func(c *ui.Context) {
		CheckboxCard(c, CheckboxCardOptions{
			Checked: &checked, Title: "Backups", Description: "Nightly snapshots",
			OnChange: func(on bool) { reported = append(reported, on) },
		})
	})
	info, ok := tester.AccessByRole(ui.RoleCheckBox, "Backups")
	if !ok || info.Checked || info.Description != "Nightly snapshots" {
		t.Fatalf("card access = %+v (found %v)", info, ok)
	}
	// The description is far from the control, yet part of the check box.
	formTextClick(t, tester, "Nightly snapshots")
	if !checked {
		t.Fatal("a click on the card's description did not check it")
	}
	tester.Key(0, ui.KeySpace)
	if checked || len(reported) != 2 || !reported[0] || reported[1] {
		t.Fatalf("Space left checked=%v after reports %v", checked, reported)
	}
	if info, _ := tester.AccessByRole(ui.RoleCheckBox, "Backups"); info.Checked {
		t.Fatal("assistive technology still reads the card as checked")
	}
}

func TestCheckboxCardDisabledKeepsValue(t *testing.T) {
	checked := true
	tester := formTextTester(t, func(c *ui.Context) {
		CheckboxCard(c, CheckboxCardOptions{Checked: &checked, Title: "Backups", Disabled: true})
	})
	_ = tester.Click("Backups")
	tester.Frame()
	if !checked {
		t.Fatal("a disabled card changed its value")
	}
}

func TestRadioCardsShareOneChoiceInRadioGroup(t *testing.T) {
	plan := "free"
	var reported []string
	state := &RadioGroupState[string]{Selected: &plan, OnChange: func(v string) { reported = append(reported, v) }}
	tester := formTextTester(t, func(c *ui.Context) {
		RadioGroup(c, RadioGroupOptions{Label: "Plan"}, func() {
			for _, v := range []string{"free", "team", "scale"} {
				RadioCard(c, RadioCardOptions[string]{State: state, Value: v, Title: v})
			}
		})
	})
	formTextClick(t, tester, "team")
	if plan != "team" {
		t.Fatalf("a click selected %q, want team", plan)
	}
	tester.Key(0, ui.KeyDown)
	if plan != "scale" {
		t.Fatalf("Down selected %q, want scale", plan)
	}
	if len(reported) != 2 || reported[0] != "team" || reported[1] != "scale" {
		t.Fatalf("reports = %q", reported)
	}
	for v, want := range map[string]bool{"free": false, "team": false, "scale": true} {
		info, ok := tester.AccessByRole(ui.RoleRadio, v)
		if !ok || info.Checked != want {
			t.Fatalf("radio %q: found %v, checked %v, want %v", v, ok, info.Checked, want)
		}
	}
}

func TestChoiceCardContentShowsAndKeepsTheCardOneControl(t *testing.T) {
	checked := false
	tester := formTextTester(t, func(c *ui.Context) {
		CheckboxCard(c, CheckboxCardOptions{
			Checked: &checked, Title: "Backups", HideIndicator: true,
			Content: func() { ui.Text(c, "$9 / month") },
		})
	})
	if !tester.HasText("$9 / month") {
		t.Fatal("the card content is not shown")
	}
	formTextClick(t, tester, "$9 / month")
	if !checked {
		t.Fatal("a click on the card's content did not check it")
	}
	if info, ok := tester.AccessByRole(ui.RoleCheckBox, "Backups"); !ok || !info.Checked {
		t.Fatalf("the card without an indicator reads as %+v (found %v)", info, ok)
	}
}

func TestChoiceCardsExposeInvalid(t *testing.T) {
	on := false
	state := &RadioGroupState[string]{Selected: new(string)}
	tt := formTextTester(t, func(c *ui.Context) {
		CheckboxCard(c, CheckboxCardOptions{Checked: &on, Title: "Terms", Invalid: true})
		RadioGroup(c, RadioGroupOptions{Label: "Plan"}, func() {
			RadioCard(c, RadioCardOptions[string]{State: state, Value: "pro", Title: "Pro", Invalid: true})
		})
	})
	if info, ok := tt.AccessByRole(ui.RoleCheckBox, "Terms"); !ok || !info.Invalid {
		t.Fatalf("check box card: found %v, invalid %v", ok, info.Invalid)
	}
	if info, ok := tt.AccessByRole(ui.RoleRadio, "Pro"); !ok || !info.Invalid {
		t.Fatalf("radio card: found %v, invalid %v", ok, info.Invalid)
	}
}
