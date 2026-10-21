package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestAlertActionUsesNativeKeyboardAndDisabledStateAtLargeText(t *testing.T) {
	calls, disabled := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Alert(c, mygoui.AlertOptions{Title: "Transfer failed", Description: "The file is still available. Try transferring it again.", Variant: mygoui.AlertDestructive,
			Action: &mygoui.ButtonOptions{Label: "Retry transfer", Disabled: disabled, OnClick: func() { calls++ }},
		}).Width(280)
	}, 300, 300)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if calls != 1 {
		t.Fatal("callout action is not a keyboard action")
	}
	disabled = true
	tt.Frame()
	if err := tt.Click("Retry transfer"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatal("disabled callout action ran")
	}
}
