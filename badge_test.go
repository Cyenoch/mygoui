package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	myui "github.com/egoist/mygo/ui"
)

func TestBadgeOnlyActionsParticipateInFocusOrder(t *testing.T) {
	calls := 0
	disabled := false
	tt := myui.NewTester(func(c *myui.Context) {
		myui.Row(c).Children(func() {
			mygoui.Badge(c, mygoui.BadgeOptions{Label: "Status"})
			mygoui.Badge(c, mygoui.BadgeOptions{
				Label: "Details", Variant: mygoui.BadgeOutline, Disabled: disabled,
				OnClick: func() { calls++ },
			})
		})
	}, 320, 120)
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("Details") || tt.Focused("Status") {
		t.Fatal("plain badge incorrectly became an action target")
	}
	tt.Key(0, myui.KeyEnter)
	if calls != 1 {
		t.Fatalf("badge activation count=%d", calls)
	}
	disabled = true
	tt.Frame()
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatal("disabled badge action activated")
	}
}
