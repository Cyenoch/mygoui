package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	myui "github.com/egoist/mygo/ui"
)

func TestButtonGroupRetainsIndividualTabStopsAndActions(t *testing.T) {
	back, forward := 0, 0
	tt := myui.NewTester(func(c *myui.Context) {
		mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{Label: "History"},
			mygoui.ButtonOptions{Label: "Back", OnClick: func() { back++ }},
			mygoui.ButtonOptions{Label: "Forward", OnClick: func() { forward++ }},
		)
	}, 320, 160)
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("Back") {
		t.Fatal("first action not focusable")
	}
	tt.Key(0, myui.KeyEnter)
	tt.Key(0, myui.KeyTab)
	if !tt.Focused("Forward") {
		t.Fatal("action group incorrectly used roving toggle focus")
	}
	tt.Key(0, myui.KeySpace)
	if back != 1 || forward != 1 {
		t.Fatalf("actions back=%d forward=%d", back, forward)
	}
}
